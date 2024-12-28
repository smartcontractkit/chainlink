package changeset

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

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

	// Construct the configuration
	cfg := changeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: map[uint64][]common.Address{
			chainSelector: contracts.TransferableContracts(),
		},
		MinDelay: minDelay,
	}

	// Create and return the changeset
	return changeset.TransferToMCMSWithTimelock(e, cfg)
}
