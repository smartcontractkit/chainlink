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
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

// TransactionStrategy interface for executing transactions with different strategies
type TransactionStrategy interface {
	BuildOperation(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (mcmstypes.BatchOperation, error)
	BuildProposal(operations []mcmstypes.BatchOperation) (mcmslib.TimelockProposal, error)
}

// SimpleTransaction executes a transaction directly without MCMS
type SimpleTransaction struct {
	Chain cldf_evm.Chain
}

func (s *SimpleTransaction) BuildOperation(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (mcmstypes.BatchOperation, error) {
	tx, err := callFn(s.Chain.DeployerKey)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}

	_, err = s.Chain.Confirm(tx)
	return mcmstypes.BatchOperation{}, err
}

func (s *SimpleTransaction) BuildProposal(_ []mcmstypes.BatchOperation) (mcmslib.TimelockProposal, error) {
	return mcmslib.TimelockProposal{}, nil
}

// MCMSTransaction executes a transaction through MCMS timelock
type MCMSTransaction struct {
	Config        *contracts.MCMSConfig
	Description   string
	Address       common.Address
	ChainSel      uint64
	MCMSContracts *commonchangeset.MCMSWithTimelockState
	Env           cldf.Environment
}

func (m *MCMSTransaction) BuildOperation(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (mcmstypes.BatchOperation, error) {
	opts := cldf.SimTransactOpts()

	tx, err := callFn(opts)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}

	op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), tx.Data(), big.NewInt(0), "", nil)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}

	return op, nil
}

func (m *MCMSTransaction) BuildProposal(operations []mcmstypes.BatchOperation) (mcmslib.TimelockProposal, error) {
	if m.MCMSContracts.Timelock == nil || m.MCMSContracts.ProposerMcm == nil {
		return mcmslib.TimelockProposal{}, errors.New("MCMS contracts are not properly initialized, missing Timelock or Proposer")
	}

	timelocksPerChain := map[uint64]string{
		m.ChainSel: m.MCMSContracts.Timelock.Address().Hex(),
	}
	proposerMCMSes := map[uint64]string{
		m.ChainSel: m.MCMSContracts.ProposerMcm.Address().Hex(),
	}
	inspector, err := proposalutils.McmsInspectorForChain(m.Env, m.ChainSel)
	if err != nil {
		return mcmslib.TimelockProposal{}, err
	}
	inspectorPerChain := map[uint64]sdk.Inspector{
		m.ChainSel: inspector,
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		m.Env,
		timelocksPerChain,
		proposerMCMSes,
		inspectorPerChain,
		operations,
		m.Description,
		proposalutils.TimelockConfig{MinDelay: m.Config.MinDuration},
	)
	if err != nil {
		return mcmslib.TimelockProposal{}, err
	}

	return *proposal, nil
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
