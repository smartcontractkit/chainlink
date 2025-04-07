package datastore

import (
	"sync"
)

// EnvMetadataStore is an interface that defines the methods for a store that manages environment metadata.
type EnvMetadataStore[DM Cloneable[DM]] interface {
	UnaryStore[EnvMetadataKey, EnvMetadata[DM]]
}

// MutableEnvMetadataStore is an interface that defines the methods for a mutable store that manages environment metadata.
type MutableEnvMetadataStore[DM Cloneable[DM]] interface {
	MutableUnaryStore[EnvMetadataKey, EnvMetadata[DM]]
}

// MemoryEnvMetadataStore is a concrete implementation of the EnvMetadataStore interface.
type MemoryEnvMetadataStore[DM Cloneable[DM]] struct {
	mu      sync.RWMutex
	Records []EnvMetadata[DM] `json:"records"`
}

// MemoryEnvMetadataStore implements EnvMetadataStore interface.
var _ EnvMetadataStore[DefaultMetadata] = &MemoryEnvMetadataStore[DefaultMetadata]{}

// MemoryEnvMetadataStore implements MutableEnvMetadataStore interface.
var _ MutableEnvMetadataStore[DefaultMetadata] = &MemoryEnvMetadataStore[DefaultMetadata]{}

// NewMemoryEnvMetadataStore creates a new MemoryEnvMetadataStore instance.
func NewMemoryEnvMetadataStore[M Cloneable[M]]() *MemoryEnvMetadataStore[M] {
	return &MemoryEnvMetadataStore[M]{Records: []EnvMetadata[M]{}}
}

// Get returns the EnvMetadata record if it exists, and a boolean indicating its existence
// and an error if any occurred.
// If no records exist, it returns an empty EnvMetadata and ErrEnvMetadataNotFound.
// If the record exists, it returns a clone of the record and a nil error.
func (s *MemoryEnvMetadataStore[M]) Get() (EnvMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Records) == 0 {
		return EnvMetadata[M]{}, ErrEnvMetadataNotFound
	}

	return s.Records[0].Clone(), nil
}

// Update replaces the existing record if present or adds it to the slice if empty.
// The record is always stored at index 0 of the slice to maintain a single record.
func (s *MemoryEnvMetadataStore[DM]) Update(record EnvMetadata[DM]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Records) == 0 {
		s.Records = append(s.Records, record)
	} else {
		s.Records[0] = record
	}

	return nil
}
