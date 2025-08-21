package seqs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms/ops"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

type SeqDeployMCMWithConfigInput struct {
	ContractType   cldf.ContractType           `json:"contractType"`
	MCMConfig      mcmsTypes.Config            `json:"mcmConfig"`
	ChainSelector  uint64                      `json:"chainSelector"`
	GasBoostConfig *commontypes.GasBoostConfig `json:"gasBoostConfig"`
}

type SeqDeployMCMWithConfigOutput struct {
	Address common.Address `json:"address"`
}

var SeqEVMDeployMCMWithConfig = operations.NewSequence(
	"seq-deploy-mcm-with-config",
	semver.MustParse("1.0.0"),
	"Deploys MCM contract & sets config",
	func(b operations.Bundle, deps cldf_evm.Chain, in SeqDeployMCMWithConfigInput) (opsutils.EVMDeployOutput, error) {
		// Deploy MCM contract
		var deployReport operations.Report[opsutils.EVMDeployInput[any], opsutils.EVMDeployOutput]
		var deployErr error
		switch in.ContractType {
		case commontypes.BypasserManyChainMultisig:
			deployReport, deployErr = operations.ExecuteOperation(b, ops.OpEVMDeployBypasserMCM, deps, opsutils.EVMDeployInput[any]{
				ChainSelector: in.ChainSelector,
			}, opsutils.RetryDeploymentWithGasBoost[any](in.GasBoostConfig))
		case commontypes.ProposerManyChainMultisig:
			deployReport, deployErr = operations.ExecuteOperation(b, ops.OpEVMDeployProposerMCM, deps, opsutils.EVMDeployInput[any]{
				ChainSelector: in.ChainSelector,
			}, opsutils.RetryDeploymentWithGasBoost[any](in.GasBoostConfig))
		case commontypes.CancellerManyChainMultisig:
			deployReport, deployErr = operations.ExecuteOperation(b, ops.OpEVMDeployCancellerMCM, deps, opsutils.EVMDeployInput[any]{
				ChainSelector: in.ChainSelector,
			}, opsutils.RetryDeploymentWithGasBoost[any](in.GasBoostConfig))
		default:
			return opsutils.EVMDeployOutput{}, fmt.Errorf("unsupported contract type for seq-deploy-mcm-with-config: %s", in.ContractType)
		}
		if deployErr != nil {
			return opsutils.EVMDeployOutput{}, fmt.Errorf("failed to deploy %s: %w", in.ContractType, deployErr)
		}

		// Set config
		groupQuorums, groupParents, signerAddresses, signerGroups, err := evm.ExtractSetConfigInputs(&in.MCMConfig)
		if err != nil {
			return opsutils.EVMDeployOutput{}, err
		}
		_, err = operations.ExecuteOperation(b, ops.OpEVMSetConfigMCM,
			deps,
			opsutils.EVMCallInput[ops.OpEVMSetConfigMCMInput]{
				ChainSelector: in.ChainSelector,
				Address:       deployReport.Output.Address,
				NoSend:        false,
				CallInput: ops.OpEVMSetConfigMCMInput{
					SignerAddresses: signerAddresses,
					SignerGroups:    signerGroups,
					GroupQuorums:    groupQuorums,
					GroupParents:    groupParents,
				},
			},
			opsutils.RetryCallWithGasBoost[ops.OpEVMSetConfigMCMInput](in.GasBoostConfig),
		)
		if err != nil {
			return opsutils.EVMDeployOutput{}, err
		}

		return deployReport.Output, nil
	},
)
