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

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
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
	capRegistry       core.CapabilitiesRegistry
	binaryURL         string
	binaryHash        []byte
	workflowID        string
	workflowOwner     string
	workflowName      string
	workflowTag       string
	resolveOrgID      func(ctx context.Context, owner string) (string, error)
	lggr              logger.Logger
	requirements      sync.Map
	restritions       sync.Map
	infoOnce          sync.Once
	provider          func(tee *sdkpb.Tee) bool
	executionHandlers *confidentialrelay.ExecutionHandlers
	enabledGate       limits.GateLimiter
	creSettingsGetter settings.Getter
}

var _ host.RequirementEnforcingModule = (*ConfidentialModule)(nil)
var _ host.RestrictionAwareModule = (*ConfidentialModule)(nil)

func NewConfidentialModule(capRegistry core.CapabilitiesRegistry, executionHandlers *confidentialrelay.ExecutionHandlers, binaryURL string, binaryHash []byte, workflowID, workflowOwner, workflowName, workflowTag string, resolveOrgID func(ctx context.Context, owner string) (string, error), enabledGate limits.GateLimiter, creSettingsGetter settings.Getter, lggr logger.Logger) (*ConfidentialModule, error) {
	if enabledGate == nil {
		return nil, errors.New("enabledGate must not be nil")
	}
	if resolveOrgID == nil {
		return nil, errors.New("resolveOrgID must not be nil")
	}
	return &ConfidentialModule{
		capRegistry:       capRegistry,
		executionHandlers: executionHandlers,
		binaryURL:         binaryURL,
		binaryHash:        binaryHash,
		workflowID:        workflowID,
		workflowOwner:     workflowOwner,
		workflowName:      workflowName,
		workflowTag:       workflowTag,
		resolveOrgID:      resolveOrgID,
		enabledGate:       enabledGate,
		creSettingsGetter: creSettingsGetter,
		lggr:              lggr,
	}, nil
}

func (m *ConfidentialModule) Start()            {}
func (m *ConfidentialModule) Close()            {}
func (m *ConfidentialModule) IsLegacyDAG() bool { return false }

func (m *ConfidentialModule) Execute(
	ctx context.Context,
	request *sdkpb.ExecuteRequest,
	helper host.ExecutionHelper,
) (*sdkpb.ExecutionResult, error) {
	if err := m.enabledGate.AllowErr(ctx); err != nil {
		return nil, fmt.Errorf("confidential-workflows capability is disabled by settings: %w", err)
	}

	workflowExecutionID := helper.GetWorkflowExecutionID()
	rawSecretsHelper, ok := helper.(host.ExecutionHelperWithRawSecrets)
	if !ok {
		return nil, fmt.Errorf("%T is not a host.executionHelperWithRawSecrets and is not safe to use for confidential compute", helper)
	}
	m.executionHandlers.AddExecution(m.workflowID, workflowExecutionID, rawSecretsHelper)
	defer m.executionHandlers.RemoveExecution(m.workflowID, workflowExecutionID)

	requirements := loadAndDelete[*sdkpb.Requirements](&m.requirements, workflowExecutionID)
	restrictions := loadAndDelete[*sdkpb.Restrictions](&m.restritions, workflowExecutionID)

	orgID, orgErr := m.resolveOrgID(ctx, m.workflowOwner)
	if orgErr != nil {
		m.lggr.Warnw("failed to resolve organization ID", "error", orgErr)
	}

	capInput := &confworkflowtypes.ConfidentialWorkflowRequest{
		Execution: &confworkflowtypes.WorkflowExecution{
			WorkflowId:        m.workflowID,
			BinaryHash:        m.binaryHash,
			SdkExecuteRequest: request,
			Owner:             m.workflowOwner,
			ExecutionId:       workflowExecutionID,
			OrgId:             orgID,
			Requirements:      requirements,
			BinaryUrl:         m.binaryURL,
			Restrictions:      restrictions,
		},
	}

	capOutput := &confworkflowtypes.ConfidentialWorkflowResponse{}
	if err := doRequest(ctx, m, helper.GetWorkflowExecutionID(), "Execute", capInput, capOutput, orgID); err != nil {
		return nil, err
	}

	return capOutput.SdkExecutionResult, nil
}

func (m *ConfidentialModule) SetRequirements(executionID string, requirements *sdkpb.Requirements) {
	m.requirements.Store(executionID, requirements)
}

func (m *ConfidentialModule) SetRestrictions(executionID string, restrictions *sdkpb.Restrictions) {
	m.restritions.Store(executionID, restrictions)
}

func (m *ConfidentialModule) providedTees(ctx context.Context) []*sdkpb.TeeTypeAndRegions {
	capOutput := &confworkflowtypes.ProvidedTeesResponse{}
	// use an empty execution ID, it's not during an execution.
	if err := doRequest(ctx, m, "", "ProvidedTees", &emptypb.Empty{}, capOutput, ""); err != nil {
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
	capOutput O,
	orgID string) error {
	payload, err := anypb.New(capInput)
	if err != nil {
		return fmt.Errorf("failed to marshal capability payload: %w", err)
	}

	executable, err := m.capRegistry.GetExecutable(ctx, confidentialWorkflowsCapabilityID)
	if err != nil {
		return fmt.Errorf("failed to get confidential-workflows capability: %w", err)
	}

	config, _ := anypb.New(&emptypb.Empty{})

	metadata := capabilities.RequestMetadata{
		WorkflowID:          m.workflowID,
		WorkflowOwner:       m.workflowOwner,
		WorkflowName:        m.workflowName,
		WorkflowTag:         m.workflowTag,
		WorkflowExecutionID: execID,
		OrgID:               orgID,
	}
	// When the WorkflowDonID binding gate is enabled, the remote executable
	// capability server rejects requests whose RequestMetadata.WorkflowDonID
	// does not match the authenticated calling DON. Set it to the calling
	// (local) workflow DON ID so the request is accepted. Any failure here must
	// fail the whole call rather than silently sending a zero WorkflowDonID.
	bindingEnabled, err := cresettings.Default.RemoteExecutableWorkflowDONBindingEnabled.GetOrDefault(ctx, m.creSettingsGetter)
	if err != nil {
		return fmt.Errorf("failed to read RemoteExecutableWorkflowDONBindingEnabled setting: %w", err)
	}
	if bindingEnabled {
		localNode, lnErr := m.capRegistry.LocalNode(ctx)
		if lnErr != nil {
			return fmt.Errorf("failed to get local node for confidential-workflows request metadata: %w", lnErr)
		}
		metadata.WorkflowDonID = localNode.WorkflowDON.ID
	}

	capReq := capabilities.CapabilityRequest{
		Payload:       payload,
		ConfigPayload: config,
		Method:        method,
		CapabilityId:  confidentialWorkflowsCapabilityID,
		Metadata:      metadata,
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
