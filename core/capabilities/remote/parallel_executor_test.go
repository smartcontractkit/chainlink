package remote_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote"
)

func Test_CancellingContext_StopsTask(t *testing.T) {
	t.Parallel()
	tp := remote.NewParallelExecutor(10, "test_parallel_executor")
	servicetest.Run(t, tp)

	cancelFns := make([]context.CancelFunc, 0, 10)

	var counter atomic.Int32
	for range 10 {
		ctx, cancel := context.WithCancel(t.Context())
		t.Cleanup(cancel)
		cancelFns = append(cancelFns, cancel)
		err := tp.ExecuteTask(ctx, func(ctx context.Context) {
			counter.Add(1)
			<-ctx.Done()
			counter.Add(-1)
		})

		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool { return counter.Load() == 10 }, 5*time.Second, 10*time.Millisecond)

	for _, cancel := range cancelFns {
		cancel()
	}

	assert.Eventually(t, func() bool { return counter.Load() == 0 }, 5*time.Second, 10*time.Millisecond)
}

func Test_ExecuteRequestTimesOutWhenParallelExecutionLimitReached(t *testing.T) {
	t.Parallel()
	tp := remote.NewParallelExecutor(3, "test_parallel_executor")
	servicetest.Run(t, tp)

	for range 3 {
		err := tp.ExecuteTask(t.Context(), func(ctx context.Context) {
			<-ctx.Done()
		})
		require.NoError(t, err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	err := tp.ExecuteTask(ctx, func(ctx context.Context) {
	})
	cancel()
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func Test_ExecutingMultipleTasksInParallel(t *testing.T) {
	t.Parallel()
	tp := remote.NewParallelExecutor(10, "test_parallel_executor")
	servicetest.Run(t, tp)

	var counter atomic.Int32
	for range 10 {
		err := tp.ExecuteTask(t.Context(), func(ctx context.Context) {
			counter.Add(1)
			<-ctx.Done()
			counter.Add(-1)
		})
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool { return counter.Load() == 10 }, 5*time.Second, 10*time.Millisecond)
}

func Test_StopsExecutingMultipleParallelTasksWhenClosed(t *testing.T) {
	t.Parallel()
	tp := remote.NewParallelExecutor(10, "test_parallel_executor")
	var counter atomic.Int32
	t.Cleanup(func() {
		assert.Eventually(t, func() bool { return counter.Load() == 0 }, 5*time.Second, 10*time.Millisecond)
	})

	servicetest.Run(t, tp)

	for range 10 {
		err := tp.ExecuteTask(t.Context(), func(ctx context.Context) {
			counter.Add(1)
			<-ctx.Done()
			counter.Add(-1)
		})
		require.NoError(t, err)
	}

	assert.Eventually(t, func() bool { return counter.Load() == 10 }, 5*time.Second, 10*time.Millisecond)
}
