package cltest

import (
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func NewSimulatedBackend(t *testing.T, alloc core.GenesisAlloc, gasLimit uint32) *backends.SimulatedBackend {
	backend := backends.NewSimulatedBackend(alloc, uint64(gasLimit))
	// NOTE: Make sure to finish closing any application/client before
	// backend.Close or they can hang
	t.Cleanup(func() {
		logger.TestLogger(t).ErrorIfFn(backend.Close, "Error closing simulated backend")
	})
	return backend
}
func NewApplicationWithConfigV2OnSimulatedBlockchain(
	t testing.TB,
	cfg chainlink.GeneralConfig,
	backend *backends.SimulatedBackend,
	flagsAndDeps ...interface{},
) *TestApplication {
	if bid := backend.Blockchain().Config().ChainID; bid.Cmp(testutils.SimulatedChainID) != 0 {
		t.Fatalf("expected backend chain ID to be %s but it was %s", testutils.SimulatedChainID.String(), bid.String())
	}
	defID := cfg.DefaultChainID()
	require.Zero(t, defID.Cmp(testutils.SimulatedChainID))
	chainID := utils.NewBig(testutils.SimulatedChainID)
	client := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)
	eventBroadcaster := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, logger.TestLogger(t), uuid.NewV4())

	flagsAndDeps = append(flagsAndDeps, client, eventBroadcaster, chainID)

	//  app.Stop() will call client.Close on the simulated backend
	app := NewApplicationWithConfig(t, cfg, flagsAndDeps...)

	return app
}

// NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain is like NewApplicationWithConfigAndKeyOnSimulatedBlockchain
// but cfg should be v2, and configtest.NewGeneralConfigSimulated used to include the simulated chain (testutils.SimulatedChainID).
func NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(
	t testing.TB,
	cfg chainlink.GeneralConfig,
	backend *backends.SimulatedBackend,
	flagsAndDeps ...interface{},
) *TestApplication {
	if bid := backend.Blockchain().Config().ChainID; bid.Cmp(testutils.SimulatedChainID) != 0 {
		t.Fatalf("expected backend chain ID to be %s but it was %s", testutils.SimulatedChainID.String(), bid.String())
	}
	defID := cfg.DefaultChainID()
	require.Zero(t, defID.Cmp(testutils.SimulatedChainID))
	chainID := utils.NewBig(testutils.SimulatedChainID)
	client := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)
	eventBroadcaster := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, logger.TestLogger(t), uuid.NewV4())

	flagsAndDeps = append(flagsAndDeps, client, eventBroadcaster, chainID)

	//  app.Stop() will call client.Close on the simulated backend
	return NewApplicationWithConfigAndKey(t, cfg, flagsAndDeps...)
}

// Mine forces the simulated backend to produce a new block every 2 seconds
func Mine(backend *backends.SimulatedBackend, blockTime time.Duration) (stopMining func()) {
	timer := time.NewTicker(blockTime)
	chStop := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-timer.C:
				backend.Commit()
			case <-chStop:
				wg.Done()
				return
			}
		}
	}()
	return func() { close(chStop); timer.Stop(); wg.Wait() }
}
