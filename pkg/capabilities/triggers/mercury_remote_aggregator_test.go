package triggers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	eventId    = "ev_id_1"
	timestamp  = 1000
	rawReport1 = "abcd"
	rawReport2 = "efgh"
)

func TestMercuryRemoteAggregator(t *testing.T) {
	agg := NewMercuryRemoteAggregator(logger.Nop())

	feed1Old := mercury.FeedReport{
		FeedID:               feedOne,
		BenchmarkPrice:       100,
		ObservationTimestamp: 100,
		FullReport:           []byte(rawReport1),
	}
	feed1New := mercury.FeedReport{
		FeedID:               feedOne,
		BenchmarkPrice:       200,
		ObservationTimestamp: 200,
		FullReport:           []byte(rawReport1),
	}
	feed2Old := mercury.FeedReport{
		FeedID:               feedTwo,
		BenchmarkPrice:       300,
		ObservationTimestamp: 300,
		FullReport:           []byte(rawReport2),
	}
	feed2New := mercury.FeedReport{
		FeedID:               feedTwo,
		BenchmarkPrice:       400,
		ObservationTimestamp: 400,
		FullReport:           []byte(rawReport2),
	}

	node1Resp, err := wrapReports([]mercury.FeedReport{feed1Old, feed2New}, eventId, 400)
	require.NoError(t, err)
	rawNode1Resp, err := pb.MarshalCapabilityResponse(node1Resp)
	require.NoError(t, err)
	node2Resp, err := wrapReports([]mercury.FeedReport{feed1New, feed2Old}, eventId, 300)
	require.NoError(t, err)
	rawNode2Resp, err := pb.MarshalCapabilityResponse(node2Resp)
	require.NoError(t, err)

	// aggregator should return latest value for each feedID
	aggResponse, err := agg.Aggregate(eventId, [][]byte{rawNode1Resp, rawNode2Resp})
	require.NoError(t, err)
	aggEvent := capabilities.TriggerEvent{}
	require.NoError(t, aggResponse.Value.UnwrapTo(&aggEvent))
	decodedReports, err := mercury.NewCodec().Unwrap(aggEvent.Payload)
	require.NoError(t, err)

	require.Len(t, decodedReports, 2)
	require.Equal(t, feed1New, decodedReports[0])
	require.Equal(t, feed2New, decodedReports[1])
}
