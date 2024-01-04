package ocr3

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var (
	// See https://github.com/smartcontractkit/ccip/compare/ccip-develop...CCIP-1438-op-stack-bridge-adapter-l-1#diff-2fe14bb9d1ecbc62f43cef26daff5d1f86275f16e1296fc9827b934a518d3f4cR20
	ConfigSet common.Hash

	defaultABI abi.ABI

	_ ocrtypes.ContractConfigTracker = &multichainConfigTracker{}
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

type CombinerFn func(masterConfig ocrtypes.ContractConfig, followerConfigs []ocrtypes.ContractConfig) (ocrtypes.ContractConfig, error)

type multichainConfigTracker struct {
	services.StateMachine

	// masterChain is the chain that contains the "master" OCR3 configuration
	// contract. This is the chain that the config tracker will listen to for
	// ConfigSet events. All other chains will have their config set in a similar
	// way however will not contain offchain or onchain config, just signers and
	// transmitters. These events will be read by the multichain config tracker
	// and used to construct the final multi-chain config.
	masterChain       relay.ID
	lggr              logger.Logger
	logPollers        map[relay.ID]logpoller.LogPoller
	clients           map[relay.ID]evmclient.Client
	contractAddresses map[relay.ID]common.Address
	contracts         map[relay.ID]no_op_ocr3.NoOpOCR3Interface
	combiner          CombinerFn
}

func NewMultichainConfigTracker(
	masterChain relay.ID,
	lggr logger.Logger,
	logPollers map[relay.ID]logpoller.LogPoller,
	clients map[relay.ID]evmclient.Client,
	contractAddresses map[relay.ID]common.Address,
	combiner CombinerFn,
) (*multichainConfigTracker, error) {
	// Ensure master chain is in the log pollers
	if _, ok := logPollers[masterChain]; !ok {
		return nil, fmt.Errorf("master chain %s not in log pollers", masterChain)
	}

	// Ensure master chain is in the clients
	if _, ok := clients[masterChain]; !ok {
		return nil, fmt.Errorf("master chain %s not in clients", masterChain)
	}

	// Ensure combiner is not nil
	if combiner == nil {
		return nil, fmt.Errorf("provide non-nil combiner")
	}

	// Register filters on all log pollers
	contracts := make(map[relay.ID]no_op_ocr3.NoOpOCR3Interface)
	for id, lp := range logPollers {
		fName := configTrackerFilterName(id, contractAddresses[id])
		err := lp.RegisterFilter(logpoller.Filter{
			Name:      fName,
			EventSigs: []common.Hash{ConfigSet},
			Addresses: []common.Address{contractAddresses[id]},
		})
		if err != nil {
			return nil, err
		}
		wrapper, err := no_op_ocr3.NewNoOpOCR3(contractAddresses[id], clients[id])
		if err != nil {
			return nil, err
		}
		contracts[id] = wrapper
	}
	return &multichainConfigTracker{
		lggr:              lggr,
		logPollers:        logPollers,
		clients:           clients,
		contractAddresses: contractAddresses,
		contracts:         contracts,
		masterChain:       masterChain,
	}, nil
}

func (m *multichainConfigTracker) Start() {}

func (m *multichainConfigTracker) Close() error {
	return nil
}

// Notify noop method
func (m *multichainConfigTracker) Notify() <-chan struct{} {
	return nil
}

// Replay abstracts the logpoller.LogPoller Replay() implementation
func (m *multichainConfigTracker) Replay(ctx context.Context, id relay.ID, fromBlock int64) error {
	return m.logPollers[id].Replay(ctx, fromBlock)
}

// LatestBlockHeight implements types.ContractConfigTracker.
// Returns the block height of the master chain.
func (m *multichainConfigTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	latestBlock, err := m.logPollers[m.masterChain].LatestBlock(pg.WithParentCtx(ctx))
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
	lgs, err := m.logPollers[m.masterChain].Logs(int64(changedInBlock), int64(changedInBlock), ConfigSet, m.contractAddresses[m.masterChain], pg.WithParentCtx(ctx))
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	if len(lgs) == 0 {
		// TODO: fallback to RPC call?
		return ocrtypes.ContractConfig{}, fmt.Errorf("no logs found for config on contract %s (chain %s) at block %d", m.contractAddresses[m.masterChain].Hex(), m.clients[m.masterChain].ConfiguredChainID().String(), changedInBlock)
	}
	masterConfig, err := configFromLog(lgs[len(lgs)-1].Data)
	if err != nil {
		return ocrtypes.ContractConfig{}, err
	}
	m.lggr.Infow("LatestConfig from master chain", "latestConfig", masterConfig)

	// check all other chains for their config
	var followerConfigs []ocrtypes.ContractConfig
	for id, lp := range m.logPollers {
		if id == m.masterChain {
			continue
		}

		lgs, err := lp.Logs(int64(changedInBlock), int64(changedInBlock), ConfigSet, m.contractAddresses[id], pg.WithParentCtx(ctx))
		if err != nil {
			return ocrtypes.ContractConfig{}, err
		}

		if len(lgs) == 0 {
			return ocrtypes.ContractConfig{}, fmt.Errorf("no logs found for config on contract %s (chain %s) at block %d", m.contractAddresses[id].Hex(), m.clients[id].ConfiguredChainID().String(), changedInBlock)
		}

		configSet, err := configFromLog(lgs[len(lgs)-1].Data)
		if err != nil {
			return ocrtypes.ContractConfig{}, err
		}
		followerConfigs = append(followerConfigs, configSet)
	}

	// at this point we can combine the configs into a single one
	combined, err := m.combiner(masterConfig, followerConfigs)
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("error combining configs: %w", err)
	}

	return combined, nil
}

// LatestConfigDetails implements types.ContractConfigTracker.
func (m *multichainConfigTracker) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ocrtypes.ConfigDigest, err error) {
	latest, err := m.logPollers[m.masterChain].LatestLogByEventSigWithConfs(ConfigSet, m.contractAddresses[m.masterChain], 1, pg.WithParentCtx(ctx))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO: try eth-call to get the config
			return 0, ocrtypes.ConfigDigest{}, err
		}
		return 0, ocrtypes.ConfigDigest{}, err
	}
	masterConfig, err := configFromLog(latest.Data)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}

	// check all other chains for their config
	var followerConfigs []ocrtypes.ContractConfig
	for id, lp := range m.logPollers {
		if id == m.masterChain {
			continue
		}

		latest, err := lp.LatestLogByEventSigWithConfs(ConfigSet, m.contractAddresses[id], 1, pg.WithParentCtx(ctx))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// TODO: try RPC call for config
				return 0, ocrtypes.ConfigDigest{}, err
			}
			return 0, ocrtypes.ConfigDigest{}, err
		}

		followerConfig, err := configFromLog(latest.Data)
		if err != nil {
			return 0, ocrtypes.ConfigDigest{}, err
		}

		followerConfigs = append(followerConfigs, followerConfig)
	}

	// at this point we can combine the configs into a single one
	combined, err := m.combiner(masterConfig, followerConfigs)
	if err != nil {
		return 0, ocrtypes.ConfigDigest{}, err
	}

	return uint64(latest.BlockNumber), masterConfig.ConfigDigest, nil
}
