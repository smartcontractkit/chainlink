package changeset

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

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
	var addrsToTransfer []common.Address
	for _, consumer := range consumers {
		addrsToTransfer = append(addrsToTransfer, consumer.Address())
	}
	for _, o := range ocr3 {
		addrsToTransfer = append(addrsToTransfer, o.Address())
	}
	for _, f := range forwarders {
		addrsToTransfer = append(addrsToTransfer, f.Address())
	}
	for _, c := range capRegs {
		addrsToTransfer = append(addrsToTransfer, c.Address())
	}
	// Construct the configuration
	cfg := changeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: map[uint64][]common.Address{
			chainSelector: addrsToTransfer,
		},
		MinDelay: minDelay,
	}

	// Create and return the changeset
	return changeset.TransferToMCMSWithTimelock(e, cfg)
}
