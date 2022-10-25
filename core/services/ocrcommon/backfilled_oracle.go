package ocrcommon

import (
	"context"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

type BackfilledPoller struct {
	FromBlock uint64
	Poller    logpoller.LogPoller
}

type BackfilledOracle struct {
	backfilledPollers []BackfilledPoller
	oracle            job.ServiceCtx
	lggr              logger.Logger
}

func NewBackfilledOracle(lggr logger.Logger, backfilledPollers []BackfilledPoller, oracle job.ServiceCtx) *BackfilledOracle {
	return &BackfilledOracle{
		backfilledPollers: backfilledPollers,
		oracle:            oracle,
		lggr:              lggr,
	}
}

func (r *BackfilledOracle) Start(ctx context.Context) error {
	go func() {
		var err error
		var errMu sync.Mutex
		var wg sync.WaitGroup
		for _, backfilledPoller := range r.backfilledPollers {
			if backfilledPoller.FromBlock != 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s := time.Now()
					r.lggr.Infow("start replaying chain", "fromBlock", backfilledPoller.FromBlock)
					srcReplayErr := backfilledPoller.Poller.Replay(context.Background(), int64(backfilledPoller.FromBlock))
					errMu.Lock()
					err = multierr.Combine(err, srcReplayErr)
					errMu.Unlock()
					r.lggr.Infow("finished replaying chain", "time", time.Since(s))
				}()
			}

		}
		wg.Wait()
		if err != nil {
			r.lggr.Errorw("unexpected error replaying", "err", err)
			return
		}
		// Start oracle with all logs present.
		if err := r.oracle.Start(ctx); err != nil {
			// Should never happen.
			r.lggr.Errorw("unexpected error starting oracle", "err", err)
		}
	}()
	return nil
}

func (r *BackfilledOracle) Close() error {
	return r.oracle.Close()
}
