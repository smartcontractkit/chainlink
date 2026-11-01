package remote

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types/mocks"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// TestTriggerSubscriber_SendDueBuckets verifies sendDueBuckets catches up on every bucket
// elapsed since the last send, so dropped or late ticks never skip a bucket.
func TestTriggerSubscriber_SendDueBuckets(t *testing.T) {
	t.Parallel()
	const nBuckets = 10
	const nRegistrations = 200

	var peer p2ptypes.PeerID
	require.NoError(t, peer.UnmarshalText([]byte(utils.MustNewPeerID())))
	capInfo := commoncap.CapabilityInfo{ID: "cap_id@1", CapabilityType: commoncap.CapabilityTypeTrigger}
	cfg := &commoncap.RemoteTriggerConfig{
		RegistrationRefresh:     nBuckets * registrationBucketDuration,
		MinResponsesToAggregate: 1,
	}

	dispatcher := mocks.NewDispatcher(t)
	dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()
	s := NewTriggerSubscriber(capInfo.ID, "LogTrigger", dispatcher, limits.NewGateLimiter(false), logger.Test(t))
	require.NoError(t, s.initMetrics())
	require.NoError(t, s.SetConfig(cfg, capInfo, 2, commoncap.DON{ID: 1, Members: []p2ptypes.PeerID{peer}}, aggregation.NewDefaultModeAggregator(1)))
	for i := range nRegistrations {
		_, err := s.RegisterTrigger(t.Context(), commoncap.TriggerRegistrationRequest{
			TriggerID: fmt.Sprintf("trigger_%d", i),
			Metadata:  commoncap.RequestMetadata{WorkflowID: fmt.Sprintf("%064x", i)},
		})
		require.NoError(t, err)
	}
	require.Len(t, s.registrationBuckets, nBuckets)

	bucketSizes := func(secs ...int64) int {
		total := 0
		for _, sec := range secs {
			total += len(s.registrationBuckets[sec%nBuckets])
		}
		return total
	}

	tests := []struct {
		name         string
		lastSentSec  int64
		nowSec       int64
		wantLastSent int64
		wantSends    int
	}{
		{"same second sends nothing", 100, 100, 100, 0},
		{"next second sends one bucket", 100, 101, 101, bucketSizes(101)},
		{"dropped ticks catch up on skipped buckets", 100, 105, 105, bucketSizes(101, 102, 103, 104, 105)},
		{"long stall sends each registration once", 100, 1000, 1000, nRegistrations},
		{"clock moving backwards resyncs", 100, 50, 50, bucketSizes(50)},
	}
	//nolint:paralleltest // subtests share the dispatcher's recorded calls
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dispatcher.Calls = nil
			got := s.sendDueBuckets(t.Context(), tc.lastSentSec, tc.nowSec)
			require.Equal(t, tc.wantLastSent, got)
			sends := 0
			for _, c := range dispatcher.Calls {
				if c.Method == "Send" && c.Arguments.Get(1).(*types.MessageBody).Method == types.MethodRegisterTrigger {
					sends++
				}
			}
			require.Equal(t, tc.wantSends, sends)
		})
	}
}

func TestTriggerSubscriber_RegistrationOrgIDMetric(t *testing.T) {
	t.Parallel()
	const withOrgID, withoutOrgID = 3, 2

	for _, gateOpen := range []bool{false, true} {
		t.Run(fmt.Sprintf("gateOpen=%t", gateOpen), func(t *testing.T) {
			t.Parallel()
			var peer p2ptypes.PeerID
			require.NoError(t, peer.UnmarshalText([]byte(utils.MustNewPeerID())))
			capInfo := commoncap.CapabilityInfo{ID: "cap_id@1", CapabilityType: commoncap.CapabilityTypeTrigger}
			cfg := &commoncap.RemoteTriggerConfig{
				RegistrationRefresh:     registrationBucketDuration,
				MinResponsesToAggregate: 1,
			}

			dispatcher := mocks.NewDispatcher(t)
			dispatcher.On("Send", mock.Anything, mock.Anything).Return(nil).Maybe()
			s := NewTriggerSubscriber(capInfo.ID, "LogTrigger", dispatcher, limits.NewGateLimiter(gateOpen), logger.Test(t))

			// A local provider avoids mutating the global beholder client, which parallel tests share.
			reader := sdkmetric.NewManualReader()
			counter, err := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader)).Meter("test").Int64Counter("registration_org_id")
			require.NoError(t, err)
			s.metrics.registrationOrgIDCounter = counter

			require.NoError(t, s.SetConfig(cfg, capInfo, 2, commoncap.DON{ID: 1, Members: []p2ptypes.PeerID{peer}}, aggregation.NewDefaultModeAggregator(1)))
			for i := range withOrgID + withoutOrgID {
				md := commoncap.RequestMetadata{WorkflowID: fmt.Sprintf("%064x", i)}
				if i < withOrgID {
					md.OrgID = "org_1"
				}
				_, err := s.RegisterTrigger(t.Context(), commoncap.TriggerRegistrationRequest{
					TriggerID: fmt.Sprintf("trigger_%d", i),
					Metadata:  md,
				})
				require.NoError(t, err)
			}

			s.sendDueBuckets(t.Context(), 100, 1000)

			var rm metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(t.Context(), &rm))
			require.Len(t, rm.ScopeMetrics, 1)
			require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
			sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
			require.True(t, ok)

			got := map[[2]string]int64{}
			for _, dp := range sum.DataPoints {
				hasOrgID, _ := dp.Attributes.Value(attribute.Key("hasOrgID"))
				propagated, _ := dp.Attributes.Value(attribute.Key("orgIDPropagateEnabled"))
				got[[2]string{hasOrgID.AsString(), propagated.AsString()}] = dp.Value
			}
			require.Equal(t, map[[2]string]int64{
				{"true", strconv.FormatBool(gateOpen)}:  withOrgID,
				{"false", strconv.FormatBool(gateOpen)}: withoutOrgID,
			}, got)
		})
	}
}

func TestSubRegState_DeliverAfterClose(t *testing.T) {
	t.Parallel()
	for range 100 {
		r := &subRegState{callback: make(chan commoncap.TriggerResponse, 1), done: make(chan struct{})}
		r.close()
		require.False(t, r.deliver(commoncap.TriggerResponse{}, make(chan struct{})))
	}
}
