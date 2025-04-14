package datastore

import (
	"sync"
)

// EnvMetadataStore is an interface that defines the methods for a store that manages environment metadata.
type EnvMetadataStore[M Cloneable[M]] interface {
	UnaryStore[EnvMetadataKey, EnvMetadata[M]]
}

// MutableEnvMetadataStore is an interface that defines the methods for a mutable store that manages environment metadata.
type MutableEnvMetadataStore[M Cloneable[M]] interface {
	MutableUnaryStore[EnvMetadataKey, EnvMetadata[M]]
}

// MemoryEnvMetadataStore is a concrete implementation of the EnvMetadataStore interface.
type MemoryEnvMetadataStore[M Cloneable[M]] struct {
	mu     sync.RWMutex
	Record *EnvMetadata[M] `json:"record"`
}

// MemoryEnvMetadataStore implements EnvMetadataStore interface.
var _ EnvMetadataStore[DefaultMetadata] = &MemoryEnvMetadataStore[DefaultMetadata]{}

// MemoryEnvMetadataStore implements MutableEnvMetadataStore interface.
var _ MutableEnvMetadataStore[DefaultMetadata] = &MemoryEnvMetadataStore[DefaultMetadata]{}

// NewMemoryEnvMetadataStore creates a new MemoryEnvMetadataStore instance.
func NewMemoryEnvMetadataStore[M Cloneable[M]]() *MemoryEnvMetadataStore[M] {
	return &MemoryEnvMetadataStore[M]{Record: nil}
}

// Get returns a copy of the stored EnvMetadata record if it exists or an error if any occurred.
// If no record exist, it returns an empty EnvMetadata and ErrEnvMetadataNotSet.
// If the record exists, it returns a copy of the record and a nil error.
func (s *MemoryEnvMetadataStore[M]) Get() (EnvMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Record == nil {
		return EnvMetadata[M]{}, ErrEnvMetadataNotSet
	}

	return s.Record.Clone(), nil
}

// Set sets the EnvMetadata record in the store. If the record already exists, it will be replaced.
func (s *MemoryEnvMetadataStore[M]) Set(record EnvMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Record = &record

	return nil
}
