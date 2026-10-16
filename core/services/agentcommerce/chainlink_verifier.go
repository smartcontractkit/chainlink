package agentcommerce

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidAttestationSignature = errors.New("invalid attestation signature")
	ErrUnknownAttestationKey       = errors.New("unknown attestation key")
)

type AttestationSignatureVerifier interface {
	VerifySignature(context.Context, AttestationKey, []byte, []byte) error
}

type Ed25519AttestationVerifier struct {
	keys map[AttestationKey]ed25519.PublicKey
}

func NewEd25519AttestationVerifier(keys map[AttestationKey]ed25519.PublicKey) (*Ed25519AttestationVerifier, error) {
	if len(keys) == 0 {
		return nil, errors.New("at least one attestation public key is required")
	}
	copied := make(map[AttestationKey]ed25519.PublicKey, len(keys))
	for key, publicKey := range keys {
		if err := requireCleanValue("attestation issuer", key.Issuer); err != nil {
			return nil, err
		}
		if err := requireCleanValue("attestation key id", key.KeyID); err != nil {
			return nil, err
		}
		if len(publicKey) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("invalid public key for issuer %q key %q", key.Issuer, key.KeyID)
		}
		copied[key] = append(ed25519.PublicKey(nil), publicKey...)
	}
	return &Ed25519AttestationVerifier{keys: copied}, nil
}

func (v *Ed25519AttestationVerifier) VerifySignature(ctx context.Context, key AttestationKey, payload, signature []byte) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	publicKey, exists := v.keys[key]
	if !exists {
		return ErrUnknownAttestationKey
	}
	if len(payload) == 0 {
		return errors.New("attestation payload is required")
	}
	if len(signature) != ed25519.SignatureSize || !ed25519.Verify(publicKey, payload, signature) {
		return ErrInvalidAttestationSignature
	}
	return nil
}

type ChainlinkProofVerifierConfig struct {
	SignatureVerifier AttestationSignatureVerifier
	Policy            VerificationPolicy
	ReplayGuard       ReplayGuard
	Clock             Clock
}

type ChainlinkProofVerifier struct {
	signatureVerifier AttestationSignatureVerifier
	policy            VerificationPolicy
	replayGuard       ReplayGuard
	clock             Clock
}

func NewChainlinkProofVerifier(cfg ChainlinkProofVerifierConfig) (*ChainlinkProofVerifier, error) {
	if cfg.SignatureVerifier == nil {
		return nil, errors.New("attestation signature verifier is required")
	}
	if cfg.ReplayGuard == nil {
		return nil, errors.New("replay guard is required")
	}
	if cfg.Clock == nil {
		return nil, errors.New("clock is required")
	}
	if err := cfg.Policy.ValidateConfiguration(); err != nil {
		return nil, err
	}
	return &ChainlinkProofVerifier{cfg.SignatureVerifier, cfg.Policy, cfg.ReplayGuard, cfg.Clock}, nil
}

