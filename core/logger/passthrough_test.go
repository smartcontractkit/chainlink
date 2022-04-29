package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
	"go.uber.org/zap/zapcore"
)

type TestingLogger interface {
	prometheusLogger | sentryLogger
}

var errTest error = errors.New("error")

func TestLogger_Passthrough(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		create func(t *testing.T, passthrough Logger) Logger
	}{
		{"prometheus", createTestLogger[prometheusLogger]},
		{"sentry", createTestLogger[sentryLogger]},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			m := setupMockLogger()
			l := test.create(t, m)

			l.With()
			l.Named("xxx")
			l.NewRootLogger(zapcore.DebugLevel)
			l.SetLogLevel(zapcore.DebugLevel)

			l.Trace()
			l.Debug()
			l.Info()
			l.Warn()
			l.Error()
			l.Critical()
			l.Panic()
			l.Fatal()

			l.Tracef("msg")
			l.Debugf("msg")
			l.Infof("msg")
			l.Warnf("msg")
			l.Errorf("msg")
			l.Criticalf("msg")
			l.Panicf("msg")
			l.Fatalf("msg")

			l.Tracew("msg")
			l.Debugw("msg")
			l.Infow("msg")
			l.Warnw("msg")
			l.Errorw("msg")
			l.Criticalw("msg")
			l.Panicw("msg")
			l.Fatalw("msg")

			err := l.Sync()
			assert.ErrorIs(t, err, errTest)

			l.Recover(errTest)

			ok := m.AssertExpectations(t)
			assert.True(t, ok)
		})
	}
}

func createTestLogger[TL TestingLogger](t *testing.T, passthrough Logger) Logger {
	var ret TL
	switch any(&ret).(type) {
	case *prometheusLogger:
		return newPrometheusLogger(passthrough)
	case *sentryLogger:
		return newSentryLogger(passthrough)
	}
	t.Fatal("unsupported logger")
	return nil
}

func setupMockLogger() *MockLogger {
	ml := &MockLogger{}

	ml.On("Helper", 1).Return(ml).Once()
	ml.On("With", mock.Anything, mock.Anything).Return(ml)
	ml.On("Named", "xxx").Return(ml).Once()
	ml.On("NewRootLogger", zapcore.DebugLevel).Return(ml, nil).Once()
	ml.On("SetLogLevel", zapcore.DebugLevel).Once()

	ml.On("Trace").Once()
	ml.On("Debug").Once()
	ml.On("Info").Once()
	ml.On("Warn").Once()
	ml.On("Error").Once()
	ml.On("Critical").Once()
	ml.On("Panic").Once()
	ml.On("Fatal").Once()

	ml.On("Tracef", mock.Anything).Once()
	ml.On("Debugf", mock.Anything).Once()
	ml.On("Infof", mock.Anything).Once()
	ml.On("Warnf", mock.Anything).Once()
	ml.On("Errorf", mock.Anything).Once()
	ml.On("Criticalf", mock.Anything).Once()
	ml.On("Panicf", mock.Anything).Once()
	ml.On("Fatalf", mock.Anything).Once()

	ml.On("Tracew", mock.Anything).Once()
	ml.On("Debugw", mock.Anything).Once()
	ml.On("Infow", mock.Anything).Once()
	ml.On("Warnw", mock.Anything).Once()
	ml.On("Errorw", mock.Anything, mock.Anything, mock.Anything).Once()
	ml.On("Criticalw", mock.Anything, mock.Anything, mock.Anything).Once()
	ml.On("Panicw", mock.Anything, mock.Anything, mock.Anything).Once()
	ml.On("Fatalw", mock.Anything, mock.Anything, mock.Anything).Once()

	ml.On("Sync").Return(errTest).Once()
	ml.On("Recover", errTest).Once()

	return ml
}
