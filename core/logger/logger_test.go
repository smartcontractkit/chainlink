package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTestLogger(t *testing.T) {
	logger := CreateTestLogger(zapcore.DebugLevel)
	logger.Warn("this is a log")
	require.Contains(t, MemoryLogTestingOnly().String(), "this is a log")
}
