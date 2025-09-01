package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
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

func TestConfig_New_WithAdditionalCores(t *testing.T) {
	// Create multiple test cores
	observedCore1, observedLogs1 := observer.New(zapcore.InfoLevel)
	observedCore2, observedLogs2 := observer.New(zapcore.DebugLevel)

	config := Config{
		LogLevel:        zapcore.DebugLevel,
		JsonConsole:     true,
		UnixTS:          true,
		AdditionalCores: []zapcore.Core{observedCore1, observedCore2},
	}

	logger, closeFn := config.New()
	defer func() {
		if err := closeFn(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()

	// Log messages at different levels
	logger.Info("info message")
	logger.Debug("debug message")

	// Verify both cores received appropriate messages
	assert.Equal(t, 1, observedLogs1.Len()) // Info level core should only get info message
	assert.Equal(t, "info message", observedLogs1.All()[0].Message)

	assert.Equal(t, 2, observedLogs2.Len()) // Debug level core should get both messages
	logs2 := observedLogs2.All()
	assert.Equal(t, "info message", logs2[0].Message)
	assert.Equal(t, "debug message", logs2[1].Message)
}

func TestConfig_New_WithoutAdditionalCores(t *testing.T) {
	config := Config{
		LogLevel:    zapcore.InfoLevel,
		JsonConsole: true,
		UnixTS:      true,
		// No AdditionalCores specified
	}

	logger, closeFn := config.New()
	defer func() {
		if err := closeFn(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()

	// Should work normally without additional cores
	assert.NotNil(t, logger)
	logger.Info("test message without additional cores")
}

func TestNewOtelCore(t *testing.T) {
	// Test that NewOtelCore returns a valid core
	core := NewOtelCore()
	assert.NotNil(t, core)

	// Test that the core can be used in a logger
	observedCore, observedLogs := observer.New(zapcore.InfoLevel)
	teeCore := zapcore.NewTee(core, observedCore)

	testLogger := zap.New(teeCore).Sugar()
	testLogger.Info("test otel core")

	// The observed core should capture the message
	assert.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, "test otel core", observedLogs.All()[0].Message)
}
