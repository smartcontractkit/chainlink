package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ocr2key"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

const nativeEvidenceVersion = "chainlink-ocr-capability-v1"

type NativeOCRPolicy struct {
	Signers       []ocrtypes.OnchainPublicKey
	MinSignatures int
	Issuer        string
	MaxLifetime   time.Duration
}

type nativeOCRSignature struct {
	Signer    uint32 `json:"signer"`
	Signature string `json:"signature"`
}

type nativeOCREvidence struct {
	Version          string                     `json:"version"`
	IntentHash       string                     `json:"intent_hash"`
	ExecutionID      string                     `json:"execution_id"`
	Network          string                     `json:"network"`
	WorkflowID       string                     `json:"workflow_id"`
	CapabilityID     string                     `json:"capability_id"`
	ReferenceID      string                     `json:"reference_id"`
	Output           string                     `json:"output"`
	ConfigDigest     string                     `json:"config_digest"`
	SequenceNumber   uint64                     `json:"sequence_number"`
	Signatures       []nativeOCRSignature       `json:"signatures"`
	ResponseMetadata commoncap.ResponseMetadata `json:"response_metadata"`
}

type NativeCapabilityBackendConfig struct {
	Executable    commoncap.Executable
	Policy        NativeOCRPolicy
	Network       string
	WorkflowID    string
	WorkflowOwner string
	CapabilityID  string
	Method        string
}

type NativeCapabilityBackend struct{ config NativeCapabilityBackendConfig }

var _ ChainlinkExecutorBackend = (*NativeCapabilityBackend)(nil)

func NewNativeCapabilityBackend(cfg NativeCapabilityBackendConfig) (*NativeCapabilityBackend, error) {
	if cfg.Executable == nil {
		return nil, errors.New("native Chainlink executable is required")
	}
	for name, value := range map[string]string{
		"network": cfg.Network, "workflow id": cfg.WorkflowID, "workflow owner": cfg.WorkflowOwner,
		"capability id": cfg.CapabilityID, "capability method": cfg.Method, "OCR issuer": cfg.Policy.Issuer,
	} {
		if err := requireCleanValue(name, value); err != nil {
			return nil, err
		}
	}
	if cfg.Policy.MinSignatures <= 0 || cfg.Policy.MinSignatures > len(cfg.Policy.Signers) {
		return nil, errors.New("invalid native OCR signature threshold")
	}
	if cfg.Policy.MaxLifetime <= 0 {
		return nil, errors.New("native evidence lifetime must be positive")
	}
	for i, signer := range cfg.Policy.Signers {
		if len(signer) == 0 {
			return nil, fmt.Errorf("native OCR signer %d is empty", i)
		}
		cfg.Policy.Signers[i] = append(ocrtypes.OnchainPublicKey(nil), signer...)
	}
	return &NativeCapabilityBackend{config: cfg}, nil
}

