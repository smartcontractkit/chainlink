package datastore

import (
	"sync"
)

// ContractMetadataStore is an interface that represents an immutable view over a set
// of ContractMetadata records identified by ContractMetadataKey.
type ContractMetadataStore[M Cloneable[M]] interface {
	Store[ContractMetadataKey, ContractMetadata[M]]
}

// MutableContractMetadataStore is an interface that represents a mutable ContractMetadataStore
// of ContractMetadata records identified by ContractMetadataKey.
type MutableContractMetadataStore[M Cloneable[M]] interface {
	MutableStore[ContractMetadataKey, ContractMetadata[M]]
}

// MemoryContractMetadataStore is an in-memory implementation of the ContractMetadataStore and
// MutableContractMetadataStore interfaces.
type MemoryContractMetadataStore[M Cloneable[M]] struct {
	mu      sync.RWMutex
	Records []ContractMetadata[M] `json:"records"`
}

// MemoryContractMetadataStore implements ContractMetadataStore interface.
var _ ContractMetadataStore[DefaultMetadata] = &MemoryContractMetadataStore[DefaultMetadata]{}

// MemoryContractMetadataStore implements MutableContractMetadataStore interface.
var _ MutableContractMetadataStore[DefaultMetadata] = &MemoryContractMetadataStore[DefaultMetadata]{}

// NewMemoryContractMetadataStore creates a new MemoryContractMetadataStore instance.
// It is a generic function that takes a type parameter M which must implement the Cloneable interface.
func NewMemoryContractMetadataStore[M Cloneable[M]]() *MemoryContractMetadataStore[M] {
	return &MemoryContractMetadataStore[M]{Records: []ContractMetadata[M]{}}
}

// Get returns the ContractMetadata for the provided key, or an error if no such record exists.
func (s *MemoryContractMetadataStore[M]) Get(key ContractMetadataKey) (ContractMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ContractMetadata[M]{}, ErrContractMetadataNotFound
	}
	return s.Records[idx].Clone(), nil
}

// Fetch returns a copy of all ContractMetadata in the store.
func (s *MemoryContractMetadataStore[M]) Fetch() ([]ContractMetadata[M], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := []ContractMetadata[M]{}
	for _, record := range s.Records {
		records = append(records, record.Clone())
	}
	return records, nil
}

// Filter returns a copy of all ContractMetadata in the store that pass all of the provided filters.
// Filters are applied in the order they are provided.
// If no filters are provided, all records are returned.
func (s *MemoryContractMetadataStore[M]) Filter(filters ...FilterFunc[ContractMetadataKey, ContractMetadata[M]]) []ContractMetadata[M] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := append([]ContractMetadata[M]{}, s.Records...)
	for _, filter := range filters {
		records = filter(records)
	}
	return records
}

// indexOf returns the index of the record with the provided key, or -1 if no such record exists.
func (s *MemoryContractMetadataStore[M]) indexOf(key ContractMetadataKey) int {
	for i, record := range s.Records {
		if record.Key().Equals(key) {
			return i
		}
	}
	return -1
}

// Add inserts a new record into the store.
// If a record with the same key already exists, an error is returned.
func (s *MemoryContractMetadataStore[M]) Add(record ContractMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		return ErrContractMetadataExists
	}
	s.Records = append(s.Records, record)
	return nil
}

// AddOrUpdate inserts a new record into the store if no record with the same key already exists.
// If a record with the same key already exists, it is updated.
func (s *MemoryContractMetadataStore[M]) AddOrUpdate(record ContractMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		s.Records = append(s.Records, record)
		return nil
	}
	s.Records[idx] = record
	return nil
}

// Update edits an existing record whose fields match the primary key elements of the supplied ContractMetadata, with
// the non-primary-key values of the supplied ContractMetadata.
// If no such record exists, an error is returned.
func (s *MemoryContractMetadataStore[M]) Update(record ContractMetadata[M]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		return ErrContractMetadataNotFound
	}
	s.Records[idx] = record
	return nil
}

// Delete deletes an existing record whose primary key elements match the supplied ContractMetadata, returning an error if no
// such record exists.
func (s *MemoryContractMetadataStore[M]) Delete(key ContractMetadataKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ErrContractMetadataNotFound
	}
	s.Records = append(s.Records[:idx], s.Records[idx+1:]...)
	return nil
}
