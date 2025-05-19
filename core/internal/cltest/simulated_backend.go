package cltest

import (
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func NewSimulatedBackend(t *testing.T, alloc types.GenesisAlloc, gasLimit uint64) evmtypes.Backend {
	backend := simulated.NewBackend(alloc, simulated.WithBlockGasLimit(gasLimit))
	// NOTE: Make sure to finish closing any application/client before
	// backend.Close or they can hang
	t.Cleanup(func() {
		logger.TestLogger(t).ErrorIfFn(backend.Close, "Error closing simulated backend")
	})

	return &syncBackend{Backend: backend}
}

type syncBackend struct {
	evmtypes.Backend
	mu sync.Mutex
}

func (s *syncBackend) Commit() common.Hash {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Backend.Commit()
}

func NewApplicationWithConfigV2OnSimulatedBlockchain(
	t testing.TB,
	cfg chainlink.GeneralConfig,
	backend evmtypes.Backend,
	flagsAndDeps ...interface{},
) *TestApplication {
	bid, err := backend.Client().ChainID(testutils.Context(t))
	require.NoError(t, err)
	if bid.Cmp(testutils.SimulatedChainID) != 0 {
		t.Fatalf("expected backend chain ID to be %s but it was %s", testutils.SimulatedChainID.String(), bid.String())
	}

	require.Zero(t, evmtest.MustGetDefaultChainID(t, cfg.EVMConfigs()).Cmp(testutils.SimulatedChainID))
	chainID := big.New(testutils.SimulatedChainID)
	client := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	flagsAndDeps = append(flagsAndDeps, client, chainID)

	//  app.Stop() will call client.Close on the simulated backend
	app := NewApplicationWithConfig(t, cfg, flagsAndDeps...)

	return app
}

// NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain is like NewApplicationWithConfigAndKeyOnSimulatedBlockchain
// but cfg should be v2, and configtest.NewGeneralConfigSimulated used to include the simulated chain (testutils.SimulatedChainID).
func NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(
	t testing.TB,
	cfg chainlink.GeneralConfig,
	backend evmtypes.Backend,
	flagsAndDeps ...interface{},
) *TestApplication {
	bid, err := backend.Client().ChainID(testutils.Context(t))
	require.NoError(t, err)
	if bid.Cmp(testutils.SimulatedChainID) != 0 {
		t.Fatalf("expected backend chain ID to be %s but it was %s", testutils.SimulatedChainID.String(), bid.String())
	}

	require.Zero(t, evmtest.MustGetDefaultChainID(t, cfg.EVMConfigs()).Cmp(testutils.SimulatedChainID))
	chainID := big.New(testutils.SimulatedChainID)
	client := client.NewSimulatedBackendClient(t, backend, testutils.SimulatedChainID)

	flagsAndDeps = append(flagsAndDeps, client, chainID)

	//  app.Stop() will call client.Close on the simulated backend
	return NewApplicationWithConfigAndKey(t, cfg, flagsAndDeps...)
}

// Mine forces the simulated backend to produce a new block every X seconds
// If you need to manually commit blocks, you must use the returned commit func, rather than calling Commit() directly,
// which will race.
func Mine(backend evmtypes.Backend, blockTime time.Duration) (commit func() common.Hash, stopMining func()) {
	timer := time.NewTicker(blockTime)
	chStop := make(chan struct{})
	commitCh := make(chan chan common.Hash)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-timer.C:
				backend.Commit()
			case hash := <-commitCh:
				hash <- backend.Commit()
			case <-chStop:
				return
			}
		}
	}()
	return func() common.Hash {
			hash := make(chan common.Hash)
			select {
			case <-chStop:
				return common.Hash{}
			case commitCh <- hash:
				return <-hash
			}
		}, func() {
			close(chStop)
			timer.Stop()
			<-done
		}
}
