package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func init() {
	InitColor(false)
}

func TestTestLogger(t *testing.T) {
	lgr := CreateTestLogger(t)
	lgr.SetLogLevel(zapcore.DebugLevel)

	const (
		testName    = "TestTestLogger"
		testMessage = "Test message"
	)
	lgr.Warn(testMessage)
	// [WARN]  Test message		logger/test_logger_test.go:23    logger=TestTestLogger
	require.Contains(t, MemoryLogTestingOnly().String(), "[WARN]")
	require.Contains(t, MemoryLogTestingOnly().String(), testMessage)
	require.Contains(t, MemoryLogTestingOnly().String(), fmt.Sprintf("logger=%s", testName))

	const (
		serviceName    = "ServiceName"
		serviceMessage = "Service message"
		key, value     = "key", "value"
	)
	srvLgr := lgr.Named(serviceName)
	srvLgr.Infow(serviceMessage, key, value)
	// [INFO]  Service message		logger/test_logger_test.go:35    key=value logger=TestTestLogger.ServiceName
	require.Contains(t, MemoryLogTestingOnly().String(), "[INFO]")
	require.Contains(t, MemoryLogTestingOnly().String(), serviceMessage)
	require.Contains(t, MemoryLogTestingOnly().String(), fmt.Sprintf("%s=%s", key, value))
	require.Contains(t, MemoryLogTestingOnly().String(), fmt.Sprintf("logger=%s.%s", testName, serviceName))

	const (
		workerName           = "WorkerName"
		workerMessage        = "Did some work"
		idKey, workerId      = "workerId", "42"
		resultKey, resultVal = "result", "success"
	)
	wrkLgr := srvLgr.Named(workerName).With(idKey, workerId)
	wrkLgr.Debugw(workerMessage, resultKey, resultVal)
	// [DEBUG]	Did some work		logger/test_logger_test.go:49    logger=TestTestLogger.ServiceName.WorkerName result=success workerId=42
	require.Contains(t, MemoryLogTestingOnly().String(), "[DEBUG]")
	require.Contains(t, MemoryLogTestingOnly().String(), workerMessage)
	require.Contains(t, MemoryLogTestingOnly().String(), fmt.Sprintf("%s=%s", idKey, workerId))
	require.Contains(t, MemoryLogTestingOnly().String(), fmt.Sprintf("%s=%s", resultKey, resultVal))
	require.Contains(t, MemoryLogTestingOnly().String(), fmt.Sprintf("logger=%s.%s.%s", testName, serviceName, workerName))

}
