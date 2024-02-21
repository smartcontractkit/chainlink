package liquiditymanager

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

var _ Rebalancer = &EvmRebalancer{}

type EvmRebalancer struct {
	rebalancer rebalancer.RebalancerInterface
	addr       common.Address
	networkSel models.NetworkSelector
	lggr       logger.Logger
}

func NewEvmRebalancer(
	address models.Address,
	net models.NetworkSelector,
	ec client.Client,
	lggr logger.Logger,
) (*EvmRebalancer, error) {
	rebal, err := rebalancer.NewRebalancer(common.Address(address), ec)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate rebalancer wrapper: %w", err)
	}

	return &EvmRebalancer{
		rebalancer: rebal,
		addr:       common.Address(address),
		networkSel: net,
		lggr:       lggr.Named("EvmRebalancer"),
	}, nil
}

func (e *EvmRebalancer) GetRebalancers(ctx context.Context) (map[models.NetworkSelector]models.Address, error) {
	lms, err := e.rebalancer.GetAllCrossChainRebalancers(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("get all cross chain rebalancers: %w", err)
	}
	ret := make(map[models.NetworkSelector]models.Address)
	for _, lm := range lms {
		ret[models.NetworkSelector(lm.RemoteChainSelector)] = models.Address(lm.RemoteRebalancer)
	}
	return ret, nil
}

func (e *EvmRebalancer) GetBalance(ctx context.Context) (*big.Int, error) {
	return e.rebalancer.GetLiquidity(&bind.CallOpts{Context: ctx})
}

func (e *EvmRebalancer) Close(ctx context.Context) error {
	return nil
}

// ConfigDigest implements Rebalancer.
func (e *EvmRebalancer) ConfigDigest(ctx context.Context) (types.ConfigDigest, error) {
	cdae, err := e.rebalancer.LatestConfigDigestAndEpoch(&bind.CallOpts{Context: ctx})
	if err != nil {
		return ocrtypes.ConfigDigest{}, fmt.Errorf("latest config digest and epoch: %w", err)
	}
	return ocrtypes.ConfigDigest(cdae.ConfigDigest), nil
}

func (e *EvmRebalancer) GetTokenAddress(ctx context.Context) (models.Address, error) {
	tokenAddress, err := e.rebalancer.ILocalToken(&bind.CallOpts{
		Context: ctx,
	})
	return models.Address(tokenAddress), err
}

func (e *EvmRebalancer) GetLatestSequenceNumber(ctx context.Context) (uint64, error) {
	cdae, err := e.rebalancer.LatestConfigDigestAndEpoch(&bind.CallOpts{Context: ctx})
	return cdae.SequenceNumber, err
}
