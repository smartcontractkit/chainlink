package whitespace_test

import (
	"bytes"
	"testing"

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

	for i := 0; i < b.N; i++ {
		_, _, _ = whitespace.FixContent("main.go", code)
	}
}

func BenchmarkFixContent_Markdown(b *testing.B) {
	doc := bytes.Repeat([]byte("This line has two trailing spaces.  \nNormal line.\n"), 50)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, _ = whitespace.FixContent("README.md", doc)
	}
}
