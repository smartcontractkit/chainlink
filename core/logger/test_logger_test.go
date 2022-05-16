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
	lgr := TestLogger(t)
	lgr.SetLogLevel(zapcore.InfoLevel)
	requireContains := func(cs ...string) {
		t.Helper()
		logs := MemoryLogTestingOnly().String()
		for _, c := range cs {
			require.Contains(t, logs, c)
		}
	}

	const (
		testName    = "TestTestLogger"
		testMessage = "Test message"
	)
	lgr.Warn(testMessage)
	// [WARN]  Test message		logger/test_logger_test.go:23    logger=1.0.0@sHaValue.TestLogger
	requireContains("[WARN]", testMessage, fmt.Sprintf("logger=%s.%s", verShaNameStatic(), testName))

	const (
		serviceName    = "ServiceName"
		serviceMessage = "Service message"
		key, value     = "key", "value"
		omittedMessage = "Don't log me"
	)
	srvLgr := lgr.Named(serviceName)
	srvLgr.SetLogLevel(zapcore.DebugLevel)
	srvLgr.Debugw(serviceMessage, key, value)
	// [DEBUG]  Service message		logger/test_logger_test.go:35    key=value logger=1.0.0@sHaValue.TestLogger.ServiceName
	requireContains("[DEBUG]", serviceMessage, fmt.Sprintf("%s=%s", key, value),
		fmt.Sprintf("logger=%s.%s.%s", verShaNameStatic(), testName, serviceName))

	const (
		workerName           = "WorkerName"
		workerMessage        = "Did some work"
		idKey, workerId      = "workerId", "42"
		resultKey, resultVal = "result", "success"
	)
	wrkLgr := srvLgr.Named(workerName).With(idKey, workerId)
	wrkLgr.Infow(workerMessage, resultKey, resultVal)
	// [INFO]	Did some work		logger/test_logger_test.go:49    logger=1.0.0@sHaValue.TestLogger.ServiceName.WorkerName result=success workerId=42
	requireContains("[INFO]", workerMessage, fmt.Sprintf("%s=%s", idKey, workerId),
		fmt.Sprintf("%s=%s", resultKey, resultVal), fmt.Sprintf("logger=%s.%s.%s.%s", verShaNameStatic(), testName, serviceName, workerName))

	const (
		critMsg = "Critical error"
	)
	lgr.Critical(critMsg)
	requireContains("[CRIT]", critMsg, fmt.Sprintf("logger=%s.%s", verShaNameStatic(), testName))
}
