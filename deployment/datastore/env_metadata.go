package datastore

import "errors"

var ErrEnvMetadataNotSet = errors.New("no environment metadata set")

// EnvMetadata implements the Record interface
var _ Record[EnvMetadataKey, EnvMetadata[DefaultMetadata]] = EnvMetadata[DefaultMetadata]{}

type EnvMetadata[M Cloneable[M]] struct {
	// Domain is the domain to which the metadata belongs.
	Domain string `json:"domain"`
	// Environment is the environment to which the metadata belongs.
	Environment string `json:"environment"`
	// Metadata is the metadata associated with the domain and environment.
	// It is a generic type that can be of any type that implements the Cloneable interface.
	Metadata M `json:"metadata"`
}

// Clone creates a copy of the EnvMetadata.
// The Metadata field is cloned using the Clone method of the Cloneable interface.
func (r EnvMetadata[M]) Clone() EnvMetadata[M] {
	return EnvMetadata[M]{
		Domain:      r.Domain,
		Environment: r.Environment,
		Metadata:    r.Metadata.Clone(),
	}
}

// Key returns the EnvMetadataKey for the EnvMetadata.
func (r EnvMetadata[M]) Key() EnvMetadataKey {
	return NewEnvMetadataKey(r.Domain, r.Environment)
}
