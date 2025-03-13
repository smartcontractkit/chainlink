package executable

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CancellingContext_StopsTask(t *testing.T) {
	tp := newParallelExecutor(10)
	defer tp.Close()

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
	defer tp.Close()

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
	defer tp.Close()

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

	tp.Close()

	assert.Eventually(t, func() bool {
		return atomic.LoadInt32(&counter) == 0
	}, 10*time.Second, 10*time.Millisecond)
}
