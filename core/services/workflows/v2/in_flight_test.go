package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInFlight(t *testing.T) {
	t.Parallel()

	t.Run("counts entered executions until they exit", func(t *testing.T) {
		t.Parallel()
		var f inFlight
		require.True(t, f.enter())
		require.True(t, f.enter())
		assert.Equal(t, int32(2), f.count())
		f.exit()
		assert.Equal(t, int32(1), f.count())
		f.exit()
		assert.Equal(t, int32(0), f.count())
	})

	t.Run("closeAndWait returns at once when idle and refuses new executions", func(t *testing.T) {
		t.Parallel()
		var f inFlight
		f.closeAndWait()
		assert.False(t, f.enter())
		assert.Equal(t, int32(0), f.count())
	})

	t.Run("closeAndWait refuses new executions and blocks until running ones exit", func(t *testing.T) {
		t.Parallel()
		var f inFlight
		require.True(t, f.enter())

		done := make(chan struct{})
		go func() {
			f.closeAndWait()
			close(done)
		}()

		require.Eventually(t, func() bool {
			f.mu.Lock()
			defer f.mu.Unlock()
			return f.closed
		}, time.Second, time.Millisecond)
		assert.False(t, f.enter(), "a closing tracker must refuse new executions")
		require.Never(t, func() bool {
			select {
			case <-done:
				return true
			default:
				return false
			}
		}, 50*time.Millisecond, time.Millisecond, "closeAndWait must block while an execution is running")

		f.exit()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("closeAndWait did not return after the last execution exited")
		}
	})
}
