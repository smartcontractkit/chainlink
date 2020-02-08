package logger

// Based on https://stackoverflow.com/a/52737940

import (
	"bytes"
	"net/url"
	"sync"

	"go.uber.org/zap"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

var _ zap.Sink = &MemorySink{}

// Close is a dummy method to satisfy the zap.Sink interface
func (s *MemorySink) Close() error { return nil }

// Sync is a dummy method to satisfy the zap.Sink interface
func (s *MemorySink) Sync() error { return nil }

var testMemoryLog *MemorySink
var createSinkOnce sync.Once

func tmlCallback(*url.URL) (zap.Sink, error) { return testMemoryLog, nil }

func registerMemorySink() {
	testMemoryLog = &MemorySink{new(bytes.Buffer)}
	if err := zap.RegisterSink("memory", prettyConsoleSink(testMemoryLog)); err != nil {
		panic(err)
	}
}

func TestMemoryLog() *MemorySink {
	createSinkOnce.Do(registerMemorySink)
	return testMemoryLog
}
