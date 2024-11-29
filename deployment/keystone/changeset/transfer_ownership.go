package changeset

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

// TransferAllOwnership transfers ownership of all keystone contracts in the address book to the existing timelock.
func TransferAllOwnership(e deployment.Environment, chainSelector uint64) (deployment.ChangesetOutput, error) {
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
	consumer, err := feedsConsumerFromAddrBook(e.ExistingAddresses, e.Chains[chainSelector])
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	cfg := changeset.TransferOwnershipConfig{
		OwnersPerChain: map[uint64]common.Address{
			chainSelector: timelock.Address(),
		},
		Contracts: map[uint64][]changeset.OwnershipTransferrer{
			chainSelector: {capReg, ocr3, forwarder, consumer},
		},
	}
	return changeset.NewTransferOwnershipChangeset(e, cfg)
}
