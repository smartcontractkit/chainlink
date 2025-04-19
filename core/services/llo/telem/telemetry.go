package telem

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink-data-streams/llo/reportcodecs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/eautils"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	legacytelem "github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

const adapterLWBAErrorName = "AdapterLWBAError"

type Telemeter interface {
	EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error)
	MakeObservationScopedTelemetryCh(opts llo.DSOpts, size int) (ch chan<- interface{})
	GetOutcomeTelemetryCh() chan<- *datastreamsllo.LLOOutcomeTelemetry
	GetReportTelemetryCh() chan<- *datastreamsllo.LLOReportTelemetry
	CaptureEATelemetry() bool
	CaptureObservationTelemetry() bool
}

type TelemeterService interface {
	Telemeter
	services.Service
}

type TelemeterParams struct {
	Logger                      logger.Logger
	MonitoringEndpoint          telemetry.MultitypeMonitoringEndpoint
	DonID                       uint32
	CaptureEATelemetry          bool
	CaptureObservationTelemetry bool
	CaptureOutcomeTelemetry     bool
	CaptureReportTelemetry      bool
}

func NewTelemeterService(params TelemeterParams) TelemeterService {
	if params.MonitoringEndpoint == nil {
		return NullTelemeter
	}
	return newTelemeter(params)
}

func newTelemeter(params TelemeterParams) *telemeter {
	// NOTE: This channel must take multiple telemetry packets per round (1 per
	// feed) so we need to make sure the buffer is large enough.
	//
	// 2000 feeds * 5s/250ms = 40_000 should hold ~5s of buffer in the worst case.
	chTelemetryPipeline := make(chan *TelemetryPipeline, 40_000)
	t := &telemeter{
		chTelemetryPipeline:         chTelemetryPipeline,
		monitoringEndpoint:          params.MonitoringEndpoint,
		donID:                       params.DonID,
		chch:                        make(chan telemetryCollectionContext, 1000), // chch should be consumed from very quickly so we don't need a large buffer, but it also won't hurt
		captureEATelemetry:          params.CaptureEATelemetry,
		captureObservationTelemetry: params.CaptureObservationTelemetry,
	}
	if params.CaptureOutcomeTelemetry {
		t.chOutcomeTelemetry = make(chan *datastreamsllo.LLOOutcomeTelemetry, 100) // only one per round so 100 buffer should be more than enough even for very fast rounds
	}
	if params.CaptureReportTelemetry {
		t.chReportTelemetry = make(chan *datastreamsllo.LLOReportTelemetry, (2+2)*datastreamsllo.MaxReportCount) // 2 instances+2x size safety buffer
	}
	t.Service, t.eng = services.Config{
		Name:  "LLOTelemeterService",
		Start: t.start,
		Close: func() error {
			// Enforce that nothing is transmitted after the service is closed
			// Doing so is a programming error
			close(t.chch)
			if t.chTelemetryPipeline != nil {
				close(t.chTelemetryPipeline)
			}
			if t.chOutcomeTelemetry != nil {
				close(t.chOutcomeTelemetry)
			}
			if t.chReportTelemetry != nil {
				close(t.chReportTelemetry)
			}
			return nil
		},
	}.NewServiceEngine(params.Logger)

	return t
}

type telemeter struct {
	services.Service
	eng *services.Engine

	monitoringEndpoint  telemetry.MultitypeMonitoringEndpoint
	chTelemetryPipeline chan *TelemetryPipeline

	donID                       uint32
	captureEATelemetry          bool
	captureObservationTelemetry bool
	chch                        chan telemetryCollectionContext
	chOutcomeTelemetry          chan *datastreamsllo.LLOOutcomeTelemetry
	chReportTelemetry           chan *datastreamsllo.LLOReportTelemetry
}

func (t *telemeter) EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error) {
	if t.Service.Ready() != nil {
		// This should never happen, telemeter should always be started BEFORE
		// the oracle and closed AFTER it
		t.eng.Errorw("Telemeter not ready, dropping observation", "run", run, "streamID", streamID, "opts", opts, "val", val, "err", err)
		return
	}
	var adapterError *eautils.AdapterError
	var dpInvariantViolationDetected bool
	if errors.As(err, &adapterError) && adapterError.Name == adapterLWBAErrorName {
		dpInvariantViolationDetected = true
	} else if err != nil {
		// ignore errors
		return
	}
	tp := &TelemetryPipeline{run, trrs, streamID, opts, val, dpInvariantViolationDetected}
	select {
	case t.chTelemetryPipeline <- tp:
	default:
	}
}

type telemetryCollectionContext struct {
	in   <-chan interface{}
	opts llo.DSOpts
}

