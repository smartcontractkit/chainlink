package liquiditymanager

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var _ Rebalancer = &EvmRebalancer{}

type EvmRebalancer struct {
	client      OnchainRebalancer
	lp          logpoller.LogPoller
	lmAbi       abi.ABI
	addr        common.Address
	networkSel  models.NetworkSelector
	ec          client.Client
	cleanupFunc func() error
	lggr        logger.Logger
}

func NewEvmRebalancer(
	address models.Address,
	net models.NetworkSelector,
	ec client.Client,
	lp logpoller.LogPoller,
	lggr logger.Logger) (*EvmRebalancer, error) {
	client, err := NewConcreteRebalancer(common.Address(address), ec)
	if err != nil {
		return nil, fmt.Errorf("new concrete rebalancer: %w", err)
	}

	lmAbi, err := abi.JSON(strings.NewReader(rebalancer.RebalancerABI))
	if err != nil {
		return nil, fmt.Errorf("new rebalancer abi: %w", err)
	}

	lpFilter := logpoller.Filter{
		Name: fmt.Sprintf("%d-lm-liquidity-transferred-%s", time.Now().UnixNano(), common.Address(address)),
		EventSigs: []common.Hash{
			lmAbi.Events["LiquidityTransferred"].ID,
		},
		Addresses: []common.Address{common.Address(address)},
	}

	if err := lp.RegisterFilter(lpFilter); err != nil {
		return nil, fmt.Errorf("register filter: %w", err)
	}

	return &EvmRebalancer{
		client:     client,
		lp:         lp,
		lmAbi:      lmAbi,
		ec:         ec,
		addr:       common.Address(address),
		networkSel: net,
		cleanupFunc: func() error {
			return lp.UnregisterFilter(lpFilter.Name)
		},
		lggr: lggr.Named("EvmRebalancer"),
	}, nil
}

func (e *EvmRebalancer) GetRebalancers(ctx context.Context) (map[models.NetworkSelector]models.Address, error) {
	return e.client.GetAllCrossChainRebalancers(ctx)
}

func (e *EvmRebalancer) GetBalance(ctx context.Context) (*big.Int, error) {
	return e.client.GetLiquidity(ctx)
}

func (e *EvmRebalancer) GetPendingTransfers(ctx context.Context, since time.Time) ([]models.PendingTransfer, error) {
	logs, err := e.lp.LogsCreatedAfter(
		e.lmAbi.Events["LiquidityTransferred"].ID,
		e.addr,
		since,
		logpoller.Finalized,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get logs created after: %w", err)
	}

	pendingTransfers := make([]models.PendingTransfer, 0, len(logs))

	for _, log := range logs {
		liqTransferred, err2 := e.client.ParseLiquidityTransferred(log.ToGethLog())
		if err2 != nil {
			return nil, fmt.Errorf("invalid log: %w", err2)
		}

		tr := models.NewPendingTransfer(models.NewTransfer(
			models.NetworkSelector(liqTransferred.FromChainSelector()),
			models.NetworkSelector(liqTransferred.ToChainSelector()),
			liqTransferred.Amount(),
			log.BlockTimestamp,
			[]byte{}, // TODO: fill in bridge data
		))
		// tr.Status = models.TransferStatusExecuted // todo: determine the status
		pendingTransfers = append(pendingTransfers, tr)
	}

	return pendingTransfers, nil
}

// bfsItem is an item in the BFS queue.
type bfsItem struct {
	networkSel models.NetworkSelector
	lmAddress  common.Address
}

func (q bfsItem) String() string {
	return fmt.Sprintf("(NetworkSelector:%v,LMAddress:%v)", q.networkSel, q.lmAddress)
}

