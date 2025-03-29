package datastore

import "errors"

var ErrContractMetadataNotFound = errors.New("contract metadata record not found")
var ErrContractMetadataExists = errors.New("contract metadata record already exists")

var _ Record[ContractMetadataKey, ContractMetadata] = ContractMetadata{}

type ContractMetadata struct {
	ChainSelector uint64
	Address       string
	Metadata      string
}

func (r ContractMetadata) Clone() ContractMetadata {
	return ContractMetadata{
		ChainSelector: r.ChainSelector,
		Address:       r.Address,
		Metadata:      r.Metadata,
	}
}

func (r ContractMetadata) Key() ContractMetadataKey {
	return NewContractMetadataKey(r.ChainSelector, r.Address)
}
