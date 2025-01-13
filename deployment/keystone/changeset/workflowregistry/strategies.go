package workflowregistry

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type strategy interface {
	Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (deployment.ChangesetOutput, error)
}

type simpleTransaction struct {
	chain deployment.Chain
}

func (s *simpleTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (deployment.ChangesetOutput, error) {
	tx, err := callFn(s.chain.DeployerKey)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	_, err = s.chain.Confirm(tx)
	return deployment.ChangesetOutput{}, err
}

type mcmsTransaction struct {
	Config      *changeset.MCMSConfig
	Description string
	Address     common.Address
	ChainSel    uint64
	ContractSet *changeset.ContractSet
}

func (m *mcmsTransaction) Apply(callFn func(opts *bind.TransactOpts) (*types.Transaction, error)) (deployment.ChangesetOutput, error) {
	opts := deployment.SimTransactOpts()

	tx, err := callFn(opts)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	op := timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(m.ChainSel),
		Batch: []mcms.Operation{
			{
				Data:  tx.Data(),
				To:    m.Address,
				Value: big.NewInt(0),
			},
		},
	}

	timelocksPerChain := map[uint64]common.Address{
		m.ChainSel: m.ContractSet.Timelock.Address(),
	}
	proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
		m.ChainSel: m.ContractSet.ProposerMcm,
	}

	proposal, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		[]timelock.BatchChainOperation{op},
		m.Description,
		m.Config.MinDuration,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}
