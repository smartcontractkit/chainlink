package lint_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/lint"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

type mockExecutor struct {
	runs []mockRun
	err  error
	// errs, when set, maps working directory to the error returned for it.
	errs map[string]error
}

type mockRun struct {
	dir  string
	name string
	args []string
}

func (m *mockExecutor) Run(ctx context.Context, dir, name string, args ...string) error {
	m.runs = append(m.runs, mockRun{dir: dir, name: name, args: args})
	if m.errs != nil {
		return m.errs[dir]
	}
	return m.err
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("runs linter on specific packages in each affected module with allow-parallel-runners", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := lint.Config{
			RepoRoot: "/repo",
			Targets: []modules.ModulePackages{
				{
					Module:   ".",
					Packages: []string{"./core/logger", "./core/services"},
				},
				{
					Module:   "deployment",
					Packages: []string{"./environment"},
				},
			},
			Fix:      true,
			Rev:      "HEAD",
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := lint.Run(t.Context(), cfg)
		require.NoError(t, err)

		require.Len(t, mock.runs, 2)

		assert.Equal(t, "/repo", mock.runs[0].dir)
		assert.Equal(t, "golangci-lint", mock.runs[0].name)
		assert.Equal(t, []string{"run", "--allow-parallel-runners", "--new-from-rev=HEAD", "--fix", "./core/logger", "./core/services"}, mock.runs[0].args)

		assert.Equal(t, "/repo/deployment", mock.runs[1].dir)
		assert.Equal(t, "golangci-lint", mock.runs[1].name)
		assert.Equal(t, []string{"run", "--allow-parallel-runners", "--new-from-rev=HEAD", "--fix", "./environment"}, mock.runs[1].args)
	})

	t.Run("returns error when linter fails", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{err: errors.New("lint failed")}
		var out bytes.Buffer

		cfg := lint.Config{
			RepoRoot: "/repo",
			Targets: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/logger"}},
			},
			Fix:      false,
			Rev:      "HEAD",
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := lint.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "golangci-lint failed in .")
	})

	t.Run("continues remaining modules when one fails and aggregates errors", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{errs: map[string]error{
			"/repo": errors.New("lint failed"),
		}}
		var out bytes.Buffer

		cfg := lint.Config{
			RepoRoot: "/repo",
			Targets: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/logger"}},
				{Module: "deployment", Packages: []string{"./environment"}},
			},
			Fix:      true,
			Rev:      "HEAD",
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := lint.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "golangci-lint failed in .")

		// The healthy module must still be linted (and fixed) despite the failure.
		require.Len(t, mock.runs, 2)
		assert.Equal(t, "/repo/deployment", mock.runs[1].dir)
	})

	t.Run("no targets to lint", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := lint.Config{
			RepoRoot: "/repo",
			Targets:  []modules.ModulePackages{},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := lint.Run(t.Context(), cfg)
		require.NoError(t, err)
		assert.Empty(t, mock.runs)
	})

	t.Run("nil stdout and stderr are replaced with discard", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}

		cfg := lint.Config{
			RepoRoot: "/repo",
			Targets: []modules.ModulePackages{
				{
					Module:   "deployment",
					Packages: []string{"./environment"},
				},
			},
			Executor: mock,
		}

		err := lint.Run(t.Context(), cfg)
		require.NoError(t, err)
		require.Len(t, mock.runs, 1)
	})
}
