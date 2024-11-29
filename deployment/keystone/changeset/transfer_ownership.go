package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// TransferAllOwnership transfers ownership of all Keystone contracts in the address book to the existing timelock.
func TransferAllOwnership(e deployment.Environment, chainSelector uint64) (deployment.ChangesetOutput, error) {
	chain := e.Chains[chainSelector]
	addrBook := e.ExistingAddresses

	// Fetch timelocks for the specified chain.
	timelocks, err := timelocksFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch timelocks: %w", err)
	}
	if len(timelocks) == 0 {
		return deployment.ChangesetOutput{}, fmt.Errorf("no timelocks found for chain %d", chainSelector)
	}

	// Fetch contracts from the address book.
	capRegs, err := capRegistriesFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch capabilities registries: %w", err)
	}

	ocr3s, err := ocr3FromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch OCR3 capabilities: %w", err)
	}

	forwarders, err := forwardersFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch forwarders: %w", err)
	}

	consumers, err := feedsConsumersFromAddrBook(addrBook, chain)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to fetch feeds consumers: %w", err)
	}

	// Initialize the Contracts slice
	var ownershipTransferrers []changeset.OwnershipTransferrer

	// Append CapabilitiesRegistry contracts
	for _, capReg := range capRegs {
		ownershipTransferrers = append(ownershipTransferrers, capReg)
	}

	// Append OCR3Capability contracts
	for _, ocr := range ocr3s {
		ownershipTransferrers = append(ownershipTransferrers, ocr)
	}

	// Append KeystoneForwarder contracts
	for _, forwarder := range forwarders {
		ownershipTransferrers = append(ownershipTransferrers, forwarder)
	}

	// Append FeedsConsumer contracts
	for _, consumer := range consumers {
		ownershipTransferrers = append(ownershipTransferrers, consumer)
	}

	// Construct the configuration
	cfg := changeset.TransferOwnershipConfig{
		OwnersPerChain: map[uint64]common.Address{
			// Assuming there is only one timelock per chain.
			chainSelector: timelocks[0].Address(),
		},
		Contracts: map[uint64][]changeset.OwnershipTransferrer{
			chainSelector: ownershipTransferrers,
		},
	}

	// Create and return the changeset
	return changeset.NewTransferOwnershipChangeset(e, cfg)
}
