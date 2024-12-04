package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

func toOwnershipTransferrer[T changeset.OwnershipTransferrer](items []T) []changeset.OwnershipTransferrer {
	ownershipAcceptors := make([]changeset.OwnershipTransferrer, len(items))
	for i, item := range items {
		ownershipAcceptors[i] = item
	}
	return ownershipAcceptors
}

type TransferAllOwnershipRequest struct {
	ChainSelector uint64
}

var _ deployment.ChangeSet[*TransferAllOwnershipRequest] = TransferAllOwnership

// TransferAllOwnership transfers ownership of all Keystone contracts in the address book to the existing timelock.
func TransferAllOwnership(e deployment.Environment, req *TransferAllOwnershipRequest) (deployment.ChangesetOutput, error) {
	chainSelector := req.ChainSelector
	chain := e.Chains[chainSelector]
	addrBook := e.ExistingAddresses

	r, err := kslib.GetContractSets(e.Logger, &kslib.GetContractSetsRequest{
		Chains: map[uint64]deployment.Chain{
			req.ChainSelector: chain,
		},
		AddressBook: addrBook,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}
	contracts := r.ContractSets[chainSelector]
	timelock := contracts.Timelock

	if timelock == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelocks: %w", err)
	}

	// Initialize the Contracts slice
	ownershipTransferrers := contracts.OwnershipTransferrers()

	// Construct the configuration
	cfg := changeset.TransferOwnershipConfig{
		OwnersPerChain: map[uint64]common.Address{
			// Assuming there is only one timelock per chain.
			chainSelector: timelock.Address(),
		},
		Contracts: map[uint64][]changeset.OwnershipTransferrer{
			chainSelector: ownershipTransferrers,
		},
	}

	// Create and return the changeset
	return changeset.NewTransferOwnershipChangeset(e, cfg)
}
