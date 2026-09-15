package generate_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/generate"
)

type benchRunner struct{}

func (b *benchRunner) Run(ctx context.Context, dir string, args ...string) error {
	return nil
}

func BenchmarkGenerateRun(b *testing.B) {
	repoRoot := "/test/repo"
	cfg := generate.Config{Runner: (&benchRunner{}).Run}

	protoFiles := []string{
		"core/capabilities/remote/types/messages.proto",
		"core/services/llo/telem/telem_streams.proto",
		"core/services/nodestatusreporter/bridgestatus/events/bridge_status.proto",
	}

	modFiles := []string{
		"go.mod",
		"deployment/go.sum",
	}

	nonGenFiles := []string{
		"core/logger/logger.go",
		"core/services/cron/cron.go",
		"tools/ci-testshard/main.go",
	}

	b.Run("ProtoMatched", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			err := generate.Run(b.Context(), repoRoot, protoFiles, cfg)
			require.NoError(b, err)
		}
	})

	b.Run("ModGraphMatched", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			err := generate.Run(b.Context(), repoRoot, modFiles, cfg)
			require.NoError(b, err)
		}
	})

	b.Run("NonGenerateSkipped", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			err := generate.Run(b.Context(), repoRoot, nonGenFiles, cfg)
			require.NoError(b, err)
		}
	})
}
