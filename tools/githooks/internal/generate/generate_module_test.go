package generate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/generate"
)

func TestRun_GoGenerateInNestedModule(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/smartcontractkit/chainlink/v2\n"), 0o600))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "tools/githooks/internal/generate"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(root, "tools/githooks/go.mod"), []byte("module github.com/smartcontractkit/chainlink/v2/tools/githooks\n"), 0o600))

	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"tools/githooks/internal/generate/generate.go"}, cfg)
	require.NoError(t, err)

	generateRuns := rec.runArgs("generate")
	require.Len(t, generateRuns, 1)
	assert.Equal(t, []string{"generate", "./internal/generate"}, generateRuns[0])

	rec.mu.Lock()
	dir := rec.dirs[0]
	rec.mu.Unlock()
	assert.Equal(t, filepath.Join(root, "tools/githooks"), dir)
}

func TestRun_GoGenerateForRootModuleFile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/smartcontractkit/chainlink/v2\n"), 0o600))

	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"core/capabilities/remote/types/messages.proto"}, cfg)
	require.NoError(t, err)

	generateRuns := rec.runArgs("generate")
	require.Len(t, generateRuns, 1)
	assert.Equal(t, []string{"generate", "./core/capabilities/remote/types"}, generateRuns[0])

	rec.mu.Lock()
	dir := rec.dirs[0]
	rec.mu.Unlock()
	assert.Equal(t, root, dir)
}
