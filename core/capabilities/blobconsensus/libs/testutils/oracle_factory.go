package testutils

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ core.OracleFactory = (*oracleFactory)(nil)

type oracleFactory struct {
	t    *testing.T
	lggr logger.Logger
}

func NewOracleFactory(t *testing.T, lggr logger.Logger) *oracleFactory {
	return &oracleFactory{
		t:    t,
		lggr: lggr,
	}
}

// NewOracle(ctx context.Context, args OracleArgs) (Oracle, error)
func (of *oracleFactory) NewOracle(ctx context.Context, args core.OracleArgs) (core.Oracle, error) {
	return &oracle{
		config: args,
		wg:     &sync.WaitGroup{},
		lggr:   of.lggr,
	}, nil
}

type oracle struct {
	config core.OracleArgs
	wg     *sync.WaitGroup
	cancel context.CancelFunc
	lggr   logger.Logger
}

func (o *oracle) Start(ctx context.Context) error {
	ctx, o.cancel = context.WithCancel(ctx)

	config := ocr3types.ReportingPluginConfig{
		F: 1,
		N: 4,
	}

	reportingPlugin, _, err := o.config.ReportingPluginFactoryService.NewReportingPlugin(
		ctx,
		config,
	)
	if err != nil {
		return err
	}

	outcomeCtx := ocr3types.OutcomeContext{
		SeqNr: 1,
	}
	o.wg.Add(1)
	go func() {
		defer o.wg.Done()
		for {
			select {
			case <-ctx.Done():
				o.lggr.Info("context canceled, stopping the routine")
				return
			case <-time.After(250 * time.Millisecond):
				// Perform the logic of a reporting plugin
				query, err := reportingPlugin.Query(ctx, outcomeCtx)
				if err != nil {
					o.lggr.Errorf("failed to generate query: %v", err)
					return
				}
				observation, err := reportingPlugin.Observation(
					ctx,
					outcomeCtx,
					query,
				)

				if err != nil {
					o.lggr.Errorf("failed to generate observation: %v", err)
					return
				}

				attributedObservation := types.AttributedObservation{
					Observation: observation,
					Observer:    commontypes.OracleID(1),
				}
				if err = reportingPlugin.ValidateObservation(ctx, outcomeCtx, query, attributedObservation); err != nil {
					o.lggr.Errorf("failed to validate observation: %v", err)
					return
				}

				// Duplicate the observation N times; happy path expectation.
				// More complex cases can be added later.
				attributedObservations := make([]types.AttributedObservation, config.N)
				n := uint8(math.Min(float64(config.N), 4))
				for i := range n {
					attributedObservations[i] = types.AttributedObservation{
						Observation: observation,
						Observer:    commontypes.OracleID(i),
					}
				}

				newOutcome, err := reportingPlugin.Outcome(ctx, outcomeCtx, query, attributedObservations)
				if err != nil {
					o.lggr.Errorf("failed to generate outcome: %v", err)
					return
				}

				reports, err := reportingPlugin.Reports(ctx, outcomeCtx.SeqNr, newOutcome)
				if err != nil {
					o.lggr.Errorf("failed to generate reports: %v", err)
					return
				}

				for _, report := range reports {
					shouldAccept, err := reportingPlugin.ShouldAcceptAttestedReport(
						ctx,
						outcomeCtx.SeqNr,
						report.ReportWithInfo,
					)
					if err != nil {
						o.lggr.Errorf("failed when checking if should accept attested report: %v", err)
						return
					}

					if !shouldAccept {
						continue
					}

					shouldTransmit, err := reportingPlugin.ShouldTransmitAcceptedReport(
						ctx,
						outcomeCtx.SeqNr,
						report.ReportWithInfo,
					)

					if err != nil {
						o.lggr.Errorf("failed when checking if should transmit accepted report: %v", err)
						return
					}

					if !shouldTransmit {
						continue
					}

					err = o.config.ContractTransmitter.Transmit(
						ctx,
						types.ConfigDigest{},
						outcomeCtx.SeqNr,
						report.ReportWithInfo,
						[]types.AttributedOnchainSignature{},
					)
					if err != nil {
						o.lggr.Errorf("failed to transmit report: %v", err)
						return
					}
				}

				// Progress rounds
				outcomeCtx.SeqNr++
				outcomeCtx.PreviousOutcome = newOutcome
			}
		}
	}()

	return nil
}

func (o *oracle) Close(ctx context.Context) error {
	if o.cancel != nil {
		o.cancel()
	}
	o.wg.Wait() // Wait for the goroutine to finish
	return nil
}
