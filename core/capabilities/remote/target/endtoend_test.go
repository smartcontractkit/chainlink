package target_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	remoteutils "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/testutils"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_RemoteTargetCapability_InsufficientCapabilityResponses(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
		assert.NotNil(t, responseError)
	}

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testRemoteTarget(ctx, t, capability, 10, 9, 10*time.Millisecond, 10, 10, 10*time.Minute, transmissionSchedule, responseTest)
}

func Test_RemoteTargetCapability_InsufficientWorkflowRequests(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
		assert.NotNil(t, responseError)
	}

	timeOut := 10 * time.Minute

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testRemoteTarget(ctx, t, capability, 10, 10, 10*time.Millisecond, 10, 9, timeOut, transmissionSchedule, responseTest)
}

func Test_RemoteTargetCapability_TransmissionSchedules(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		mp, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	capability := &TestCapability{}

	testRemoteTarget(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, transmissionSchedule, responseTest)

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testRemoteTarget(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, transmissionSchedule, responseTest)
}

func Test_RemoteTargetCapability_DonTopologies(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		mp, err := response.Value.Unwrap()
		require.NoError(t, err)
		assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	capability := &TestCapability{}

	// Test scenarios where the number of submissions is greater than or equal to F + 1
	testRemoteTarget(ctx, t, capability, 1, 0, timeOut, 1, 0, timeOut, transmissionSchedule, responseTest)
	testRemoteTarget(ctx, t, capability, 4, 3, timeOut, 1, 0, timeOut, transmissionSchedule, responseTest)
	testRemoteTarget(ctx, t, capability, 10, 3, timeOut, 1, 0, timeOut, transmissionSchedule, responseTest)

	testRemoteTarget(ctx, t, capability, 1, 0, timeOut, 1, 0, timeOut, transmissionSchedule, responseTest)
	testRemoteTarget(ctx, t, capability, 1, 0, timeOut, 4, 3, timeOut, transmissionSchedule, responseTest)
	testRemoteTarget(ctx, t, capability, 1, 0, timeOut, 10, 3, timeOut, transmissionSchedule, responseTest)

	testRemoteTarget(ctx, t, capability, 4, 3, timeOut, 4, 3, timeOut, transmissionSchedule, responseTest)
	testRemoteTarget(ctx, t, capability, 10, 3, timeOut, 10, 3, timeOut, transmissionSchedule, responseTest)
	testRemoteTarget(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, transmissionSchedule, responseTest)
}

func Test_RemoteTargetCapability_CapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
		assert.Equal(t, "failed to execute capability: an error", responseError.Error())
	}

	capability := &TestErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testRemoteTarget(ctx, t, capability, 10, 9, 10*time.Minute, 10, 9, 10*time.Minute, transmissionSchedule, responseTest)
}

func Test_RemoteTargetCapability_RandomCapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		assert.Equal(t, "request expired", responseError.Error())
	}

	capability := &TestRandomErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testRemoteTarget(ctx, t, capability, 10, 9, 10*time.Millisecond, 10, 9, 10*time.Minute, transmissionSchedule, responseTest)
}

func testRemoteTarget(ctx context.Context, t *testing.T, underlying commoncap.TargetCapability, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration, transmissionSchedule *values.Map,
	responseTest func(t *testing.T, response commoncap.CapabilityResponse, responseError error)) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(remoteutils.NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(remoteutils.NewPeerID())))

	capDonInfo := commoncap.DON{
		ID:      2,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTarget,
		Description:    "Remote Target",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(remoteutils.NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := remoteutils.NewTestAsyncMessageBroker(t, 1000)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := target.NewServer(&commoncap.RemoteTargetConfig{RequestHashExcludedAttributes: []string{}}, capabilityPeer, underlying, capInfo, capDonInfo, workflowDONs, capabilityDispatcher,
			capabilityNodeResponseTimeout, lggr)
		servicetest.Run(t, capabilityNode)
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
	}

	workflowNodes := make([]commoncap.TargetCapability, numWorkflowPeers)
	for i := 0; i < numWorkflowPeers; i++ {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := target.NewClient(capInfo, workflowDonInfo, workflowPeerDispatcher, workflowNodeTimeout, lggr)
		servicetest.Run(t, workflowNode)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	servicetest.Run(t, broker)

	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)

	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(len(workflowNodes))

	for _, caller := range workflowNodes {
		go func(caller commoncap.TargetCapability) {
			defer wg.Done()
			response, err := caller.Execute(ctx,
				commoncap.CapabilityRequest{
					Metadata: commoncap.RequestMetadata{
						WorkflowID:          remoteutils.WorkflowID1,
						WorkflowExecutionID: remoteutils.WorkflowExecutionID1,
					},
					Config: transmissionSchedule,
					Inputs: executeInputs,
				})

			responseTest(t, response, err)
		}(caller)
	}

	wg.Wait()
}

type abstractTestCapability struct {
}

func (t abstractTestCapability) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return commoncap.CapabilityInfo{}, nil
}

func (t abstractTestCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return nil
}

func (t abstractTestCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return nil
}

type TestCapability struct {
	abstractTestCapability
}

func (t TestCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	value := request.Inputs.Underlying["executeValue1"]
	response, err := values.NewMap(map[string]any{"response": value})
	if err != nil {
		return commoncap.CapabilityResponse{}, err
	}
	return commoncap.CapabilityResponse{
		Value: response,
	}, nil
}

type TestErrorCapability struct {
	abstractTestCapability
}

func (t TestErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, errors.New("an error")
}

type TestRandomErrorCapability struct {
	abstractTestCapability
}

func (t TestRandomErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, errors.New(uuid.New().String())
}
