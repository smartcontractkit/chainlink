package bm

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

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

func NewTransmitter(lggr logger.Logger, fromAccount ed25519.PublicKey) Transmitter {
	return &transmitter{
		lggr.Named("DummyTransmitter"),
		fmt.Sprintf("%x", fromAccount),
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
	transmitSuccessCount.Inc()
	t.lggr.Debugw("Transmit", "digest", digest, "seqNr", seqNr, "report.Report", report.Report, "report.Info", report.Info, "sigs", sigs)
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
