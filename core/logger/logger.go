package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/getsentry/sentry-go"
	"github.com/smartcontractkit/chainlink/core/static"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const SentryFlushDeadline = 5 * time.Second

var envLvl = zapcore.InfoLevel

func init() {
	_ = envLvl.Set(os.Getenv("LOG_LEVEL"))

	// If SENTRY_ENVIRONMENT is set, it will override everything. Otherwise infers from CHAINLINK_DEV.
	var sentryenv string
	if env := os.Getenv("SENTRY_ENVIRONMENT"); env != "" {
		sentryenv = env
	} else if os.Getenv("CHAINLINK_DEV") == "true" {
		sentryenv = "dev"
	} else {
		sentryenv = "prod"
	}
	// If SENTRY_DSN is set, it will override everything. Otherwise static.SentryDSN will be used.
	// If neither are set, sentry is disabled.
	var sentrydsn string
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		sentrydsn = dsn
	} else {
		sentrydsn = static.SentryDSN
	}
	// If SENTRY_RELEASE is set, it will override everything. Otherwise, static.Version will be used.
	var sentryrelease string
	if release := os.Getenv("SENTRY_RELEASE"); release != "" {
		sentryrelease = release
	} else {
		sentryrelease = static.Version
	}
	err := sentry.Init(sentry.ClientOptions{
		// AttachStacktrace is needed to send stacktrace alongside panics
		AttachStacktrace: true,
		Dsn:              sentrydsn,
		Environment:      sentryenv,
		Release:          sentryrelease,
		// Enable printing of SDK debug messages.
		// Uncomment line below to debug sentry
		// Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
//
// The package-level helper functions are being phased out. Loggers should be injected
// instead (and usually Named as well): e.g. lggr.Named("<service name>")
//
// Tips
//  - Tests should use a TestLogger, with NewLogger being reserved for actual
//    runtime and limited direct testing.
//  - Critical level logs should only be used when user intervention is required.
//  - Trace level logs are omitted unless compiled with the trace tag. For example: go test -tags trace ...

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

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Critical(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Tracef(format string, values ...interface{})
	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Criticalf(format string, values ...interface{})
	Panicf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})

	Tracew(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	CriticalW(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	// ErrorIf logs the error if present.
	ErrorIf(err error, msg string)

	// ErrorIfClosing calls c.Close() and logs any returned error along with name.
	ErrorIfClosing(c io.Closer, name string)

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
	newLogger.SugaredLogger = l.sugaredWithCallerSkip(skip)
	return &newLogger
}

func (l *zapLogger) sugaredWithCallerSkip(skip int) *zap.SugaredLogger {
	return l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
}

func (l *zapLogger) ErrorIf(err error, msg string) {
	if err != nil {
		sentry.CaptureException(err)
		l.withCallerSkip(1).Errorw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		sentry.CaptureException(err)
		l.withCallerSkip(1).Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
	}
}

func (l *zapLogger) Error(args ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelError)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	l.sugaredWithCallerSkip(1).Error(args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	l.sugaredWithCallerSkip(1).Fatal(args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"args": args,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf("%v", args))
	l.sugaredWithCallerSkip(1).Panic(args...)
}

func (l *zapLogger) Errorf(format string, values ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelError)
	})
	hub.CaptureMessage(fmt.Sprintf(format, values...))
	l.sugaredWithCallerSkip(1).Errorf(format, values...)
}

func (l *zapLogger) Fatalf(format string, values ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", map[string]interface{}{
			"values": values,
		})
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(fmt.Sprintf(format, values...))
	l.sugaredWithCallerSkip(1).Fatalf(format, values...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelError)
	})
	hub.CaptureMessage(msg)
	l.sugaredWithCallerSkip(1).Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(msg)
	l.sugaredWithCallerSkip(1).Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	defer sentry.Flush(SentryFlushDeadline)
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("logger", toMap(keysAndValues))
		scope.SetLevel(sentry.LevelFatal)
	})
	hub.CaptureMessage(msg)
	l.sugaredWithCallerSkip(1).Panicw(msg, keysAndValues...)
}

func toMap(args ...interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args); {
		// Make sure this element isn't a dangling key
		if i == len(args)-1 {
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string ignore it
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); ok {
			m[keyStr] = val
		}
		i += 2
	}
	return m
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
	return newLogger(c.LogLevel(), c.RootDir(), c.JSONConsole(), c.LogToDisk(), c.LogUnixTimestamps())
}

func newLogger(logLevel zapcore.Level, dir string, jsonConsole bool, toDisk bool, unixTS bool) Logger {
	cfg := newProductionConfig(dir, jsonConsole, toDisk, unixTS)
	cfg.Level.SetLevel(logLevel)
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

// newBaseConfig returns a zap.NewProductionConfig with sampling disabled and a modified level encoder.
func newBaseConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.EncoderConfig.EncodeLevel = encodeLevel
	return cfg
}
