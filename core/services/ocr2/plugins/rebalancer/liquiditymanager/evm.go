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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditygraph"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager/liquidity_manager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type EvmLiquidityManager struct {
	client      OnchainLiquidityManager
	lp          logpoller.LogPoller
	lmAbi       abi.ABI
	addr        common.Address
	networkSel  models.NetworkSelector
	ec          client.Client
	cleanupFunc func() error
}

func NewEvmLiquidityManager(address models.Address, net models.NetworkSelector, ec client.Client, lp logpoller.LogPoller) (*EvmLiquidityManager, error) {
	// uncomment when we implement the concrete lm
	// client, err := NewConcreteLiquidityManager(common.Address(address), ec)
	// if err != nil {
	// 	return nil, fmt.Errorf("new concrete lm: %w", err)
	// }
	dummyClient, err := NewDummyLiquidityManager(common.Address(address), ec, big.NewInt(1000))
	if err != nil {
		return nil, fmt.Errorf("new dummy lm: %w", err)
	}

	lmAbi, err := abi.JSON(strings.NewReader(liquidity_manager.LiquidityManagerABI))
	if err != nil {
		return nil, fmt.Errorf("new lm abi: %w", err)
	}

	lpFilter := logpoller.Filter{
		Name: fmt.Sprintf("lm-liquidity-transferred-%s", common.Address(address)),
		EventSigs: []common.Hash{
			lmAbi.Events["LiquidityTransferred"].ID,
		},
		Addresses: []common.Address{common.Address(address)},
	}

	if err := lp.RegisterFilter(lpFilter); err != nil {
		return nil, fmt.Errorf("register filter: %w", err)
	}

	return &EvmLiquidityManager{
		client:     dummyClient,
		lp:         lp,
		lmAbi:      lmAbi,
		ec:         ec,
		addr:       common.Address(address),
		networkSel: net,
		cleanupFunc: func() error {
			return lp.UnregisterFilter(lpFilter.Name)
		},
	}, nil
}

func (e *EvmLiquidityManager) GetLiquidityManagers(ctx context.Context) (map[models.NetworkSelector]models.Address, error) {
	return e.client.GetAllCrossChainLiquidityMangers(ctx)
}

func (e *EvmLiquidityManager) GetBalance(ctx context.Context) (*big.Int, error) {
	return e.client.GetLiquidity(ctx)
}

func (e *EvmLiquidityManager) GetPendingTransfers(ctx context.Context, since time.Time) ([]models.PendingTransfer, error) {
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
		))
		// tr.Status = models.TransferStatusExecuted // todo: determine the status
		pendingTransfers = append(pendingTransfers, tr)
	}

	return pendingTransfers, nil
}
func (e EvmLiquidityManager) Discover(ctx context.Context, lmFactory Factory) (*Registry, liquiditygraph.LiquidityGraph, error) {
	g := liquiditygraph.NewGraph()
	lms := NewRegistry()

	type qItem struct {
		networkID models.NetworkSelector
		lmAddress common.Address
	}

	seen := mapset.NewSet[qItem]()
	queue := mapset.NewSet[qItem]()

	elem := qItem{networkID: e.networkSel, lmAddress: e.addr}
	queue.Add(elem)
	seen.Add(elem)

	for queue.Cardinality() > 0 {
		elem, ok := queue.Pop()
		if !ok {
			return nil, nil, fmt.Errorf("unexpected internal error, there is a bug in the algorithm")
		}

		// TODO: investigate fetching the balance here.
		g.AddNetwork(elem.networkID, big.NewInt(0))

		lm, err := lmFactory.NewLiquidityManager(elem.networkID, models.Address(elem.lmAddress))
		if err != nil {
			return nil, nil, fmt.Errorf("init liquidity manager: %w", err)
		}

		lms.Add(elem.networkID, models.Address(elem.lmAddress))

		destinationLMs, err := lm.GetLiquidityManagers(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("get %v destination liquidity managers: %w", elem.networkID, err)
		}

		if destinationLMs == nil {
			continue
		}

		for destNetworkID, lmAddr := range destinationLMs {
			g.AddConnection(elem.networkID, destNetworkID)

			newElem := qItem{networkID: destNetworkID, lmAddress: common.Address(lmAddr)}
			if !seen.Contains(newElem) {
				queue.Add(newElem)
				seen.Add(newElem)

				if _, exists := lms.Get(destNetworkID); !exists {
					lms.Add(destNetworkID, lmAddr)
				}
			}
		}
	}

	return lms, g, nil
}

func (e *EvmLiquidityManager) Close(ctx context.Context) error {
	return e.cleanupFunc()
}
