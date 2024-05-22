package triggers

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	eventID    = "ev_id_1"
	timestamp  = 1000
	rawReport1 = "abcd"
	rawReport2 = "efgh"
)

type testMercuryCodec struct {
}

func (c testMercuryCodec) UnwrapValid(wrapped values.Value, _ [][]byte, _ int) ([]datastreams.FeedReport, error) {
	dest := []datastreams.FeedReport{}
	err := wrapped.UnwrapTo(&dest)
	return dest, err
}

func (c testMercuryCodec) Wrap(reports []datastreams.FeedReport) (values.Value, error) {
	return values.Wrap(reports)
}

func TestMercuryRemoteAggregator(t *testing.T) {
	agg := NewMercuryRemoteAggregator(testMercuryCodec{}, nil, 0, logger.Nop())
	signatures := [][]byte{{1, 2, 3}}

	feed1Old := datastreams.FeedReport{
		FeedID:               feedOne,
		BenchmarkPrice:       big.NewInt(100).Bytes(),
		ObservationTimestamp: 100,
		FullReport:           []byte(rawReport1),
		ReportContext:        []byte{},
		Signatures:           signatures,
	}
	feed1New := datastreams.FeedReport{
		FeedID:               feedOne,
		BenchmarkPrice:       big.NewInt(200).Bytes(),
		ObservationTimestamp: 200,
		FullReport:           []byte(rawReport1),
		ReportContext:        []byte{},
		Signatures:           signatures,
	}
	feed2Old := datastreams.FeedReport{
		FeedID:               feedTwo,
		BenchmarkPrice:       big.NewInt(300).Bytes(),
		ObservationTimestamp: 300,
		FullReport:           []byte(rawReport2),
		ReportContext:        []byte{},
		Signatures:           signatures,
	}
	feed2New := datastreams.FeedReport{
		FeedID:               feedTwo,
		BenchmarkPrice:       big.NewInt(400).Bytes(),
		ObservationTimestamp: 400,
		FullReport:           []byte(rawReport2),
		ReportContext:        []byte{},
		Signatures:           signatures,
	}

	node1Resp, err := wrapReports([]datastreams.FeedReport{feed1Old, feed2New}, eventID, 400, datastreams.SignersMetadata{})
	require.NoError(t, err)
	rawNode1Resp, err := pb.MarshalCapabilityResponse(node1Resp)
	require.NoError(t, err)
	node2Resp, err := wrapReports([]datastreams.FeedReport{feed1New, feed2Old}, eventID, 300, datastreams.SignersMetadata{})
	require.NoError(t, err)
	rawNode2Resp, err := pb.MarshalCapabilityResponse(node2Resp)
	require.NoError(t, err)

	// aggregator should return latest value for each feedID
	aggResponse, err := agg.Aggregate(eventID, [][]byte{rawNode1Resp, rawNode2Resp})
	require.NoError(t, err)
	aggEvent := capabilities.TriggerEvent{}
	require.NoError(t, aggResponse.Value.UnwrapTo(&aggEvent))
	decodedReports, err := testMercuryCodec{}.UnwrapValid(aggEvent.Payload, nil, 0)
	require.NoError(t, err)

	require.Len(t, decodedReports, 2)
	require.Equal(t, feed1New, decodedReports[0])
	require.Equal(t, feed2New, decodedReports[1])
}
