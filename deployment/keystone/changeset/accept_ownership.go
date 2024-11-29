package changeset

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	ccipowner "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// AcceptAllOwnershipsProposal creates a MCMS proposal to call accept ownership on all the Keystone contracts in the address book.
func AcceptAllOwnershipsProposal(e deployment.Environment, chainSelector uint64, minDelay time.Duration) (deployment.ChangesetOutput, error) {
	chain := e.Chains[chainSelector]
	addrBook := e.ExistingAddresses

	// Fetch contracts from the address book.
	timelocks, err := timelocksFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	capRegs, err := capRegistriesFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	ocr3, err := ocr3FromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	forwarders, err := forwardersFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	consumers, err := feedsConsumersFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	mcmsProposers, err := proposersFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Initialize the Contracts slice
	var ownershipAcceptors []changeset.OwnershipAcceptor

	// Append CapabilitiesRegistry contracts
	for _, capReg := range capRegs {
		ownershipAcceptors = append(ownershipAcceptors, capReg)
	}

	// Append OCR3Capability contracts
	for _, ocr := range ocr3 {
		ownershipAcceptors = append(ownershipAcceptors, ocr)
	}

	// Append KeystoneForwarder contracts
	for _, forwarder := range forwarders {
		ownershipAcceptors = append(ownershipAcceptors, forwarder)
	}

	// Append FeedsConsumer contracts
	for _, consumer := range consumers {
		ownershipAcceptors = append(ownershipAcceptors, consumer)
	}

	// Construct the configuration
	cfg := changeset.AcceptOwnershipConfig{
		OwnersPerChain: map[uint64]common.Address{
			// Assuming there is only one timelock per chain.
			chainSelector: timelocks[0].Address(),
		},
		ProposerMCMSes: map[uint64]*ccipowner.ManyChainMultiSig{
			// Assuming there is only one MCMS proposer per chain.
			chainSelector: mcmsProposers[0],
		},
		Contracts: map[uint64][]changeset.OwnershipAcceptor{
			chainSelector: ownershipAcceptors,
		},
		MinDelay: minDelay,
	}

	// Create and return the changeset
	return changeset.NewAcceptOwnershipChangeset(e, cfg)
}
