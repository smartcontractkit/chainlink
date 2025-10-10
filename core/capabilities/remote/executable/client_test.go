package executable_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

const (
	stepReferenceID1     = "step1"
	workflowID1          = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
	workflowID2          = "25c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce1"
	workflowExecutionID1 = "95ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeed"
	workflowExecutionID2 = "85ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0abbadeee"
	workflowOwnerID      = "0xAA"
)

func Test_Client_DonTopologies(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CAPPL-322")

	ctx := testutils.Context(t)

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
			}
		}
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	var methods []func(caller commoncap.ExecutableCapability)

	methods = append(methods, func(caller commoncap.ExecutableCapability) {
		executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
		if assert.NoError(t, err) {
			executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
		}
	})

	for _, method := range methods {
		testClient(t, 1, responseTimeOut, 1, 0,
			capability, method)

		testClient(t, 10, responseTimeOut, 1, 0,
			capability, method)

		testClient(t, 1, responseTimeOut, 10, 3,
			capability, method)

		testClient(t, 10, responseTimeOut, 10, 3,
			capability, method)

		testClient(t, 10, responseTimeOut, 10, 9,
			capability, method)
	}
}

func Test_Client_TransmissionSchedules(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-104")
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		if assert.NoError(t, responseError) {
			mp, err := response.Value.Unwrap()
			if assert.NoError(t, err) {
				assert.Equal(t, "aValue1", mp.(map[string]any)["response"].(string))
			}
		}
	}

	capability := &TestCapability{}

	responseTimeOut := 10 * time.Minute

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_OneAtATime,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(t, 1, responseTimeOut, 1, 0,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err2) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
	testClient(t, 10, responseTimeOut, 10, 3,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err2) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})

	transmissionSchedule, err = values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	testClient(t, 1, responseTimeOut, 1, 0,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
	testClient(t, 10, responseTimeOut, 10, 3,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
}

func Test_Client_TimesOutIfInsufficientCapabilityPeerResponses(t *testing.T) {
	ctx := testutils.Context(t)

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		assert.ErrorIs(t, responseError, executable.ErrRequestExpired)
	}

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	// number of capability peers is less than F + 1

	testClient(t, 10, 1*time.Second, 10, 11,
		capability,
		func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
}

func Test_Client_ContextCanceledBeforeQuorumReached(t *testing.T) {
	ctx, cancel := context.WithCancel(testutils.Context(t))

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		assert.ErrorIs(t, responseError, executable.ErrContextDoneBeforeResponseQuorum)
	}

	capability := &TestCapability{}
	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "20s",
	})
	require.NoError(t, err)

	cancel()
	testClient(t, 2, 20*time.Second, 2, 2,
		capability,
		func(caller commoncap.ExecutableCapability) {
			executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		})
}

func testClient(t *testing.T, numWorkflowPeers int, workflowNodeResponseTimeout time.Duration,
	numCapabilityPeers int, capabilityDonF uint8, underlying commoncap.ExecutableCapability,
	method func(caller commoncap.ExecutableCapability)) {
	lggr := logger.Test(t)

	capabilityPeers := make([]p2ptypes.PeerID, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityPeers[i] = NewP2PPeerID(t)
	}

	capDonInfo := commoncap.DON{
		ID:      1,
		Members: capabilityPeers,
		F:       capabilityDonF,
	}

	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1.0.0",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Executable Capability",
		DON:            &capDonInfo,
	}

	workflowPeers := make([]p2ptypes.PeerID, numWorkflowPeers)
	for i := range numWorkflowPeers {
		workflowPeers[i] = NewP2PPeerID(t)
	}

	workflowDonInfo := commoncap.DON{
		Members: workflowPeers,
		ID:      2,
	}

	broker := newTestAsyncMessageBroker(t, 100)

	receivers := make([]remotetypes.Receiver, numCapabilityPeers)
	for i := range numCapabilityPeers {
		capabilityDispatcher := broker.NewDispatcherForNode(capabilityPeers[i])
		receiver := newTestServer(capabilityPeers[i], capabilityDispatcher, workflowDonInfo, underlying)
		broker.RegisterReceiverNode(capabilityPeers[i], receiver)
		receivers[i] = receiver
	}

	callers := make([]commoncap.ExecutableCapability, numWorkflowPeers)

	for i := range numWorkflowPeers {
		workflowPeerDispatcher := broker.NewDispatcherForNode(workflowPeers[i])
		caller := executable.NewClient(capInfo.ID, "", workflowPeerDispatcher, lggr)
		err := caller.SetConfig(capInfo, workflowDonInfo, workflowNodeResponseTimeout, nil)
		require.NoError(t, err)
		servicetest.Run(t, caller)
		broker.RegisterReceiverNode(workflowPeers[i], caller)
		callers[i] = caller
	}

	servicetest.Run(t, broker)

	wg := &sync.WaitGroup{}
	wg.Add(len(callers))

	// Fire off all the requests
	for _, caller := range callers {
		go func(caller commoncap.ExecutableCapability) {
			defer wg.Done()
			method(caller)
		}(caller)
	}

	wg.Wait()
}

