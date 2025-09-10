package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

type UpdateAllowedDonsOpDeps struct {
	Env *cldf.Environment
}

type UpdateAllowedDonsOpInput struct {
	ContractAddress  string
	RegistryChainSel uint64
	DonIDs           []uint32
	Allowed          bool
	MCMSConfig       *changeset.MCMSConfig
}

type UpdateAllowedDonsOpOutput struct {
	Proposals []mcms.TimelockProposal
}

var UpdateAllowedDonsOp = operations.NewOperation[UpdateAllowedDonsOpInput, UpdateAllowedDonsOpOutput, UpdateAllowedDonsOpDeps](
	"update-allowed-dons-op",
	semver.MustParse("1.0.0"),
	"Update Allowed DONs in Workflow Registry",
	func(b operations.Bundle, deps UpdateAllowedDonsOpDeps, input UpdateAllowedDonsOpInput) (UpdateAllowedDonsOpOutput, error) {
		if len(input.DonIDs) == 0 {
			return UpdateAllowedDonsOpOutput{}, errors.New("must provide at least one DonID")
		}

		evmChains := deps.Env.BlockChains.EVMChains()
		chain, ok := evmChains[input.RegistryChainSel]
		if !ok {
			return UpdateAllowedDonsOpOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", input.RegistryChainSel)
		}

		registry, err := contracts.GetOwnedContractV2[*workflow_registry.WorkflowRegistry](deps.Env.DataStore.Addresses(), chain, input.ContractAddress)
		if err != nil {
			return UpdateAllowedDonsOpOutput{}, fmt.Errorf("failed to get workflow registry contract: %w", err)
		}

		var s workflowregistry.StrategyV2
		if input.MCMSConfig != nil {
			if registry.McmsContracts == nil {
				return UpdateAllowedDonsOpOutput{}, errors.New("registry must be owned by MCMS")
			}

			s = &workflowregistry.MCMSTransactionV2{
				Config:        input.MCMSConfig,
				Description:   "proposal to update allowed dons",
				Address:       registry.Contract.Address(),
				ChainSel:      input.RegistryChainSel,
				MCMSContracts: registry.McmsContracts,
				Env:           *deps.Env,
			}
		} else {
			s = &workflowregistry.SimpleTransactionV2{
				Chain: chain,
			}
		}

		proposals, err := s.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.Contract.UpdateAllowedDONs(opts, input.DonIDs, input.Allowed)
			if err != nil {
				err = cldf.DecodeErr(workflow_registry.WorkflowRegistryABI, err)
			}
			return tx, err
		})
		if err != nil {
			return UpdateAllowedDonsOpOutput{}, fmt.Errorf("failed to update allowed dons: %w", err)
		}

		return UpdateAllowedDonsOpOutput{Proposals: proposals}, nil
	},
)
