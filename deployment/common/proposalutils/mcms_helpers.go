package proposalutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

func verboseDebug(lggr logger.Logger, event *owner_helpers.RBACTimelockCallScheduled) {
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	lggr.Debugw("scheduled", "event", string(b))
}

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

// MaybeLoadMCMSWithTimelockContracts looks for the addresses corresponding to
// contracts deployed with DeployMCMSWithTimelock and loads them into a
// MCMSWithTimelockState struct. If none of the contracts are found, the state struct will be nil.
// An error indicates:
// - Found but was unable to load a contract
// - It only found part of the bundle of contracts
// - If found more than one instance of a contract (we expect one bundle in the given addresses)
func MaybeLoadMCMSWithTimelockContracts(chain cldf_evm.Chain, addresses map[string]cldf.TypeAndVersion) (*MCMSWithTimelockContracts, error) {
	state := MCMSWithTimelockContracts{}
	// We expect one of each contract on the chain.
	timelock := cldf.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
	callProxy := cldf.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
	proposer := cldf.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
	canceller := cldf.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
	bypasser := cldf.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)

	// Convert map keys to a slice
	wantTypes := []cldf.TypeAndVersion{timelock, proposer, canceller, bypasser, callProxy}
	_, err := cldf.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", chain.Name(), err)
	}

	for address, tvStr := range addresses {
		switch {
		case tvStr.Type == timelock.Type && tvStr.Version.String() == timelock.Version.String():
			tl, err := owner_helpers.NewRBACTimelock(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.Timelock = tl
		case tvStr.Type == callProxy.Type && tvStr.Version.String() == callProxy.Version.String():
			cp, err := owner_helpers.NewCallProxy(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CallProxy = cp
		case tvStr.Type == proposer.Type && tvStr.Version.String() == proposer.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.ProposerMcm = mcms
		case tvStr.Type == bypasser.Type && tvStr.Version.String() == bypasser.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.BypasserMcm = mcms
		case tvStr.Type == canceller.Type && tvStr.Version.String() == canceller.Version.String():
			mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress(address), chain.Client)
			if err != nil {
				return nil, err
			}
			state.CancellerMcm = mcms
		}
	}
	return &state, nil
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

func McmsTimelockConverters(env cldf.Environment) (map[uint64]mcmssdk.TimelockConverter, error) {
	evmChains := env.BlockChains.EVMChains()
	solanaChains := env.BlockChains.SolanaChains()
	converters := make(map[uint64]mcmssdk.TimelockConverter, len(evmChains)+len(solanaChains))

	for _, chain := range evmChains {
		var err error
		converters[chain.Selector], err = McmsTimelockConverterForChain(chain.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms inspector for chain %s: %w", chain.String(), err)
		}
	}

	for _, chain := range solanaChains {
		var err error
		converters[chain.Selector], err = McmsTimelockConverterForChain(chain.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to get mcms inspector for chain %s: %w", chain.String(), err)
		}
	}

	return converters, nil
}

func McmsInspectorForChain(env cldf.Environment, chain uint64) (mcmssdk.Inspector, error) {
	chainFamily, err := mcmstypes.GetChainSelectorFamily(mcmstypes.ChainSelector(chain))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for chain %d: %w", chain, err)
	}

	switch chainFamily {
	case chain_selectors.FamilyEVM:
		return mcmsevmsdk.NewInspector(env.BlockChains.EVMChains()[chain].Client), nil
	case chain_selectors.FamilySolana:
		return mcmssolanasdk.NewInspector(env.BlockChains.SolanaChains()[chain].Client), nil
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
