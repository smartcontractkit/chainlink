//go:build !trace

package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTrace(t *testing.T) {
	lgr := TestLogger(t)
	lgr.SetLogLevel(zapcore.InfoLevel)
	requireNotContains := func(ns ...string) {
		t.Helper()
		logs := MemoryLogTestingOnly().String()
		for _, n := range ns {
			require.NotContains(t, logs, n)
		}
	}

	const (
		testName    = "TestTrace"
		testMessage = "Trace message"
	)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message		logger/test_logger_test.go:23    logger=TestLogger
	requireNotContains("[TRACE]", testMessage, fmt.Sprintf("logger=%s", testName))

	lgr.SetLogLevel(zapcore.DebugLevel)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message		logger/test_logger_test.go:23    logger=TestLogger
	requireNotContains("[TRACE]", testMessage, fmt.Sprintf("logger=%s", testName))
}
