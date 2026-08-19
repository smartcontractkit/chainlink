package testrunner_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/testrunner"
)

type benchRunner struct{}

func (b *benchRunner) Run(ctx context.Context, dir string, name string, args ...string) error {
	return nil
}

func BenchmarkTestrunnerRun(b *testing.B) {
	packages := []string{
		"./core/logger",
		"./core/services/cron",
		"./core/services/telemetry",
		"./tools/ci-testshard",
		"./tools/githooks/internal/modules",
	}

	cfg := testrunner.Config{
		RepoRoot: "/test/repo",
		Packages: packages,
		Short:    true,
		Executor: &benchRunner{},
		Stdout:   io.Discard,
		Stderr:   io.Discard,
	}

	b.ReportAllocs()

	for b.Loop() {
		err := testrunner.Run(b.Context(), cfg)
		require.NoError(b, err)
	}
}
