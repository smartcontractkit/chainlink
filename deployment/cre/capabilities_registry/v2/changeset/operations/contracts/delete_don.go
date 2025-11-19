package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type DeleteDONDeps struct {
	Env                  *cldf.Environment
	Strategy             strategies.TransactionStrategy
	CapabilitiesRegistry *capabilities_registry_v2.CapabilitiesRegistry
}

// DeleteDONInput is the user-provided input
type DeleteDONInput struct {
	ChainSelector uint64
	DonNames      []string

	// Optional MCMS config: if provided, the tx will be wrapped in a proposal and not immediately executed
	MCMSConfig *contracts.MCMSConfig
}

func (r *DeleteDONInput) Validate() error {
	if len(r.DonNames) == 0 {
		return errors.New("must specify don names")
	}
	for _, name := range r.DonNames {
		if name == "" {
			return errors.New("don names cannot contain empty string")
		}
	}
	return nil
}

// DeleteDONOutput returns the resulting MCMS operation (if any) and basic confirmation
type DeleteDONOutput struct {
	Operation *mcmstypes.BatchOperation
	// DeletedNames echoes the requested DON names
	DeletedNames []string
}

// DeleteDON implements a simple, safe-by-default removal of one or more DONs by name.
var DeleteDON = operations.NewOperation[DeleteDONInput, DeleteDONOutput, DeleteDONDeps](
	"delete-don-op",
	semver.MustParse("1.0.0"),
	"Delete DON(s) from Capabilities Registry by name",
	func(b operations.Bundle, deps DeleteDONDeps, input DeleteDONInput) (DeleteDONOutput, error) {
		if err := input.Validate(); err != nil {
			return DeleteDONOutput{}, err
		}

		registry := deps.CapabilitiesRegistry
		for _, name := range input.DonNames {
			if _, err := registry.GetDONByName(&bind.CallOpts{}, name); err != nil {
				return DeleteDONOutput{}, cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			}
		}

		// Execute the transaction using the strategy; delete all provided names in one call
		operation, tx, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return registry.RemoveDONsByName(opts, input.DonNames)
		})
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return DeleteDONOutput{}, fmt.Errorf("failed to execute RemoveDONsByName: %w", err)
		}

		// If using MCMS, return the prepared operation without mining
		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for DeleteDON %v on chain %d", input.DonNames, input.ChainSelector)
			return DeleteDONOutput{Operation: operation, DeletedNames: input.DonNames}, nil
		}

		txHash := ""
		if tx != nil {
			txHash = tx.Hash().String()
		}
		deps.Env.Logger.Infof("Submitted DeleteDON for %v on chain %d, tx hash %q", input.DonNames, input.ChainSelector, txHash)

		return DeleteDONOutput{
			Operation:    operation,
			DeletedNames: input.DonNames,
		}, nil
	},
)
