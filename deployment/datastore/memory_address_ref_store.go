package datastore

import (
	"sync"
)

// AddressRefStore is an interface that represents an immutable view over a set
// of AddressRef records  identified by AddressRefKey.
type AddressRefStore interface {
	Store[AddressRefKey, AddressRef]
}

// MutableAddressRefStore is an interface that represents a mutable AddressRefStore
// of AddressRef records identified by AddressRefKey.
type MutableAddressRefStore interface {
	MutableStore[AddressRefKey, AddressRef]
}

// MemoryAddressRefStore is an in-memory implementation of the AddressRefStore and
// MutableAddressRefStore interfaces.
type MemoryAddressRefStore struct {
	Records []AddressRef `json:"records"`
	mu      sync.RWMutex
}

// MemoryAddressRefStore implements AddressRefStore interface.
var _ AddressRefStore = &MemoryAddressRefStore{}

// MemoryAddressRefStore implements MutableAddressRefStore interface.
var _ MutableAddressRefStore = &MemoryAddressRefStore{}

// NewMemoryAddressRefStore creates a new MemoryAddressRefStore instance.
func NewMemoryAddressRefStore() *MemoryAddressRefStore {
	return &MemoryAddressRefStore{Records: []AddressRef{}}
}

// Get returns the AddressRef for the provided key, or an error if no such record exists.
func (s *MemoryAddressRefStore) Get(key AddressRefKey) (AddressRef, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return AddressRef{}, ErrAddressRefNotFound
	}
	return s.Records[idx].Clone(), nil
}

// Fetch returns a copy of all AddressRef in the store.
func (s *MemoryAddressRefStore) Fetch() ([]AddressRef, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := []AddressRef{}
	for _, record := range s.Records {
		records = append(records, record.Clone())
	}
	return records, nil
}

// Filter returns a copy of all AddressRef in the store that pass all of the provided filters.
// Filters are applied in the order they are provided.
// If no filters are provided, all records are returned.
func (s *MemoryAddressRefStore) Filter(filters ...FilterFunc[AddressRefKey, AddressRef]) []AddressRef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := append([]AddressRef{}, s.Records...)
	for _, filter := range filters {
		records = filter(records)
	}

	return records
}

// indexOf returns the index of the record with the provided key, or -1 if no such record exists.
func (s *MemoryAddressRefStore) indexOf(key AddressRefKey) int {
	for idx, record := range s.Records {
		if record.Key().Equals(key) {
			return idx
		}
	}
	return -1
}

// Add inserts a new record into the store.
// If a record with the same key already exists, an error is returned.
func (s *MemoryAddressRefStore) Add(record AddressRef) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		return ErrAddressRefExists
	}
	s.Records = append(s.Records, record)
	return nil
}

// AddOrUpdate inserts a new record into the store if no record with the same key already exists.
// If a record with the same key already exists, it is updated.
func (s *MemoryAddressRefStore) AddOrUpdate(record AddressRef) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		s.Records[idx] = record
		return nil
	}
	s.Records = append(s.Records, record)
	return nil
}

// Update edits an existing record whose fields match the primary key elements of the supplied AddressRecord, with
// the non-primary-key values of the supplied AddressRecord.
// If no such record exists, an error is returned.
func (s *MemoryAddressRefStore) Update(record AddressRef) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		return ErrAddressRefNotFound
	}
	s.Records[idx] = record
	return nil
}

// Delete deletes record whose primary key elements match the supplied AddressRecord, returning an error if no
// such record exists to be deleted.
func (s *MemoryAddressRefStore) Delete(key AddressRefKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ErrAddressRefNotFound
	}
	s.Records = append(s.Records[:idx], s.Records[idx+1:]...)
	return nil
}
