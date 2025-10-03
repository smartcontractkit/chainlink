package proposalutils

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsaptossdk "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// MCMSWithTimelockContracts holds the Go bindings
// for a MCMSWithTimelock contract deployment.
// It is public for use in product specific packages.
// Either all fields are nil or all fields are non-nil.
type MCMSWithTimelockContracts struct {
	CancellerMcm *owner_helpers.ManyChainMultiSig
	BypasserMcm  *owner_helpers.ManyChainMultiSig
	ProposerMcm  *owner_helpers.ManyChainMultiSig
	Timelock     *owner_helpers.RBACTimelock
	CallProxy    *owner_helpers.CallProxy
}

// Validate checks that all fields are non-nil, ensuring it's ready
// for use generating views or interactions.
func (state MCMSWithTimelockContracts) Validate() error {
	if state.Timelock == nil {
		return errors.New("timelock not found")
	}
	if state.CancellerMcm == nil {
		return errors.New("canceller not found")
	}
	if state.ProposerMcm == nil {
		return errors.New("proposer not found")
	}
	if state.BypasserMcm == nil {
		return errors.New("bypasser not found")
	}
	if state.CallProxy == nil {
		return errors.New("call proxy not found")
	}
	return nil
}

func McmsTimelockConverterForChain(chain uint64) (mcmssdk.TimelockConverter, error) {
	chainFamily, err := mcmstypes.GetChainSelectorFamily(mcmstypes.ChainSelector(chain))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for chain %d: %w", chain, err)
	}

	switch chainFamily {
	case chain_selectors.FamilyEVM:
		return &mcmsevmsdk.TimelockConverter{}, nil
	case chain_selectors.FamilySolana:
		return mcmssolanasdk.TimelockConverter{}, nil
	default:
		return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
	}
}

type mcmsInspectorOptions struct {
	AptosRole mcmsaptossdk.TimelockRole
}

type MCMSInspectorOption func(*mcmsInspectorOptions)

func WithAptosRole(role mcmsaptossdk.TimelockRole) MCMSInspectorOption {
	return func(opts *mcmsInspectorOptions) {
		opts.AptosRole = role
	}
}

func McmsInspectorForChain(env cldf.Environment, chain uint64, opts ...MCMSInspectorOption) (mcmssdk.Inspector, error) {
	var options mcmsInspectorOptions
	for _, opt := range opts {
		opt(&options)
	}

	chainFamily, err := mcmstypes.GetChainSelectorFamily(mcmstypes.ChainSelector(chain))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for chain %d: %w", chain, err)
	}

	switch chainFamily {
	case chain_selectors.FamilyEVM:
		return mcmsevmsdk.NewInspector(env.BlockChains.EVMChains()[chain].Client), nil
	case chain_selectors.FamilySolana:
		return mcmssolanasdk.NewInspector(env.BlockChains.SolanaChains()[chain].Client), nil
	case chain_selectors.FamilyAptos:
		if options.AptosRole.String() == "unknown" {
			return nil, fmt.Errorf("aptos role not properly set for chain: %d", chain)
		}
		inspector := mcmsaptossdk.NewInspector(env.BlockChains.AptosChains()[chain].Client, options.AptosRole)

		return inspector, nil
	default:
		return nil, fmt.Errorf("unsupported chain family %s", chainFamily)
	}
}

func McmsInspectors(env cldf.Environment) (map[uint64]mcmssdk.Inspector, error) {
	evmChains := env.BlockChains.EVMChains()
	solanaChains := env.BlockChains.SolanaChains()
	inspectors := make(map[uint64]mcmssdk.Inspector, len(evmChains)+len(solanaChains))

	for _, chain := range evmChains {
		var err error
		inspectors[chain.Selector], err = McmsInspectorForChain(env, chain.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms inspector for chain %s: %w", chain.String(), err)
		}
	}

	for _, chain := range solanaChains {
		var err error
		inspectors[chain.Selector], err = McmsInspectorForChain(env, chain.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms inspector for chain %s: %w", chain.String(), err)
		}
	}

	return inspectors, nil
}

func TransactionForChain(
	chain uint64, toAddress string, data []byte, value *big.Int, contractType string, tags []string,
) (mcmstypes.Transaction, error) {
	chainFamily, err := mcmstypes.GetChainSelectorFamily(mcmstypes.ChainSelector(chain))
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to get chain family for chain %d: %w", chain, err)
	}

	var tx mcmstypes.Transaction

	switch chainFamily {
	case chain_selectors.FamilyEVM:
		tx = mcmsevmsdk.NewTransaction(common.HexToAddress(toAddress), data, value, contractType, tags)

	case chain_selectors.FamilySolana:
		accounts := []*solana.AccountMeta{} // FIXME: how to pass accounts to support solana?
		var err error
		tx, err = mcmssolanasdk.NewTransaction(toAddress, data, value, accounts, contractType, tags)
		if err != nil {
			return mcmstypes.Transaction{}, fmt.Errorf("failed to create solana transaction: %w", err)
		}

	default:
		return mcmstypes.Transaction{}, fmt.Errorf("unsupported chain family %s", chainFamily)
	}

	return tx, nil
}

func BatchOperationForChain(
	chain uint64, toAddress string, data []byte, value *big.Int, contractType string, tags []string,
) (mcmstypes.BatchOperation, error) {
	tx, err := TransactionForChain(chain, toAddress, data, value, contractType, tags)
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to create transaction for chain: %w", err)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(chain),
		Transactions:  []mcmstypes.Transaction{tx},
	}, nil
}

func GetAptosRoleFromAction(action mcmstypes.TimelockAction) (mcmsaptossdk.TimelockRole, error) {
	switch action {
	case mcmstypes.TimelockActionSchedule:
		return mcmsaptossdk.TimelockRoleProposer, nil
	case mcmstypes.TimelockActionBypass:
		return mcmsaptossdk.TimelockRoleBypasser, nil
	case mcmstypes.TimelockActionCancel:
		return mcmsaptossdk.TimelockRoleCanceller, nil
	case "":
		// Default case for empty action to avoid breaking changes
		return mcmsaptossdk.TimelockRoleProposer, nil
	default:
		return mcmsaptossdk.TimelockRoleProposer, fmt.Errorf("invalid action: %s", action)
	}
}
