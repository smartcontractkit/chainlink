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

func BenchmarkFixContent_Go(b *testing.B) {
	code := []byte(`package main

const query = ` + "`" + `
SELECT *   
FROM table   
` + "`" + `

func main() {   
	x := 1   
	_ = x   
}
`)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _, _ = whitespace.FixContent("main.go", code)
	}
}

func BenchmarkFixContent_Python(b *testing.B) {
	code := []byte(`def main():   
    """
    multiline docstring with spaces    
    second line with spaces   
    """
    x = 1   
    return x
`)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _, _ = whitespace.FixContent("script.py", code)
	}
}

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
		name := fmt.Sprintf("file_%d.go", i)
		full := filepath.Join(tmpDir, name)
		_ = os.WriteFile(full, []byte("package main   \n\nfunc main() {}\n"), 0o600)
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
