package syncer

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sync(t *testing.T) {
	// Calling sync errors with not implemented
	var (
		lggr = logger.TestLogger(t)
		ctx  = context.Background()
		orm  = NewUnimplementedDS()
		wr   = NewWorkflowRegistry(lggr, orm)
		err  = wr.Sync(ctx, true)
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "not implemented")
}

func Test_StartSyncStop(t *testing.T) {
	var (
		lggr        = logger.TestLogger(t)
		ctx, cancel = context.WithCancel(context.Background())
		orm         = NewUnimplementedDS()
		wr          = NewWorkflowRegistry(lggr, orm)
		giveTicker  = make(chan struct{})
		gotErrs     = make([]error, 0)
		wantErrs    = 2 // initial sync + one tick
	)
	t.Cleanup(cancel)

	// set a forcing ticker
	wr.tickerCh = giveTicker

	// go read errors from the syncer
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case syncErr, open := <-wr.errCh:
				if !open {
					return
				}
				gotErrs = append(gotErrs, syncErr)
			}
		}
	}()

	// start the syncer
	err := wr.Start(ctx)
	require.NoError(t, err)

	// force one additional sync
	giveTicker <- struct{}{}

	// eventually read two errors from the syncer
	delta := maxTestTime(t)
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Len(c, gotErrs, wantErrs)
	}, delta, delta/10, "length of errors not two")

	require.NoError(t, wr.Close())
}

func maxTestTime(t *testing.T) time.Duration {
	now := time.Now()
	d, ok := t.Deadline()
	if !ok {
		return 5 * time.Second
	}
	delta := d.Sub(now)
	return delta
}
