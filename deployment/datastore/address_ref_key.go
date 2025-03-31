package datastore

import (
	"github.com/Masterminds/semver/v3"
)

// AddressRefKey is an interface that represents a key for AddressRef records.
// It is used to uniquely identify a record in the AddressRefStore.
type AddressRefKey interface {
	Comparable[AddressRefKey]

	// ChainSelector returns the chain-selector selector of the chain where the contract is deployed.
	ChainSelector() uint64
	// Type returns the contract type of the contract.
	// This is a simple string type for identifying contract
	Type() ContractType
	// Version returns the semantic version of the contract.
	Version() *semver.Version
	// Qualifier returns the optional qualifier for the contract.
	// This can be used to differentiate between different references of the same contract.
	Qualifier() string
}

// addressRefKey implements the AddressRefKey interface.
var _ AddressRefKey = addressRefKey{}

// addressRefKey is a struct that implements the AddressRefKey interface.
// It is used to uniquely identify a record in the AddressRefStore.
type addressRefKey struct {
	chainSelector uint64
	contractType  ContractType
	version       *semver.Version
	qualifier     string
}

// CHainSelector returns the chain-selector selector of the chain where the contract is deployed.
func (a addressRefKey) ChainSelector() uint64 { return a.chainSelector }

// Type returns the contract type of the contract.
// This is a simple string type for identifying contract
func (a addressRefKey) Type() ContractType { return a.contractType }

// Version returns the semantic version of the contract.
func (a addressRefKey) Version() *semver.Version { return a.version }

// Qualifier returns the optional qualifier for the contract.
func (a addressRefKey) Qualifier() string { return a.qualifier }

// Equals returns true if the two AddressRefKey instances are equal, false otherwise.
func (a addressRefKey) Equals(other AddressRefKey) bool {
	return a.chainSelector == other.ChainSelector() &&
		a.contractType == other.Type() &&
		a.version.Equal(other.Version()) &&
		a.qualifier == other.Qualifier()
}

// NewAddressRefKey creates a new AddressRefKey instance.
func NewAddressRefKey(chainSelector uint64, contractType ContractType, version *semver.Version, qualifier string) AddressRefKey {
	return addressRefKey{
		chainSelector: chainSelector,
		contractType:  contractType,
		version:       version,
		qualifier:     qualifier,
	}
}
