package compute

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	simpleBinaryCmd = "core/capabilities/compute/test/simple/cmd"
)

// Verify that cache evicts an expired module.
func TestCache(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-558")

	t.Parallel()
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second
	reapTicker := make(chan time.Time)

	cache := newModuleCache(clock, tick, timeout, 0)
	cache.onReaper = make(chan struct{}, 1)
	cache.reapTicker = reapTicker
	cache.start()
	defer cache.close()

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, false, t)
	hmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         logger.TestLogger(t),
		IsUncompressed: true,
	}, binary)
	require.NoError(t, err)

	id := uuid.New().String()
	mod := &module{
		module: hmod,
	}
	cache.add(id, mod)

	got, ok := cache.get(id)
	assert.True(t, ok)

	assert.Equal(t, got, mod)

	clock.Advance(15 * time.Second)
	reapTicker <- time.Now()
	<-cache.onReaper
	_, ok = cache.get(id)
	assert.False(t, ok)
}

// Verify that an expired module is not evicted because evictAfterSize is 1
func TestCache_EvictAfterSize(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second
	reapTicker := make(chan time.Time)

	cache := newModuleCache(clock, tick, timeout, 1)
	cache.onReaper = make(chan struct{}, 1)
	cache.reapTicker = reapTicker
	cache.start()
	defer cache.close()

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, false, t)
	hmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         logger.TestLogger(t),
		IsUncompressed: true,
	}, binary)
	require.NoError(t, err)

	id := uuid.New().String()
	mod := &module{
		module: hmod,
	}
	cache.add(id, mod)
	assert.Len(t, cache.m, 1)

	got, ok := cache.get(id)
	assert.True(t, ok)

	assert.Equal(t, got, mod)

	clock.Advance(15 * time.Second)
	reapTicker <- time.Now()
	select {
	case <-t.Context().Done():
		return
	case <-cache.onReaper:
	}
	_, ok = cache.get(id)
	assert.True(t, ok)
}

func TestCache_AddDuplicatedModule(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-574")

	t.Parallel()
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second
	reapTicker := make(chan time.Time)

	cache := newModuleCache(clock, tick, timeout, 0)
	cache.onReaper = make(chan struct{}, 1)
	cache.reapTicker = reapTicker
	cache.start()
	defer cache.close()

	simpleBinary := wasmtest.CreateTestBinary(simpleBinaryCmd, false, t)
	shmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         logger.TestLogger(t),
		IsUncompressed: true,
	}, simpleBinary)
	require.NoError(t, err)

	// we will use the same id for both modules, but should only be associated to the simple module
	duplicatedID := uuid.New().String()
	smod := &module{
		module: shmod,
	}
	err = cache.add(duplicatedID, smod)
	require.NoError(t, err)

	got, ok := cache.get(duplicatedID)
	assert.True(t, ok)
	assert.Equal(t, got, smod)

	// Adding a different module but with the same id should not overwrite the existing module
	fetchBinary := wasmtest.CreateTestBinary(fetchBinaryCmd, false, t)
	fhmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         logger.TestLogger(t),
		IsUncompressed: true,
	}, fetchBinary)
	require.NoError(t, err)

	fmod := &module{
		module: fhmod,
	}
	err = cache.add(duplicatedID, fmod)
	require.ErrorContains(t, err, fmt.Sprintf("module with id %q already exists in cache", duplicatedID))

	// validate that the module is still the same
	got, ok = cache.get(duplicatedID)
	assert.True(t, ok)
	assert.Equal(t, got, smod)
}
