package testrunner_test

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/testrunner"
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

	t.Run("runs test binary for root module and go test for submodules in parallel", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Modules: []modules.ModulePackages{
				{
					Module:   ".",
					Packages: []string{"./core/logger"},
				},
				{
					Module:   "tools/githooks",
					Packages: []string{"./internal/generate"},
				},
			},
			Short:    true,
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(t.Context(), cfg)
		require.NoError(t, err)

		mock.mu.Lock()
		runs := make([]mockRun, len(mock.runs))
		copy(runs, mock.runs)
		mock.mu.Unlock()

		require.Len(t, runs, 2)

		runDirs := []string{runs[0].dir, runs[1].dir}
		assert.Contains(t, runDirs, "/repo")
		assert.Contains(t, runDirs, "/repo/tools/githooks")
	})

	t.Run("aggregates errors when multiple modules fail", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{err: errors.New("test failed")}
		var out bytes.Buffer

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Modules: []modules.ModulePackages{
				{
					Module:   ".",
					Packages: []string{"./core/logger"},
				},
				{
					Module:   "tools/githooks",
					Packages: []string{"./internal/generate"},
				},
			},
			Short:    true,
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tests failed on .")
		assert.Contains(t, err.Error(), "tests failed on tools/githooks")
	})

	t.Run("no modules to test", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Modules:  []modules.ModulePackages{},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(t.Context(), cfg)
		require.NoError(t, err)
		assert.Empty(t, mock.runs)
	})

	t.Run("nil stdout and stderr are replaced with discard", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Modules: []modules.ModulePackages{
				{
					Module:   "tools/githooks",
					Packages: []string{"./internal/generate"},
				},
			},
			Executor: mock,
		}

		err := testrunner.Run(t.Context(), cfg)
		require.NoError(t, err)
		require.Len(t, mock.runs, 1)
	})
}
