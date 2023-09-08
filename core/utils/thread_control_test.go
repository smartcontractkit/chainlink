package utils

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThreadControl_Close(t *testing.T) {
	n := 10
	tc := NewThreadControl(n)

	finished := atomic.Int32{}

	for i := 0; i < n; i++ {
		require.NoError(t, tc.Go(func(ctx context.Context) {
			<-ctx.Done()
			finished.Add(1)
		}))
	}

	tc.Close()

	require.Equal(t, int32(n), finished.Load())
}
func TestThreadControl_ThreadsLimitExceeded(t *testing.T) {
	tc := NewThreadControl(1)

	finished := atomic.Int32{}

	fn := func(ctx context.Context) {
		<-ctx.Done()
		finished.Add(1)
	}
	require.NoError(t, tc.Go(fn))
	require.Error(t, tc.Go(fn))

	tc.Close()

	require.Equal(t, int32(1), finished.Load())
}
