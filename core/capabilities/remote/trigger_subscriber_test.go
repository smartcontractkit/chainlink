package remote_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	peerID1     = "12D3KooWF3dVeJ6YoT5HFnYhmwQWWMoEwVFzJQ5kKCMX3ZityxMC"
	peerID2     = "12D3KooWQsmok6aD8PZqt3RnJhQRrNzKHLficq7zYFRp7kZ1hHP8"
	peerID3     = "12D3KooWBCF1XT5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N"
	workflowID1 = "15c631d295ef5e32deb99a10ee6804bc4af13855687559d7ff6552ac6dbb2ce0"
)

var (
	triggerEvent1 = map[string]any{"event": "triggerEvent1"}
)

func TestTriggerSubscriber_RegisterAndReceive(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 1, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	// register trigger
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
	agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
	require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))
	require.NoError(t, subscriber.Start(t.Context()))

	req := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(t.Context(), req)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req))
		// calling UnregisterTrigger repeatedly is safe
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req))
		require.NoError(t, subscriber.Close())
	})

	// receive trigger event
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	triggerEvent := buildTriggerEvent(t, capDon.Members[0][:])
	subscriber.Receive(t.Context(), triggerEvent)
	response := <-triggerEventCallbackCh
	require.Equal(t, response.Event.Outputs, triggerEventValue)
}

func TestTriggerSubscriber_CorrectEventExpiryCheck(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 3, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	// register trigger
	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      10 * time.Second,
		MinResponsesToAggregate: 2,
		MessageExpiry:           10 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
	agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
	require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))

	require.NoError(t, subscriber.Start(t.Context()))
	regReq := commoncap.TriggerRegistrationRequest{
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	triggerEventCallbackCh, err := subscriber.RegisterTrigger(t.Context(), regReq)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), regReq))
		require.NoError(t, subscriber.Close())
	})

	// receive trigger events:
	// cleanup loop happens every 10 seconds, at 0:00, 0:10, 0:20, etc.
	// send the event from the first node around 0:02 (this is a bad node
	// that sends it too early)
	triggerEvent := buildTriggerEvent(t, capDon.Members[0][:])
	time.Sleep(2 * time.Second)
	subscriber.Receive(t.Context(), triggerEvent)

	// send events from nodes 2 & 3 (the good ones) around 0:15 so that
	// the diff between 0:02 and 0:15 exceeds the expiry threshold but
	// we don't hit the cleanup loop yet
	time.Sleep(13 * time.Second)
	triggerEvent.Sender = capDon.Members[1][:]
	subscriber.Receive(t.Context(), triggerEvent)
	// the aggregation shouldn't happen after events 1 and 2 as they
	// were received too far apart in time
	require.Empty(t, triggerEventCallbackCh)
	triggerEvent.Sender = capDon.Members[2][:]
	subscriber.Receive(t.Context(), triggerEvent)

	// event should be processed
	response := <-triggerEventCallbackCh
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	require.Equal(t, response.Event.Outputs, triggerEventValue)
}

func TestTriggerSubscriber_SetConfig_Basic(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 3, 1)
	agg := aggregation.NewDefaultModeAggregator(1)

	t.Run("returns error when capability info ID doesn't match subscriber's ID", func(t *testing.T) {
		dispatcher := remoteMocks.NewDispatcher(t)
		subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
		config := &commoncap.RemoteTriggerConfig{}
		mismatchedCapInfo := commoncap.CapabilityInfo{ID: "different_id", CapabilityType: commoncap.CapabilityTypeTrigger}
		err := subscriber.SetConfig(config, mismatchedCapInfo, workflowDon.ID, capDon, agg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "capability info provided does not match")
		require.Contains(t, err.Error(), "different_id")
		require.Contains(t, err.Error(), capInfo.ID)
	})

	t.Run("returns error when aggregator is nil", func(t *testing.T) {
		dispatcher := remoteMocks.NewDispatcher(t)
		subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
		config := &commoncap.RemoteTriggerConfig{}
		err := subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "aggregator not set")
	})

	t.Run("updates existing config", func(t *testing.T) {
		dispatcher := remoteMocks.NewDispatcher(t)
		subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
		// Set initial config
		initialConfig := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     100 * time.Millisecond,
			MinResponsesToAggregate: 1,
			MessageExpiry:           100 * time.Second,
		}
		err := subscriber.SetConfig(initialConfig, capInfo, workflowDon.ID, capDon, agg)
		require.NoError(t, err)

		// Update with new config
		newConfig := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     500 * time.Millisecond,
			MinResponsesToAggregate: 3,
			MessageExpiry:           500 * time.Second,
		}
		err = subscriber.SetConfig(newConfig, capInfo, workflowDon.ID, capDon, agg)
		require.NoError(t, err)

		// Verify updated config works
		require.NoError(t, subscriber.Start(t.Context()))
		require.NoError(t, subscriber.Close())
	})
	t.Run("handles nil initial config", func(t *testing.T) {
		dispatcher := remoteMocks.NewDispatcher(t)
		subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
		// Set initial config as nil
		err := subscriber.SetConfig(nil, capInfo, workflowDon.ID, capDon, agg)
		require.NoError(t, err)

		// Verify config works
		require.NoError(t, subscriber.Start(t.Context()))
		require.NoError(t, subscriber.Close())
	})
}

