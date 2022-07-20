//go:build !trace

package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTrace(t *testing.T) {
	lgr, observered := TestLoggerObserved(t, zapcore.InfoLevel)

	const (
		testName    = "TestTrace"
		testMessage = "Trace message"
	)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message		logger/test_logger_test.go:23    logger=TestLogger
	require.Empty(t, observered.TakeAll())

	lgr.SetLogLevel(zapcore.DebugLevel)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message		logger/test_logger_test.go:23    logger=TestLogger
	require.Empty(t, observered.TakeAll())
}
