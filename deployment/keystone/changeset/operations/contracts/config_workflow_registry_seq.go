package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type ConfigWorkflowRegistrySeqDeps struct {
	Env *cldf.Environment
}

type ConfigWorkflowRegistrySeqInput struct {
	ContractAddress       common.Address
	RegistryChainSelector uint64
	AllowedDonIDs         []uint32
	WorkflowOwners        []common.Address
	MCMSConfig            *changeset.MCMSConfig
}

type ConfigWorkflowRegistrySeqOutput struct {
	RegistryChainSelector uint64
	AllowedDonIDs         []uint32
	WorkflowOwners        []common.Address
}

var ConfigWorkflowRegistrySeq = operations.NewSequence[ConfigWorkflowRegistrySeqInput, ConfigWorkflowRegistrySeqOutput, ConfigWorkflowRegistrySeqDeps](
	"config-workflow-registry-seq",
	semver.MustParse("1.0.0"),
	"Configure Workflow Registry",
	func(b operations.Bundle, deps ConfigWorkflowRegistrySeqDeps, input ConfigWorkflowRegistrySeqInput) (ConfigWorkflowRegistrySeqOutput, error) {
		if len(input.AllowedDonIDs) == 0 {
			return ConfigWorkflowRegistrySeqOutput{}, errors.New("allowed don ids not set")
		}
		if len(input.WorkflowOwners) == 0 {
			return ConfigWorkflowRegistrySeqOutput{}, errors.New("workflow owners not set")
		}

		_, err := operations.ExecuteOperation(b, UpdateAllowedDonsOp, UpdateAllowedDonsOpDeps(deps), UpdateAllowedDonsOpInput{
			ContractAddress:  input.ContractAddress.Hex(),
			RegistryChainSel: input.RegistryChainSelector,
			DonIDs:           input.AllowedDonIDs,
			Allowed:          true,
			MCMSConfig:       input.MCMSConfig,
		})
		if err != nil {
			return ConfigWorkflowRegistrySeqOutput{}, fmt.Errorf("failed to update allowed Dons: %w", err)
		}

		addresses := make([]string, 0, len(input.WorkflowOwners))
		for _, owner := range input.WorkflowOwners {
			addresses = append(addresses, owner.Hex())
		}

		_, err = operations.ExecuteOperation(b, UpdateAuthorizedAddressesOp, UpdateAuthorizedAddressesOpDeps(deps), UpdateAuthorizedAddressesOpInput{
			ContractAddress:  input.ContractAddress.Hex(),
			RegistryChainSel: input.RegistryChainSelector,
			Addresses:        addresses,
			Allowed:          true,
			MCMSConfig:       input.MCMSConfig,
		})
		if err != nil {
			return ConfigWorkflowRegistrySeqOutput{}, fmt.Errorf("failed to update authorized addresses: %w", err)
		}

		return ConfigWorkflowRegistrySeqOutput{
			RegistryChainSelector: input.RegistryChainSelector,
			AllowedDonIDs:         input.AllowedDonIDs,
			WorkflowOwners:        input.WorkflowOwners,
		}, nil
	},
)
