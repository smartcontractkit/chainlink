package compute

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	simpleBinaryLocation = "test/simple/cmd/testmodule.wasm"
	simpleBinaryCmd      = "core/capabilities/compute/test/simple/cmd"
)

// Verify that cache evicts an expired module.
func TestCache(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second
	reapTicker := make(chan time.Time)
	lggr := logger.TestLogger(t)

	cache := newModuleCache(clock, tick, timeout, 0, lggr)
	cache.onReaper = make(chan struct{}, 1)
	cache.reapTicker = reapTicker
	cache.start()
	defer cache.close()

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, simpleBinaryLocation, false, t)
	hmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         lggr,
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
	ctx := tests.Context(t)
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second
	reapTicker := make(chan time.Time)
	lggr := logger.TestLogger(t)

	cache := newModuleCache(clock, tick, timeout, 1, lggr)
	cache.onReaper = make(chan struct{}, 1)
	cache.reapTicker = reapTicker
	cache.start()
	defer cache.close()

	binary := wasmtest.CreateTestBinary(simpleBinaryCmd, simpleBinaryLocation, false, t)
	hmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         lggr,
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
	case <-ctx.Done():
		return
	case <-cache.onReaper:
	}
	_, ok = cache.get(id)
	assert.True(t, ok)
}

func TestCache_AddDuplicatedModule(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second
	reapTicker := make(chan time.Time)
	lggr, logs := logger.TestLoggerObserved(t, zapcore.WarnLevel)

	cache := newModuleCache(clock, tick, timeout, 0, lggr)
	cache.onReaper = make(chan struct{}, 1)
	cache.reapTicker = reapTicker
	cache.start()
	defer cache.close()

	simpleBinary := wasmtest.CreateTestBinary(simpleBinaryCmd, simpleBinaryLocation, false, t)
	shmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         lggr,
		IsUncompressed: true,
	}, simpleBinary)
	require.NoError(t, err)

	// we will use the same id for both modules, but should only be associated to the simple module
	duplicatedID := uuid.New().String()
	smod := &module{
		module: shmod,
	}
	cache.add(duplicatedID, smod)

	got, ok := cache.get(duplicatedID)
	assert.True(t, ok)
	assert.Equal(t, got, smod)
	assert.Empty(t, logs.All())

	// Adding a different module but with the same id should not overwrite the existing module
	fetchBinary := wasmtest.CreateTestBinary(fetchBinaryCmd, fetchBinaryLocation, false, t)
	fhmod, err := host.NewModule(&host.ModuleConfig{
		Logger:         lggr,
		IsUncompressed: true,
	}, fetchBinary)
	require.NoError(t, err)

	fmod := &module{
		module: fhmod,
	}
	cache.add(duplicatedID, fmod)

	// validate that the module is still the same
	got, ok = cache.get(duplicatedID)
	assert.True(t, ok)
	assert.Equal(t, got, smod)

	// validate that the log was emitted
	assert.Equal(t, 1, len(logs.All()))
	expectedLog := fmt.Sprintf("module with id %q already exists in cache, skipping...", duplicatedID)
	assert.Equal(t, expectedLog, logs.All()[0].Entry.Message)
}
