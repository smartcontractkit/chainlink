package whitespace_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/whitespace"
)

func BenchmarkFixContent_Markdown(b *testing.B) {
	doc := bytes.Repeat([]byte("This line has two trailing spaces.  \nNormal line.\n"), 50)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _, _ = whitespace.FixContent("README.md", doc)
	}
}

func BenchmarkFixContent_Generic(b *testing.B) {
	doc := bytes.Repeat([]byte("key: value   \nlist:\n  - item 1   \n"), 50)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _, _ = whitespace.FixContent("config.yaml", doc)
	}
}

func BenchmarkRun_Parallel(b *testing.B) {
	tmpDir := b.TempDir()
	numFiles := 200
	fileList := make([]string, numFiles)

	for i := range numFiles {
		name := fmt.Sprintf("file_%d.yaml", i)
		full := filepath.Join(tmpDir, name)
		_ = os.WriteFile(full, []byte("key: value   \n\nkey2: value2\n"), 0o600)
		fileList[i] = name
	}

	ctx := context.Background()
	cfg := whitespace.Config{CheckOnly: true}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, err := whitespace.Run(ctx, tmpDir, fileList, cfg)
		require.NoError(b, err)
	}
}
