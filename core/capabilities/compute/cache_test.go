package compute

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/wasmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	binaryLocation = "test/simple/cmd/testmodule.wasm"
	binaryCmd      = "core/capabilities/compute/test/simple/cmd"
)

func TestCache(t *testing.T) {
	clock := clockwork.NewFakeClock()
	tick := 1 * time.Second
	timeout := 1 * time.Second

	cache := newModuleCache(clock, tick, timeout)
	cache.onReaper = make(chan struct{}, 1)
	cache.start()
	defer cache.close()

	var binary []byte
	binary = wasmtest.CreateTestBinary(binaryCmd, binaryLocation, false, t)
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

	got, err := cache.get(id)
	require.NoError(t, err)

	assert.Equal(t, got, mod)

	clock.Advance(15 * time.Second)
	<-cache.onReaper
	m, err := cache.get(id)
	assert.ErrorContains(t, err, "could not find module", m)
}
