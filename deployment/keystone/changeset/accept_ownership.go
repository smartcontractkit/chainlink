package changeset

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	ccipowner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func toOwnershipAcceptors[T changeset.OwnershipAcceptor](items []T) []changeset.OwnershipAcceptor {
	ownershipAcceptors := make([]changeset.OwnershipAcceptor, len(items))
	for i, item := range items {
		ownershipAcceptors[i] = item
	}
	return ownershipAcceptors
}

type AcceptAllOwnershipRequest struct {
	ChainSelector uint64
	MinDelay      time.Duration
}

var _ deployment.ChangeSet[*AcceptAllOwnershipRequest] = AcceptAllOwnershipsProposal

// AcceptAllOwnershipsProposal creates a MCMS proposal to call accept ownership on all the Keystone contracts in the address book.
func AcceptAllOwnershipsProposal(e deployment.Environment, req *AcceptAllOwnershipRequest) (deployment.ChangesetOutput, error) {
	chainSelector := req.ChainSelector
	minDelay := req.MinDelay
	chain := e.Chains[chainSelector]
	addrBook := e.ExistingAddresses

	r, err := kslib.GetContractSets(e.Logger, &kslib.GetContractSetsRequest{
		Chains: map[uint64]deployment.Chain{
			req.ChainSelector: chain,
		},
		AddressBook: addrBook,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	contracts := r.ContractSets[chainSelector]
	ownershipAcceptors := contracts.OwnershipAcceptors()
	// Construct the configuration
	cfg := changeset.AcceptOwnershipConfig{
		OwnersPerChain: map[uint64]common.Address{
			chainSelector: contracts.Timelock.Address(),
		},
		ProposerMCMSes: map[uint64]*ccipowner.ManyChainMultiSig{
			chainSelector: contracts.ProposerMcm,
		},
		Contracts: map[uint64][]changeset.OwnershipAcceptor{
			chainSelector: ownershipAcceptors,
		},
		MinDelay: minDelay,
	}

	// Create and return the changeset
	return changeset.NewAcceptOwnershipChangeset(e, cfg)
}
