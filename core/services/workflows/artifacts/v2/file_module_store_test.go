package v2

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_RoundTrip(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	binary := []byte("fake-wasm-binary-content")
	require.NoError(t, s.StoreModule("wf-1", "bin-abc", binary))

	p, ok, err := s.GetModulePath("wf-1")
	require.NoError(t, err)
	require.True(t, ok)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, binary, got)
}

func TestStore_GetBinaryID(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", "abc123", []byte("bin")))

	id, ok, err := s.GetBinaryID("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "abc123", id)
}

func TestStore_Overwrite(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", "old-id", []byte("old")))
	require.NoError(t, s.StoreModule("wf-1", "new-id", []byte("new")))

	id, ok, err := s.GetBinaryID("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "new-id", id)

	p, ok, err := s.GetModulePath("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, []byte("new"), got)
}

func TestStore_MissingModule(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	p, ok, err := s.GetModulePath("nonexistent")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Empty(t, p)

	id, ok, err := s.GetBinaryID("nonexistent")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Empty(t, id)
}

func TestStore_DeleteModule(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", "bin-1", []byte("data")))
	require.NoError(t, s.DeleteModule("wf-1"))

	_, ok, err := s.GetModulePath("wf-1")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestStore_DeleteNonExistent(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	assert.NoError(t, s.DeleteModule("never-stored"))
}

func TestStore_AtomicWrite(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileModuleStore(dir)
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", "bin-1", []byte("good")))

	// Simulate a crash by creating the tmp file but not renaming it
	tmpPath := filepath.Join(s.workflowDir("wf-1"), binaryFileName+".tmp")
	require.NoError(t, os.WriteFile(tmpPath, []byte("partial"), 0o600))

	// The original binary should still be intact
	p, ok, err := s.GetModulePath("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, []byte("good"), got)
}

func TestStore_ConcurrentAccess(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir())
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			wfID := "wf-" + string(rune('A'+idx))
			assert.NoError(t, s.StoreModule(wfID, "bin", []byte("data")))
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			wfID := "wf-" + string(rune('A'+idx))
			_, _, err := s.GetModulePath(wfID)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
}
