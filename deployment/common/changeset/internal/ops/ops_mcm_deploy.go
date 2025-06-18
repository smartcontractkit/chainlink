package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/ethereum/go-ethereum/core/types"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmsnew_zksync "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/evm/zksync"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type OpEVMDeployMCMOutput struct {
	Address common.Address `json:"address"`
}

var mcmDeployers = opsutils.VMDeployers[any]{
	DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, _ any) (common.Address, *types.Transaction, error) {
		addr, tx, _, err := bindings.DeployManyChainMultiSig(
			opts,
			backend,
		)
		return addr, tx, err
	},
	DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, _ any) (common.Address, error) {
		addr, _, _, err := mcmsnew_zksync.DeployManyChainMultiSigZk(
			opts,
			client,
			wallet,
			backend,
		)
		return addr, err
	},
}

var OpEVMDeployProposerMCM = opsutils.NewEVMDeployOperation(
	"evm-proposer-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Proposer MCM contract",
	cldf.NewTypeAndVersion(commontypes.ProposerManyChainMultisig, deployment.Version1_0_0),
	mcmDeployers,
)

var OpEVMDeployBypasserMCM = opsutils.NewEVMDeployOperation(
	"evm-bypasser-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Bypasser MCM contract",
	cldf.NewTypeAndVersion(commontypes.BypasserManyChainMultisig, deployment.Version1_0_0),
	mcmDeployers,
)

var OpEVMDeployCancellerMCM = opsutils.NewEVMDeployOperation(
	"evm-canceller-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Canceller MCM contract",
	cldf.NewTypeAndVersion(commontypes.CancellerManyChainMultisig, deployment.Version1_0_0),
	mcmDeployers,
)
