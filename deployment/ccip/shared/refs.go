package shared

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

// ErrRefReservation is returned when the refs a changeset intends to write cannot all be held in
// a datastore at once.
var ErrRefReservation = errors.New("cannot reserve datastore refs")

// ReserveRefs claims a datastore key for every ref a changeset is about to write and returns the
// store holding the reservations.
//
// It is meant to run before anything is deployed. A ref is keyed on
// (chain, type, version, qualifier) and none of that depends on the address, so the whole key set
// is known up front: a changeset reserves its keys, deploys, then fills in each address with an
// upsert into a key it already owns. That upsert cannot conflict, which is what removes the
// post-deploy datastore step that could otherwise fail once the contracts already exist on chain.
//
// Passing the same refs to Validate and to the deployment is the point: the check and the write
// are then the same computation rather than two that can drift apart.
func ReserveRefs(refs []datastore.AddressRef) (*datastore.MemoryDataStore, error) {
	// An empty qualifier is only valid for a chain singleton. datastore lookups that filter on
	// type and version read an empty qualifier as an omitted filter, so an unqualified member of
	// a multi-instance group would make every lookup for that group ambiguous. Count the group
	// sizes first: reacting while writing would be too late, because by then one arbitrary member
	// has already taken the key.
	type typeVersion struct {
		chain   uint64
		ct      datastore.ContractType
		version string
	}
	instances := make(map[typeVersion]int, len(refs))
	for _, ref := range refs {
		if ref.Version == nil {
			continue
		}
		instances[typeVersion{ref.ChainSelector, ref.Type, ref.Version.String()}]++
	}

	ds := datastore.NewMemoryDataStore()
	for _, ref := range refs {
		if ref.Version == nil {
			return nil, fmt.Errorf(
				"%w: %s on chain %d has no version, so it has no datastore key",
				ErrRefReservation, ref.Type, ref.ChainSelector)
		}

		if ref.Qualifier == "" && instances[typeVersion{ref.ChainSelector, ref.Type, ref.Version.String()}] > 1 {
			return nil, fmt.Errorf(
				"%w: chain %d has more than one %s %s and one of them has no qualifier; "+
					"an empty qualifier reads as 'any' when resolving, so every instance needs a name",
				ErrRefReservation, ref.ChainSelector, ref.Type, ref.Version)
		}

		if err := ds.Addresses().Add(ref); err != nil {
			return nil, fmt.Errorf(
				"%w: two refs claim (chain %d, %s %s, qualifier %q): %w",
				ErrRefReservation, ref.ChainSelector, ref.Type, ref.Version, ref.Qualifier, err)
		}
	}

	return ds, nil
}
