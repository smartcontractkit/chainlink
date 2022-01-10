package logger

import (
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
)

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

	// Helper creates a new logger with the number of callers skipped by caller annotation increased by skip.
	// This allows wrappers and helpers to point higher up the stack (like testing.T.Helper()).
	Helper(skip int) Logger

	// Recover reports recovered panics; this is useful because it avoids
	// double-reporting to sentry
	Recover(panicErr interface{})
}

// NewLogger returns a new Logger configured from environment variables, and logs any parsing errors.
// Tests should use TestLogger.
func NewLogger() Logger {
	var c Config
	var parseErrs []string

	var invalid string
	c.LogLevel, invalid = envvar.LogLevel.ParseLogLevel()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	c.Dir = os.Getenv("LOG_FILE_DIR")
	if c.Dir == "" {
		var invalid2 string
		c.Dir, invalid2 = envvar.RootDir.ParseString()
		if invalid2 != "" {
			parseErrs = append(parseErrs, invalid2)
		}
	}

	c.JsonConsole, invalid = envvar.JSONConsole.ParseBool()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	c.ToDisk, invalid = envvar.LogToDisk.ParseBool()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	c.UnixTS, invalid = envvar.LogUnixTS.ParseBool()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	l := c.New()
	for _, msg := range parseErrs {
		l.Error(msg)
	}
	return l
}

type Config struct {
	LogLevel                   zapcore.Level
	Dir                        string
	JsonConsole                bool
	UnixTS                     bool
	ToDisk                     bool // if false, the Logger will only log to stdout
	DiskMaxSizeBeforeRotate    int  // megabytes
	DiskMaxAgeBeforeDelete     int  // days
	DiskMaxBackupsBeforeDelete int  // files
}

// New returns a new Logger with pretty printing to stdout, prometeus counters, and sentry forwarding.
// Tests should use TestLogger.
func (c *Config) New() Logger {
	l, err := newZapLogger(c)
	if err != nil {
		log.Fatal(err)
	}
	l = newSentryLogger(l)
	return newPrometheusLogger(l)
}

// InitColor explicitly sets the global color.NoColor option.
// Not safe for concurrent use. Only to be called from init().
func InitColor(c bool) {
	color.NoColor = !c
}

// Constants for service names for package specific logging configuration
const (
	HeadTracker     = "HeadTracker"
	HeadListener    = "HeadListener"
	HeadSaver       = "HeadSaver"
	HeadBroadcaster = "HeadBroadcaster"
	FluxMonitor     = "FluxMonitor"
	Keeper          = "Keeper"
)

func GetLogServices() []string {
	return []string{
		HeadTracker,
		FluxMonitor,
		Keeper,
	}
}