func TestTriggerSubscriber_MultipleTriggersSameWorkflow(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 1, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
	agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
	require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))
	require.NoError(t, subscriber.Start(t.Context()))

	// Register two triggers for the same workflow with different triggerIDs
	req1 := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger1",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	req2 := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger2",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	callbackCh1, err := subscriber.RegisterTrigger(t.Context(), req1)
	require.NoError(t, err)
	callbackCh2, err := subscriber.RegisterTrigger(t.Context(), req2)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req1))
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req2))
		require.NoError(t, subscriber.Close())
	})

	// Send message for trigger1 - should only go to callbackCh1
	triggerEvent1Msg := buildTriggerEventWithTriggerID(t, capDon.Members[0][:], workflowID1, "trigger1", "event1")
	subscriber.Receive(t.Context(), triggerEvent1Msg)

	resp := <-callbackCh1
	require.NotNil(t, resp.Event.Outputs)

	select {
	case <-callbackCh2:
		t.Fatal("did not expect message on callbackCh2")
	default:
		// expected - no message on callbackCh2
	}

	// Send message for trigger2 - should only go to callbackCh2
	triggerEvent2Msg := buildTriggerEventWithTriggerID(t, capDon.Members[0][:], workflowID1, "trigger2", "event2")
	subscriber.Receive(t.Context(), triggerEvent2Msg)

	resp = <-callbackCh2
	require.NotNil(t, resp.Event.Outputs)

	select {
	case <-callbackCh1:
		t.Fatal("did not expect another message on callbackCh1")
	default:
		// expected - no message on callbackCh1
	}
}

func TestTriggerSubscriber_LegacyMessageWithoutTriggerID(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 1, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
	agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
	require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))
	require.NoError(t, subscriber.Start(t.Context()))

	// Register single trigger (legacy style, with triggerID but receiving messages without it)
	req := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger1",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	callbackCh, err := subscriber.RegisterTrigger(t.Context(), req)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req))
		require.NoError(t, subscriber.Close())
	})

	// Send legacy message without triggerID - should still route to the single registered trigger
	legacyMsg := buildTriggerEvent(t, capDon.Members[0][:])
	subscriber.Receive(t.Context(), legacyMsg)

	resp := <-callbackCh
	require.NotNil(t, resp.Event.Outputs)
}

func TestTriggerSubscriber_AckReplayOnDuplicateReceive(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 3, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
	agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
	require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))
	require.NoError(t, subscriber.Start(t.Context()))

	triggerRegID := fmt.Sprintf("trigger_reg_%s_%d", workflowID1, 0)
	req := commoncap.TriggerRegistrationRequest{
		TriggerID: triggerRegID,
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	callbackCh, err := subscriber.RegisterTrigger(t.Context(), req)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req))
		require.NoError(t, subscriber.Close())
	})

	eventID := "event-ack-replay-1"
	msg := buildTriggerEventWithTriggerID(t, capDon.Members[0][:], workflowID1, triggerRegID, eventID)
	subscriber.Receive(t.Context(), msg)
	<-callbackCh

	require.NoError(t, subscriber.AckEvent(t.Context(), triggerRegID, eventID, "method"))

	dispatcher.Calls = nil
	subscriber.Receive(t.Context(), msg)

	ackSends := 0
	for _, call := range dispatcher.Calls {
		if call.Method != "Send" {
			continue
		}
		m := call.Arguments.Get(1).(*remotetypes.MessageBody)
		if m.Method == remotetypes.MethodTriggerEventAck {
			ackSends++
		}
	}
	require.Equal(t, len(capDon.Members), ackSends, "duplicate receive should fan out ACK to all capability DON members")

	select {
	case r := <-callbackCh:
		t.Fatalf("expected no second aggregated delivery to engine, got %+v", r)
	default:
	}
}

