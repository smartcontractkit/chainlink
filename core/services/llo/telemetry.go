package llo

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/eautils"
	mercuryutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization/telem"
)

const adapterLWBAErrorName = "AdapterLWBAError"

type Telemeter interface {
	EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error)
}

type TelemeterService interface {
	Telemeter
	services.Service
}

func NewTelemeterService(lggr logger.Logger, monitoringEndpoint commontypes.MonitoringEndpoint, donID uint32) TelemeterService {
	if monitoringEndpoint == nil {
		return NullTelemeter
	}
	return newTelemeter(lggr, monitoringEndpoint, donID)
}

func newTelemeter(lggr logger.Logger, monitoringEndpoint commontypes.MonitoringEndpoint, donID uint32) *telemeter {
	// NOTE: This channel must take multiple telemetry packets per round (1 per
	// feed) so we need to make sure the buffer is large enough.
	//
	// 2000 feeds * 5s/250ms = 40_000 should hold ~5s of buffer in the worst case.
	chTelemetryObservation := make(chan TelemetryObservation, 40_000)
	t := &telemeter{
		chTelemetryObservation: chTelemetryObservation,
		monitoringEndpoint:     monitoringEndpoint,
		donID:                  donID,
	}
	t.Service, t.eng = services.Config{
		Name:  "LLOTelemeterService",
		Start: t.start,
	}.NewServiceEngine(lggr)

	return t
}

type telemeter struct {
	services.Service
	eng *services.Engine

	monitoringEndpoint     commontypes.MonitoringEndpoint
	chTelemetryObservation chan TelemetryObservation
	donID                  uint32
}

func (t *telemeter) EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error) {
	if t.Service.Ready() != nil {
		// This should never happen, telemeter should always be started BEFORE
		// the oracle and closed AFTER it
		t.eng.SugaredLogger.Errorw("Telemeter not ready, dropping observation", "run", run, "streamID", streamID, "opts", opts, "val", val, "err", err)
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
	tObs := TelemetryObservation{run, trrs, streamID, opts, val, dpInvariantViolationDetected}
	select {
	case t.chTelemetryObservation <- tObs:
	default:
	}
}

func (t *telemeter) start(_ context.Context) error {
	t.eng.Go(func(ctx context.Context) {
		for {
			select {
			case tObs := <-t.chTelemetryObservation:
				t.collectV3PremiumLegacyTelemetry(tObs)
			case <-ctx.Done():
				return
			}
		}
	})
	return nil
}

func (t *telemeter) collectV3PremiumLegacyTelemetry(d TelemetryObservation) {
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
		tea := &telem.EnhancedEAMercury{
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
			t.eng.SugaredLogger.Warnw("Failed to convert sequence number to epoch and round", "err", err)
		} else {
			tea.Round = int64(round)
			tea.Epoch = int64(epoch)
		}

		bytes, err := proto.Marshal(tea)
		if err != nil {
			t.eng.SugaredLogger.Warnf("protobuf marshal failed %v", err.Error())
			continue
		}

		t.monitoringEndpoint.SendLog(bytes)
	}
}

type TelemetryObservation struct {
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
