package actionlint_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/actionlint"
)

type benchExecutor struct{}

func (b *benchExecutor) Run(ctx context.Context, dir, name string, args ...string) error {
	return nil
}

func (b *benchExecutor) Output(ctx context.Context, dir, name string, args ...string) ([]byte, error) {
	return []byte(validKjanatHelpOutput), nil
}

func BenchmarkActionlintRun(b *testing.B) {
	files := []string{
		".github/workflows/ci.yml",
		".github/workflows/deploy.yaml",
		"core/services/app.go",
		"README.md",
	}

	cfg := actionlint.Config{
		RepoRoot: "/test/repo",
		Files:    files,
		Executor: &benchExecutor{},
		Stdout:   io.Discard,
		Stderr:   io.Discard,
	}

	b.ReportAllocs()

	for b.Loop() {
		err := actionlint.Run(b.Context(), cfg)
		require.NoError(b, err)
	}
}