func executeMethod(ctx context.Context, caller commoncap.ExecutableCapability, transmissionSchedule *values.Map,
	executeInputs *values.Map, responseTest func(t *testing.T, responseCh commoncap.CapabilityResponse, responseError error), t *testing.T) {
	responseCh, err := caller.Execute(ctx,
		commoncap.CapabilityRequest{
			Metadata: commoncap.RequestMetadata{
				WorkflowID:          workflowID1,
				WorkflowExecutionID: workflowExecutionID1,
				WorkflowOwner:       workflowOwnerID,
			},
			Config: transmissionSchedule,
			Inputs: executeInputs,
		})

	responseTest(t, responseCh, err)
}

// Simple client that only responds once it has received a message from each workflow peer
type clientTestServer struct {
	peerID             p2ptypes.PeerID
	dispatcher         remotetypes.Dispatcher
	workflowDonInfo    commoncap.DON
	messageIDToSenders map[string]map[p2ptypes.PeerID]bool

	executableCapability commoncap.ExecutableCapability

	mux sync.Mutex
}

func newTestServer(peerID p2ptypes.PeerID, dispatcher remotetypes.Dispatcher, workflowDonInfo commoncap.DON,
	executableCapability commoncap.ExecutableCapability) *clientTestServer {
	return &clientTestServer{
		dispatcher:           dispatcher,
		workflowDonInfo:      workflowDonInfo,
		peerID:               peerID,
		messageIDToSenders:   make(map[string]map[p2ptypes.PeerID]bool),
		executableCapability: executableCapability,
	}
}

func (t *clientTestServer) Receive(_ context.Context, msg *remotetypes.MessageBody) {
	t.mux.Lock()
	defer t.mux.Unlock()

	sender := toPeerID(msg.Sender)
	messageID, err := executable.GetMessageID(msg)
	if err != nil {
		panic(err)
	}

	if t.messageIDToSenders[messageID] == nil {
		t.messageIDToSenders[messageID] = make(map[p2ptypes.PeerID]bool)
	}

	sendersOfMessageID := t.messageIDToSenders[messageID]
	if sendersOfMessageID[sender] {
		panic("received duplicate message")
	}

	sendersOfMessageID[sender] = true

	if len(t.messageIDToSenders[messageID]) == len(t.workflowDonInfo.Members) {
		switch msg.Method {
		case remotetypes.MethodExecute:
			capabilityRequest, err := pb.UnmarshalCapabilityRequest(msg.Payload)
			if err != nil {
				panic(err)
			}
			resp, responseErr := t.executableCapability.Execute(context.Background(), capabilityRequest)
			payload, marshalErr := pb.MarshalCapabilityResponse(resp)
			t.sendResponse(messageID, responseErr, payload, marshalErr)
		default:
			panic("unknown method")
		}
	}
}

func (t *clientTestServer) sendResponse(messageID string, responseErr error,
	payload []byte, marshalErr error) {
	for receiver := range t.messageIDToSenders[messageID] {
		var responseMsg = &remotetypes.MessageBody{
			CapabilityId:    "cap_id@1.0.0",
			CapabilityDonId: 1,
			CallerDonId:     t.workflowDonInfo.ID,
			Method:          remotetypes.MethodExecute,
			MessageId:       []byte(messageID),
			Sender:          t.peerID[:],
			Receiver:        receiver[:],
		}

		if responseErr != nil {
			responseMsg.Error = remotetypes.Error_INTERNAL_ERROR
		} else {
			if marshalErr != nil {
				panic(marshalErr)
			}
			responseMsg.Payload = payload
		}

		err := t.dispatcher.Send(receiver, responseMsg)
		if err != nil {
			panic(err)
		}
	}
}

