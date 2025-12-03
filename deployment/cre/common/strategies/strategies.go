package strategies

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

// TransactionStrategy interface for executing transactions with different strategies
type TransactionStrategy interface {
	// Apply executes the provided call function and returns the resulting MCMS batch operation if applicable.
	// The callFn should accept transaction options and return a transaction or an error.
	// If using MCMS, the returned BatchOperation can be used to build a proposal.
	// If no MCMS is used, the returned BatchOperation will be nil, and the transaction will be confirmed.
	Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (*mcmstypes.BatchOperation, *types.Transaction, error)

	// BuildProposal constructs a TimelockProposal from the provided batch operations.
	// This is only applicable when using MCMS; otherwise, it returns an empty proposal.
	BuildProposal(operations []mcmstypes.BatchOperation) (*mcmslib.TimelockProposal, error)
}

// CreateStrategy is a factory function to create the appropriate strategy based on configuration
func CreateStrategy(
	chain cldf_evm.Chain,
	env cldf.Environment,
	mcmsConfig *contracts.MCMSConfig,
	mcmsContracts *commonchangeset.MCMSWithTimelockState,
	targetAddress common.Address,
	description string,
) (TransactionStrategy, error) {
	if mcmsConfig != nil {
		if mcmsContracts == nil {
			return nil, errors.New("MCMS contracts are required when mcmsConfig is not nil")
		}

		err := mcmsConfig.Validate(chain, *mcmsContracts)
		if err != nil {
			return nil, fmt.Errorf("invalid MCMS configuration: %w", err)
		}

		return &MCMSTransaction{
			Config:        mcmsConfig,
			Description:   description,
			Address:       targetAddress,
			ChainSel:      chain.ChainSelector(),
			MCMSContracts: mcmsContracts,
			Env:           env,
		}, nil
	}

	return &SimpleTransaction{Chain: chain}, nil
}

// Legacy aliases for backward compatibility with existing CRE modules
// Deprecated: Use TransactionStrategy instead
type StrategyV2 = TransactionStrategy

// Deprecated: Use SimpleTransaction instead
type SimpleTransactionV2 = SimpleTransaction

// Deprecated: Use MCMSTransaction instead
type MCMSTransactionV2 = MCMSTransaction

// GetMCMSContracts retrieves MCMS contracts from the environment using merged approach (both DataStore and AddressBook)
func GetMCMSContracts(e cldf.Environment, chainSelector uint64, mcmsConfig contracts.MCMSConfig) (*commonchangeset.MCMSWithTimelockState, error) {
	var qualifier string

	qualifier, ok := mcmsConfig.TimelockQualifierPerChain[chainSelector]
	if !ok {
		return nil, fmt.Errorf("no timelock qualifier found for chain selector %d in MCMSConfig", chainSelector)
	}

	states, err := commonchangeset.MaybeLoadMCMSWithTimelockStateWithQualifier(e, []uint64{chainSelector}, qualifier)
	if err != nil {
		return nil, fmt.Errorf("failed to load MCMS contracts for chain %d: %w", chainSelector, err)
	}

	state, ok := states[chainSelector]
	if !ok {
		return nil, fmt.Errorf("MCMS contracts not found for chain %d", chainSelector)
	}

	return state, nil
}
