package oraclelib

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func NewChainAgnosticBackFilledOracle(lggr logger.Logger, srcProvider services.ServiceCtx, dstProvider services.ServiceCtx, oracle job.ServiceCtx) *ChainAgnosticBackFilledOracle {
	return &ChainAgnosticBackFilledOracle{
		srcProvider: srcProvider,
		dstProvider: dstProvider,
		oracle:      oracle,
		lggr:        lggr,
		stopCh:      make(chan struct{}),
		done:        make(chan struct{}),
	}
}

type ChainAgnosticBackFilledOracle struct {
	srcProvider   services.ServiceCtx
	dstProvider   services.ServiceCtx
	oracle        job.ServiceCtx
	lggr          logger.Logger
	oracleStarted atomic.Bool
	stopCh        commonservices.StopChan
	done          chan struct{}
}

func (r *ChainAgnosticBackFilledOracle) Start(_ context.Context) error {
	go r.run()
	return nil
}

func (r *ChainAgnosticBackFilledOracle) run() {
	defer close(r.done)
	ctx, cancel := r.stopCh.NewCtx()
	defer cancel()

	var err error
	var errMu sync.Mutex
	var wg sync.WaitGroup
	// Replay in parallel if both requested.
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		srcReplayErr := r.srcProvider.Start(ctx)
		errMu.Lock()
		err = multierr.Combine(err, srcReplayErr)
		errMu.Unlock()
		r.lggr.Infow("finished replaying src chain", "time", time.Since(s))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s := time.Now()
		dstReplayErr := r.dstProvider.Start(ctx)
		errMu.Lock()
		err = multierr.Combine(err, dstReplayErr)
		errMu.Unlock()
		r.lggr.Infow("finished replaying dst chain", "time", time.Since(s))
	}()

	wg.Wait()
	if err != nil {
		r.lggr.Criticalw("unexpected error replaying, continuing plugin boot without all the logs backfilled", "err", err)
	}
	if err := ctx.Err(); err != nil {
		r.lggr.Errorw("context already cancelled", "err", err)
	}
	// Start oracle with all logs present from dstStartBlock on dst and
	// all logs from srcStartBlock on src.
	if err := r.oracle.Start(ctx); err != nil {
		// Should never happen.
		r.lggr.Errorw("unexpected error starting oracle", "err", err)
	} else {
		r.oracleStarted.Store(true)
	}
}

func (r *ChainAgnosticBackFilledOracle) Close() error {
	close(r.stopCh)
	<-r.done
	if r.oracleStarted.Load() {
		// If the oracle is running, it must be Closed/stopped
		// TODO: Close should be safe to call in either case?
		if err := r.oracle.Close(); err != nil {
			r.lggr.Errorw("unexpected error stopping oracle", "err", err)
			return err
		}
		// Flag the oracle as closed with our internal variable that keeps track
		// of its state.  This will allow to re-start the process
		r.oracleStarted.Store(false)
	}
	return nil
}
