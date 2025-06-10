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

type OpEVMDeployCallProxyInput struct {
	ContractType  cldf.ContractType
	Timelock      common.Address `json:"timelock"`
	ChainSelector uint64         // Needed to distinguish different input for Operations API
}

type OpEVMDeployCallProxyOutput struct {
	Address common.Address `json:"address"`
}

var OpEVMDeployCallProxy = operations.NewOperation(
	"evm-call-proxy-deploy",
	semver.MustParse("1.0.0"),
	"Deploys CallProxy contract on the specified EVM chains",
	func(b operations.Bundle, deps OpEVMMCMSDeps, input OpEVMDeployCallProxyInput) (OpEVMDeployCallProxyOutput, error) {
		timelock, err := cldf.DeployContract(b.Logger, deps.Chain, deps.AddrBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.CallProxy] {
				var (
					timelock common.Address
					tx2      *types.Transaction
					cc       *bindings.CallProxy
					err2     error
				)
				if chain.IsZkSyncVM {
					timelock, _, cc, err2 = mcmsnew_zksync.DeployCallProxyZk(
						nil,
						chain.ClientZkSyncVM,
						chain.DeployerKeyZkSyncVM,
						chain.Client,
						input.Timelock,
					)
				} else {
					timelock, tx2, cc, err2 = bindings.DeployCallProxy(
						chain.DeployerKey,
						chain.Client,
						input.Timelock,
					)
				}

				tv := cldf.NewTypeAndVersion(input.ContractType, deployment.Version1_0_0)
				for _, option := range deps.Options {
					option(&tv)
				}

				return cldf.ContractDeploy[*bindings.CallProxy]{
					Address: timelock, Contract: cc, Tx: tx2, Tv: tv, Err: err2,
				}
			})

		if err != nil {
			b.Logger.Errorw("Failed to deploy CallProxy",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"err", err,
			)
			return OpEVMDeployCallProxyOutput{}, err
		}

		return OpEVMDeployCallProxyOutput{
			Address: timelock.Address,
		}, nil
	})
