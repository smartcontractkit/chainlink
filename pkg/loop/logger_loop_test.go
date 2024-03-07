package loop_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
)

func TestHCLogLoggerPanic(t *testing.T) {
	lggr, ol := logger.TestObservedSugared(t, zapcore.DebugLevel)

	type testCase struct {
		name                string
		level               int
		expectedMessage     string
		expectedCustomKey   string
		expectedCustomValue string
		expectedLogLevel    zapcore.Level
	}

	tests := []testCase{
		{
			level:               test.PANIC,
			expectedMessage:     "[PANIC] panic: random panic",
			expectedCustomKey:   "",
			expectedCustomValue: "",
			expectedLogLevel:    zapcore.DPanicLevel,
		},
		{
			level:               test.FATAL,
			expectedMessage:     "[FATAL] some panic log",
			expectedCustomKey:   "custom-name-panic",
			expectedCustomValue: "custom-value-panic",
			expectedLogLevel:    zapcore.DPanicLevel,
		},
		{
			level:               test.CRITICAL,
			expectedMessage:     "some critical error log",
			expectedCustomKey:   "custom-name-critical",
			expectedCustomValue: "custom-value-critical",
			expectedLogLevel:    zapcore.DPanicLevel,
		}, {
			level:               test.ERROR,
			expectedMessage:     "some error log",
			expectedCustomKey:   "custom-name-error",
			expectedCustomValue: "custom-value-error",
			expectedLogLevel:    zapcore.ErrorLevel,
		},
		{
			level:               test.INFO,
			expectedMessage:     "some info log",
			expectedCustomKey:   "custom-name-info",
			expectedCustomValue: "custom-value-info",
			expectedLogLevel:    zapcore.InfoLevel,
		},
		{
			level:               test.WARN,
			expectedMessage:     "some warn log",
			expectedCustomKey:   "custom-name-warn",
			expectedCustomValue: "custom-value-warn",
			expectedLogLevel:    zapcore.WarnLevel,
		},
		{
			level:               test.DEBUG,
			expectedMessage:     "some debug log",
			expectedCustomKey:   "custom-name-debug",
			expectedCustomValue: "custom-value-debug",
			expectedLogLevel:    zapcore.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerTest := &test.GRPCPluginLoggerTest{SugaredLogger: lggr}
			cc := loggerTest.ClientConfig()
			cc.Cmd = NewHelperProcessCommand(test.PluginLoggerTestName, false, tt.level)
			c := plugin.NewClient(cc)
			_, err := c.Client()
			require.NoError(t, err)
			time.Sleep(time.Second * 2) //wait for log to sync
			c.Kill()
			fLogs := ol.FilterMessage(tt.expectedMessage)
			logs := fLogs.TakeAll()
			require.Equal(t, len(logs), 1, fmt.Sprintf("could not find expected log %q", tt.expectedMessage))
			require.Equal(t, tt.expectedMessage, logs[0].Message)
			require.Equal(t, tt.expectedLogLevel, logs[0].Level)
			if tt.expectedCustomKey != "" {
				found := false
				for _, e := range logs[0].Context {
					if e.Key == tt.expectedCustomKey && e.String == tt.expectedCustomValue {
						found = true
						break
					}
				}
				require.True(t, found, fmt.Sprintf("could not find expected values %s=%s", tt.expectedCustomKey, tt.expectedCustomValue))
			}
		})
	}
}