// MakeObservationScopedTelemetryCh reads telem packets from the returned channel and sends them
// to the monitoring endpoint. Stops reading when channel closed or when
// telemeter is stopped
//
// CALLER IS RESPONSIBLE FOR CLOSING THE RETURNED CHANNEL TO AVOID MEMORY
// LEAKS.
//
// It is necessary to make a new channel for every Observation call because it
// closes over DSOpts which is scoped to that call only.
func (t *telemeter) MakeObservationScopedTelemetryCh(opts llo.DSOpts, size int) chan<- interface{} {
	if !(t.captureObservationTelemetry || t.captureEATelemetry) {
		return nil
	}
	ch := make(chan interface{}, size)
	tcc := telemetryCollectionContext{
		in:   ch,
		opts: opts,
	}

	select {
	case t.chch <- tcc:
	default:
		// This should be performant enough with buffer of t.chch large enough
		// that we never hit this case, however, we should NEVER block
		// observations on telemetry even if something pathological happens.
		t.eng.Errorw("Telemeter chch full, will not record telemetry", "seqNr", opts.SeqNr())
		return nil
	}
	return ch
}

func (t *telemeter) GetOutcomeTelemetryCh() chan<- *datastreamsllo.LLOOutcomeTelemetry {
	return t.chOutcomeTelemetry
}

func (t *telemeter) GetReportTelemetryCh() chan<- *datastreamsllo.LLOReportTelemetry {
	return t.chReportTelemetry
}

func (t *telemeter) CaptureEATelemetry() bool {
	return t.captureEATelemetry
}

func (t *telemeter) CaptureObservationTelemetry() bool {
	return t.captureObservationTelemetry
}

func (t *telemeter) start(_ context.Context) error {
	t.eng.Go(func(ctx context.Context) {
		wg := sync.WaitGroup{}
		for {
			select {
			case tcc := <-t.chch:
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						select {
						case p, ok := <-tcc.in:
							if !ok {
								// channel closed by producer
								return
							}
							t.collectObservationTelemetry(p, tcc.opts)
						case <-ctx.Done():
							return
						}
					}
				}()

			case rt := <-t.chOutcomeTelemetry:
				t.collectOutcomeTelemetry(rt)
			case rt := <-t.chReportTelemetry:
				t.collectReportTelemetry(rt)
			case p := <-t.chTelemetryPipeline:
				t.collectV3PremiumLegacyTelemetry(p)
			case <-ctx.Done():
				wg.Wait()
				return
			}
		}
	})
	return nil
}

func (t *telemeter) collectObservationTelemetry(p interface{}, opts llo.DSOpts) {
	var telemType synchronization.TelemetryType
	var msg proto.Message
	switch v := p.(type) {
	case *pipeline.BridgeTelemetry:
		telemType = synchronization.PipelineBridge
		cd := opts.ConfigDigest()
		msg = &LLOBridgeTelemetry{
			BridgeAdapterName:        v.Name,
			BridgeRequestData:        v.RequestData,
			BridgeResponseData:       v.ResponseData,
			BridgeResponseError:      v.ResponseError,
			BridgeResponseStatusCode: int32(v.ResponseStatusCode), //nolint:gosec // G115 // even if overflow does happen, its harmless
			RequestStartTimestamp:    v.RequestStartTimestamp.UnixNano(),
			RequestFinishTimestamp:   v.RequestFinishTimestamp.UnixNano(),
			LocalCacheHit:            v.LocalCacheHit,
			SpecId:                   v.SpecID,
			StreamId:                 v.StreamID,
			DotId:                    v.DotID,
			DonId:                    t.donID,
			SeqNr:                    opts.SeqNr(),
			ConfigDigest:             cd[:],
			ObservationTimestamp:     opts.ObservationTimestamp().UnixNano(),
		}
		if opts.VerboseLogging() {
			t.eng.Infow("Sending LLOBridgeTelemetry telemetry", "StreamId", v.StreamID, "BridgeAdapterName", v.Name, "BridgeResponseError", v.ResponseError)
		}
	case *LLOObservationTelemetry:
		telemType = synchronization.LLOObservation
		v.DonId = t.donID
		msg = v
		if opts.VerboseLogging() {
			t.eng.Infow("Sending LLOObservationTelemetry telemetry", "StreamId", v.StreamId, "ObservationTimestamp", v.ObservationTimestamp)
		}
	default:
		t.eng.Warnw("Unknown telemetry type", "type", fmt.Sprintf("%T", p))
		return
	}
	bytes, err := proto.Marshal(msg)
	if err != nil {
		t.eng.Warnf("protobuf marshal failed %v", err.Error())
		return
	}

	t.monitoringEndpoint.SendTypedLog(telemType, bytes)
}

func (t *telemeter) collectOutcomeTelemetry(rt *datastreamsllo.LLOOutcomeTelemetry) {
	bytes, err := proto.Marshal(rt)
	if err != nil {
		t.eng.Warnf("protobuf marshal failed %v", err.Error())
		return
	}
	t.monitoringEndpoint.SendTypedLog(synchronization.LLOOutcome, bytes)
}

