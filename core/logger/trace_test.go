//go:build trace

package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTrace(t *testing.T) {
	lgr := TestLogger(t)
	lgr.SetLogLevel(zapcore.InfoLevel)
	requireContains := func(cs ...string) {
		t.Helper()
		logs := MemoryLogTestingOnly().String()
		for _, c := range cs {
			require.Contains(t, logs, c)
		}
	}
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
	// [DEBUG] [TRACE] Trace message
	requireNotContains(testMessage)

	lgr.SetLogLevel(zapcore.DebugLevel)
	lgr.Trace(testMessage)
	// [DEBUG] [TRACE] Trace message
	requireContains("[DEBUG]", "[TRACE]", testMessage)
}
