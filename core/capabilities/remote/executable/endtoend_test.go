package executable_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/e2etesting"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_RemoteExecutableCapability_ExecutionNotBlockedBySlowCapabilityExecution_AllAtOnce(t *testing.T) {
	ctx := testutils.Context(t)

	workflowIDToPause := map[string]time.Duration{}
	workflowIDToPause[workflowID1] = 1 * time.Minute
	workflowIDToPause[workflowID2] = 1 * time.Second
	capability := &TestSlowExecutionCapability{
		workflowIDToPause: workflowIDToPause,
	}

	var callCount int64

	numWorkflowPeers := int64(10)
	var testShuttingDown int32

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		shuttingDown := atomic.LoadInt32(&testShuttingDown)
		if shuttingDown != 0 {
			return
		}

		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "1s", mp.(map[string]any)["response"].(string))
			}

			atomic.AddInt64(&callCount, 1)
		}
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	method := func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest, workflowID1)
	}
	testRemoteExecutableCapability(ctx, t, capability, int(numWorkflowPeers), 9, timeOut, 10,
		9, timeOut, method, false)

	method = func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest, workflowID2)
	}
	testRemoteExecutableCapability(ctx, t, capability, int(numWorkflowPeers), 9, timeOut, 10,
		9, timeOut, method, false)

	require.Eventually(t, func() bool {
		count := atomic.LoadInt64(&callCount)

		if count == numWorkflowPeers {
			atomic.AddInt32(&testShuttingDown, 1)
			return true
		}

		return false
	}, 1*time.Minute, 10*time.Millisecond, "require 10 callbacks from 1s delay capability")
}

func Test_RemoteExecutableCapability_ExecutionNotBlockedBySlowCapabilityExecution_OneAtATime(t *testing.T) {
	ctx := testutils.Context(t)

	workflowIDToPause := map[string]time.Duration{}
	workflowIDToPause[workflowID1] = 1 * time.Minute
	workflowIDToPause[workflowID2] = 1 * time.Second
	capability := &TestSlowExecutionCapability{
		workflowIDToPause: workflowIDToPause,
	}

	var callCount int64

	numWorkflowPeers := int64(10)
	var testShuttingDown int32

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		shuttingDown := atomic.LoadInt32(&testShuttingDown)
		if shuttingDown != 0 {
			return
		}

		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "1s", mp.(map[string]any)["response"].(string))
			}

			atomic.AddInt64(&callCount, 1)
		}
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	method := func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest, workflowID1)
	}
	testRemoteExecutableCapability(ctx, t, capability, int(numWorkflowPeers), 9, timeOut, 10,
		9, timeOut, method, false)

	method = func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest, workflowID2)
	}
	testRemoteExecutableCapability(ctx, t, capability, int(numWorkflowPeers), 9, timeOut, 10,
		9, timeOut, method, false)

	require.Eventually(t, func() bool {
		count := atomic.LoadInt64(&callCount)

		if count == numWorkflowPeers {
			atomic.AddInt32(&testShuttingDown, 1)
			return true
		}

		return false
	}, 1*time.Minute, 10*time.Millisecond, "require 10 callbacks from 1s delay capability")
}

func Test_RemoteExecutableCapability_TransmissionSchedules(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-108")
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
			}
		}
	}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	timeOut := 10 * time.Minute

	capability := &TestCapability{}

	method := func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest, workflowID1)
	}
	testRemoteExecutableCapability(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, method, true)

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)
	method = func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, responseTest, workflowID1)
	}

	testRemoteExecutableCapability(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, method, true)
}

func Test_RemoteExecutionCapability_CapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	capability := &TestErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
			assert.ErrorContains(t, responseError, "failed to execute capability")
		}, workflowID1)
	})

	for _, method := range methods {
		testRemoteExecutableCapability(ctx, t, capability, 10, 9, 10*time.Minute, 10, 9, 10*time.Minute, method, true)
	}
}

func Test_RemoteExecutableCapability_RandomCapabilityError(t *testing.T) {
	ctx := testutils.Context(t)

	capability := &TestRandomErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
			assert.ErrorContains(t, responseError, "failed to execute capability")
		}, workflowID1)
	})

	for _, method := range methods {
		testRemoteExecutableCapability(ctx, t, capability, 10, 9, 1*time.Second, 10, 9, 10*time.Minute,
			method, true)
	}
}

func testRemoteExecutableCapability(ctx context.Context, t *testing.T, underlying commoncap.ExecutableCapability, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration,
	method func(ctx context.Context, caller commoncap.ExecutableCapability), waitForExecuteCalls bool) {
	lggr := logger.Test(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(e2etesting.NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(e2etesting.NewPeerID())))

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
	for i := range numWorkflowPeers {
		workflowPeerID := p2ptypes.PeerID{}
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(e2etesting.NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := e2etesting.NewTestAsyncMessageBroker(t, 1000)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := executable.NewServer(capInfo.ID, "", capabilityPeer, capabilityDispatcher, lggr)
		cfg := &commoncap.RemoteExecutableConfig{
			RequestHashExcludedAttributes: []string{},
			RequestTimeout:                capabilityNodeResponseTimeout,
			ServerMaxParallelRequests:     10,
		}
		require.NoError(t, capabilityNode.SetConfig(cfg, underlying, capInfo, capDonInfo, workflowDONs, nil))
		servicetest.Run(t, capabilityNode)
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
		capabilityNodes[i] = capabilityNode
	}

	workflowNodes := make([]commoncap.ExecutableCapability, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := executable.NewClient(capInfo.ID, "", workflowPeerDispatcher, lggr)
		err := workflowNode.SetConfig(capInfo, workflowDonInfo, workflowNodeTimeout, nil)
		require.NoError(t, err)
		servicetest.Run(t, workflowNode)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	servicetest.Run(t, broker)

	wg := &sync.WaitGroup{}
	wg.Add(len(workflowNodes))

	for _, caller := range workflowNodes {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(ctx, caller)
		}(caller)
	}
	if waitForExecuteCalls {
		wg.Wait()
	}
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

type TestSlowExecutionCapability struct {
	abstractTestCapability
	workflowIDToPause map[string]time.Duration
}

func (t *TestSlowExecutionCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	var delay time.Duration

	delay, ok := t.workflowIDToPause[request.Metadata.WorkflowID]
	if !ok {
		panic("workflowID not found")
	}

	select {
	case <-time.After(delay):
		break
	case <-ctx.Done():
		return commoncap.CapabilityResponse{}, nil
	}

	response, err := values.NewMap(map[string]any{"response": delay.String()})
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

func (t TestErrorCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New("an error")
}

func (t TestErrorCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New("an error")
}

type TestRandomErrorCapability struct {
	abstractTestCapability
}

func (t TestRandomErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (commoncap.CapabilityResponse, error) {
	return commoncap.CapabilityResponse{}, errors.New(uuid.New().String())
}

func (t TestRandomErrorCapability) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New(uuid.New().String())
}

func (t TestRandomErrorCapability) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New(uuid.New().String())
}

func executeCapability(ctx context.Context, t *testing.T, caller commoncap.ExecutableCapability, transmissionSchedule *values.Map, responseTest func(t *testing.T, response commoncap.CapabilityResponse, responseError error),
	workflowID string) {
	executeInputs, err := values.NewMap(
		map[string]any{
			"executeValue1": "aValue1",
		},
	)
	require.NoError(t, err)
	response, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          workflowID,
				WorkflowExecutionID: workflowExecutionID1,
			},
			Config: transmissionSchedule,
			Inputs: executeInputs,
		})

	responseTest(t, response, err)
}
