package bridge

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/arb"
	bridgecommon "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/opstack"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/testonlybridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

var (
	ErrBridgeNotFound = errors.New("bridge not found")
)

// Bridge provides a way to get pending transfers from one chain to another
// for transfers that are using the native bridge for the (source, dest) chain pair.
// For example, if ethereum is the source, and arbitrum is the destination, the bridge
// would be able to get pending transfers from ethereum to arbitrum via the standard arbitrum
// bridge.
type Bridge interface {
	// GetTransfers returns all of the pending transfers from the source chain to the destination chain
	// for the given local and remote token addresses.
	// Pending transfers that are ready to finalize have the appropriate bridge data set.
	GetTransfers(ctx context.Context, localToken, remoteToken models.Address) ([]models.PendingTransfer, error)

	// GetBridgePayloadAndFee returns the bridge specific payload for the given transfer.
	// This payload must always be correctly ABI-encoded.
	// Note that this payload is not directly provided to the bridge but the bridge adapter
	// contracts. The bridge adapter may slightly alter the payload before sending it to the bridge.
	// For example, for an L1 to L2 transfer using Arbitrum's bridge, this will return the
	// fees required for the transfer to succeed reliably.
	// This should only be called when we want to trigger a transfer (i.e, there is no transfer in flight)
	// Bridge specific payloads for pending transfers are returned by GetTransfers.
	GetBridgePayloadAndFee(ctx context.Context, transfer models.Transfer) ([]byte, *big.Int, error)

	// QuorumizedBridgePayload returns a single bridge payload given the slice of bridge payloads.
	// The liquidityManager oracles all generate bridge payloads separately, and this function is used to
	// "collapse" all of them into a single payload in a pure way.
	// For example, if the bridge payload is a cost parameter, one implementation of this method
	// could either produce the median of all the costs, or take the maximum cost, or the minimum
	// cost. The choice of implementation is up to the bridge.
	QuorumizedBridgePayload(payloads [][]byte, f int) ([]byte, error)

	Close(ctx context.Context) error
}

type Factory interface {
	NewBridge(ctx context.Context, source, dest models.NetworkSelector) (Bridge, error)
	GetBridge(source, dest models.NetworkSelector) (Bridge, error)
}

type Opt func(c *factory)

type evmDep struct {
	lp                      logpoller.LogPoller
	ethClient               client.Client
	liquidityManagerAddress models.Address
	bridgeAdapters          map[models.NetworkSelector]models.Address
}

type factory struct {
	evmDeps       map[models.NetworkSelector]evmDep
	cachedBridges sync.Map
	lggr          logger.Logger
}

