package executable_test

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/synctest"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_RemoteExecutableCapability_ExecutionNotBlockedBySlowCapabilityExecution(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		schedule string
	}{
		{"AllAtOnce", transmission.Schedule_AllAtOnce},
		{"OneAtATime", transmission.Schedule_OneAtATime},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				ctx := t.Context()

				capability := &TestSlowExecutionCapability{
					workflowIDToPause: map[string]time.Duration{
						workflowID1: synctest.SlowDelay(),
						workflowID2: synctest.FastDelay(),
					},
				}

				const numWorkflowPeers = 10
				// Slow capability delay is 1m; timeout must exceed that for phase-1 assertions after synctest.Wait.
				requestTimeout := 2 * time.Minute

				harness := setupRemoteExecutableHarness(t, capability, numWorkflowPeers, 9, requestTimeout, 10, 9, requestTimeout)

				transmissionSchedule, err := values.NewMap(map[string]any{
					"schedule":   tc.schedule,
					"deltaStage": "10ms",
				})
				require.NoError(t, err)

				executeInputs, err := values.NewMap(map[string]any{
					"executeValue1": "aValue1",
				})
				require.NoError(t, err)

				var wgSlow sync.WaitGroup
				wgSlow.Add(len(harness.workflowNodes))
				for _, caller := range harness.workflowNodes {
					go func(caller commoncap.ExecutableCapability) {
						defer wgSlow.Done()
						executeCapability(ctx, t, caller, transmissionSchedule, executeInputs, func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
							if assert.NoError(t, responseError) {
								mp, err := response.Value.Unwrap()
								if assert.NoError(t, err) {
									assert.Equal(t, synctest.SlowDelay().String(), mp.(map[string]any)["response"].(string))
								}
							}
						}, workflowID1, workflowExecutionID1)
					}(caller)
				}

				var wgFast sync.WaitGroup
				wgFast.Add(len(harness.workflowNodes))
				for _, caller := range harness.workflowNodes {
					go func(caller commoncap.ExecutableCapability) {
						defer wgFast.Done()
						executeCapability(ctx, t, caller, transmissionSchedule, executeInputs, func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
							if assert.NoError(t, responseError) {
								mp, err := response.Value.Unwrap()
								if assert.NoError(t, err) {
									assert.Equal(t, synctest.FastDelay().String(), mp.(map[string]any)["response"].(string))
								}
							}
						}, workflowID2, workflowExecutionID2)
					}(caller)
				}

				wgFast.Wait()
				synctest.Wait()
				wgSlow.Wait()
			})
		})
	}
}

func Test_RemoteExecutableCapability_TransmissionSchedules(t *testing.T) {
	t.Parallel()

	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-108")
	ctx := t.Context()

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

	executeInputs, err := values.NewMap(map[string]any{
		"executeValue1": "aValue1",
	})
	require.NoError(t, err)

	method := func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, executeInputs, responseTest, workflowID1, workflowExecutionID1)
	}
	testRemoteExecutableCapability(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, method, true)

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)
	method = func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, executeInputs, responseTest, workflowID1, workflowExecutionID1)
	}

	testRemoteExecutableCapability(ctx, t, capability, 10, 9, timeOut, 10, 9, timeOut, method, true)
}

func Test_RemoteExecutionCapability_CapabilityError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	capability := &TestErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	executeInputs, err := values.NewMap(map[string]any{
		"executeValue1": "aValue1",
	})
	require.NoError(t, err)

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = make([]func(ctx context.Context, caller commoncap.ExecutableCapability), 0, 1)
	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, executeInputs, func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
			assert.ErrorContains(t, responseError, "failed to execute capability")
		}, workflowID1, workflowExecutionID1)
	})

	for _, method := range methods {
		testRemoteExecutableCapability(ctx, t, capability, 10, 9, 10*time.Minute, 10, 9, 10*time.Minute, method, true)
	}
}

