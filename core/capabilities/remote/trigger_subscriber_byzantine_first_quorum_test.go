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

	// numMembers = 3f+1 with f=1: member 0 is byzantine, members 1-3 are honest.
	// Production configures MinResponsesToAggregate = f+1 (system-tests set it to
	// faultyNodes+1), so the first aggregation attempt fires after only f+1
	// responses; 2f+1 is the stricter quorum. Both must survive a byzantine
	// response landing in the first quorum.
	for _, tc := range []struct {
		name         string
		minResponses uint32
	}{
		{name: "f+1", minResponses: 2},
		{name: "2f+1", minResponses: 3},
	} {
		t.Run(fmt.Sprintf("minResponses=%d(%s)", tc.minResponses, tc.name), func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				const (
					numMembers = 4
					eventID    = "byz-first-quorum-event"
					triggerID  = "trigger-byz-first-quorum"
				)

				// ── Setup: subscriber with 4 cap-DON members and MinResponsesToAggregate=tc.minResponses ──

				lggr := logger.Test(t)
				capInfo, capDon, workflowDon := buildTwoTestDONs(t, numMembers, 1)
				dispatcher := remoteMocks.NewDispatcher(t)
				dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()

				config := &commoncap.RemoteTriggerConfig{
					RegistrationRefresh:     time.Hour,
					RegistrationExpiry:      100 * time.Second,
					MinResponsesToAggregate: tc.minResponses,
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
				// 1-byzantine / (minResponses-1)-honest split.
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

				// ── Phase 1: Deliver the mixed first quorum [byzantine, honest×(minResponses-1)] ──
				//
				//   distinct senders = minResponses
				//     => Ready(..., once=true) returns true
				//   payloads = [byzantine, honest×(minResponses-1)]
				//     => only minResponses-1 identical => Aggregate() fails
				//   wasReady must stay false so aggregation can be retried
				//   (before the fix, Ready() itself committed wasReady=true)
				subscriber.Receive(t.Context(), msgFrom(0, byzantinePayload))
				for i := uint32(1); i < tc.minResponses; i++ {
					subscriber.Receive(t.Context(), msgFrom(int(i), honestPayload))
				}

				// Aggregation failed on the mixed quorum — no callback expected.
				// Receive is synchronous, so the callback channel is deterministically
				// empty at this point; the non-blocking probe is stable under synctest.
				select {
				case <-callbackCh:
					t.Fatal("unexpected callback from mixed first quorum [byzantine, honest×(minResponses-1)]")
				default:
				}

				// ── Phase 2: Deliver one more honest response ──
				//
				//   The cache now holds [byzantine, honest×minResponses].
				//   minResponses identical honest payloads >= MinResponsesToAggregate
				//   => aggregation WOULD succeed if Ready() were called again.
				//   Before the fix, wasReady=true made Ready(..., once=true) return
				//   false immediately, so aggregation was never retried.
				subscriber.Receive(t.Context(), msgFrom(int(tc.minResponses), honestPayload))

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
					t.Fatalf(
						"DEFECT [report 85306]: %d matching honest responses are in the cache "+
							"after the mixed first quorum, but MessageCache.wasReady=true "+
							"prevents aggregation retry. Fix: only commit wasReady after "+
							"Aggregate() succeeds — call Ready(..., once=false) and mark "+
							"delivered post-aggregation, or reset wasReady=false on "+
							"aggregation error so the next response can re-attempt.",
						tc.minResponses)
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
				for i := range tc.minResponses {
					subscriber.Receive(t.Context(), recoveryMsg(int(i)))
				}

				select {
				case resp := <-callbackCh:
					require.NotNil(t, resp.Event.Outputs,
						"recovery event must aggregate and deliver to callback")
				case <-time.After(2 * time.Second):
					t.Fatal("recovery event unexpectedly failed to aggregate — subscriber is not functional")
				}
			})
		})
	}
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
