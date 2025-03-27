package datastore

import (
	"sync"
)

type ContractMetadataStore interface {
	Store[ContractMetadataKey, ContractMetadataRecord]
}

type MutableContractMetadataStore interface {
	MutableStore[ContractMetadataKey, ContractMetadataRecord]
}

var _ ContractMetadataStore = &InMemoryContractMetadataStore{}
var _ MutableContractMetadataStore = &InMemoryContractMetadataStore{}

type InMemoryContractMetadataStore struct {
	mu      sync.RWMutex
	records []ContractMetadataRecord
}

func NewInMemoryContractMetadataStore() InMemoryContractMetadataStore {
	return InMemoryContractMetadataStore{records: []ContractMetadataRecord{}}
}

func (s *InMemoryContractMetadataStore) indexOf(key ContractMetadataKey) int {
	for i, record := range s.records {
		if record.Key().Equals(key) {
			return i
		}
	}
	return -1
}

func (s *InMemoryContractMetadataStore) Get(key ContractMetadataKey) (ContractMetadataRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ContractMetadataRecord{}, ErrContractMetadataRecordNotFound
	}
	return s.records[idx].Clone(), nil
}

func (s *InMemoryContractMetadataStore) Fetch() ([]ContractMetadataRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := []ContractMetadataRecord{}
	for _, record := range s.records {
		records = append(records, record.Clone())
	}
	return records, nil
}

func (s *InMemoryContractMetadataStore) Filter(filters ...FilterFunc[ContractMetadataKey, ContractMetadataRecord]) []ContractMetadataRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := append([]ContractMetadataRecord{}, s.records...)
	for _, filter := range filters {
		records = filter(records)
	}
	return records
}

func (s *InMemoryContractMetadataStore) Add(record ContractMetadataRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		return ErrContractMetadataRecordExists
	}
	s.records = append(s.records, record)
	return nil
}

func (s *InMemoryContractMetadataStore) AddOrUpdate(record ContractMetadataRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		s.records = append(s.records, record)
		return nil
	}
	s.records[idx] = record
	return nil
}

func (s *InMemoryContractMetadataStore) Update(record ContractMetadataRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		return ErrContractMetadataRecordNotFound
	}
	s.records[idx] = record
	return nil
}

func (s *InMemoryContractMetadataStore) Delete(key ContractMetadataKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ErrContractMetadataRecordNotFound
	}
	s.records = append(s.records[:idx], s.records[idx+1:]...)
	return nil
}
