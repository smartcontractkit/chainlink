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
	feedIDA                 = datastreams.FeedID("0x0001013ebd4ed3f5889fb5a8a52b42675c60c1a8c42bc79eaa72dcd922ac4292")
	deviationA              = decimal.NewFromFloat(0.1)
	heartbeatA              = 60
	feedIDB                 = datastreams.FeedID("0x0003c317fec7fad514c67aacc6366bf2f007ce37100e3cddcacd0ccaa1f3746d")
	deviationB              = decimal.NewFromFloat(0.01)
	heartbeatB              = 360
	mercuryFullReportA      = []byte("report")
	allowedPartialStaleness = 0.2
)

func TestDataFeedsAggregator_Aggregate_TwoRounds(t *testing.T) {
	metaVal, err := values.Wrap(datastreams.SignersMetadata{
		Signers:               [][]byte{newSigner(t), newSigner(t)},
		MinRequiredSignatures: 1,
	})
	require.NoError(t, err)
	mockTriggerEvent, err := values.Wrap(capabilities.TriggerEvent{
		Metadata: metaVal,
		Payload:  &values.Map{},
	})
	require.NoError(t, err)
	config := getConfig(t, feedIDA.String(), "0.1", heartbeatA)
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
	require.Equal(t, 2, len(newState.FeedInfo))
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
	outcome, err = agg.Aggregate(outcome, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent}, 2: {mockTriggerEvent}}, 1)
	require.NoError(t, err)
	require.True(t, outcome.ShouldReport)

	// validate metadata
	err = proto.Unmarshal(outcome.Metadata, newState)
	require.NoError(t, err)
	require.Equal(t, 2, len(newState.FeedInfo))
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

func TestDataFeedsAggregator_Aggregate_AllowedPartialStaleness(t *testing.T) {
	metaVal, err := values.Wrap(datastreams.SignersMetadata{
		Signers:               [][]byte{newSigner(t), newSigner(t)},
		MinRequiredSignatures: 1,
	})
	require.NoError(t, err)
	mockTriggerEvent, err := values.Wrap(capabilities.TriggerEvent{
		Metadata: metaVal,
		Payload:  &values.Map{},
	})
	require.NoError(t, err)
	config := getConfig(t, feedIDA.String(), "0.1", heartbeatA)
	codec := mocks.NewReportCodec(t)
	agg, err := datafeeds.NewDataFeedsAggregator(*config, codec, logger.Nop())
	require.NoError(t, err)

	// first round, both feeds are stale
	latestReportsRound1 := []datastreams.FeedReport{
		{
			FeedID:               feedIDA.String(),
			ObservationTimestamp: 1000,
			BenchmarkPrice:       big.NewInt(100).Bytes(),
		},
		{
			FeedID:               feedIDB.String(),
			ObservationTimestamp: 1100,
			BenchmarkPrice:       big.NewInt(200).Bytes(),
		},
	}
	codec.On("UnwrapValid", mock.Anything, mock.Anything, mock.Anything).Return(latestReportsRound1, nil).Twice()
	outcome, err := agg.Aggregate(nil, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent}, 2: {mockTriggerEvent}}, 1)
	require.NoError(t, err)
	require.True(t, outcome.ShouldReport)
	require.Equal(t, 2, len(outcome.EncodableOutcome.Fields[datafeeds.TopLevelListOutputFieldName].GetListValue().Fields))

	// second round, B hits deviation, A is not stale
	latestReportsRound2 := []datastreams.FeedReport{
		{
			FeedID:               feedIDA.String(),
			ObservationTimestamp: 1010,
			BenchmarkPrice:       big.NewInt(100).Bytes(),
		},
		{
			FeedID:               feedIDB.String(),
			ObservationTimestamp: 1110,
			BenchmarkPrice:       big.NewInt(400).Bytes(),
		},
	}
	codec.On("UnwrapValid", mock.Anything, mock.Anything, mock.Anything).Return(latestReportsRound2, nil).Twice()
	outcome, err = agg.Aggregate(outcome, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent}, 2: {mockTriggerEvent}}, 1)
	require.NoError(t, err)
	require.True(t, outcome.ShouldReport)
	require.Equal(t, 1, len(outcome.EncodableOutcome.Fields[datafeeds.TopLevelListOutputFieldName].GetListValue().Fields))

	// third round, B hits deviation, A is within allowed partial staleness threshold
	latestReportsRound3 := []datastreams.FeedReport{
		{
			FeedID:               feedIDA.String(),
			ObservationTimestamp: 1055,
			BenchmarkPrice:       big.NewInt(100).Bytes(),
		},
		{
			FeedID:               feedIDB.String(),
			ObservationTimestamp: 1150,
			BenchmarkPrice:       big.NewInt(600).Bytes(),
		},
	}
	codec.On("UnwrapValid", mock.Anything, mock.Anything, mock.Anything).Return(latestReportsRound3, nil).Twice()
	outcome, err = agg.Aggregate(outcome, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent}, 2: {mockTriggerEvent}}, 1)
	require.NoError(t, err)
	require.True(t, outcome.ShouldReport)
	require.Equal(t, 2, len(outcome.EncodableOutcome.Fields[datafeeds.TopLevelListOutputFieldName].GetListValue().Fields))
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

	config := getConfig(t, feedIDA.String(), "0.1", heartbeatA)
	codec := mocks.NewReportCodec(t)
	agg, err := datafeeds.NewDataFeedsAggregator(*config, codec, logger.Nop())
	require.NoError(t, err)

	// no valid signers - each one should appear at least twice to be valid
	_, err = agg.Aggregate(nil, map[commontypes.OracleID][]values.Value{1: {mockTriggerEvent}}, 1)
	require.Error(t, err)
}

func TestDataFeedsAggregator_ParseConfig(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		config := getConfig(t, feedIDA.String(), "0.1", heartbeatA)
		parsedConfig, err := datafeeds.ParseConfig(*config)
		require.NoError(t, err)
		require.Equal(t, deviationA, parsedConfig.Feeds[feedIDA].Deviation)
		require.Equal(t, heartbeatA, parsedConfig.Feeds[feedIDA].Heartbeat)
		require.Equal(t, deviationB, parsedConfig.Feeds[feedIDB].Deviation)
		require.Equal(t, heartbeatB, parsedConfig.Feeds[feedIDB].Heartbeat)
		require.Equal(t, allowedPartialStaleness, parsedConfig.AllowedPartialStaleness)
	})

	t.Run("invalid ID", func(t *testing.T) {
		config := getConfig(t, "bad_id", "0.1", heartbeatA)
		_, err := datafeeds.ParseConfig(*config)
		require.Error(t, err)
	})

	t.Run("invalid deviation string", func(t *testing.T) {
		config := getConfig(t, feedIDA.String(), "bad_number", heartbeatA)
		_, err := datafeeds.ParseConfig(*config)
		require.Error(t, err)
	})
}

func getConfig(t *testing.T, feedID string, deviation string, heartbeat int) *values.Map {
	unwrappedConfig := map[string]any{
		"feeds": map[string]any{
			feedID: map[string]any{
				"deviation": deviation,
				"heartbeat": heartbeat,
			},
			feedIDB.String(): map[string]any{
				"deviation": deviationB.String(),
				"heartbeat": heartbeatB,
			},
		},
		"allowedPartialStaleness": "0.2",
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
