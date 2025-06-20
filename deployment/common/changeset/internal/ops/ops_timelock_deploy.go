package ops

import (
	"math/big"

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

type OpEVMDeployTimelockInput struct {
	TimelockMinDelay *big.Int         `json:"timelockMinDelay"`
	Admin            common.Address   `json:"admin"`      // Admin of the timelock contract, usually the deployer key
	Proposers        []common.Address `json:"proposers"`  // Proposer of the timelock contract, usually the deployer key
	Executors        []common.Address `json:"executors"`  // Executor of the timelock contract, usually the call proxy
	Cancellers       []common.Address `json:"cancellers"` // Canceller of the timelock contract, usually the deployer key
	Bypassers        []common.Address `json:"bypassers"`  // Bypasser of the timelock contract, usually the deployer key
}

var OpEVMDeployTimelock = opsutils.NewEVMDeployOperation(
	"evm-timelock-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Timelock contract on the specified EVM chains",
	cldf.NewTypeAndVersion(commontypes.RBACTimelock, deployment.Version1_0_0),
	opsutils.VMDeployers[OpEVMDeployTimelockInput]{
		DeployEVM: func(opts *bind.TransactOpts, backend bind.ContractBackend, deployInput OpEVMDeployTimelockInput) (common.Address, *types.Transaction, error) {
			addr, tx, _, err := bindings.DeployRBACTimelock(
				opts,
				backend,
				deployInput.TimelockMinDelay,
				deployInput.Admin,
				deployInput.Proposers,
				deployInput.Executors,
				deployInput.Cancellers,
				deployInput.Bypassers,
			)
			return addr, tx, err
		},
		DeployZksyncVM: func(opts *accounts.TransactOpts, client *clients.Client, wallet *accounts.Wallet, backend bind.ContractBackend, deployInput OpEVMDeployTimelockInput) (common.Address, error) {
			addr, _, _, err := mcmsnew_zksync.DeployRBACTimelockZk(
				opts,
				client,
				wallet,
				backend,
				deployInput.TimelockMinDelay,
				deployInput.Admin,
				deployInput.Proposers,
				deployInput.Executors,
				deployInput.Cancellers,
				deployInput.Bypassers,
			)
			return addr, err
		},
	})
