package llo

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/mercurytransmitter"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

// LLO Transmitter implementation, based on
// core/services/relay/evm/mercury/transmitter.go
//
// If you need to "fan-out" transmits and send reports to a new destination,
// add a new subTransmitter

const (
	// Mercury server error codes
	DuplicateReport = 2
)

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

type TransmitterRetirementReportCacheWriter interface {
	StoreAttestedRetirementReport(ctx context.Context, cd ocrtypes.ConfigDigest, seqNr uint64, retirementReport []byte, sigs []types.AttributedOnchainSignature) error
}

type transmitter struct {
	services.StateMachine
	lggr           logger.Logger
	verboseLogging bool
	fromAccount    string

	subTransmitters       []Transmitter
	retirementReportCache TransmitterRetirementReportCacheWriter
}

type TransmitterOpts struct {
	Lggr                   logger.Logger
	VerboseLogging         bool
	FromAccount            string
	MercuryTransmitterOpts mercurytransmitter.Opts
	RetirementReportCache  TransmitterRetirementReportCacheWriter
}

// The transmitter will handle starting and stopping the subtransmitters
func NewTransmitter(opts TransmitterOpts) Transmitter {
	subTransmitters := []Transmitter{
		mercurytransmitter.New(opts.MercuryTransmitterOpts),
	}
	return &transmitter{
		services.StateMachine{},
		opts.Lggr,
		opts.VerboseLogging,
		opts.FromAccount,
		subTransmitters,
		opts.RetirementReportCache,
	}
}

func (t *transmitter) Start(ctx context.Context) error {
	return t.StartOnce("llo.Transmitter", func() error {
		for _, st := range t.subTransmitters {
			if err := st.Start(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *transmitter) Close() error {
	return t.StopOnce("llo.Transmitter", func() error {
		for _, st := range t.subTransmitters {
			if err := st.Close(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (t *transmitter) HealthReport() map[string]error {
	report := map[string]error{t.Name(): t.Healthy()}
	for _, st := range t.subTransmitters {
		services.CopyHealth(report, st.HealthReport())
	}
	return report
}

func (t *transmitter) Name() string { return t.lggr.Name() }

func (t *transmitter) Transmit(
	ctx context.Context,
	digest types.ConfigDigest,
	seqNr uint64,
	report ocr3types.ReportWithInfo[llotypes.ReportInfo],
	sigs []types.AttributedOnchainSignature,
) (err error) {
	if t.verboseLogging {
		t.lggr.Debugw("Transmit report", "digest", digest, "seqNr", seqNr, "report", report, "sigs", sigs)
	}

	if report.Info.ReportFormat == llotypes.ReportFormatRetirement {
		// Retirement reports don't get transmitted; rather, they are stored in
		// the RetirementReportCache
		t.lggr.Debugw("Storing retirement report", "digest", digest, "seqNr", seqNr)
		if err := t.retirementReportCache.StoreAttestedRetirementReport(ctx, digest, seqNr, report.Report, sigs); err != nil {
			return fmt.Errorf("failed to write retirement report to cache: %w", err)
		}
		return nil
	}
	g := new(errgroup.Group)
	for _, st := range t.subTransmitters {
		st := st
		g.Go(func() error {
			return st.Transmit(ctx, digest, seqNr, report, sigs)
		})
	}
	return g.Wait()
}

// FromAccount returns the stringified (hex) CSA public key
func (t *transmitter) FromAccount(ctx context.Context) (ocr2types.Account, error) {
	return ocr2types.Account(t.fromAccount), nil
}
