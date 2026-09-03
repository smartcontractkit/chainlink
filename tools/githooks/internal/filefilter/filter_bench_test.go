package filefilter_test

import (
	"bytes"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/filefilter"
)

func BenchmarkIsEligiblePath(b *testing.B) {
	paths := []string{
		"core/services/app.go",
		"README.md",
		"tools/githooks/.bin/githooks",
		"assets/logo.png",
		"deployment/environment/env.go",
		"contracts/Token.sol",
		"Dockerfile",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		for _, p := range paths {
			_ = filefilter.IsEligiblePath(p)
		}
	}
}

func BenchmarkIsBinary(b *testing.B) {
	data := bytes.Repeat([]byte("func main() { fmt.Println(\"hello world\") }\n"), 50)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = filefilter.IsBinary(data)
	}
}
