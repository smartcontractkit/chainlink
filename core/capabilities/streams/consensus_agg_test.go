package streams_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/datafeeds"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Integration/load test that combines Data Feeds Consensus Aggregator and Streams Codec.
// For more meaningful measurements, increase the values of parameters P and T.
func TestStreamsConsensusAggregator(t *testing.T) {
	Nt := 10 // trigger DON nodes
	Ft := 3  // trigger DON faulty nodes
	Nw := 10 // workflow DON nodes
	Fw := 3  // workflow DON faulty nodes
	P := 40  // feeds
	T := 2   // test iterations

	triggerNodes := newNodes(t, Nt)
	feeds := newFeedsWithSignedReports(t, triggerNodes, Nt, P, 1)
	allowedSigners := make([][]byte, Nt)
	for i := 0; i < Nt; i++ {
		allowedSigners[i] = triggerNodes[i].bundle.PublicKey() // bad name - see comment on evmKeyring.PublicKey
	}

	config := newAggConfig(t, feeds)
	lggr := logger.TestLogger(t)
	codec := streams.NewCodec(lggr)
	agg, err := datafeeds.NewDataFeedsAggregator(*config, codec, lggr)
	require.NoError(t, err)

	// init round - empty previous Outcome, empty observations
	outcome, err := agg.Aggregate(nil, map[commontypes.OracleID][]values.Value{}, Fw)
	require.NoError(t, err)
	require.False(t, outcome.ShouldReport)

	// validate metadata
	newState := &datafeeds.DataFeedsOutcomeMetadata{}
	err = proto.Unmarshal(outcome.Metadata, newState)
	require.NoError(t, err)
	require.Equal(t, P, len(newState.FeedInfo))

	// run test aggregations
	startTs := time.Now().UnixMilli()
	processingTime := int64(0)
	for c := 0; c < T; c++ {
		obs := newObservations(t, Nw, feeds, Ft+1, allowedSigners)
		processingStart := time.Now().UnixMilli()
		outcome, err = agg.Aggregate(outcome, obs, Fw)
		processingTime += time.Now().UnixMilli() - processingStart
		require.NoError(t, err)
	}
	totalTime := time.Now().UnixMilli() - startTs
	lggr.Infow("elapsed", "totalMs", totalTime, "processingMs", processingTime)
}

func newAggConfig(t *testing.T, feeds []feed) *values.Map {
	feedMap := map[string]any{}
	for _, feed := range feeds {
		feedMap[feed.feedIDStr] = map[string]any{
			"deviation":  "0.1",
			"heartbeat":  1,
			"remappedID": feed.feedIDStr,
		}
	}
	unwrappedConfig := map[string]any{
		"feeds":                   feedMap,
		"allowedPartialStaleness": "0.2",
	}
	config, err := values.NewMap(unwrappedConfig)
	require.NoError(t, err)
	return config
}

func newObservations(t *testing.T, nNodes int, feeds []feed, minRequiredSignatures int, allowedSigners [][]byte) map[commontypes.OracleID][]values.Value {
	observations := map[commontypes.OracleID][]values.Value{}
	for i := 0; i < nNodes; i++ {
		reportList := []datastreams.FeedReport{}
		for j := 0; j < len(feeds); j++ {
			reportIdx := 0
			signedStreamsReport := datastreams.FeedReport{
				FeedID:        feeds[j].feedIDStr,
				FullReport:    feeds[j].reports[reportIdx].rawReport,
				ReportContext: feeds[j].reports[reportIdx].reportCtx,
				Signatures:    feeds[j].reports[reportIdx].signatures,
			}
			reportList = append(reportList, signedStreamsReport)
		}

		meta := datastreams.Metadata{
			Signers:               allowedSigners,
			MinRequiredSignatures: minRequiredSignatures,
		}
		p := datastreams.StreamsTriggerEvent{
			Payload:  reportList,
			Metadata: meta,
		}
		outputs, err := values.WrapMap(p)
		require.NoError(t, err)

		observations[commontypes.OracleID(i)] = []values.Value{outputs}
	}
	return observations
}
