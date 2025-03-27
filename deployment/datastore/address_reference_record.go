package datastore

import (
	"errors"
	"maps"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	ErrAddressReferenceRecordNotFound = errors.New("no such reference record can be found for the provided key")
	ErrAddressReferenceRecordExists   = errors.New("a reference record with the supplied key already exists")
)

var _ Record[AddressReferenceKey, AddressReferenceRecord] = AddressReferenceRecord{}

type AddressReferenceRecord struct {
	Address   string
	Chain     uint64
	Labels    deployment.LabelSet
	Qualifier string
	Type      deployment.ContractType
	Version   *semver.Version
}

func (r AddressReferenceRecord) Clone() AddressReferenceRecord {
	return AddressReferenceRecord{
		Chain:     r.Chain,
		Type:      r.Type,
		Version:   r.Version,
		Qualifier: r.Qualifier,
		Address:   r.Address,
		Labels:    maps.Clone(r.Labels),
	}
}

func (r AddressReferenceRecord) Key() AddressReferenceKey {
	return NewAddressReferenceKey(r.Chain, r.Type, r.Version, r.Qualifier)
}
