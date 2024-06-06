package capability

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var _ capabilities.BaseCapability = (*mockBaseCapability)(nil)

type mockBaseCapability struct {
	info capabilities.CapabilityInfo
}

func (f *mockBaseCapability) Info(ctx context.Context) (capabilities.CapabilityInfo, error) {
	return f.info, nil
}

var _ capabilities.TriggerExecutable = (*mockTriggerExecutable)(nil)

type mockTriggerExecutable struct {
	callback chan capabilities.CapabilityResponse
}

func (f *mockTriggerExecutable) XXXTestingPushToCallbackChan(cr capabilities.CapabilityResponse) {
	f.callback <- cr
}

func (f *mockTriggerExecutable) RegisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	return f.callback, nil
}

func (f *mockTriggerExecutable) UnregisterTrigger(ctx context.Context, request capabilities.CapabilityRequest) error {
	f.callback = nil
	return nil
}

var _ capabilities.CallbackExecutable = (*mockCallbackExecutable)(nil)

type mockCallbackExecutable struct {
	registeredWorkflowRequest *capabilities.RegisterToWorkflowRequest
	callback                  chan capabilities.CapabilityResponse
}

func (f *mockCallbackExecutable) RegisterToWorkflow(ctx context.Context, request capabilities.RegisterToWorkflowRequest) error {
	f.registeredWorkflowRequest = &request
	return nil
}

func (f *mockCallbackExecutable) UnregisterFromWorkflow(ctx context.Context, request capabilities.UnregisterFromWorkflowRequest) error {
	f.registeredWorkflowRequest = nil
	return nil
}

func (f *mockCallbackExecutable) Execute(ctx context.Context, request capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	f.callback <- capabilities.CapabilityResponse{
		Value: nil,
		Err:   errors.New("some-error"),
	}
	return f.callback, nil
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

type mockTriggerCapability struct {
	*mockBaseCapability
	*mockTriggerExecutable
}

type mockActionCapability struct {
	*mockBaseCapability
	*mockCallbackExecutable
}

type mockConsensusCapability struct {
	*mockBaseCapability
	*mockCallbackExecutable
}

type mockTargetCapability struct {
	*mockBaseCapability
	*mockCallbackExecutable
}

type testRegistryPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	brokerExt *net.BrokerExt
	impl      *mocks.CapabilitiesRegistry
}

func (r *testRegistryPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, client *grpc.ClientConn) (any, error) {
	r.brokerExt.Broker = broker
	return NewCapabilitiesRegistryClient(client, r.brokerExt), nil
}

func (r *testRegistryPlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	r.brokerExt.Broker = broker
	pb.RegisterCapabilitiesRegistryServer(server, NewCapabilitiesRegistryServer(r.brokerExt, r.impl))
	return nil
}

