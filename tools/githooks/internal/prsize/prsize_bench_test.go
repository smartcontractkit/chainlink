package prsize_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/prsize"
)

func BenchmarkCalculateEffectiveLines(b *testing.B) {
	files := make([]prsize.FileStat, 100)
	for i := range files {
		files[i] = prsize.FileStat{
			Path:      "core/services/feature/file.go",
			Additions: i * 2,
			Deletions: i,
		}
	}

	strategies := []prsize.Strategy{
		prsize.StrategyPerFileMax,
		prsize.StrategySum,
		prsize.StrategyMax,
		prsize.StrategyWeighted,
	}

	for _, strategy := range strategies {
		b.Run(string(strategy), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = prsize.CalculateEffectiveLines(files, strategy)
			}
		})
	}
}

func BenchmarkFormatReport(b *testing.B) {
	stat := &prsize.DiffStat{
		FilesChanged:   25,
		Additions:      650,
		Deletions:      120,
		EffectiveLines: 700,
		MergeBase:      "0123456789abcdef",
		IgnoredFiles:   []string{"go.sum", "package-lock.json"},
	}
	cfg := prsize.Config{
		SmallLimit:  200,
		MediumLimit: 500,
		Strategy:    prsize.StrategyPerFileMax,
	}

	b.ReportAllocs()
	for b.Loop() {
		_ = prsize.FormatReport(stat, prsize.SizeLarge, cfg)
	}
}
