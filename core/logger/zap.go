package logger

import (
	"os"
	"sync"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AtomicCore provides thread-safe core swapping using atomic operations.
// It starts as a noop core and can be atomically swapped to include additional cores.
var _ zapcore.Core = &CoreWrapper{}

type AtomicCore struct {
	root     *CoreWrapper
	register *Register[CoreWrapper]
}

func NewAtomicCore() *AtomicCore {
	register := &Register[CoreWrapper]{
		stopCleanup: make(chan struct{}),
	}
	register.startPeriodicCleanup()
	root := &CoreWrapper{register: register, core: zapcore.NewNopCore()}
	return &AtomicCore{root, register}
}

func (a *AtomicCore) Root() zapcore.Core {
	return a.root
}

func (a *AtomicCore) Store(core zapcore.Core) {
	a.register.Update(func(cw *CoreWrapper) {
		cw.Store(core)
	})
}

func (a *AtomicCore) Close() {
	a.register.Close()
}

type CoreWrapper struct {
	mu       sync.RWMutex
	core     zapcore.Core
	fields   []zapcore.Field
	register *Register[CoreWrapper]
}

func (c *CoreWrapper) Store(core zapcore.Core) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.core = core.With(c.fields)
}

func (c *CoreWrapper) With(fields []zapcore.Field) zapcore.Core {
	c.mu.RLock()
	defer c.mu.RUnlock()
	combined := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	combined = append(combined, c.fields...)
	combined = append(combined, fields...)

	cw := &CoreWrapper{register: c.register, fields: combined, core: c.core.With(fields)}
	c.register.Add(cw)
	return cw
}

func (c *CoreWrapper) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.core.Check(ent, ce)
}

func (c *CoreWrapper) Enabled(lvl zapcore.Level) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.core.Enabled(lvl)
}

func (c *CoreWrapper) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.core.Write(ent, fields)
}

func (c *CoreWrapper) Sync() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.core.Sync()
}

var _ Logger = &zapLogger{}

type zapLogger struct {
	*zap.SugaredLogger
	level      zap.AtomicLevel
	fields     []any
	callerSkip int
}

func makeEncoderConfig(unixTS bool) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	if !unixTS {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	encoderConfig.EncodeLevel = encodeLevel

	return encoderConfig
}

func (l *zapLogger) SetLogLevel(lvl zapcore.Level) {
	l.level.SetLevel(lvl)
}

func (l *zapLogger) With(args ...any) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.With(args...)
	newLogger.fields = copyFields(l.fields, args...)
	return &newLogger
}

// copyFields returns a copy of fields with add appended.
func copyFields(fields []any, add ...any) []any {
	f := make([]any, 0, len(fields)+len(add))
	f = append(f, fields...)
	f = append(f, add...)
	return f
}

func (l *zapLogger) Named(name string) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.SugaredLogger.Named(name)
	newLogger.Trace("Named logger created")
	return &newLogger
}

func (l *zapLogger) Helper(skip int) Logger {
	newLogger := *l
	newLogger.SugaredLogger = l.sugaredHelper(skip)
	newLogger.callerSkip += skip
	return &newLogger
}

func (l *zapLogger) Name() string {
	return l.Desugar().Name()
}

func (l *zapLogger) sugaredHelper(skip int) *zap.SugaredLogger {
	return l.WithOptions(zap.AddCallerSkip(skip))
}

func (l *zapLogger) Sync() error {
	err := l.SugaredLogger.Sync()
	if err == nil {
		return nil
	}
	var msg string
	if uw := pkgerrors.Unwrap(err); uw != nil {
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

func (l *zapLogger) Recover(panicErr any) {
	l.Criticalw("Recovered goroutine panic", "panic", panicErr)
}
