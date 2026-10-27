package logger

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	InitColor(false)
}

// callerLine returns the line number of the statement that called it, so tests can
// assert on caller reporting without hard-coding a line that edits above them shift.
func callerLine(t *testing.T) int {
	t.Helper()
	_, _, line, ok := runtime.Caller(1)
	require.True(t, ok)
	return line
}

func TestTestLogger(t *testing.T) {
	t.Parallel()
	lgr, observed := TestLoggerObserved(t, zapcore.DebugLevel)

	const (
		testMessage = "Test message"
	)
	lgr.Warn(testMessage)
	// [WARN]  Test message		logger/test_logger_test.go:<line>	version=unset@unset
	logs := observed.TakeAll()
	require.Len(t, logs, 1)
	log := logs[0]
	assert.Equal(t, zap.WarnLevel, log.Level)
	assert.Equal(t, testMessage, log.Message)
	assert.Empty(t, log.LoggerName)
	ver := log.ContextMap()["version"]
	assert.Contains(t, ver, "@")

	const (
		serviceName    = "ServiceName"
		serviceMessage = "Service message"
		key, value     = "key", "value"
	)
	srvLgr := lgr.Named(serviceName)
	srvLgr.SetLogLevel(zapcore.DebugLevel)
	expectedLine := callerLine(t)
	srvLgr.Debugw(serviceMessage, key, value)
	// [DEBUG]  Service message		logger/test_logger_test.go:<line> version=unset@unset key=value logger=ServiceName
	logs = observed.TakeAll()
	require.Len(t, logs, 1)
	log = logs[0]
	assert.Equal(t, zap.DebugLevel, log.Level)
	assert.Equal(t, serviceMessage, log.Message)
	assert.Equal(t, serviceName, log.LoggerName)
	assert.Equal(t, value, log.ContextMap()[key])
	assert.Contains(t, log.Caller.String(), "core/logger/test_logger_test.go")
	assert.Equal(t, expectedLine+1, log.Caller.Line)

	const (
		workerName           = "WorkerName"
		workerMessage        = "Did some work"
		idKey, workerID      = "workerId", "42"
		resultKey, resultVal = "result", "success"
	)
	workLgr := srvLgr.Named(workerName).With(idKey, workerID)
	workLgr.Infow(workerMessage, resultKey, resultVal)
	// [INFO]	Did some work		logger/test_logger_test.go:<line> version=unset@unset logger=ServiceName.WorkerName result=success workerId=42
	logs = observed.TakeAll()
	require.Len(t, logs, 1)
	log = logs[0]
	assert.Equal(t, zap.InfoLevel, log.Level)
	assert.Equal(t, workerMessage, log.Message)
	assert.Equal(t, fmt.Sprintf("%s.%s", serviceName, workerName), log.LoggerName)
	assert.Equal(t, workerID, log.ContextMap()[idKey])
	assert.Equal(t, resultVal, log.ContextMap()[resultKey])

	const (
		critMsg = "Critical error"
	)
	lgr.Critical(critMsg)
	logs = observed.TakeAll()
	require.Len(t, logs, 1)
	log = logs[0]
	assert.Equal(t, zap.DPanicLevel, log.Level)
	assert.Equal(t, critMsg, log.Message)
	assert.Empty(t, log.LoggerName)
}
