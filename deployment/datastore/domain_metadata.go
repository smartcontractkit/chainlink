package datastore

import "errors"

var ErrDomainMetadataNotFound = errors.New("no domain metadata record can be found")

// DomainMetadata implements the Record interface
var _ Record[DomainMetadataKey, DomainMetadata[DefaultMetadata]] = DomainMetadata[DefaultMetadata]{}

type DomainMetadata[M Cloneable[M]] struct {
	// Domain is the domain to which the metadata belongs.
	Domain string `json:"domain"`
	// Environment is the environment to which the metadata belongs.
	Environment string `json:"environment"`
	// Metadata is the metadata associated with the domain and environment.
	// It is a generic type that can be of any type that implements the Cloneable interface.
	Metadata M `json:"metadata"`
}

// Clone creates a copy of the DomainMetadata.
// The Metadata field is cloned using the Clone method of the Cloneable interface.
func (r DomainMetadata[M]) Clone() DomainMetadata[M] {
	return DomainMetadata[M]{
		Domain:      r.Domain,
		Environment: r.Environment,
		Metadata:    r.Metadata.Clone(),
	}
}

// Key returns the DomainMetadataKey for the DomainMetadata.
func (r DomainMetadata[M]) Key() DomainMetadataKey {
	return NewDomainMetadataKey(r.Domain, r.Environment)
}
