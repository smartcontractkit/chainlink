package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = &zapLogger{}

// zapLoggerConfig defines the struct that serves as config when spinning up a the zap logger
type zapLoggerConfig struct {
	zap.Config
	local          Config
	diskLogLevel   zap.AtomicLevel
	diskStats      utils.DiskStatsProvider
	diskPollConfig zapDiskPollConfig

	// This is for tests only
	testDiskLogLvlChan chan zapcore.Level
}

type zapLogger struct {
	*zap.SugaredLogger
	config            zapLoggerConfig
	name              string
	fields            []interface{}
	callerSkip        int
	pollDiskSpaceStop chan struct{}
	pollDiskSpaceDone chan struct{}
}

func (cfg zapLoggerConfig) newLogger(cores ...zapcore.Core) (Logger, func() error, error) {
	cfg.diskLogLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	newCore, errWriter, err := cfg.newCore()
	if err != nil {
		return nil, nil, err
	}
	cores = append(cores, newCore)
	if cfg.local.DebugLogsToDisk() {
		diskCore, diskErr := cfg.newDiskCore()
		if diskErr != nil {
			return nil, nil, diskErr
		}
		cores = append(cores, diskCore)
	}

	core := zapcore.NewTee(cores...)
	lggr := &zapLogger{
		config:            cfg,
		pollDiskSpaceStop: make(chan struct{}),
		pollDiskSpaceDone: make(chan struct{}),
		SugaredLogger:     zap.New(core, zap.ErrorOutput(errWriter)).Sugar(),
	}

	if cfg.local.DebugLogsToDisk() {
		go lggr.pollDiskSpace()
	}

	var once sync.Once
	close := func() error {
		once.Do(func() {
			if cfg.local.DebugLogsToDisk() {
				close(lggr.pollDiskSpaceStop)
				<-lggr.pollDiskSpaceDone
			}
		})

		return lggr.Sync()
	}

	return lggr, close, err
}

func (cfg zapLoggerConfig) newCore() (zapcore.Core, zapcore.WriteSyncer, error) {
	encoder := zapcore.NewJSONEncoder(makeEncoderConfig(cfg.local))

	sink, closeOut, err := zap.Open(cfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}

	errSink, _, err := zap.Open(cfg.ErrorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}

	if cfg.Level == (zap.AtomicLevel{}) {
		return nil, nil, errors.New("missing Level")
	}

	filteredLogLevels := zap.LevelEnablerFunc(cfg.Level.Enabled)

	return zapcore.NewCore(encoder, sink, filteredLogLevels), errSink, nil
}

func makeEncoderConfig(cfg Config) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	if !cfg.UnixTS {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	encoderConfig.EncodeLevel = encodeLevel

	return encoderConfig
}

func (l *zapLogger) SetLogLevel(lvl zapcore.Level) {
	l.config.Level.SetLevel(lvl)
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
	newLogger.Trace("Named logger created")
	return &newLogger
}

func (l *zapLogger) NewRootLogger(lvl zapcore.Level) (Logger, error) {
	newLogger := *l
	newLogger.config.Level = zap.NewAtomicLevelAt(lvl)
	newCore, errWriter, err := newLogger.config.newCore()
	if err != nil {
		return nil, err
	}
	cores := []zapcore.Core{
		// The console core is what we want to be unique per root, so we spin a new one here
		newCore,
	}
	if newLogger.config.local.DebugLogsToDisk() {
		diskCore, diskErr := newLogger.config.newDiskCore()
		if diskErr != nil {
			return nil, diskErr
		}
		cores = append(cores, diskCore)
	}
	core := zap.New(zapcore.NewTee(cores...), zap.ErrorOutput(errWriter), zap.AddCallerSkip(l.callerSkip))

	newLogger.SugaredLogger = core.Named(l.name).Sugar().With(l.fields...)

	return &newLogger, err
}

func (l *zapLogger) Helper(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.sugaredHelper(skip)
	newLogger.callerSkip += skip
	return &newLogger
}

func (l *zapLogger) sugaredHelper(skip int) *zap.SugaredLogger {
	return l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
}

func (l *zapLogger) ErrorIf(err error, msg string) {
	if err != nil {
		l.Helper(1).Errorw(msg, "err", err)
	}
}

func (l *zapLogger) ErrorIfClosing(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		l.Helper(1).Errorw(fmt.Sprintf("Error closing %s", name), "err", err)
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

func (l *zapLogger) Recover(panicErr interface{}) {
	l.Criticalw("Recovered goroutine panic", "panic", panicErr)
}
