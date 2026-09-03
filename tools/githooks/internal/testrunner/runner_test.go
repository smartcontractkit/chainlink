package testrunner_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/modules"
	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/testrunner"
)

type mockExecutor struct {
	runs []mockRun
	err  error
}

type mockRun struct {
	dir  string
	name string
	args []string
}

func (m *mockExecutor) Run(ctx context.Context, dir, name string, args ...string) error {
	m.runs = append(m.runs, mockRun{dir: dir, name: name, args: args})
	return m.err
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("runs test binary for root module and go test for submodules", func(t *testing.T) {
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

		require.Len(t, mock.runs, 2)

		assert.Equal(t, "/repo", mock.runs[0].dir)
		assert.Equal(t, "/repo/tools/test/.bin/test", mock.runs[0].name)
		assert.Equal(t, []string{"-short", "./core/logger"}, mock.runs[0].args)

		assert.Equal(t, "/repo/tools/githooks", mock.runs[1].dir)
		assert.Equal(t, "go", mock.runs[1].name)
		assert.Equal(t, []string{"test", "-short", "./internal/generate"}, mock.runs[1].args)
	})

	t.Run("returns error when test runner fails", func(t *testing.T) {
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
			},
			Short:    true,
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(t.Context(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tests failed on .")
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
