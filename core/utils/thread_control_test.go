package utils

import (
	"context"
	"sync/atomic"
	"testing"

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
