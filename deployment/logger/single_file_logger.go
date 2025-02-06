package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	// The Chainlink logger interface we're implementing:
	corelogger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

// SingleFileLogger is a wrapper around a zap SugaredLogger that implements
// Chainlink’s corelogger.Logger interface.
type SingleFileLogger struct {
	*zap.SugaredLogger
}

// Ensure it truly implements the interface at compile time:
var _ corelogger.Logger = (*SingleFileLogger)(nil)

// NewSingleFileLogger creates a zap-based logger that writes everything to one file.
// The file name includes the test name + timestamp so that parallel tests don’t collide.
func NewSingleFileLogger(tb testing.TB) *SingleFileLogger {
	// Our logs will go here so GH can upload them:
	baseDir := "logs"

	// For uniqueness, include test name + timestamp
	filename := fmt.Sprintf("%s_%d.log", tb.Name(), time.Now().UnixNano())
	dirOfFilename := filepath.Dir(filename)

	dir, err := filepath.Abs(filepath.Join(baseDir, dirOfFilename))
	if err != nil {
		log.Fatalf("Failed to get absolute path for %q: %v", filepath.Join(baseDir, dirOfFilename), err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("Failed to create logs dir %q: %v", dir, err)
	}

	fullPath, err := filepath.Abs(filepath.Join(baseDir, filename))
	if err != nil {
		log.Fatalf("Failed to get absolute path for %q: %v", fullPath, err)
	}
	// Log the full path
	log.Printf("Logging to %q", fullPath)
	f, err := os.Create(fullPath)
	if err != nil {
		log.Fatalf("Failed to create %q: %v", fullPath, err)
	}

	// Auto-flush writer so we don't lose logs if the test fails abruptly
	writer := &autoFlushWriter{file: f}

	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000000000")

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(writer),
		zap.DebugLevel,
	)

	zapLogger := zap.New(core, zap.AddCaller())
	return &SingleFileLogger{zapLogger.Sugar()}
}

// autoFlushWriter flushes on every log write, so we don’t lose logs if the test crashes.
type autoFlushWriter struct {
	file *os.File
}

func (w *autoFlushWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.Write(p)
	if err == nil {
		_ = w.file.Sync()
	}
	return n, err
}

func (w *autoFlushWriter) Sync() error {
	return w.file.Sync()
}

// ----------------------------------------------------------------------------
// Implement the corelogger.Logger interface
// ----------------------------------------------------------------------------

func (l *SingleFileLogger) Name() string {
	return l.Desugar().Name()
}

func (l *SingleFileLogger) Named(name string) corelogger.Logger {
	// zap’s Named(...) returns a new *zap.Logger. We return a new SingleFileLogger wrapping it.
	return &SingleFileLogger{l.SugaredLogger.Named(name)}
}

func (l *SingleFileLogger) Trace(args ...interface{}) {
	// “Trace” in Chainlink’s interface is basically a debug-level call
	l.Debug(args...)
}

func (l *SingleFileLogger) Tracef(format string, values ...interface{}) {
	l.Debugf(format, values...)
}

func (l *SingleFileLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.Debugw(msg, keysAndValues...)
}

func (l *SingleFileLogger) Critical(args ...interface{}) {
	// “Critical” is typically mapped to an error-level call or fatal if you prefer
	l.Error(args...)
}

func (l *SingleFileLogger) Criticalf(format string, values ...interface{}) {
	l.Errorf(format, values...)
}

func (l *SingleFileLogger) Criticalw(msg string, keysAndValues ...interface{}) {
	l.Errorw(msg, keysAndValues...)
}

func (l *SingleFileLogger) Helper(skip int) corelogger.Logger {
	// Helper() is for skipping extra stack frames in error logs. We can do AddCallerSkip here:
	return &SingleFileLogger{l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()}
}

func (l *SingleFileLogger) Recover(panicErr interface{}) {
	// Called on panics, typically you might add custom logic or stacktrace
	if panicErr != nil {
		l.Errorf("Recovering from panic: %v", panicErr)
		debug.PrintStack()
	}
}

func (l *SingleFileLogger) SetLogLevel(zapLevel zapcore.Level) {
	// If you need dynamic level changes, you can store an AtomicLevel.
	// If not, you can just do nothing or rebuild a new logger.
}

func (l *SingleFileLogger) With(args ...interface{}) corelogger.Logger {
	// Adds extra fields to the logger. Return a new instance with them.
	return &SingleFileLogger{l.SugaredLogger.With(args...)}
}
