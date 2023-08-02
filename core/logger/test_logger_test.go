package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	InitColor(false)
}

func TestTestLogger(t *testing.T) {
	lgr, observed := TestLoggerObserved(t, zapcore.DebugLevel)

	const (
		testMessage = "Test message"
	)
	lgr.Warn(testMessage)
	// [WARN]  Test message		logger/test_logger_test.go:23	version=unset@unset
	logs := observed.TakeAll()
	require.Len(t, logs, 1)
	log := logs[0]
	assert.Equal(t, zap.WarnLevel, log.Level)
	assert.Equal(t, testMessage, log.Message)
	assert.Equal(t, "", log.LoggerName)
	ver := log.ContextMap()["version"]
	assert.Contains(t, ver, "@")

	const (
		serviceName    = "ServiceName"
		serviceMessage = "Service message"
		key, value     = "key", "value"
	)
	srvLgr := lgr.Named(serviceName)
	srvLgr.SetLogLevel(zapcore.DebugLevel)
	srvLgr.Debugw(serviceMessage, key, value)
	// [DEBUG]  Service message		logger/test_logger_test.go:35 version=unset@unset key=value logger=ServiceName
	logs = observed.TakeAll()
	require.Len(t, logs, 1)
	log = logs[0]
	assert.Equal(t, zap.DebugLevel, log.Level)
	assert.Equal(t, serviceMessage, log.Message)
	assert.Equal(t, serviceName, log.LoggerName)
	assert.Equal(t, value, log.ContextMap()[key])
	assert.Contains(t, log.Caller.String(), "core/logger/test_logger_test.go")
	assert.Equal(t, log.Caller.Line, 41)

	const (
		workerName           = "WorkerName"
		workerMessage        = "Did some work"
		idKey, workerId      = "workerId", "42"
		resultKey, resultVal = "result", "success"
	)
	wrkLgr := srvLgr.Named(workerName).With(idKey, workerId)
	wrkLgr.Infow(workerMessage, resultKey, resultVal)
	// [INFO]	Did some work		logger/test_logger_test.go:49 version=unset@unset logger=ServiceName.WorkerName result=success workerId=42
	logs = observed.TakeAll()
	require.Len(t, logs, 1)
	log = logs[0]
	assert.Equal(t, zap.InfoLevel, log.Level)
	assert.Equal(t, workerMessage, log.Message)
	assert.Equal(t, fmt.Sprintf("%s.%s", serviceName, workerName), log.LoggerName)
	assert.Equal(t, workerId, log.ContextMap()[idKey])
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
	assert.Equal(t, "", log.LoggerName)
}
