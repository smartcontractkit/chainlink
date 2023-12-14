package oraclelib

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

type BackfilledOracle struct {
	srcStartBlock, dstStartBlock uint64
	oracleStarted                atomic.Bool
	cancelFn                     context.CancelFunc
	src, dst                     logpoller.LogPoller
	oracle                       job.ServiceCtx
	lggr                         logger.Logger
}

func NewBackfilledOracle(lggr logger.Logger, src, dst logpoller.LogPoller, srcStartBlock, dstStartBlock uint64, oracle job.ServiceCtx) *BackfilledOracle {
	return &BackfilledOracle{
		srcStartBlock: srcStartBlock,
		dstStartBlock: dstStartBlock,
		oracleStarted: atomic.Bool{},
		cancelFn:      nil,
		src:           src,
		dst:           dst,
		oracle:        oracle,
		lggr:          lggr,
	}
}

func (r *BackfilledOracle) Start(_ context.Context) error {
	go r.Run()
	return nil
}

func (r *BackfilledOracle) IsRunning() bool {
	return r.oracleStarted.Load()
}

func (r *BackfilledOracle) Run() {
	ctx, cancelFn := context.WithCancel(context.Background())
	r.cancelFn = cancelFn
	var err error
	var errMu sync.Mutex
	var wg sync.WaitGroup
	// Replay in parallel if both requested.
	if r.srcStartBlock != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := time.Now()
			r.lggr.Infow("start replaying src chain", "fromBlock", r.srcStartBlock)
			srcReplayErr := r.src.Replay(ctx, int64(r.srcStartBlock))
			errMu.Lock()
			err = multierr.Combine(err, srcReplayErr)
			errMu.Unlock()
			r.lggr.Infow("finished replaying src chain", "time", time.Since(s))
		}()
	}
	if r.dstStartBlock != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := time.Now()
			r.lggr.Infow("start replaying dst chain", "fromBlock", r.dstStartBlock)
			dstReplayErr := r.dst.Replay(ctx, int64(r.dstStartBlock))
			errMu.Lock()
			err = multierr.Combine(err, dstReplayErr)
			errMu.Unlock()
			r.lggr.Infow("finished replaying  dst chain", "time", time.Since(s))
		}()
	}
	wg.Wait()
	if err != nil {
		r.lggr.Errorw("unexpected error replaying", "err", err)
		return
	}
	if err := ctx.Err(); err != nil {
		r.lggr.Errorw("context already cancelled", "err", err)
		return
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

func (r *BackfilledOracle) Close() error {
	if r.oracleStarted.Load() {
		// If the oracle is running, it must be Closed/stopped
		if err := r.oracle.Close(); err != nil {
			r.lggr.Errorw("unexpected error stopping oracle", "err", err)
			return err
		}
		// Flag the oracle as closed with our internal variable that keeps track
		// of its state.  This will allow to re-start the process
		r.oracleStarted.Store(false)
	}
	if r.cancelFn != nil {
		// This is useful to step the previous tasks that are spawned in
		// parallel before starting the Oracle. This will use the context to
		// signal them to exit immediately.
		//
		// It can be possible this is the only way to stop the Start() async
		// flow, specially when the previusly task are running (the replays) and
		// `oracleStarted` would be false in that example. Calling `cancelFn()`
		// will stop the replays and will prevent the oracle to start
		r.cancelFn()
	}
	return nil
}
