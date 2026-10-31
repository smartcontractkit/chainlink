package agentcommerce

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type TrustedAttestationIdentity struct {
	Issuer string
	KeyID  string
}

type VerificationPolicy struct {
	Version      string
	Domain       string
	Network      string
	WorkflowID   string
	CapabilityID string

	TrustedIdentities []TrustedAttestationIdentity
	MaxFutureSkew     time.Duration
	MaxAge            time.Duration
	MaxLifetime       time.Duration
}

func (p VerificationPolicy) ValidateConfiguration() error {
	required := map[string]string{
		"version": p.Version, "domain": p.Domain, "network": p.Network,
		"workflow id": p.WorkflowID, "capability id": p.CapabilityID,
	}
	for name, value := range required {
		if err := requireCleanValue(name, value); err != nil {
			return err
		}
	}
	if len(p.TrustedIdentities) == 0 {
		return errors.New("at least one trusted attestation identity is required")
	}
	for i, identity := range p.TrustedIdentities {
		if err := requireCleanValue(fmt.Sprintf("trusted identity %d issuer", i), identity.Issuer); err != nil {
			return err
		}
		if err := requireCleanValue(fmt.Sprintf("trusted identity %d key id", i), identity.KeyID); err != nil {
			return err
		}
	}
	if p.MaxFutureSkew < 0 {
		return errors.New("max future skew must not be negative")
	}
	if p.MaxAge <= 0 {
		return errors.New("max age must be positive")
	}
	if p.MaxLifetime <= 0 {
		return errors.New("max lifetime must be positive")
	}
	return nil
}

func (p VerificationPolicy) ValidateAttestation(a Attestation, now time.Time) error {
	if err := p.ValidateConfiguration(); err != nil {
		return fmt.Errorf("invalid verification policy: %w", err)
	}
	if err := validateAttestationShape(a, true); err != nil {
		return err
	}
	if now.IsZero() {
		return errors.New("verification time is required")
	}
	now = now.UTC()
	if a.Version != p.Version {
		return errors.New("unexpected attestation version")
	}
	if a.Domain != p.Domain {
		return errors.New("unexpected attestation domain")
	}
	if a.Network != p.Network {
		return errors.New("unexpected attestation network")
	}
	if a.WorkflowID != p.WorkflowID {
		return errors.New("unexpected attestation workflow")
	}
	if a.CapabilityID != p.CapabilityID {
		return errors.New("unexpected attestation capability")
	}
	if !p.trusts(a.Issuer, a.KeyID) {
		return errors.New("untrusted attestation identity")
	}
	emptySum := sha256.Sum256(nil)
	emptyHash := hex.EncodeToString(emptySum[:])
	if a.OutputHash == emptyHash {
		return errors.New("attestation binds an empty output")
	}
	if a.ReportHash == emptyHash {
		return errors.New("attestation binds an empty report")
	}
	if a.IssuedAt.After(now.Add(p.MaxFutureSkew)) {
		return errors.New("attestation issued_at is too far in the future")
	}
	if !a.ExpiresAt.After(now) {
		return errors.New("attestation has expired")
	}
	if now.Sub(a.IssuedAt) > p.MaxAge {
		return errors.New("attestation is stale")
	}
	if a.ExpiresAt.Sub(a.IssuedAt) > p.MaxLifetime {
		return errors.New("attestation lifetime exceeds policy")
	}
	return nil
}

func (p VerificationPolicy) trusts(issuer, keyID string) bool {
	for _, identity := range p.TrustedIdentities {
		if identity.Issuer == issuer && identity.KeyID == keyID {
			return true
		}
	}
	return false
}
