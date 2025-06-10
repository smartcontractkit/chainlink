package seqs

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/ops"
)

type SeqDeployMCMWithConfigDeps struct {
	Chain    cldf_evm.Chain
	AddrBook cldf.AddressBook
	Options  []func(*cldf.TypeAndVersion)
}

type SeqDeployMCMWithConfigInput struct {
	ContractType  cldf.ContractType `json:"contractType"`
	MCMConfig     mcmsTypes.Config  `json:"mcmConfig"`
	ChainSelector uint64            `json:"chainSelector"`
}

type SeqDeployMCMWithConfigOutput struct {
	Address common.Address `json:"address"`
}

var SeqEVMDeployMCMWithConfig = operations.NewSequence(
	"seq-deploy-mcm-with-config",
	semver.MustParse("1.0.0"),
	"Deploys MCM contract & sets config",
	func(b operations.Bundle, deps SeqDeployMCMWithConfigDeps, in SeqDeployMCMWithConfigInput) (SeqDeployMCMWithConfigOutput, error) {
		seqOp := SeqDeployMCMWithConfigOutput{}
		// Deploy MCM contract
		deployReport, err := operations.ExecuteOperation(b, ops.OpEVMDeployMCM,
			ops.OpEVMMCMSDeps{
				Chain:    deps.Chain,
				Options:  deps.Options,
				AddrBook: deps.AddrBook,
			},
			ops.OpEVMDeployMCMInput{
				ContractType:  in.ContractType,
				ChainSelector: in.ChainSelector,
			},
		)
		if err != nil {
			return seqOp, err
		}

		seqOp.Address = deployReport.Output.Address

		// Set config
		_, err = operations.ExecuteOperation(b, ops.OpEVMSetConfigMCM,
			ops.OpEVMMCMSDeps{
				Chain:   deps.Chain,
				Options: deps.Options,
			},
			ops.OpEVMSetConfigMCMInput{
				Address:      deployReport.Output.Address,
				ContractType: in.ContractType,
				MCMConfig:    in.MCMConfig,
			},
		)
		if err != nil {
			return seqOp, err
		}

		return seqOp, nil
	},
)
