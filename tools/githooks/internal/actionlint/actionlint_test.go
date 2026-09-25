package actionlint_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/actionlint"
)

type mockExecutor struct {
	helpOutput []byte
	helpErr    error
	runErr     error
	runs       []mockRun
}

type mockRun struct {
	dir  string
	name string
	args []string
}

func (m *mockExecutor) Run(ctx context.Context, dir, name string, args ...string) error {
	m.runs = append(m.runs, mockRun{dir: dir, name: name, args: args})
	return m.runErr
}

func (m *mockExecutor) Output(ctx context.Context, dir, name string, args ...string) ([]byte, error) {
	m.runs = append(m.runs, mockRun{dir: dir, name: name, args: args})
	if m.helpErr != nil {
		return nil, m.helpErr
	}
	return m.helpOutput, nil
}

const validKjanatHelpOutput = `Usage: actionlint [FLAGS] [FILES...] [-]

  actionlint is a linter for GitHub Actions workflow files.

Documents:

  - List of checks: https://github.com/kjanat/actionlint/tree/main/docs/checks.md
  - Usage:          https://github.com/kjanat/actionlint/tree/main/docs/usage.md
  - Configuration:  https://github.com/kjanat/actionlint/tree/main/docs/config.md
`

const unmaintainedRhysdHelpOutput = `Usage: actionlint [FLAGS] [FILES...] [-]

  actionlint is a linter for GitHub Actions workflow files.

Documents:

  - List of checks: https://github.com/rhysd/actionlint/tree/main/docs/checks.md
  - Usage:          https://github.com/rhysd/actionlint/tree/main/docs/usage.md
  - Configuration:  https://github.com/rhysd/actionlint/tree/main/docs/config.md
`

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("skips execution when no .github YAML files are changed", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(validKjanatHelpOutput),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				"core/services/app.go",
				"README.md",
				".github/CODEOWNERS",
				".github/in-memory-tests.json",
				".github/scripts/test.sh",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.NoError(t, err)
		assert.Empty(t, mock.runs)
	})

	t.Run("fails when unmaintained rhysd/actionlint is detected", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(unmaintainedRhysdHelpOutput),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/workflows/ci.yml",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmaintained")
		assert.Contains(t, err.Error(), "rhysd/actionlint")
		assert.Contains(t, err.Error(), "https://github.com/kjanat/actionlint")
	})

	t.Run("fails with install instructions when actionlint binary is missing", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpErr: errors.New("executable file not found in $PATH"),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/workflows/ci.yml",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.Contains(t, err.Error(), "https://github.com/kjanat/actionlint")
	})

	t.Run("runs actionlint with specific files when only workflow files changed", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(validKjanatHelpOutput),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/workflows/ci.yml",
				".github/workflows/deploy.yaml",
				"core/services/app.go",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.NoError(t, err)

		// 1 output call for help validation, 1 run call for linting
		require.Len(t, mock.runs, 2)
		assert.Equal(t, "actionlint", mock.runs[1].name)
		assert.Equal(t, "/repo", mock.runs[1].dir)
		assert.Equal(t, []string{".github/workflows/ci.yml", ".github/workflows/deploy.yaml"}, mock.runs[1].args)
	})

	t.Run("runs actionlint on all workflows when non-workflow .github YAML changed", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(validKjanatHelpOutput),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/actionlint.yml",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.NoError(t, err)

		require.Len(t, mock.runs, 2)
		assert.Equal(t, "actionlint", mock.runs[1].name)
		assert.Equal(t, "/repo", mock.runs[1].dir)
		assert.Empty(t, mock.runs[1].args) // No file args = actionlint scans all workflows
	})

	t.Run("runs actionlint on all workflows when composite action YAML changed", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(validKjanatHelpOutput),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/actions/setup-go/action.yml",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.NoError(t, err)

		require.Len(t, mock.runs, 2)
		assert.Equal(t, "actionlint", mock.runs[1].name)
		assert.Empty(t, mock.runs[1].args)
	})

	t.Run("runs actionlint on all workflows when mixed workflow and non-workflow YAML changed", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(validKjanatHelpOutput),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/workflows/ci.yml",
				".github/actionlint.yml",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.NoError(t, err)

		require.Len(t, mock.runs, 2)
		assert.Equal(t, "actionlint", mock.runs[1].name)
		assert.Empty(t, mock.runs[1].args)
	})

	t.Run("propagates actionlint execution errors", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{
			helpOutput: []byte(validKjanatHelpOutput),
			runErr:     errors.New("syntax error in workflow"),
		}
		var out bytes.Buffer

		cfg := actionlint.Config{
			RepoRoot: "/repo",
			Files: []string{
				".github/workflows/ci.yml",
			},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := actionlint.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "syntax error in workflow")
	})
}
