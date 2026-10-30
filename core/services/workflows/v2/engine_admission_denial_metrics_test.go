package v2_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	metmocks "github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering/mocks"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// TestEngine_AdmissionDenial_ShardDenialCounters pins the metrics the
// admission-denial branch emits per outcome: the shard-denial sentinels map to
// their own dropReason values and shard-denial counters, while any other
// admission error keeps the generic admission_denied dropReason and fires no
// shard-denial counter.
func TestEngine_AdmissionDenial_ShardDenialCounters(t *testing.T) {
	t.Parallel()

	const (
		droppedTotal       = "platform_engine_trigger_event_dropped_total"
		shardNotOwnerTotal = "platform_engine_shard_execution_denied_not_owner"
		shardOrchErrTotal  = "platform_engine_shard_execution_denied_orchestrator_error"
	)

	testCases := []struct {
		name              string
		admissionErr      error
		wantDropReason    string
		wantShardNotOwner int64
		wantShardOrchErr  int64
	}{
		{
			name:              "not owner",
			admissionErr:      v2.ErrShardDeniedNotOwner,
			wantDropReason:    "shard_denied_not_owner",
			wantShardNotOwner: 1,
		},
		{
			name:             "orchestrator error",
			admissionErr:     v2.ErrShardDeniedOrchestrator,
			wantDropReason:   "shard_denied_orchestrator",
			wantShardOrchErr: 1,
		},
		{
			name:           "other admission error keeps generic reason",
			admissionErr:   errors.New("not owned by this node"),
			wantDropReason: "admission_denied",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			meter, reader := newTestMeter(t)

			// No module.Execute expectation: a denied event must never reach WASM.
			harness := newDropPathHarness(t, dropPathLegacy, metmocks.NewBillingClient(t), func(cfg *v2.EngineConfig) {
				cfg.Meter = meter
				cfg.Hooks.OnTriggerAdmission = func(context.Context, triggers.CoordinatedEvent) error {
					return tc.admissionErr
				}
			})

			harness.deliver(capabilities.TriggerResponse{
				Event: capabilities.TriggerEvent{TriggerType: "basic-trigger@1.0.0", ID: "event_admission_" + tc.name},
			})

			require.EventuallyWithT(t, func(c *assert.CollectT) {
				dropped, dropReasons := collectCounterByAttr(t, reader, droppedTotal, "dropReason")
				if !assert.Equal(c, int64(1), dropped, "the denied event must be counted as dropped exactly once") {
					return
				}
				assert.Equal(c, []string{tc.wantDropReason}, dropReasons, "dropReason must identify the admission outcome")
				// Both shard-denial counters are emitted before the drop total, so
				// they are final by the time the drop is observed.
				notOwner, _ := collectCounterByAttr(t, reader, shardNotOwnerTotal, "triggerID")
				assert.Equal(c, tc.wantShardNotOwner, notOwner, shardNotOwnerTotal)
				orchErr, _ := collectCounterByAttr(t, reader, shardOrchErrTotal, "triggerID")
				assert.Equal(c, tc.wantShardOrchErr, orchErr, shardOrchErrTotal)
			}, 5*time.Second, 50*time.Millisecond)
		})
	}
}
