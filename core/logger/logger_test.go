package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log/noop"
	"go.uber.org/zap/zapcore"
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

func TestSetOtelCore(t *testing.T) {
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
			// Use TestLoggerObserved to create a logger that we can access
			logger, _ := TestLoggerObserved(t, zapcore.InfoLevel)
			require.NotNil(t, logger)

			if tc.enableOtel {
				// Create a no-op OTel logger for testing
				noopLogger := noop.NewLoggerProvider().Logger("test")

				// Enable OTel integration via SetOtelCore
				// Check if we can access OtelCore through the interface
				if otelCore, ok := logger.(OtelCore); ok {
					otelCore.SetOtelCore(noopLogger)
					// Test that logger works with otel core
					logger.Info("test log message with otel")
				} else {
					t.Skip("Logger does not implement OtelCore interface")
				}
			} else {
				// Test that regular logger works
				logger.Info("test log message without otel")
			}

			// Test that the logger was created successfully
			assert.NotNil(t, logger)
		})
	}
}
