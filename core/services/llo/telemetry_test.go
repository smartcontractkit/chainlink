package llo

import (
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/eautils"
	legacytelem "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

var _ commontypes.MonitoringEndpoint = &mockMonitoringEndpoint{}

type mockMonitoringEndpoint struct {
	chLogs chan []byte
}

func (m *mockMonitoringEndpoint) SendLog(log []byte) {
	m.chLogs <- log
}

const bridgeResponse = `{
			"meta":{
				"adapterName":"data-source-name"
			},
			"timestamps":{
				"providerDataRequestedUnixMs":92233720368547760,
				"providerDataReceivedUnixMs":-92233720368547760,
				"providerDataStreamEstablishedUnixMs":1,
				"providerIndicatedTimeUnixMs":-123456789
			}
		}`

var trrs = pipeline.TaskRunResults{
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "test-bridge-1",
			BaseTask:    pipeline.NewBaseTask(0, "ds1", nil, nil, 0),
			RequestData: `{"data":{"from":"eth", "to":"usd"}}`,
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
		CreatedAt:  time.Unix(0, 0),
		FinishedAt: null.TimeFrom(time.Unix(0, 0)),
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds1_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "123456.123456789",
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "test-bridge-2",
			BaseTask:    pipeline.NewBaseTask(0, "ds2", nil, nil, 0),
			RequestData: `{"data":{"from":"eth", "to":"usd"}}`,
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
		CreatedAt:  time.Unix(1, 0),
		FinishedAt: null.TimeFrom(time.Unix(10, 0)),
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds2_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "12345678",
		},
	},
	pipeline.TaskRunResult{
		Task: &pipeline.BridgeTask{
			Name:        "test-bridge-3",
			BaseTask:    pipeline.NewBaseTask(0, "ds3", nil, nil, 0),
			RequestData: `{"data":{"from":"eth", "to":"usd"}}`,
		},
		Result: pipeline.Result{
			Value: bridgeResponse,
		},
		CreatedAt:  time.Unix(2, 0),
		FinishedAt: null.TimeFrom(time.Unix(20, 0)),
	},
	pipeline.TaskRunResult{
		Task: &pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "ds3_parse", nil, nil, 1),
		},
		Result: pipeline.Result{
			Value: "1234567890",
		},
	},
}

