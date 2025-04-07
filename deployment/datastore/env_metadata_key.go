package datastore

type EnvMetadataKey interface {
	Comparable[EnvMetadataKey]

	// Domain returns the domain to which the metadata belongs.
	Domain() string
	// Environment returns the environment to which the metadata belongs.
	Environment() string
}

// envMetadataKey implements the EnvMetadataKey interface.
var _ EnvMetadataKey = envMetadataKey{}

// envMetadataKey is a struct that implements the EnvMetadataKey interface.
type envMetadataKey struct {
	domain      string
	environment string
}

// Domain returns the domain to which the metadata belongs.
func (d envMetadataKey) Domain() string { return d.domain }

// Environment returns the environment to which the metadata belongs.
func (d envMetadataKey) Environment() string { return d.environment }

// Equals returns true if the two EnvMetadataKey instances are equal, false otherwise.
func (d envMetadataKey) Equals(other EnvMetadataKey) bool {
	return d.Domain() == other.Domain() &&
		d.Environment() == other.Environment()
}

// NewEnvMetadataKey creates a new EnvMetadataKey instance.
func NewEnvMetadataKey(domain string, environment string) EnvMetadataKey {
	return envMetadataKey{
		domain:      domain,
		environment: environment,
	}
}
