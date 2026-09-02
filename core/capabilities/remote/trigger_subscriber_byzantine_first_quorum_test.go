package remote_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	remotetypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	remoteMocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/synctest"
)

// reproduces the one-shot readiness defect
func TestTriggerSubscriber_ByzantineFirstQuorumSuppressesLaterHonestAggregation(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		const (
			numMembers   = 4
			minResponses = 3
			eventID      = "byz-first-quorum-event"
			triggerID    = "trigger-byz-first-quorum"
		)

		// ── Setup: subscriber with 4 cap-DON members and MinResponsesToAggregate=3 ──

		lggr := logger.Test(t)
		capInfo, capDon, workflowDon := buildTwoTestDONs(t, numMembers, 1)
		dispatcher := remoteMocks.NewDispatcher(t)
		dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

		config := &commoncap.RemoteTriggerConfig{
			RegistrationRefresh:     time.Hour,
			RegistrationExpiry:      100 * time.Second,
			MinResponsesToAggregate: minResponses,
			MessageExpiry:           100 * time.Second,
		}
		subscriber := remote.NewTriggerSubscriber(capInfo.ID, "method", dispatcher, lggr)
		agg := aggregation.NewDefaultModeAggregator(config.MinResponsesToAggregate)
		require.NoError(t, subscriber.SetConfig(config, capInfo, workflowDon.ID, capDon, agg))
		require.NoError(t, subscriber.Start(t.Context()))

		regReq := commoncap.TriggerRegistrationRequest{
			TriggerID: triggerID,
			Metadata:  commoncap.RequestMetadata{WorkflowID: workflowID1},
		}
		callbackCh, err := subscriber.RegisterTrigger(t.Context(), regReq)
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, subscriber.UnregisterTrigger(t.Context(), regReq))
			require.NoError(t, subscriber.Close())
		})

		// Two byte-distinct payloads.  Both are valid TriggerResponses but they
		// differ in Outputs, so the default-mode aggregator (which requires
		// minResponses byte-identical payloads) cannot find a quorum among a
		// 1-byzantine / 2-honest split.
		honestPayload := marshalTriggerResponseOutputs(t, eventID, map[string]any{"price": "8300"})
		byzantinePayload := marshalTriggerResponseOutputs(t, eventID, map[string]any{"price": "8300-tampered"})
		require.NotEqual(t, honestPayload, byzantinePayload,
			"honest and byzantine payloads must be byte-distinct for the aggregator to split")

		// Helper: build a MethodTriggerEvent message from a specific DON member
		// with a pre-marshaled payload.  All messages share the same eventID so
		// they land in the same MessageCache entry.
		msgFrom := func(memberIdx int, payload []byte) *remotetypes.MessageBody {
			return &remotetypes.MessageBody{
				Sender: capDon.Members[memberIdx][:],
				Method: remotetypes.MethodTriggerEvent,
				Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
					TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
						TriggerEventId: eventID,
						WorkflowIds:    []string{workflowID1},
						TriggerIds:     []string{triggerID},
					},
				},
				Payload: payload,
			}
		}

		// ── Phase 1: Deliver the mixed first quorum [byzantine, honest, honest] ──
		//
		//   distinct senders = 3 >= MinResponsesToAggregate
		//     => Ready(..., once=true) returns true and commits wasReady=true
		//   payloads = [byzantine, honest, honest]
		//     => only 2 identical => Aggregate() fails
		//   wasReady stays true (the defect)
		subscriber.Receive(t.Context(), msgFrom(0, byzantinePayload))
		subscriber.Receive(t.Context(), msgFrom(1, honestPayload))
		subscriber.Receive(t.Context(), msgFrom(2, honestPayload))

		// Aggregation failed on the mixed quorum — no callback expected.
		// Receive is synchronous, so the callback channel is deterministically
		// empty at this point; the non-blocking probe is stable under synctest.
		select {
		case <-callbackCh:
			t.Fatal("unexpected callback from mixed first quorum [byzantine, honest, honest]")
		default:
		}

		// ── Phase 2: Deliver the fourth (honest) response ──
		//
		//   The cache now holds [byzantine, honest, honest, honest].
		//   3 identical honest payloads >= MinResponsesToAggregate => aggregation
		//   WOULD succeed if Ready() were called again.  But wasReady=true makes
		//   Ready(..., once=true) return false immediately, so aggregation is
		//   never retried.
		subscriber.Receive(t.Context(), msgFrom(3, honestPayload))

		// Before the fix: the callback never fires (wasReady blocks the retry).
		// After the fix:  the callback fires with the honest aggregated response.
		// Under synctest the time.After fires deterministically (fake clock
		// advances instantly once all goroutines block) instead of waiting real
		// seconds.
		select {
		case resp := <-callbackCh:
			require.NotNil(t, resp.Event.Outputs,
				"aggregated response must carry outputs after the fix")
		case <-time.After(2 * time.Second):
			t.Fatal(
				"DEFECT [report 85306]: 3 matching honest responses are in the cache " +
					"after the mixed first quorum, but MessageCache.wasReady=true " +
					"prevents aggregation retry. Fix: only commit wasReady after " +
					"Aggregate() succeeds — call Ready(..., once=false) and mark " +
					"delivered post-aggregation, or reset wasReady=false on " +
					"aggregation error so the next response can re-attempt.",
			)
		}

		// ── Phase 3: Prove the subscriber is still functional with a fresh event ──
		//
		//   A new all-honest event should aggregate normally, confirming the
		//   subscriber/cache/aggregator pipeline is not permanently broken —
		//   only the attacked event was suppressed.
		const recoveryEventID = "recovery-event"
		recoveryPayload := marshalTriggerResponseOutputs(t, recoveryEventID, map[string]any{"price": "8400"})
		recoveryMsg := func(memberIdx int) *remotetypes.MessageBody {
			return &remotetypes.MessageBody{
				Sender: capDon.Members[memberIdx][:],
				Method: remotetypes.MethodTriggerEvent,
				Metadata: &remotetypes.MessageBody_TriggerEventMetadata{
					TriggerEventMetadata: &remotetypes.TriggerEventMetadata{
						TriggerEventId: recoveryEventID,
						WorkflowIds:    []string{workflowID1},
						TriggerIds:     []string{triggerID},
					},
				},
				Payload: recoveryPayload,
			}
		}
		subscriber.Receive(t.Context(), recoveryMsg(0))
		subscriber.Receive(t.Context(), recoveryMsg(1))
		subscriber.Receive(t.Context(), recoveryMsg(2))

		select {
		case resp := <-callbackCh:
			require.NotNil(t, resp.Event.Outputs,
				"recovery event must aggregate and deliver to callback")
		case <-time.After(2 * time.Second):
			t.Fatal("recovery event unexpectedly failed to aggregate — subscriber is not functional")
		}
	})
}

// marshalTriggerResponseOutputs creates a marshaled TriggerResponse with the
// given event ID and outputs map.  Two calls with different outputs produce
// byte-distinct payloads that the default-mode aggregator treats as different.
func marshalTriggerResponseOutputs(t *testing.T, eventID string, outputs map[string]any) []byte {
	t.Helper()
	val, err := values.NewMap(outputs)
	require.NoError(t, err)
	resp := commoncap.TriggerResponse{
		Event: commoncap.TriggerEvent{ID: eventID, Outputs: val},
	}
	marshaled, err := pb.MarshalTriggerResponse(resp)
	require.NoError(t, err)
	return marshaled
}
