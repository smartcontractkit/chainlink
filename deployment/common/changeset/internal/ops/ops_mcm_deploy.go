package ops

import (
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

// common dependency for MCMS Ops
type OpEVMMCMSDeps struct {
	Chain    cldf_evm.Chain
	Options  []func(*cldf.TypeAndVersion)
	AddrBook cldf.AddressBook
}

type OpEVMDeployMCMInput struct {
	ContractType  cldf.ContractType
	ChainSelector uint64 // Needed to distinguish different input for Operations API
}

type OpEVMDeployMCMOutput struct {
	Address common.Address `json:"address"`
}

var OpEVMDeployMCM = operations.NewOperation(
	"evm-mcm-deploy",
	semver.MustParse("1.0.0"),
	"Deploys MCM contracts on the specified EVM chains",
	func(b operations.Bundle, deps OpEVMMCMSDeps, input OpEVMDeployMCMInput) (OpEVMDeployMCMOutput, error) {
		out := OpEVMDeployMCMOutput{}

		mcm, err := cldf.DeployContract(b.Logger, deps.Chain, deps.AddrBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.ManyChainMultiSig] {
				var (
					mcmAddr common.Address
					tx      *types.Transaction
					mcm     *bindings.ManyChainMultiSig
					err2    error
				)
				if chain.IsZkSyncVM {
					mcmAddr, _, mcm, err2 = mcmsnew_zksync.DeployManyChainMultiSigZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
					)
				} else {
					mcmAddr, tx, mcm, err2 = bindings.DeployManyChainMultiSig(
						deps.Chain.DeployerKey,
						deps.Chain.Client,
					)
				}

				tv := cldf.NewTypeAndVersion(input.ContractType, deployment.Version1_0_0)
				for _, option := range deps.Options {
					option(&tv)
				}

				return cldf.ContractDeploy[*bindings.ManyChainMultiSig]{
					Address: mcmAddr, Contract: mcm, Tx: tx, Tv: tv, Err: err2,
				}
			})

		if err != nil {
			b.Logger.Errorw("Failed to deploy MCM",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"err", err,
			)
			return out, err
		}

		return OpEVMDeployMCMOutput{
			Address: mcm.Address,
		}, nil
	})
