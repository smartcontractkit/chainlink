package ops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ccip_owner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"

	evmMcms "github.com/smartcontractkit/mcms/sdk/evm"
	"github.com/smartcontractkit/mcms/sdk/evm/bindings"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type OpEVMMCMSDeps struct {
	Chain    cldf_evm.Chain
	Backend  bind.ContractBackend
	Options  []func(*cldf.TypeAndVersion)
	AddrBook cldf.AddressBook
}

type OpEVMDeployMCMSInput struct {
	ContractType  cldf.ContractType
	ChainSelector uint64 // Needed to distinguish different input for Operations API
}

type OpEVMSetConfigMCMSInput struct {
	Address      common.Address
	ContractType cldf.ContractType
	MCMConfig    mcmsTypes.Config
}

type OpEVMDeployMCMSOutput struct {
	Address common.Address `json:"address"`
	Tx      common.Hash    `json:"tx"`
}

type OpEVMSetConfigMCMSOutput struct {
	Tx common.Hash `json:"tx"`
}

var OpEVMDeployMCMS = operations.NewOperation(
	"evm-mcms-deploy",
	semver.MustParse("1.0.0"),
	"Deploys MCMS contracts on the specified EVM chains",
	func(b operations.Bundle, deps OpEVMMCMSDeps, input OpEVMDeployMCMSInput) (OpEVMDeployMCMSOutput, error) {
		out := OpEVMDeployMCMSOutput{}

		mcm, err := cldf.DeployContract(b.Logger, deps.Chain, deps.AddrBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*bindings.ManyChainMultiSig] {
				addr, tx, _, err2 := ccip_owner.DeployManyChainMultiSig(
					deps.Chain.DeployerKey,
					deps.Backend,
				)

				tv := cldf.NewTypeAndVersion(input.ContractType, deployment.Version1_0_0)
				for _, option := range deps.Options {
					option(&tv)
				}

				return cldf.ContractDeploy[*bindings.ManyChainMultiSig]{
					Address: addr, Contract: nil, Tx: tx, Tv: tv, Err: err2,
				}
			})

		if err != nil {
			b.Logger.Errorw("Failed to deploy MCMS",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"err", err,
			)
			return out, err
		}

		// Confirm the transaction
		if _, err = deps.Chain.Confirm(mcm.Tx); err != nil {
			b.Logger.Errorw("Failed to confirm deployment",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"contractAddr", mcm.Address.String(),
				"err", err,
			)

			return out, err
		}

		tv := cldf.NewTypeAndVersion(input.ContractType, deployment.Version1_0_0)
		for _, option := range deps.Options {
			option(&tv)
		}

		return OpEVMDeployMCMSOutput{
			Address: mcm.Address,
			Tx:      mcm.Tx.Hash(),
		}, nil
	})

var OpEVMSetConfig = operations.NewOperation(
	"evm-mcms-set-config",
	semver.MustParse("1.0.0"),
	"Sets Config on the deployed MCMS contracts",
	func(b operations.Bundle, deps OpEVMMCMSDeps, input OpEVMSetConfigMCMSInput) (OpEVMSetConfigMCMSOutput, error) {
		out := OpEVMSetConfigMCMSOutput{}

		groupQuorums, groupParents, signerAddresses, signerGroups, err := evmMcms.ExtractSetConfigInputs(&input.MCMConfig)
		if err != nil {
			b.Logger.Errorw("Failed to extract set config inputs", "chain", deps.Chain.Name(), "err", err)
			return out, err
		}

		mcm, err := ccip_owner.NewManyChainMultiSig(input.Address, deps.Backend)
		if err != nil {
			b.Logger.Errorw("Failed to create ManyChainMultiSig instance",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"contractAddr", input.Address.String(),
				"err", err,
			)
			return out, err
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
			return out, err
		}

		// Confirm the transaction
		if _, err = deps.Chain.Confirm(tx); err != nil {
			b.Logger.Errorw("Failed to confirm deployment",
				"chainSelector", deps.Chain.ChainSelector(),
				"chainName", deps.Chain.Name(),
				"contractAddr", input.Address.String(),
				"err", err,
			)

			return out, err
		}

		return OpEVMSetConfigMCMSOutput{
			Tx: tx.Hash(),
		}, nil
	})
