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

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

// MCMSTransaction executes a transaction through MCMS timelock
type MCMSTransaction struct {
	Env           cldf.Environment
	ChainSel      uint64
	Description   string
	Address       common.Address
	Config        *contracts.MCMSConfig
	MCMSContracts *commonchangeset.MCMSWithTimelockState
}

func (m *MCMSTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (*mcmstypes.BatchOperation, *types.Transaction, error) {
	opts := cldf.SimTransactOpts()

	tx, err := callFn(opts)
	if err != nil {
		return nil, nil, err
	}

	op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), tx.Data(), big.NewInt(0), "", nil)
	if err != nil {
		return nil, tx, err
	}

	return &op, tx, nil
}

func (m *MCMSTransaction) BuildProposal(operations []mcmstypes.BatchOperation) (*mcmslib.TimelockProposal, error) {
	if m.Config == nil || m.MCMSContracts == nil {
		return nil, errors.New("MCMS configuration or contracts are not provided")
	}

	if m.MCMSContracts.Timelock == nil || m.MCMSContracts.ProposerMcm == nil {
		return nil, errors.New("MCMS contracts are not properly initialized, missing Timelock or Proposer")
	}

	if len(operations) == 0 {
		return nil, errors.New("no operations provided to build proposal")
	}

	mcmContract, err := m.Config.MCMBasedOnAction(*m.MCMSContracts)
	if err != nil {
		return nil, fmt.Errorf("failed to get mcms contract by action '%s' for config %v : %w", m.Config.MCMSAction, m.Config, err)
	}

	timelocksPerChain := map[uint64]string{
		m.ChainSel: m.MCMSContracts.Timelock.Address().Hex(),
	}
	mcmsAddressesPerChain := map[uint64]string{
		m.ChainSel: mcmContract.Address().Hex(),
	}
	inspector, err := proposalutils.McmsInspectorForChain(m.Env, m.ChainSel)
	if err != nil {
		return nil, err
	}
	inspectorPerChain := map[uint64]sdk.Inspector{
		m.ChainSel: inspector,
	}

	return proposalutils.BuildProposalFromBatchesV2(
		m.Env,
		timelocksPerChain,
		mcmsAddressesPerChain,
		inspectorPerChain,
		operations,
		m.Description,
		*m.Config,
	)
}
