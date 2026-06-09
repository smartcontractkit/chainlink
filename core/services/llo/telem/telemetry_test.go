package telem

import (
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/eautils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	legacytelem "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

var _ telemetry.MultitypeMonitoringEndpoint = &mockMonitoringEndpoint{}

type typedLog struct {
	log       []byte
	telemType synchronization.TelemetryType
}

type mockMonitoringEndpoint struct {
	chTypedLogs chan typedLog
}

func (m *mockMonitoringEndpoint) SendTypedLog(telemType synchronization.TelemetryType, log []byte) {
	m.chTypedLogs <- typedLog{log, telemType}
}

type mockOpts struct {
	verboseLogging bool
}

func (m *mockOpts) VerboseLogging() bool { return m.verboseLogging }
func (m *mockOpts) SeqNr() uint64        { return 1042 }
func (m *mockOpts) OutCtx() ocr3types.OutcomeContext {
	return ocr3types.OutcomeContext{SeqNr: 1042, PreviousOutcome: ocr3types.Outcome([]byte("foo"))}
}
func (m *mockOpts) ConfigDigest() ocr2types.ConfigDigest {
	return ocr2types.ConfigDigest{6, 5, 4}
}
func (m *mockOpts) ObservationTimestamp() time.Time {
	return time.Unix(1737936858, 0)
}

func (m *mockOpts) OutcomeCodec() llo.OutcomeCodec {
	return nil
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
		tm := newTelemeter(TelemeterParams{
			Logger:             lggr,
			MonitoringEndpoint: m,
			DonID:              donID,
		})
		servicetest.Run(t, tm)

		t.Run("if error is some random failure returns immediately", func(t *testing.T) {
			// should return immediately and not even send on the channel
			m.chTypedLogs = nil
			tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, nil, errors.New("test error"))
		})
		t.Run("if error is dp invariant violation, sets this flag", func(t *testing.T) {
			m.chTypedLogs = make(chan typedLog, 100)
			adapterError := new(eautils.AdapterError)
			adapterError.Name = adapterLWBAErrorName
			tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, nil, adapterError)
			tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

			var i int
			for tLog := range m.chTypedLogs {
				assert.Equal(t, synchronization.EnhancedEAMercury, tLog.telemType)
				decoded := &legacytelem.EnhancedEAMercury{}
				require.NoError(t, proto.Unmarshal(tLog.log, decoded))
				assert.True(t, decoded.DpInvariantViolationDetected)
				if i == 2 {
					return
				}
				i++
			}
		})
	})
	t.Run("with decimal value, sets all values correctly", func(t *testing.T) {
		tm := newTelemeter(TelemeterParams{
			Logger:             lggr,
			MonitoringEndpoint: m,
			DonID:              donID,
		})
		val := llo.ToDecimal(decimal.NewFromFloat32(102.12))
		servicetest.Run(t, tm)
		tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, nil)
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		var i int
		for tLog := range m.chTypedLogs {
			assert.Equal(t, synchronization.EnhancedEAMercury, tLog.telemType)
			decoded := &legacytelem.EnhancedEAMercury{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
			assert.Equal(t, int(1003), int(decoded.Version))
			assert.InDelta(t, float64(123456.123456789), decoded.DpBenchmarkPrice, 0.0000000001)
			assert.Zero(t, decoded.DpBid)
			assert.Zero(t, decoded.DpAsk)
			assert.False(t, decoded.DpInvariantViolationDetected)
			assert.Zero(t, decoded.CurrentBlockNumber)
			assert.Empty(t, decoded.CurrentBlockHash)
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
			assert.Empty(t, decoded.ObservationBidString)
			assert.Zero(t, decoded.ObservationAsk)
			assert.Empty(t, decoded.ObservationAskString)
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
		tm := newTelemeter(TelemeterParams{
			Logger:             lggr,
			MonitoringEndpoint: m,
			DonID:              donID,
		})
		val := &llo.Quote{Bid: decimal.NewFromFloat32(102.12), Benchmark: decimal.NewFromFloat32(103.32), Ask: decimal.NewFromFloat32(104.25)}
		servicetest.Run(t, tm)
		tm.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, nil)
		time.Sleep(10 * time.Millisecond)
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		var i int
		for tLog := range m.chTypedLogs {
			assert.Equal(t, synchronization.EnhancedEAMercury, tLog.telemType)
			decoded := &legacytelem.EnhancedEAMercury{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
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

func Test_Telemeter_observationScopedTelemetry(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	donID := uint32(1)
	opts := &mockOpts{}

	t.Run("if both CaptureEATelemetry and CaptureObservationTelemetry are false, returns nil channel", func(t *testing.T) {
		tm := newTelemeter(TelemeterParams{
			Logger: lggr,
			DonID:  donID,
		})
		ch := tm.MakeObservationScopedTelemetryCh(opts, 100)
		assert.Nil(t, ch)
	})

	t.Run("transmits *pipeline.BridgeTelemetry", func(t *testing.T) {
		t.Parallel()
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:             lggr,
			MonitoringEndpoint: m,
			DonID:              donID,
			CaptureEATelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.MakeObservationScopedTelemetryCh(opts, 100)
		require.NotNil(t, ch)

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
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		tLog := <-m.chTypedLogs
		assert.Equal(t, synchronization.PipelineBridge, tLog.telemType)
		decoded := &LLOBridgeTelemetry{}
		require.NoError(t, proto.Unmarshal(tLog.log, decoded))
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
	t.Run("transmits *LLOObservationTelemetry", func(t *testing.T) {
		t.Parallel()
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                      lggr,
			MonitoringEndpoint:          m,
			DonID:                       donID,
			CaptureObservationTelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.MakeObservationScopedTelemetryCh(opts, 100)
		require.NotNil(t, ch)

		ch <- &LLOObservationTelemetry{
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
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		tLog := <-m.chTypedLogs
		assert.Equal(t, synchronization.LLOObservation, tLog.telemType)
		decoded := &LLOObservationTelemetry{}
		require.NoError(t, proto.Unmarshal(tLog.log, decoded))
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
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		obsLggr, observedLogs := logger.TestLoggerObserved(t, zapcore.WarnLevel)
		tm := newTelemeter(TelemeterParams{
			Logger:                      obsLggr,
			MonitoringEndpoint:          m,
			DonID:                       donID,
			CaptureEATelemetry:          true,
			CaptureObservationTelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.MakeObservationScopedTelemetryCh(opts, 100)
		require.NotNil(t, ch)

		ch <- struct{}{}
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		testutils.WaitForLogMessage(t, observedLogs, "Unknown telemetry type")
	})
}

func Test_Telemeter_outcomeTelemetry(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	donID := uint32(1)

	t.Run("returns nil channel if CaptureOutcomeTelemetry is false", func(t *testing.T) {
		tm := newTelemeter(TelemeterParams{
			Logger: lggr,
			DonID:  donID,
		})
		ch := tm.GetOutcomeTelemetryCh()
		assert.Nil(t, ch)
	})

	t.Run("transmits *datastreamsllo.LLOOutcomeTelemetry", func(t *testing.T) {
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                  lggr,
			MonitoringEndpoint:      m,
			DonID:                   donID,
			CaptureOutcomeTelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetOutcomeTelemetryCh()
		require.NotNil(t, ch)

		t.Run("zero values", func(t *testing.T) {
			opts := &mockOpts{}
			cd := opts.ConfigDigest()
			orig := &datastreamsllo.LLOOutcomeTelemetry{SeqNr: opts.SeqNr(), ConfigDigest: cd[:]}
			ch <- orig

			// Wait until the telemetry is buffered.
			testutils.RequireEventually(t, func() bool {
				tm.telemetryBufferMu.Lock()
				defer tm.telemetryBufferMu.Unlock()
				return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) > 0
			})

			tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

			tLog := <-m.chTypedLogs
			assert.Equal(t, synchronization.LLOOutcome, tLog.telemType)
			decoded := &datastreamsllo.LLOOutcomeTelemetry{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
			assert.Empty(t, decoded.LifeCycleStage)
			assert.Zero(t, decoded.ObservationTimestampNanoseconds)
			assert.Zero(t, decoded.ChannelDefinitions)
			assert.Zero(t, decoded.ValidAfterNanoseconds)
			assert.Zero(t, decoded.StreamAggregates)
			assert.Equal(t, opts.SeqNr(), decoded.SeqNr)
			assert.Equal(t, cd[:], decoded.ConfigDigest)
			assert.Zero(t, decoded.DonId)
		})

		t.Run("with values", func(t *testing.T) {
			opts := &mockOpts{}
			cd := opts.ConfigDigest()
			orig := &datastreamsllo.LLOOutcomeTelemetry{
				LifeCycleStage:                  "foo",
				ObservationTimestampNanoseconds: 2,
				ChannelDefinitions: map[uint32]*datastreamsllo.LLOChannelDefinitionProto{
					3: {
						ReportFormat: 4,
						Streams: []*datastreamsllo.LLOStreamDefinition{
							{
								StreamID:   5,
								Aggregator: 6,
							},
						},
						Opts: []byte{7},
					},
				},
				ValidAfterNanoseconds: map[uint32]uint64{
					8: 9,
				},
				StreamAggregates: map[uint32]*datastreamsllo.LLOAggregatorStreamValue{
					10: {
						AggregatorValues: map[uint32]*datastreamsllo.LLOStreamValue{
							11: {
								Type:  12,
								Value: []byte{13},
							},
						},
					},
				},
				SeqNr:        opts.SeqNr(),
				ConfigDigest: cd[:],
				DonId:        10,
			}
			ch <- orig

			// Wait until the telemetry is buffered.
			testutils.RequireEventually(t, func() bool {
				tm.telemetryBufferMu.Lock()
				defer tm.telemetryBufferMu.Unlock()
				return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) > 0
			})

			tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

			tLog := <-m.chTypedLogs
			assert.Equal(t, synchronization.LLOOutcome, tLog.telemType)
			decoded := &datastreamsllo.LLOOutcomeTelemetry{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
			assert.Equal(t, "foo", decoded.LifeCycleStage)
			assert.Equal(t, uint64(2), decoded.ObservationTimestampNanoseconds)
			assert.Len(t, decoded.ChannelDefinitions, 1)
			assert.Equal(t, uint32(4), decoded.ChannelDefinitions[3].ReportFormat)
			assert.Len(t, decoded.ChannelDefinitions[3].Streams, 1)
			assert.Equal(t, uint32(5), decoded.ChannelDefinitions[3].Streams[0].StreamID)
			assert.Equal(t, uint32(6), decoded.ChannelDefinitions[3].Streams[0].Aggregator)
			assert.Equal(t, []byte{7}, decoded.ChannelDefinitions[3].Opts)
			assert.Len(t, decoded.ValidAfterNanoseconds, 1)
			assert.Equal(t, uint64(9), decoded.ValidAfterNanoseconds[8])
			assert.Len(t, decoded.StreamAggregates, 1)
			assert.Len(t, decoded.StreamAggregates[10].AggregatorValues, 1)
			assert.Equal(t, llo.LLOStreamValue_Type(12), decoded.StreamAggregates[10].AggregatorValues[11].Type)
			assert.Equal(t, []byte{13}, decoded.StreamAggregates[10].AggregatorValues[11].Value)
			assert.Equal(t, opts.SeqNr(), decoded.SeqNr)
			assert.Equal(t, cd[:], decoded.ConfigDigest)
			assert.Equal(t, uint32(10), decoded.DonId)
		})
	})

	t.Run("overwrites previous outcome for same seqNr (epoch transition)", func(t *testing.T) {
		// Simulates the bug scenario: Outcome() is called twice for the
		// same seqNr across two epochs (first epoch fails prepare-quorum,
		// second commits). Only the last outcome should survive in the
		// buffer and be flushed.
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                  lggr,
			MonitoringEndpoint:      m,
			DonID:                   donID,
			CaptureOutcomeTelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetOutcomeTelemetryCh()
		require.NotNil(t, ch)

		opts := &mockOpts{}
		cd := opts.ConfigDigest()

		// First outcome (from failed epoch)
		epoch1Outcome := &datastreamsllo.LLOOutcomeTelemetry{
			LifeCycleStage:                  "production",
			ObservationTimestampNanoseconds: 1000000001,
			SeqNr:                           opts.SeqNr(),
			ConfigDigest:                    cd[:],
			DonId:                           10,
		}
		ch <- epoch1Outcome

		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) > 0
		})

		// Second outcome (from committed epoch) — different observation timestamp
		epoch2Outcome := &datastreamsllo.LLOOutcomeTelemetry{
			LifeCycleStage:                  "production",
			ObservationTimestampNanoseconds: 2000000002,
			SeqNr:                           opts.SeqNr(),
			ConfigDigest:                    cd[:],
			DonId:                           10,
		}
		ch <- epoch2Outcome

		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			buf := tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]
			if len(buf) == 0 {
				return false
			}
			// Wait until the buffer contains the second outcome
			msg := buf[0].msg.(*datastreamsllo.LLOOutcomeTelemetry)
			return msg.ObservationTimestampNanoseconds == 2000000002
		})

		// Verify buffer has exactly 1 entry (overwritten, not appended)
		tm.telemetryBufferMu.Lock()
		bufLen := len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()])
		tm.telemetryBufferMu.Unlock()
		assert.Equal(t, 1, bufLen, "expected exactly 1 outcome in buffer (overwrite), got %d", bufLen)

		// Flush and verify only one message is sent
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		tLog := <-m.chTypedLogs
		assert.Equal(t, synchronization.LLOOutcome, tLog.telemType)
		decoded := &datastreamsllo.LLOOutcomeTelemetry{}
		require.NoError(t, proto.Unmarshal(tLog.log, decoded))

		// The flushed outcome should be from the second (committed) epoch
		assert.Equal(t, uint64(2000000002), decoded.ObservationTimestampNanoseconds,
			"expected observation timestamp from committed epoch")
		assert.Equal(t, opts.SeqNr(), decoded.SeqNr)

		// Verify no additional messages were sent
		select {
		case extra := <-m.chTypedLogs:
			t.Fatalf("expected no more outcome messages, but got one with type %v", extra.telemType)
		default:
			// good — no extra messages
		}
	})
}

func Test_Telemeter_reportTelemetry(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	donID := uint32(1)

	t.Run("returns nil channel if CaptureReportTelemetry is false", func(t *testing.T) {
		tm := newTelemeter(TelemeterParams{
			Logger: lggr,
			DonID:  donID,
		})
		ch := tm.GetReportTelemetryCh()
		assert.Nil(t, ch)
	})

	t.Run("transmits *datastreamsllo.LLOReportTelemetry", func(t *testing.T) {
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                 lggr,
			MonitoringEndpoint:     m,
			DonID:                  donID,
			CaptureReportTelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetReportTelemetryCh()
		require.NotNil(t, ch)

		t.Run("zero values", func(t *testing.T) {
			opts := &mockOpts{}
			cd := opts.ConfigDigest()
			orig := &datastreamsllo.LLOReportTelemetry{SeqNr: opts.SeqNr(), ConfigDigest: cd[:]}
			ch <- orig

			// Wait until the telemetry is buffered.
			testutils.RequireEventually(t, func() bool {
				tm.telemetryBufferMu.Lock()
				defer tm.telemetryBufferMu.Unlock()
				return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) > 0
			})

			tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

			tLog := <-m.chTypedLogs
			assert.Equal(t, synchronization.LLOReport, tLog.telemType)
			decoded := &datastreamsllo.LLOReportTelemetry{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
			assert.Zero(t, decoded.ChannelId)
			assert.Zero(t, decoded.ValidAfterNanoseconds)
			assert.Zero(t, decoded.ObservationTimestampNanoseconds)
			assert.Zero(t, decoded.ReportFormat)
			assert.False(t, decoded.Specimen)
			assert.Empty(t, decoded.StreamDefinitions)
			assert.Empty(t, decoded.StreamValues)
			assert.Empty(t, decoded.ChannelOpts)
			assert.Equal(t, opts.SeqNr(), decoded.SeqNr)
			assert.Equal(t, cd[:], decoded.ConfigDigest)
		})

		t.Run("with values", func(t *testing.T) {
			opts := &mockOpts{}
			cd := opts.ConfigDigest()
			orig := &datastreamsllo.LLOReportTelemetry{
				ChannelId:                       1,
				ValidAfterNanoseconds:           2,
				ObservationTimestampNanoseconds: 3,
				ReportFormat:                    4,
				Specimen:                        true,
				StreamDefinitions: []*datastreamsllo.LLOStreamDefinition{
					{
						StreamID:   5,
						Aggregator: 6,
					},
				},
				StreamValues: []*datastreamsllo.LLOStreamValue{
					{
						Type:  7,
						Value: []byte{8},
					},
				},
				ChannelOpts:  []byte{9},
				SeqNr:        opts.SeqNr(),
				ConfigDigest: cd[:],
			}
			ch <- orig

			// Wait until the telemetry is buffered.
			testutils.RequireEventually(t, func() bool {
				tm.telemetryBufferMu.Lock()
				defer tm.telemetryBufferMu.Unlock()
				return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) > 0
			})

			tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

			tLog := <-m.chTypedLogs
			assert.Equal(t, synchronization.LLOReport, tLog.telemType)
			decoded := &datastreamsllo.LLOReportTelemetry{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
			assert.Equal(t, uint32(1), decoded.ChannelId)
			assert.Equal(t, uint64(2), decoded.ValidAfterNanoseconds)
			assert.Equal(t, uint64(3), decoded.ObservationTimestampNanoseconds)
			assert.Equal(t, uint32(4), decoded.ReportFormat)
			assert.True(t, decoded.Specimen)
			assert.Len(t, decoded.StreamDefinitions, 1)
			assert.Equal(t, uint32(5), decoded.StreamDefinitions[0].StreamID)
			assert.Equal(t, uint32(6), decoded.StreamDefinitions[0].Aggregator)
			assert.Len(t, decoded.StreamValues, 1)
			assert.Equal(t, llo.LLOStreamValue_Type(7), decoded.StreamValues[0].Type)
			assert.Equal(t, []byte{8}, decoded.StreamValues[0].Value)
			assert.Equal(t, []byte{9}, decoded.ChannelOpts)
			assert.Equal(t, opts.SeqNr(), decoded.SeqNr)
			assert.Equal(t, cd[:], decoded.ConfigDigest)
		})
	})

	t.Run("appends multiple reports for same seqNr (one per channel)", func(t *testing.T) {
		// Report telemetry must append, not overwrite. Reports() generates
		// one report per reportable channel (up to ~1000+), all sharing
		// the same seqNr. Overwriting would lose all but the last channel.
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                 lggr,
			MonitoringEndpoint:     m,
			DonID:                  donID,
			CaptureReportTelemetry: true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetReportTelemetryCh()
		require.NotNil(t, ch)

		opts := &mockOpts{}
		cd := opts.ConfigDigest()

		// Send 3 reports for different channels, all with the same seqNr
		for i := uint32(1); i <= 3; i++ {
			ch <- &datastreamsllo.LLOReportTelemetry{
				ChannelId:    i,
				SeqNr:        opts.SeqNr(),
				ConfigDigest: cd[:],
			}
		}

		// Wait until all 3 are buffered
		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) == 3
		})

		// Verify buffer has all 3 entries (appended, not overwritten)
		tm.telemetryBufferMu.Lock()
		bufLen := len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()])
		tm.telemetryBufferMu.Unlock()
		assert.Equal(t, 3, bufLen, "expected 3 reports in buffer (append), got %d", bufLen)

		// Flush and verify all 3 messages are sent
		tm.TrackSeqNr(opts.ConfigDigest(), opts.SeqNr())

		receivedChannels := make([]uint32, 0, 3)
		for i := 0; i < 3; i++ {
			tLog := <-m.chTypedLogs
			assert.Equal(t, synchronization.LLOReport, tLog.telemType)
			decoded := &datastreamsllo.LLOReportTelemetry{}
			require.NoError(t, proto.Unmarshal(tLog.log, decoded))
			receivedChannels = append(receivedChannels, decoded.ChannelId)
		}
		assert.ElementsMatch(t, []uint32{1, 2, 3}, receivedChannels,
			"all 3 channel reports should be flushed")

		// Verify no additional messages were sent
		select {
		case extra := <-m.chTypedLogs:
			t.Fatalf("expected no more report messages, but got one with type %v", extra.telemType)
		default:
			// good — no extra messages
		}
	})
}

