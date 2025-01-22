package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

type SplitTestLogger struct {
	*zap.SugaredLogger
}

func (l SplitTestLogger) Name() string {
	return l.Desugar().Name()
}

func (l SplitTestLogger) Named(name string) corelogger.Logger {
	return SplitTestLogger{l.SugaredLogger.Named(name)}
}

func (l SplitTestLogger) Trace(args ...interface{}) {
	l.Debug(args...)
}

func (l SplitTestLogger) Tracef(format string, values ...interface{}) {
	l.Debugf(format, values...)
}

func (l SplitTestLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.Debugw(msg, keysAndValues...)
}

func (l SplitTestLogger) Critical(args ...interface{}) {
	l.Error(args...)
}

func (l SplitTestLogger) Criticalf(format string, values ...interface{}) {
	l.Errorf(format, values...)
}

func (l SplitTestLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	l.Errorw(msg, keysAndValues...)
}

func (l SplitTestLogger) Helper(skip int) corelogger.Logger {
	return SplitTestLogger{l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()}
}

func (l SplitTestLogger) Recover(panicErr interface{}) {

}

func (l SplitTestLogger) SetLogLevel(zapcore.Level) {}

func (l SplitTestLogger) With(args ...interface{}) corelogger.Logger {
	return SplitTestLogger{l.SugaredLogger.With(args...)}
}

func NewNamedTestLogger(tb testing.TB) *SplitTestLogger {
	nameFn := func(level string) string {
		return fmt.Sprintf("%s_%s.log", tb.Name(), level)
	}

	logger := zap.New(zapcore.NewTee(getCoresWithNamedFunction(nameFn)...))

	return &SplitTestLogger{logger.Sugar()}
}

func getCoresWithNamedFunction(nameFn func(level string) string) []zapcore.Core {
	sharedEncoderConfig := zap.NewDevelopmentEncoderConfig()
	sharedEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000000000")

	dir := "logs"

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	debugFile, err := os.Create(filepath.Join(dir, nameFn("debug")))
	if err != nil {
		log.Fatalf("Failed to create debug.log: %v", err)
	}
	infoFile, err := os.Create(filepath.Join(dir, nameFn("info")))
	if err != nil {
		log.Fatalf("Failed to create info.log: %v", err)
	}
	warnFile, err := os.Create(filepath.Join(dir, nameFn("warn")))
	if err != nil {
		log.Fatalf("Failed to create warn.log: %v", err)
	}
	errorFile, err := os.Create(filepath.Join(dir, nameFn("error")))
	if err != nil {
		log.Fatalf("Failed to create error.log: %v", err)
	}

	debugLevelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.DebugLevel
	})
	infoLevelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.InfoLevel
	})
	warnLevelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.WarnLevel
	})
	errorLevelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.ErrorLevel || level == zap.DPanicLevel || level == zap.PanicLevel || level == zap.FatalLevel
	})

	debugFileWriter := &autoFlushWriter{writer: debugFile}
	infoFileWriter := &autoFlushWriter{writer: infoFile}
	warnFileWriter := &autoFlushWriter{writer: warnFile}
	errorFileWriter := &autoFlushWriter{writer: errorFile}

	consoleEncoder := zapcore.NewConsoleEncoder(sharedEncoderConfig)
	consoleWriter := zapcore.Lock(os.Stdout)

	// send only Error
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		consoleWriter,
		zap.ErrorLevel,
	)

	// add trace?

	debugFileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(sharedEncoderConfig),
		zapcore.AddSync(debugFileWriter),
		debugLevelEnabler,
	)
	infoFileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(sharedEncoderConfig),
		zapcore.AddSync(infoFileWriter),
		infoLevelEnabler,
	)
	warnFileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(sharedEncoderConfig),
		zapcore.AddSync(warnFileWriter),
		warnLevelEnabler,
	)
	errorFileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(sharedEncoderConfig),
		zapcore.AddSync(errorFileWriter),
		errorLevelEnabler,
	)

	return []zapcore.Core{debugFileCore, infoFileCore, warnFileCore, errorFileCore, consoleCore}
}

type autoFlushWriter struct {
	writer *os.File
}

func (afw *autoFlushWriter) Write(p []byte) (n int, err error) {
	n, err = afw.writer.Write(p)
	if err == nil && afw.canSync() {
		_ = afw.writer.Sync() // Flush after every write
	}
	return
}

func (afw *autoFlushWriter) Sync() error {
	if afw.canSync() {
		return afw.writer.Sync()
	}
	return nil
}

func (afw *autoFlushWriter) canSync() bool {
	// Check if the writer supports syncing
	stat, err := afw.writer.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeNamedPipe) == 0 // Pipes cannot be synced
}
