package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const benchSlowThreshold = 30 * time.Second

// realScenarioIterationBytes loads iteration-{0,1,2}.log.jsonl from testdata.
func realScenarioIterationBytes(b *testing.B) [][]byte {
	b.Helper()
	out := make([][]byte, 3)
	for i := range 3 {
		p := filepath.Join("testdata", fmt.Sprintf("iteration-%d.log.jsonl", i))
		data, err := os.ReadFile(p)
		if err != nil {
			b.Fatalf("read %s: %v", p, err)
		}
		out[i] = data
	}
	return out
}

func BenchmarkAnalyze_RealThreeIterations(b *testing.B) {
	payloads := realScenarioIterationBytes(b)
	var total int64
	for _, p := range payloads {
		total += int64(len(p))
	}
	b.ReportAllocs()
	b.SetBytes(total)

	for b.Loop() {
		rs := make([]io.Reader, len(payloads))
		for i, p := range payloads {
			rs[i] = bytes.NewReader(p)
		}
		_, _, err := Analyze(rs, benchSlowThreshold)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDigestIterationJSONL(b *testing.B) {
	payloads := realScenarioIterationBytes(b)
	for i, payload := range payloads {
		b.Run(fmt.Sprintf("iter_%d", i), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(payload)))
			for b.Loop() {
				_, err := DigestIterationJSONL(bytes.NewReader(payload), benchSlowThreshold)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkAnalyzeResults_RealDir measures glob, sort, open per call, plus Analyze.
// Warm-up runs before the timed b.Loop body so first-open cold start does not skew results; each op still opens files.
func BenchmarkAnalyzeResults_RealDir(b *testing.B) {
	dir := "testdata"
	if _, err := os.Stat(dir); err != nil {
		b.Fatalf("testdata dir: %v", err)
	}
	payloads := realScenarioIterationBytes(b)
	var total int64
	for _, p := range payloads {
		total += int64(len(p))
	}
	if _, _, err := AnalyzeResults(dir, benchSlowThreshold); err != nil {
		b.Fatalf("warm-up AnalyzeResults: %v", err)
	}
	b.ReportAllocs()
	b.SetBytes(total)
	for b.Loop() {
		_, _, err := AnalyzeResults(dir, benchSlowThreshold)
		if err != nil {
			b.Fatal(err)
		}
	}
}
