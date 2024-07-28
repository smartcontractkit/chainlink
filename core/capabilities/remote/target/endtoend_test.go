package target_test

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
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

func Test_RemoteTargetCapability_InsufficientCapabilityResponses(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		assert.NotNil(t, response.Err)
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

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		assert.NotNil(t, response.Err)
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

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
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

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
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

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		assert.Equal(t, "failed to execute capability: an error", response.Err.Error())
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

	responseTest := func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error) {
		require.NoError(t, responseError)
		response := <-responseCh
		assert.Equal(t, "request expired", response.Err.Error())
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
	responseTest func(t *testing.T, responseCh <-chan commoncap.CapabilityResponse, responseError error)) {
	lggr := logger.TestLogger(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeerID := p2ptypes.PeerID{}
		require.NoError(t, capabilityPeerID.UnmarshalText([]byte(NewPeerID())))
		capabilityPeers[i] = capabilityPeerID
	}

	capabilityPeerID := p2ptypes.PeerID{}
	require.NoError(t, capabilityPeerID.UnmarshalText([]byte(NewPeerID())))

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

	capabilityNodes := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := 0; i < numCapabilityPeers; i++ {
		capabilityPeer := capabilityPeers[i]
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeer)
		capabilityNode := target.NewServer(capabilityPeer, underlying, capInfo, capDonInfo, workflowDONs, capabilityDispatcher,
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
			responseCh, err := caller.Execute(ctx,
				commoncap.CapabilityRequest{
					Metadata: commoncap.RequestMetadata{
						WorkflowID:          "workflowID",
						WorkflowExecutionID: "workflowExecutionID",
					},
					Config: transmissionSchedule,
					Inputs: executeInputs,
				})

			responseTest(t, responseCh, err)
		}(caller)
	}

	wg.Wait()
}

type testAsyncMessageBroker struct {
	services.StateMachine
	t *testing.T

	nodes map[p2ptypes.PeerID]remotetypes.Receiver

	sendCh chan *remotetypes.MessageBody

	stopCh services.StopChan
	wg     sync.WaitGroup
}

func (a *testAsyncMessageBroker) HealthReport() map[string]error {
	return nil
}

func (a *testAsyncMessageBroker) Name() string {
	return "testAsyncMessageBroker"
}

func newTestAsyncMessageBroker(t *testing.T, sendChBufferSize int) *testAsyncMessageBroker {
	return &testAsyncMessageBroker{
		t:      t,
		nodes:  make(map[p2ptypes.PeerID]remotetypes.Receiver),
		stopCh: make(services.StopChan),
		sendCh: make(chan *remotetypes.MessageBody, sendChBufferSize),
	}
}

func (a *testAsyncMessageBroker) Start(ctx context.Context) error {
	return a.StartOnce("testAsyncMessageBroker", func() error {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for {
				select {
				case <-a.stopCh:
					return
				case msg := <-a.sendCh:
					receiverId := toPeerID(msg.Receiver)

					receiver, ok := a.nodes[receiverId]
					if !ok {
						panic("server not found for peer id")
					}

					receiver.Receive(tests.Context(a.t), msg)
				}
			}
		}()
		return nil
	})
}

func (a *testAsyncMessageBroker) Close() error {
	return a.StopOnce("testAsyncMessageBroker", func() error {
		close(a.stopCh)

		a.wg.Wait()
		return nil
	})
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

func (t *nodeDispatcher) SetReceiver(capabilityId string, donId uint32, receiver remotetypes.Receiver) error {
	return nil
}
func (t *nodeDispatcher) RemoveReceiver(capabilityId string, donId uint32) {}

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

func (t TestCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	ch := make(chan commoncap.CapabilityResponse, 1)

	value := request.Inputs.Underlying["executeValue1"]

	response, err := values.NewMap(map[string]any{"response": value})
	if err != nil {
		return nil, err
	}
	ch <- commoncap.CapabilityResponse{
		Value: response,
	}

	return ch, nil
}

type TestErrorCapability struct {
	abstractTestCapability
}

func (t TestErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	return nil, errors.New("an error")
}

type TestRandomErrorCapability struct {
	abstractTestCapability
}

func (t TestRandomErrorCapability) Execute(ctx context.Context, request commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	return nil, errors.New(uuid.New().String())
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

	return base58.Encode(peerID[:])
}

func libp2pMagic() []byte {
	return []byte{0x00, 0x24, 0x08, 0x01, 0x12, 0x20}
}
