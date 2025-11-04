package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log/noop"
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

// TestAtomicCoreSwap tests the atomic core swap functionality after logger creation.
func TestAtomicCoreSwap(t *testing.T) {
	atomicCore := NewAtomicCore()
	setOtelCore := atomicCore.Store

	lggrCfg := Config{
		LogLevel:       zapcore.InfoLevel,
		JsonConsole:    true,
		UnixTS:         false,
		FileMaxSizeMB:  0,
		FileMaxAgeDays: 0,
		FileMaxBackups: 0,
		SentryEnabled:  false,
	}

	lggr, closeFn := lggrCfg.NewWithCores(atomicCore)
	defer func() {
		err := closeFn()
		require.NoError(t, err)
	}()

	// Create observer to capture logs
	otelCore, otelLogs := observer.New(zapcore.InfoLevel)

	lggr.Info("before swap")

	assert.Equal(t, 0, otelLogs.Len(), "Expected no logs before core swap")

	// Swap to the observer core
	setOtelCore(otelCore)

	lggr.Info("after swap")

	assert.Equal(t, 1, otelLogs.Len(), "Expected 1 log after core swap")
	assert.Equal(t, "after swap", otelLogs.All()[0].Message)
	assert.Equal(t, zapcore.InfoLevel, otelLogs.All()[0].Level)

}
