package testing

import (
	"fmt"
	"sync/atomic"
)

type Logger interface {
	Error(msg string, args ...any)
}

type T struct {
	lggr  Logger
	msgID atomic.Int32
}

func New(lggr Logger) *T {
	return &T{lggr: lggr}
}

func (t *T) Errorf(format string, args ...interface{}) {
	// if the log was produced by require/assert we need to split it, as engine does not allow logs longer than 1k bytes
	logLine := fmt.Sprintf(format, args...)
	const maxLength = 1000
	if len(logLine) <= maxLength {
		t.lggr.Error(fmt.Sprintf(format, args...))
	}
	id := t.msgID.Add(1)
	prefix := fmt.Sprintf("MsgID=%d ", id)
	prefix = prefix + " Part=%d/%d "
	prefixLength := len(prefix) // since parts counter in string are represented by %d we allocate 2 digits for each, which should be enough
	partLen := maxLength - prefixLength
	parts := len(logLine) / (partLen)
	if len(logLine)%partLen != 0 {
		parts++
	}
	i := 0
	start := 0
	for start < len(logLine) {
		i++
		partPrefix := fmt.Sprintf(prefix, id, i, parts)
		end := min(start+partLen, len(logLine))
		t.lggr.Error(fmt.Sprintf(partPrefix+format, logLine[start:end]))
		start = end
	}
}

func (t *T) FailNow() {
	panic("Test failed. Panic to stop execution")
}
