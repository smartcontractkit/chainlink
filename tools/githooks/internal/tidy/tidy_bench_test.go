package tidy_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/tidy"
)

func BenchmarkTidyRun(b *testing.B) {
	repoRoot := "/test/repo"

	mockRunner := func(ctx context.Context, dir string, args ...string) error {
		// simulate minimal IO delay
		time.Sleep(100 * time.Microsecond)
		return nil
	}

	cfg := tidy.Config{Runner: mockRunner}
	modules3 := []string{".", "deployment", "tools/githooks"}
	modules10 := []string{
		".", "deployment", "tools/githooks", "tools/test",
		"core/scripts", "devenv", "system-tests/lib",
		"integration-tests", "integration-tests/load", "system-tests/tests/smoke",
	}

	b.Run("Parallel3Modules", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			err := tidy.Run(b.Context(), repoRoot, modules3, cfg)
			require.NoError(b, err)
		}
	})

	b.Run("Parallel10Modules", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			err := tidy.Run(b.Context(), repoRoot, modules10, cfg)
			require.NoError(b, err)
		}
	})
}
