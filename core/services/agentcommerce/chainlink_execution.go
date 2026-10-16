package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type ChainlinkExecutionRequest struct {
	IntentHash, ExecutionID, Service, Buyer, Seller string
	Network, WorkflowID, CapabilityID, Chain        string
	Deliverables                                    []string
	Metadata                                        map[string]string
	RequestedAt                                     string
}

type ChainlinkExecutionReport struct {
	Output      []byte
	Report      []byte
	Attestation Attestation
	Metadata    map[string]string
}

type ChainlinkExecutorBackend interface {
	Execute(context.Context, ChainlinkExecutionRequest) (ChainlinkExecutionReport, error)
}

type ChainlinkExecutorConfig struct {
	Backend      ChainlinkExecutorBackend
	Clock        Clock
	Network      string
	WorkflowID   string
	CapabilityID string
}

type ChainlinkExecutor struct {
	backend      ChainlinkExecutorBackend
	clock        Clock
	network      string
	workflowID   string
	capabilityID string
}

func NewChainlinkExecutor(cfg ChainlinkExecutorConfig) (*ChainlinkExecutor, error) {
	if cfg.Backend == nil {
		return nil, errors.New("chainlink execution backend is required")
	}
	if cfg.Clock == nil {
		return nil, errors.New("clock is required")
	}
	if err := requireCleanValue("network", cfg.Network); err != nil {
		return nil, err
	}
	if err := requireCleanValue("workflow id", cfg.WorkflowID); err != nil {
		return nil, err
	}
	if err := requireCleanValue("capability id", cfg.CapabilityID); err != nil {
		return nil, err
	}
	return &ChainlinkExecutor{cfg.Backend, cfg.Clock, cfg.Network, cfg.WorkflowID, cfg.CapabilityID}, nil
}

func (e *ChainlinkExecutor) Execute(ctx context.Context, intent SignedIntent) ([]byte, map[string]string, error) {
	if ctx == nil {
		return nil, nil, errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if err := ValidateIntentTerms(intent.Terms); err != nil {
		return nil, nil, fmt.Errorf("invalid intent: %w", err)
	}
	if err := ValidateSHA256Hex("intent hash", intent.Hash); err != nil {
		return nil, nil, err
	}
	calculated, err := IntentHash(intent.Terms)
	if err != nil {
		return nil, nil, fmt.Errorf("calculate intent hash: %w", err)
	}
	if calculated != intent.Hash {
		return nil, nil, errors.New("intent hash mismatch")
	}
	executionID, err := ExecutionIDForIntent(intent.Hash)
	if err != nil {
		return nil, nil, err
	}
	now := e.clock.Now().UTC()
	if now.IsZero() {
		return nil, nil, errors.New("clock returned zero time")
	}
	request := ChainlinkExecutionRequest{
		IntentHash: intent.Hash, ExecutionID: executionID,
		Service: intent.Terms.Service, Buyer: intent.Terms.Buyer, Seller: intent.Terms.Seller,
		Network: e.network, WorkflowID: e.workflowID, CapabilityID: e.capabilityID,
		Chain:        intent.Terms.SettlementChain,
		Deliverables: append([]string(nil), intent.Terms.Deliverables...),
		Metadata:     cloneStringMap(intent.Terms.Metadata), RequestedAt: canonicalTimestamp(now),
	}
	result, err := e.backend.Execute(ctx, request)
	if err != nil {
		return nil, nil, fmt.Errorf("chainlink execution failed: %w", err)
	}
	if len(result.Output) == 0 {
		return nil, nil, errors.New("chainlink execution returned empty output")
	}
	if len(result.Report) == 0 {
		return nil, nil, errors.New("chainlink execution returned empty report")
	}
	outputSum, reportSum := sha256.Sum256(result.Output), sha256.Sum256(result.Report)
	outputHash, reportHash := hex.EncodeToString(outputSum[:]), hex.EncodeToString(reportSum[:])
	a := result.Attestation
	if err := validateAttestationShape(a, true); err != nil {
		return nil, nil, fmt.Errorf("invalid chainlink attestation: %w", err)
	}
	checks := []struct{ got, want, message string }{
		{a.Version, ChainlinkAttestationVersion, "unexpected attestation version"},
		{a.Domain, ChainlinkAttestationDomain, "unexpected attestation domain"},
		{a.IntentHash, intent.Hash, "attestation intent hash mismatch"},
		{a.ExecutionID, executionID, "attestation execution id mismatch"},
		{a.OutputHash, outputHash, "attestation output hash mismatch"},
		{a.ReportHash, reportHash, "attestation report hash mismatch"},
		{a.Network, e.network, "attestation network mismatch"},
		{a.WorkflowID, e.workflowID, "attestation workflow mismatch"},
		{a.CapabilityID, e.capabilityID, "attestation capability mismatch"},
	}
	for _, check := range checks {
		if check.got != check.want {
			return nil, nil, errors.New("chainlink execution returned " + check.message)
		}
	}
	metadata := cloneStringMap(result.Metadata)
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata[metadataAttestationVersion] = a.Version
	metadata[metadataAttestationDomain] = a.Domain
	metadata[metadataIntentHash] = a.IntentHash
	metadata[metadataExecutionID] = a.ExecutionID
	metadata[metadataOutputSHA256] = a.OutputHash
	metadata[metadataReportSHA256] = a.ReportHash
	metadata[metadataNetwork] = a.Network
	metadata[metadataWorkflowID] = a.WorkflowID
	metadata[metadataCapabilityID] = a.CapabilityID
	metadata[metadataIssuer] = a.Issuer
	metadata[metadataKeyID] = a.KeyID
	metadata[metadataIssuedAt] = canonicalTimestamp(a.IssuedAt)
	metadata[metadataExpiresAt] = canonicalTimestamp(a.ExpiresAt)
	metadata[metadataAttestationSignature] = hex.EncodeToString(a.Signature)
	return append([]byte(nil), result.Output...), metadata, nil
}

func ExecutionIDForIntent(intentHash string) (string, error) {
	if err := ValidateSHA256Hex("intent hash", intentHash); err != nil {
		return "", err
	}
	return "exec-" + intentHash[:16], nil
}

func cloneStringMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}
