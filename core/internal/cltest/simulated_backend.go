package cltest

import (
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func NewSimulatedBackend(t *testing.T, alloc core.GenesisAlloc, gasLimit uint32) *backends.SimulatedBackend {
	backend := backends.NewSimulatedBackend(alloc, uint64(gasLimit))
	// NOTE: Make sure to finish closing any application/client before
	// backend.Close or they can hang
	t.Cleanup(func() {
		logger.TestLogger(t).ErrorIfClosing(backend, "simulated backend")
	})
	return backend
}

const SimulatedBackendEVMChainID int64 = 1337

func NewApplicationWithConfigAndKeyOnSimulatedBlockchain(
	t testing.TB,
	cfg *configtest.TestGeneralConfig,
	backend *backends.SimulatedBackend,
	flagsAndDeps ...interface{},
) *TestApplication {
	chainId := backend.Blockchain().Config().ChainID
	cfg.Overrides.DefaultChainID = chainId

	// Only set P2PEnabled override to false if it wasn't set by calling test
	if !cfg.Overrides.P2PEnabled.Valid {
		cfg.Overrides.P2PEnabled = null.BoolFrom(false)
	}

	client := client.NewSimulatedBackendClient(t, backend, chainId)
	eventBroadcaster := pg.NewEventBroadcaster(cfg.DatabaseURL(), 0, 0, logger.TestLogger(t), uuid.NewV4())

	zero := models.MustMakeDuration(0 * time.Millisecond)
	reaperThreshold := models.MustMakeDuration(100 * time.Millisecond)
	simulatedBackendChain := evmtypes.DBChain{
		ID: *utils.NewBigI(SimulatedBackendEVMChainID),
		Cfg: &evmtypes.ChainCfg{
			GasEstimatorMode:               null.StringFrom("FixedPrice"),
			EvmHeadTrackerMaxBufferSize:    null.IntFrom(100),
			EvmHeadTrackerSamplingInterval: &zero, // Head sampling disabled
			EthTxResendAfterThreshold:      &zero,
			EvmFinalityDepth:               null.IntFrom(15),
			EthTxReaperThreshold:           &reaperThreshold,
			MinIncomingConfirmations:       null.IntFrom(1),
			MinimumContractPayment:         assets.NewLinkFromJuels(100),
		},
		Enabled: true,
	}

	flagsAndDeps = append(flagsAndDeps, client, eventBroadcaster, simulatedBackendChain)

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
