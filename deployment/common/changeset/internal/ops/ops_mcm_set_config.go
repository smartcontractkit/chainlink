package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	bindings "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	evmMcms "github.com/smartcontractkit/mcms/sdk/evm"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type OpEVMSetConfigMCMInput struct {
	Address      common.Address    `json:"address"`
	ContractType cldf.ContractType `json:"contractType"`
	MCMConfig    mcmsTypes.Config  `json:"mcmConfig"` // Config for the ManyChainMultiSig contract
}

type OpEVMSetConfigMCMOutput struct {
	Tx common.Hash `json:"tx"`
}

var OpEVMSetConfigMCM = operations.NewOperation(
	"evm-mcm-set-config",
	semver.MustParse("1.0.0"),
	"Sets Config on the deployed MCM contract",
	func(b operations.Bundle, deps OpEVMMCMSDeps, input OpEVMSetConfigMCMInput) (OpEVMSetConfigMCMOutput, error) {

		groupQuorums, groupParents, signerAddresses, signerGroups, err := evmMcms.ExtractSetConfigInputs(&input.MCMConfig)
		if err != nil {
			b.Logger.Errorw("Failed to extract set config inputs", "chain", deps.Chain.Name(), "err", err)
			return OpEVMSetConfigMCMOutput{}, err
		}

		mcm, err := bindings.NewManyChainMultiSig(input.Address, deps.Chain.Client)
		if err != nil {
			b.Logger.Errorw("Failed to create ManyChainMultiSig instance",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"contractAddr", input.Address.String(),
				"err", err,
			)
			return OpEVMSetConfigMCMOutput{}, err
		}

		tx, err := mcm.SetConfig(deps.Chain.DeployerKey,
			signerAddresses,
			// Signer 1 is int group 0 (root group) with quorum 1.
			signerGroups,
			groupQuorums,
			groupParents,
			false,
		)
		if err != nil {
			b.Logger.Errorw("Failed to Set MCM config",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"err", err,
			)
			return OpEVMSetConfigMCMOutput{}, err
		}

		// Confirm the transaction
		if _, err = deps.Chain.Confirm(tx); err != nil {
			b.Logger.Errorw("Failed to confirm deployment",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"contractAddr", input.Address.String(),
				"err", err,
			)

			return OpEVMSetConfigMCMOutput{Tx: tx.Hash()}, err
		}

		return OpEVMSetConfigMCMOutput{
			Tx: tx.Hash(),
		}, nil
	})