func Test_RemoteExecutableCapability_RandomCapabilityError(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	capability := &TestRandomErrorCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	executeInputs, err := values.NewMap(map[string]any{
		"executeValue1": "aValue1",
	})
	require.NoError(t, err)

	var methods []func(ctx context.Context, caller commoncap.ExecutableCapability)

	methods = make([]func(ctx context.Context, caller commoncap.ExecutableCapability), 0, 1)
	methods = append(methods, func(ctx context.Context, caller commoncap.ExecutableCapability) {
		executeCapability(ctx, t, caller, transmissionSchedule, executeInputs, func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error) {
			assert.ErrorContains(t, responseError, "failed to execute capability")
		}, workflowID1, workflowExecutionID1)
	})

	for _, method := range methods {
		testRemoteExecutableCapability(ctx, t, capability, 10, 9, 1*time.Second, 10, 9, 10*time.Minute,
			method, true)
	}
}

type remoteExecutableHarness struct {
	workflowNodes []commoncap.ExecutableCapability
}

func setupRemoteExecutableHarness(t *testing.T, underlying commoncap.ExecutableCapability, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration) remoteExecutableHarness {
	t.Helper()

	lggr := logger.Test(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

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
		require.NoError(t, workflowPeerID.UnmarshalText([]byte(NewPeerID())))
		workflowPeers[i] = workflowPeerID
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      1,
		F:       workflowDonF,
	}

	broker := newTestAsyncMessageBroker(t, 1000)

	workflowDONs := map[uint32]commoncap.DON{
		workflowDonInfo.ID: workflowDonInfo,
	}

	for i := range numCapabilityPeers {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := executable.NewServer(capInfo.ID, "", capabilityPeer, capabilityDispatcher, limits.NewGateLimiter(false), lggr)
		cfg := &commoncap.RemoteExecutableConfig{
			RequestHashExcludedAttributes: []string{},
			RequestTimeout:                capabilityNodeResponseTimeout,
			ServerMaxParallelRequests:     10,
		}
		require.NoError(t, capabilityNode.SetConfig(cfg, underlying, capInfo, capDonInfo, workflowDONs, nil))
		servicetest.Run(t, capabilityNode)
		broker.RegisterReceiverNode(capabilityPeer, capabilityNode)
	}

	workflowNodes := make([]commoncap.ExecutableCapability, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		workflowNode := executable.NewClient(capInfo.ID, "", workflowPeerDispatcher, lggr)
		err := workflowNode.SetConfig(capInfo, workflowDonInfo, workflowNodeTimeout, nil, nil, 0)
		require.NoError(t, err)
		servicetest.Run(t, workflowNode)
		broker.RegisterReceiverNode(workflowPeers[i], workflowNode)
		workflowNodes[i] = workflowNode
	}

	servicetest.Run(t, broker)

	return remoteExecutableHarness{workflowNodes: workflowNodes}
}

func testRemoteExecutableCapability(ctx context.Context, t *testing.T, underlying commoncap.ExecutableCapability, numWorkflowPeers int, workflowDonF uint8, workflowNodeTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, capabilityNodeResponseTimeout time.Duration,
	method func(ctx context.Context, caller commoncap.ExecutableCapability), waitForExecuteCalls bool) {
	harness := setupRemoteExecutableHarness(t, underlying, numWorkflowPeers, workflowDonF, workflowNodeTimeout, numCapabilityPeers, capabilityDonF, capabilityNodeResponseTimeout)

	wg := &sync.WaitGroup{}
	wg.Add(len(harness.workflowNodes))

	for _, caller := range harness.workflowNodes {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(ctx, caller)
		}(caller)
	}
	if waitForExecuteCalls {
		wg.Wait()
	}
}

type testAsyncMessageBroker struct {
	services.Service
	eng *services.Engine
	t   *testing.T

	nodes map[p2ptypes.PeerID]remotetypes.Receiver

	sendCh chan *remotetypes.MessageBody
}

