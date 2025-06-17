package v1_2

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type DeployRouterInput struct {
	RMNProxy      common.Address
	WethAddress   common.Address
	ChainSelector uint64
}

type RouterApplyRampUpdatesOpInput struct {
	OnRampUpdates  []router.RouterOnRamp
	OffRampRemoves []router.RouterOffRamp
	OffRampAdds    []router.RouterOffRamp
}

var (
	routerDeployers = opsutil.VMDeployers[DeployRouterInput]{
		DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, input DeployRouterInput) (common.Address, *types.Transaction, error) {
			addr, tx, _, err := router.DeployRouter(
				opts,
				backend,
				input.WethAddress,
				input.RMNProxy,
			)
			return addr, tx, err
		},
		DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, input DeployRouterInput) (common.Address, error) {
			addr, _, _, err := router.DeployRouterZk(
				opts,
				client,
				wallet,
				backend,
				input.WethAddress,
				input.RMNProxy,
			)
			return addr, err
		},
	}

	DeployRouter = opsutil.NewEVMDeployOperation(
		"DeployRouter",
		semver.MustParse("1.0.0"),
		"Deploys Router 1.2 contract on the specified evm chain",
		cldf.NewTypeAndVersion(shared.Router, deployment.Version1_2_0),
		routerDeployers,
	)

	DeployTestRouter = opsutil.NewEVMDeployOperation(
		"DeployTestRouter",
		semver.MustParse("1.0.0"),
		"Deploys TestRouter 1.2 contract on the specified evm chain",
		cldf.NewTypeAndVersion(shared.TestRouter, deployment.Version1_2_0),
		routerDeployers,
	)

	RouterApplyRampUpdatesOp = opsutil.NewEVMCallOperation(
		"RouterApplyRampUpdatesOp",
		semver.MustParse("1.0.0"),
		"Updates OnRamps and OffRamps on the Router contract",
		router.RouterABI,
		shared.Router,
		router.NewRouter,
		func(router *router.Router, opts *bind.TransactOpts, input RouterApplyRampUpdatesOpInput) (*types.Transaction, error) {
			return router.ApplyRampUpdates(opts, input.OnRampUpdates, input.OffRampRemoves, input.OffRampAdds)
		},
	)
)
