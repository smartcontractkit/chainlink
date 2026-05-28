package vault

import (
	"errors"
	"fmt"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestUserFacingError(t *testing.T) {
	t.Run("returns error message for userError", func(t *testing.T) {
		err := newUserError("key does not exist")
		assert.Equal(t, "key does not exist", userFacingError(err, "fallback"))
	})

	t.Run("returns fallback for non-userError", func(t *testing.T) {
		err := errors.New("internal failure")
		assert.Equal(t, "fallback msg", userFacingError(err, "fallback msg"))
	})

	t.Run("returns wrapped error message for wrapped userError", func(t *testing.T) {
		err := fmt.Errorf("context: %w", newUserError("bad key"))
		assert.Equal(t, "context: bad key", userFacingError(err, "fallback"))
	})
}

func TestLogUserErrorAware(t *testing.T) {
	t.Run("logs at debug level for userError", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := newUserError("key does not exist")

		logUserErrorAware(lggr, "failed to observe request", err, "id", "req-1")

		debugLogs := observed.FilterLevelExact(zapcore.DebugLevel)
		errorLogs := observed.FilterLevelExact(zapcore.ErrorLevel)
		assert.Equal(t, 1, debugLogs.Len())
		assert.Equal(t, 0, errorLogs.Len())

		entry := debugLogs.All()[0]
		assert.Equal(t, "failed to observe request", entry.Message)
		fields := entry.ContextMap()
		assert.Equal(t, "req-1", fields["id"])
		assert.Contains(t, fmt.Sprint(fields["error"]), "key does not exist")
	})

	t.Run("logs at error level for internal error", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := errors.New("database connection lost")

		logUserErrorAware(lggr, "failed to observe request", err, "id", "req-2")

		debugLogs := observed.FilterLevelExact(zapcore.DebugLevel)
		errorLogs := observed.FilterLevelExact(zapcore.ErrorLevel)
		assert.Equal(t, 0, debugLogs.Len())
		assert.Equal(t, 1, errorLogs.Len())

		entry := errorLogs.All()[0]
		assert.Equal(t, "failed to observe request", entry.Message)
		fields := entry.ContextMap()
		assert.Equal(t, "req-2", fields["id"])
		assert.Contains(t, fmt.Sprint(fields["error"]), "database connection lost")
	})

	t.Run("logs at debug level for wrapped userError", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := fmt.Errorf("validation: %w", newUserError("bad input"))

		logUserErrorAware(lggr, "request failed", err, "op", "create")

		debugLogs := observed.FilterLevelExact(zapcore.DebugLevel)
		errorLogs := observed.FilterLevelExact(zapcore.ErrorLevel)
		assert.Equal(t, 1, debugLogs.Len())
		assert.Equal(t, 0, errorLogs.Len())
	})

	t.Run("includes all key-value pairs in log entry", func(t *testing.T) {
		lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)
		err := errors.New("internal error")

		logUserErrorAware(lggr, "op failed", err, "id", "req-3", "requestID", "abc-123")

		entry := observed.FilterLevelExact(zapcore.ErrorLevel).All()[0]
		fields := entry.ContextMap()
		assert.Equal(t, "req-3", fields["id"])
		assert.Equal(t, "abc-123", fields["requestID"])
		assert.Contains(t, fmt.Sprint(fields["error"]), "internal error")
	})
}
