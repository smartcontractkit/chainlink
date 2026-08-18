package testrunner_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (m *mockExecutor) Run(ctx context.Context, dir string, name string, args ...string) error {
	m.runs = append(m.runs, mockRun{dir: dir, name: name, args: args})
	return m.err
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("runs test binary with short flag and packages", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Packages: []string{"./core/logger", "./core/services/telemetry"},
			Short:    true,
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(context.Background(), cfg)
		require.NoError(t, err)

		require.Len(t, mock.runs, 1)
		assert.Equal(t, "/repo", mock.runs[0].dir)
		assert.Equal(t, "/repo/tools/test/.bin/test", mock.runs[0].name)
		assert.Equal(t, []string{"-short", "./core/logger", "./core/services/telemetry"}, mock.runs[0].args)
	})

	t.Run("returns error when test runner fails", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{err: errors.New("test failed")}
		var out bytes.Buffer

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Packages: []string{"./core/logger"},
			Short:    true,
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(context.Background(), cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "test failed")
	})

	t.Run("no packages to test", func(t *testing.T) {
		t.Parallel()

		mock := &mockExecutor{}
		var out bytes.Buffer

		cfg := testrunner.Config{
			RepoRoot: "/repo",
			Packages: []string{},
			Executor: mock,
			Stdout:   &out,
			Stderr:   &out,
		}

		err := testrunner.Run(context.Background(), cfg)
		require.NoError(t, err)
		assert.Empty(t, mock.runs)
	})
}
