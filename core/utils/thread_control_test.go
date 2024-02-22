package utils

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	tc.GoCtx(ctx, func(c context.Context) {
		defer wg.Done()
		<-c.Done()
		finished.Add(1)
	})

	go func() {
		time.After(1 * time.Millisecond)
		cancel()
	}()

	wg.Wait()
	require.Equal(t, int32(1), finished.Load())
}
