package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
)

const (
	entityCapabilityExecutionStartedV2  = "workflows.v2.CapabilityExecutionStarted"
	entityCapabilityExecutionFinishedV2 = "workflows.v2.CapabilityExecutionFinished"
)

// stubSecretsFetcher is a minimal SecretsFetcher/RawSecretsFetcher for exercising
// ExecutionHelper.GetSecrets and executionHelperWithRawSecrets.GetRawSecretsResponse.
type stubSecretsFetcher struct {
	resp    []*sdkpb.SecretResponse
	err     error
	rawResp *vaultcommon.GetSecretsResponse
	rawErr  error
}

func (s stubSecretsFetcher) GetSecrets(context.Context, *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	return s.resp, s.err
}

func (s stubSecretsFetcher) GetRawSecretsResponse(context.Context, *sdkpb.GetSecretsRequest, host.EncryptionKeyFetcher) (*vaultcommon.GetSecretsResponse, error) {
	return s.rawResp, s.rawErr
}

func (s stubSecretsFetcher) GetOwner() string { return "stub-owner" }

// newSecretsEventTestHelper builds an ExecutionHelper whose base engine only
// carries the workflow event labels, which is all the wrappers read from it.
func newSecretsEventTestHelper(t *testing.T, executionID string, fetcher SecretsFetcher) *ExecutionHelper {
	t.Helper()

	labels := map[string]string{
		platform.KeyWorkflowOwner: "0x0000000000000000000000000000000000000001",
		platform.KeyWorkflowID:    "wf-id",
		platform.KeyWorkflowName:  "wf-name",
	}
	be := &baseEngine{}
	be.loggerLabels.Store(&labels)

	return &ExecutionHelper{
		base:                be,
		WorkflowExecutionID: executionID,
		SecretsFetcher:      fetcher,
	}
}

func decodeCapabilityExecutionStarted(t *testing.T, observer beholdertest.Observer) *eventsv2.CapabilityExecutionStarted {
	t.Helper()
	msgs := observer.Messages(t, beholder.AttrKeyEntity, entityCapabilityExecutionStartedV2)
	require.Len(t, msgs, 1, "expected exactly one CapabilityExecutionStarted event")
	var out eventsv2.CapabilityExecutionStarted
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &out))
	return &out
}

func decodeCapabilityExecutionFinished(t *testing.T, observer beholdertest.Observer) *eventsv2.CapabilityExecutionFinished {
	t.Helper()
	msgs := observer.Messages(t, beholder.AttrKeyEntity, entityCapabilityExecutionFinishedV2)
	require.Len(t, msgs, 1, "expected exactly one CapabilityExecutionFinished event")
	var out eventsv2.CapabilityExecutionFinished
	require.NoError(t, proto.Unmarshal(msgs[0].Body, &out))
	return &out
}

// TestExecutionHelper_GetSecretsEmitsCapabilityExecutionEvents verifies that a
// successful vault get-secrets call from a normal workflow emits a Started/Finished
// pair keyed to the real workflow execution ID, so it surfaces in the CRE UI like
// any other capability call.
//
// NewObserver mutates process-global beholder state, so this test must not be parallel.
func TestExecutionHelper_GetSecretsEmitsCapabilityExecutionEvents(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	observer := beholdertest.NewObserver(t)

	const executionID = "exec-1"
	callbackID := int32(7)
	h := newSecretsEventTestHelper(t, executionID, stubSecretsFetcher{})

	_, err := h.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{CallbackId: callbackID})
	require.NoError(t, err)

	started := decodeCapabilityExecutionStarted(t, observer)
	assert.Equal(t, vaultcommon.CapabilityID, started.GetCapabilityID())
	assert.Equal(t, vaultcommon.MethodGetSecrets, started.GetMethod())
	assert.Equal(t, callbackID, started.GetStepRef())
	assert.Equal(t, executionID, started.GetWorkflowExecutionID())
	require.NotNil(t, started.GetWorkflow())
	assert.Equal(t, "wf-id", started.GetWorkflow().GetWorkflowID())
	assert.Equal(t, "wf-name", started.GetWorkflow().GetWorkflowName())
	assert.Equal(t, "0x0000000000000000000000000000000000000001", started.GetWorkflow().GetWorkflowOwner())

	finished := decodeCapabilityExecutionFinished(t, observer)
	assert.Equal(t, vaultcommon.CapabilityID, finished.GetCapabilityID())
	assert.Equal(t, vaultcommon.MethodGetSecrets, finished.GetMethod())
	assert.Equal(t, callbackID, finished.GetStepRef())
	assert.Equal(t, executionID, finished.GetWorkflowExecutionID())
	assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED, finished.GetStatus())
	assert.Empty(t, finished.GetError())
}

