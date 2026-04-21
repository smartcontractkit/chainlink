package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpStateNilHandle(t *testing.T) {
	t.Parallel()
	var h *Handle
	require.NoError(t, h.DumpState(context.Background(), t.TempDir(), 0))
}

func TestDumpStateNoContainer(t *testing.T) {
	t.Parallel()
	h := &Handle{}
	dir := t.TempDir()
	require.NoError(t, h.DumpState(context.Background(), dir, 0))
	_, err := os.Stat(filepath.Join(dir, "postgres-state-0.md"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}
