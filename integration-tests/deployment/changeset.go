package deployment

import (
	owner_wrappers "github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"
)

// TODO: Move to real MCM structs once available.
type Proposal struct {
	// keccak256(abi.encode(root, validUntil)) is what is signed by MCMS
	// signers.
	ValidUntil uint32
	// Leaves are the items in the proposal.
	// Uses these to generate the root as well as display whats in the root.
	// These Ops may be destined for distinct chains.
	Ops []owner_wrappers.ManyChainMultiSigOp
}

func (p Proposal) String() string {
	// TODO
	return ""
}

// Services as input to CI/Async tasks
type ChangesetOutput struct {
	JobSpecs    map[string][]string
	Proposals   []Proposal
	AddressBook AddressBook
}
