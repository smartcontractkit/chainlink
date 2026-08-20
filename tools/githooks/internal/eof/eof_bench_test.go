package eof_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/eof"
)

func BenchmarkFixContent(b *testing.B) {
	sample := bytes.Repeat([]byte("line of code here\n"), 100)
	sample = append(sample, []byte("\n\n\n")...)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = eof.FixContent(sample)
	}
}