func TestTriggerSubscriber_UnregisterOneTriggerKeepsOther(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 1, 1)
	dispatcher := remoteMocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

	config := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     100 * time.Millisecond,
		RegistrationExpiry:      100 * time.Second,
		MinResponsesToAggregate: 1,
		MessageExpiry:           100 * time.Second,
	}
	subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
	agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
	require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))
	require.NoError(t, subscriber.Start(t.Context()))

	// Register two triggers
	req1 := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger1",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}
	req2 := commoncap.TriggerRegistrationRequest{
		TriggerID: "trigger2",
		Metadata: commoncap.RequestMetadata{
			WorkflowID: workflowID1,
		},
	}

	_, err := subscriber.RegisterTrigger(t.Context(), req1)
	require.NoError(t, err)
	callbackCh2, err := subscriber.RegisterTrigger(t.Context(), req2)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req2))
		require.NoError(t, subscriber.Close())
	})

	// Unregister trigger1
	require.NoError(t, subscriber.UnregisterTrigger(t.Context(), req1))

	// trigger2 should still work
	triggerEvent2Msg := buildTriggerEventWithTriggerID(t, capDon.Members[0][:], workflowID1, "trigger2", "event2")
	subscriber.Receive(t.Context(), triggerEvent2Msg)

	resp := <-callbackCh2
	require.NotNil(t, resp.Event.Outputs)
}

type subscriberSvc interface {
	remote.TriggerSubscriber
	services.Service
}

func TestTriggerSubscriber_RegistrationCheck(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	capInfo, capDon, workflowDon := buildTwoTestDONs(t, 1, 1)

	newSubscriber := func(t *testing.T) (subscriberSvc, *remoteMocks.Dispatcher) {
		dispatcher := remoteMocks.NewDispatcher(t)
		dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

		cfg := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     time.Hour,
			MinResponsesToAggregate: 1,
			MessageExpiry:           time.Minute,
		}

		sub := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
		agg := aggregation.NewDefaultModeAggregator(1)

		require.NoError(t, sub.SetConfig(cfg, capInfo, workflowDon.ID, capDon, agg))
		require.NoError(t, sub.Start(t.Context()))

		t.Cleanup(func() {
			require.NoError(t, sub.Close())
		})

		return sub, dispatcher
	}

	buildCheckMsg := func(workflowID, triggerID string) *remotetypes.MessageBody {
		return &remotetypes.MessageBody{
			Sender:      capDon.Members[0][:],
			Method:      remotetypes.MethodTriggerRegistrationCheck,
			CallerDonId: workflowDon.ID,
			Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
					WorkflowIds: []string{workflowID},
					TriggerIds:  []string{triggerID},
				},
			},
		}
	}

	t.Run("does not resend registration when trigger exists", func(t *testing.T) {
		sub, dispatcher := newSubscriber(t)

		_, err := sub.RegisterTrigger(t.Context(), commoncap.TriggerRegistrationRequest{
			TriggerID: "triggerA",
			Metadata: commoncap.RequestMetadata{
				WorkflowID: workflowID1,
			},
		})
		require.NoError(t, err)

		dispatcher.Calls = nil
		sub.Receive(t.Context(), buildCheckMsg(workflowID1, "triggerA"))

		// With the batching optimization, the subscriber no longer resends
		// MethodRegisterTrigger in response to a check for an existing
		// registration — the periodic registrationLoop handles that, and
		// the publisher ignores duplicate re-registrations anyway. Only
		// unregisters are sent for missing registrations.
		for _, call := range dispatcher.Calls {
			if call.Method != "Send" {
				continue
			}
			msg := call.Arguments.Get(1).(*remotetypes.MessageBody)
			require.NotEqual(t, remotetypes.MethodRegisterTrigger, msg.Method,
				"should not resend registration for existing trigger")
		}
	})

	t.Run("sends unregister when trigger missing with correct metadata", func(t *testing.T) {
		sub, dispatcher := newSubscriber(t)

		dispatcher.Calls = nil
		sub.Receive(t.Context(), buildCheckMsg(workflowID1, "triggerA"))

		var found bool
		for _, call := range dispatcher.Calls {
			if call.Method != "Send" {
				continue
			}
			msg := call.Arguments.Get(1).(*remotetypes.MessageBody)
			if msg.Method == remotetypes.MethodUnregisterTrigger {
				require.Equal(t, capInfo.ID, msg.CapabilityId)
				require.Equal(t, capDon.ID, msg.CapabilityDonId)
				require.Equal(t, workflowDon.ID, msg.CallerDonId)
				meta := msg.GetTriggerEventMetadata()
				require.NotNil(t, meta)
				require.Equal(t, []string{workflowID1}, meta.WorkflowIds)
				require.Equal(t, []string{"triggerA"}, meta.TriggerIds)
				found = true
				break
			}
		}
		require.True(t, found, "expected a MethodUnregisterTrigger Send call")
	})

	t.Run("sends unregister after trigger is unregistered locally", func(t *testing.T) {
		sub, dispatcher := newSubscriber(t)

		req := commoncap.TriggerRegistrationRequest{
			TriggerID: "triggerA",
			Metadata: commoncap.RequestMetadata{
				WorkflowID: workflowID1,
			},
		}
		_, err := sub.RegisterTrigger(t.Context(), req)
		require.NoError(t, err)

		// While registered, a check should NOT send anything (no resend).
		dispatcher.Calls = nil
		sub.Receive(t.Context(), buildCheckMsg(workflowID1, "triggerA"))
		for _, call := range dispatcher.Calls {
			if call.Method != "Send" {
				continue
			}
			msg := call.Arguments.Get(1).(*remotetypes.MessageBody)
			require.NotEqual(t, remotetypes.MethodRegisterTrigger, msg.Method,
				"should not resend registration for existing trigger")
		}

		// Unregister locally
		require.NoError(t, sub.UnregisterTrigger(t.Context(), req))

		// Now the same check should result in unregister
		dispatcher.Calls = nil
		sub.Receive(t.Context(), buildCheckMsg(workflowID1, "triggerA"))
		dispatcher.AssertCalled(t, "Send", mock.Anything, mock.MatchedBy(func(m *remotetypes.MessageBody) bool {
			return m.Method == remotetypes.MethodUnregisterTrigger
		}))
	})

	t.Run("ignores check from unknown sender", func(t *testing.T) {
		sub, dispatcher := newSubscriber(t)

		unknownPeer := p2ptypes.PeerID{0xaa}
		checkMsg := &remotetypes.MessageBody{
			Sender:      unknownPeer[:],
			Method:      remotetypes.MethodTriggerRegistrationCheck,
			CallerDonId: workflowDon.ID,
			Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
				TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
					WorkflowIds: []string{workflowID1},
					TriggerIds:  []string{"triggerA"},
				},
			},
		}

		dispatcher.Calls = nil
		sub.Receive(t.Context(), checkMsg)

		// No registration or unregistration calls should have been made
		for _, call := range dispatcher.Calls {
			if call.Method == "Send" {
				msg := call.Arguments.Get(1).(*remotetypes.MessageBody)
				require.NotEqual(t, remotetypes.MethodRegisterTrigger, msg.Method, "should not re-register from unknown sender")
				require.NotEqual(t, remotetypes.MethodUnregisterTrigger, msg.Method, "should not unregister from unknown sender")
			}
		}
	})
}

