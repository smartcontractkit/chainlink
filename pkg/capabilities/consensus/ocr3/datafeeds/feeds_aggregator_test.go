package datafeeds_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/datafeeds"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

var (
	feedIDA            = datastreams.FeedID("0x0001013ebd4ed3f5889fb5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	deviationA         = decimal.NewFromFloat(0.1)
	heartbeatA         = 60
	mercuryFullReportA = []byte("report")
)

func TestDataFeedsAggregator_Aggregate_TwoRounds(t *testing.T) {
	meta := datastreams.SignersMetadata{
		Signers:               [][]byte{newSigner(t), newSigner(t)},
		MinRequiredSignatures: 1,
	}
	metaVal, err := values.Wrap(meta)
	require.NoError(t, err)
	mockTriggerEvent1, err := values.Wrap(capabilities.TriggerEvent{
		Metadata: metaVal,
		Payload:  &values.Map{},
	})
	require.NoError(t, err)
	mockTriggerEvent2, err := values.Wrap(capabilities.TriggerEvent{
		Metadata: metaVal,
		Payload:  &values.Map{},
	})
	require.NoError(t, err)
	config := getConfig(t, feedIDA.String(), deviationA, heartbeatA)
	codec := mocks.NewReportCodec(t)
	agg, err := datafeeds.NewDataFeedsAggregator(*config, codec, logger.Nop())
	require.NoError(t, err)

	// first round, empty previous Outcome, empty observations
	outcome, err := agg.Aggregate(nil, map[commontypes.OracleID][]values.Value{}, 1)
	require.NoError(t, err)
	require.False(t, outcome.ShouldReport)

	// validate metadata
	newState := &datafeeds.DataFeedsOutcomeMetadata{}
	err = proto.Unmarshal(outcome.Metadata, newState)
	require.NoError(t, err)
	require.Equal(t, 1, len(newState.FeedInfo))
	_, ok := newState.FeedInfo[feedIDA.String()]
	require.True(t, ok)
	require.Equal(t, []byte(nil), newState.FeedInfo[feedIDA.String()].BenchmarkPrice)

	// second round, non-empty previous Outcome, one observation
	latestMercuryReports := []datastreams.FeedReport{
		{
			FeedID:               feedIDA.String(),
			ObservationTimestamp: 1,
			BenchmarkPrice:       big.NewInt(100).Bytes(),
			FullReport:           mercuryFullReportA,
		},
	}
	codec.On("UnwrapValid", mock.Anything, mock.Anything, mock.Anything).Return(latestMercuryReports, nil)
	outcome, err = agg.Aggregate(outcome, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent1}, 2: {mockTriggerEvent2}}, 1)
	require.NoError(t, err)
	require.True(t, outcome.ShouldReport)

	// validate metadata
	err = proto.Unmarshal(outcome.Metadata, newState)
	require.NoError(t, err)
	require.Equal(t, 1, len(newState.FeedInfo))
	_, ok = newState.FeedInfo[feedIDA.String()]
	require.True(t, ok)
	require.Equal(t, big.NewInt(100).Bytes(), newState.FeedInfo[feedIDA.String()].BenchmarkPrice)

	// validate encodable outcome
	val := values.FromMapValueProto(outcome.EncodableOutcome)
	require.NoError(t, err)
	topLevelMap, err := val.Unwrap()
	require.NoError(t, err)
	mm, ok := topLevelMap.(map[string]any)
	require.True(t, ok)

	idBytes := feedIDA.Bytes()
	expected := map[string]any{
		datafeeds.TopLevelListOutputFieldName: []any{
			map[string]any{
				datafeeds.FeedIDOutputFieldName:    idBytes[:],
				datafeeds.RawReportOutputFieldName: mercuryFullReportA,
				datafeeds.TimestampOutputFieldName: int64(1),
				datafeeds.PriceOutputFieldName:     big.NewInt(100),
			},
		},
	}
	require.Equal(t, expected, mm)
}

func TestDataFeedsAggregator_Aggregate_Failures(t *testing.T) {
	meta := datastreams.SignersMetadata{
		Signers:               [][]byte{newSigner(t), newSigner(t)},
		MinRequiredSignatures: 1,
	}
	metaVal, err := values.Wrap(meta)
	require.NoError(t, err)
	mockTriggerEvent, err := values.Wrap(capabilities.TriggerEvent{
		Metadata: metaVal,
		Payload:  &values.Map{},
	})
	require.NoError(t, err)

	config := getConfig(t, feedIDA.String(), deviationA, heartbeatA)
	codec := mocks.NewReportCodec(t)
	agg, err := datafeeds.NewDataFeedsAggregator(*config, codec, logger.Nop())
	require.NoError(t, err)

	// no valid signers - each one should appear at least twice to be valid
	_, err = agg.Aggregate(nil, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent}}, 1)
	require.Error(t, err)
}

func TestDataFeedsAggregator_ParseConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		config := getConfig(t, feedIDA.String(), deviationA, heartbeatA)
		parsedConfig, err := datafeeds.ParseConfig(*config)
		require.NoError(t, err)
		require.Equal(t, deviationA, parsedConfig.Feeds[feedIDA].Deviation)
		require.Equal(t, heartbeatA, parsedConfig.Feeds[feedIDA].Heartbeat)
	})

	t.Run("invalid ID", func(t *testing.T) {
		config := getConfig(t, "bad_id", deviationA, heartbeatA)
		parsedConfig, err := datafeeds.ParseConfig(*config)
		require.NoError(t, err)
		require.Equal(t, deviationA, parsedConfig.Feeds[feedIDA].Deviation)
		require.Equal(t, heartbeatA, parsedConfig.Feeds[feedIDA].Heartbeat)
	})

	t.Run("parsed workflow config", func(t *testing.T) {
		fdID := datastreams.FeedID("0x1111111111111111111100000000000000000000000000000000000000000000")
		cfg, err := values.NewMap(map[string]any{
			fdID.String(): map[string]any{
				"deviation": "0.1",
				"heartbeat": 60,
			},
		})
		require.NoError(t, err)
		parsedConfig, err := datafeeds.ParseConfig(*cfg)
		require.NoError(t, err)
		require.Equal(t, deviationA, parsedConfig.Feeds[fdID].Deviation, parsedConfig)
		require.Equal(t, heartbeatA, parsedConfig.Feeds[fdID].Heartbeat)
	})
}

func getConfig(t *testing.T, feedID string, deviation decimal.Decimal, heartbeat int) *values.Map {
	unwrappedConfig := map[string]any{
		feedIDA.String(): map[string]any{
			"deviation": deviation.String(),
			"heartbeat": heartbeat,
		},
	}
	config, err := values.NewMap(unwrappedConfig)
	require.NoError(t, err)
	return config
}

func newSigner(t *testing.T) []byte {
	buf := make([]byte, 20)
	_, err := rand.Read(buf)
	require.NoError(t, err)
	return buf
}
