package v2

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

type mockLimiter[T any] struct {
	limit T
}

func (m *mockLimiter[T]) Limit(ctx context.Context) (T, error)     { return m.limit, nil }
func (m *mockLimiter[T]) Check(ctx context.Context, value T) error { return nil }
func (m *mockLimiter[T]) Close() error                             { return nil }

func TestEngine_emitUserLogs_ProcessesLogsWhenContextCancelled(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		limiters := &EngineLimiters{
			LogEvent: &mockLimiter[int]{limit: 100},
			LogLine:  &mockLimiter[config.Size]{limit: 1000},
		}

		e := &Engine{
			cfg: &EngineConfig{
				DebugMode:     true,
				LocalLimiters: limiters,
			},
		}
		e.setLogger(logger.Sugared(logger.Test(t)))

		userLogChan := make(chan *protoevents.LogLine, 10)

		// Enqueue some logs before starting the engine
		userLogChan <- &protoevents.LogLine{Message: "log 1"}
		userLogChan <- &protoevents.LogLine{Message: "log 2"}
		userLogChan <- &protoevents.LogLine{Message: "log 3"}

		// Close the channel immediately so we can observe complete teardown
		close(userLogChan)

		// Create a context that is ALREADY cancelled.
		// The OLD approach with `select` would randomly (or always) pick <-ctx.Done()
		// and return before processing all logs from the closed userLogChan.
		// The NEW approach reads the channel until it's closed, ignoring ctx.Done().
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		done := make(chan struct{})
		go func() {
			e.emitUserLogs(ctx, userLogChan, "test-exec-id", nil)
			close(done)
		}()

		synctest.Wait()

		select {
		case <-done:
			// Success: the function should return once the channel is fully drained.
		default:
			t.Fatal("emitUserLogs hung, likely failed to tear down properly")
		}

		// Verify the channel was fully drained.
		// (Since it's closed, a read on a drained channel returns !ok)
		_, ok := <-userLogChan
		require.False(t, ok, "Expected userLogChan to be fully drained but it was not")
	})
}

func TestEngine_emitUserLogs_TeardownOnChannelClose(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		limiters := &EngineLimiters{
			LogEvent: &mockLimiter[int]{limit: 100},
			LogLine:  &mockLimiter[config.Size]{limit: 1000},
		}

		e := &Engine{
			cfg: &EngineConfig{
				DebugMode:     true,
				LocalLimiters: limiters,
			},
		}
		e.setLogger(logger.Sugared(logger.Test(t)))

		userLogChan := make(chan *protoevents.LogLine, 10)

		// Create a context that is NEVER cancelled.
		// This proves that closing the channel is sufficient to tear down the loop.
		ctx := context.Background()

		done := make(chan struct{})
		go func() {
			e.emitUserLogs(ctx, userLogChan, "test-exec-id", nil)
			close(done)
		}()

		// Send a log and then close the channel.
		userLogChan <- &protoevents.LogLine{Message: "log 1"}
		close(userLogChan)

		synctest.Wait()

		select {
		case <-done:
			// Success: the function should return once the channel is closed and drained.
		default:
			t.Fatal("emitUserLogs hung, failed to tear down properly when userLogChan was closed")
		}

		_, ok := <-userLogChan
		require.False(t, ok, "Expected userLogChan to be fully drained but it was not")
	})
}
