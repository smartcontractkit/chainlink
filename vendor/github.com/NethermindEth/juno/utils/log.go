package utils

import (
	"encoding"
	"errors"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrUnknownLogLevel = errors.New("unknown log level (known: debug, info, warn, error)")

type LogLevel int

// The following are necessary for Cobra and Viper, respectively, to unmarshal log level
// CLI/config parameters properly.
var (
	_ pflag.Value              = (*LogLevel)(nil)
	_ encoding.TextUnmarshaler = (*LogLevel)(nil)
)

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	default:
		// Should not happen.
		panic(ErrUnknownLogLevel)
	}
}

func (l *LogLevel) Set(s string) error {
	switch s {
	case "DEBUG", "debug":
		*l = DEBUG
	case "INFO", "info":
		*l = INFO
	case "WARN", "warn":
		*l = WARN
	case "ERROR", "error":
		*l = ERROR
	default:
		return ErrUnknownLogLevel
	}
	return nil
}

func (l *LogLevel) Type() string {
	return "LogLevel"
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	return l.Set(string(text))
}

type Logger interface {
	SimpleLogger
	pebble.Logger
}

type SimpleLogger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
}

type ZapLogger struct {
	*zap.SugaredLogger
}

var _ Logger = (*ZapLogger)(nil)

func NewNopZapLogger() *ZapLogger {
	return &ZapLogger{zap.NewNop().Sugar()}
}

func NewZapLogger(logLevel LogLevel, colour bool) (*ZapLogger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if !colour {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("15:04:05.000 02/01/2006 -07:00"))
	}
	level, err := zapcore.ParseLevel(logLevel.String())
	if err != nil {
		return nil, err
	}
	config.Level.SetLevel(level)
	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{log.Sugar()}, nil
}

func (l *ZapLogger) Warningf(msg string, args ...any) {
	l.Warnf(msg, args)
}