func (b *NativeCapabilityBackend) Execute(ctx context.Context, req ChainlinkExecutionRequest) (ChainlinkExecutionReport, error) {
	if ctx == nil {
		return ChainlinkExecutionReport{}, errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return ChainlinkExecutionReport{}, err
	}
	wantExecutionID, err := ExecutionIDForIntent(req.IntentHash)
	if err != nil {
		return ChainlinkExecutionReport{}, err
	}
	if req.ExecutionID != wantExecutionID {
		return ChainlinkExecutionReport{}, errors.New("execution id mismatch")
	}
	if req.Network != b.config.Network || req.WorkflowID != b.config.WorkflowID || req.CapabilityID != b.config.CapabilityID {
		return ChainlinkExecutionReport{}, errors.New("native capability authority mismatch")
	}
	requestedAt, err := parseCanonicalTimestamp("requested_at", req.RequestedAt)
	if err != nil {
		return ChainlinkExecutionReport{}, err
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return ChainlinkExecutionReport{}, err
	}
	nativeReq := commoncap.CapabilityRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: b.config.WorkflowID, WorkflowOwner: b.config.WorkflowOwner,
			WorkflowExecutionID: req.ExecutionID, ReferenceID: req.IntentHash,
			ExecutionTimestamp: requestedAt,
		},
		Payload: &anypb.Any{TypeUrl: "type.googleapis.com/agentcommerce.v1.ExecutionRequest", Value: payload},
		Method:  b.config.Method, CapabilityId: b.config.CapabilityID,
	}
	resp, err := b.config.Executable.Execute(ctx, nativeReq)
	if err != nil {
		return ChainlinkExecutionReport{}, fmt.Errorf("execute native Chainlink capability: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return ChainlinkExecutionReport{}, err
	}
	if resp.Payload == nil || len(resp.Payload.Value) == 0 {
		return ChainlinkExecutionReport{}, errors.New("native capability returned empty output")
	}
	if resp.OCRAttestation == nil {
		return ChainlinkExecutionReport{}, errors.New("native capability returned no OCR attestation")
	}
	evidence := nativeOCREvidence{
		Version: nativeEvidenceVersion, IntentHash: req.IntentHash, ExecutionID: req.ExecutionID,
		Network: b.config.Network, WorkflowID: b.config.WorkflowID, CapabilityID: b.config.CapabilityID,
		ReferenceID: req.IntentHash, Output: hex.EncodeToString(resp.Payload.Value),
		ConfigDigest:   hex.EncodeToString(resp.OCRAttestation.ConfigDigest[:]),
		SequenceNumber: resp.OCRAttestation.SequenceNumber, ResponseMetadata: resp.Metadata,
	}
	for _, signature := range resp.OCRAttestation.Sigs {
		evidence.Signatures = append(evidence.Signatures, nativeOCRSignature{Signer: signature.Signer, Signature: hex.EncodeToString(signature.Signature)})
	}
	if err := verifyNativeOCREvidence(evidence, b.config.Policy); err != nil {
		return ChainlinkExecutionReport{}, err
	}
	report, err := json.Marshal(evidence)
	if err != nil {
		return ChainlinkExecutionReport{}, err
	}
	outputHash, reportHash := sha256.Sum256(resp.Payload.Value), sha256.Sum256(report)
	attestation := Attestation{
		Version: ChainlinkAttestationVersion, Domain: ChainlinkAttestationDomain,
		IntentHash: req.IntentHash, ExecutionID: req.ExecutionID,
		OutputHash: hex.EncodeToString(outputHash[:]), ReportHash: hex.EncodeToString(reportHash[:]),
		Network: b.config.Network, WorkflowID: b.config.WorkflowID, CapabilityID: b.config.CapabilityID,
		Issuer: b.config.Policy.Issuer, KeyID: evidence.ConfigDigest,
		IssuedAt: requestedAt, ExpiresAt: requestedAt.Add(b.config.Policy.MaxLifetime), Signature: report,
	}
	return ChainlinkExecutionReport{Output: append([]byte(nil), resp.Payload.Value...), Report: report, Attestation: attestation}, nil
}

func verifyNativeOCREvidence(e nativeOCREvidence, policy NativeOCRPolicy) error {
	if e.Version != nativeEvidenceVersion || e.ReferenceID != e.IntentHash {
		return errors.New("invalid native OCR evidence identity")
	}
	if err := ValidateSHA256Hex("native intent hash", e.IntentHash); err != nil {
		return err
	}
	wantExecutionID, err := ExecutionIDForIntent(e.IntentHash)
	if err != nil || e.ExecutionID != wantExecutionID {
		return errors.New("native OCR execution identity mismatch")
	}
	output, err := hex.DecodeString(e.Output)
	if err != nil || len(output) == 0 {
		return errors.New("native OCR output is malformed")
	}
	digestBytes, err := hex.DecodeString(e.ConfigDigest)
	if err != nil || len(digestBytes) != len(ocrtypes.ConfigDigest{}) {
		return errors.New("native OCR config digest is malformed")
	}
	var digest ocrtypes.ConfigDigest
	copy(digest[:], digestBytes)
	reportData, err := commoncap.ResponseToReportData(e.ExecutionID, e.ReferenceID, output, e.ResponseMetadata)
	if err != nil {
		return fmt.Errorf("construct native OCR report data: %w", err)
	}
	sigData := ocr2key.ReportToSigData3(digest, e.SequenceNumber, reportData[:])
	seen := make(map[uint32]struct{}, len(e.Signatures))
	valid := 0
	for _, attributed := range e.Signatures {
		if int(attributed.Signer) >= len(policy.Signers) {
			return errors.New("native OCR signer index is out of range")
		}
		if _, duplicate := seen[attributed.Signer]; duplicate {
			return errors.New("native OCR signer is duplicated")
		}
		seen[attributed.Signer] = struct{}{}
		signature, err := hex.DecodeString(attributed.Signature)
		if err != nil || !ocr2key.EvmVerifyBlob(policy.Signers[attributed.Signer], sigData, signature) {
			return errors.New("native OCR signature verification failed")
		}
		valid++
	}
	if valid < policy.MinSignatures {
		return errors.New("native OCR signature threshold not met")
	}
	return nil
}

type NativeOCRProofVerifierConfig struct {
	Policy        NativeOCRPolicy
	ReplayGuard   ReplayGuard
	Clock         Clock
	Network       string
	WorkflowID    string
	CapabilityID  string
	MaxFutureSkew time.Duration
	MaxAge        time.Duration
}

