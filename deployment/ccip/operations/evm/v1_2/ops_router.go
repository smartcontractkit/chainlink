package v1_2

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type DeployRouterInput struct {
	IsTestRouter bool
	opsutil.DeployContractInput
}

var (
	DeployRouter = operations.NewOperation(
		"DeployRouter",
		semver.MustParse("1.0.0"),
		"Deploys Router contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.OpDependencies, input DeployRouterInput) (cldf.ContractDeploy[*router.Router], error) {
			state := deps.CurrentState
			e := deps.Env
			ab := input.AB
			chainState, chainExists := state.Chains[input.Chain.Selector]
			if !chainExists {
				return cldf.ContractDeploy[*router.Router]{}, fmt.Errorf("chain %s not found in existing state, "+
					"deploy the prerequisites first", input.Chain.String())
			}
			chain := input.Chain
			rmnProxy := chainState.RMNProxy
			if chainState.RMNProxy == nil {
				e.Logger.Errorw("RMNProxy not found", "chain", chain.String())
				return cldf.ContractDeploy[*router.Router]{}, fmt.Errorf("rmn proxy not found for chain %s, deploy the prerequisites first", chain.String())
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
					return cldf.ContractDeploy[*router.Router]{}, nil
				}
				r, err := deployFn(chain, cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0))
				if err != nil {
					return cldf.ContractDeploy[*router.Router]{}, err
				}
				e.Logger.Infow("deployed test router", "chain", chain.String(), "addr", chainState.TestRouter.Address)
				return r, nil
			}
			if chainState.Router != nil {
				e.Logger.Infow("router already deployed, no-op", "chain", chain.String(), "addr", chainState.Router.Address)
				return cldf.ContractDeploy[*router.Router]{}, nil
			}
			r, err := deployFn(chain, cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0))
			if err != nil {
				return cldf.ContractDeploy[*router.Router]{}, err
			}
			e.Logger.Infow("deployed router", "chain", chain.String(), "addr", chainState.Router.Address)
			return r, err
		})
)
