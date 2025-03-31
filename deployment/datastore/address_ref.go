package datastore

import (
	"errors"
	"maps"

	"github.com/Masterminds/semver/v3"
)

var (
	ErrAddressRefNotFound = errors.New("no such address ref can be found for the provided key")
	ErrAddressRefExists   = errors.New("an address ref with the supplied key already exists")
)

// TODO: ContractType is defined in many places in the codebase. Once it is moved
// to a common package, it is currently defined here to avoid circular dependencies
// between the datastore and the deployment packages. This should be removed once
// ContractType is mved to /deployment/common/types/types.go

// ContractType is a simple string type for identifying contract types.
type ContractType string

// String returns the string representation of the ContractType.
func (ct ContractType) String() string {
	return string(ct)
}

// AddressRef implements the Record interface
var _ Record[AddressRefKey, AddressRef] = AddressRef{}

type AddressRef struct {
	// Address is the address of the contract on the chain.
	Address string `json:"address"`
	// ChainSelector is the chain-selector of the chain where the contract is deployed.
	ChainSelector uint64 `json:"chainSelector"`
	// Labels are the labels associated with the contract.
	Labels LabelSet `json:"labels"`
	// Qualifier is an optional qualifier for the contract.
	Qualifier string `json:"qualifier"`
	// ContractType is a simple string type for identifying contract types.
	Type ContractType `json:"type"`
	// Version is the version of the contract.
	Version *semver.Version `json:"version"`
}

// Clone creates a copy of the AddressRefRecord.
func (r AddressRef) Clone() AddressRef {
	return AddressRef{
		ChainSelector: r.ChainSelector,
		Type:          r.Type,
		Version:       r.Version,
		Qualifier:     r.Qualifier,
		Address:       r.Address,
		Labels:        maps.Clone(r.Labels),
	}
}

// Key returns the AddressRefKey for the AddressRefRecord.
func (r AddressRef) Key() AddressRefKey {
	return NewAddressRefKey(r.ChainSelector, r.Type, r.Version, r.Qualifier)
}
