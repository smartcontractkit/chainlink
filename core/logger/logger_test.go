package logger

import (
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

			atomicCore := NewAtomicCore()
			logger, closeFn := cfg.NewWithCores()
			defer closeFn()
			require.NotNil(t, logger)

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
