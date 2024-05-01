package capreg

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"go.uber.org/multierr"
)

var _ services.ServiceCtx = (*syncer)(nil)

type syncer struct {
	locals []Local
	lggr   logger.Logger
}

// HealthReport implements services.Service.
func (s *syncer) HealthReport() map[string]error {
	return map[string]error{} // TODO: implement
}

// Name implements services.Service.
func (s *syncer) Name() string {
	return "CapabilityRegistrySyncer"
}

// Ready implements services.Service.
func (s *syncer) Ready() error {
	return nil // TODO: implement
}

func NewSyncer(locals []Local, lggr logger.Logger) *syncer {
	return &syncer{
		locals: locals,
		lggr:   lggr.Named("capreg_syncer"),
	}
}

// Close implements services.Service.
func (s *syncer) Close() error {
	var errs error
	for _, local := range s.locals {
		if err := local.Close(); err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	return errs
}

// Start implements services.Service.
func (s *syncer) Start(ctx context.Context) error {
	tick := time.NewTicker(12 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			latestState := s.refreshOnchainState(ctx)
			for _, local := range s.locals {
				if err := local.Sync(ctx, latestState); err != nil {
					s.lggr.Errorw("failed to sync chain state to local state", "err", err)
				}
			}
		}
	}
}

// refreshOnchainState fetches the capability registry state from the blockchain.
func (s *syncer) refreshOnchainState(_ context.Context) State {
	return State{} // TODO: implement
}
