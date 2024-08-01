package logging

import (
	"bufio"
	"bytes"

	testing "github.com/mitchellh/go-testing-interface"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TestLogger struct {
	instance *zap.Logger
	stdout   *bytes.Buffer
}

// Instance returns the actual *zap.Logger you should pass to your dependency
// to accumulate log lines and inspect them later on.
func (l *TestLogger) Instance() *zap.Logger {
	return l.instance
}

// RecordedLines returns all the logged line seen, each line being a full fledged JSON value
func (w *TestLogger) RecordedLines(t testing.T) (out []string) {
	scanner := bufio.NewScanner(w.stdout)
	for succeed, line := next(scanner); succeed && scanner.Err() == nil; succeed, line = next(scanner) {
		out = append(out, line)
	}

	if scanner.Err() != nil {
		t.Errorf("test logger scanning logged lines fail unexpectedly: %w", scanner.Err())
	}
	return
}

func next(scanner *bufio.Scanner) (succeed bool, line string) {
	succeed = scanner.Scan()
	line = scanner.Text()
	return
}

type bufferSyncer bytes.Buffer

func (b *bufferSyncer) Write(p []byte) (n int, err error) {
	return (*bytes.Buffer)(b).Write(p)
}

func (b *bufferSyncer) Sync() error {
	return nil
}

func NewTestLogger(t testing.T) *TestLogger {
	stdout := bytes.NewBuffer(nil)
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), (*bufferSyncer)(stdout), zap.DebugLevel)

	return &TestLogger{
		instance: zap.New(core),
		stdout:   stdout,
	}
}
