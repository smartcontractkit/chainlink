package liquiditymanager

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/dummy_liquidity_manager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/liquiditymanager/liquidity_manager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type OnchainLiquidityManager interface {
	GetAllCrossChainLiquidityMangers(ctx context.Context) (map[models.NetworkSelector]models.Address, error)
	GetLiquidity(ctx context.Context) (*big.Int, error)
	ParseLiquidityTransferred(log gethtypes.Log) (LiquidityTransferredEvent, error)
}

type LiquidityTransferredEvent interface {
	FromChainSelector() uint64
	ToChainSelector() uint64
	Amount() *big.Int
}

var _ OnchainLiquidityManager = &concreteLiquidityManager{}

// concreteLiquidityManager implements OnchainLiquidityManager
// using the actual liquidity manager contract's generated go bindings.
// i.e, business model is in full effect here.
type concreteLiquidityManager struct {
	client liquidity_manager.LiquidityManagerInterface
}

// GetAllCrossChainLiquidityMangers implements OnchainLiquidityManager.
func (c *concreteLiquidityManager) GetAllCrossChainLiquidityMangers(ctx context.Context) (map[models.NetworkSelector]models.Address, error) {
	lms, err := c.client.GetAllCrossChainLiquidityMangers(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get all cross chain liquidity managers: %w", err)
	}
	ret := make(map[models.NetworkSelector]models.Address)
	for _, lm := range lms {
		ret[models.NetworkSelector(lm.RemoteChainSelector)] = models.Address(lm.RemoteLiquidityManager)
	}
	return ret, nil
}

// GetLiquidity implements OnchainLiquidityManager.
func (c *concreteLiquidityManager) GetLiquidity(ctx context.Context) (*big.Int, error) {
	return c.client.GetLiquidity(&bind.CallOpts{Context: ctx})
}

// ParseLiquidityTransferred implements OnchainLiquidityManager.
func (c *concreteLiquidityManager) ParseLiquidityTransferred(log gethtypes.Log) (LiquidityTransferredEvent, error) {
	e, err := c.client.ParseLiquidityTransferred(log)
	if err != nil {
		return nil, fmt.Errorf("parse liquidity transferred: %w", err)
	}
	return &concreteLiquidityTransferredEvent{e: e}, nil
}

func NewConcreteLiquidityManager(address common.Address, backend bind.ContractBackend) (*concreteLiquidityManager, error) {
	client, err := liquidity_manager.NewLiquidityManager(address, backend)
	if err != nil {
		return nil, fmt.Errorf("init concrete liquidity manager: %w", err)
	}
	return &concreteLiquidityManager{client: client}, nil
}

type concreteLiquidityTransferredEvent struct {
	e *liquidity_manager.LiquidityManagerLiquidityTransferred
}

var _ LiquidityTransferredEvent = &concreteLiquidityTransferredEvent{}

func (c *concreteLiquidityTransferredEvent) FromChainSelector() uint64 {
	return c.e.FromChainSelector
}

func (c *concreteLiquidityTransferredEvent) ToChainSelector() uint64 {
	return c.e.ToChainSelector
}

func (c *concreteLiquidityTransferredEvent) Amount() *big.Int {
	return c.e.Amount
}

// dummyLiquidityManager implements OnchainLiquidityManager
// using the dummy liquidity manager gethwrapper,
// which only manages neighbors and has no business logic.
type dummyLiquidityManager struct {
	client       dummy_liquidity_manager.DummyLiquidityManagerInterface
	maxLiquidity *big.Int
}

// GetAllCrossChainLiquidityMangers implements OnchainLiquidityManager.
func (d *dummyLiquidityManager) GetAllCrossChainLiquidityMangers(ctx context.Context) (map[models.NetworkSelector]models.Address, error) {
	lms, err := d.client.GetAllCrossChainLiquidityMangers(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get all cross chain liquidity managers: %w", err)
	}
	ret := make(map[models.NetworkSelector]models.Address)
	for _, lm := range lms {
		ret[models.NetworkSelector(lm.RemoteChainSelector)] = models.Address(lm.RemoteLiquidityManager)
	}
	return ret, nil
}

// GetLiquidity implements OnchainLiquidityManager.
func (d *dummyLiquidityManager) GetLiquidity(ctx context.Context) (*big.Int, error) {
	liq, err := rand.Int(rand.Reader, d.maxLiquidity)
	if err != nil {
		return nil, fmt.Errorf("dummyLiquidityManager: generate random liquidity: %w", err)
	}
	return liq, nil
}

// ParseLiquidityTransferred implements OnchainLiquidityManager.
func (d *dummyLiquidityManager) ParseLiquidityTransferred(log gethtypes.Log) (LiquidityTransferredEvent, error) {
	e, err := d.client.ParseLiquidityTransferred(log)
	if err != nil {
		return nil, fmt.Errorf("parse liquidity transferred: %w", err)
	}
	return &dummyLiquidityTransferredEvent{e: e}, nil
}

var _ OnchainLiquidityManager = &dummyLiquidityManager{}

func NewDummyLiquidityManager(address common.Address, backend bind.ContractBackend, maxLiquidity *big.Int) (*dummyLiquidityManager, error) {
	client, err := dummy_liquidity_manager.NewDummyLiquidityManager(address, backend)
	if err != nil {
		return nil, fmt.Errorf("init dummy liquidity manager: %w", err)
	}
	return &dummyLiquidityManager{client: client, maxLiquidity: maxLiquidity}, nil
}

type dummyLiquidityTransferredEvent struct {
	e *dummy_liquidity_manager.DummyLiquidityManagerLiquidityTransferred
}

var _ LiquidityTransferredEvent = &dummyLiquidityTransferredEvent{}

func (d *dummyLiquidityTransferredEvent) FromChainSelector() uint64 {
	return d.e.FromChainSelector
}

func (d *dummyLiquidityTransferredEvent) ToChainSelector() uint64 {
	return d.e.ToChainSelector
}

func (d *dummyLiquidityTransferredEvent) Amount() *big.Int {
	return d.e.Amount
}
