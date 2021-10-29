// Logger is the main interface of this package.
//
// The package-level helper functions are being phased out. Loggers should be injected
// instead (and usually Named as well): e.g. lggr.Named("<service name>")
//
// Tests should use a TestLogger, with NewLogger being reserved for actual
// runtime and limited direct testing.
package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var envLvl = zapcore.InfoLevel

func init() {
	_ = envLvl.Set(os.Getenv("LOG_LEVEL"))
}

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
// TestLogger should be used in tests.
type Logger interface {
	// With creates a new Logger with the given arguments
	With(args ...interface{}) Logger
	// Named creates a new Logger sub-scoped with name.
	// Names are inherited and dot-separated.
	//   a := l.Named("a") // logger=a
	//   b := a.Named("b") // logger=a.b
	Named(name string) Logger

	// NewRootLogger creates a new root Logger with an independent log level
	// unaffected by upstream calls to SetLogLevel.
	NewRootLogger(lvl zapcore.Level) (Logger, error)

	// SetLogLevel changes the log level for this and all connected Loggers.
	SetLogLevel(zapcore.Level)

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	// WarnIf logs the error if present.
	WarnIf(err error, msg string)
	// ErrorIf logs the error if present.
	ErrorIf(err error, msg string)
	PanicIf(err error, msg string)

	// ErrorIfCalling calls fn and logs any returned error along with func name.
	ErrorIfCalling(fn func() error)

	// Sync flushes any buffered log entries.
	// Some insignificant errors are suppressed.
	Sync() error

	// withCallerSkip creates a new logger with the number of callers skipped by
	// caller annotation increased by add. For wrappers to use internally.
	withCallerSkip(add int) Logger
}

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	config zap.Config
	name   string
	fields []interface{}
}

func newZapLogger(cfg zap.Config) (Logger, error) {
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &zapLogger{config: cfg, SugaredLogger: zl.Sugar()}, nil
}

func (l *zapLogger) SetLogLevel(lvl zapcore.Level) {
	l.config.Level.SetLevel(lvl)
}

// Constants for service names for package specific logging configuration
const (
	HeadTracker = "head_tracker"
	FluxMonitor = "fluxmonitor"
	Keeper      = "keeper"
)

func GetLogServices() []string {
	return []string{
		HeadTracker,
		FluxMonitor,
		Keeper,
	}
}

func (l *zapLogger) With(args ...interface{}) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.With(args...)
	newLogger.fields = copyFields(l.fields, args...)
	return &newLogger
}

// copyFields returns a copy of fields with add appended.
func copyFields(fields []interface{}, add ...interface{}) []interface{} {
	f := make([]interface{}, 0, len(fields)+len(add))
	f = append(f, fields...)
	f = append(f, add...)
	return f
}

func joinName(old, new string) string {
	if old == "" {
		return new
	}
	return old + "." + new
}

func (l *zapLogger) Named(name string) Logger {
	newLogger := *l
	newLogger.name = joinName(l.name, name)
	newLogger.SugaredLogger = l.SugaredLogger.Named(name)
	return &newLogger
}

func (l *zapLogger) NewRootLogger(lvl zapcore.Level) (Logger, error) {
	newLogger := *l
	newLogger.config.Level = zap.NewAtomicLevelAt(lvl)
	zl, err := newLogger.config.Build()
	if err != nil {
		return nil, err
	}
	newLogger.SugaredLogger = zl.Named(l.name).Sugar().With(l.fields...)
	return &newLogger, nil
}

func (l *zapLogger) withCallerSkip(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return &newLogger
}

func (l *zapLogger) WarnIf(err error, msg string) {
	if err != nil {
		l.withCallerSkip(1).Warnw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.withCallerSkip(1).Errorw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIfCalling(fn func() error) {
	err := fn()
	if err != nil {
		fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		l.withCallerSkip(1).Errorw(fmt.Sprintf("Error calling %s", fnName), "err", err)
	}
}

func (l *zapLogger) PanicIf(err error, msg string) {
	if err != nil {
		l.withCallerSkip(1).Panicw(msg, "err", err)
	}
}

func (l *zapLogger) Sync() error {
	err := l.SugaredLogger.Sync()
	if err == nil {
		return nil
	}
	var msg string
	if uw := errors.Unwrap(err); uw != nil {
		msg = uw.Error()
	} else {
		msg = err.Error()
	}
	switch msg {
	case os.ErrInvalid.Error(), "bad file descriptor",
		"inappropriate ioctl for device":
		return nil
	}
	return err
}

// newProductionConfig returns a new production zap.Config.
func newProductionConfig(dir string, jsonConsole bool, toDisk bool, unixTS bool) zap.Config {
	config := newBaseConfig()
	if !unixTS {
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	if !jsonConsole {
		config.OutputPaths = []string{"pretty://console"}
	}
	if toDisk {
		destination := logFileURI(dir)
		config.OutputPaths = append(config.OutputPaths, destination)
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, destination)
	}
	return config
}

type Config interface {
	RootDir() string
	JSONConsole() bool
	LogToDisk() bool
	LogLevel() zapcore.Level
	LogUnixTimestamps() bool
}

// NewLogger returns a new Logger configured by c with pretty printing to stdout.
// If LogToDisk is false, the Logger will only log to stdout.
// Tests should use TestLogger instead.
func NewLogger(c Config) Logger {
	cfg := newProductionConfig(c.RootDir(), c.JSONConsole(), c.LogToDisk(), c.LogUnixTimestamps())
	cfg.Level.SetLevel(c.LogLevel())
	l, err := newZapLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

// InitColor explicitly sets the global color.NoColor option.
// Not safe for concurrent use. Only to be called from init().
func InitColor(c bool) {
	color.NoColor = !c
}

func newBaseConfig() zap.Config {
	// Copied from zap.NewProductionConfig with sampling disabled
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:      false,
		Sampling:         nil,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
