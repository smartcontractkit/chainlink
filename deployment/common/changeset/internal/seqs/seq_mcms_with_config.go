package seqs

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/ops"
)

type SeqDeployMCMSWithConfigDeps struct {
	Chain    cldf_evm.Chain
	AddrBook cldf.AddressBook
	Backend  bind.ContractBackend
	Options  []func(*cldf.TypeAndVersion)
}

type SeqDeployMCMSWithConfigInput struct {
	ContractType cldf.ContractType
	MCMConfig    mcmsTypes.Config
	ChainSelector uint64
}

type SeqDeployMCMSWithConfigOutput struct {
	Address common.Address `json:"address"`
}

var SeqEVMDeployMCMSWithConfig = operations.NewSequence(
	"seq-deploy-mcms-with-config",
	semver.MustParse("1.0.0"),
	"Deploys MCMS contract & sets config",
	func(b operations.Bundle, deps SeqDeployMCMSWithConfigDeps, in SeqDeployMCMSWithConfigInput) (SeqDeployMCMSWithConfigOutput, error) {
		out := SeqDeployMCMSWithConfigOutput{}
		// Deploy MCMS contract
		deployReport, err := operations.ExecuteOperation(b, ops.OpEVMDeployMCMS,
			ops.OpEVMMCMSDeps{
				Chain:   deps.Chain,
				Backend: deps.Backend,
				Options: deps.Options,
				AddrBook: deps.AddrBook,
			},
			ops.OpEVMDeployMCMSInput{
				ContractType: in.ContractType,
				ChainSelector: in.ChainSelector,
			},
		)
		if err != nil {
			return out, err
		}

		out.Address = deployReport.Output.Address

		// Set config
		_, err = operations.ExecuteOperation(b, ops.OpEVMSetConfig,
			ops.OpEVMMCMSDeps{
				Chain:   deps.Chain,
				Backend: deps.Backend,
				Options: deps.Options,
			},
			ops.OpEVMSetConfigMCMSInput{
				Address:      deployReport.Output.Address,
				ContractType: in.ContractType,
				MCMConfig:    in.MCMConfig,
			},
		)
		if err != nil {
			return out, err
		}

		return out, nil
	},
)
