package utils

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThreadControl_Close(t *testing.T) {
	n := 10
	tc := NewThreadControl()

	finished := atomic.Int32{}

	for i := 0; i < n; i++ {
		tc.Go(func(ctx context.Context) {
			<-ctx.Done()
			finished.Add(1)
		})
	}

	tc.Close()

	require.Equal(t, int32(n), finished.Load())
}

func TestThreadControl_GoCtx(t *testing.T) {
	tc := NewThreadControl()
	defer tc.Close()

	var wg sync.WaitGroup
	finished := atomic.Int32{}

	timeout := 100 * time.Millisecond

	start := time.Now()
	ctx, cancel := context.WithDeadline(context.Background(), start.Add(timeout))
	defer cancel()

	wg.Add(1)
	tc.GoCtx(ctx, func(c context.Context) {
		defer wg.Done()
		<-c.Done()
		finished.Add(1)
	})

	wg.Wait()
	elapsed := time.Since(start)
	assert.GreaterOrEqual(t, elapsed, timeout)
	assert.Less(t, elapsed, 2*timeout)
	require.Equal(t, int32(1), finished.Load())
}
