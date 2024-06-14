//go:build trace

package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTrace(t *testing.T) {
	lgr, observed := TestLoggerObserved(t, zapcore.DebugLevel)
	lgr.SetLogLevel(zapcore.InfoLevel)

	const (
		testMessage = "Trace message"
	)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message
	require.Empty(t, observed.TakeAll())

	lgr.SetLogLevel(zapcore.DebugLevel)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message
	logs := observed.TakeAll()
	require.Len(t, logs, 1)
	log := logs[0]
	require.Equal(t, zapcore.DebugLevel, log.Level)
	require.Equal(t, "[TRACE] "+testMessage, log.Message)
}