func Test_Telemeter_v3PremiumLegacy(t *testing.T) {
	lggr := logger.TestLogger(t)
	m := &mockMonitoringEndpoint{}

	run := &pipeline.Run{ID: 42}
	streamID := uint32(135)
	donID := uint32(1)
	opts := &mockOpts{}

	t.Run("with error", func(t *testing.T) {
		tm := newTelemeter(lggr, m, donID)
		servicetest.Run(t, tm)

		t.Run("if error is some random failure returns immediately", func(t *testing.T) {
			// should return immediately and not even send on the channel
			m.chLogs = nil
			tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, nil, errors.New("test error"))
		})
		t.Run("if error is dp invariant violation, sets this flag", func(t *testing.T) {
			m.chLogs = make(chan []byte, 100)
			adapterError := new(eautils.AdapterError)
			adapterError.Name = adapterLWBAErrorName
			tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, nil, adapterError)

			var i int
			for log := range m.chLogs {
				decoded := &legacytelem.EnhancedEAMercury{}
				require.NoError(t, proto.Unmarshal(log, decoded))
				assert.True(t, decoded.DpInvariantViolationDetected)
				if i == 2 {
					return
				}
				i++
			}
		})
	})
	t.Run("with decimal value, sets all values correctly", func(t *testing.T) {
		tm := newTelemeter(lggr, m, donID)
		val := llo.ToDecimal(decimal.NewFromFloat32(102.12))
		servicetest.Run(t, tm)
		tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, nil)

		var i int
		for log := range m.chLogs {
			decoded := &legacytelem.EnhancedEAMercury{}
			require.NoError(t, proto.Unmarshal(log, decoded))
			assert.Equal(t, int(1003), int(decoded.Version))
			assert.InDelta(t, float64(123456.123456789), decoded.DpBenchmarkPrice, 0.0000000001)
			assert.Zero(t, decoded.DpBid)
			assert.Zero(t, decoded.DpAsk)
			assert.False(t, decoded.DpInvariantViolationDetected)
			assert.Zero(t, decoded.CurrentBlockNumber)
			assert.Zero(t, decoded.CurrentBlockHash)
			assert.Zero(t, decoded.CurrentBlockTimestamp)
			assert.Zero(t, decoded.FetchMaxFinalizedTimestamp)
			assert.Zero(t, decoded.MaxFinalizedTimestamp)
			assert.Zero(t, decoded.ObservationTimestamp)
			assert.False(t, decoded.IsLinkFeed)
			assert.Zero(t, decoded.LinkPrice)
			assert.False(t, decoded.IsNativeFeed)
			assert.Zero(t, decoded.NativePrice)
			assert.Equal(t, int64(i*1000), decoded.BridgeTaskRunStartedTimestamp)
			assert.Equal(t, int64(i*10000), decoded.BridgeTaskRunEndedTimestamp)
			assert.Equal(t, int64(92233720368547760), decoded.ProviderRequestedTimestamp)
			assert.Equal(t, int64(-92233720368547760), decoded.ProviderReceivedTimestamp)
			assert.Equal(t, int64(1), decoded.ProviderDataStreamEstablished)
			assert.Equal(t, int64(-123456789), decoded.ProviderIndicatedTime)
			assert.Equal(t, "streamID:135", decoded.Feed)
			assert.Equal(t, int64(102), decoded.ObservationBenchmarkPrice)
			assert.Equal(t, "102.12", decoded.ObservationBenchmarkPriceString)
			assert.Zero(t, decoded.ObservationBid)
			assert.Zero(t, decoded.ObservationBidString)
			assert.Zero(t, decoded.ObservationAsk)
			assert.Zero(t, decoded.ObservationAskString)
			assert.Zero(t, decoded.ObservationMarketStatus)
			assert.Equal(t, "0605040000000000000000000000000000000000000000000000000000000000", decoded.ConfigDigest)
			assert.Equal(t, int64(18), decoded.Round)
			assert.Equal(t, int64(4), decoded.Epoch)
			assert.Equal(t, "eth/usd", decoded.AssetSymbol)
			assert.Equal(t, uint32(1), decoded.DonId)
			if i == 2 {
				return
			}
			i++
		}
	})
	t.Run("with quote value", func(t *testing.T) {
		tm := newTelemeter(lggr, m, donID)
		val := &llo.Quote{Bid: decimal.NewFromFloat32(102.12), Benchmark: decimal.NewFromFloat32(103.32), Ask: decimal.NewFromFloat32(104.25)}
		servicetest.Run(t, tm)
		tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, nil)

		var i int
		for log := range m.chLogs {
			decoded := &legacytelem.EnhancedEAMercury{}
			require.NoError(t, proto.Unmarshal(log, decoded))
			assert.Equal(t, int64(103), decoded.ObservationBenchmarkPrice)
			assert.Equal(t, "103.32", decoded.ObservationBenchmarkPriceString)
			assert.Equal(t, int64(102), decoded.ObservationBid)
			assert.Equal(t, "102.12", decoded.ObservationBidString)
			assert.Equal(t, int64(104), decoded.ObservationAsk)
			assert.Equal(t, "104.25", decoded.ObservationAskString)
			assert.Zero(t, decoded.ObservationMarketStatus)
			if i == 2 {
				return
			}
			i++
		}
	})
}

