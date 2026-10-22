package v2

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
)

func newLazyModuleTestHandler(t *testing.T, withStore bool) *eventHandler {
	t.Helper()
	h := &eventHandler{lggr: logger.TestLogger(t), moduleEngineVersion: "v1"}
	if withStore {
		store, err := artifacts.NewFileModuleStore(t.TempDir(), false)
		require.NoError(t, err)
		h.moduleStore = store
		h.moduleLRU = NewModuleLRU(clockwork.NewFakeClock())
	}
	return h
}

func TestCreateEngineModule_Lazy(t *testing.T) {
	t.Parallel()
	// not a valid WASM binary: any attempt to compile it fails
	binary := []byte("not-a-wasm-binary")
	cfg := func(t *testing.T) *host.ModuleConfig { return &host.ModuleConfig{Logger: logger.TestLogger(t)} }

	t.Run("defers compilation and registers an unloaded evictable module", func(t *testing.T) {
		t.Parallel()
		h := newLazyModuleTestHandler(t, true)

		m, err := h.createEngineModule(t.Context(), "wf-lazy", binary, cfg(t), true)
		require.NoError(t, err)

		evictable, ok := m.(*EvictableModule)
		require.True(t, ok)
		assert.False(t, evictable.IsLoaded())
		assert.True(t, h.moduleLRU.Contains("wf-lazy"))
		_, version, found, err := h.moduleStore.GetModule("wf-lazy")
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "v1", version)
	})

	t.Run("first execution surfaces the compile error", func(t *testing.T) {
		t.Parallel()
		h := newLazyModuleTestHandler(t, true)

		m, err := h.createEngineModule(t.Context(), "wf-lazy", binary, cfg(t), true)
		require.NoError(t, err)
		t.Cleanup(m.Close)

		_, err = m.Execute(t.Context(), nil, nil)
		require.ErrorContains(t, err, "failed to create module on reload")
	})

	t.Run("compiles eagerly when not lazy", func(t *testing.T) {
		t.Parallel()
		h := newLazyModuleTestHandler(t, true)

		_, err := h.createEngineModule(t.Context(), "wf-eager", binary, cfg(t), false)
		require.Error(t, err)
		assert.False(t, h.moduleLRU.Contains("wf-eager"))
	})

	t.Run("compiles eagerly without a module store", func(t *testing.T) {
		t.Parallel()
		h := newLazyModuleTestHandler(t, false)

		_, err := h.createEngineModule(t.Context(), "wf-nostore", binary, cfg(t), true)
		require.Error(t, err)
	})
}
