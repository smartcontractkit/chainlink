package datastore

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/deployment"
)

type AddressReferenceKey interface {
	Comparable[AddressReferenceKey]

	Chain() uint64
	Type() deployment.ContractType
	Version() *semver.Version
	Qualifier() string
}

type addressReferenceKey struct {
	chain        uint64
	contractType deployment.ContractType
	version      *semver.Version
	qualifier    string
}

func (a addressReferenceKey) Chain() uint64                 { return a.chain }
func (a addressReferenceKey) Type() deployment.ContractType { return a.contractType }
func (a addressReferenceKey) Version() *semver.Version      { return a.version }
func (a addressReferenceKey) Qualifier() string             { return a.qualifier }

func (a addressReferenceKey) Equals(other AddressReferenceKey) bool {
	return a.chain == other.Chain() &&
		a.contractType == other.Type() &&
		a.version.Equal(other.Version()) &&
		a.qualifier == other.Qualifier()
}

func NewAddressReferenceKey(chain uint64, contractType deployment.ContractType, version *semver.Version, qualifier string) AddressReferenceKey {
	return addressReferenceKey{
		chain:        chain,
		contractType: contractType,
		version:      version,
		qualifier:    qualifier,
	}
}
