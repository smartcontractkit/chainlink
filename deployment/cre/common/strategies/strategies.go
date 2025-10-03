package strategies

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

// TransactionStrategy interface for executing transactions with different strategies
type TransactionStrategy interface {
	Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) ([]mcmslib.TimelockProposal, error)
}

// SimpleTransaction executes a transaction directly without MCMS
type SimpleTransaction struct {
	Chain cldf_evm.Chain
}

func (s *SimpleTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) ([]mcmslib.TimelockProposal, error) {
	tx, err := callFn(s.Chain.DeployerKey)
	if err != nil {
		return nil, err
	}

	_, err = s.Chain.Confirm(tx)
	return []mcmslib.TimelockProposal{}, err
}

// MCMSTransaction executes a transaction through MCMS timelock
type MCMSTransaction struct {
	Config        *ocr3.MCMSConfig
	Description   string
	Address       common.Address
	ChainSel      uint64
	MCMSContracts *commonchangeset.MCMSWithTimelockState
	Env           cldf.Environment
}

func (m *MCMSTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) ([]mcmslib.TimelockProposal, error) {
	opts := cldf.SimTransactOpts()

	tx, err := callFn(opts)
	if err != nil {
		return nil, err
	}

	op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), tx.Data(), big.NewInt(0), "", nil)
	if err != nil {
		return nil, err
	}

	timelocksPerChain := map[uint64]string{
		m.ChainSel: m.MCMSContracts.Timelock.Address().Hex(),
	}
	proposerMCMSes := map[uint64]string{
		m.ChainSel: m.MCMSContracts.ProposerMcm.Address().Hex(),
	}
	inspector, err := proposalutils.McmsInspectorForChain(m.Env, m.ChainSel)
	if err != nil {
		return nil, err
	}
	inspectorPerChain := map[uint64]sdk.Inspector{
		m.ChainSel: inspector,
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		m.Env,
		timelocksPerChain,
		proposerMCMSes,
		inspectorPerChain,
		[]mcmstypes.BatchOperation{op},
		m.Description,
		proposalutils.TimelockConfig{MinDelay: m.Config.MinDuration},
	)
	if err != nil {
		return nil, err
	}

	return []mcmslib.TimelockProposal{*proposal}, nil
}

// CreateStrategy is a factory function to create the appropriate strategy based on configuration
func CreateStrategy(
	chain cldf_evm.Chain,
	env cldf.Environment,
	mcmsConfig *ocr3.MCMSConfig,
	mcmsContracts *commonchangeset.MCMSWithTimelockState,
	targetAddress common.Address,
	description string,
) (TransactionStrategy, error) {
	if mcmsConfig != nil {
		if mcmsConfig == nil {
			return nil, errors.New("MCMS config is required when mcmsConfig is not nil")
		}
		if mcmsContracts == nil {
			return nil, errors.New("MCMS contracts are required when mcmsConfig is not nil")
		}
		return &MCMSTransaction{
			Config:        mcmsConfig,
			Description:   description,
			Address:       targetAddress,
			ChainSel:      chain.Selector,
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
func GetMCMSContracts(e cldf.Environment, chainSelector uint64, qualifier string) (*commonchangeset.MCMSWithTimelockState, error) {
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