func TestCapabilitiesRegistry(t *testing.T) {
	stopCh := make(chan struct{})
	logger := logger.Test(t)
	reg := mocks.NewCapabilitiesRegistry(t)

	capabilityResponse := capabilities.CapabilityResponse{
		Value: nil,
		Err:   errors.New("some-error"),
	}

	pluginName := "registry-test"
	client, server := plugin.TestPluginGRPCConn(
		t,
		true,
		map[string]plugin.Plugin{
			pluginName: &testRegistryPlugin{
				impl: reg,
				brokerExt: &net.BrokerExt{
					BrokerConfig: net.BrokerConfig{
						StopCh: stopCh,
						Logger: logger,
					},
				},
			},
		},
	)

	defer client.Close()
	defer server.Stop()

	regClient, err := client.Dispense(pluginName)
	require.NoError(t, err)

	rc, ok := regClient.(*capabilitiesRegistryClient)
	require.True(t, ok)

	// No capabilities in register
	reg.On("Get", mock.Anything, "some-id").Return(nil, errors.New("capability not found"))
	_, err = rc.Get(tests.Context(t), "some-id")
	require.ErrorContains(t, err, "capability not found")

	reg.On("GetAction", mock.Anything, "some-id").Return(nil, errors.New("capability not found"))
	_, err = rc.GetAction(tests.Context(t), "some-id")
	require.ErrorContains(t, err, "capability not found")

	reg.On("GetConsensus", mock.Anything, "some-id").Return(nil, errors.New("capability not found"))
	_, err = rc.GetConsensus(tests.Context(t), "some-id")
	require.ErrorContains(t, err, "capability not found")

	reg.On("GetTarget", mock.Anything, "some-id").Return(nil, errors.New("capability not found"))
	_, err = rc.GetTarget(tests.Context(t), "some-id")
	require.ErrorContains(t, err, "capability not found")

	reg.On("GetTrigger", mock.Anything, "some-id").Return(nil, errors.New("capability not found"))
	_, err = rc.GetTrigger(tests.Context(t), "some-id")
	require.ErrorContains(t, err, "capability not found")

	reg.On("List", mock.Anything).Return([]capabilities.BaseCapability{}, nil)
	list, err := rc.List(tests.Context(t))
	require.NoError(t, err)
	require.Len(t, list, 0)

	// Add capability Trigger
	triggerInfo := capabilities.CapabilityInfo{
		ID:             "trigger-1@1.0.0",
		CapabilityType: capabilities.CapabilityTypeTrigger,
		Description:    "trigger-1-description",
	}
	testTrigger := mockTriggerCapability{
		mockBaseCapability:    &mockBaseCapability{info: triggerInfo},
		mockTriggerExecutable: &mockTriggerExecutable{callback: make(chan capabilities.CapabilityResponse, 10)},
	}

	// After adding the trigger, we'll expect something wrapped by the internal client type below.
	reg.On("Add", mock.Anything, mock.AnythingOfType("*capability.TriggerCapabilityClient")).Return(nil)
	err = rc.Add(tests.Context(t), testTrigger)
	require.NoError(t, err)

	reg.On("GetTrigger", mock.Anything, "trigger-1@1.0.0").Return(testTrigger, nil)
	triggerCap, err := rc.GetTrigger(tests.Context(t), "trigger-1@1.0.0")
	require.NoError(t, err)

	// Test trigger Info()
	testCapabilityInfo(t, triggerInfo, triggerCap)

	// Test TriggerExecutable
	callbackChan, err := triggerCap.RegisterTrigger(tests.Context(t), capabilities.CapabilityRequest{
		Inputs: &values.Map{},
		Config: &values.Map{},
	})
	require.NoError(t, err)

	testTrigger.XXXTestingPushToCallbackChan(capabilityResponse)
	require.Equal(t, capabilityResponse, <-callbackChan)

	err = triggerCap.UnregisterTrigger(tests.Context(t), capabilities.CapabilityRequest{
		Inputs: &values.Map{},
		Config: &values.Map{},
	})
	require.NoError(t, err)
	require.Nil(t, testTrigger.callback)

	// Add capability Trigger
	actionInfo := capabilities.CapabilityInfo{
		ID:             "action-1@2.0.0",
		CapabilityType: capabilities.CapabilityTypeAction,
		Description:    "action-1-description",
	}

	actionCallbackChan := make(chan capabilities.CapabilityResponse, 10)
	testAction := mockActionCapability{
		mockBaseCapability:     &mockBaseCapability{info: actionInfo},
		mockCallbackExecutable: &mockCallbackExecutable{callback: actionCallbackChan},
	}
	reg.On("GetAction", mock.Anything, "action-1@2.0.0").Return(testAction, nil)
	actionCap, err := rc.GetAction(tests.Context(t), "action-1@2.0.0")
	require.NoError(t, err)

	testCapabilityInfo(t, actionInfo, actionCap)

	// Test Executable
	workflowRequest := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: "workflow-ID",
		},
	}
	err = actionCap.RegisterToWorkflow(tests.Context(t), workflowRequest)
	require.NoError(t, err)
	require.Equal(t, workflowRequest.Metadata.WorkflowID, testAction.registeredWorkflowRequest.Metadata.WorkflowID)

	actionCallbackChan <- capabilityResponse
	callbackChan, err = actionCap.Execute(tests.Context(t), capabilities.CapabilityRequest{})
	require.NoError(t, err)
	require.Equal(t, capabilityResponse, <-callbackChan)
	err = actionCap.UnregisterFromWorkflow(tests.Context(t), capabilities.UnregisterFromWorkflowRequest{})
	require.NoError(t, err)
	require.Nil(t, testAction.registeredWorkflowRequest)

	// Add capability Consensus
	consensusInfo := capabilities.CapabilityInfo{
		ID:             "consensus-1@3.0.0",
		CapabilityType: capabilities.CapabilityTypeConsensus,
		Description:    "consensus-1-description",
	}
	testConsensus := mockConsensusCapability{
		mockBaseCapability:     &mockBaseCapability{info: consensusInfo},
		mockCallbackExecutable: &mockCallbackExecutable{},
	}
	reg.On("GetConsensus", mock.Anything, "consensus-1@3.0.0").Return(testConsensus, nil)
	consensusCap, err := rc.GetConsensus(tests.Context(t), "consensus-1@3.0.0")
	require.NoError(t, err)

	testCapabilityInfo(t, consensusInfo, consensusCap)

	// Add capability Target
	targetInfo := capabilities.CapabilityInfo{
		ID:             "target-1@1.0.0",
		CapabilityType: capabilities.CapabilityTypeTarget,
		Description:    "target-1-description",
	}
	testTarget := mockTargetCapability{
		mockBaseCapability:     &mockBaseCapability{info: targetInfo},
		mockCallbackExecutable: &mockCallbackExecutable{},
	}
	reg.On("GetTarget", mock.Anything, "target-1@1.0.0").Return(testTarget, nil)
	targetCap, err := rc.GetTarget(tests.Context(t), "target-1@1.0.0")
	require.NoError(t, err)

	testCapabilityInfo(t, targetInfo, targetCap)
}

func testCapabilityInfo(t *testing.T, expectedInfo capabilities.CapabilityInfo, cap capabilities.BaseCapability) {
	gotInfo, err := cap.Info(tests.Context(t))
	require.NoError(t, err)
	require.Equal(t, expectedInfo.ID, gotInfo.ID)
	require.Equal(t, expectedInfo.CapabilityType, gotInfo.CapabilityType)
	require.Equal(t, expectedInfo.Description, gotInfo.Description)
	require.Equal(t, expectedInfo.Version(), gotInfo.Version())
}
