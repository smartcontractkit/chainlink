package loop

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// HCLogLogger returns an [hclog.Logger] backed by the given [logger.Logger].
func HCLogLogger(l logger.Logger) hclog.Logger {
	hcl := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Output: io.Discard, // only write through p.Logger Sink
	})
	hcl.RegisterSink(&hclSinkAdapter{l: logger.Sugared(l)})
	return hcl
}

var _ hclog.SinkAdapter = (*hclSinkAdapter)(nil)

// hclSinkAdapter implements [hclog.SinkAdapter] with a [logger.Logger].
type hclSinkAdapter struct {
	l logger.SugaredLogger
	m sync.Map // [string]func() l.Logger
}

func (h *hclSinkAdapter) named(name string) logger.SugaredLogger {
	onceVal := onceValue(func() logger.SugaredLogger {
		return h.l.Named(name)
	})
	v, _ := h.m.LoadOrStore(name, onceVal)
	return v.(func() logger.SugaredLogger)()
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

// logMessage is the JSON payload that gets sent to Stderr from the plugin to the host
type logMessage struct {
	Message   string                 `json:"@message"`
	Level     string                 `json:"@level"`
	Timestamp time.Time              `json:"timestamp"`
	ExtraArgs []*LogMessageExtraArgs `json:"extra_args"`
}

// LogMessageExtraArgs is a key value pair within the Output payload
type LogMessageExtraArgs struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// flattenExtraArgs is used to flatten arguments of the log message
func flattenExtraArgs(le *logMessage) []interface{} {
	var result []interface{}
	result = append(result, "level")
	result = append(result, le.Level)
	result = append(result, "timestamp")
	result = append(result, le.Timestamp)
	for _, kv := range le.ExtraArgs {
		result = append(result, kv.Key)
		result = append(result, kv.Value)
	}

	return result
}

func parseJSON(input string) (*logMessage, error) {
	var raw map[string]interface{}
	entry := &logMessage{}

	err := json.Unmarshal([]byte(input), &raw)
	if err != nil {
		return nil, err
	}

	if v, ok := raw["@message"]; ok {
		entry.Message = v.(string)
		delete(raw, "@message")
	}

	if v, ok := raw["@level"]; ok {
		entry.Level = v.(string)
		delete(raw, "@level")
	}

	for k, v := range raw {
		entry.ExtraArgs = append(entry.ExtraArgs, &LogMessageExtraArgs{
			Key:   k,
			Value: v,
		})
	}

	return entry, nil
}

// logDebug will parse msg and figure out if it's a panic, fatal or critical log message, this is done here because the hashicorp plugin will push any
// unrecognizable message from stderr as a debug statement
func logDebug(msg string, l logger.SugaredLogger, args ...interface{}) {
	if strings.HasPrefix(msg, "panic:") {
		l.Criticalw(fmt.Sprintf("[PANIC] %s", msg), args...)
	} else if log, err := parseJSON(msg); err == nil {
		switch log.Level {
		case "fatal":
			l.Criticalw(fmt.Sprintf("[FATAL] %s", log.Message), flattenExtraArgs(log)...)
		case "critical", "dpanic":
			l.Criticalw(log.Message, flattenExtraArgs(log)...)
		default:
			l.Debugw(log.Message, flattenExtraArgs(log)...)
		}
	} else {
		l.Debugw(msg, args...)
	}
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
	case hclog.Trace:
		l.Debugw(msg, args...)
	case hclog.Debug:
		logDebug(msg, l, args...)
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
