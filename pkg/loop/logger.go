package loop

import (
	"io"
	"sync"

	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

// HCLogLogger returns an [hclog.Logger] backed by the given [logger.Logger].
func HCLogLogger(l logger.Logger) hclog.Logger {
	hcl := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Output: io.Discard, // only write through p.Logger Sink
	})
	hcl.RegisterSink(&hclSinkAdapter{l: l})
	return hcl
}

var _ hclog.SinkAdapter = (*hclSinkAdapter)(nil)

// hclSinkAdapter implements [hclog.SinkAdapter] with a [logger.Logger].
type hclSinkAdapter struct {
	l logger.Logger
	m sync.Map // [string]func() l.Logger
}

func (h *hclSinkAdapter) named(name string) logger.Logger {
	onceVal := onceValue(func() logger.Logger {
		return logger.Named(h.l, name)
	})
	v, _ := h.m.LoadOrStore(name, onceVal)
	return v.(func() logger.Logger)()
}

func removeArg(args []interface{}, key string) ([]interface{}, string) {
	if len(args) < 2 {
		return args, ""
	}
	for i := 0; i+1 < len(args); i += 2 {
		if args[i] == key {
			if v, ok := (args)[i+1].(string); ok {
				return slices.Delete(args, i, i+2), v
			}
			break
		}
	}
	return args, ""
}

func (h *hclSinkAdapter) Accept(_ string, level hclog.Level, msg string, args ...interface{}) {
	if level == hclog.Off {
		return
	}

	l := h.l
	var name string
	if args, name = removeArg(args, "logger"); name != "" {
		l = h.named(name)
	}
	switch level {
	case hclog.Off:
		return // unreachable, but satisfies linter
	case hclog.NoLevel:
	case hclog.Debug, hclog.Trace:
		l.Debugw(msg, args...)
	case hclog.Info:
		l.Infow(msg, args...)
	case hclog.Warn:
		l.Warnw(msg, args...)
	case hclog.Error:
		l.Errorw(msg, args...)
	}
}

// NewLogger returns a new [logger.Logger] configured to encode [hclog] compatible JSON.
func NewLogger() (logger.Logger, error) {
	return logger.NewWith(func(cfg *zap.Config) {
		cfg.Level.SetLevel(zap.DebugLevel)
		cfg.EncoderConfig.LevelKey = "@level"
		cfg.EncoderConfig.MessageKey = "@message"
		cfg.EncoderConfig.TimeKey = "@timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000Z07:00")
	})
}

// onceValue returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
//
// Note: Copied from sync.OnceValue in upcoming 1.21 release. Can be removed after upgrading.
func onceValue[T any](f func() T) func() T {
	var (
		once   sync.Once
		valid  bool
		p      any
		result T
	)
	g := func() {
		defer func() {
			p = recover()
			if !valid {
				panic(p)
			}
		}()
		result = f()
		valid = true
	}
	return func() T {
		once.Do(g)
		if !valid {
			panic(p)
		}
		return result
	}
}
