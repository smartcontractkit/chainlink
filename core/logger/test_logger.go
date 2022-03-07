package logger

// Based on https://stackoverflow.com/a/52737940

import (
	"bytes"
	"log"
	"net/url"
	"sync"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/config/envvar"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	m sync.Mutex
	b bytes.Buffer
}

var _ zap.Sink = &MemorySink{}

func (s *MemorySink) Write(p []byte) (n int, err error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.b.Write(p)
}

// Close is a dummy method to satisfy the zap.Sink interface
func (s *MemorySink) Close() error { return nil }

// Sync is a dummy method to satisfy the zap.Sink interface
func (s *MemorySink) Sync() error { return nil }

// String returns the full log contents, as a string
func (s *MemorySink) String() string {
	s.m.Lock()
	defer s.m.Unlock()
	return s.b.String()
}

func (s *MemorySink) Reset() {
	s.m.Lock()
	defer s.m.Unlock()
	s.b.Reset()
}

var testMemoryLog MemorySink
var createSinkOnce sync.Once

func registerMemorySink() {
	if err := zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return PrettyConsole{Sink: &testMemoryLog}, nil
	}); err != nil {
		panic(err)
	}
}

func MemoryLogTestingOnly() *MemorySink {
	createSinkOnce.Do(registerMemorySink)
	return &testMemoryLog
}

// TestLogger creates a logger that directs output to PrettyConsole configured
// for test output, and to the buffer testMemoryLog. t is optional.
// Log level is derived from the LOG_LEVEL env var.
func TestLogger(t T) Logger {
	cfg := newZapConfigTest()
	ll, invalid := envvar.LogLevel.ParseLogLevel()
	cfg.Level.SetLevel(ll)
	l, close, err := zapLoggerConfig{Config: cfg}.newLogger()
	if err != nil {
		if t == nil {
			log.Fatal(err)
		}
		t.Fatal(err)
	}
	if invalid != "" {
		l.Error(invalid)
	}
	if t != nil {
		t.Cleanup(func() {
			assert.NoError(t, close())
		})
	}
	if t == nil {
		return l
	}
	return l.Named(verShaNameStatic()).Named(t.Name())
}

func newZapConfigTest() zap.Config {
	_ = MemoryLogTestingOnly() // Make sure memory log is created
	config := newZapConfigBase()
	config.OutputPaths = []string{"pretty://console", "memory://"}
	return config
}

type T interface {
	Name() string
	Cleanup(f func())
	Fatal(...interface{})
	Errorf(format string, args ...interface{})
}
