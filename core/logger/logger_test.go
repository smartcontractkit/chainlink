package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log/noop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink-common/pkg/logger/otelzap"
)

func TestConfig(t *testing.T) {
	// no sampling
	assert.Nil(t, newZapConfigBase().Sampling)
	assert.Nil(t, newZapConfigProd(false, false).Sampling)

	// not development, which would trigger panics for Critical level
	assert.False(t, newZapConfigBase().Development)
	assert.False(t, newZapConfigProd(false, false).Development)
}

func TestStderrWriter(t *testing.T) {
	sw := stderrWriter{}

	// Test Write
	n, err := sw.Write([]byte("Hello, World!"))
	require.NoError(t, err)
	assert.Equal(t, 13, n, "Expected 13 bytes written")

	// Test Close
	err = sw.Close()
	require.NoError(t, err)
}

func TestOtelCore(t *testing.T) {
	testCases := []struct {
		name       string
		enableOtel bool
	}{
		{
			name:       "otel integration enabled",
			enableOtel: true,
		},
		{
			name:       "otel integration disabled",
			enableOtel: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{
				LogLevel: zapcore.InfoLevel,
			}

			var logger Logger
			var closeFn func() error

			if tc.enableOtel {
				// Create a no-op OTel logger for testing
				noopLogger := noop.NewLoggerProvider().Logger("test")
				otelCore := otelzap.NewCore(noopLogger, otelzap.WithLevel(zapcore.DebugLevel))

				logger, closeFn = cfg.NewWithCores(otelCore)
				defer func() {
					err := closeFn()
					require.NoError(t, err)
				}()
				require.NotNil(t, logger)

				// Test that logger works with otel core
				logger.Info("test log message with otel")
			} else {
				// Test that regular logger works
				logger, closeFn = cfg.NewWithCores()
				defer func() {
					err := closeFn()
					require.NoError(t, err)
				}()
				require.NotNil(t, logger)

				logger.Info("test log message without otel")
			}

			// Test that the logger was created successfully
			assert.NotNil(t, logger)
		})
	}
}

func TestAtomicCoreSwap(t *testing.T) {
	// This test simulates two processes:
	// 1. Process 1 creates a logger with an AtomicCore (initially noop)
	// 2. Process 2 swaps in a new OTel core using SetOtelCore
	// 3. Verify that subsequent logs go to both original and new cores

	// Create test cores that capture log entries
	observedCore1, observedLogs1 := observer.New(zapcore.InfoLevel)
	observedCore2, observedLogs2 := observer.New(zapcore.InfoLevel)

	// Process 1: Create logger with AtomicCore
	atomicCore := NewAtomicCore()

	// Build the logger manually to have full control over cores
	zcfg := newZapConfigBase()
	zcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Create a tee with observedCore1 and atomicCore
	teeCore := zapcore.NewTee(observedCore1, atomicCore)

	// Create logger with the combined core
	errSink, _, err := zap.Open(zcfg.ErrorOutputPaths...)
	require.NoError(t, err)

	zapLog := zap.New(teeCore, zap.ErrorOutput(errSink), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger := &zapLogger{
		level:         zcfg.Level,
		SugaredLogger: zapLog.Sugar(),
	}

	// Log before swapping - should only go to observedCore1 (atomicCore is noop)
	logger.Info("before swap")
	assert.Equal(t, 1, observedLogs1.Len(), "Expected 1 log in observedCore1 before swap")
	assert.Equal(t, 0, observedLogs2.Len(), "Expected 0 logs in observedCore2 before swap")

	// Process 2: Swap in a new core (simulating SetOtelCore)
	// In production, this would be an OTel core, but we use observedCore2 for verification
	core2 := observedCore2
	atomicCore.Store(&core2)

	// Log after swapping - should go to both observedCore1 and observedCore2
	logger.Info("after swap")
	assert.Equal(t, 2, observedLogs1.Len(), "Expected 2 logs in observedCore1 after swap")
	assert.Equal(t, 1, observedLogs2.Len(), "Expected 1 log in observedCore2 after swap")

	// Verify the message in observedCore2
	entries := observedLogs2.All()
	require.Len(t, entries, 1)
	assert.Equal(t, "after swap", entries[0].Message)
	assert.Equal(t, zapcore.InfoLevel, entries[0].Level)

	// Test with different log levels
	logger.Debug("debug message")
	// Debug is below InfoLevel, so shouldn't be logged
	assert.Equal(t, 2, observedLogs1.Len(), "Debug should not be logged at Info level")
	assert.Equal(t, 1, observedLogs2.Len(), "Debug should not be logged at Info level")

	logger.Warn("warning message")
	assert.Equal(t, 3, observedLogs1.Len(), "Warn should be logged")
	assert.Equal(t, 2, observedLogs2.Len(), "Warn should be logged in both cores")

	// Verify the second message in observedCore2
	entries = observedLogs2.All()
	require.Len(t, entries, 2)
	assert.Equal(t, "warning message", entries[1].Message)
	assert.Equal(t, zapcore.WarnLevel, entries[1].Level)
}
