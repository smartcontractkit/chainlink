package datastore

import "errors"

var ErrContractMetadataNotFound = errors.New("no contract metadata record can be found for the provided key")
var ErrContractMetadataExists = errors.New("a contract metadata record with the supplied key already exists")

// ContractMetadata implements the Record interface
var _ Record[ContractMetadataKey, ContractMetadata[DefaultMetadata]] = ContractMetadata[DefaultMetadata]{}

// DefaultMetadata is a default implementation of the custom contract metadata that can be used to store the
// contract metadata in the datastore as a simple string.
type DefaultMetadata string

// DefaultMetadata implements the Cloneable interface
func (d DefaultMetadata) Clone() DefaultMetadata { return d }

// ContractMetadata is a generic struct that holds the metadata for a contract on a specific chain.
// It implements the Record interface and is used to store contract metadata in the datastore.
// The metadata is generic and can be of any type that implements the Cloneable interface.
type ContractMetadata[M Cloneable[M]] struct {
	// Address is the address of the contract on the chain.
	Address string `json:"address"`
	// ChainSelector is the chain-selector of the chain where the contract is deployed.
	ChainSelector uint64 `json:"chainSelector"`
	// Metadata is the metadata associated with the contract.
	// It is a generic type that can be of any type that implements the Cloneable interface.
	Metadata M `json:"metadata"`
}

// Clone creates a copy of the ContractMetadata.
// The Metadata field is cloned using the Clone method of the Cloneable interface.
func (r ContractMetadata[M]) Clone() ContractMetadata[M] {
	return ContractMetadata[M]{
		ChainSelector: r.ChainSelector,
		Address:       r.Address,
		Metadata:      r.Metadata.Clone(),
	}
}

// Key returns the ContractMetadataKey for the ContractMetadata.
// It is used to uniquely identify the contract metadata in the datastore.
func (r ContractMetadata[M]) Key() ContractMetadataKey {
	return NewContractMetadataKey(r.ChainSelector, r.Address)
}
