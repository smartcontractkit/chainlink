package lint_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/lint"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

type benchExecutor struct{}

func (b *benchExecutor) Run(ctx context.Context, dir string, name string, args ...string) error {
	return nil
}

func BenchmarkLintRun(b *testing.B) {
	targets := []modules.ModulePackages{
		{Module: ".", Packages: []string{"./core/logger", "./core/services/cron"}},
		{Module: "deployment", Packages: []string{"./environment"}},
		{Module: "tools/githooks", Packages: []string{"./internal/modules", "./internal/lint"}},
	}

	cfg := lint.Config{
		RepoRoot: "/test/repo",
		Targets:  targets,
		Fix:      true,
		Rev:      "HEAD",
		Executor: &benchExecutor{},
		Stdout:   io.Discard,
		Stderr:   io.Discard,
	}

	b.ReportAllocs()

	for b.Loop() {
		err := lint.Run(b.Context(), cfg)
		require.NoError(b, err)
	}
}
