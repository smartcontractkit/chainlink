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
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"

	confworkflowtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialworkflow"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"

	workflowtypes "github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const confidentialWorkflowsCapabilityID = "confidential-workflows@1.0.0-alpha"

// WorkflowAttributes is the JSON structure stored in WorkflowSpec.Attributes.
type WorkflowAttributes struct {
	Confidential bool `json:"confidential"`
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
	capRegistry   core.CapabilitiesRegistry
	binaryHash    []byte
	workflowID    string
	workflowOwner string
	workflowName  string
	workflowTag   string
	// retrieveURL mints a fresh pre-signed CloudFront URL via the storage
	// service's NodeService.DownloadArtifact at every Execute call. Each
	// workflow DON node gets its own URL with its own signature and expiry,
	// which is why the URL is carried on ConfidentialWorkflowRequest (outside
	// the hashed WorkflowExecution) rather than inside it.
	retrieveURL       workflowtypes.LocationRetrieverFunc
	creSettingsGetter settings.Getter
	orgID             string // resolved via org resolver at construction; stable for module lifetime
	lggr              logger.Logger
}

var _ host.ModuleV2 = (*ConfidentialModule)(nil)

func NewConfidentialModule(
	capRegistry core.CapabilitiesRegistry,
	binaryHash []byte,
	workflowID, workflowOwner, workflowName, workflowTag string,
	retrieveURL workflowtypes.LocationRetrieverFunc,
	creSettingsGetter settings.Getter,
	orgID string,
	lggr logger.Logger,
) *ConfidentialModule {
	return &ConfidentialModule{
		capRegistry:       capRegistry,
		binaryHash:        binaryHash,
		workflowID:        workflowID,
		workflowOwner:     workflowOwner,
		workflowName:      workflowName,
		workflowTag:       workflowTag,
		retrieveURL:       retrieveURL,
		creSettingsGetter: creSettingsGetter,
		orgID:             orgID,
		lggr:              lggr,
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

	if m.retrieveURL == nil {
		return nil, errors.New("confidential module is missing a URL retriever; cannot fetch binary from storage service")
	}
	binaryURL, err := m.retrieveURL(ctx, &storage_service.DownloadArtifactRequest{
		Id:   m.workflowID,
		Type: storage_service.ArtifactType_ARTIFACT_TYPE_BINARY,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to mint pre-signed binary URL from storage service: %w", err)
	}

	capInput := &confworkflowtypes.ConfidentialWorkflowRequest{
		// binary_url rides at the request level, outside the hashed
		// WorkflowExecution. WorkflowExecution.binary_url is deprecated and
		// intentionally left unset: it is inside ComputeRequest.PublicData
		// (hashed for F+1 quorum), so a per-node value there would break quorum.
		BinaryUrl: binaryURL,
		Execution: &confworkflowtypes.WorkflowExecution{
			WorkflowId:     m.workflowID,
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
			WorkflowOwner:       m.workflowOwner,
			WorkflowID:          m.workflowID,
			WorkflowName:        m.workflowName,
			WorkflowTag:         m.workflowTag,
			WorkflowExecutionID: helper.GetWorkflowExecutionID(),
		},
	}
	propagateOrgIDMeta, _ := cresettings.Default.PropagateOrgIDInRequestMetadata.GetOrDefault(ctx, m.creSettingsGetter)
	if propagateOrgIDMeta && m.orgID != "" {
		capReq.Metadata.OrgID = m.orgID
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
