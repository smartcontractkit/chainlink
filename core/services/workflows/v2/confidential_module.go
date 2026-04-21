package v2

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	confworkflowtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
)

const confidentialWorkflowsCapabilityID = "confidential-workflows@1.0.0-alpha"

// WorkflowAttributes is the JSON structure stored in WorkflowSpec.Attributes.
type WorkflowAttributes struct {
	Confidential    bool               `json:"confidential"`
	VaultDonSecrets []SecretIdentifier `json:"vault_don_secrets"`
}

// SecretIdentifier identifies a secret in VaultDON.
type SecretIdentifier struct {
	Key       string `json:"key"`
	Namespace string `json:"namespace,omitempty"`
}

// ParseWorkflowAttributes parses the Attributes JSON from a WorkflowSpec.
// Returns a zero-value struct if data is nil or empty.
func ParseWorkflowAttributes(data []byte) (WorkflowAttributes, error) {
	var attrs WorkflowAttributes
	if len(data) == 0 {
		return attrs, nil
	}
	if err := json.Unmarshal(data, &attrs); err != nil {
		return attrs, fmt.Errorf("failed to parse workflow attributes: %w", err)
	}
	return attrs, nil
}

// IsConfidential returns true if the Attributes JSON has "confidential": true.
// Returns an error if the attributes contain malformed JSON, so callers can
// fail loudly rather than silently falling through to non-confidential execution.
func IsConfidential(data []byte) (bool, error) {
	attrs, err := ParseWorkflowAttributes(data)
	if err != nil {
		return false, err
	}
	return attrs.Confidential, nil
}

// ConfidentialModule implements host.ModuleV2 for confidential workflows.
// Instead of running WASM locally, it delegates execution to the
// confidential-workflows capability via the CapabilitiesRegistry.
type ConfidentialModule struct {
	capRegistry     core.CapabilitiesRegistry
	binaryURL       string
	binaryHash      []byte
	workflowID      string
	workflowOwner   string
	workflowName    string
	workflowTag     string
	vaultDonSecrets []SecretIdentifier
	lggr            logger.Logger
}

var _ host.ModuleV2 = (*ConfidentialModule)(nil)

func NewConfidentialModule(
	capRegistry core.CapabilitiesRegistry,
	binaryURL string,
	binaryHash []byte,
	workflowID, workflowOwner, workflowName, workflowTag string,
	vaultDonSecrets []SecretIdentifier,
	lggr logger.Logger,
) *ConfidentialModule {
	return &ConfidentialModule{
		capRegistry:     capRegistry,
		binaryURL:       binaryURL,
		binaryHash:      binaryHash,
		workflowID:      workflowID,
		workflowOwner:   workflowOwner,
		workflowName:    workflowName,
		workflowTag:     workflowTag,
		vaultDonSecrets: vaultDonSecrets,
		lggr:            lggr,
	}
}

func (m *ConfidentialModule) Start()            {}
func (m *ConfidentialModule) Close()            {}
func (m *ConfidentialModule) IsLegacyDAG() bool { return false }

func (m *ConfidentialModule) Execute(
	ctx context.Context,
	request *sdkpb.ExecuteRequest,
	helper host.ExecutionHelper,
) (*sdkpb.ExecutionResult, error) {
	execReqBytes, err := proto.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ExecuteRequest: %w", err)
	}

	protoSecrets := make([]*confworkflowtypes.SecretIdentifier, len(m.vaultDonSecrets))
	for i, s := range m.vaultDonSecrets {
		// VaultDON treats "main" as the default namespace for secrets.
		ns := s.Namespace
		if ns == "" {
			ns = "main"
		}
		protoSecrets[i] = &confworkflowtypes.SecretIdentifier{
			Key:       s.Key,
			Namespace: &ns,
		}
	}

	capInput := &confworkflowtypes.ConfidentialWorkflowRequest{
		VaultDonSecrets: protoSecrets,
		Execution: &confworkflowtypes.WorkflowExecution{
			WorkflowId:     m.workflowID,
			BinaryUrl:      m.binaryURL,
			BinaryHash:     m.binaryHash,
			ExecuteRequest: execReqBytes,
			Owner:          m.workflowOwner,
			ExecutionId:    helper.GetWorkflowExecutionID(),
			OrgId:          contexts.CREValue(ctx).Org,
		},
	}

	payload, err := anypb.New(capInput)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability payload: %w", err)
	}

	executable, err := m.capRegistry.GetExecutable(ctx, confidentialWorkflowsCapabilityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get confidential-workflows capability: %w", err)
	}

	capReq := capabilities.CapabilityRequest{
		Payload:      payload,
		Method:       "Execute",
		CapabilityId: confidentialWorkflowsCapabilityID,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          m.workflowID,
			WorkflowOwner:       m.workflowOwner,
			WorkflowName:        m.workflowName,
			WorkflowTag:         m.workflowTag,
			WorkflowExecutionID: helper.GetWorkflowExecutionID(),
		},
	}

	capResp, err := executable.Execute(ctx, capReq)
	if err != nil {
		return nil, fmt.Errorf("confidential-workflows capability execution failed: %w", err)
	}

	if capResp.Payload == nil {
		return nil, errors.New("confidential-workflows capability returned nil payload")
	}

	var confResp confworkflowtypes.ConfidentialWorkflowResponse
	if err := capResp.Payload.UnmarshalTo(&confResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal capability response: %w", err)
	}

	var result sdkpb.ExecutionResult
	if err := proto.Unmarshal(confResp.ExecutionResult, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ExecutionResult: %w", err)
	}

	return &result, nil
}

// ComputeBinaryHash returns the SHA-256 hash of the given binary.
func ComputeBinaryHash(binary []byte) []byte {
	h := sha256.Sum256(binary)
	return h[:]
}
