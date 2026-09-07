package eof_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/eof"
)

func BenchmarkFixContent(b *testing.B) {
	sample := bytes.Repeat([]byte("line of code here\n"), 100)
	sample = append(sample, []byte("\n\n\n")...)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _ = eof.FixContent(sample)
	}
}

func BenchmarkRun_Parallel(b *testing.B) {
	tmpDir := b.TempDir()
	numFiles := 200
	fileList := make([]string, numFiles)

	for i := range numFiles {
		name := fmt.Sprintf("file_%d.go", i)
		full := filepath.Join(tmpDir, name)
		_ = os.WriteFile(full, []byte("package main\n\nfunc main() {}"), 0o600)
		fileList[i] = name
	}

	ctx := context.Background()
	cfg := eof.Config{CheckOnly: true}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, err := eof.Run(ctx, tmpDir, fileList, cfg)
		require.NoError(b, err)
	}
}
