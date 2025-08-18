package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/ethereum/go-ethereum/core/types"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	mcmsnew_zksync "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/evm/zksync"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMDeployCallProxyInput struct {
	Timelock common.Address `json:"timelock"`
}

var OpEVMDeployCallProxy = opsutils.NewEVMDeployOperation(
	"evm-call-proxy-deploy",
	semver.MustParse("1.0.0"),
	"Deploys CallProxy contract on the specified EVM chains",
	cldf.NewTypeAndVersion(commontypes.CallProxy, deployment.Version1_0_0),
	opsutils.VMDeployers[OpEVMDeployCallProxyInput]{
		DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput OpEVMDeployCallProxyInput) (common.Address, *types.Transaction, error) {
			addr, tx, _, err := bindings.DeployCallProxy(
				opts,
				backend,
				deployInput.Timelock,
			)
			return addr, tx, err
		},
		DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, deployInput OpEVMDeployCallProxyInput) (common.Address, error) {
			addr, _, _, err := mcmsnew_zksync.DeployCallProxyZk(
				opts,
				client,
				wallet,
				backend,
				deployInput.Timelock,
			)
			return addr, err
		},
	},
)