func (e EvmRebalancer) Discover(ctx context.Context, lmFactory Factory) (*Registry, liquiditygraph.LiquidityGraph, error) {
	lggr := e.lggr.With("func", "Discover", "rebalancer", e.addr)
	lggr.Debugw("Starting discovery")

	g := liquiditygraph.NewGraph()
	lms := NewRegistry()

	seen := mapset.NewSet[bfsItem]()
	queue := mapset.NewSet[bfsItem]()

	elem := bfsItem{networkSel: e.networkSel, lmAddress: e.addr}
	queue.Add(elem)
	seen.Add(elem)

	lggr.Debugw("Starting BFS", "queue", queue, "seen", seen)
	for queue.Cardinality() > 0 {
		elem, ok := queue.Pop()
		if !ok {
			return nil, nil, fmt.Errorf("unexpected internal error, there is a bug in the algorithm")
		}

		lggr.Debugw("Popped element from queue", "elem", elem)
		// TODO: investigate fetching the balance here.
		// TODO: make use of returned value?
		g.AddNetwork(elem.networkSel, big.NewInt(0))
		lggr.Debugw("Added elem to network", "elem", elem)

		lggr.Debugw("Creating new rebalancer object", "elem", elem)
		lm, err := lmFactory.NewRebalancer(elem.networkSel, models.Address(elem.lmAddress))
		if err != nil {
			lggr.Errorw("Failed to create new rebalancer", "err", err)
			return nil, nil, fmt.Errorf("init liquidity manager: %w", err)
		}

		lggr.Debugw("Adding rebalancer to registry", "elem", elem)
		lms.Add(elem.networkSel, models.Address(elem.lmAddress))

		lggr.Debugw("Getting destination liquidity managers", "elem", elem)
		destinationLMs, err := lm.GetRebalancers(ctx)
		if err != nil {
			lggr.Errorw("Failed to get destination liquidity managers", "err", err)
			return nil, nil, fmt.Errorf("get %v destination liquidity managers: %w", elem.networkSel, err)
		}

		lggr.Debugw("Got destination liquidity managers", "destinationLMs", destinationLMs, "elem", elem)
		if destinationLMs == nil {
			lggr.Debugw("No destination liquidity managers found", "destinationLMs", destinationLMs, "elem", elem)
			continue
		}

		lggr.Debugw("Adding connections", "elem", elem, "destinationLMs", destinationLMs)
		for destNetworkSel, lmAddr := range destinationLMs {
			if !g.HasNetwork(destNetworkSel) {
				lggr.Debugw("Adding new network to graph not yet seen", "destNetworkSel", destNetworkSel)
				g.AddNetwork(destNetworkSel, big.NewInt(0))
			} else {
				lggr.Debugw("Network already seen, not adding", "destNetworkSel", destNetworkSel)
			}
			lggr.Debugw("Adding connection", "elem", elem, "destNetworkSel", destNetworkSel, "sourceNetworkID", elem.networkSel)
			if err := g.AddConnection(elem.networkSel, destNetworkSel); err != nil {
				return nil, nil, fmt.Errorf("add connection: %w", err)
			}

			newElem := bfsItem{networkSel: destNetworkSel, lmAddress: common.Address(lmAddr)}
			lggr.Debugw("Deciding whether to add new element to queue", "newElem", newElem)
			if !seen.Contains(newElem) {
				lggr.Debugw("Not seen before, adding to queue", "newElem", newElem)
				queue.Add(newElem)
				seen.Add(newElem)

				if _, exists := lms.Get(destNetworkSel); !exists {
					lggr.Debugw("Not seen before, adding to registry", "newElem", newElem)
					lms.Add(destNetworkSel, lmAddr)
				} else {
					lggr.Debugw("Already seen, not adding to registry", "newElem", newElem)
				}
			} else {
				lggr.Debugw("Already seen, not adding to queue", "newElem", newElem)
			}
		}
	}

	return lms, g, nil
}

func (e *EvmRebalancer) Close(ctx context.Context) error {
	return e.cleanupFunc()
}

// ConfigDigest implements Rebalancer.
func (e *EvmRebalancer) ConfigDigest(ctx context.Context) (types.ConfigDigest, error) {
	return e.client.GetConfigDigest(ctx)
}
