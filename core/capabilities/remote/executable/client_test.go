package executable_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/executable"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/synctest"
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
	t.Parallel()

	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/CAPPL-322")

	ctx := t.Context()

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

	methods = make([]func(caller commoncap.ExecutableCapability), 0, 1)
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
	t.Parallel()

	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-104")
	ctx := t.Context()

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
		},
	)
	testClient(t, 10, responseTimeOut, 10, 3,
		capability, func(caller commoncap.ExecutableCapability) {
			executeInputs, err2 := values.NewMap(map[string]any{"executeValue1": "aValue1"})
			if assert.NoError(t, err2) {
				executeMethod(ctx, caller, transmissionSchedule, executeInputs, responseTest, t)
			}
		},
	)

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

func Test_Client_ResponseAggregationGrace(t *testing.T) {
	t.Parallel()
	const (
		requestTimeout         = 500 * time.Millisecond
		capabilityExecuteDelay = 2 * time.Second
	)

	for _, tc := range []struct {
		name               string
		numWorkflowPeers   int
		numCapabilityPeers int
		capabilityDonF     uint8
		responseTest       func(t *testing.T, response commoncap.CapabilityResponse, responseError error)
	}{
		// Client cancels at requestTimeout + responseAggregationGrace (~10.5s).
		// Ten workflow peers force clientTestServer to wait for the full workflow DON
		// and block the broker during Execute; quorum (F+1=4) is theoretically
		// reachable but responses arrive after the grace window.
		{
			name:               "TimesOutWhenMockBlockingExceedsGrace",
			numWorkflowPeers:   10,
			numCapabilityPeers: 10,
			capabilityDonF:     3,
			responseTest:       assertClientRequestExpired,
		},
		// capabilityExecuteDelay exceeds requestTimeout but fits within
		// requestTimeout + responseAggregationGrace.
		// Use a single workflow peer so the mock cap server executes as soon as it
		// receives one request; with multiple workflow peers the mock waits for all
		// peers and blocks the test broker while executing, which exceeds the grace window.
		{
			name:               "WaitsForResponsesWithinGrace",
			numWorkflowPeers:   1,
			numCapabilityPeers: 4,
			capabilityDonF:     3,
			responseTest:       assertClientResponseWithDelay,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				ctx := t.Context()

				capability := &TestSlowExecutionCapability{
					workflowIDToPause: map[string]time.Duration{
						workflowID1: capabilityExecuteDelay,
					},
				}

				transmissionSchedule, err := values.NewMap(map[string]any{
					"schedule":   transmission.Schedule_AllAtOnce,
					"deltaStage": "10ms",
				})
				require.NoError(t, err)

				executeInputs, err := values.NewMap(map[string]any{"executeValue1": "aValue1"})
				require.NoError(t, err)

				testClient(t, tc.numWorkflowPeers, requestTimeout, tc.numCapabilityPeers, tc.capabilityDonF,
					capability,
					func(caller commoncap.ExecutableCapability) {
						executeMethod(ctx, caller, transmissionSchedule, executeInputs, tc.responseTest, t)
					},
				)
			})
		})
	}
}

func assertClientRequestExpired(t *testing.T, _ commoncap.CapabilityResponse, responseError error) {
	assert.ErrorIs(t, responseError, executable.ErrRequestExpired)
}

func assertClientResponseWithDelay(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
	if assert.NoError(t, responseError) {
		mp, err := response.Value.Unwrap()
		if assert.NoError(t, err) {
			assert.Equal(t, (2 * time.Second).String(), mp.(map[string]any)["response"].(string))
		}
	}
}

func Test_Client_ConsensusFailedIfInsufficientCapabilityPeerResponses(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	responseTest := func(t *testing.T, response commoncap.CapabilityResponse, responseError error) {
		var capErr caperrors.Error
		require.ErrorAs(t, responseError, &capErr)
		require.Equal(t, caperrors.ConsensusFailed, capErr.Code())
		require.Contains(t, capErr.Error(), "[100]ConsensusFailed: response quorum unreachable: not enough matching capability responses: received 1/10 peer responses with 1 unique payloads; best match count 1, need 12 (9 responses pending)")
	}

	capability := &TestCapability{}

	transmissionSchedule, err := values.NewMap(map[string]any{
		"schedule":   transmission.Schedule_AllAtOnce,
		"deltaStage": "10ms",
	})
	require.NoError(t, err)

	// F+1 exceeds peer count; first divergent response makes quorum unreachable.

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
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())

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