func TestClient_SetConfig(t *testing.T) {
	lggr := logger.Test(t)
	capabilityID := "test_capability@1.0.0"

	// Create broker and dispatcher like other tests
	broker := newTestAsyncMessageBroker(t, 100)
	peerID := NewP2PPeerID(t)
	dispatcher := broker.NewDispatcherForNode(peerID)
	client := executable.NewClient(capabilityID, "execute", dispatcher, lggr)

	// Create valid test data
	validCapInfo := commoncap.CapabilityInfo{
		ID:             capabilityID,
		CapabilityType: commoncap.CapabilityTypeAction,
		Description:    "Test capability",
	}

	validDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
		F:       0,
	}

	validTimeout := 30 * time.Second

	t.Run("successful config set", func(t *testing.T) {
		transmissionConfig := &transmission.TransmissionConfig{
			Schedule:   transmission.Schedule_OneAtATime,
			DeltaStage: 10 * time.Millisecond,
		}

		err := client.SetConfig(validCapInfo, validDonInfo, validTimeout, transmissionConfig)
		require.NoError(t, err)

		// Verify config was set
		info, err := client.Info(context.Background())
		require.NoError(t, err)
		assert.Equal(t, validCapInfo.ID, info.ID)
	})

	t.Run("mismatched capability ID", func(t *testing.T) {
		invalidCapInfo := commoncap.CapabilityInfo{
			ID:             "different_capability@1.0.0",
			CapabilityType: commoncap.CapabilityTypeAction,
		}

		err := client.SetConfig(invalidCapInfo, validDonInfo, validTimeout, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "capability info provided does not match the client's capabilityID")
		assert.Contains(t, err.Error(), "different_capability@1.0.0 != test_capability@1.0.0")
	})

	t.Run("empty DON members", func(t *testing.T) {
		invalidDonInfo := commoncap.DON{
			ID:      1,
			Members: []p2ptypes.PeerID{},
			F:       0,
		}

		err := client.SetConfig(validCapInfo, invalidDonInfo, validTimeout, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty localDonInfo provided")
	})

	t.Run("successful config update", func(t *testing.T) {
		// Set initial config
		initialTimeout := 10 * time.Second
		err := client.SetConfig(validCapInfo, validDonInfo, initialTimeout, nil)
		require.NoError(t, err)

		// Replace with new config
		newTimeout := 60 * time.Second
		newDonInfo := commoncap.DON{
			ID:      2,
			Members: []p2ptypes.PeerID{NewP2PPeerID(t), NewP2PPeerID(t)},
			F:       1,
		}

		err = client.SetConfig(validCapInfo, newDonInfo, newTimeout, nil)
		require.NoError(t, err)

		// Verify the config was completely replaced
		info, err := client.Info(context.Background())
		require.NoError(t, err)
		assert.Equal(t, validCapInfo.ID, info.ID)
	})
}

func TestClient_SetConfig_StartClose(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	capabilityID := "test_capability@1.0.0"

	// Create broker and dispatcher like other tests
	broker := newTestAsyncMessageBroker(t, 100)
	peerID := NewP2PPeerID(t)
	dispatcher := broker.NewDispatcherForNode(peerID)
	client := executable.NewClient(capabilityID, "execute", dispatcher, lggr)

	validCapInfo := commoncap.CapabilityInfo{
		ID:             capabilityID,
		CapabilityType: commoncap.CapabilityTypeAction,
		Description:    "Test capability",
	}

	validDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
		F:       0,
	}

	validTimeout := 30 * time.Second

	t.Run("start fails without config", func(t *testing.T) {
		clientWithoutConfig := executable.NewClient(capabilityID, "execute", dispatcher, lggr)
		err := clientWithoutConfig.Start(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config not set - call SetConfig() before Start()")
	})

	t.Run("start succeeds after config set", func(t *testing.T) {
		require.NoError(t, client.SetConfig(validCapInfo, validDonInfo, validTimeout, nil))
		require.NoError(t, client.Start(ctx))
		require.NoError(t, client.Close())
	})

	t.Run("config can be updated after start", func(t *testing.T) {
		// Create a fresh client for this test since services can only be started once
		freshClient := executable.NewClient(capabilityID, "execute", dispatcher, lggr)

		// Set initial config and start
		require.NoError(t, freshClient.SetConfig(validCapInfo, validDonInfo, validTimeout, nil))
		require.NoError(t, freshClient.Start(ctx))

		// Update config while running
		validCapInfo.Description = "new description"
		require.NoError(t, freshClient.SetConfig(validCapInfo, validDonInfo, validTimeout, nil))

		// Verify config was updated
		info, err := freshClient.Info(ctx)
		require.NoError(t, err)
		assert.Equal(t, validCapInfo.Description, info.Description)

		// Clean up
		require.NoError(t, freshClient.Close())
	})
}
