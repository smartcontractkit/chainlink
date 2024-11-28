package changeset

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	ccipowner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// AcceptAllOwnershipsProposal creates a mcms proposal to call accept ownership on all the keystone contracts in the address book.
func AcceptAllOwnershipsProposal(e deployment.Environment, chainSelector uint64, minDelay time.Duration) (deployment.ChangesetOutput, error) {
	timelock, err := timelockFromAddrBook(e.ExistingAddresses, e.Chains[chainSelector])
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	capReg, err := capRegistryFromAddrBook(e.ExistingAddresses, e.Chains[chainSelector])
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	ocr3, err := ocr3FromAddrBook(e.ExistingAddresses, e.Chains[chainSelector])
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	forwarder, err := forwarderFromAddrBook(e.ExistingAddresses, e.Chains[chainSelector])
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	mcmsProposer, err := proposerFromAddrBook(e.ExistingAddresses, e.Chains[chainSelector])
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	cfg := changeset.AcceptOwnershipConfig{
		OwnersPerChain: map[uint64]common.Address{
			chainSelector: timelock.Address(),
		},
		ProposerMCMSes: map[uint64]*ccipowner.ManyChainMultiSig{
			chainSelector: mcmsProposer,
		},
		Contracts: map[uint64][]changeset.OwnershipAcceptor{
			chainSelector: {capReg, ocr3, forwarder},
		},
		MinDelay: minDelay,
	}
	return changeset.NewAcceptOwnershipChangeset(e, cfg)
}
