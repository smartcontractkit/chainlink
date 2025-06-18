package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	common "github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// logsFile describes the logs file name
const logsFile = "chainlink_debug.log"

// Create a standard error writer to avoid test issues around os.Stderr being
// reassigned when verbose logging is enabled
type stderrWriter struct{}

func (sw stderrWriter) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}
func (sw stderrWriter) Close() error {
	return nil // never close stderr
}
func (sw stderrWriter) Sync() error {
	return os.Stderr.Sync()
}

func init() {
	err := zap.RegisterSink("pretty", prettyConsoleSink(stderrWriter{}))
	if err != nil {
		log.Fatalf("failed to register pretty printer %+v", err)
	}
	err = registerOSSinks()
	if err != nil {
		log.Fatalf("failed to register os specific sinks %+v", err)
	}
	if os.Getenv("CL_LOG_COLOR") != "true" {
		InitColor(false)
	}
}

var _ common.Logger = (Logger)(nil)

// Logger is the main interface of this package.
// It implements uber/zap's SugaredLogger interface and adds conditional logging helpers.
//
// Loggers should be injected (and usually Named as well): e.g. lggr.Named("<service name>")
//
// Tests
//   - Tests should use a TestLogger, with NewLogger being reserved for actual
//     runtime and limited direct testing.
//
// Levels
//   - Fatal: Logs and then calls os.Exit(1). Be careful about using this since it does NOT unwind the stack and may exit uncleanly.
//   - Panic: Unrecoverable error. Example: invariant violation, programmer error
//   - Critical: Requires quick action from the node op, obviously these should happen extremely rarely. Example: failed to listen on TCP port
//   - Error: Something bad happened, and it was clearly on the node op side. No need for immediate action though. Example: database write timed out
//   - Warn: Something bad happened, not clear who/what is at fault. Node ops should have a rough look at these once in a while to see whether anything stands out. Example: connection to peer was closed unexpectedly. observation timed out.
//   - Info: High level information. First level we’d expect node ops to look at. Example: entered new epoch with leader, made an observation with value, etc.
//   - Debug: Useful for forensic debugging, but we don't expect nops to look at this. Example: Got a message, dropped a message, ...
//   - Trace: Only included if compiled with the trace tag. For example: go test -tags trace ...
//
// Node Operator Docs: https://docs.chain.link/docs/configuration-variables/#log_level
// Deprecated: use [common.Logger] & [common.SugaredLogger]
type Logger interface {
	// With creates a new Logger with the given arguments
	With(args ...interface{}) Logger
	// Named creates a new Logger sub-scoped with name.
	// Names are inherited and dot-separated.
	//   a := l.Named("A") // logger=A
	//   b := a.Named("A") // logger=A.B
	// Names are generally `MixedCaps`, without spaces, like Go names.
	Named(name string) Logger

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

	// Sync flushes any buffered log entries.
	// Some insignificant errors are suppressed.
	Sync() error

	// Helper creates a new logger with the number of callers skipped by caller annotation increased by skip.
	// This allows wrappers and helpers to point higher up the stack (like testing.T.Helper()).
	Helper(skip int) Logger

	// Name returns the fully qualified name of the logger.
	Name() string

	// Recover reports recovered panics; this is useful because it avoids
	// double-reporting to sentry
	Recover(panicErr interface{})
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
	sha, ver := static.Short()
	return fmt.Sprintf("%s@%s", ver, sha)
}

// NewLogger returns a new Logger with default configuration.
// Tests should use TestLogger.
func NewLogger() (Logger, func() error) {
	var c Config
	return c.New()
}

type Config struct {
	LogLevel       zapcore.Level
	Dir            string
	JsonConsole    bool
	UnixTS         bool
	FileMaxSizeMB  int
	FileMaxAgeDays int
	FileMaxBackups int // files
	SentryEnabled  bool

	diskSpaceAvailableFn diskSpaceAvailableFn
	diskPollConfig       zapDiskPollConfig
	// This is for tests only
	testDiskLogLvlChan chan zapcore.Level
}

