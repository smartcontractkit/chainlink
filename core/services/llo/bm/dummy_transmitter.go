package bm

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-data-streams/llo/reportcodecs/evm"

	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

// A dummy transmitter useful for benchmarking and testing

var (
	promTransmitSuccessCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "dummytransmitter",
		Name:      "transmit_success_count",
		Help:      "Running count of successful transmits",
	})
)

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

type transmitter struct {
	lggr        logger.Logger
	fromAccount string
}

func NewTransmitter(lggr logger.Logger, fromAccount string) Transmitter {
	return &transmitter{
		logger.Named(lggr, "DummyTransmitter"),
		fromAccount,
	}
}

func (t *transmitter) Start(context.Context) error {
	return nil
}

func (t *transmitter) Close() error {
	return nil
}

func (t *transmitter) Transmit(
	ctx context.Context,
	digest ocr2types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []ocr2types.AttributedOnchainSignature,
) error {
	lggr := t.lggr
	{
		switch report.Info.ReportFormat {
		case llotypes.ReportFormatJSON:
			r, err := (llo.JSONReportCodec{}).Decode(report.Report)
			if err != nil {
				lggr.Debugw(fmt.Sprintf("Failed to decode report with type %s", report.Info.ReportFormat), "err", err)
			} else if r.SeqNr > 0 {
				lggr = logger.With(lggr,
					"report.Report.ConfigDigest", r.ConfigDigest,
					"report.Report.SeqNr", r.SeqNr,
					"report.Report.ChannelID", r.ChannelID,
					"report.Report.ValidAfterNanoseconds", r.ValidAfterNanoseconds,
					"report.Report.ObservationTimestampNanoseconds", r.ObservationTimestampNanoseconds,
					"report.Report.Values", r.Values,
					"report.Report.Specimen", r.Specimen,
				)
			}
		case llotypes.ReportFormatEVMPremiumLegacy:
			r, err := (evm.ReportCodecPremiumLegacy{}).Decode(report.Report)
			if err != nil {
				lggr.Debugw(fmt.Sprintf("Failed to decode report with type %s", report.Info.ReportFormat), "err", err)
			} else if r.ObservationsTimestamp > 0 {
				lggr = logger.With(lggr,
					"report.Report.FeedId", r.FeedId,
					"report.Report.ObservationsTimestamp", r.ObservationsTimestamp,
					"report.Report.BenchmarkPrice", r.BenchmarkPrice,
					"report.Report.Bid", r.Bid,
					"report.Report.Ask", r.Ask,
					"report.Report.ValidFromTimestamp", r.ValidFromTimestamp,
					"report.Report.ExpiresAt", r.ExpiresAt,
					"report.Report.LinkFee", r.LinkFee,
					"report.Report.NativeFee", r.NativeFee,
				)
			}
		default:
			err := fmt.Errorf("unhandled report format: %s", report.Info.ReportFormat)
			lggr.Debugw(fmt.Sprintf("Failed to decode report with type %s", report.Info.ReportFormat), "err", err)
		}
	}
	promTransmitSuccessCount.Inc()
	lggr.Infow("Transmit (dummy)", "digest", digest, "seqNr", seqNr, "report.Report", report.Report, "report.Info", report.Info, "sigs", sigs)
	return nil
}

// FromAccount returns the stringified (hex) CSA public key
func (t *transmitter) FromAccount(context.Context) (ocr2types.Account, error) {
	return ocr2types.Account(t.fromAccount), nil
}

func (t *transmitter) Ready() error { return nil }

func (t *transmitter) HealthReport() map[string]error {
	report := map[string]error{t.Name(): nil}
	return report
}

func (t *transmitter) Name() string { return t.lggr.Name() }
