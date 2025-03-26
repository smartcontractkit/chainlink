package datastore

import (
	"sync"
)

// AddressReferenceStore is an interface that represents an immutable view over a set
// of AddressReferenceRecords identified by AddressReferenceKeys.
type AddressReferenceStore interface {
	Store[AddressReferenceKey, AddressReferenceRecord]
}

// MutableAddressReferenceStore is an interface that represents a mutable AddressReferenceStore
// of AddressReferenceRecords identified by AddressReferenceKeys.
type MutableAddressReferenceStore interface {
	MutableStore[AddressReferenceKey, AddressReferenceRecord]
}

// InMemoryAddressReferenceStore is an in-memory implementation of the AddressReferenceStore and
// MutableAddressReferenceStore interfaces.
var _ AddressReferenceStore = &InMemoryAddressReferenceStore{}
var _ MutableAddressReferenceStore = &InMemoryAddressReferenceStore{}

type InMemoryAddressReferenceStore struct {
	records []AddressReferenceRecord
	mu      sync.RWMutex
}

func NewInMemoryAddressReferenceStore() InMemoryAddressReferenceStore {
	return InMemoryAddressReferenceStore{records: []AddressReferenceRecord{}}
}

// Store interface methods implementation

// Get returns the AddressReferenceRecord for the provided key, or an error if no such record exists.
func (s *InMemoryAddressReferenceStore) Get(key AddressReferenceKey) (AddressReferenceRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return AddressReferenceRecord{}, ErrAddressReferenceRecordNotFound
	}
	return s.records[idx].Clone(), nil
}

// Fetch returns a copy of all AddressReferenceRecords in the store.
func (s *InMemoryAddressReferenceStore) Fetch() ([]AddressReferenceRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := []AddressReferenceRecord{}
	for _, record := range s.records {
		records = append(records, record.Clone())
	}
	return records, nil
}

// Filter returns a copy of all AddressReferenceRecords in the store that pass all of the provided filters.
// Filters are applied in the order they are provided.
// If no filters are provided, all records are returned.
func (s *InMemoryAddressReferenceStore) Filter(filters ...FilterFunc[AddressReferenceKey, AddressReferenceRecord]) []AddressReferenceRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := append([]AddressReferenceRecord{}, s.records...)
	for _, filter := range filters {
		records = filter(records)
	}

	return records
}

// MutableStore interface methods implementation

// indexOf returns the index of the record with the provided key, or -1 if no such record exists.
func (s *InMemoryAddressReferenceStore) indexOf(key AddressReferenceKey) int {
	for idx, record := range s.records {
		if record.Key().Equals(key) {
			return idx
		}
	}
	return -1
}

// Add inserts a new record into the store.
// If a record with the same key already exists, an error is returned.
func (s *InMemoryAddressReferenceStore) Add(record AddressReferenceRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		return ErrAddressReferenceRecordExists
	}
	s.records = append(s.records, record)
	return nil
}

// AddOrUpdate inserts a new record into the store if no record with the same key already exists.
// If a record with the same key already exists, it is updated.
func (s *InMemoryAddressReferenceStore) AddOrUpdate(record AddressReferenceRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		s.records[idx] = record
		return nil
	}
	s.records = append(s.records, record)
	return nil
}

// Update edits an existing record whose fields match the primary key elements of the supplied AddressRecord, with
// the non-primary-key values of the supplied AddressRecord.
// If no such record exists, an error is returned.
func (s *InMemoryAddressReferenceStore) Update(record AddressReferenceRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		return ErrAddressReferenceRecordNotFound
	}
	s.records[idx] = record
	return nil
}

// Delete deletes record whose primary key elements match the supplied AddressRecord, returning an error if no
// such record exists to be deleted.
func (s *InMemoryAddressReferenceStore) Delete(key AddressReferenceKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ErrAddressReferenceRecordNotFound
	}
	s.records = append(s.records[:idx], s.records[idx+1:]...)
	return nil
}
