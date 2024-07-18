package bm

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// A dummy transmitter useful for benchmarking and testing

var (
	transmitSuccessCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "llo_transmit_success_count",
		Help: "Running count of successful transmits",
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
		lggr.Named("DummyTransmitter"),
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
	digest types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) error {
	lggr := t.lggr
	switch report.Info.ReportFormat {
	case llotypes.ReportFormatJSON:
		r, err := (llo.JSONReportCodec{}).Decode(report.Report)
		if err != nil {
			lggr.Debugw("Failed to decode JSON report", "err", err)
		}
		lggr = lggr.With(
			"report.Report.ConfigDigest", r.ConfigDigest,
			"report.Report.SeqNr", r.SeqNr,
			"report.Report.ChannelID", r.ChannelID,
			"report.Report.ValidAfterSeconds", r.ValidAfterSeconds,
			"report.Report.Values", r.Values,
			"report.Report.Specimen", r.Specimen,
		)
	default:
	}
	transmitSuccessCount.Inc()
	lggr.Infow("Transmit (dummy)", "digest", digest, "seqNr", seqNr, "report.Report", report.Report, "report.Info", report.Info, "sigs", sigs)
	return nil
}

// FromAccount returns the stringified (hex) CSA public key
func (t *transmitter) FromAccount() (ocr2types.Account, error) {
	return ocr2types.Account(t.fromAccount), nil
}

func (t *transmitter) Ready() error { return nil }

func (t *transmitter) HealthReport() map[string]error {
	report := map[string]error{t.Name(): nil}
	return report
}

func (t *transmitter) Name() string { return t.lggr.Name() }
