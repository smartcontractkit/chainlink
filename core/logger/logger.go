// Package logger is used to store details of events in the node.
// Events can be categorized by Trace, Debug, Info, Error, Fatal, and Panic.
package logger

import (
	"log"
	"reflect"
	"runtime"

	"gorm.io/gorm"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
type Logger interface {
	// Named creates a new named logger with the given name
	Named(name string) Logger
	// With creates a new logger with the given arguments
	With(args ...interface{}) Logger
	// WithDB creates a new logger with the given db.
	WithDB(db *gorm.DB) Logger

	// GetORM returns the underlying ORM, if one is set.
	GetORM() ORM

	// InitServiceLevelLogger builds a service level logger with a given logging level & serviceName
	InitServiceLevelLogger(serviceName string, logLevel string) (Logger, error)
	// ServiceLogLevel is the log level set for a specified package
	ServiceLogLevel(serviceName string) (string, error)
	// GetServiceLogLevels retrieves all service log levels from the db
	GetServiceLogLevels() (map[string]string, error)

	Info(args ...interface{})
	Infof(format string, values ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Debug(args ...interface{})
	Debugf(format string, values ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Warn(args ...interface{})
	Warnf(format string, values ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	// WarnIf logs the error if present.
	WarnIf(err error)

	Error(args ...interface{})
	Errorf(format string, values ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	// ErrorIf logs the error if present.
	ErrorIf(err error, optionalMsg ...string)
	// ErrorIfCalling calls fn and logs any returned error along with func name.
	ErrorIfCalling(fn func() error, optionalMsg ...string)

	Fatal(values ...interface{})
	Fatalf(format string, values ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Panic(args ...interface{})
	PanicIf(err error)

	Sync() error

	// withCallerSkip creates a new logger with the number of callers skipped by
	// caller annotation increased by add. For wrappers to use internally.
	withCallerSkip(add int) Logger
}

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	Orm         ORM
	lvl         zapcore.Level
	dir         string
	jsonConsole bool
	toDisk      bool
	fields      []interface{}
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

func (l *zapLogger) Write(b []byte) (int, error) {
	l.Info(string(b))
	return len(b), nil
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

func (l *zapLogger) Named(name string) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Named(name).With("id", name)
	newLogger.fields = copyFields(l.fields, "id", name)
	return &newLogger
}

func (l *zapLogger) withCallerSkip(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return &newLogger
}

func (l *zapLogger) WarnIf(err error) {
	if err != nil {
		l.Warn(err)
	}
}

func (l *zapLogger) ErrorIf(err error, optionalMsg ...string) {
	if err != nil {
		if len(optionalMsg) > 0 {
			l.Error(errors.Wrap(err, optionalMsg[0]))
		} else {
			l.Error(err)
		}
	}
}

func (l *zapLogger) ErrorIfCalling(fn func() error, optionalMsg ...string) {
	err := fn()
	if err != nil {
		fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		e := errors.Wrap(err, fnName)
		if len(optionalMsg) > 0 {
			l.Error(errors.Wrap(e, optionalMsg[0]))
		} else {
			l.Error(e)
		}
	}
}

func (l *zapLogger) PanicIf(err error) {
	if err != nil {
		l.Panic(err)
	}
}

func (l *zapLogger) WithDB(db *gorm.DB) Logger {
	newLogger := *l
	newLogger.Orm = NewORM(db)
	return &newLogger
}

func (l *zapLogger) GetORM() ORM {
	return l.Orm
}

func (l *zapLogger) GetServiceLogLevels() (map[string]string, error) {
	serviceLogLevels := make(map[string]string)

	for _, svcName := range GetLogServices() {
		svc, err := l.ServiceLogLevel(svcName)
		if err != nil {
			Fatalf("error getting service log levels: %v", err)
		}
		serviceLogLevels[svcName] = svc
	}

	return serviceLogLevels, nil
}

// initLogConfig builds a zap.Config for a logger
func initLogConfig(dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) zap.Config {
	config := zap.NewProductionConfig()
	if !jsonConsole {
		config.OutputPaths = []string{"pretty://console"}
	}
	if toDisk {
		destination := logFileURI(dir)
		config.OutputPaths = append(config.OutputPaths, destination)
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, destination)
	}
	config.Level.SetLevel(lvl)
	return config
}

// CreateProductionLogger returns a log config for the passed directory
// with the given LogLevel and customizes stdout for pretty printing.
func CreateProductionLogger(
	dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) Logger {
	config := initLogConfig(dir, jsonConsole, lvl, toDisk)

	zl, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
	return &zapLogger{
		SugaredLogger: zl.Sugar(),
		lvl:           lvl,
		dir:           dir,
		jsonConsole:   jsonConsole,
		toDisk:        toDisk,
	}
}

func (l *zapLogger) InitServiceLevelLogger(serviceName string, logLevel string) (Logger, error) {
	var ll zapcore.Level
	if err := ll.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	config := initLogConfig(l.dir, l.jsonConsole, ll, l.toDisk)

	zl, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	newLogger := *l
	newLogger.SugaredLogger = zl.Named(serviceName).Sugar().With(l.fields...)
	newLogger.fields = copyFields(l.fields)
	return &newLogger, nil
}

func (l *zapLogger) ServiceLogLevel(serviceName string) (string, error) {
	if l.Orm != nil {
		level, err := l.Orm.GetServiceLogLevel(serviceName)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			Warnf("Error while trying to fetch %s service log level: %v", serviceName, err)
		} else if err == nil {
			return level, nil
		}
	}
	return l.lvl.String(), nil
}

// NewProductionConfig returns a production logging config
func NewProductionConfig(lvl zapcore.Level, dir string, jsonConsole, toDisk bool) (c zap.Config) {
	var outputPath string
	if jsonConsole {
		outputPath = "stderr"
	} else {
		outputPath = "pretty://console"
	}
	// Mostly copied from zap.NewProductionConfig with sampling disabled
	c = zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Development:      false,
		Sampling:         nil,
		Encoding:         "json",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{outputPath},
		ErrorOutputPaths: []string{"stderr"},
	}
	if toDisk {
		destination := logFileURI(dir)
		c.OutputPaths = append(c.OutputPaths, destination)
		c.ErrorOutputPaths = append(c.ErrorOutputPaths, destination)
	}
	return
}

// NewProductionEncoderConfig returns a production encoder config
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	// Copied from zap.NewProductionEncoderConfig but with ISO timestamps instead of Unix
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
