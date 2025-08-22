package changeset

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type AcceptAllOwnershipRequest struct {
	ChainSelector uint64
	MinDelay      time.Duration
}

var _ cldf.ChangeSet[*AcceptAllOwnershipRequest] = AcceptAllOwnershipsProposal

// AcceptAllOwnershipsProposal creates a MCMS proposal to call accept ownership on all the Keystone contracts in the address book.
func AcceptAllOwnershipsProposal(e cldf.Environment, req *AcceptAllOwnershipRequest) (cldf.ChangesetOutput, error) {
	chainSelector := req.ChainSelector

	// Construct the configuration
	cfg := changeset.TransferToMCMSWithTimelockConfig{
		ContractsByChain: map[uint64][]common.Address{
			chainSelector: getTransferableContracts(e.DataStore.Addresses(), chainSelector),
		},
		MCMSConfig: proposalutils.TimelockConfig{MinDelay: req.MinDelay},
	}

	// Create and return the changeset
	return changeset.TransferToMCMSWithTimelockV2(e, cfg)
}

func getTransferableContracts(addressStore datastore.AddressRefStore, chainSelector uint64) []common.Address {
	var transferableContracts []common.Address

	addresses := addressStore.Filter(datastore.AddressRefByChainSelector(chainSelector))
	for _, addr := range addresses {
		isOCR3Capability := addr.Type == datastore.ContractType(OCR3Capability)
		isWorkflowRegistry := addr.Type == datastore.ContractType(WorkflowRegistry)
		isKeystoneForwarder := addr.Type == datastore.ContractType(KeystoneForwarder)
		isCapabilityRegistry := addr.Type == datastore.ContractType(CapabilitiesRegistry)

		if isCapabilityRegistry || isWorkflowRegistry || isKeystoneForwarder || isOCR3Capability {
			transferableContracts = append(transferableContracts, common.HexToAddress(addr.Address))
		}
	}

	return transferableContracts
}
