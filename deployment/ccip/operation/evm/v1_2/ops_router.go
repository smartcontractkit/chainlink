package v1_2

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"

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
	DeployRouter = opsutil.NewEVMDeployOperation(
		"DeployRouter",
		semver.MustParse("1.0.0"),
		"Deploys Router 1.2 contract on the specified evm chain",
		shared.Router,
		router.RouterMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_2_0,
			EVMBytecode:      common.FromHex(router.RouterBin),
			ZkSyncVMBytecode: router.RouterZkBytecode,
		},
		func(input DeployRouterInput) []any {
			return []any{input.WethAddress, input.RMNProxy}
		},
	)

	DeployTestRouter = opsutil.NewEVMDeployOperation(
		"DeployTestRouter",
		semver.MustParse("1.0.0"),
		"Deploys TestRouter 1.2 contract on the specified evm chain",
		shared.TestRouter,
		router.RouterMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_2_0,
			EVMBytecode:      common.FromHex(router.RouterBin),
			ZkSyncVMBytecode: router.RouterZkBytecode,
		},
		func(input DeployRouterInput) []any {
			return []any{input.WethAddress, input.RMNProxy}
		},
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

	UpdateWrappedNativeAddressOnRouterOp = opsutil.NewEVMCallOperation(
		"UpdateWrappedNativeAddressOnRouterOp",
		semver.MustParse("1.0.0"),
		"Updates Wrapped Native token address on Router contract for a chain",
		router.RouterABI,
		shared.Router,
		router.NewRouter,
		func(router *router.Router, opts *bind.TransactOpts, wrappedNativeAddress common.Address) (*types.Transaction, error) {
			return router.SetWrappedNative(opts, wrappedNativeAddress)
		},
	)
)
