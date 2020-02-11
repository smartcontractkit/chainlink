package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTestLogger(t *testing.T) {
	logger := CreateTestLogger(zapcore.DebugLevel)
	msg := "this is a test of the logging system"
	logger.Debug(msg)
	require.Contains(t, TestMemoryLog().String(), msg)
}