func (v *ChainlinkProofVerifier) Verify(ctx context.Context, intent SignedIntent, proof Proof) (VerificationResult, error) {
	if ctx == nil {
		return VerificationResult{}, errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return VerificationResult{}, err
	}
	now := v.clock.Now().UTC()
	if now.IsZero() {
		return VerificationResult{}, errors.New("clock returned zero time")
	}
	if err := ValidateIntentTerms(intent.Terms); err != nil {
		return VerificationResult{}, fmt.Errorf("invalid intent: %w", err)
	}
	if err := ValidateSHA256Hex("intent hash", intent.Hash); err != nil {
		return verificationRejected(now, "invalid intent hash"), nil
	}
	calculated, err := IntentHash(intent.Terms)
	if err != nil {
		return VerificationResult{}, fmt.Errorf("calculate intent hash: %w", err)
	}
	if calculated != intent.Hash {
		return verificationRejected(now, "intent hash mismatch"), nil
	}
	if proof.ExecutionID == "" {
		return verificationRejected(now, "missing execution id"), nil
	}
	if proof.Timestamp.IsZero() {
		return verificationRejected(now, "missing proof timestamp"), nil
	}
	if proof.Timestamp.After(now.Add(v.policy.MaxFutureSkew)) {
		return verificationRejected(now, "proof timestamp is too far in the future"), nil
	}
	if err := ValidateSHA256Hex("output hash", proof.OutputHash); err != nil {
		return verificationRejected(now, "invalid output hash"), nil
	}
	a, err := attestationFromProof(proof)
	if err != nil {
		return verificationRejected(now, "invalid attestation evidence"), nil
	}
	if a.IntentHash != intent.Hash {
		return verificationRejected(now, "attestation intent mismatch"), nil
	}
	if a.ExecutionID != proof.ExecutionID {
		return verificationRejected(now, "attestation execution mismatch"), nil
	}
	if a.OutputHash != proof.OutputHash {
		return verificationRejected(now, "attestation output mismatch"), nil
	}
	if err := v.policy.ValidateAttestation(a, now); err != nil {
		return verificationRejected(now, "attestation rejected by policy"), nil
	}
	payload, err := CanonicalAttestationPayload(a)
	if err != nil {
		return verificationRejected(now, "invalid attestation payload"), nil
	}
	err = v.signatureVerifier.VerifySignature(ctx, AttestationKey{a.Issuer, a.KeyID}, payload, a.Signature)
	if err != nil {
		if errors.Is(err, ErrInvalidAttestationSignature) || errors.Is(err, ErrUnknownAttestationKey) {
			return verificationRejected(now, "attestation signature rejected"), nil
		}
		return VerificationResult{}, fmt.Errorf("attestation signature verification failed: %w", err)
	}
	replayKey, err := AttestationReplayKey(a)
	if err != nil {
		return VerificationResult{}, fmt.Errorf("create replay key: %w", err)
	}
	if err := v.replayGuard.Claim(ctx, replayKey); err != nil {
		if errors.Is(err, ErrAttestationReplay) {
			return verificationRejected(now, "attestation replay rejected"), nil
		}
		return VerificationResult{}, fmt.Errorf("claim replay key: %w", err)
	}
	return VerificationResult{Verified: true, Method: "chainlink-attestation-v1", Time: now}, nil
}

func attestationFromProof(proof Proof) (Attestation, error) {
	if proof.Metadata == nil {
		return Attestation{}, errors.New("proof metadata is required")
	}
	issuedAt, err := parseCanonicalTimestamp("attestation issued_at", proof.Metadata[metadataIssuedAt])
	if err != nil {
		return Attestation{}, err
	}
	expiresAt, err := parseCanonicalTimestamp("attestation expires_at", proof.Metadata[metadataExpiresAt])
	if err != nil {
		return Attestation{}, err
	}
	signature, err := hex.DecodeString(proof.Metadata[metadataAttestationSignature])
	if err != nil || len(signature) == 0 {
		return Attestation{}, errors.New("invalid attestation signature encoding")
	}
	a := Attestation{
		Version: proof.Metadata[metadataAttestationVersion], Domain: proof.Metadata[metadataAttestationDomain],
		IntentHash: proof.Metadata[metadataIntentHash], ExecutionID: proof.Metadata[metadataExecutionID],
		OutputHash: proof.Metadata[metadataOutputSHA256], ReportHash: proof.Metadata[metadataReportSHA256],
		Network: proof.Metadata[metadataNetwork], WorkflowID: proof.Metadata[metadataWorkflowID],
		CapabilityID: proof.Metadata[metadataCapabilityID], Issuer: proof.Metadata[metadataIssuer],
		KeyID: proof.Metadata[metadataKeyID], IssuedAt: issuedAt, ExpiresAt: expiresAt, Signature: signature,
	}
	if err := validateAttestationShape(a, true); err != nil {
		return Attestation{}, err
	}
	return a, nil
}

func verificationRejected(now time.Time, reason string) VerificationResult {
	return VerificationResult{Verified: false, Method: "chainlink-attestation-v1", Reason: reason, Time: now}
}
