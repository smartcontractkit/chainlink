package logger

// Based on https://stackoverflow.com/a/52737940

import (
	"bytes"
	"log"
	"net/url"
	"sync"

	"go.uber.org/zap"
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
	l, err := newZapLogger(newTestConfig())
	if err != nil {
		if t == nil {
			log.Fatal(err)
		}
		t.Fatal(err)
	}
	if t == nil {
		return l
	}
	return l.Named(t.Name())
}

func newTestConfig() zap.Config {
	_ = MemoryLogTestingOnly() // Make sure memory log is created
	config := newBaseConfig()
	config.Level.SetLevel(envLvl)
	config.OutputPaths = []string{"pretty://console", "memory://"}
	return config
}

type T interface {
	Name() string
	Fatal(...interface{})
}
