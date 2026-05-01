package db

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/config"
	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/output"
)

func TestDumpDiagnosticsNilHandle(t *testing.T) {
	t.Parallel()
	var h *Handle
	require.NoError(t, h.DumpDiagnostics(context.Background(), t.TempDir(), 0))
}

func TestDumpDiagnosticsNoContainer(t *testing.T) {
	t.Parallel()
	h := &Handle{}
	dir := t.TempDir()
	require.NoError(t, h.DumpDiagnostics(context.Background(), dir, 0))
	_, err := os.Stat(filepath.Join(dir, "postgres-state-0.md"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestEnsureDatabaseURLSetsCLDatabaseURL(t *testing.T) {
	t.Setenv("CL_DATABASE_URL", "")
	want := "postgres://user@host/chainlink_test"
	h, err := Ensure(context.Background(), &config.App{
		PostgresVersion: "15",
		DatabaseURL:     want,
		AIOutput:        true,
	}, output.New(true, io.Discard, io.Discard, output.SkipFD))
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, want, os.Getenv("CL_DATABASE_URL"))
}

func TestEnsureDatabaseURLConflictsWithEnv(t *testing.T) {
	t.Setenv("CL_DATABASE_URL", "postgres://already/set")
	_, err := Ensure(context.Background(), &config.App{
		PostgresVersion: "15",
		DatabaseURL:     "postgres://other/db",
		AIOutput:        true,
	}, output.New(true, io.Discard, io.Discard, output.SkipFD))
	require.Error(t, err)
	assert.ErrorContains(t, err, "CL_DATABASE_URL")
}

func TestEnsureRequiresPostgresVersion(t *testing.T) {
	t.Parallel()
	_, err := Ensure(context.Background(), &config.App{
		PostgresVersion: "",
		AIOutput:        true,
	}, output.New(true, io.Discard, io.Discard, output.SkipFD))
	require.Error(t, err)
	assert.ErrorContains(t, err, "postgres version is required")
}

func TestHandleEnvIncludesConnectionString(t *testing.T) {
	t.Parallel()
	h := &Handle{connStr: "postgres://worker/db"}
	assert.Equal(t, []string{"CL_DATABASE_URL=postgres://worker/db"}, h.Env())
}

func TestEnsurePoolRejectsExternalDatabaseWhenParallel(t *testing.T) {
	t.Parallel()
	_, err := EnsurePool(context.Background(), &config.App{
		PostgresVersion:    "15",
		DatabaseURL:        "postgres://example/db",
		ParallelIterations: 2,
		AIOutput:           true,
	}, output.New(true, io.Discard, io.Discard, output.SkipFD), 2)
	require.Error(t, err)
	assert.ErrorContains(t, err, "--parallel-iterations > 1 cannot be used with --database-url")
}

func TestEnsurePoolCreatesSingleExternalDatabaseHandle(t *testing.T) {
	t.Setenv("CL_DATABASE_URL", "")
	pool, err := EnsurePool(context.Background(), &config.App{
		PostgresVersion:    "15",
		DatabaseURL:        "postgres://example/db",
		ParallelIterations: 1,
		AIOutput:           true,
	}, output.New(true, io.Discard, io.Discard, output.SkipFD), 1)
	require.NoError(t, err)
	require.Len(t, pool.Handles(), 1)
	assert.Equal(t, []string{"CL_DATABASE_URL=postgres://example/db"}, pool.Handles()[0].Env())
	assert.Equal(t, "postgres://example/db", os.Getenv("CL_DATABASE_URL"))
}
