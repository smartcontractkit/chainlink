package ocr3impls

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"go.uber.org/multierr"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/discoverer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

var (
	// See https://github.com/smartcontractkit/ccip/compare/ccip-develop...CCIP-1438-op-stack-bridge-adapter-l-1#diff-2fe14bb9d1ecbc62f43cef26daff5d1f86275f16e1296fc9827b934a518d3f4cR20
	ConfigSet common.Hash

	defaultABI abi.ABI

	_ ocrtypes.ContractConfigTracker = &multichainConfigTracker{}

	defaultTimeout = 1 * time.Minute

	configTrackerWorkers = 4
)

func init() {
	var err error
	tabi, err := no_op_ocr3.NoOpOCR3MetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	defaultABI = *tabi
	ConfigSet = defaultABI.Events["ConfigSet"].ID
}

type CombinerFn func(masterChain commontypes.RelayID, contractConfigs map[commontypes.RelayID]ocrtypes.ContractConfig) (ocrtypes.ContractConfig, error)

type multichainConfigTracker struct {
	services.StateMachine

	// masterChain is the chain that contains the "master" OCR3 configuration
	// contract.
	masterChain       commontypes.RelayID
	lggr              logger.Logger
	logPollers        map[commontypes.RelayID]logpoller.LogPoller
	masterClient      evmclient.Client
	contractAddresses map[commontypes.RelayID]common.Address
	masterContract    no_op_ocr3.NoOpOCR3Interface
	combiner          CombinerFn
	fromBlocks        map[string]int64
	replaying         atomic.Bool
}

func NewMultichainConfigTracker(
	masterChain commontypes.RelayID,
	lggr logger.Logger,
	logPollers map[commontypes.RelayID]logpoller.LogPoller,
	masterClient evmclient.Client,
	masterContract common.Address,
	discovererFactory discoverer.Factory,
	combiner CombinerFn,
	fromBlocks map[string]int64,
) (*multichainConfigTracker, error) {
	// Ensure master chain is in the log pollers
	if _, ok := logPollers[masterChain]; !ok {
		return nil, fmt.Errorf("master chain %s not in log pollers", masterChain)
	}

	// Ensure combiner is not nil
	if combiner == nil {
		return nil, fmt.Errorf("provide non-nil combiner")
	}

	// Ensure factory is not nil
	if discovererFactory == nil {
		return nil, fmt.Errorf("provide non-nil liquidity manager factory")
	}

	// before we register filters we need to discover all the available liquidity managers
	// on all the chains
	masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network ID %s: %w", masterChain, err)
	}

	chain, exists := chainsel.ChainByEvmChainID(masterChainID)
	if !exists {
		return nil, fmt.Errorf("chain selector for chain %d not found", masterChainID)
	}

	discoverer, err := discovererFactory.NewDiscoverer(
		models.NetworkSelector(chain.Selector),
		models.Address(masterContract),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create discoverer: %w", err)
	}

	// Discover all the liquidity managers
	lggr.Infow("Discovering all liquidity managers", "masterLM", masterContract.Hex())
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	grph, err := discoverer.Discover(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to discover liquidity managers: %w", err)
	}
	lggr.Infow("Finished discovering all liquidity managers", "graph", grph)

	// sanity check, there should be only one liquidity manager per-chain per-asset
	if grph.Len() != len(logPollers) {
		return nil, fmt.Errorf("expected %d liquidity managers but found %d", len(logPollers), grph.Len())
	}

	// Register filters on all log pollers
	contracts := make(map[commontypes.RelayID]common.Address)
	for id, lp := range logPollers {
		nid, err2 := strconv.ParseUint(id.ChainID, 10, 64)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse network ID %s: %w", id, err2)
		}

		ch, exists := chainsel.ChainByEvmChainID(nid)
		if !exists {
			return nil, fmt.Errorf("chain %d not found", nid)
		}

		address, err2 := grph.GetLiquidityManagerAddress(models.NetworkSelector(ch.Selector))
		if err2 != nil {
			return nil, fmt.Errorf("no rebalancer found for network selector %d", ch.Selector)
		}
		fName := configTrackerFilterName(id, common.Address(address))
		err2 = lp.RegisterFilter(ctx, logpoller.Filter{
			Name:      fName,
			EventSigs: []common.Hash{ConfigSet},
			Addresses: []common.Address{common.Address(address)},
		})
		if err2 != nil {
			return nil, err2
		}
		contracts[id] = common.Address(address)
	}

	wrapper, err := no_op_ocr3.NewNoOpOCR3(masterContract, masterClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create wrapper: %w", err)
	}

	return &multichainConfigTracker{
		lggr:              lggr,
		logPollers:        logPollers,
		masterChain:       masterChain,
		combiner:          combiner,
		masterClient:      masterClient,
		contractAddresses: contracts,
		masterContract:    wrapper,
		fromBlocks:        fromBlocks,
	}, nil
}

