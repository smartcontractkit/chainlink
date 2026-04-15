package proposalutils

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	mcmschainwrappers "github.com/smartcontractkit/mcms/chainwrappers"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"

	cldfmcmsadapters "github.com/smartcontractkit/chainlink-deployments-framework/chain/mcms/adapters"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmsaptossdk "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
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

	action := mcmstypes.TimelockActionSchedule
	if options.AptosRole.String() != "unknown" {
		var err error
		action, err = mcmsaptossdk.ActionFromAptosRole(options.AptosRole)
		if err != nil {
			return nil, fmt.Errorf("failed to get action from aptos role %s: %w", options.AptosRole, err)
		}
	}

	chainAccessor := cldfmcmsadapters.Wrap(env.BlockChains)

	return mcmschainwrappers.BuildInspector(&chainAccessor, mcmstypes.ChainSelector(chain), action,
		mcmstypes.ChainMetadata{})
}

func McmsInspectors(env cldf.Environment) (map[uint64]mcmssdk.Inspector, error) {
	chainsMetadata := map[mcmstypes.ChainSelector]mcmstypes.ChainMetadata{}
	for chainSelector := range env.BlockChains.All() {
		chainsMetadata[mcmstypes.ChainSelector(chainSelector)] = mcmstypes.ChainMetadata{}
	}

	chainAccessor := cldfmcmsadapters.Wrap(env.BlockChains)

	mcmsInspectors, err := mcmschainwrappers.BuildInspectors(&chainAccessor, chainsMetadata, mcmstypes.TimelockActionSchedule)
	if err != nil {
		return nil, fmt.Errorf("failed to build inspectors: %w", err)
	}

	inspectors := make(map[uint64]mcmssdk.Inspector, len(mcmsInspectors))
	for chainSelector, inspector := range mcmsInspectors {
		inspectors[uint64(chainSelector)] = inspector
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
	if action == "" {
		return mcmsaptossdk.TimelockRoleProposer, nil
	}
	return mcmsaptossdk.AptosRoleFromAction(action)
}