func newTestAsyncMessageBroker(t *testing.T, sendChBufferSize int) *testAsyncMessageBroker {
	b := &testAsyncMessageBroker{
		t:      t,
		nodes:  make(map[p2ptypes.PeerID]remotetypes.Receiver),
		sendCh: make(chan *remotetypes.MessageBody, sendChBufferSize),
	}
	b.Service, b.eng = services.Config{
		Name:  "testAsyncMessageBroker",
		Start: b.start,
	}.NewServiceEngine(logger.Test(t))
	return b
}

func (a *testAsyncMessageBroker) start(ctx context.Context) error {
	a.eng.Go(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-a.sendCh:
				receiverID := toPeerID(msg.Receiver)

				receiver, ok := a.nodes[receiverID]
				if !ok {
					panic("server not found for peer id")
				}

				receiver.Receive(a.t.Context(), msg)
			}
		}
	})
	return nil
}

func (a *testAsyncMessageBroker) NewDispatcherForNode(nodePeerID p2ptypes.PeerID) remotetypes.Dispatcher {
	return &nodeDispatcher{
		callerPeerID: nodePeerID,
		broker:       a,
	}
}

func (a *testAsyncMessageBroker) RegisterReceiverNode(nodePeerID p2ptypes.PeerID, node remotetypes.Receiver) {
	if _, ok := a.nodes[nodePeerID]; ok {
		panic("node already registered")
	}

	a.nodes[nodePeerID] = node
}

func (a *testAsyncMessageBroker) Send(msg *remotetypes.MessageBody) {
	a.sendCh <- msg
}

func toPeerID(id []byte) p2ptypes.PeerID {
	return [32]byte(id)
}

type broker interface {
	Send(msg *remotetypes.MessageBody)
}

type nodeDispatcher struct {
	callerPeerID p2ptypes.PeerID
	broker       broker
}

func (t *nodeDispatcher) Name() string {
	return "nodeDispatcher"
}

func (t *nodeDispatcher) Start(ctx context.Context) error {
	return nil
}

func (t *nodeDispatcher) Close() error {
	return nil
}

func (t *nodeDispatcher) Ready() error {
	return nil
}

func (t *nodeDispatcher) HealthReport() map[string]error {
	return nil
}

func (t *nodeDispatcher) Send(peerID p2ptypes.PeerID, msgBody *remotetypes.MessageBody) error {
	msgBody.Version = 1
	msgBody.Sender = t.callerPeerID[:]
	msgBody.Receiver = peerID[:]
	msgBody.Timestamp = time.Now().UnixMilli()
	t.broker.Send(msgBody)
	return nil
}

func (t *nodeDispatcher) SetReceiver(capabilityID string, donID uint32, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityID string, donID uint32) {}

func (t *nodeDispatcher) SetReceiverForMethod(capabilityID string, donID uint32, methodName string, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiverForMethod(capabilityID string, donID uint32, methodName string) {
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

func NewP2PPeerID(t *testing.T) p2ptypes.PeerID {
	id := p2ptypes.PeerID{}
	require.NoError(t, id.UnmarshalText([]byte(NewPeerID())))
	return id
}

func NewPeerID() string {
	var privKey [32]byte
	_, err := rand.Read(privKey[:])
	if err != nil {
		panic(err)
	}

	peerID := append(libp2pMagic(), privKey[:]...)

	return base58.Encode(peerID)
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}

func executeCapability(ctx context.Context, t *testing.T, caller commoncap.ExecutableCapability, transmissionSchedule *values.Map, executeInputs *values.Map, responseTest func(t *testing.T, response commoncap.CapabilityResponse, responseError error),
	workflowID, workflowExecutionID string) {
	response, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          workflowID,
				WorkflowExecutionID: workflowExecutionID,
			},
			Config: transmissionSchedule,
			Inputs: executeInputs,
		})

	responseTest(t, response, err)
}