func (m *multichainConfigTracker) GetContractAddresses() map[commontypes.RelayID]common.Address {
	return m.contractAddresses
}

func (m *multichainConfigTracker) Start() {
	_ = m.StartOnce("MultichainConfigTracker", func() error {
		if len(m.fromBlocks) == 0 {
			return nil
		}
		m.replaying.Store(true)
		defer m.replaying.Store(false)

		networks := len(m.fromBlocks)
		ctx := context.Background() // TODO: deadline?
		running := make(chan struct{}, configTrackerWorkers)
		results := make(chan error, networks)
		go func() {
			for id, fromBlock := range m.fromBlocks {
				running <- struct{}{}
				go func(id string, fromBlock int64) {
					defer func() { <-running }()
					err := m.ReplayChain(ctx, commontypes.NewRelayID("evm", id), fromBlock)
					if err != nil {
						m.lggr.Errorw("failed to replay chain", "chain", id, "fromBlock", fromBlock, "err", err)
					} else {
						m.lggr.Infow("successfully replayed chain", "chain", id, "fromBlock", fromBlock)
					}
					results <- err
				}(id, fromBlock)
			}
		}()

		// wait for results, we expect the same number of results as networks
		var errs error
		for i := 0; i < networks; i++ {
			err := <-results
			if err != nil {
				errs = multierr.Append(errs, err)
			}
		}
		if errs != nil {
			m.lggr.Errorw("failed to replay some chains", "err", errs)
		}
		return errs
	})
}

func (m *multichainConfigTracker) Close() error {
	return nil
}

// Notify noop method
func (m *multichainConfigTracker) Notify() <-chan struct{} {
	return nil
}

// ReplayChain replays the log poller for the provided chain
func (m *multichainConfigTracker) ReplayChain(ctx context.Context, id commontypes.RelayID, fromBlock int64) error {
	if _, ok := m.logPollers[id]; !ok {
		return fmt.Errorf("chain %s not in log pollers", id)
	}
	return m.logPollers[id].Replay(ctx, fromBlock)
}

// Replay replays the log poller for the master chain
func (m *multichainConfigTracker) Replay(ctx context.Context, fromBlock int64) error {
	return m.logPollers[m.masterChain].Replay(ctx, fromBlock)
}

// LatestBlockHeight implements types.ContractConfigTracker.
// Returns the block height of the master chain.
func (m *multichainConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latestBlock, err := m.logPollers[m.masterChain].LatestBlock(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return uint64(latestBlock.BlockNumber), nil
}

// LatestConfig implements types.ContractConfigTracker.
// LatestConfig fetches the config from the master chain and then fetches the
// remaining configurations from all the other chains.
func (m *multichainConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	// if we're still replaying the follower chains we won't have their configs
	if m.replaying.Load() {
		return ocrtypes.ContractConfig{}, errors.New("cannot call LatestConfig while replaying")
	}

	lgs, err := m.logPollers[m.masterChain].Logs(ctx, int64(changedInBlock), int64(changedInBlock), ConfigSet, m.contractAddresses[m.masterChain])
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		return ocrtypes.ContractConfig{}, fmt.Errorf("no logs found for config on contract %s (chain %s) at block %d",
			m.contractAddresses[m.masterChain].Hex(), m.masterChain.String(), changedInBlock)
	}
	masterConfig, err := configFromLog(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	m.lggr.Infow("LatestConfig from master chain", "latestConfig", masterConfig)

	// check all other chains for their config
	contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
		m.masterChain: masterConfig,
	}
	for id, lp := range m.logPollers {
		if id == m.masterChain {
			continue
		}

		lg, err2 := lp.LatestLogByEventSigWithConfs(ctx, ConfigSet, m.contractAddresses[id], 1)
		if err2 != nil {
			return ocrtypes.ContractConfig{}, err2
		}

		configSet, err2 := configFromLog(lg.Data)
		if err2 != nil {
			return ocrtypes.ContractConfig{}, err2
		}
		contractConfigs[id] = configSet
	}

	// at this point we can combine the configs into a single one
	combined, err := m.combiner(m.masterChain, contractConfigs)
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("error combining configs: %w", err)
	}

	return combined, nil
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (m *multichainConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latest, err := m.logPollers[m.masterChain].LatestLogByEventSigWithConfs(ctx, ConfigSet, m.contractAddresses[m.masterChain], 1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.callLatestConfigDetails(ctx)
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	masterConfig, err := configFromLog(latest.Data)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, fmt.Errorf("failed to unpack latest config details: %w", err)
	}

	return uint64(latest.BlockNumber), masterConfig.ConfigDigest, nil
}

func (m *multichainConfigTracker) callLatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	lcd, err := m.masterContract.LatestConfigDetails(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, fmt.Errorf("failed to get latest config details: %w", err)
	}
	return uint64(lcd.BlockNumber), lcd.ConfigDigest, nil
}
