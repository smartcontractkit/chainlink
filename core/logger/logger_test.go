package logger

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log/noop"
	"go.uber.org/zap/zapcore"

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
			logger, closeFn, atomicCore := cfg.NewWithAtomicCore()
			defer closeFn()
			require.NotNil(t, logger)
			require.NotNil(t, atomicCore)

			if tc.enableOtel {
				// Create a no-op OTel logger for testing
				noopLogger := noop.NewLoggerProvider().Logger("test")

				otelCore := otelzap.NewCore(noopLogger, otelzap.WithLevel(zapcore.DebugLevel))
				atomicCore.Store(&otelCore)
				// Test that logger works with otel core
				logger.Info("test log message with otel")
			} else {
				// Test that regular logger works
				logger.Info("test log message without otel")
			}

			// Test that the logger was created successfully
			assert.NotNil(t, logger)
		})
	}
}

func TestNewWithAtomicCore_RotatingFileLogger(t *testing.T) {
	// Test that rotating file logger works with AtomicCore approach
	tempDir := t.TempDir()

	cfg := Config{
		LogLevel:       zapcore.InfoLevel,
		Dir:            tempDir,
		FileMaxSizeMB:  1, // Enable rotating file logger
		FileMaxAgeDays: 1,
		FileMaxBackups: 1,
	}

	logger, closeFn, atomicCore := cfg.NewWithAtomicCore()
	defer closeFn()

	require.NotNil(t, logger)
	require.NotNil(t, atomicCore)

	// Test that the logger works
	logger.Info("test message for rotating file logger")

	// Test that we can add OTel core to the AtomicCore
	noopLogger := noop.NewLoggerProvider().Logger("test")
	otelCore := otelzap.NewCore(noopLogger, otelzap.WithLevel(zapcore.DebugLevel))
	atomicCore.Store(&otelCore)

	// Test logging with OTel integration
	logger.Info("test message with otel integration")

	// Verify log file was created
	logFile := filepath.Join(tempDir, "chainlink_debug.log")
	require.FileExists(t, logFile)
}
