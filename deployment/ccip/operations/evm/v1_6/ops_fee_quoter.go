package v1_6

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type DeployRouterInput struct {
	Chain        cldf.Chain
	IsTestRouter bool
}

var (
	DeployRouter = operations.NewOperation(
		"DeployRouter",
		semver.MustParse("1.0.0"),
		"Deploys Router contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.OpDependencies, input DeployRouterInput) (opsutil.OpOutput, error) {
			state := deps.CurrentState
			e := deps.Env
			ab := cldf.NewMemoryAddressBook()
			chainState, chainExists := state.Chains[input.Chain.Selector]
			if !chainExists {
				return opsutil.OpOutput{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", input.Chain.String())
			}
			chain := input.Chain
			rmnProxy := chainState.RMNProxy
			if chainState.RMNProxy == nil {
				e.Logger.Errorw("RMNProxy not found", "chain", chain.String())
				return opsutil.OpOutput{}, fmt.Errorf("rmn proxy not found for chain %s, deploy the prerequisites first", chain.String())
			}
			deployFn := func(chain cldf.Chain, tv cldf.TypeAndVersion) error {
				_, err := cldf.DeployContract(e.Logger, chain, ab,
					func(chain cldf.Chain) cldf.ContractDeploy[*router.Router] {
						routerAddr, tx2, routerC, err2 := router.DeployRouter(
							chain.DeployerKey,
							chain.Client,
							chainState.Weth9.Address(),
							rmnProxy.Address(),
						)
						return cldf.ContractDeploy[*router.Router]{
							Address: routerAddr, Contract: routerC, Tx: tx2, Tv: tv, Err: err2,
						}
					})
				if err != nil {
					e.Logger.Errorw("Failed to deploy router", "chain", chain.String(), "err", err)
					return err
				}
				return nil
			}
			if input.IsTestRouter {
				if chainState.TestRouter != nil {
					e.Logger.Infow("test router already deployed", "chain", chain.String(), "addr", chainState.TestRouter.Address)
					return opsutil.OpOutput{
						AddressBook: ab,
					}, nil
				}
				err := deployFn(chain, cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0))
				if err != nil {
					return opsutil.OpOutput{}, err
				}
				e.Logger.Infow("deployed test router", "chain", chain.String(), "addr", chainState.TestRouter.Address)
				return opsutil.OpOutput{
					AddressBook: ab,
				}, nil
			}
			if chainState.Router != nil {
				e.Logger.Infow("router already deployed, no-op", "chain", chain.String(), "addr", chainState.Router.Address)
				return opsutil.OpOutput{
					AddressBook: ab,
				}, nil
			}
			err := deployFn(chain, cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0))
			if err != nil {
				return opsutil.OpOutput{}, err
			}
			e.Logger.Infow("deployed router", "chain", chain.String(), "addr", chainState.Router.Address)
			return opsutil.OpOutput{
				AddressBook: ab,
			}, err
		})
)

type FeeQuoterParams struct {
	MaxFeeJuelsPerMsg              *big.Int
	TokenPriceStalenessThreshold   uint32
	LinkPremiumMultiplierWeiPerEth uint64
	WethPremiumMultiplierWeiPerEth uint64
	MorePremiumMultiplierWeiPerEth []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs
	TokenPriceFeedUpdates          []fee_quoter.FeeQuoterTokenPriceFeedUpdate
	TokenTransferFeeConfigArgs     []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs
	DestChainConfigArgs            []fee_quoter.FeeQuoterDestChainConfigArgs
}

func (c FeeQuoterParams) Validate() error {
	if c.MaxFeeJuelsPerMsg == nil {
		return errors.New("MaxFeeJuelsPerMsg is nil")
	}
	if c.MaxFeeJuelsPerMsg.Cmp(big.NewInt(0)) <= 0 {
		return errors.New("MaxFeeJuelsPerMsg must be positive")
	}
	if c.TokenPriceStalenessThreshold == 0 {
		return errors.New("TokenPriceStalenessThreshold can't be 0")
	}
	return nil
}

func DefaultFeeQuoterParams() FeeQuoterParams {
	return FeeQuoterParams{
		MaxFeeJuelsPerMsg:              big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
		TokenPriceStalenessThreshold:   uint32(24 * 60 * 60),
		LinkPremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
		WethPremiumMultiplierWeiPerEth: 1e18, // 1.0 ETH
		TokenPriceFeedUpdates:          []fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
		TokenTransferFeeConfigArgs:     []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{},
		MorePremiumMultiplierWeiPerEth: []fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{},
		DestChainConfigArgs:            []fee_quoter.FeeQuoterDestChainConfigArgs{},
	}
}
