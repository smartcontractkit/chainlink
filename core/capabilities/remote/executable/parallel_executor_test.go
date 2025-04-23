package executable

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

func Test_CancellingContext_StopsTask(t *testing.T) {
	tp := newParallelExecutor(10)
	servicetest.Run(t, tp)

	var cancelFns []context.CancelFunc

	var counter int32
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancelFns = append(cancelFns, cancel)
		err := tp.ExecuteTask(ctx, func(ctx context.Context) {
			atomic.AddInt32(&counter, 1)
			<-ctx.Done()
			atomic.AddInt32(&counter, -1)
		})

		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&counter) == 10
	}, 10*time.Second, 10*time.Millisecond)

	for _, cancel := range cancelFns {
		cancel()
	}

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&counter) == 0
	}, 10*time.Second, 10*time.Millisecond)
}

func Test_ExecuteRequestTimesOutWhenParallelExecutionLimitReached(t *testing.T) {
	tp := newParallelExecutor(3)
	servicetest.Run(t, tp)

	for i := 0; i < 3; i++ {
		err := tp.ExecuteTask(context.Background(), func(ctx context.Context) {
			<-ctx.Done()
		})
		require.NoError(t, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	err := tp.ExecuteTask(ctx, func(ctx context.Context) {
	})
	cancel()
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func Test_ExecutingMultipleTasksInParallel(t *testing.T) {
	tp := newParallelExecutor(10)
	servicetest.Run(t, tp)

	var counter int32
	for i := 0; i < 10; i++ {
		err := tp.ExecuteTask(context.Background(), func(ctx context.Context) {
			atomic.AddInt32(&counter, 1)
			<-ctx.Done()
			atomic.AddInt32(&counter, -1)
		})
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&counter) == 10
	}, 10*time.Second, 10*time.Millisecond)
}

func Test_StopsExecutingMultipleParallelTasksWhenClosed(t *testing.T) {
	tp := newParallelExecutor(10)
	var counter int32
	t.Cleanup(func() {
		assert.Equal(t, int32(0), atomic.LoadInt32(&counter))
	})

	servicetest.Run(t, tp)

	for i := 0; i < 10; i++ {
		err := tp.ExecuteTask(context.Background(), func(ctx context.Context) {
			atomic.AddInt32(&counter, 1)
			<-ctx.Done()
			atomic.AddInt32(&counter, -1)
		})
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&counter) == 10
	}, 10*time.Second, 10*time.Millisecond)
}
