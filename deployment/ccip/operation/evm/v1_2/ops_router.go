package v1_2

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
)

type DeployRouterInput struct {
	IsTestRouter bool
	RMNProxy     common.Address
	WethAddress  common.Address
}

var (
	DeployRouter = operations.NewOperation(
		"DeployRouter",
		semver.MustParse("1.0.0"),
		"Deploys Router 1.2 contract on the specified evm chain",
		func(b operations.Bundle, deps opsutil.DeployContractDependencies, input DeployRouterInput) (common.Address, error) {
			ab := deps.AddressBook
			chain := deps.Chain

			deployFn := func(chain cldf_evm.Chain, tv cldf.TypeAndVersion) (cldf.ContractDeploy[*router.Router], error) {
				r, err := cldf.DeployContract(b.Logger, chain, ab,
					func(chain cldf_evm.Chain) cldf.ContractDeploy[*router.Router] {
						var (
							routerAddr common.Address
							tx2        *types.Transaction
							routerC    *router.Router
							err2       error
						)
						if chain.IsZkSyncVM {
							routerAddr, _, routerC, err2 = router.DeployRouterZk(
								nil,
								chain.ClientZkSyncVM,
								chain.DeployerKeyZkSyncVM,
								chain.Client,
								input.WethAddress,
								input.RMNProxy,
							)
						} else {
							routerAddr, tx2, routerC, err2 = router.DeployRouter(
								chain.DeployerKey,
								chain.Client,
								input.WethAddress,
								input.RMNProxy,
							)
						}
						return cldf.ContractDeploy[*router.Router]{
							Address: routerAddr, Contract: routerC, Tx: tx2, Tv: tv, Err: err2,
						}
					})
				if err != nil {
					b.Logger.Errorw("Failed to deploy router", "chain", chain.String(), "err", err)
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
				r, err := deployFn(chain, cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0))
				if err != nil {
					return common.Address{}, err
				}
				b.Logger.Infow("deployed test router", "chain", chain.String(), "addr", r.Address)
				return r.Address, nil
			}

			r, err := deployFn(chain, cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0))
			if err != nil {
				return common.Address{}, err
			}
			b.Logger.Infow("deployed router", "chain", chain.String(), "addr", r.Address)
			return r.Address, err
		})
)
