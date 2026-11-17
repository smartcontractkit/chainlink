package remote_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
)

func createTestCapabilityInfo(id string, capType commoncap.CapabilityType) commoncap.CapabilityInfo {
	return commoncap.CapabilityInfo{
		ID:             id,
		CapabilityType: capType,
		Description:    "Test capability",
	}
}

func TestCombinedClient_Info(t *testing.T) {
	t.Parallel()

	info := createTestCapabilityInfo("test-capability-info", commoncap.CapabilityTypeAction)
	returnedInfo, err := remote.NewCombinedClient(info).Info(t.Context())
	require.NoError(t, err)
	assert.Equal(t, info, returnedInfo)
}

// forwardOp describes one forwarded API call: setup installs a mock for the
// method (returning err, or the wanted result on success); invoke calls the client.
type forwardOp struct {
	name   string
	setup  func(c remote.CombinedClient, method string, err error) (want any)
	invoke func(ctx context.Context, c remote.CombinedClient, method string) (any, error)
}

func TestCombinedClient_Forwarding(t *testing.T) {
	t.Parallel()

	ops := []forwardOp{
		{
			name: "RegisterTrigger",
			setup: func(c remote.CombinedClient, method string, err error) any {
				m := &mocks.TriggerCapability{}
				c.SetTriggerSubscriber(method, m)
				ch := make(<-chan commoncap.TriggerResponse, 1)
				if err != nil {
					m.On("RegisterTrigger", mock.Anything, mock.Anything).Return((<-chan commoncap.TriggerResponse)(nil), err)
					return nil
				}
				m.On("RegisterTrigger", mock.Anything, mock.Anything).Return(ch, nil)
				return ch
			},
			invoke: func(ctx context.Context, c remote.CombinedClient, method string) (any, error) {
				return c.RegisterTrigger(ctx, commoncap.TriggerRegistrationRequest{TriggerID: "id", Method: method})
			},
		},
		{
			name: "UnregisterTrigger",
			setup: func(c remote.CombinedClient, method string, err error) any {
				m := &mocks.TriggerCapability{}
				c.SetTriggerSubscriber(method, m)
				m.On("UnregisterTrigger", mock.Anything, mock.Anything).Return(err)
				return nil
			},
			invoke: func(ctx context.Context, c remote.CombinedClient, method string) (any, error) {
				return nil, c.UnregisterTrigger(ctx, commoncap.TriggerRegistrationRequest{TriggerID: "id", Method: method})
			},
		},
		{
			name: "Execute",
			setup: func(c remote.CombinedClient, method string, err error) any {
				m := &mocks.ExecutableCapability{}
				c.SetExecutableClient(method, m)
				m.On("Execute", mock.Anything, mock.Anything).Return(commoncap.CapabilityResponse{}, err)
				return commoncap.CapabilityResponse{}
			},
			invoke: func(ctx context.Context, c remote.CombinedClient, method string) (any, error) {
				return c.Execute(ctx, commoncap.CapabilityRequest{Method: method})
			},
		},
	}

	for _, op := range ops {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			newClient := func() remote.CombinedClient {
				return remote.NewCombinedClient(createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeAction))
			}

			t.Run("success", func(t *testing.T) {
				t.Parallel()
				c := newClient()
				want := op.setup(c, "method", nil)
				got, err := op.invoke(t.Context(), c, "method")
				require.NoError(t, err)
				assert.Equal(t, want, got)
			})

			t.Run("method not defined", func(t *testing.T) {
				t.Parallel()
				got, err := op.invoke(t.Context(), newClient(), "undefined-method")
				require.ErrorContains(t, err, "method undefined-method not defined")
				assert.Zero(t, got)
			})

			t.Run("error from underlying", func(t *testing.T) {
				t.Parallel()
				c := newClient()
				wantErr := errors.New("underlying failed")
				op.setup(c, "method", wantErr)
				got, err := op.invoke(t.Context(), c, "method")
				require.ErrorIs(t, err, wantErr)
				assert.Zero(t, got)
			})
		})
	}
}

func TestCombinedClient_WorkflowRegistrationNotSupported(t *testing.T) {
	t.Parallel()

	client := remote.NewCombinedClient(createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeTrigger))
	meta := commoncap.RegistrationMetadata{WorkflowID: "test-workflow"}

	err := client.RegisterToWorkflow(t.Context(), commoncap.RegisterToWorkflowRequest{Metadata: meta})
	require.ErrorContains(t, err, "RegisterToWorkflow is not supported by remote capabilities")

	err = client.UnregisterFromWorkflow(t.Context(), commoncap.UnregisterFromWorkflowRequest{Metadata: meta})
	require.ErrorContains(t, err, "UnregisterFromWorkflow is not supported by remote capabilities")
}

func TestCombinedClient_MultipleMethods(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	client := remote.NewCombinedClient(createTestCapabilityInfo("test-multi-capability", commoncap.CapabilityTypeAction))

	triggers := map[string]*mocks.TriggerCapability{"trigger-1": {}, "trigger-2": {}}
	executables := map[string]*mocks.ExecutableCapability{"exec-1": {}, "exec-2": {}}
	for method, m := range triggers {
		client.SetTriggerSubscriber(method, m)
		m.On("RegisterTrigger", ctx, commoncap.TriggerRegistrationRequest{TriggerID: method, Method: method}).
			Return((<-chan commoncap.TriggerResponse)(make(chan commoncap.TriggerResponse, 1)), nil).Once()
	}
	for method, m := range executables {
		client.SetExecutableClient(method, m)
		m.On("Execute", ctx, commoncap.CapabilityRequest{Method: method}).Return(commoncap.CapabilityResponse{}, nil).Once()
	}

	// each call reaches only the mock registered for its method
	for method := range triggers {
		_, err := client.RegisterTrigger(ctx, commoncap.TriggerRegistrationRequest{TriggerID: method, Method: method})
		require.NoError(t, err)
	}
	for method := range executables {
		_, err := client.Execute(ctx, commoncap.CapabilityRequest{Method: method})
		require.NoError(t, err)
	}
	for _, m := range triggers {
		m.AssertExpectations(t)
	}
	for _, m := range executables {
		m.AssertExpectations(t)
	}
}

func TestCombinedClient_RemoveAndRestoreMethods(t *testing.T) {
	t.Parallel()
	client := remote.NewCombinedClient(createTestCapabilityInfo("test-capability", commoncap.CapabilityTypeAction))
	trigger := &mocks.TriggerCapability{}
	executable := &mocks.ExecutableCapability{}
	client.SetTriggerSubscriber("Trigger", trigger)
	client.SetExecutableClient("Write", executable)
	client.SetExecutableClient("Keep", executable)
	client.SetTriggerSubscriber("Trigger", nil)
	client.SetExecutableClient("Write", nil)
	_, err := client.RegisterTrigger(t.Context(), commoncap.TriggerRegistrationRequest{Method: "Trigger"})
	require.ErrorContains(t, err, "method Trigger not defined")
	_, err = client.Execute(t.Context(), commoncap.CapabilityRequest{Method: "Write"})
	require.ErrorContains(t, err, "method Write not defined")
	assert.Same(t, executable, client.GetExecutableClient("Keep"))
	client.SetTriggerSubscriber("Trigger", trigger)
	client.SetExecutableClient("Write", executable)
	assert.Same(t, trigger, client.GetTriggerSubscriber("Trigger"))
	assert.Same(t, executable, client.GetExecutableClient("Write"))
}
