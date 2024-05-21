//go:build !trace

package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTrace(t *testing.T) {
	lgr, observed := TestLoggerObserved(t, zapcore.InfoLevel)

	const (
		testName    = "TestTrace"
		testMessage = "Trace message"
	)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message		logger/test_logger_test.go:23    logger=TestLogger
	require.Empty(t, observed.TakeAll())

	lgr.SetLogLevel(zapcore.DebugLevel)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message		logger/test_logger_test.go:23    logger=TestLogger
	require.Empty(t, observed.TakeAll())
}
