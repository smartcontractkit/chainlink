package logger

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// LogsFile describes the logs file name
const LogsFile = "chainlink_debug.log"

func init() {
	err := zap.RegisterSink("pretty", prettyConsoleSink(os.Stderr))
	if err != nil {
		log.Fatalf("failed to register pretty printer %+v", err)
	}
	err = registerOSSinks()
	if err != nil {
		log.Fatalf("failed to register os specific sinks %+v", err)
	}
	if os.Getenv("LOG_COLOR") != "true" {
		InitColor(false)
	}
}

//go:generate mockery --name Logger --output . --filename logger_mock_test.go --inpackage --case=underscore

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
//
// Loggers should be injected (and usually Named as well): e.g. lggr.Named("<service name>")
//
// Tests
//  - Tests should use a TestLogger, with NewLogger being reserved for actual
//    runtime and limited direct testing.
//
// Levels
//  - Fatal: Logs and then calls os.Exit(1). Be careful about using this since it does NOT unwind the stack and may exit uncleanly.
//  - Panic: Unrecoverable error. Example: invariant violation, programmer error
//  - Critical: Requires quick action from the node op, obviously these should happen extremely rarely. Example: failed to listen on TCP port
//  - Error: Something bad happened, and it was clearly on the node op side. No need for immediate action though. Example: database write timed out
//  - Warn: Something bad happened, not clear who/what is at fault. Node ops should have a rough look at these once in a while to see whether anything stands out. Example: connection to peer was closed unexpectedly. observation timed out.
//  - Info: High level information. First level we’d expect node ops to look at. Example: entered new epoch with leader, made an observation with value, etc.
//  - Debug: Useful for forensic debugging, but we don't expect nops to look at this. Example: Got a message, dropped a message, ...
//  - Trace: Only included if compiled with the trace tag. For example: go test -tags trace ...
//
// Node Operator Docs: https://docs.chain.link/docs/configuration-variables/#log_level
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
	// Fatal logs and then calls os.Exit(1)
	// Be careful about using this since it does NOT unwind the stack and may
	// exit uncleanly
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
	Criticalw(msg string, keysAndValues ...interface{})
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

// newZapConfigProd returns a new production zap.Config.
func newZapConfigProd(jsonConsole bool, unixTS bool) zap.Config {
	config := newZapConfigBase()
	if !unixTS {
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	if !jsonConsole {
		config.OutputPaths = []string{"pretty://console"}
	}
	return config
}

func verShaNameStatic() string {
	return verShaName(static.Version, static.Sha)
}

func verShaName(ver, sha string) string {
	if sha == "" {
		sha = "unset"
	} else if len(sha) > 7 {
		sha = sha[:7]
	}
	if ver == "" {
		ver = "unset"
	}
	return fmt.Sprintf("%s@%s", ver, sha)
}

// NewLogger returns a new Logger configured from environment variables, and logs any parsing errors.
// Tests should use TestLogger.
func NewLogger() (Logger, func() error) {
	var c Config
	var parseErrs []string
	var warnings []string

	var invalid string
	c.LogLevel, invalid = envvar.LogLevel.Parse()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	c.Dir = os.Getenv("LOG_FILE_DIR")
	if c.Dir == "" {
		var invalid2 string
		c.Dir, invalid2 = envvar.RootDir.Parse()
		if invalid2 != "" {
			parseErrs = append(parseErrs, invalid2)
		}
	}

	c.JsonConsole, invalid = envvar.JSONConsole.Parse()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	var fileMaxSize utils.FileSize
	fileMaxSize, invalid = envvar.LogFileMaxSize.Parse()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}
	if fileMaxSize <= 0 {
		c.FileMaxSizeMB = 0 // disabled
	} else if fileMaxSize < utils.MB {
		c.FileMaxSizeMB = 1 // 1Mb is the minimum accepted by logging backend
		warnings = append(warnings, fmt.Sprintf("LogFileMaxSize %s is too small: using default %s", fileMaxSize, utils.FileSize(utils.MB)))
	} else {
		c.FileMaxSizeMB = int(fileMaxSize / utils.MB)
	}

	if c.DebugLogsToDisk() {
		var (
			fileMaxAge int64
			maxBackups int64
		)

		fileMaxAge, invalid = envvar.LogFileMaxAge.Parse()
		c.FileMaxAgeDays = int(fileMaxAge)
		if invalid != "" {
			parseErrs = append(parseErrs, invalid)
		}

		maxBackups, invalid = envvar.LogFileMaxBackups.Parse()
		c.FileMaxBackups = int(maxBackups)
		if invalid != "" {
			parseErrs = append(parseErrs, invalid)
		}
	}

	c.UnixTS, invalid = envvar.LogUnixTS.Parse()
	if invalid != "" {
		parseErrs = append(parseErrs, invalid)
	}

	l, close := c.New()
	for _, msg := range parseErrs {
		l.Error(msg)
	}
	for _, msg := range warnings {
		l.Warn(msg)
	}
	return l.Named(verShaNameStatic()), close
}

type Config struct {
	LogLevel       zapcore.Level
	Dir            string
	JsonConsole    bool
	UnixTS         bool
	FileMaxSizeMB  int
	FileMaxAgeDays int
	FileMaxBackups int // files
}

// New returns a new Logger with pretty printing to stdout, prometheus counters, and sentry forwarding.
// Tests should use TestLogger.
func (c *Config) New() (Logger, func() error) {
	cfg := newZapConfigProd(c.JsonConsole, c.UnixTS)
	cfg.Level.SetLevel(c.LogLevel)
	l, close, err := zapLoggerConfig{
		local:          *c,
		Config:         cfg,
		diskStats:      utils.NewDiskStatsProvider(),
		diskPollConfig: newDiskPollConfig(diskPollInterval),
	}.newLogger()
	if err != nil {
		log.Fatal(err)
	}
	l = newSentryLogger(l)
	return newPrometheusLogger(l), close
}

// DebugLogsToDisk returns whether debug logs should be stored in disk
func (c Config) DebugLogsToDisk() bool {
	return c.FileMaxSizeMB > 0
}

// RequiredDiskSpace returns the required disk space in order to allow debug logs to be stored in disk
func (c Config) RequiredDiskSpace() utils.FileSize {
	return utils.FileSize(c.FileMaxSizeMB * utils.MB * (c.FileMaxBackups + 1))
}

// InitColor explicitly sets the global color.NoColor option.
// Not safe for concurrent use. Only to be called from init().
func InitColor(c bool) {
	color.NoColor = !c
}

// newZapConfigBase returns a zap.NewProductionConfig with sampling disabled and a modified level encoder.
func newZapConfigBase() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	cfg.EncoderConfig.EncodeLevel = encodeLevel
	return cfg
}
