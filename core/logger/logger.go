// Package logger is used to store details of events in the node.
// Events can be categorized by Trace, Debug, Info, Error, Fatal, and Panic.
package logger

import (
	"log"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
type logger struct {
	*zap.SugaredLogger
	Orm         ORM
	lvl         zapcore.Level
	dir         string
	jsonConsole bool
	toDisk      bool
}

type ServiceName string

// Constants for service names for package specific logging configuration
var (
	HeadTracker          = ServiceName("head_tracker")
	FluxMonitor          = ServiceName("fluxmonitor")
	JobSubscriber        = ServiceName("job_subscriber")
	RunQueue             = ServiceName("run_queue")
	BalanceMonitor       = ServiceName("balance_monitor")
	TxManager            = ServiceName("tx_manager")
	GasUpdater           = ServiceName("gas_updater")
	HeadListener         = ServiceName("head_listener")
	HeadBroadcaster      = ServiceName("head_broadcaster")
	EventBroadcaster     = ServiceName("event_broadcaster")
	DatabaseBackup       = ServiceName("database_backup")
	PromReporter         = ServiceName("prom_reporter")
	SingletonPeerWrapper = ServiceName("singleton_peer_wrapper")
	OCRContractTracker   = ServiceName("ocr_contract_tracker")
	ExplorerClient       = ServiceName("explorer_client")
	StatsPusher          = ServiceName("stats_pusher")
)

// func GetLogServices() []ServiceName {
// 	return []ServiceName{HeadTracker, FluxMonitor}
// }

// Write logs a message at the Info level and returns the length
// of the given bytes.
func (l *logger) Write(b []byte) (int, error) {
	l.Info(string(b))
	return len(b), nil
}

// WarnIf logs the error if present.
func (l *logger) WarnIf(err error) {
	if err != nil {
		l.Warn(err)
	}
}

// ErrorIf logs the error if present.
func (l *logger) ErrorIf(err error, optionalMsg ...string) {
	if err != nil {
		if len(optionalMsg) > 0 {
			l.Error(errors.Wrap(err, optionalMsg[0]))
		} else {
			l.Error(err)
		}
	}
}

// ErrorIfCalling calls the given function and logs the error of it if there is.
func (l *logger) ErrorIfCalling(f func() error, optionalMsg ...string) {
	err := f()
	if err != nil {
		e := errors.Wrap(err, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		if len(optionalMsg) > 0 {
			Default.Error(errors.Wrap(e, optionalMsg[0]))
		} else {
			Default.Error(e)
		}
	}
}

func (l *logger) PanicIf(err error) {
	if err != nil {
		l.Panic(err)
	}
}

// func (l *logger) SetDB(db *gorm.DB) {
// 	l.Orm = NewORM(db)
// }

// GetServiceLogLevels retrieves all service log levels from the db
// func (l *logger) GetServiceLogLevels() (map[ServiceName]string, error) {
// 	serviceLogLevels := make(map[ServiceName]string)

// 	headTracker, err := l.ServiceLogLevel(HeadTracker)
// 	if err != nil {
// 		Fatalf("error getting service log levels: %v", err)
// 	}

// 	serviceLogLevels[HeadTracker] = headTracker

// 	fluxMonitor, err := l.ServiceLogLevel(FluxMonitor)
// 	if err != nil {
// 		Fatalf("error getting service log levels: %v", err)
// 	}

// 	serviceLogLevels[FluxMonitor] = fluxMonitor

// 	return serviceLogLevels, nil
// }

// CreateLogger dwisott
// func CreateLogger(zl *zap.SugaredLogger) Logger {
// 	return &logger{SugaredLogger: zl}
// }

// func CreateLoggerWithConfig(zl *zap.SugaredLogger, lvl zapcore.Level, dir string, jsonConsole bool, toDisk bool) Logger {
// 	return &logger{
// 		SugaredLogger: zl,
// 		lvl:           lvl,
// 		dir:           dir,
// 		jsonConsole:   jsonConsole,
// 		toDisk:        toDisk,
// 	}
// }

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
func CreateProductionLogger(dir string, jsonConsole bool, lvl zapcore.Level, toDisk bool) Logger {
	config := initLogConfig(dir, jsonConsole, lvl, toDisk)

	zl, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
	return &logger{
		SugaredLogger: zl.Sugar(),
		lvl:           lvl,
		dir:           dir,
		jsonConsole:   jsonConsole,
		toDisk:        toDisk,
	}
}

// InitServiceLevelLogger builds a service level logger with a given logging level & serviceName
// func (l Logger) InitServiceLevelLogger(serviceName ServiceName, logLevel string) (Logger, error) {
// 	var ll zapcore.Level
// 	if err := ll.UnmarshalText([]byte(logLevel)); err != nil {
// 		return nil, err
// 	}

// 	config := initLogConfig(l.dir, l.jsonConsole, ll, l.toDisk)

// 	zl, err := config.Build(zap.AddCallerSkip(1))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return CreateLoggerWithConfig(zl.Named(string(serviceName)).Sugar(), ll, l.dir, l.jsonConsole, l.toDisk), nil
// }

// ServiceLogLevel is the log level set for a specified package
// func (l *logger) ServiceLogLevel(serviceName ServiceName) (string, error) {
// 	if l.Orm != nil {
// 		level, err := l.Orm.GetServiceLogLevel(serviceName)
// 		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 			Warnf("Error while trying to fetch %s service log level: %v", serviceName, err)
// 		} else if err == nil {
// 			return level, nil
// 		}
// 	}
// 	return l.lvl.String(), nil
// }

func (l *logger) With(args ...interface{}) *logger {
	return &logger{
		l.SugaredLogger.With(args),
		l.Orm,
		l.lvl,
		l.dir,
		l.jsonConsole,
		l.toDisk,
	}
}

// NewProductionConfig returns a production logging config
// func NewProductionConfig(lvl zapcore.Level, dir string, jsonConsole, toDisk bool) (c zap.Config) {
// 	var outputPath string
// 	if jsonConsole {
// 		outputPath = "stderr"
// 	} else {
// 		outputPath = "pretty://console"
// 	}
// 	// Mostly copied from zap.NewProductionConfig with sampling disabled
// 	c = zap.Config{
// 		Level:            zap.NewAtomicLevelAt(lvl),
// 		Development:      false,
// 		Sampling:         nil,
// 		Encoding:         "json",
// 		EncoderConfig:    NewProductionEncoderConfig(),
// 		OutputPaths:      []string{outputPath},
// 		ErrorOutputPaths: []string{"stderr"},
// 	}
// 	if toDisk {
// 		destination := logFileURI(dir)
// 		c.OutputPaths = append(c.OutputPaths, destination)
// 		c.ErrorOutputPaths = append(c.ErrorOutputPaths, destination)
// 	}
// 	return
// }

// NewProductionEncoderConfig returns a production encoder config
// func NewProductionEncoderConfig() zapcore.EncoderConfig {
// 	// Copied from zap.NewProductionEncoderConfig but with ISO timestamps instead of Unix
// 	return zapcore.EncoderConfig{
// 		TimeKey:        "ts",
// 		LevelKey:       "level",
// 		NameKey:        "logger",
// 		CallerKey:      "caller",
// 		FunctionKey:    zapcore.OmitKey,
// 		MessageKey:     "msg",
// 		StacktraceKey:  "stacktrace",
// 		LineEnding:     zapcore.DefaultLineEnding,
// 		EncodeLevel:    zapcore.LowercaseLevelEncoder,
// 		EncodeTime:     zapcore.ISO8601TimeEncoder,
// 		EncodeDuration: zapcore.SecondsDurationEncoder,
// 		EncodeCaller:   zapcore.ShortCallerEncoder,
// 	}
// }
