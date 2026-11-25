package strategies

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

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

// ownerABI is the one function every target of CreateStrategy must expose (Ownable / ConfirmedOwner).
// A target without owner() is a wiring bug and fails loudly here rather than as a reverted proposal
// or transaction later. Calls that are not owner-only use CreateStrategyForNonOwnerCall.
const ownerABI = `[{"type":"function","name":"owner","inputs":[],"outputs":[{"name":"","type":"address"}],"stateMutability":"view"}]`

// chainLabel renders the chain for error messages, without repeating the selector when it is the name.
func chainLabel(chain cldf_evm.Chain) string {
	sel := strconv.FormatUint(chain.Selector, 10)
	if chain.Name() == sel {
		return sel
	}
	return fmt.Sprintf("%s (%s)", chain.Name(), sel)
}

// readOwner reads owner() of the contract at addr from live chain state.
func readOwner(chain cldf_evm.Chain, addr common.Address) (common.Address, error) {
	parsed, err := abi.JSON(strings.NewReader(ownerABI))
	if err != nil {
		return common.Address{}, err
	}
	var out []any
	err = bind.NewBoundContract(addr, parsed, chain.Client, chain.Client, chain.Client).
		Call(&bind.CallOpts{Context: context.Background()}, &out, "owner")
	if err != nil {
		return common.Address{}, err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// verifyOwner checks that the live owner of the target contract matches the way the caller
// wants to send the transaction: the timelock when mcmsConfig is set, the deployer key otherwise.
func verifyOwner(
	chain cldf_evm.Chain,
	mcmsConfig *contracts.MCMSConfig,
	mcmsContracts *evmstate.MCMSWithTimelockState,
	target common.Address,
) error {
	owner, err := readOwner(chain, target)
	if err != nil {
		return fmt.Errorf("failed to read owner() of contract %s on chain %s: %w", target, chainLabel(chain), err)
	}

	if mcmsConfig != nil {
		if mcmsContracts.Timelock == nil {
			return errors.New("MCMS contracts are missing the timelock")
		}
		timelock := mcmsContracts.Timelock.Address()
		if owner == timelock {
			return nil
		}
		if chain.DeployerKey != nil && owner == chain.DeployerKey.From {
			return fmt.Errorf("contract %s on chain %s is owned by the deployer key %s but mcmsConfig names timelock %s: "+
				"remove mcmsConfig, or, if ownership was handed to the timelock, the accept-ownership proposal has not executed yet",
				target, chainLabel(chain), owner, timelock)
		}
		return fmt.Errorf("mcmsConfig names timelock %s but contract %s on chain %s is owned by %s: fix mcmsConfig (timelock qualifier) to match the owner",
			timelock, target, chainLabel(chain), owner)
	}

	if chain.DeployerKey == nil {
		return errors.New("chain has no deployer key and mcmsConfig is not set")
	}
	if owner == chain.DeployerKey.From {
		return nil
	}
	return fmt.Errorf("contract %s on chain %s is owned by %s but the deployer key is %s: "+
		"if %s is a timelock, provide mcmsConfig; otherwise the deployer key cannot send this transaction",
		target, chainLabel(chain), owner, chain.DeployerKey.From, owner)
}

// CreateStrategy is a factory function to create the appropriate strategy based on configuration.
// It reads owner() of targetAddress and fails if the owner does not match the chosen strategy.
func CreateStrategy(
	chain cldf_evm.Chain,
	env cldf.Environment,
	mcmsConfig *contracts.MCMSConfig,
	mcmsContracts *evmstate.MCMSWithTimelockState,
	targetAddress common.Address,
	description string,
) (TransactionStrategy, error) {
	return buildStrategy(chain, env, mcmsConfig, mcmsContracts, targetAddress, description, true)
}

// CreateStrategyForNonOwnerCall is CreateStrategy without the owner check. It is only for calls the
// contract accepts from accounts other than its owner, such as a workflow owner registering or
// pausing its own workflows on the workflow registry: there the sender is deliberately not the
// registry's owner, so comparing it with owner() would reject a valid call. Owner-only calls must
// use CreateStrategy.
func CreateStrategyForNonOwnerCall(
	chain cldf_evm.Chain,
	env cldf.Environment,
	mcmsConfig *contracts.MCMSConfig,
	mcmsContracts *evmstate.MCMSWithTimelockState,
	targetAddress common.Address,
	description string,
) (TransactionStrategy, error) {
	return buildStrategy(chain, env, mcmsConfig, mcmsContracts, targetAddress, description, false)
}

func buildStrategy(
	chain cldf_evm.Chain,
	env cldf.Environment,
	mcmsConfig *contracts.MCMSConfig,
	mcmsContracts *evmstate.MCMSWithTimelockState,
	targetAddress common.Address,
	description string,
	checkOwner bool,
) (TransactionStrategy, error) {
	if mcmsConfig != nil {
		if mcmsContracts == nil {
			return nil, errors.New("MCMS contracts are required when mcmsConfig is not nil")
		}

		err := mcmsConfig.Validate(chain, *mcmsContracts)
		if err != nil {
			return nil, fmt.Errorf("invalid MCMS configuration: %w", err)
		}

		if checkOwner {
			if err = verifyOwner(chain, mcmsConfig, mcmsContracts, targetAddress); err != nil {
				return nil, err
			}
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

	if checkOwner {
		if err := verifyOwner(chain, nil, nil, targetAddress); err != nil {
			return nil, err
		}
	}

	return &SimpleTransaction{Chain: chain}, nil
}

// Legacy aliases for backward compatibility with existing CRE modules
//
// Deprecated: Use TransactionStrategy instead
type StrategyV2 = TransactionStrategy

// Deprecated: Use SimpleTransaction instead
type SimpleTransactionV2 = SimpleTransaction

// Deprecated: Use MCMSTransaction instead
type MCMSTransactionV2 = MCMSTransaction

// GetMCMSContracts retrieves MCMS contracts from the environment using merged approach (both DataStore and AddressBook)
func GetMCMSContracts(e cldf.Environment, chainSelector uint64, mcmsConfig contracts.MCMSConfig) (*evmstate.MCMSWithTimelockState, error) {
	var qualifier string

	qualifier, ok := mcmsConfig.TimelockQualifierPerChain[chainSelector]
	if !ok {
		return nil, fmt.Errorf("no timelock qualifier found for chain selector %d in MCMSConfig", chainSelector)
	}

	states, err := evmstate.MaybeLoadMCMSWithTimelockStateWithQualifier(e, []uint64{chainSelector}, qualifier)
	if err != nil {
		return nil, fmt.Errorf("failed to load MCMS contracts for chain %d: %w", chainSelector, err)
	}

	state, ok := states[chainSelector]
	if !ok {
		return nil, fmt.Errorf("MCMS contracts not found for chain %d", chainSelector)
	}

	return state, nil
}
