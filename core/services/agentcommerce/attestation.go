package agentcommerce

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ChainlinkAttestationVersion = "1"
	ChainlinkAttestationDomain  = "agentcommerce.chainlink.attestation"

	metadataAttestationVersion   = "chainlink_attestation_version"
	metadataAttestationDomain    = "chainlink_attestation_domain"
	metadataIntentHash           = "chainlink_intent_hash"
	metadataExecutionID          = "chainlink_execution_id"
	metadataOutputSHA256         = "chainlink_output_sha256"
	metadataReportSHA256         = "chainlink_report_sha256"
	metadataNetwork              = "chainlink_network"
	metadataWorkflowID           = "chainlink_workflow_id"
	metadataCapabilityID         = "chainlink_capability_id"
	metadataIssuer               = "chainlink_attestation_issuer"
	metadataKeyID                = "chainlink_attestation_key_id"
	metadataIssuedAt             = "chainlink_attestation_issued_at"
	metadataExpiresAt            = "chainlink_attestation_expires_at"
	metadataAttestationSignature = "chainlink_attestation_signature"
)

type Clock interface {
	Now() time.Time
}

type ClockFunc func() time.Time

func (f ClockFunc) Now() time.Time { return f() }

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }

type AttestationKey struct {
	Issuer string
	KeyID  string
}

// Attestation binds an ACP intent and execution to authenticated evidence.
// Signature is excluded from the canonical signing payload.
type Attestation struct {
	Version      string
	Domain       string
	IntentHash   string
	ExecutionID  string
	OutputHash   string
	ReportHash   string
	Network      string
	WorkflowID   string
	CapabilityID string
	Issuer       string
	KeyID        string
	IssuedAt     time.Time
	ExpiresAt    time.Time
	Signature    []byte
}

type attestationSigningPayload struct {
	Version      string `json:"version"`
	Domain       string `json:"domain"`
	IntentHash   string `json:"intent_hash"`
	ExecutionID  string `json:"execution_id"`
	OutputHash   string `json:"output_hash"`
	ReportHash   string `json:"report_hash"`
	Network      string `json:"network"`
	WorkflowID   string `json:"workflow_id"`
	CapabilityID string `json:"capability_id"`
	Issuer       string `json:"issuer"`
	KeyID        string `json:"key_id"`
	IssuedAt     string `json:"issued_at"`
	ExpiresAt    string `json:"expires_at"`
}

func CanonicalAttestationPayload(a Attestation) ([]byte, error) {
	if err := validateAttestationShape(a, false); err != nil {
		return nil, err
	}
	payload := attestationSigningPayload{
		Version: a.Version, Domain: a.Domain, IntentHash: a.IntentHash,
		ExecutionID: a.ExecutionID, OutputHash: a.OutputHash,
		ReportHash: a.ReportHash, Network: a.Network,
		WorkflowID: a.WorkflowID, CapabilityID: a.CapabilityID,
		Issuer: a.Issuer, KeyID: a.KeyID,
		IssuedAt:  canonicalTimestamp(a.IssuedAt),
		ExpiresAt: canonicalTimestamp(a.ExpiresAt),
	}
	return json.Marshal(payload)
}

func AttestationDigest(a Attestation) (string, error) {
	payload, err := CanonicalAttestationPayload(a)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func validateAttestationShape(a Attestation, requireSignature bool) error {
	required := map[string]string{
		"version": a.Version, "domain": a.Domain, "intent hash": a.IntentHash,
		"execution id": a.ExecutionID, "output hash": a.OutputHash,
		"report hash": a.ReportHash, "network": a.Network,
		"workflow id": a.WorkflowID, "capability id": a.CapabilityID,
		"issuer": a.Issuer, "key id": a.KeyID,
	}
	for name, value := range required {
		if err := requireCleanValue(name, value); err != nil {
			return err
		}
	}
	if err := ValidateSHA256Hex("intent hash", a.IntentHash); err != nil {
		return err
	}
	if err := ValidateSHA256Hex("output hash", a.OutputHash); err != nil {
		return err
	}
	if err := ValidateSHA256Hex("report hash", a.ReportHash); err != nil {
		return err
	}
	if a.IssuedAt.IsZero() {
		return errors.New("attestation issued_at is required")
	}
	if a.ExpiresAt.IsZero() {
		return errors.New("attestation expires_at is required")
	}
	if !a.ExpiresAt.After(a.IssuedAt) {
		return errors.New("attestation expires_at must be after issued_at")
	}
	if requireSignature && len(a.Signature) == 0 {
		return errors.New("attestation signature is required")
	}
	return nil
}

func ValidateSHA256Hex(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	if value != strings.ToLower(value) {
		return fmt.Errorf("%s must use lowercase hexadecimal", name)
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return fmt.Errorf("%s is not valid hexadecimal", name)
	}
	if len(decoded) != sha256.Size {
		return fmt.Errorf("%s must encode exactly %d bytes", name, sha256.Size)
	}
	return nil
}

func requireCleanValue(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	if strings.TrimSpace(value) != value {
		return fmt.Errorf("%s must not contain leading or trailing whitespace", name)
	}
	return nil
}

func canonicalTimestamp(value time.Time) string {
	return value.UTC().Round(0).Format(time.RFC3339Nano)
}

func parseCanonicalTimestamp(name, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required", name)
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s is invalid", name)
	}
	if canonicalTimestamp(parsed) != value {
		return time.Time{}, fmt.Errorf("%s is not canonical", name)
	}
	return parsed.UTC(), nil
}