func Test_Telemeter_observationTelemetry(t *testing.T) {
	lggr := logger.TestLogger(t)

	donID := uint32(1)

	opts := &mockOpts{}

	t.Run("transmits *pipeline.BridgeTelemetry", func(t *testing.T) {
		t.Parallel()
		m := &mockMonitoringEndpoint{chLogs: make(chan []byte, 100)}
		tm := newTelemeter(lggr, m, donID)
		servicetest.Run(t, tm)
		ch := tm.MakeTelemChannel(opts, 100)

		ch <- &pipeline.BridgeTelemetry{
			Name:                   "test-bridge-1",
			RequestData:            []byte(`foo`),
			ResponseData:           []byte(`bar`),
			ResponseError:          ptr("test error"),
			ResponseStatusCode:     200,
			RequestStartTimestamp:  time.Unix(1, 1),
			RequestFinishTimestamp: time.Unix(2, 1),
			LocalCacheHit:          true,
			SpecID:                 3,
			StreamID:               ptr(uint32(135)),
			DotID:                  "ds1",
		}

		log := <-m.chLogs
		decoded := &telem.LLOBridgeTelemetry{}
		require.NoError(t, proto.Unmarshal(log, decoded))
		assert.Equal(t, "test-bridge-1", decoded.BridgeAdapterName)
		assert.Equal(t, []byte(`foo`), decoded.BridgeRequestData)
		assert.Equal(t, []byte(`bar`), decoded.BridgeResponseData)
		require.NotNil(t, decoded.BridgeResponseError)
		assert.Equal(t, "test error", *decoded.BridgeResponseError)
		assert.Equal(t, int32(200), decoded.BridgeResponseStatusCode)
		assert.Equal(t, int64(1000000001), decoded.RequestStartTimestamp)
		assert.Equal(t, int64(2000000001), decoded.RequestFinishTimestamp)
		assert.True(t, decoded.LocalCacheHit)
		assert.Equal(t, int32(3), decoded.SpecId)
		require.NotNil(t, decoded.StreamId)
		assert.Equal(t, uint32(135), *decoded.StreamId)
		assert.Equal(t, "ds1", decoded.DotId)

		// added by telemeter
		assert.Equal(t, donID, decoded.DonId)
		assert.Equal(t, opts.SeqNr(), decoded.SeqNr)
		assert.Equal(t, opts.ConfigDigest().Hex(), hex.EncodeToString(decoded.ConfigDigest))
		assert.Equal(t, opts.ObservationTimestamp().UnixNano(), decoded.ObservationTimestamp)
	})
	t.Run("transmits *telem.LLOObservationTelemetry", func(t *testing.T) {
		t.Parallel()
		m := &mockMonitoringEndpoint{chLogs: make(chan []byte, 100)}
		tm := newTelemeter(lggr, m, donID)
		servicetest.Run(t, tm)
		ch := tm.MakeTelemChannel(opts, 100)

		ch <- &telem.LLOObservationTelemetry{
			StreamId:              135,
			StreamValueType:       1,
			StreamValueBinary:     []byte{0x01, 0x02, 0x03},
			StreamValueText:       "stream value text",
			ObservationError:      ptr("test error"),
			ObservationTimestamp:  time.Unix(1, 1).UnixNano(),
			ObservationFinishedAt: time.Unix(2, 1).UnixNano(),
			SeqNr:                 42,
			ConfigDigest:          []byte{0x01, 0x02, 0x03},
		}

		log := <-m.chLogs
		decoded := &telem.LLOObservationTelemetry{}
		require.NoError(t, proto.Unmarshal(log, decoded))
		assert.Equal(t, uint32(135), decoded.StreamId)
		assert.Equal(t, int32(1), decoded.StreamValueType)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, decoded.StreamValueBinary)
		assert.Equal(t, "stream value text", decoded.StreamValueText)
		require.NotNil(t, decoded.ObservationError)
		assert.Equal(t, "test error", *decoded.ObservationError)
		assert.Equal(t, int64(1000000001), decoded.ObservationTimestamp)
		assert.Equal(t, int64(2000000001), decoded.ObservationFinishedAt)
		assert.Equal(t, uint64(42), decoded.SeqNr)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, decoded.ConfigDigest)

		// telemeter adds don ID
		assert.Equal(t, donID, decoded.DonId)
	})

	t.Run("ignores unknown telemetry type", func(t *testing.T) {
		t.Parallel()
		m := &mockMonitoringEndpoint{chLogs: make(chan []byte, 100)}
		obsLggr, observedLogs := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		tm := newTelemeter(obsLggr, m, donID)
		servicetest.Run(t, tm)
		ch := tm.MakeTelemChannel(opts, 100)

		ch <- struct{}{}

		testutils.WaitForLogMessage(t, observedLogs, "Unknown telemetry type")
	})
}

func ptr[T any](t T) *T { return &t }