// New returns a new Logger with pretty printing to stdout, prometheus counters, and sentry forwarding.
// Tests should use TestLogger.
func (c *Config) New() (Logger, func() error) {
	if c.diskSpaceAvailableFn == nil {
		c.diskSpaceAvailableFn = diskSpaceAvailable
	}
	if !c.diskPollConfig.isSet() {
		c.diskPollConfig = newDiskPollConfig(diskPollInterval)
	}

	cfg := newZapConfigProd(c.JsonConsole, c.UnixTS)
	cfg.Level.SetLevel(c.LogLevel)
	var (
		l           Logger
		closeLogger func() error
		err         error
	)
	if !c.DebugLogsToDisk() {
		l, closeLogger, err = newDefaultLogger(cfg, c.UnixTS)
	} else {
		l, closeLogger, err = newRotatingFileLogger(cfg, *c)
	}
	if err != nil {
		log.Fatal(err)
	}

	if c.SentryEnabled {
		l = newSentryLogger(l)
	}
	l = newPrometheusLogger(l)
	l = l.With("version", verShaNameStatic())
	return l, closeLogger
}

// DebugLogsToDisk returns whether debug logs should be stored in disk
func (c Config) DebugLogsToDisk() bool {
	return c.FileMaxSizeMB > 0
}

// RequiredDiskSpace returns the required disk space in order to allow debug logs to be stored in disk
func (c Config) RequiredDiskSpace() utils.FileSize {
	return utils.FileSize(c.FileMaxSizeMB * utils.MB * (c.FileMaxBackups + 1))
}

func (c *Config) DiskSpaceAvailable(path string) (utils.FileSize, error) {
	if c.diskSpaceAvailableFn == nil {
		c.diskSpaceAvailableFn = diskSpaceAvailable
	}

	return c.diskSpaceAvailableFn(path)
}

func (c Config) LogsFile() string {
	return filepath.Join(c.Dir, logsFile)
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

func newDefaultLogger(zcfg zap.Config, unixTS bool) (Logger, func() error, error) {
	core, coreCloseFn, err := newDefaultLoggingCore(zcfg, unixTS)
	if err != nil {
		return nil, nil, err
	}

	l, loggerCloseFn, err := newLoggerForCore(zcfg, core)
	if err != nil {
		coreCloseFn()
		return nil, nil, err
	}

	return l, func() error {
		coreCloseFn()
		loggerCloseFn()
		return nil
	}, nil
}

func newLoggerForCore(zcfg zap.Config, core zapcore.Core) (*zapLogger, func(), error) {
	errSink, closeFn, err := zap.Open(zcfg.ErrorOutputPaths...)
	if err != nil {
		return nil, nil, err
	}

	return &zapLogger{
		level:         zcfg.Level,
		SugaredLogger: zap.New(core, zap.ErrorOutput(errSink), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar(),
	}, closeFn, nil
}

func newDefaultLoggingCore(zcfg zap.Config, unixTS bool) (zapcore.Core, func(), error) {
	encoder := zapcore.NewJSONEncoder(makeEncoderConfig(unixTS))

	sink, closeOut, err := zap.Open(zcfg.OutputPaths...)
	if err != nil {
		return nil, nil, err
	}

	if zcfg.Level == (zap.AtomicLevel{}) {
		return nil, nil, errors.New("missing Level")
	}

	filteredLogLevels := zap.LevelEnablerFunc(zcfg.Level.Enabled)

	core := zapcore.NewCore(encoder, sink, filteredLogLevels)
	return core, closeOut, nil
}

func newDiskCore(diskLogLevel zap.AtomicLevel, local Config) (zapcore.Core, error) {
	diskUsage, err := local.DiskSpaceAvailable(local.Dir)
	if err != nil || diskUsage < local.RequiredDiskSpace() {
		diskLogLevel.SetLevel(disabledLevel)
	}

	var (
		encoder = zapcore.NewConsoleEncoder(makeEncoderConfig(local.UnixTS))
		sink    = zapcore.AddSync(&lumberjack.Logger{
			Filename:   local.logFileURI(),
			MaxSize:    local.FileMaxSizeMB,
			MaxAge:     local.FileMaxAgeDays,
			MaxBackups: local.FileMaxBackups,
			Compress:   true,
		})
		allLogLevels = zap.LevelEnablerFunc(diskLogLevel.Enabled)
	)

	return zapcore.NewCore(encoder, sink, allLogLevels), nil
}
