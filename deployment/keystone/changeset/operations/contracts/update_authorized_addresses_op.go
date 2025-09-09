package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

type UpdateAuthorizedAddressesOpDeps struct {
	Env *cldf.Environment
}

type UpdateAuthorizedAddressesOpInput struct {
	ContractAddress  string
	RegistryChainSel uint64
	Addresses        []string
	Allowed          bool
	MCMSConfig       *changeset.MCMSConfig
}

type UpdateAuthorizedAddressesOpOutput struct {
	Proposals []mcms.TimelockProposal
}

var UpdateAuthorizedAddressesOp = operations.NewOperation[UpdateAuthorizedAddressesOpInput, UpdateAuthorizedAddressesOpOutput, UpdateAuthorizedAddressesOpDeps](
	"update-authorized-addresses-op",
	semver.MustParse("1.0.0"),
	"Update Authorized Addresses in Workflow Registry",
	func(b operations.Bundle, deps UpdateAuthorizedAddressesOpDeps, input UpdateAuthorizedAddressesOpInput) (UpdateAuthorizedAddressesOpOutput, error) {
		if len(input.Addresses) == 0 {
			return UpdateAuthorizedAddressesOpOutput{}, errors.New("must provide at least 1 address")
		}

		evmChains := deps.Env.BlockChains.EVMChains()
		chain, ok := evmChains[input.RegistryChainSel]
		if !ok {
			return UpdateAuthorizedAddressesOpOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", input.RegistryChainSel)
		}

		registry, err := contracts.GetOwnedContractV2[*workflow_registry.WorkflowRegistry](deps.Env.DataStore.Addresses(), chain, input.ContractAddress)
		if err != nil {
			return UpdateAuthorizedAddressesOpOutput{}, fmt.Errorf("failed to get workflow registry contract: %w", err)
		}

		var addr []common.Address
		for _, a := range input.Addresses {
			addr = append(addr, common.HexToAddress(a))
		}

		var s workflowregistry.StrategyV2
		if input.MCMSConfig != nil {
			if registry.McmsContracts == nil {
				return UpdateAuthorizedAddressesOpOutput{}, errors.New("registry must be owned by MCMS")
			}

			s = &workflowregistry.MCMSTransactionV2{
				Config:        input.MCMSConfig,
				Description:   "proposal to update authorized addresses",
				Address:       registry.Contract.Address(),
				ChainSel:      chain.Selector,
				MCMSContracts: registry.McmsContracts,
				Env:           *deps.Env,
			}
		} else {
			s = &workflowregistry.SimpleTransactionV2{
				Chain: chain,
			}
		}

		proposals, err := s.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.Contract.UpdateAuthorizedAddresses(opts, addr, input.Allowed)
			if err != nil {
				err = cldf.DecodeErr(workflow_registry.WorkflowRegistryABI, err)
			}
			return tx, err
		})
		if err != nil {
			return UpdateAuthorizedAddressesOpOutput{}, fmt.Errorf("failed to update authorized addresses: %w", err)
		}

		return UpdateAuthorizedAddressesOpOutput{Proposals: proposals}, nil
	},
)
