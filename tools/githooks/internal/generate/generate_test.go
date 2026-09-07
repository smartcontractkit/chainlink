package generate_test

import (
	"context"
	"errors"
	"os"
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

	t.Run("runs go generate on package with changed .go file containing go:generate directive", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/smartcontractkit/chainlink/v2\n"), 0o600))
		typesPkg := filepath.Join(root, "core/capabilities/remote/types")
		require.NoError(t, os.MkdirAll(typesPkg, 0o750))
		typesFile := filepath.Join(typesPkg, "types.go")
		require.NoError(t, os.WriteFile(typesFile, []byte("package types\n\n//go:generate protoc --proto_path=. messages.proto\n"), 0o600))

		var calledCmds [][]string
		runner := func(ctx context.Context, dir string, args ...string) error {
			calledCmds = append(calledCmds, args)
			return nil
		}

		err := generate.Run(t.Context(), root, []string{"core/capabilities/remote/types/types.go"}, generate.Config{Runner: runner})
		require.NoError(t, err)

		require.Len(t, calledCmds, 1)
		assert.Equal(t, []string{"generate", "./core/capabilities/remote/types"}, calledCmds[0])
	})

	t.Run("runs go generate on core/web when operator_ui/TAG changes", func(t *testing.T) {
		t.Parallel()

		var calledCmds [][]string
		runner := func(ctx context.Context, dir string, args ...string) error {
			calledCmds = append(calledCmds, args)
			return nil
		}

		files := []string{"operator_ui/TAG"}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)

		require.Len(t, calledCmds, 1)
		assert.Equal(t, []string{"generate", "./core/web"}, calledCmds[0])
	})

	t.Run("skips non-code assets without calling runner", func(t *testing.T) {
		t.Parallel()

		called := false
		runner := func(ctx context.Context, dir string, args ...string) error {
			called = true
			return nil
		}

		files := []string{"README.md", "docs/index.html", "core/web/assets/app.css", "package.json"}

		err := generate.Run(t.Context(), repoRoot, files, generate.Config{Runner: runner})
		require.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("runs go generate concurrently across multiple packages", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/smartcontractkit/chainlink/v2\n"), 0o600))

		pkg1 := filepath.Join(root, "core/services/pkg1")
		pkg2 := filepath.Join(root, "core/services/pkg2")
		require.NoError(t, os.MkdirAll(pkg1, 0o750))
		require.NoError(t, os.MkdirAll(pkg2, 0o750))

		require.NoError(t, os.WriteFile(filepath.Join(pkg1, "a.go"), []byte("package pkg1\n\n//go:generate echo pkg1\n"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(pkg2, "b.go"), []byte("package pkg2\n\n//go:generate echo pkg2\n"), 0o600))

		var (
			mu         sync.Mutex
			calledCmds [][]string
		)
		runner := func(ctx context.Context, dir string, args ...string) error {
			mu.Lock()
			calledCmds = append(calledCmds, args)
			mu.Unlock()
			return nil
		}

		err := generate.Run(t.Context(), root, []string{"core/services/pkg1/a.go", "core/services/pkg2/b.go"}, generate.Config{Runner: runner})
		require.NoError(t, err)

		mu.Lock()
		defer mu.Unlock()
		assert.Len(t, calledCmds, 2)
		assert.Contains(t, calledCmds, []string{"generate", "./core/services/pkg1"})
		assert.Contains(t, calledCmds, []string{"generate", "./core/services/pkg2"})
	})
}
