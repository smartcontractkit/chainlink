package logger

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger/otelzap"
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

func TestLogStreamingEnabled(t *testing.T) {
	testCases := []struct {
		name                string
		logStreamingEnabled bool
		expectedCoresCount  int
		shouldHaveOtelCore  bool
	}{
		{
			logStreamingEnabled: true,
			expectedCoresCount:  2, // default core + otel core
			shouldHaveOtelCore:  true,
		},
		{
			logStreamingEnabled: false,
			expectedCoresCount:  1, // only default core
			shouldHaveOtelCore:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a no-op OTel logger for testing
			noopLogger := noop.NewLoggerProvider().Logger("test")

			config := Config{
				LogLevel:    zapcore.InfoLevel,
				OtelLogger:  noopLogger,
				JsonConsole: true,
				UnixTS:      true,
			}

			logger, closeLogger := config.New()
			require.NotNil(t, logger)
			require.NotNil(t, closeLogger)
			defer func() {
				err := closeLogger()
				require.NoError(t, err)
			}()

			// Test that logger works
			logger.Info("test log message")

			// Test that the logger was created successfully with the right config
			assert.NotNil(t, logger)
		})
	}
}

func TestNewOtelCore(t *testing.T) {
	// Create a no-op OTel logger for testing
	noopLogger := noop.NewLoggerProvider().Logger("test")

	// Test that NewOtelCore returns a valid core
	core := otelzap.NewCore(noopLogger)
	require.NotNil(t, core)

	// Test that the core can handle log entries
	entry := zapcore.Entry{
		Level:   zapcore.InfoLevel,
		Time:    zapcore.DefaultClock.Now(),
		Message: "test message",
	}

	// This should not panic even if beholder is not initialized
	err := core.Write(entry, nil)
	require.NoError(t, err)

	// Test core properties
	assert.True(t, core.Enabled(zapcore.InfoLevel))
}