func NewFactory(lggr logger.Logger, opts ...Opt) Factory {
	c := &factory{
		evmDeps: make(map[models.NetworkSelector]evmDep),
		lggr:    lggr,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithEvmDep(
	networkID models.NetworkSelector,
	lp logpoller.LogPoller,
	ethClient client.Client,
	liquidityManagerAddress models.Address,
	bridgeAdapters map[models.NetworkSelector]models.Address,
) Opt {
	return func(f *factory) {
		f.evmDeps[networkID] = evmDep{
			lp:                      lp,
			ethClient:               ethClient,
			liquidityManagerAddress: liquidityManagerAddress,
			bridgeAdapters:          bridgeAdapters,
		}
	}
}

func (f *factory) NewBridge(ctx context.Context, source, dest models.NetworkSelector) (Bridge, error) {
	if source == dest {
		return nil, fmt.Errorf("no bridge between the same network and itself: %d", source)
	}

	bridge, err := f.GetBridge(source, dest)
	if errors.Is(err, ErrBridgeNotFound) {
		f.lggr.Infow("Bridge not found, initializing new bridge", "source", source, "dest", dest)
		return f.initBridge(ctx, source, dest)
	}
	return bridge, err
}

func (f *factory) initBridge(ctx context.Context, source, dest models.NetworkSelector) (Bridge, error) {
	f.lggr.Debugw("Initializing bridge", "source", source, "dest", dest)

	var bridge Bridge
	var err error

	switch source {
	// Arbitrum L2 --> Ethereum L1 bridge
	case models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
		models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector):
		if !bridgecommon.Supports(source, dest) {
			return nil, fmt.Errorf("unsupported destination for arbitrum l2 -> l1 bridge: %d", dest)
		}
		l2Deps, ok := f.evmDeps[source]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for source selector %d", source)
		}
		l1Deps, ok := f.evmDeps[dest]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for dest selector %d", dest)
		}
		l1BridgeAdapter, ok := l1Deps.bridgeAdapters[source]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for source selector %d in deps for dest selector %d", dest, source)
		}
		l2BridgeAdapter, ok := l2Deps.bridgeAdapters[dest]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for dest selector %d in deps for source selector %d", source, dest)
		}
		f.lggr.Infow("addresses check",
			"l1RollupAddress", arb.AllContracts[uint64(dest)]["Rollup"],
			"l1liquidityManagerAddress", l1Deps.liquidityManagerAddress,
			"l1BridgeAdapter", l1BridgeAdapter,
			"l2BridgeAdapter", l2BridgeAdapter,
		)
		bridge, err = arb.NewL2ToL1Bridge(
			ctx,
			f.lggr,
			source,
			dest,
			arb.AllContracts[uint64(dest)]["Rollup"], // l1 rollup address
			common.Address(l1Deps.liquidityManagerAddress), // l1 liquidityManager address
			common.Address(l2Deps.liquidityManagerAddress), // l2 liquidityManager address
			l2Deps.lp,        // l2 log poller
			l1Deps.lp,        // l1 log poller
			l2Deps.ethClient, // l2 eth client
			l1Deps.ethClient, // l1 eth client
		)

	// Optimism L2 --> Ethereum L1 bridge
	case models.NetworkSelector(chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector),
		models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector):
		if !bridgecommon.Supports(source, dest) {
			return nil, fmt.Errorf("unsupported destination for optimism l2 -> l1 bridge: %d", dest)
		}
		l2Deps, ok := f.evmDeps[source]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for source selector %d", source)
		}
		l1Deps, ok := f.evmDeps[dest]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for dest selector %d", dest)
		}
		l1BridgeAdapter, ok := l1Deps.bridgeAdapters[source]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for source selector %d in deps for dest selector %d", dest, source)
		}
		l2BridgeAdapter, ok := l2Deps.bridgeAdapters[dest]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for dest selector %d in deps for source selector %d", source, dest)
		}
		f.lggr.Infow("addresses check",
			"l1StandardBridgeProxyAddress", opstack.OptimismContractsByChainSelector[uint64(dest)]["L1StandardBridgeProxy"],
			"l2StandardBridgeAddress", opstack.OptimismContractsByChainSelector[uint64(source)]["L2StandardBridge"],
			"l1liquidityManagerAddress", l1Deps.liquidityManagerAddress,
			"l2liquidityManagerAddress", l2Deps.liquidityManagerAddress,
			"l1BridgeAdapter", l1BridgeAdapter,
			"l2BridgeAdapter", l2BridgeAdapter,
		)
		bridge, err = opstack.NewL2ToL1Bridge(
			ctx,
			f.lggr,
			source,
			dest,
			common.Address(l1Deps.liquidityManagerAddress), // l1 liquidityManager address
			common.Address(l2Deps.liquidityManagerAddress), // l2 liquidityManager address
			l1Deps.ethClient, // l1 eth client
			l2Deps.ethClient, // l2 eth client
			l1Deps.lp,        // l1 log poller
			l2Deps.lp,        // l2 log poller
		)
	// Ethereum L1 --> Arbitrum L2 bridge OR
	// Ethereum L1 --> Optimism L2 bridge
	case models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
		models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector):
		if !bridgecommon.Supports(source, dest) {
			return nil, fmt.Errorf("unsupported destination for eth l1 -> l2 bridge: %d", dest)
		}
		l1Deps, ok := f.evmDeps[source]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for source selector %d", source)
		}
		l2Deps, ok := f.evmDeps[dest]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for dest selector %d", dest)
		}
		l1BridgeAdapter, ok := l1Deps.bridgeAdapters[dest]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for source selector %d in deps for selector %d", source, dest)
		}
		f.lggr.Infow("addresses check",
			"l1GatewayRouterAddress", arb.AllContracts[uint64(source)]["L1GatewayRouter"],
			"inboxAddress", arb.AllContracts[uint64(source)]["L1Inbox"],
			"l1liquidityManagerAddress", l1Deps.liquidityManagerAddress,
			"l2liquidityManagerAddress", l2Deps.liquidityManagerAddress,
			"l1BridgeAdapter", l1BridgeAdapter,
		)
		switch dest {
		case models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
			models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector):
			f.lggr.Infow("dest arb addresses check",
				"l1GatewayRouterAddress", arb.AllContracts[uint64(source)]["L1GatewayRouter"],
				"inboxAddress", arb.AllContracts[uint64(source)]["L1Inbox"],
				"l1liquidityManagerAddress", l1Deps.liquidityManagerAddress,
				"l2liquidityManagerAddress", l2Deps.liquidityManagerAddress,
				"l1BridgeAdapter", l1BridgeAdapter,
			)
			bridge, err = arb.NewL1ToL2Bridge(
				ctx,
				f.lggr,
				source,
				dest,
				common.Address(l1Deps.liquidityManagerAddress),      // l1 liquidityManager address
				common.Address(l2Deps.liquidityManagerAddress),      // l2 liquidityManager address
				arb.AllContracts[uint64(source)]["L1GatewayRouter"], // l1 gateway router address
				arb.AllContracts[uint64(source)]["L1Inbox"],         // l1 inbox address
				l1Deps.ethClient, // l1 eth client
				l2Deps.ethClient, // l2 eth client
				l1Deps.lp,        // l1 log poller
				l2Deps.lp,        // l2 log poller
			)
		case models.NetworkSelector(chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector),
			models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.Selector):
			f.lggr.Infow("dest OP addresses check",
				"L1StandardBridgeProxyAddress", opstack.OptimismContractsByChainSelector[uint64(source)]["L1StandardBridgeProxy"],
				"L2StandardBridgeAddress", opstack.OptimismContractsByChainSelector[uint64(dest)]["L2StandardBridge"],
				"l1liquidityManagerAddress", l1Deps.liquidityManagerAddress,
				"l2liquidityManagerAddress", l2Deps.liquidityManagerAddress,
				"l1BridgeAdapter", l1BridgeAdapter,
			)
			bridge, err = opstack.NewL1ToL2Bridge(
				ctx,
				f.lggr,
				source,
				dest,
				common.Address(l1Deps.liquidityManagerAddress),                                    // l1 liquidityManager address
				common.Address(l2Deps.liquidityManagerAddress),                                    // l2 liquidityManager address
				opstack.OptimismContractsByChainSelector[uint64(source)]["L1StandardBridgeProxy"], // l1 standard bridge proxy address
				opstack.OptimismContractsByChainSelector[uint64(dest)]["L2StandardBridge"],        // l2 standard bridge address
				l1Deps.ethClient, // l1 eth client
				l2Deps.ethClient, // l2 eth client
				l1Deps.lp,        // l1 log poller
				l2Deps.lp,        // l2 log poller
			)
		default:
			return nil, fmt.Errorf("unsupported destination for eth l1 -> l2 bridge: %d", dest)
		}
	case models.NetworkSelector(chainsel.GETH_TESTNET.Selector),
		models.NetworkSelector(chainsel.TEST_90000001.Selector),
		models.NetworkSelector(chainsel.TEST_90000002.Selector),
		models.NetworkSelector(chainsel.TEST_90000003.Selector),
		models.NetworkSelector(chainsel.GETH_DEVNET_2.Selector):
		// these chains are only ever used for tests
		// in tests we only ever deploy the MockL1Bridge adapter
		// so this is an "L1 to L1" bridge setup, but not really
		if source == dest {
			return nil, fmt.Errorf("no bridge between the same network and itself: %d", source)
		}
		sourceDeps, ok := f.evmDeps[source]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for source selector %d", source)
		}
		destDeps, ok := f.evmDeps[dest]
		if !ok {
			return nil, fmt.Errorf("evm dependencies not found for dest selector %d", dest)
		}
		sourceAdapter, ok := sourceDeps.bridgeAdapters[dest]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for source selector %d in deps for selector %d", source, dest)
		}
		destAdapter, ok := destDeps.bridgeAdapters[source]
		if !ok {
			return nil, fmt.Errorf("bridge adapter not found for dest selector %d in deps for selector %d", dest, source)
		}
		bridge, err = testonlybridge.New(
			ctx,
			source,
			dest,
			sourceDeps.liquidityManagerAddress,
			destDeps.liquidityManagerAddress,
			sourceAdapter,
			destAdapter,
			sourceDeps.lp,
			destDeps.lp,
			sourceDeps.ethClient,
			destDeps.ethClient,
			f.lggr,
		)
	default:
		return nil, fmt.Errorf("unsupported source chain selector: %d", source)
	}

	if err != nil {
		return nil, err
	}

	f.cachedBridges.Store(f.cacheKey(source, dest), bridge)
	return bridge, nil
}

func (f *factory) GetBridge(source, dest models.NetworkSelector) (Bridge, error) {
	bridge, exists := f.cachedBridges.Load(f.cacheKey(source, dest))
	if !exists {
		return nil, ErrBridgeNotFound
	}

	b, ok := bridge.(Bridge)
	if !ok {
		return nil, fmt.Errorf("cached bridge has wrong type: %T", bridge)
	}
	return b, nil
}

func (f *factory) cacheKey(source, dest models.NetworkSelector) string {
	return fmt.Sprintf("%d-%d", source, dest)
}