// TestExecutionHelper_GetSecretsEmitsErroredOnFailure verifies that a failed vault
// get-secrets call finishes as errored with the returned error, so the failure is
// visible in the CRE UI.
//
// NewObserver mutates process-global beholder state, so this test must not be parallel.
func TestExecutionHelper_GetSecretsEmitsErroredOnFailure(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	observer := beholdertest.NewObserver(t)

	wantErr := errors.New("vault unavailable")
	h := newSecretsEventTestHelper(t, "exec-1", stubSecretsFetcher{err: wantErr})

	_, err := h.GetSecrets(t.Context(), &sdkpb.GetSecretsRequest{CallbackId: 1})
	require.ErrorIs(t, err, wantErr)

	finished := decodeCapabilityExecutionFinished(t, observer)
	assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_FAILED, finished.GetStatus())
	assert.Equal(t, wantErr.Error(), finished.GetError())
}

// TestExecutionHelper_GetRawSecretsResponseEmitsCapabilityExecutionEvents verifies
// that a successful vault get-secrets call from a confidential (TEE) workflow — the
// path the confidential relay handler drives — emits the same Started/Finished pair,
// so confidential workflows surface in the CRE UI too.
//
// NewObserver mutates process-global beholder state, so this test must not be parallel.
func TestExecutionHelper_GetRawSecretsResponseEmitsCapabilityExecutionEvents(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	observer := beholdertest.NewObserver(t)

	const executionID = "exec-tee"
	callbackID := int32(3)
	base := newSecretsEventTestHelper(t, executionID, stubSecretsFetcher{})
	raw := &executionHelperWithRawSecrets{base}

	_, err := raw.GetRawSecretsResponse(t.Context(), &sdkpb.GetSecretsRequest{CallbackId: callbackID}, nil)
	require.NoError(t, err)

	started := decodeCapabilityExecutionStarted(t, observer)
	assert.Equal(t, vaultcommon.CapabilityID, started.GetCapabilityID())
	assert.Equal(t, vaultcommon.MethodGetSecrets, started.GetMethod())
	assert.Equal(t, callbackID, started.GetStepRef())
	assert.Equal(t, executionID, started.GetWorkflowExecutionID())

	finished := decodeCapabilityExecutionFinished(t, observer)
	assert.Equal(t, vaultcommon.CapabilityID, finished.GetCapabilityID())
	assert.Equal(t, vaultcommon.MethodGetSecrets, finished.GetMethod())
	assert.Equal(t, callbackID, finished.GetStepRef())
	assert.Equal(t, executionID, finished.GetWorkflowExecutionID())
	assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_SUCCEEDED, finished.GetStatus())
	assert.Empty(t, finished.GetError())
}

// TestExecutionHelper_GetRawSecretsResponseEmitsErroredOnFailure verifies that a
// failed confidential get-secrets call finishes as errored with the returned error.
//
// NewObserver mutates process-global beholder state, so this test must not be parallel.
func TestExecutionHelper_GetRawSecretsResponseEmitsErroredOnFailure(t *testing.T) { //nolint:paralleltest // uses beholdertest.NewObserver, a global singleton swap
	observer := beholdertest.NewObserver(t)

	wantErr := errors.New("vault unavailable")
	base := newSecretsEventTestHelper(t, "exec-tee", stubSecretsFetcher{rawErr: wantErr})
	raw := &executionHelperWithRawSecrets{base}

	_, err := raw.GetRawSecretsResponse(t.Context(), &sdkpb.GetSecretsRequest{CallbackId: 1}, nil)
	require.ErrorIs(t, err, wantErr)

	finished := decodeCapabilityExecutionFinished(t, observer)
	assert.Equal(t, eventsv2.ExecutionStatus_EXECUTION_STATUS_FAILED, finished.GetStatus())
	assert.Equal(t, wantErr.Error(), finished.GetError())
}
