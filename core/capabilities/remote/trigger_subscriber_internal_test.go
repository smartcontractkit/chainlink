package remote

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	s := NewTriggerSubscriber(capInfo.ID, "LogTrigger", dispatcher, logger.Test(t))
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
			got := s.sendDueBuckets(tc.lastSentSec, tc.nowSec)
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
