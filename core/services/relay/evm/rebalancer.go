package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/discoverer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager"
	rebalancermodels "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/ocr3impls"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	ocr3ABI = evmtypes.MustGetABI(no_op_ocr3.NoOpOCR3MetaData.ABI)
)

type RebalancerProvider interface {
	commontypes.Plugin
	ContractTransmitterOCR3() ocr3types.ContractTransmitter[rebalancermodels.Report]
	LiquidityManagerFactory() liquiditymanager.Factory
	DiscovererFactory() discoverer.Factory
	BridgeFactory() bridge.Factory
}

type RebalancerRelayer interface {
	NewRebalancerProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (RebalancerProvider, error)
}

var _ RebalancerRelayer = (*rebalancerRelayer)(nil)

type rebalancerRelayer struct {
	chains      legacyevm.LegacyChainContainer
	lggr        logger.Logger
	ethKeystore keystore.Eth
}

func NewRebalancerRelayer(
	chains legacyevm.LegacyChainContainer,
	lggr logger.Logger,
	ethKeystore keystore.Eth) RebalancerRelayer {
	return &rebalancerRelayer{
		chains:      chains,
		lggr:        lggr,
		ethKeystore: ethKeystore,
	}
}

// NewRebalancerProvider implements RebalancerRelayer.
func (r *rebalancerRelayer) NewRebalancerProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (RebalancerProvider, error) {
	configWatcher, lmContracts, lmFactory, discovererFactory, bridgeFactory, err := newRebalancerConfigProvider(r.lggr, r.chains, rargs)
	if err != nil {
		return nil, fmt.Errorf("failed to create config watcher: %w", err)
	}

	var (
		transmitters = make(map[relay.ID]ocr3types.ContractTransmitter[rebalancermodels.Report])
	)
	for _, chain := range r.chains.Slice() {
		fromAddresses, err2 := r.ethKeystore.EnabledAddressesForChain(chain.ID())
		if err2 != nil {
			return nil, fmt.Errorf("failed to get enabled keys for chain %s: %w", chain.ID().String(), err2)
		}
		if len(fromAddresses) != 1 {
			return nil, fmt.Errorf("rebalancer services: expected only one enabled key for chain %s, got %d", chain.ID().String(), len(fromAddresses))
		}
		relayID := relay.NewID(relay.EVM, chain.ID().String())
		tm, err2 := ocrcommon.NewTransmitter(
			chain.TxManager(),
			fromAddresses,
			1e6, // TODO: gas limit may vary depending on tx
			fromAddresses[0],
			txmgr.NewSendEveryStrategy(),
			txmgrtypes.TransmitCheckerSpec[common.Address]{},
			chain.ID(),
			r.ethKeystore,
		)
		if err2 != nil {
			return nil, fmt.Errorf("failed to create transmitter: %w", err2)
		}
		t, err2 := ocr3impls.NewOCR3ContractTransmitter[rebalancermodels.Report](
			lmContracts[relayID],
			ocr3ABI,
			tm,
			r.lggr.Named(fmt.Sprintf("OCR3ContractTransmitter-%s", chain.ID().String())),
			nil, // TODO: implement report to evm tx metadata
		)
		if err2 != nil {
			return nil, fmt.Errorf("failed to create ocr3 contract transmitter: %w", err2)
		}
		transmitters[relayID] = t
	}
	multichainTransmitter, err := ocr3impls.NewMultichainTransmitterOCR3[rebalancermodels.Report](
		transmitters,
		r.lggr.Named("MultichainTransmitterOCR3"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create multichain transmitter: %w", err)
	}
	return &rebalancerProvider{
		configWatcher:       configWatcher,
		contractTransmitter: multichainTransmitter,
		lmFactory:           lmFactory,
		discovererFactory:   discovererFactory,
		bridgeFactory:       bridgeFactory,
	}, nil
}

var _ RebalancerProvider = (*rebalancerProvider)(nil)

type rebalancerProvider struct {
	*configWatcher
	contractTransmitter ocr3types.ContractTransmitter[rebalancermodels.Report]
	lmFactory           liquiditymanager.Factory
	discovererFactory   discoverer.Factory
	bridgeFactory       bridge.Factory
}

func (r *rebalancerProvider) Codec() commontypes.Codec {
	return nil
}

// ChainReader implements RebalancerProvider.
func (*rebalancerProvider) ChainReader() commontypes.ChainReader {
	return nil
}

// ContractTransmitter implements RebalancerProvider.
func (*rebalancerProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return nil
}

func (r *rebalancerProvider) ContractTransmitterOCR3() ocr3types.ContractTransmitter[rebalancermodels.Report] {
	return r.contractTransmitter
}

func (r *rebalancerProvider) LiquidityManagerFactory() liquiditymanager.Factory {
	return r.lmFactory
}

func (r *rebalancerProvider) DiscovererFactory() discoverer.Factory {
	return r.discovererFactory
}

func (r *rebalancerProvider) BridgeFactory() bridge.Factory {
	return r.bridgeFactory
}

func newRebalancerConfigProvider(
	lggr logger.Logger,
	chains legacyevm.LegacyChainContainer,
	rargs commontypes.RelayArgs,
) (
	*configWatcher,
	map[relay.ID]common.Address,
	liquiditymanager.Factory,
	discoverer.Factory,
	bridge.Factory,
	error,
) {
	var relayConfig types.RelayConfig
	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to unmarshal relay config (%s): %w", string(rargs.RelayConfig), err)
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, nil, nil, nil, nil, fmt.Errorf("invalid contract address %s", rargs.ContractID)
	}

	var lmFactoryOpts []liquiditymanager.Opt
	var discovererOpts []discoverer.Opt
	for _, chain := range chains.Slice() {
		ch, exists := chainsel.ChainByEvmChainID(chain.ID().Uint64())
		if !exists {
			return nil, nil, nil, nil, nil, fmt.Errorf("chain %d not found", chain.ID().Uint64())
		}

		lmFactoryOpts = append(lmFactoryOpts, liquiditymanager.WithEvmDep(
			rebalancermodels.NetworkSelector(ch.Selector),
			chain.Client(),
		))
		discovererOpts = append(discovererOpts, discoverer.WithEvmDep(
			rebalancermodels.NetworkSelector(ch.Selector),
			chain.Client(),
		))
	}
	lmFactory := liquiditymanager.NewBaseRebalancerFactory(lggr, lmFactoryOpts...)
	discovererFactory := discoverer.NewFactory(lggr, discovererOpts...)

	masterChain, err := chains.Get(relayConfig.ChainID.String())
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to get master chain %s: %w", relayConfig.ChainID, err)
	}

	logPollers := make(map[relay.ID]logpoller.LogPoller)
	for _, chain := range chains.Slice() {
		logPollers[relay.NewID(relay.EVM, chain.ID().String())] = chain.LogPoller()
	}

	// sanity check that all chains specified in RelayConfig.fromBlocks are present
	for chainID := range relayConfig.FromBlocks {
		_, err2 := chains.Get(chainID)
		if err2 != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to get chain %s specified in RelayConfig.fromBlocks: %w", chainID, err2)
		}
	}

	contractAddress := common.HexToAddress(rargs.ContractID)
	masterSelector, ok := chainsel.ChainByEvmChainID(masterChain.ID().Uint64())
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("chain selector for master chain %d not found", masterChain.ID().Uint64())
	}

	discoverer, err := discovererFactory.NewDiscoverer(
		rebalancermodels.NetworkSelector(masterSelector.Selector),
		rebalancermodels.Address(contractAddress),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to create initial discoverer for master chain %d (evm chain id: %d): %w",
			masterSelector.Selector, masterSelector.EvmChainID, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	rebalancerGraph, err := discoverer.Discover(ctx)
	cancel()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to discover rebalancer graph: %w", err)
	}

	// at this point we can instantiate the bridges
	var bridgeOpts []bridge.Opt
	for _, networkID := range rebalancerGraph.GetNetworks() {
		chain, ok := chainsel.ChainBySelector(uint64(networkID))
		if !ok {
			return nil, nil, nil, nil, nil, fmt.Errorf("chain selector for network %d not found", networkID)
		}
		legacyChain, err2 := chains.Get(strconv.FormatUint(chain.EvmChainID, 10))
		if err2 != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to get legacy chain chain %d: %w", chain.EvmChainID, err2)
		}
		rebalancerAddress, err2 := rebalancerGraph.GetRebalancerAddress(networkID)
		if err2 != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to get rebalancer address for network %d: %w", networkID, err2)
		}
		xchainRebalData, err2 := rebalancerGraph.GetXChainRebalancerData(networkID)
		if err2 != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("failed to get xchain rebalancer data for network %d: %w", networkID, err2)
		}
		bridgeAdapters := make(map[rebalancermodels.NetworkSelector]rebalancermodels.Address)
		for remoteNetworkID, data := range xchainRebalData {
			bridgeAdapters[remoteNetworkID] = data.LocalBridgeAdapterAddress
		}
		bridgeOpts = append(bridgeOpts, bridge.WithEvmDep(
			networkID,
			legacyChain.LogPoller(),
			legacyChain.Client(),
			rebalancerAddress,
			bridgeAdapters,
		))
	}
	bridgeFactory := bridge.NewFactory(lggr, bridgeOpts...)

	mcct, err := ocr3impls.NewMultichainConfigTracker(
		relay.NewID(relay.EVM, relayConfig.ChainID.String()),
		lggr.Named("MultichainConfigTracker"),
		logPollers,
		masterChain.Client(),
		contractAddress,
		discovererFactory,
		ocr3impls.TransmitterCombiner,
		relayConfig.FromBlocks,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to create multichain config tracker: %w", err)
	}

	digester := ocr3impls.MultichainConfigDigester{
		MasterChainDigester: evmutil.EVMOffchainConfigDigester{
			ChainID:         masterChain.ID().Uint64(),
			ContractAddress: contractAddress,
		},
	}

	return newConfigWatcher(
		lggr.Named("RebalancerConfigWatcher"),
		contractAddress,
		ocr3ABI,
		digester,
		mcct,
		masterChain,
		relayConfig.FromBlock,
		rargs.New,
	), mcct.GetContractAddresses(), lmFactory, discovererFactory, bridgeFactory, nil
}
