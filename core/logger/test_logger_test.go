package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	InitColor(false)
}

func TestTestLogger(t *testing.T) {
	logger := CreateTestLogger(t)
	logger.Warn("this is a log")
	require.Contains(t, MemoryLogTestingOnly().String(), "this is a log")
	require.Contains(t, MemoryLogTestingOnly().String(), "logger=TestTestLogger")
}
