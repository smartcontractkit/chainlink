package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	// TxHash is set when executed immediately (no MCMS)
	TxHash common.Hash
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
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeleteDONOutput{}, cldf.ErrChainNotFound
		}

		missing := make([]string, 0)
		for _, name := range input.DonNames {
			if _, err := registry.GetDONByName(&bind.CallOpts{}, name); err != nil {
				// Treat any revert as non-existent for this validation step
				missing = append(missing, name)
			}
		}
		if len(missing) > 0 {
			return DeleteDONOutput{}, fmt.Errorf("the following DON doesn't exist, or failed to retrieve it from the contract: %v", missing)
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

		deps.Env.Logger.Infof("Submitted DeleteDON for %v on chain %d, tx %s", input.DonNames, input.ChainSelector, tx.Hash())
		// Wait for inclusion when executing immediately
		ctx := b.GetContext()
		if _, err := bind.WaitMined(ctx, chain.Client, tx); err != nil {
			return DeleteDONOutput{}, fmt.Errorf("failed to mine RemoveDONsByName transaction %s: %w", tx.Hash(), err)
		}

		// Post condition: verify each name is now gone
		for _, name := range input.DonNames {
			if _, err := registry.GetDONByName(&bind.CallOpts{}, name); err == nil {
				return DeleteDONOutput{}, fmt.Errorf("DON '%s' still exists after deletion", name)
			}
		}

		deps.Env.Logger.Infof("Successfully deleted DONs %v on chain %d", input.DonNames, input.ChainSelector)

		return DeleteDONOutput{
			Operation:    operation,
			TxHash:       tx.Hash(),
			DeletedNames: input.DonNames,
		}, nil
	},
)
