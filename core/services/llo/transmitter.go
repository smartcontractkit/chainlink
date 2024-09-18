package llo

import (
	"context"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
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

// TODO: prom metrics (common with mercury/transmitter.go?)
// https://smartcontract-it.atlassian.net/browse/MERC-3659

const (
	// Mercury server error codes
	DuplicateReport = 2
	// TODO: revisit these values in light of parallel composition
	// https://smartcontract-it.atlassian.net/browse/MERC-3659
	// maxTransmitQueueSize = 10_000
	// maxDeleteQueueSize   = 10_000
	// transmitTimeout      = 5 * time.Second
)

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

type transmitter struct {
	services.StateMachine
	lggr        logger.Logger
	fromAccount string

	subTransmitters []Transmitter
}

type TransmitterOpts struct {
	Lggr                   logger.Logger
	FromAccount            string
	MercuryTransmitterOpts mercurytransmitter.Opts
}

// The transmitter will handle starting and stopping the subtransmitters
func NewTransmitter(opts TransmitterOpts) Transmitter {
	subTransmitters := []Transmitter{
		mercurytransmitter.New(opts.MercuryTransmitterOpts),
	}
	return &transmitter{
		services.StateMachine{},
		opts.Lggr,
		opts.FromAccount,
		subTransmitters,
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
func (t *transmitter) FromAccount() (ocr2types.Account, error) {
	return ocr2types.Account(t.fromAccount), nil
}
