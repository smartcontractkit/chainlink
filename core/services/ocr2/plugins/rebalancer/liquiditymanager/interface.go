package liquiditymanager

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type OnchainRebalancer interface {
	GetAllCrossChainRebalancers(ctx context.Context) (map[models.NetworkSelector]models.Address, error)
	GetLiquidity(ctx context.Context) (*big.Int, error)
	ParseLiquidityTransferred(log gethtypes.Log) (LiquidityTransferredEvent, error)
	GetConfigDigest(ctx context.Context) (ocrtypes.ConfigDigest, error)
}

type LiquidityTransferredEvent interface {
	FromChainSelector() uint64
	ToChainSelector() uint64
	Amount() *big.Int
}

var _ OnchainRebalancer = &concreteRebalancer{}

// concreteRebalancer implements OnchainRebalancer
// using the actual rebalancer contract's generated go bindings.
// i.e, business model is in full effect here.
type concreteRebalancer struct {
	client rebalancer.RebalancerInterface
}

// GetConfigDigest implements OnchainRebalancer.
func (c *concreteRebalancer) GetConfigDigest(ctx context.Context) (ocrtypes.ConfigDigest, error) {
	cdae, err := c.client.LatestConfigDigestAndEpoch(&bind.CallOpts{Context: ctx})
	if err != nil {
		return ocrtypes.ConfigDigest{}, fmt.Errorf("latest config digest and epoch: %w", err)
	}
	return ocrtypes.ConfigDigest(cdae.ConfigDigest), nil
}

// GetAllCrossChainRebalancers implements OnchainRebalancer.
func (c *concreteRebalancer) GetAllCrossChainRebalancers(ctx context.Context) (map[models.NetworkSelector]models.Address, error) {
	lms, err := c.client.GetAllCrossChainRebalancers(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get all cross chain rebalancers: %w", err)
	}
	ret := make(map[models.NetworkSelector]models.Address)
	for _, lm := range lms {
		ret[models.NetworkSelector(lm.RemoteChainSelector)] = models.Address(lm.RemoteRebalancer)
	}
	return ret, nil
}

// GetLiquidity implements OnchainRebalancer.
func (c *concreteRebalancer) GetLiquidity(ctx context.Context) (*big.Int, error) {
	return c.client.GetLiquidity(&bind.CallOpts{Context: ctx})
}

// ParseLiquidityTransferred implements OnchainRebalancer.
func (c *concreteRebalancer) ParseLiquidityTransferred(log gethtypes.Log) (LiquidityTransferredEvent, error) {
	e, err := c.client.ParseLiquidityTransferred(log)
	if err != nil {
		return nil, fmt.Errorf("parse liquidity transferred: %w", err)
	}
	return &concreteLiquidityTransferredEvent{e: e}, nil
}

func NewConcreteRebalancer(address common.Address, backend bind.ContractBackend) (*concreteRebalancer, error) {
	client, err := rebalancer.NewRebalancer(address, backend)
	if err != nil {
		return nil, fmt.Errorf("init concrete rebalancer: %w", err)
	}
	return &concreteRebalancer{client: client}, nil
}

type concreteLiquidityTransferredEvent struct {
	e *rebalancer.RebalancerLiquidityTransferred
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