func buildTwoTestDONs(t *testing.T, capDonSize int, workflowDonSize int) (commoncap.CapabilityInfo, commoncap.DON, commoncap.DON) {
	capInfo := commoncap.CapabilityInfo{
		ID:             "cap_id@1",
		CapabilityType: commoncap.CapabilityTypeTrigger,
		Description:    "Remote Trigger",
	}

	capDon := commoncap.DON{
		ID:      1,
		Members: []p2ptypes.PeerID{},
		F:       0,
	}
	for range capDonSize {
		pid := utils.MustNewPeerID()
		peer := p2ptypes.PeerID{}
		require.NoError(t, peer.UnmarshalText([]byte(pid)))
		capDon.Members = append(capDon.Members, peer)
	}

	workflowDon := commoncap.DON{
		ID:      2,
		Members: []p2ptypes.PeerID{},
		F:       0,
	}
	for range workflowDonSize {
		pid := utils.MustNewPeerID()
		peer := p2ptypes.PeerID{}
		require.NoError(t, peer.UnmarshalText([]byte(pid)))
		workflowDon.Members = append(workflowDon.Members, peer)
	}
	return capInfo, capDon, workflowDon
}

func buildTriggerEvent(t *testing.T, sender []byte) *remotetypes.MessageBody {
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	capResponse := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			Outputs: triggerEventValue,
		},
		Err: nil,
	}
	marshaled, err := pb.MarshalTriggerResponse(capResponse)
	require.NoError(t, err)

	return &remotetypes.MessageBody{
		Sender: sender,
		Method: remotetypes.MethodTriggerEvent,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				WorkflowIds: []string{workflowID1},
			},
		},
		Payload: marshaled,
	}
}

func buildTriggerEventWithTriggerID(t *testing.T, sender []byte, workflowID, triggerID, eventID string) *remotetypes.MessageBody {
	triggerEventValue, err := values.NewMap(triggerEvent1)
	require.NoError(t, err)
	capResponse := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{
			ID:      eventID,
			Outputs: triggerEventValue,
		},
		Err: nil,
	}
	marshaled, err := pb.MarshalTriggerResponse(capResponse)
	require.NoError(t, err)

	return &remotetypes.MessageBody{
		Sender: sender,
		Method: remotetypes.MethodTriggerEvent,
		Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
			TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
				TriggerEventId: eventID,
				WorkflowIds:    []string{workflowID},
				TriggerIds:     []string{triggerID},
			},
		},
		Payload: marshaled,
	}
}
