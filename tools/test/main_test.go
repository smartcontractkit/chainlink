package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func findRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	require.NoError(t, err)

	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "GNUmakefile")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return cwd
}

func TestNoopResourcesMatchTestrigContract(t *testing.T) {
	t.Parallel()

	resources := noopResources(3)
	require.Len(t, resources, 3)

	for i, resource := range resources {
		assert.Empty(t, resource.Env, "resource %d env", i)
		assert.Nil(t, resource.Reset, "resource %d reset", i)
		assert.Nil(t, resource.DumpDiagnostics, "resource %d dump diagnostics", i)
		assert.Nil(t, resource.Cleanup, "resource %d cleanup", i)

		// Mirror testrig's nil checks before invoking callbacks.
		if resource.Reset != nil {
			require.NoError(t, resource.Reset(context.Background()))
		}
		if resource.DumpDiagnostics != nil {
			require.NoError(t, resource.DumpDiagnostics(context.Background(), t.TempDir(), 0))
		}
		if resource.Cleanup != nil {
			require.NoError(t, resource.Cleanup())
		}
	}
}

func TestDBProviderSkipsPostgresForConfigPackages(t *testing.T) {
	repoRoot := findRepoRoot(t)
	t.Chdir(repoRoot)
	t.Setenv("CL_DATABASE_URL", "")

	resources, err := dbProviderForArgs(context.Background(), 2, []string{"./core/config/..."})
	require.NoError(t, err)
	require.Len(t, resources, 2)
	assert.Empty(t, os.Getenv("CL_DATABASE_URL"))

	for i, resource := range resources {
		assert.Empty(t, resource.Env, "resource %d env", i)
		assert.Nil(t, resource.Reset, "resource %d reset", i)
		assert.Nil(t, resource.Cleanup, "resource %d cleanup", i)
	}
}