// Test_Telemeter_outcomeTelemetry_samplingAtFlushTime is a regression test for
// the bug where enabling telemetry sampling caused LLO outcome telemetry to
// drop to ~0.2/s/node instead of the expected ~1/s/node.
//
// Root cause: sampling was applied at enqueue time. For LLOOutcome, the
// fingerprint is (donId, configDigest) bucketed per second, so only the first
// Outcome of each wall-clock second was buffered. The buffer is keyed by
// seqNr and flushed only on Transmit(seqNr=X). When the OCR3 round cadence
// (DeltaRound) is much faster than the report-emission cadence
// (DefaultMinReportIntervalNanoseconds), most seqNrs do not Transmit — so the
// sampled entry was almost always evicted by a later Transmit before being
// sent.
//
// Fix: apply sampling at flush time inside sendBufferedTelemetry, so the
// buffer always sees every Outcome and the sampler decides admission on the
// (already transmit-filtered) survivor.
func Test_Telemeter_outcomeTelemetry_samplingAtFlushTime(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	donID := uint32(1)

	// Build outcomes spanning multiple wall-clock seconds. Per second, we emit
	// several Outcomes at different seqNrs but only the LAST seqNr of each
	// second transmits — mimicking a DON where DeltaRound << report interval.
	const (
		secondsCovered      = 3
		outcomesPerSecond   = 5
		baseObservationUnix = int64(1737936858)
		baseSeqNr           = uint64(1000)
	)

	cd := (&mockOpts{}).ConfigDigest()

	t.Run("with sampling enabled emits ~1 outcome per second despite seqNr/transmit mismatch", func(t *testing.T) {
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                  lggr,
			MonitoringEndpoint:      m,
			DonID:                   donID,
			CaptureOutcomeTelemetry: true,
			SampleTelemetry:         true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetOutcomeTelemetryCh()
		require.NotNil(t, ch)

		transmittingSeqNrs := make([]uint64, 0, secondsCovered)
		for s := 0; s < secondsCovered; s++ {
			secStart := time.Unix(baseObservationUnix+int64(s), 0).UnixNano()
			for i := 0; i < outcomesPerSecond; i++ {
				seqNr := baseSeqNr + uint64(s*outcomesPerSecond+i)
				// Spread observation timestamps within the same wall-clock
				// second (different nanos, same second bucket).
				obsTs := uint64(secStart + int64(i)*int64(10*time.Millisecond))
				ch <- &datastreamsllo.LLOOutcomeTelemetry{
					LifeCycleStage:                  "production",
					ObservationTimestampNanoseconds: obsTs,
					SeqNr:                           seqNr,
					ConfigDigest:                    cd[:],
					DonId:                           donID,
				}
			}
			// Only the last seqNr of each second "transmits".
			transmittingSeqNrs = append(transmittingSeqNrs, baseSeqNr+uint64(s*outcomesPerSecond+outcomesPerSecond-1))
		}

		// Wait until every outcome is buffered.
		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			return len(tm.telemetryBuffer[cd.Hex()]) == secondsCovered*outcomesPerSecond
		})

		for _, seqNr := range transmittingSeqNrs {
			tm.TrackSeqNr(cd, seqNr)
		}

		// We expect one outcome telemetry per wall-clock second (sampler
		// admits the first survivor per second bucket).
		received := make([]uint64, 0, secondsCovered)
		for i := 0; i < secondsCovered; i++ {
			select {
			case tLog := <-m.chTypedLogs:
				assert.Equal(t, synchronization.LLOOutcome, tLog.telemType)
				decoded := &datastreamsllo.LLOOutcomeTelemetry{}
				require.NoError(t, proto.Unmarshal(tLog.log, decoded))
				received = append(received, decoded.ObservationTimestampNanoseconds)
			case <-time.After(testutils.WaitTimeout(t)):
				t.Fatalf("timed out waiting for outcome telemetry #%d (got %d so far)", i+1, len(received))
			}
		}
		assert.Len(t, received, secondsCovered)

		// Verify no additional messages — sampler should have dropped the
		// remaining survivors that fell in already-seen second buckets.
		select {
		case extra := <-m.chTypedLogs:
			decoded := &datastreamsllo.LLOOutcomeTelemetry{}
			require.NoError(t, proto.Unmarshal(extra.log, decoded))
			t.Fatalf("expected no more outcome messages, got one with ts=%d", decoded.ObservationTimestampNanoseconds)
		case <-time.After(100 * time.Millisecond):
			// good
		}
	})

	t.Run("with sampling disabled emits one outcome per transmitting seqNr", func(t *testing.T) {
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                  lggr,
			MonitoringEndpoint:      m,
			DonID:                   donID,
			CaptureOutcomeTelemetry: true,
			SampleTelemetry:         false,
		})
		servicetest.Run(t, tm)
		ch := tm.GetOutcomeTelemetryCh()
		require.NotNil(t, ch)

		transmittingSeqNrs := make([]uint64, 0, secondsCovered)
		for s := 0; s < secondsCovered; s++ {
			secStart := time.Unix(baseObservationUnix+int64(s), 0).UnixNano()
			for i := 0; i < outcomesPerSecond; i++ {
				seqNr := baseSeqNr + uint64(s*outcomesPerSecond+i)
				obsTs := uint64(secStart + int64(i)*int64(10*time.Millisecond))
				ch <- &datastreamsllo.LLOOutcomeTelemetry{
					LifeCycleStage:                  "production",
					ObservationTimestampNanoseconds: obsTs,
					SeqNr:                           seqNr,
					ConfigDigest:                    cd[:],
					DonId:                           donID,
				}
			}
			transmittingSeqNrs = append(transmittingSeqNrs, baseSeqNr+uint64(s*outcomesPerSecond+outcomesPerSecond-1))
		}

		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			return len(tm.telemetryBuffer[cd.Hex()]) == secondsCovered*outcomesPerSecond
		})

		for _, seqNr := range transmittingSeqNrs {
			tm.TrackSeqNr(cd, seqNr)
		}

		// Without sampling, every transmit produces telemetry — one per second
		// (since only one seqNr per second transmits).
		for i := 0; i < secondsCovered; i++ {
			select {
			case tLog := <-m.chTypedLogs:
				assert.Equal(t, synchronization.LLOOutcome, tLog.telemType)
			case <-time.After(testutils.WaitTimeout(t)):
				t.Fatalf("timed out waiting for outcome telemetry #%d", i+1)
			}
		}

		select {
		case <-m.chTypedLogs:
			t.Fatal("expected no more outcome messages")
		case <-time.After(100 * time.Millisecond):
		}
	})
}

