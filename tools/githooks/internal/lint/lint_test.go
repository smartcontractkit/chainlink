package lint_test

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/lint"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
)

type mockExecutor struct {
	mu   sync.Mutex
	runs []mockRun
	err  error
}

type mockRun struct {
	dir  string
	name string
	args []string
}

func (m *mockExecutor) Run(ctx context.Context, dir string, name string, args ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runs = append(m.runs, mockRun{dir: dir, name: name, args: args})
	return m.err
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("runs linter on specific packages in each affected module in parallel with allow-parallel-runners", func(t *testing.T) {
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

		mock.mu.Lock()
		runs := make([]mockRun, len(mock.runs))
		copy(runs, mock.runs)
		mock.mu.Unlock()

		require.Len(t, runs, 2)

		// Order may vary due to concurrent execution, verify both modules ran
		runDirs := []string{runs[0].dir, runs[1].dir}
		assert.Contains(t, runDirs, "/repo")
		assert.Contains(t, runDirs, "/repo/deployment")

		for _, r := range runs {
			assert.Equal(t, "golangci-lint", r.name)
			assert.Contains(t, r.args, "--allow-parallel-runners")
			assert.Contains(t, r.args, "--new-from-rev=HEAD")
			assert.Contains(t, r.args, "--fix")
		}
	})

	t.Run("uses patch file when PatchFile is provided", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := lint.Config{
			RepoRoot:  "/repo",
			Targets:   []modules.ModulePackages{{Module: ".", Packages: []string{"./core/logger"}}},
			PatchFile: "/tmp/changes.patch",
			Fix:       true,
			Executor:  mock,
			Stdout:    &out,
			Stderr:    &out,
		}

		err := lint.Run(t.Context(), cfg)
		require.NoError(t, err)

		mock.mu.Lock()
		runs := mock.runs
		mock.mu.Unlock()

		require.Len(t, runs, 1)
		assert.Equal(t, []string{"run", "--allow-parallel-runners", "--new-from-patch=/tmp/changes.patch", "--fix", "./core/logger"}, runs[0].args)
	})

	t.Run("aggregates errors when multiple modules fail", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{err: errors.New("lint failed")}
		var out bytes.Buffer

		cfg := lint.Config{
			RepoRoot: "/repo",
			Targets: []modules.ModulePackages{
				{Module: ".", Packages: []string{"./core/logger"}},
				{Module: "deployment", Packages: []string{"./environment"}},
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
		assert.Contains(t, err.Error(), "golangci-lint failed in deployment")
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
