package v1_2

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type DeployRouterInput struct {
	IsTestRouter bool
	Chain        uint64
}

var (
	DeployRouter = operations.NewOperation(
		"DeployRouter",
		semver.MustParse("1.0.0"),
		"Deploys Router 1.2 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.OpDependencies, input DeployRouterInput) (common.Address, error) {
			state := deps.CurrentState
			e := deps.Env
			ab := deps.AddressBook
			chain := deps.Env.Chains[input.Chain]
			chainState, chainExists := state.Chains[input.Chain]
			if !chainExists {
				return common.Address{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", chain.String())
			}

			rmnProxy := chainState.RMNProxy
			if chainState.RMNProxy == nil {
				e.Logger.Errorw("RMNProxy not found", "chain", chain.String())
				return common.Address{}, fmt.Errorf("rmn proxy not found for chain %s, deploy the prerequisites first", chain.String())
			}
			deployFn := func(chain cldf.Chain, tv cldf.TypeAndVersion) (cldf.ContractDeploy[*router.Router], error) {
				r, err := cldf.DeployContract(e.Logger, chain, ab,
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
					return cldf.ContractDeploy[*router.Router]{}, err
				}
				return cldf.ContractDeploy[*router.Router]{
					Address:  r.Address,
					Contract: r.Contract,
					Tx:       r.Tx,
					Tv:       r.Tv,
					Err:      nil,
				}, nil
			}
			if input.IsTestRouter {
				if chainState.TestRouter != nil {
					e.Logger.Infow("test router already deployed", "chain", chain.String(), "addr", chainState.TestRouter.Address)
					return common.Address{}, nil
				}
				r, err := deployFn(chain, cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0))
				if err != nil {
					return common.Address{}, err
				}
				e.Logger.Infow("deployed test router", "chain", chain.String(), "addr", chainState.TestRouter.Address)
				return r.Address, nil
			}
			if chainState.Router != nil {
				e.Logger.Infow("router already deployed, no-op", "chain", chain.String(), "addr", chainState.Router.Address)
				return common.Address{}, nil
			}
			r, err := deployFn(chain, cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0))
			if err != nil {
				return common.Address{}, err
			}
			e.Logger.Infow("deployed router", "chain", chain.String(), "addr", chainState.Router.Address)
			return r.Address, err
		})
)
