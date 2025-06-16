package seqs

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

// SeqPersistAddressDeps contains the dependencies for the SeqPersistAddress sequence.
type SeqPersistAddressDeps struct {
	AddrBook  cldf.AddressBook
	Datastore datastore.MutableDataStore
}

// SeqPersistAddressInput is the input to the SeqPersistAddress sequence.
type SeqPersistAddressInput struct {
	ChainSelector uint64   `json:"chainSelector"`
	Address       string   `json:"address"`
	Type          string   `json:"type"`
	Version       string   `json:"version"`
	Qualifier     string   `json:"qualifier"`
	Labels        []string `json:"labels"`
}

// SeqPersistAddressOutput is the output of the SeqPersistAddress sequence.
// It does not return any data, but it is included for consistency with other sequences.
type SeqPersistAddressOutput struct{}

// SeqPersistAddress is a sequence that persists the address of a deployed contract
// in the address book and datastore.
var SeqPersistAddress = operations.NewSequence(
	"seq-persist-address",
	semver.MustParse("1.0.0"),
	"Persist address of deployed contract in address book and datastore",
	func(b operations.Bundle, deps SeqPersistAddressDeps, in SeqPersistAddressInput) (SeqPersistAddressOutput, error) {
		out := SeqPersistAddressOutput{}

		// Store it in the legacy address book
		if _, err := operations.ExecuteOperation(b, ops.OpAddAddrBookRecord,
			ops.OpAddAddrBookRecordDeps{AddrBook: deps.AddrBook},
			ops.OpAddAddrBookRecordInput{
				ChainSelector: in.ChainSelector,
				Address:       in.Address,
				Type:          in.Type,
				Version:       in.Version,
				Labels:        in.Labels,
			},
		); err != nil {
			return out, err
		}

		// Store the address reference in the datastore
		if _, err := operations.ExecuteOperation(b, ops.OpAddDatastoreAddrRef,
			ops.OpAddDatastoreAddrRefDeps{Datastore: deps.Datastore},
			ops.OpAddDatastoreAddrRefInput{
				ChainSelector: in.ChainSelector,
				Address:       in.Address,
				Type:          in.Type,
				Version:       in.Version,
				Qualifier:     in.Qualifier,
				Labels:        in.Labels,
			},
		); err != nil {
			return out, err
		}

		return out, nil
	})
