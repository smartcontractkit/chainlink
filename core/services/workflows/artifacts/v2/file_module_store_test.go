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
	s, err := NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	binary := []byte("fake-wasm-binary-content")
	require.NoError(t, s.StoreModule("wf-1", binary, "v1.2.3"))

	p, ver, ok, err := s.GetModule("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "v1.2.3", ver)

	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, binary, got)
}

func TestStore_Overwrite(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", []byte("old"), "v1"))
	require.NoError(t, s.StoreModule("wf-1", []byte("new"), "v2"))

	p, ver, ok, err := s.GetModule("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "v2", ver)
	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, []byte("new"), got)
}

func TestStore_MissingModule(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	p, ver, ok, err := s.GetModule("nonexistent")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Empty(t, p)
	assert.Empty(t, ver)
}

func TestStore_LegacyEntryWithoutEngineVersion(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileModuleStore(dir, false)
	require.NoError(t, err)

	require.NoError(t, os.MkdirAll(s.workflowDir("wf-legacy"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(s.workflowDir("wf-legacy"), binaryFileName), []byte("legacy"), 0o600))

	p, ver, ok, err := s.GetModule("wf-legacy")
	require.NoError(t, err)
	require.True(t, ok)
	assert.NotEmpty(t, p)
	assert.Empty(t, ver)
}

func TestStore_DeleteModule(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", []byte("data"), "v1"))
	require.NoError(t, s.DeleteModule("wf-1"))

	_, _, ok, err := s.GetModule("wf-1")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestStore_DeleteNonExistent(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	assert.NoError(t, s.DeleteModule("never-stored"))
}

func TestStore_AtomicWrite(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileModuleStore(dir, false)
	require.NoError(t, err)

	require.NoError(t, s.StoreModule("wf-1", []byte("good"), "v1"))

	tmpPath := filepath.Join(s.workflowDir("wf-1"), binaryFileName+".tmp")
	require.NoError(t, os.WriteFile(tmpPath, []byte("partial"), 0o600))

	p, _, ok, err := s.GetModule("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, []byte("good"), got)
}

func TestStore_CleanOnStartup(t *testing.T) {
	dir := t.TempDir()
	stale := filepath.Join(dir, "stale-wf", binaryFileName)
	require.NoError(t, os.MkdirAll(filepath.Dir(stale), 0o755))
	require.NoError(t, os.WriteFile(stale, []byte("leftover"), 0o600))

	s, err := NewFileModuleStore(dir, true)
	require.NoError(t, err)

	_, err = os.Stat(stale)
	require.ErrorIs(t, err, os.ErrNotExist)

	require.NoError(t, s.StoreModule("wf-1", []byte("fresh"), "v1"))
	p, ver, ok, err := s.GetModule("wf-1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "v1", ver)
	got, err := os.ReadFile(p)
	require.NoError(t, err)
	assert.Equal(t, []byte("fresh"), got)
}

func TestStore_ConcurrentAccess(t *testing.T) {
	s, err := NewFileModuleStore(t.TempDir(), false)
	require.NoError(t, err)

	const wfSuffix = "ABCDEFGHIJ"
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			wfID := "wf-" + string(wfSuffix[idx])
			assert.NoError(t, s.StoreModule(wfID, []byte("data"), "v1"))
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			wfID := "wf-" + string(wfSuffix[idx])
			_, _, _, err := s.GetModule(wfID)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
}
