package datastore

import "errors"

var ErrContractMetadataRecordNotFound = errors.New("contract metadata record not found")
var ErrContractMetadataRecordExists = errors.New("contract metadata record already exists")

var _ Record[ContractMetadataKey, ContractMetadataRecord] = ContractMetadataRecord{}

type ContractMetadataRecord struct {
	Chain    uint64
	Address  string
	Metadata string
}

func (r ContractMetadataRecord) Clone() ContractMetadataRecord {
	return ContractMetadataRecord{
		Chain:    r.Chain,
		Address:  r.Address,
		Metadata: r.Metadata,
	}
}

func (r ContractMetadataRecord) Key() ContractMetadataKey {
	return NewContractMetadataKey(r.Chain, r.Address)
}