func (t *telemeter) collectReportTelemetry(rt *datastreamsllo.LLOReportTelemetry) {
	bytes, err := proto.Marshal(rt)
	if err != nil {
		t.eng.Warnf("protobuf marshal failed %v", err.Error())
		return
	}
	t.monitoringEndpoint.SendTypedLog(synchronization.LLOReport, bytes)
}

func (t *telemeter) collectV3PremiumLegacyTelemetry(d *TelemetryPipeline) {
	eaTelemetryValues := ocrcommon.ParseMercuryEATelemetry(t.eng.SugaredLogger, d.trrs, mercuryutils.REPORT_V3)
	for _, eaTelem := range eaTelemetryValues {
		var benchmarkPrice, bidPrice, askPrice int64
		var bp, bid, ask string
		switch v := d.val.(type) {
		case *llo.Decimal:
			benchmarkPrice = v.Decimal().IntPart()
			bp = v.Decimal().String()
		case *llo.Quote:
			benchmarkPrice = v.Benchmark.IntPart()
			bp = v.Benchmark.String()
			bidPrice = v.Bid.IntPart()
			bid = v.Bid.String()
			askPrice = v.Ask.IntPart()
			ask = v.Ask.String()
		}
		tea := &legacytelem.EnhancedEAMercury{
			DataSource:                      eaTelem.DataSource,
			DpBenchmarkPrice:                eaTelem.DpBenchmarkPrice,
			DpBid:                           eaTelem.DpBid,
			DpAsk:                           eaTelem.DpAsk,
			DpInvariantViolationDetected:    d.dpInvariantViolationDetected,
			BridgeTaskRunStartedTimestamp:   eaTelem.BridgeTaskRunStartedTimestamp,
			BridgeTaskRunEndedTimestamp:     eaTelem.BridgeTaskRunEndedTimestamp,
			ProviderRequestedTimestamp:      eaTelem.ProviderRequestedTimestamp,
			ProviderReceivedTimestamp:       eaTelem.ProviderReceivedTimestamp,
			ProviderDataStreamEstablished:   eaTelem.ProviderDataStreamEstablished,
			ProviderIndicatedTime:           eaTelem.ProviderIndicatedTime,
			Feed:                            fmt.Sprintf("streamID:%d", d.streamID),
			ObservationBenchmarkPrice:       benchmarkPrice,
			ObservationBid:                  bidPrice,
			ObservationAsk:                  askPrice,
			ObservationBenchmarkPriceString: bp,
			ObservationBidString:            bid,
			ObservationAskString:            ask,
			IsLinkFeed:                      false,
			IsNativeFeed:                    false,
			ConfigDigest:                    d.opts.ConfigDigest().Hex(),
			AssetSymbol:                     eaTelem.AssetSymbol,
			Version:                         uint32(1000 + mercuryutils.REPORT_V3), // add 1000 to distinguish between legacy feeds, this can be changed if necessary
			DonId:                           t.donID,
		}
		epoch, round, err := evm.SeqNrToEpochAndRound(d.opts.OutCtx().SeqNr)
		if err != nil {
			t.eng.Warnw("Failed to convert sequence number to epoch and round", "err", err)
		} else {
			tea.Round = int64(round)
			tea.Epoch = int64(epoch)
		}

		bytes, err := proto.Marshal(tea)
		if err != nil {
			t.eng.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}

		t.monitoringEndpoint.SendTypedLog(synchronization.EnhancedEAMercury, bytes)
	}
}

type TelemetryObserve struct {
	Opts      llo.DSOpts
	Telemetry interface{}
}

type TelemetryPipeline struct {
	run                          *pipeline.Run
	trrs                         pipeline.TaskRunResults
	streamID                     uint32
	opts                         llo.DSOpts
	val                          llo.StreamValue
	dpInvariantViolationDetected bool
}

var NullTelemeter TelemeterService = &nullTelemeter{}

type nullTelemeter struct{}

func (t *nullTelemeter) EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error) {
}
func (t *nullTelemeter) MakeObservationScopedTelemetryCh(opts llo.DSOpts, size int) (ch chan<- interface{}) {
	return nil
}
func (t *nullTelemeter) GetOutcomeTelemetryCh() chan<- *datastreamsllo.LLOOutcomeTelemetry {
	return nil
}
func (t *nullTelemeter) GetReportTelemetryCh() chan<- *datastreamsllo.LLOReportTelemetry {
	return nil
}
func (t *nullTelemeter) CaptureEATelemetry() bool {
	return false
}
func (t *nullTelemeter) CaptureObservationTelemetry() bool {
	return false
}
func (t *nullTelemeter) Start(context.Context) error {
	return nil
}
func (t *nullTelemeter) Close() error {
	return nil
}
func (t *nullTelemeter) Healthy() error {
	return nil
}
func (t *nullTelemeter) Unhealthy() error {
	return nil
}
func (t *nullTelemeter) HealthReport() map[string]error {
	return nil
}
func (t *nullTelemeter) Name() string {
	return "NullTelemeter"
}
func (t *nullTelemeter) Ready() error {
	return nil
}
