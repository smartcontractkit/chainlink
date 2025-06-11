package ops

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/ethereum/go-ethereum/core/types"
	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	mcmsnew_zksync "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/evm/zksync"
)

type OpEVMDeployTimelockInput struct {
	ContractType     cldf.ContractType
	TimelockMinDelay *big.Int         `json:"timelockMinDelay"`
	Admin            common.Address   `json:"admin"`      // Admin of the timelock contract, usually the deployer key
	Proposers        []common.Address `json:"proposers"`  // Proposer of the timelock contract, usually the deployer key
	Executors        []common.Address `json:"executors"`  // Executor of the timelock contract, usually the call proxy
	Cancellers       []common.Address `json:"cancellers"` // Canceller of the timelock contract, usually the deployer key
	Bypassers        []common.Address `json:"bypassers"`  // Bypasser of the timelock contract, usually the deployer key
	ChainSelector    uint64           // Needed to distinguish different input for Operations API
}

type OpEVMDeployTimelockOutput struct {
	Address common.Address `json:"address"`
}

var OpEVMDeployTimelock = operations.NewOperation(
	"evm-timelock-deploy",
	semver.MustParse("1.0.0"),
	"Deploys Timelock contract on the specified EVM chains",
	func(b operations.Bundle, deps OpEVMMCMSDeps, input OpEVMDeployTimelockInput) (OpEVMDeployTimelockOutput, error) {
		timelock, err := cldf.DeployContract(b.Logger, deps.Chain, deps.AddrBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.RBACTimelock] {
				var (
					timelock common.Address
					tx2      *types.Transaction
					cc       *bindings.RBACTimelock
					err2     error
				)
				if chain.IsZkSyncVM {
					timelock, _, cc, err2 = mcmsnew_zksync.DeployRBACTimelockZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						input.TimelockMinDelay,
						input.Admin,
						input.Proposers,
						input.Executors,
						input.Cancellers,
						input.Bypassers,
					)
				} else {
					timelock, tx2, cc, err2 = bindings.DeployRBACTimelock(
						chain.DeployerKey,
						chain.Client,
						input.TimelockMinDelay,
						input.Admin,
						input.Proposers,
						input.Executors,
						input.Cancellers,
						input.Bypassers,
					)
				}

				tv := cldf.NewTypeAndVersion(input.ContractType, deployment.Version1_0_0)
				for _, option := range deps.Options {
					option(&tv)
				}

				return cldf.ContractDeploy[*bindings.RBACTimelock]{
					Address: timelock, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
				}
			})

		if err != nil {
			b.Logger.Errorw("Failed to deploy timelock",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"err", err,
			)
			return OpEVMDeployTimelockOutput{}, err
		}

		return OpEVMDeployTimelockOutput{
			Address: timelock.Address,
		}, nil
	})
