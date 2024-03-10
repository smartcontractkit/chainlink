package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var errTest = errors.New("error")

func TestLogger_Passthrough(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		create func(passthrough Logger) Logger
	}{
		{"prometheus", newPrometheusLogger},
		{"sentry", newSentryLogger},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			m := setupMockLogger(t)
			l := test.create(m)

			l.With()
			l.Named("xxx")
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

			nm := l.Name()
			require.Equal(t, "mockLogger", nm)

			err := l.Sync()
			assert.ErrorIs(t, err, errTest)

			l.Recover(errTest)
		})
	}
}

func setupMockLogger(t *testing.T) *MockLogger {
	ml := NewMockLogger(t)

	ml.On("Helper", 1).Return(ml).Once()
	ml.On("With", mock.Anything, mock.Anything).Return(ml)
	ml.On("Named", "xxx").Return(ml).Once()
	ml.On("SetLogLevel", zapcore.DebugLevel).Once()

	ml.On("Trace").Once()
	ml.On("Debug").Once()
	ml.On("Info").Once()
	ml.On("Warn").Once()
	ml.On("Error").Once()
	ml.On("Critical").Once()
	ml.On("Panic").Once()
	ml.On("Fatal").Once()

	ml.On("Tracef", "msg").Once()
	ml.On("Debugf", "msg").Once()
	ml.On("Infof", "msg").Once()
	ml.On("Warnf", "msg").Once()
	ml.On("Errorf", "msg").Once()
	ml.On("Criticalf", "msg").Once()
	ml.On("Panicf", "msg").Once()
	ml.On("Fatalf", "msg").Once()

	ml.On("Tracew", "msg").Once()
	ml.On("Debugw", "msg").Once()
	ml.On("Infow", "msg").Once()
	ml.On("Warnw", "msg").Once()
	ml.On("Errorw", "msg", mock.Anything, mock.Anything).Once()
	ml.On("Criticalw", "msg", mock.Anything, mock.Anything).Once()
	ml.On("Panicw", "msg", mock.Anything, mock.Anything).Once()
	ml.On("Fatalw", "msg", mock.Anything, mock.Anything).Once()

	ml.On("Sync").Return(errTest).Once()
	ml.On("Name").Return("mockLogger").Once()
	ml.On("Recover", errTest).Once()

	return ml
}