type NativeOCRProofVerifier struct{ config NativeOCRProofVerifierConfig }

func NewNativeOCRProofVerifier(cfg NativeOCRProofVerifierConfig) (*NativeOCRProofVerifier, error) {
	if cfg.ReplayGuard == nil || cfg.Clock == nil {
		return nil, errors.New("native OCR verifier requires replay guard and clock")
	}
	if cfg.MaxFutureSkew < 0 || cfg.MaxAge <= 0 {
		return nil, errors.New("invalid native OCR verifier time policy")
	}
	for name, value := range map[string]string{"network": cfg.Network, "workflow id": cfg.WorkflowID, "capability id": cfg.CapabilityID} {
		if err := requireCleanValue(name, value); err != nil {
			return nil, err
		}
	}
	if cfg.Policy.MinSignatures <= 0 || cfg.Policy.MinSignatures > len(cfg.Policy.Signers) {
		return nil, errors.New("invalid native OCR verifier signature policy")
	}
	return &NativeOCRProofVerifier{config: cfg}, nil
}

func (v *NativeOCRProofVerifier) Verify(ctx context.Context, intent SignedIntent, proof Proof) (VerificationResult, error) {
	if ctx == nil {
		return VerificationResult{}, errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return VerificationResult{}, err
	}
	now := v.config.Clock.Now().UTC()
	calculated, err := IntentHash(intent.Terms)
	if err != nil || calculated != intent.Hash {
		return verificationRejected(now, "intent hash mismatch"), nil
	}
	if proof.Timestamp.IsZero() || proof.Timestamp.After(now.Add(v.config.MaxFutureSkew)) {
		return verificationRejected(now, "invalid proof timestamp"), nil
	}
	attestation, err := attestationFromProof(proof)
	if err != nil {
		return verificationRejected(now, "invalid native attestation envelope"), nil
	}
	if attestation.IntentHash != intent.Hash || attestation.ExecutionID != proof.ExecutionID || attestation.OutputHash != proof.OutputHash {
		return verificationRejected(now, "native attestation identity mismatch"), nil
	}
	if attestation.Version != ChainlinkAttestationVersion || attestation.Domain != ChainlinkAttestationDomain {
		return verificationRejected(now, "native attestation contract mismatch"), nil
	}
	if attestation.Network != v.config.Network || attestation.WorkflowID != v.config.WorkflowID || attestation.CapabilityID != v.config.CapabilityID || attestation.Issuer != v.config.Policy.Issuer {
		return verificationRejected(now, "native attestation authority mismatch"), nil
	}
	if attestation.IssuedAt.After(now.Add(v.config.MaxFutureSkew)) || now.Sub(attestation.IssuedAt) > v.config.MaxAge || !attestation.ExpiresAt.After(now) {
		return verificationRejected(now, "native attestation time policy rejected"), nil
	}
	reportHash := sha256.Sum256(attestation.Signature)
	if hex.EncodeToString(reportHash[:]) != attestation.ReportHash {
		return verificationRejected(now, "native report hash mismatch"), nil
	}
	var native nativeOCREvidence
	if err := json.Unmarshal(attestation.Signature, &native); err != nil {
		return verificationRejected(now, "malformed native OCR evidence"), nil
	}
	if attestation.KeyID != native.ConfigDigest {
		return verificationRejected(now, "native attestation key identity mismatch"), nil
	}
	if native.IntentHash != intent.Hash || native.ExecutionID != proof.ExecutionID || native.Network != v.config.Network || native.WorkflowID != v.config.WorkflowID || native.CapabilityID != v.config.CapabilityID {
		return verificationRejected(now, "native OCR evidence substitution"), nil
	}
	output, err := hex.DecodeString(native.Output)
	if err != nil {
		return verificationRejected(now, "malformed native OCR output"), nil
	}
	outputHash := sha256.Sum256(output)
	if hex.EncodeToString(outputHash[:]) != proof.OutputHash {
		return verificationRejected(now, "native OCR output substitution"), nil
	}
	if err := verifyNativeOCREvidence(native, v.config.Policy); err != nil {
		return verificationRejected(now, "native OCR authentication rejected"), nil
	}
	replayKey, err := AttestationReplayKey(attestation)
	if err != nil {
		return VerificationResult{}, err
	}
	if err := v.config.ReplayGuard.Claim(ctx, replayKey); err != nil {
		if errors.Is(err, ErrAttestationReplay) {
			return verificationRejected(now, "native OCR replay rejected"), nil
		}
		return VerificationResult{}, err
	}
	return VerificationResult{Verified: true, Method: nativeEvidenceVersion, Time: now}, nil
}
