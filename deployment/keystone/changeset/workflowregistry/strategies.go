package workflowregistry

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type strategy interface {
	Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (cldf.ChangesetOutput, error)
}

type simpleTransaction struct {
	chain cldf_evm.Chain
}

func (s *simpleTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (cldf.ChangesetOutput, error) {
	tx, err := callFn(s.chain.DeployerKey)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	_, err = s.chain.Confirm(tx)
	return cldf.ChangesetOutput{}, err
}

type mcmsTransaction struct {
	Config      *changeset.MCMSConfig
	Description string
	Address     common.Address
	ChainSel    uint64
	ContractSet *changeset.ContractSet
	Env         cldf.Environment
}

func (m *mcmsTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (cldf.ChangesetOutput, error) {
	opts := cldf.SimTransactOpts()

	tx, err := callFn(opts)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	op, err := proposalutils.BatchOperationForChain(m.ChainSel, m.Address.Hex(), tx.Data(), big.NewInt(0), "", nil)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	timelocksPerChain := map[uint64]string{
		m.ChainSel: m.ContractSet.Timelock.Address().Hex(),
	}
	proposerMCMSes := map[uint64]string{
		m.ChainSel: m.ContractSet.ProposerMcm.Address().Hex(),
	}
	inspector, err := proposalutils.McmsInspectorForChain(m.Env, m.ChainSel)
	if err != nil {
		return cldf.ChangesetOutput{}, err
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
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
	}, nil
}
