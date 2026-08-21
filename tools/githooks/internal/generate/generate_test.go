package generate_test

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/generate"
)

func TestRun(t *testing.T) {
	t.Parallel()

	repoRoot := "/test/repo"

	t.Run("runs go generate on package with changed proto file", func(t *testing.T) {
		t.Parallel()

		var (
			mu         sync.Mutex
			calledCmds [][]string
			calledDirs []string
		)

		runner := func(ctx context.Context, dir string, args ...string) error {
			mu.Lock()
			calledDirs = append(calledDirs, dir)
			calledCmds = append(calledCmds, args)
			mu.Unlock()
			return nil
		}

		files := []string{
			"core/capabilities/remote/types/messages.proto",
			"core/services/llo/telem/telem_streams.proto",
		}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)

		mu.Lock()
		defer mu.Unlock()

		assert.Len(t, calledCmds, 2)
		assert.Contains(t, calledCmds, []string{"generate", "./core/capabilities/remote/types"})
		assert.Contains(t, calledCmds, []string{"generate", "./core/services/llo/telem"})
	})

	t.Run("runs config docs generation when core/config changes", func(t *testing.T) {
		t.Parallel()

		var calledCmds [][]string
		runner := func(ctx context.Context, dir string, args ...string) error {
			calledCmds = append(calledCmds, args)
			return nil
		}

		files := []string{"core/config/toml/types.go"}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)

		require.Len(t, calledCmds, 1)
		assert.Equal(t, []string{"run", "./core/config/docs/cmd/generate", "-o", "./docs/"}, calledCmds[0])
	})

	t.Run("runs modgraph go.md generation when go.mod or go.sum changes", func(t *testing.T) {
		t.Parallel()

		var calledCmds [][]string
		runner := func(ctx context.Context, dir string, args ...string) error {
			calledCmds = append(calledCmds, args)
			return nil
		}

		files := []string{"go.mod", "deployment/go.sum"}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)

		require.Len(t, calledCmds, 1)
		assert.Equal(t, []string{"modgraph"}, calledCmds[0])
	})

	t.Run("skips execution when no generate targets matched", func(t *testing.T) {
		t.Parallel()

		called := false
		runner := func(ctx context.Context, dir string, args ...string) error {
			called = true
			return nil
		}

		files := []string{"core/logger/logger.go", "core/services/cron/cron.go"}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("returns error when generator fails", func(t *testing.T) {
		t.Parallel()

		runner := func(ctx context.Context, dir string, args ...string) error {
			return errors.New("protoc failed")
		}

		files := []string{"core/capabilities/remote/types/messages.proto"}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "protoc failed")
	})

	t.Run("handles absolute file paths", func(t *testing.T) {
		t.Parallel()

		var calledCmds [][]string
		runner := func(ctx context.Context, dir string, args ...string) error {
			calledCmds = append(calledCmds, args)
			return nil
		}

		absProto := filepath.Join(repoRoot, "core/services/nodestatusreporter/bridgestatus/events/bridge_status.proto")
		files := []string{absProto}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)

		require.Len(t, calledCmds, 1)
		assert.Equal(t, []string{"generate", "./core/services/nodestatusreporter/bridgestatus/events"}, calledCmds[0])
	})
}
