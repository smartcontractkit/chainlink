package datastore

import (
	"sync"
)

// DomainMetadataStore is an interface that defines the methods for a store that manages domain metadata.
type DomainMetadataStore[DM Cloneable[DM]] interface {
	UnaryStore[DomainMetadataKey, DomainMetadata[DM]]
}

// MutableDomainMetadataStore is an interface that defines the methods for a mutable store that manages domain metadata.
type MutableDomainMetadataStore[DM Cloneable[DM]] interface {
	MutableUnaryStore[DomainMetadataKey, DomainMetadata[DM]]
}

// MemoryDomainMetadataStore is a concrete implementation of the DomainMetadataStore interface.
type MemoryDomainMetadataStore[DM Cloneable[DM]] struct {
	mu      sync.RWMutex
	Records []DomainMetadata[DM] `json:"records"`
}

// MemoryDomainMetadataStore implements DomainMetadataStore interface.
var _ DomainMetadataStore[DefaultMetadata] = &MemoryDomainMetadataStore[DefaultMetadata]{}

// MemoryDomainMetadataStore implements MutableDomainMetadataStore interface.
var _ MutableDomainMetadataStore[DefaultMetadata] = &MemoryDomainMetadataStore[DefaultMetadata]{}

// NewMemoryDomainMetadataStore creates a new MemoryDomainMetadataStore instance.
func NewMemoryDomainMetadataStore[M Cloneable[M]]() *MemoryDomainMetadataStore[M] {
	return &MemoryDomainMetadataStore[M]{Records: []DomainMetadata[M]{}}
}

// Get returns the DomainMetadata record if it exists, and a boolean indicating its existence
// and an error if any occurred.
// If no records exist, it returns an empty DomainMetadata and ErrDomainMetadataNotFound.
// If the record exists, it returns a clone of the record and a nil error.
func (s *MemoryDomainMetadataStore[M]) Get() (DomainMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Records) == 0 {
		return DomainMetadata[M]{}, ErrDomainMetadataNotFound
	}

	return s.Records[0].Clone(), nil
}

// Update replaces the existing record if present or adds it to the slice if empty.
// The record is always stored at index 0 of the slice to maintain a single record.
func (s *MemoryDomainMetadataStore[DM]) Update(record DomainMetadata[DM]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Records) == 0 {
		s.Records = append(s.Records, record)
	} else {
		s.Records[0] = record
	}

	return nil
}
