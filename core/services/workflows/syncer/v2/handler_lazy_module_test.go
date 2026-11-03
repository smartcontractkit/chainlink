package v2

import (
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	capreg "github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/confidentialrelay"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
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

func TestBuildEngineConfig_SdkLabelIsPerEngine(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	limiters, err := v2.NewLimiters(limits.Factory{Logger: lggr}, nil)
	require.NoError(t, err)

	h := newLazyModuleTestHandler(t, true)
	h.emitter = custmsg.NewLabeler()
	h.engineLimiters = limiters
	h.capRegistry = capreg.NewRegistry(lggr)
	h.executionHandlers = &confidentialrelay.ExecutionHandlers{}
	h.cachedTriggerSubscriptionsEnabled = true

	subs, err := proto.Marshal(&sdk.TriggerSubscriptionRequest{Subscriptions: []*sdk.TriggerSubscription{{Id: "id_0", Method: "method"}}})
	require.NoError(t, err)
	name, err := types.NewWorkflowName("name")
	require.NoError(t, err)

	build := func(wfID, sdkName string) *v2.EngineConfig {
		cfg, buildErr := h.buildEngineConfig(t.Context(), wfID, "owner", name, "tag", nil, []byte("not-wasm"), "url", subs, sdkName, nil)
		require.NoError(t, buildErr)
		return cfg
	}
	cfgA := build("wf-a", "sdk-a")
	cfgB := build("wf-b", "sdk-b")

	assert.Equal(t, "sdk-a", cfgA.SdkName)
	assert.Equal(t, "sdk-a", cfgA.BeholderEmitter.Labels()[platform.KeySDK])
	assert.Equal(t, "sdk-b", cfgB.SdkName)
	assert.Equal(t, "sdk-b", cfgB.BeholderEmitter.Labels()[platform.KeySDK])
	assert.NotContains(t, h.emitter.Labels(), platform.KeySDK, "shared emitter must not be mutated")
}
