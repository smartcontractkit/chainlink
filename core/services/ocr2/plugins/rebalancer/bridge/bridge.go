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
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge/arb"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge/testonlybridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
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
	// The rebalancer oracles all generate bridge payloads separately, and this function is used to
	// "collapse" all of them into a single payload in a pure way.
	// For example, if the bridge payload is a cost parameter, one implementation of this method
	// could either produce the median of all the costs, or take the maximum cost, or the minimum
	// cost. The choice of implementation is up to the bridge.
	QuorumizedBridgePayload(payloads [][]byte, f int) ([]byte, error)

	Close(ctx context.Context) error
}

//go:generate mockery --name Factory --output ./mocks --filename bridge_factory_mock.go --case=underscore
type Factory interface {
	NewBridge(source, dest models.NetworkSelector) (Bridge, error)
}

type Opt func(c *factory)

type evmDep struct {
	lp                logpoller.LogPoller
	ethClient         client.Client
	rebalancerAddress models.Address
	bridgeAdapters    map[models.NetworkSelector]models.Address
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
	rebalancerAddress models.Address,
	bridgeAdapters map[models.NetworkSelector]models.Address,
) Opt {
	return func(f *factory) {
		f.evmDeps[networkID] = evmDep{
			lp:                lp,
			ethClient:         ethClient,
			rebalancerAddress: rebalancerAddress,
			bridgeAdapters:    bridgeAdapters,
		}
	}
}

func (f *factory) NewBridge(source, dest models.NetworkSelector) (Bridge, error) {
	if source == dest {
		return nil, fmt.Errorf("no bridge between the same network and itself: %d", source)
	}

	bridge, err := f.GetBridge(source, dest)
	if errors.Is(err, ErrBridgeNotFound) {
		f.lggr.Infow("Bridge not found, initializing new bridge", "source", source, "dest", dest)
		return f.initBridge(source, dest)
	}
	return bridge, err
}

func (f *factory) initBridge(source, dest models.NetworkSelector) (Bridge, error) {
	f.lggr.Debugw("Initializing bridge", "source", source, "dest", dest)

	var bridge Bridge
	var err error

	switch source {
	case models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector),
		models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector):
		// source: arbitrum l2 -> dest: ethereum l1
		// only dest that is supported is eth mainnet if source == arb mainnet
		// only dest that is supported is eth sepolia if source == arb sepolia
		if source == models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector) &&
			dest != models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector) {
			return nil, fmt.Errorf("unsupported destination for arbitrum mainnet l1 -> l2 bridge: %d, must be eth mainnet", dest)
		}
		if source == models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector) &&
			dest != models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector) {
			return nil, fmt.Errorf("unsupported destination for arbitrum sepolia l1 -> l2 bridge: %d, must be eth sepolia", dest)
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
			"l1RebalancerAddress", l1Deps.rebalancerAddress,
			"l1BridgeAdapter", l1BridgeAdapter,
			"l2BridgeAdapter", l2BridgeAdapter,
		)
		bridge, err = arb.NewL2ToL1Bridge(
			f.lggr,
			source,
			dest,
			arb.AllContracts[uint64(dest)]["Rollup"], // l1 rollup address
			common.Address(l1Deps.rebalancerAddress), // l1 rebalancer address
			common.Address(l2Deps.rebalancerAddress), // l2 rebalancer address
			l2Deps.lp,                                // l2 log poller
			l1Deps.lp,                                // l1 log poller
			l2Deps.ethClient,                         // l2 eth client
			l1Deps.ethClient,                         // l1 eth client
		)
	case models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector),
		models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector):
		// source: Ethereum L1 -> dest: Arbitrum L2
		// only dest that is supported is arbitrum mainnet if source == eth mainnet
		// only dest that is supported is arbitrum sepolia if source == eth sepolia
		if source == models.NetworkSelector(chainsel.ETHEREUM_MAINNET.Selector) &&
			dest != models.NetworkSelector(chainsel.ETHEREUM_MAINNET_ARBITRUM_1.Selector) {
			return nil, fmt.Errorf("unsupported destination for eth mainnet l1 -> l2 bridge: %d, must be arb mainnet", dest)
		}
		if source == models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector) &&
			dest != models.NetworkSelector(chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector) {
			return nil, fmt.Errorf("unsupported destination for eth sepolia l1 -> l2 bridge: %d, must be arb sepolia", dest)
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
			"l1RebalancerAddress", l1Deps.rebalancerAddress,
			"l2RebalancerAddress", l2Deps.rebalancerAddress,
			"l1BridgeAdapter", l1BridgeAdapter,
		)
		bridge, err = arb.NewL1ToL2Bridge(
			f.lggr,
			source,
			dest,
			common.Address(l1Deps.rebalancerAddress), // l1 rebalancer address
			common.Address(l2Deps.rebalancerAddress), // l2 rebalancer address
			arb.AllContracts[uint64(source)]["L1GatewayRouter"], // l1 gateway router address
			arb.AllContracts[uint64(source)]["L1Inbox"],         // l1 inbox address
			l1Deps.ethClient, // l1 eth client
			l2Deps.ethClient, // l2 eth client
			l1Deps.lp,        // l1 log poller
			l2Deps.lp,        // l2 log poller
		)
	case models.NetworkSelector(chainsel.GETH_TESTNET.Selector),
		models.NetworkSelector(chainsel.TEST_90000001.Selector),
		models.NetworkSelector(chainsel.TEST_90000002.Selector),
		models.NetworkSelector(chainsel.TEST_90000003.Selector):
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
			source,
			dest,
			sourceDeps.rebalancerAddress,
			destDeps.rebalancerAddress,
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
