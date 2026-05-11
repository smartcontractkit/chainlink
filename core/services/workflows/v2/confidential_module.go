package v2

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"

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

// ConfidentialModule implements host.ModuleV2 for confidential workflows.
// Instead of running WASM locally, it delegates execution to the
// confidential-workflows capability via the CapabilitiesRegistry.
type ConfidentialModule struct {
	capRegistry   core.CapabilitiesRegistry
	binaryURL     string
	binaryHash    []byte
	workflowID    string
	workflowOwner string
	workflowName  string
	workflowTag   string
	lggr          logger.Logger
	requirements  sync.Map
	restritions   sync.Map
	infoOnce      sync.Once
	provider      func(tee *sdkpb.Tee) bool
}

var _ host.RequirementEnforcingModule = (*ConfidentialModule)(nil)
var _ host.RestrictionAwareModule = (*ConfidentialModule)(nil)

func NewConfidentialModule(capRegistry core.CapabilitiesRegistry, binaryURL string, binaryHash []byte, workflowID, workflowOwner, workflowName, workflowTag string, lggr logger.Logger) *ConfidentialModule {
	return &ConfidentialModule{
		capRegistry:   capRegistry,
		binaryURL:     binaryURL,
		binaryHash:    binaryHash,
		workflowID:    workflowID,
		workflowOwner: workflowOwner,
		workflowName:  workflowName,
		workflowTag:   workflowTag,
		lggr:          lggr,
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
	workflowExecutionID := helper.GetWorkflowExecutionID()

	requirements := loadAndDelete[*sdkpb.Requirements](&m.requirements, workflowExecutionID)
	restrictions := loadAndDelete[*sdkpb.Restrictions](&m.restritions, workflowExecutionID)

	capInput := &confworkflowtypes.ConfidentialWorkflowRequest{
		Execution: &confworkflowtypes.WorkflowExecution{
			WorkflowId:     m.workflowID,
			BinaryUrl:      m.binaryURL,
			BinaryHash:     m.binaryHash,
			ExecuteRequest: request,
			Owner:          m.workflowOwner,
			ExecutionId:    workflowExecutionID,
			OrgId:          contexts.CREValue(ctx).Org,
			Requirements:   requirements,
			Restrictions:   restrictions,
		},
	}

	capOutput := &confworkflowtypes.ConfidentialWorkflowResponse{}
	if err := doRequest(ctx, m, helper.GetWorkflowExecutionID(), "Execute", capInput, capOutput); err != nil {
		return nil, err
	}

	return capOutput.ExecutionResult, nil
}

func (m *ConfidentialModule) SetRequirements(executionId string, requirements *sdkpb.Requirements) {
	m.requirements.Store(executionId, requirements)
}

func (m *ConfidentialModule) SetRestrictions(executionId string, restrictions *sdkpb.Restrictions) {
	m.restritions.Store(executionId, restrictions)
}

func (m *ConfidentialModule) providedTees(ctx context.Context) []*sdkpb.TeeTypeAndRegions {
	capOutput := &confworkflowtypes.ProvidedTeesResponse{}
	// use an empty execution ID, it's not during an execution.
	if err := doRequest(ctx, m, "", "ProvidedTees", &emptypb.Empty{}, capOutput); err != nil {
		m.lggr.Errorf("failed to get regions from confidential-workflows capability, assuming no supported regions: %v", err)
		return []*sdkpb.TeeTypeAndRegions{}
	}

	return capOutput.Tee
}

func (m *ConfidentialModule) Tee(ctx context.Context, tee *sdkpb.Tee) bool {
	m.infoOnce.Do(func() {
		m.provider = host.NewProviderFromSelection(m.providedTees(ctx))
	})

	return m.provider(tee)
}

func loadAndDelete[T any](m *sync.Map, key string) T {
	var t T
	raw, loaded := m.LoadAndDelete(key)
	if loaded {
		t = raw.(T)
	}
	return t
}

func doRequest[I, O proto.Message](
	ctx context.Context,
	m *ConfidentialModule,
	execID string,
	method string,
	capInput I,
	capOutput O) error {
	payload, err := anypb.New(capInput)
	if err != nil {
		return fmt.Errorf("failed to marshal capability payload: %w", err)
	}

	executable, err := m.capRegistry.GetExecutable(ctx, confidentialWorkflowsCapabilityID)
	if err != nil {
		return fmt.Errorf("failed to get confidential-workflows capability: %w", err)
	}

	config, _ := anypb.New(&emptypb.Empty{})

	capReq := capabilities.CapabilityRequest{
		Payload:       payload,
		ConfigPayload: config,
		Method:        method,
		CapabilityId:  confidentialWorkflowsCapabilityID,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          m.workflowID,
			WorkflowOwner:       m.workflowOwner,
			WorkflowName:        m.workflowName,
			WorkflowTag:         m.workflowTag,
			WorkflowExecutionID: execID,
		},
	}

	capResp, err := executable.Execute(ctx, capReq)
	if err != nil {
		return fmt.Errorf("confidential-workflows capability execution failed: %w", err)
	}

	if capResp.Payload == nil {
		return errors.New("confidential-workflows capability returned nil payload")
	}

	if err = capResp.Payload.UnmarshalTo(capOutput); err != nil {
		return fmt.Errorf("failed to unmarshal capability response: %w", err)
	}

	return nil
}

// ComputeBinaryHash returns the SHA-256 hash of the given binary.
func ComputeBinaryHash(binary []byte) []byte {
	h := sha256.Sum256(binary)
	return h[:]
}
