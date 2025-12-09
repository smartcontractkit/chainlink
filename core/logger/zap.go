package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
	"weak"

	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AtomicCore provides thread-safe core swapping using atomic operations.
// It starts as a noop core and can be atomically swapped to include additional cores.
var _ zapcore.Core = &AtomicCore{}

const cleanupInterval = time.Minute * 5

type AtomicCore struct {
	mu          sync.RWMutex
	core        zapcore.Core
	children    map[weak.Pointer[withCore]]struct{}
	stopCleanup chan struct{}
	cleanupWg   sync.WaitGroup
}

// NewAtomicCore creates a new AtomicCore initialized with a noop core
func NewAtomicCore() *AtomicCore {
	ac := &AtomicCore{
		core:        zapcore.NewNopCore(),
		children:    make(map[weak.Pointer[withCore]]struct{}),
		stopCleanup: make(chan struct{}),
	}
	ac.startPeriodicCleanup()
	fmt.Println("DEBUG AtomicCore cleanup", "step", "new")
	return ac
}

func (d *AtomicCore) Store(core zapcore.Core) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.core = core
	for p := range d.children {
		c := p.Value()
		if c != nil {
			c.Store(d.core)
		}
	}
}

func (d *AtomicCore) load() zapcore.Core {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.core
}

func (d *AtomicCore) Enabled(l zapcore.Level) bool { return d.load().Enabled(l) }

func (d *AtomicCore) With(fs []zapcore.Field) zapcore.Core {
	d.mu.Lock()
	defer d.mu.Unlock()
	coreWithFields := d.core.With(fs)
	w := &withCore{fields: fs, AtomicCore: AtomicCore{
		core:     coreWithFields,
		children: make(map[weak.Pointer[withCore]]struct{}),
	}}
	if d.children == nil {
		d.children = make(map[weak.Pointer[withCore]]struct{})
	}
	d.children[weak.Make(w)] = struct{}{}
	return w
}

func (d *AtomicCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return d.load().Check(e, ce)
}

func (d *AtomicCore) Write(e zapcore.Entry, fs []zapcore.Field) error { return d.load().Write(e, fs) }

func (d *AtomicCore) Sync() error { return d.load().Sync() }

func (d *AtomicCore) Close() {
	close(d.stopCleanup)
	d.cleanupWg.Wait()
}

func (d *AtomicCore) cleanup() {
	fmt.Println("DEBUG AtomicCore cleanup", "step", "enter")
	var wg sync.WaitGroup
	defer wg.Wait()
	d.mu.Lock()
	defer d.mu.Unlock()
	fmt.Println("DEBUG AtomicCore cleanup", "step", "init loop")
	initialChildrenCount := len(d.children)
	for p := range d.children {
		c := p.Value()
		if c == nil {
			delete(d.children, p)
			continue
		}
		go c.cleanup()
	}
	fmt.Println("DEBUG AtomicCore cleanup", "deleted_children_count", initialChildrenCount-len(d.children))
}

func (d *AtomicCore) startPeriodicCleanup() {
	d.cleanupWg.Go(func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("DEBUG AtomicCore cleanup", "step", "trigger")
				d.cleanup()
			case <-d.stopCleanup:
				return
			}
		}
	})
}

type withCore struct {
	fields []zapcore.Field
	AtomicCore
}

func (w *withCore) Store(core zapcore.Core) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.core = core.With(w.fields)
	for p := range w.children {
		c := p.Value()
		if c != nil {
			c.Store(w.core)
		}
	}
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