func testClient(t *testing.T, numWorkflowPeers int, workflowNodeResponseTimeout time.Duration, numCapabilityPeers int, capabilityDonF uint8, underlying commoncap.ExecutableCapability, method func(caller commoncap.ExecutableCapability)) {
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
		err := caller.SetConfig(capInfo, workflowDonInfo, workflowNodeResponseTimeout, nil, nil, 0)
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

func (t *clientTestServer) Receive(ctx context.Context, msg *remotetypes.MessageBody) {
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
			resp, responseErr := t.executableCapability.Execute(ctx, capabilityRequest)
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

type clientSetConfigTestFixture struct {
	Client interface {
		SetConfig(commoncap.CapabilityInfo, commoncap.DON, time.Duration, *transmission.TransmissionConfig, [][]byte, uint32) error
		Info(context.Context) (commoncap.CapabilityInfo, error)
		Start(context.Context) error
		Close() error
	}
	ValidCapInfo commoncap.CapabilityInfo
	ValidDonInfo commoncap.DON
	ValidTimeout time.Duration
	CapabilityID string
}

func newClientSetConfigTestFixture(t *testing.T) clientSetConfigTestFixture {
	t.Helper()

	capabilityID := "test_capability@1.0.0"
	broker := newTestAsyncMessageBroker(t, 100)
	peerID := NewP2PPeerID(t)
	dispatcher := broker.NewDispatcherForNode(peerID)
	client := executable.NewClient(capabilityID, "execute", dispatcher, logger.Test(t))

	validDonInfo := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{NewP2PPeerID(t)},
		F:       0,
	}

	validCapInfo := commoncap.CapabilityInfo{
		ID:             capabilityID,
		CapabilityType: commoncap.CapabilityTypeAction,
		Description:    "Test capability",
		DON:            &validDonInfo,
	}

	return clientSetConfigTestFixture{
		Client:       client,
		ValidCapInfo: validCapInfo,
		ValidDonInfo: validDonInfo,
		ValidTimeout: 30 * time.Second,
		CapabilityID: capabilityID,
	}
}

func TestClient_SetConfig(t *testing.T) {
	t.Parallel()

	t.Run("successful config set", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)

		transmissionConfig := &transmission.TransmissionConfig{
			Schedule:   transmission.Schedule_OneAtATime,
			DeltaStage: 10 * time.Millisecond,
		}

		err := fixture.Client.SetConfig(fixture.ValidCapInfo, fixture.ValidDonInfo, fixture.ValidTimeout, transmissionConfig, nil, 0)
		require.NoError(t, err)

		info, err := fixture.Client.Info(t.Context())
		require.NoError(t, err)
		assert.Equal(t, fixture.ValidCapInfo.ID, info.ID)
	})

	t.Run("mismatched capability ID", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)

		invalidCapInfo := commoncap.CapabilityInfo{
			ID:             "different_capability@1.0.0",
			CapabilityType: commoncap.CapabilityTypeAction,
		}

		err := fixture.Client.SetConfig(invalidCapInfo, fixture.ValidDonInfo, fixture.ValidTimeout, nil, nil, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "capability info provided does not match the client's capabilityID")
		assert.Contains(t, err.Error(), "different_capability@1.0.0 != test_capability@1.0.0")
	})

	t.Run("empty DON members", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)

		invalidDonInfo := commoncap.DON{
			ID:      1,
			Members: []p2ptypes.PeerID{},
			F:       0,
		}

		err := fixture.Client.SetConfig(fixture.ValidCapInfo, invalidDonInfo, fixture.ValidTimeout, nil, nil, 0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty localDonInfo provided")
	})

	t.Run("successful config update", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)

		initialTimeout := 10 * time.Second
		err := fixture.Client.SetConfig(fixture.ValidCapInfo, fixture.ValidDonInfo, initialTimeout, nil, nil, 0)
		require.NoError(t, err)

		newTimeout := 60 * time.Second
		newDonInfo := commoncap.DON{
			ID:      2,
			Members: []p2ptypes.PeerID{NewP2PPeerID(t), NewP2PPeerID(t)},
			F:       1,
		}

		err = fixture.Client.SetConfig(fixture.ValidCapInfo, newDonInfo, newTimeout, nil, nil, 0)
		require.NoError(t, err)

		info, err := fixture.Client.Info(t.Context())
		require.NoError(t, err)
		assert.Equal(t, fixture.ValidCapInfo.ID, info.ID)
	})
}

func TestClient_SetConfig_StartClose(t *testing.T) {
	t.Parallel()

	t.Run("start fails without config", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)

		err := fixture.Client.Start(t.Context())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "config not set - call SetConfig() before Start()")
	})

	t.Run("start succeeds after config set", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)
		ctx := t.Context()

		require.NoError(t, fixture.Client.SetConfig(fixture.ValidCapInfo, fixture.ValidDonInfo, fixture.ValidTimeout, nil, nil, 0))
		require.NoError(t, fixture.Client.Start(ctx))
		require.NoError(t, fixture.Client.Close())
	})

	t.Run("config can be updated after start", func(t *testing.T) {
		t.Parallel()

		fixture := newClientSetConfigTestFixture(t)
		ctx := t.Context()

		require.NoError(t, fixture.Client.SetConfig(fixture.ValidCapInfo, fixture.ValidDonInfo, fixture.ValidTimeout, nil, nil, 0))
		require.NoError(t, fixture.Client.Start(ctx))

		newCapInfo := fixture.ValidCapInfo
		newCapInfo.Description = "new description"
		require.NoError(t, fixture.Client.SetConfig(newCapInfo, fixture.ValidDonInfo, fixture.ValidTimeout, nil, nil, 0))

		info, err := fixture.Client.Info(ctx)
		require.NoError(t, err)
		assert.Equal(t, newCapInfo.Description, info.Description)

		require.NoError(t, fixture.Client.Close())
	})
}
