package logging

import (
	"fmt"
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

// This v2 version of `core.go` is a work in progress without any backwar compatibility
// version. It might not made it to an official version of the library so you can depend
// at your own risk.

type loggerOptions struct {
	autoStartSwitcherServer bool
	encoderVerbosity        int
	level                   zap.AtomicLevel
	loggerName              string
	reportAllErrors         bool
	serviceName             string
	zapOptions              []zap.Option
}

type LoggerOption interface {
	apply(o *loggerOptions)
}

type loggerFuncOption func(o *loggerOptions)

func (f loggerFuncOption) apply(o *loggerOptions) {
	f(o)
}

func WithAutoStartSwitcherServer() LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.autoStartSwitcherServer = true
	})
}

func WithAtomicLevel(level zap.AtomicLevel) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.level = level
	})
}

func WithLoggerName(name string) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.loggerName = name
	})
}

func WithReportAllErrors() LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.reportAllErrors = true
	})
}

func WithServiceName(name string) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.serviceName = name
	})
}

func WithZapOption(zapOption zap.Option) LoggerOption {
	return loggerFuncOption(func(o *loggerOptions) {
		o.zapOptions = append(o.zapOptions, zapOption)
	})
}

// NewNopLogger creates a new no-op logger (via `zap.NewNop`) and automatically registered it
// withing the logging registry.
func NewNopLogger(shortName string, packageID string) *zap.Logger {
	logger := zap.NewNop()

	Register(packageID, &logger)
	return logger
}

func NewSimpleLogger(shortName string, packageID string) *zap.Logger {
	opts := []LoggerOption{
		WithAtomicLevel(inferLevel()),
		WithReportAllErrors(),
		WithServiceName(shortName),
	}

	if isProductionEnvironment() {
		opts = append(opts, WithAutoStartSwitcherServer())
	}

	return NewLogger(shortName, packageID, opts...)
}

// NewLogger creates a new logger with sane defaults based on a varity of rules described
// below and automatically registered withing the logging registry.
func NewLogger(shortName string, packageID string, opts ...LoggerOption) *zap.Logger {
	logger, err := MaybeNewLogger(shortName, packageID, opts...)
	if err != nil {
		panic(fmt.Errorf("unable to create logger (in production? %t): %w", isProductionEnvironment(), err))
	}

	return logger
}

func MaybeNewLogger(shortName string, packageID string, opts ...LoggerOption) (*zap.Logger, error) {
	options := newDefaultLoggerOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	logger, err := newLogger(options)
	if err != nil {
		return nil, err
	}

	if options.loggerName != "" {
		logger = logger.Named(options.loggerName)
	}

	Register(packageID, &logger)
	return logger, nil
}

func newLogger(opts *loggerOptions) (*zap.Logger, error) {
	zapOptions := opts.zapOptions

	if isProductionEnvironment() {
		reportAllErrors := opts.reportAllErrors
		serviceName := opts.serviceName

		if reportAllErrors && opts.serviceName != "" {
			zapOptions = append(zapOptions, zapdriver.WrapCore(zapdriver.ReportAllErrors(true), zapdriver.ServiceName(serviceName)))
		} else if reportAllErrors {
			zapOptions = append(zapOptions, zapdriver.WrapCore(zapdriver.ReportAllErrors(true)))
		} else if opts.serviceName != "" {
			zapOptions = append(zapOptions, zapdriver.WrapCore(zapdriver.ServiceName(serviceName)))
		}

		return zapdriver.NewProductionConfig().Build(zapOptions...)
	}

	// Development logger
	isTTY := terminal.IsTerminal(int(os.Stderr.Fd()))
	logStdoutWriter := zapcore.Lock(os.Stderr)
	core := zapcore.NewCore(NewEncoder(opts.encoderVerbosity, isTTY), logStdoutWriter, opts.level)

	return zap.New(core), nil
}

func newDefaultLoggerOptions() (o *loggerOptions) {
	return &loggerOptions{
		encoderVerbosity: inferEncoderVerbosity(),
		level:            inferLevel(),
	}
}

func inferLevel() zap.AtomicLevel {
	if os.Getenv("DEBUG") != "" || os.Getenv("TRACE") != "" {
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	return zap.NewAtomicLevelAt(zapcore.InfoLevel)
}

func inferEncoderVerbosity() int {
	if os.Getenv("DEBUG") != "" || os.Getenv("TRACE") != "" {
		return 3
	}

	return 1
}

func isProductionEnvironment() bool {
	_, err := os.Stat("/.dockerenv")

	return !os.IsNotExist(err)
}