// Test_Telemeter_reportTelemetry_samplingAtFlushTime is the analogous
// regression test to Test_Telemeter_outcomeTelemetry_samplingAtFlushTime, but
// for LLOReport. The same root cause applies: sampling was applied at enqueue
// time, but the buffer is seqNr-keyed and only the transmitting seqNr's entry
// survives flush. When most seqNrs do not transmit, sampled report entries
// were evicted before they could be sent.
//
// LLOReport's sampler fingerprint additionally keys on channelId, so each
// channel has its own per-second bucket.
func Test_Telemeter_reportTelemetry_samplingAtFlushTime(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	donID := uint32(1)

	const (
		secondsCovered      = 3
		seqNrsPerSecond     = 5
		baseObservationUnix = int64(1737936858)
		baseSeqNr           = uint64(2000)
	)

	channels := []uint32{1, 2}
	cd := (&mockOpts{}).ConfigDigest()

	t.Run("with sampling enabled emits ~1 report per channel per second despite seqNr/transmit mismatch", func(t *testing.T) {
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                 lggr,
			MonitoringEndpoint:     m,
			DonID:                  donID,
			CaptureReportTelemetry: true,
			SampleTelemetry:        true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetReportTelemetryCh()
		require.NotNil(t, ch)

		transmittingSeqNrs := make([]uint64, 0, secondsCovered)
		expectedBufferEntries := 0
		for s := 0; s < secondsCovered; s++ {
			secStart := time.Unix(baseObservationUnix+int64(s), 0).UnixNano()
			for i := 0; i < seqNrsPerSecond; i++ {
				seqNr := baseSeqNr + uint64(s*seqNrsPerSecond+i)
				obsTs := uint64(secStart + int64(i)*int64(10*time.Millisecond))
				// Each seqNr emits a report per channel — mimics the
				// Reports() call shape in LLO (one report per channel).
				for _, channelID := range channels {
					ch <- &datastreamsllo.LLOReportTelemetry{
						ChannelId:                       channelID,
						ObservationTimestampNanoseconds: obsTs,
						SeqNr:                           seqNr,
						ConfigDigest:                    cd[:],
					}
					expectedBufferEntries++
				}
			}
			// Only the last seqNr of each second transmits.
			transmittingSeqNrs = append(transmittingSeqNrs, baseSeqNr+uint64(s*seqNrsPerSecond+seqNrsPerSecond-1))
		}

		// Wait until every report is buffered (no enqueue-time sampling).
		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			total := 0
			for _, m := range tm.telemetryBuffer[cd.Hex()] {
				total += len(m)
			}
			return total == expectedBufferEntries
		})

		for _, seqNr := range transmittingSeqNrs {
			tm.TrackSeqNr(cd, seqNr)
		}

		// Expect one report per channel per second bucket (sampler admits the
		// first survivor per (channelId, second) fingerprint).
		expected := secondsCovered * len(channels)
		seen := make(map[uint32]int)
		for i := 0; i < expected; i++ {
			select {
			case tLog := <-m.chTypedLogs:
				assert.Equal(t, synchronization.LLOReport, tLog.telemType)
				decoded := &datastreamsllo.LLOReportTelemetry{}
				require.NoError(t, proto.Unmarshal(tLog.log, decoded))
				seen[decoded.ChannelId]++
			case <-time.After(testutils.WaitTimeout(t)):
				t.Fatalf("timed out waiting for report telemetry #%d (got %d so far)", i+1, len(seen))
			}
		}
		for _, channelID := range channels {
			assert.Equal(t, secondsCovered, seen[channelID],
				"expected one report per second for channel %d", channelID)
		}

		select {
		case extra := <-m.chTypedLogs:
			decoded := &datastreamsllo.LLOReportTelemetry{}
			require.NoError(t, proto.Unmarshal(extra.log, decoded))
			t.Fatalf("expected no more report messages, got one with channel=%d ts=%d",
				decoded.ChannelId, decoded.ObservationTimestampNanoseconds)
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("with sampling disabled emits all reports for transmitting seqNrs", func(t *testing.T) {
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                 lggr,
			MonitoringEndpoint:     m,
			DonID:                  donID,
			CaptureReportTelemetry: true,
			SampleTelemetry:        false,
		})
		servicetest.Run(t, tm)
		ch := tm.GetReportTelemetryCh()
		require.NotNil(t, ch)

		transmittingSeqNrs := make([]uint64, 0, secondsCovered)
		expectedBufferEntries := 0
		for s := 0; s < secondsCovered; s++ {
			secStart := time.Unix(baseObservationUnix+int64(s), 0).UnixNano()
			for i := 0; i < seqNrsPerSecond; i++ {
				seqNr := baseSeqNr + uint64(s*seqNrsPerSecond+i)
				obsTs := uint64(secStart + int64(i)*int64(10*time.Millisecond))
				for _, channelID := range channels {
					ch <- &datastreamsllo.LLOReportTelemetry{
						ChannelId:                       channelID,
						ObservationTimestampNanoseconds: obsTs,
						SeqNr:                           seqNr,
						ConfigDigest:                    cd[:],
					}
					expectedBufferEntries++
				}
			}
			transmittingSeqNrs = append(transmittingSeqNrs, baseSeqNr+uint64(s*seqNrsPerSecond+seqNrsPerSecond-1))
		}

		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			total := 0
			for _, m := range tm.telemetryBuffer[cd.Hex()] {
				total += len(m)
			}
			return total == expectedBufferEntries
		})

		for _, seqNr := range transmittingSeqNrs {
			tm.TrackSeqNr(cd, seqNr)
		}

		// Without sampling, every report at a transmitting seqNr flushes.
		// One transmitting seqNr per second × len(channels) reports per seqNr.
		expected := secondsCovered * len(channels)
		for i := 0; i < expected; i++ {
			select {
			case tLog := <-m.chTypedLogs:
				assert.Equal(t, synchronization.LLOReport, tLog.telemType)
			case <-time.After(testutils.WaitTimeout(t)):
				t.Fatalf("timed out waiting for report telemetry #%d", i+1)
			}
		}

		select {
		case <-m.chTypedLogs:
			t.Fatal("expected no more report messages")
		case <-time.After(100 * time.Millisecond):
		}
	})

	t.Run("with sampling enabled multiple reports per channel at same transmitting seqNr keep one per second", func(t *testing.T) {
		// Reports() emits one report per channel for a given seqNr. When
		// that seqNr does transmit, the per-channel reports should be
		// admitted (different fingerprints). This test verifies the
		// append-semantics interaction with flush-time sampling.
		m := &mockMonitoringEndpoint{chTypedLogs: make(chan typedLog, 100)}
		tm := newTelemeter(TelemeterParams{
			Logger:                 lggr,
			MonitoringEndpoint:     m,
			DonID:                  donID,
			CaptureReportTelemetry: true,
			SampleTelemetry:        true,
		})
		servicetest.Run(t, tm)
		ch := tm.GetReportTelemetryCh()
		require.NotNil(t, ch)

		opts := &mockOpts{}
		cd := opts.ConfigDigest()
		obsTs := uint64(time.Unix(baseObservationUnix, 0).UnixNano())

		// Append 3 reports at the same seqNr for distinct channels.
		for _, channelID := range []uint32{10, 20, 30} {
			ch <- &datastreamsllo.LLOReportTelemetry{
				ChannelId:                       channelID,
				ObservationTimestampNanoseconds: obsTs,
				SeqNr:                           opts.SeqNr(),
				ConfigDigest:                    cd[:],
			}
		}

		testutils.RequireEventually(t, func() bool {
			tm.telemetryBufferMu.Lock()
			defer tm.telemetryBufferMu.Unlock()
			return len(tm.telemetryBuffer[cd.Hex()][opts.SeqNr()]) == 3
		})

		tm.TrackSeqNr(cd, opts.SeqNr())

		received := make(map[uint32]struct{})
		for i := 0; i < 3; i++ {
			select {
			case tLog := <-m.chTypedLogs:
				decoded := &datastreamsllo.LLOReportTelemetry{}
				require.NoError(t, proto.Unmarshal(tLog.log, decoded))
				received[decoded.ChannelId] = struct{}{}
			case <-time.After(testutils.WaitTimeout(t)):
				t.Fatalf("timed out waiting for report #%d", i+1)
			}
		}
		assert.Equal(t, map[uint32]struct{}{10: {}, 20: {}, 30: {}}, received,
			"each per-channel report should be admitted (distinct sampler fingerprints)")
	})
}

func ptr[T any](t T) *T { return &t }
