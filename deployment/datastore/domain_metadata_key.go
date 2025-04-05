package datastore

type DomainMetadataKey interface {
	Comparable[DomainMetadataKey]

	// Domain returns the domain to which the metadata belongs.
	Domain() string
	// Environment returns the environment to which the metadata belongs.
	Environment() string
}

// domainMetadataKey implements the DomainMetadataKey interface.
var _ DomainMetadataKey = domainMetadataKey{}

// domainMetadataKey is a struct that implements the DomainMetadataKey interface.
type domainMetadataKey struct {
	domain      string
	environment string
}

// Domain returns the domain to which the metadata belongs.
func (d domainMetadataKey) Domain() string { return d.domain }

// Environment returns the environment to which the metadata belongs.
func (d domainMetadataKey) Environment() string { return d.environment }

// Equals returns true if the two DomainMetadataKey instances are equal, false otherwise.
func (d domainMetadataKey) Equals(other DomainMetadataKey) bool {
	return d.domain == other.Domain() &&
		d.environment == other.Environment()
}

// NewDomainMetadataKey creates a new DomainMetadataKey instance.
func NewDomainMetadataKey(domain string, environment string) DomainMetadataKey {
	return domainMetadataKey{
		domain:      domain,
		environment: environment,
	}
}
