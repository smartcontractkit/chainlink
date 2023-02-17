package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
			logger: &other{zaptest.NewLogger(t).Sugar(), ""},
		},
		{
			name:   "different",
			logger: &different{zaptest.NewLogger(t).Sugar(), ""},
		},
		{
			name:    "missing",
			logger:  &mismatch{zaptest.NewLogger(t).Sugar(), ""},
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

func TestNamed(t *testing.T) {
	prod, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		logger       Logger
		expectedName string
	}{
		{
			expectedName: "test.test1",
			logger:       Test(t).Named("test").Named("test1"),
		},
		{
			expectedName: "nop.nested",
			logger:       Nop().Named("nop").Named("nested"),
		},
		{
			expectedName: "prod",
			logger:       prod.Named("prod"),
		},
		{
			expectedName: "",
			logger:       &other{zaptest.NewLogger(t).Sugar(), ""},
		},
	} {
		t.Run(fmt.Sprintf("test_logger_name_expect_%s", tt.expectedName), func(t *testing.T) {
			require.Equal(t, tt.expectedName, tt.logger.Name())
		})
	}

}

type other struct {
	*zap.SugaredLogger
	name string
}

func (o *other) With(args ...interface{}) Logger {
	return &other{o.SugaredLogger.With(args...), ""}
}

func (o *other) Name() string {
	return o.name
}

func (o *other) Named(name string) Logger {
	newLogger := *o
	newLogger.name = joinName(o.name, name)
	newLogger.SugaredLogger = o.SugaredLogger.Named(name)
	return &newLogger
}

type different struct {
	*zap.SugaredLogger
	name string
}

func (d *different) With(args ...interface{}) differentLogger {
	return &different{d.SugaredLogger.With(args...), ""}
}

func (d *different) Name() string {
	return d.name
}

func (d *different) Named(name string) Logger {
	newLogger := *d
	newLogger.name = joinName(d.name, name)
	newLogger.SugaredLogger = d.SugaredLogger.Named(name)
	return &newLogger
}

type mismatch struct {
	*zap.SugaredLogger
	name string
}

func (m *mismatch) With(args ...interface{}) interface{} {
	return &mismatch{m.SugaredLogger.With(args...), ""}
}

func (m *mismatch) Name() string {
	return m.name
}

func (m *mismatch) Named(name string) Logger {
	newLogger := *m
	newLogger.name = joinName(m.name, name)
	newLogger.SugaredLogger = m.SugaredLogger.Named(name)
	return &newLogger
}

type differentLogger interface {
	Name() string
	Named(string) Logger

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
