package logger

import (
	"bytes"
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

func MemoryLogTestingOnly() *MemorySink {
	return &testMemoryLog
}
