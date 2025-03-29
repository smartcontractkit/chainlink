package datastore

import (
	"sync"
)

type ContractMetadataStore interface {
	Store[ContractMetadataKey, ContractMetadata]
}

type MutableContractMetadataStore interface {
	MutableStore[ContractMetadataKey, ContractMetadata]
}

var _ ContractMetadataStore = &MemoryContractMetadataStore{}
var _ MutableContractMetadataStore = &MemoryContractMetadataStore{}

type MemoryContractMetadataStore struct {
	mu      sync.RWMutex
	records []ContractMetadata
}

func NewInMemoryContractMetadataStore() MemoryContractMetadataStore {
	return MemoryContractMetadataStore{records: []ContractMetadata{}}
}

func (s *MemoryContractMetadataStore) indexOf(key ContractMetadataKey) int {
	for i, record := range s.records {
		if record.Key().Equals(key) {
			return i
		}
	}
	return -1
}

func (s *MemoryContractMetadataStore) Get(key ContractMetadataKey) (ContractMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ContractMetadata{}, ErrContractMetadataNotFound
	}
	return s.records[idx].Clone(), nil
}

func (s *MemoryContractMetadataStore) Fetch() ([]ContractMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := []ContractMetadata{}
	for _, record := range s.records {
		records = append(records, record.Clone())
	}
	return records, nil
}

func (s *MemoryContractMetadataStore) Filter(filters ...FilterFunc[ContractMetadataKey, ContractMetadata]) []ContractMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := append([]ContractMetadata{}, s.records...)
	for _, filter := range filters {
		records = filter(records)
	}
	return records
}

func (s *MemoryContractMetadataStore) Add(record ContractMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx != -1 {
		return ErrContractMetadataExists
	}
	s.records = append(s.records, record)
	return nil
}

func (s *MemoryContractMetadataStore) AddOrUpdate(record ContractMetadata) error {
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

func (s *MemoryContractMetadataStore) Update(record ContractMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(record.Key())
	if idx == -1 {
		return ErrContractMetadataNotFound
	}
	s.records[idx] = record
	return nil
}

func (s *MemoryContractMetadataStore) Delete(key ContractMetadataKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.indexOf(key)
	if idx == -1 {
		return ErrContractMetadataNotFound
	}
	s.records = append(s.records[:idx], s.records[idx+1:]...)
	return nil
}
