package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestWith(t *testing.T) {
	prod, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name    string
		logger  Logger
		expSame bool
	}{
		{
			name:   "test",
			logger: Test(t),
		},
		{
			name:   "nop",
			logger: Nop(),
		},
		{
			name:   "prod",
			logger: prod,
		},
		{
			name:   "other",
			logger: &other{zaptest.NewLogger(t).Sugar()},
		},
		{
			name:   "different",
			logger: &different{zaptest.NewLogger(t).Sugar()},
		},
		{
			name:    "missing",
			logger:  &mismatch{zaptest.NewLogger(t).Sugar()},
			expSame: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := With(tt.logger, "foo", "bar")
			same := got == tt.logger
			if same && !tt.expSame {
				t.Error("expected a new logger with foo==bar, but got same")
			} else if tt.expSame && !same {
				t.Errorf("expected the same logger %v, w/o foo=bar, but got %v", tt.logger, got)
			}
		})
	}

}

type other struct{ *zap.SugaredLogger }

func (o *other) With(args ...interface{}) Logger {
	return &other{o.SugaredLogger.With(args...)}
}

type different struct{ *zap.SugaredLogger }

func (d *different) With(args ...interface{}) differentLogger {
	return &different{d.SugaredLogger.With(args...)}
}

type mismatch struct{ *zap.SugaredLogger }

func (m *mismatch) With(args ...interface{}) interface{} {
	return &mismatch{m.SugaredLogger.With(args...)}
}

type differentLogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Panicf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Sync() error
}
