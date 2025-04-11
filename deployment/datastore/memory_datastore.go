package datastore

import "errors"

// Merger is an interface that defines a method for merging two data stores.
type Merger[T any] interface {
	// Merge merges the given data into the current data store.
	// It should return an error if the merge fails.
	Merge(other T) error
}

// Sealer is an interface that defines a method for sealing a data store.
// A sealed data store cannot be modified further.
type Sealer[T any] interface {
	// Seal seals the data store, preventing further modifications.
	Seal() T
}

// BaseDataStore is an interface that defines the basic operations for a data store.
// It is parameterized by the type of address reference store and contract metadata store it uses.
type BaseDataStore[
	T Cloneable[T],
	U Cloneable[U],
	R AddressRefStore, CM ContractMetadataStore[T], EM EnvMetadataStore[U],
] interface {
	Addresses() R
	ContractMetadata() CM
	EnvMetadata() EM
}

// DataStore is an interface that defines the operations for a read-only data store.
type DataStore[T Cloneable[T], U Cloneable[U]] interface {
	BaseDataStore[T, U, AddressRefStore, ContractMetadataStore[T], EnvMetadataStore[U]]
}

// MutableDataStore is an interface that defines the operations for a mutable data store.
type MutableDataStore[T Cloneable[T], U Cloneable[U]] interface {
	Merger[DataStore[T, U]]
	Sealer[DataStore[T, U]]

	BaseDataStore[T, U, MutableAddressRefStore, MutableContractMetadataStore[T], MutableEnvMetadataStore[U]]
}

// MemoryDataStore is a concrete implementation of the MutableDataStore interface.
var _ MutableDataStore[DefaultMetadata, DefaultMetadata] = &MemoryDataStore[DefaultMetadata, DefaultMetadata]{}

type MemoryDataStore[CM Cloneable[CM], EM Cloneable[EM]] struct {
	AddressRefStore       *MemoryAddressRefStore           `json:"addressRefStore"`
	ContractMetadataStore *MemoryContractMetadataStore[CM] `json:"contractMetadataStore"`
	EnvMetadataStore      *MemoryEnvMetadataStore[EM]      `json:"envMetadataStore"`
}

// NewMemoryDataStore creates a new instance of MemoryDataStore.
// NOTE: The instance returned is mutable and can be modified.
func NewMemoryDataStore[CM Cloneable[CM], EM Cloneable[EM]]() *MemoryDataStore[CM, EM] {
	return &MemoryDataStore[CM, EM]{
		AddressRefStore:       NewMemoryAddressRefStore(),
		ContractMetadataStore: NewMemoryContractMetadataStore[CM](),
		EnvMetadataStore:      NewMemoryEnvMetadataStore[EM](),
	}
}

// Seal seals the MemoryDataStore, by returning a new instance of sealedMemoryDataStore.
func (s *MemoryDataStore[CM, EM]) Seal() DataStore[CM, EM] {
	return &sealedMemoryDataStore[CM, EM]{
		AddressRefStore:       s.AddressRefStore,
		ContractMetadataStore: s.ContractMetadataStore,
		EnvMetadataStore:      s.EnvMetadataStore,
	}
}

// Addresses returns the AddressRefStore of the MemoryDataStore.
func (s *MemoryDataStore[CM, EM]) Addresses() MutableAddressRefStore {
	return s.AddressRefStore
}

// ContractMetadata returns the ContractMetadataStore of the MemoryDataStore.
func (s *MemoryDataStore[CM, EM]) ContractMetadata() MutableContractMetadataStore[CM] {
	return s.ContractMetadataStore
}

// EnvMetadata returns the EnvMetadataStore of the MutableEnvMetadataStore.
func (s *MemoryDataStore[CM, EM]) EnvMetadata() MutableEnvMetadataStore[EM] {
	return s.EnvMetadataStore
}

// Merge merges the given mutable data store into the current MemoryDataStore.
func (s *MemoryDataStore[CM, EM]) Merge(other DataStore[CM, EM]) error {
	addressRefs, err := other.Addresses().Fetch()
	if err != nil {
		return err
	}

	for _, addressRef := range addressRefs {
		if err := s.AddressRefStore.AddOrUpdate(addressRef); err != nil {
			return err
		}
	}

	contractMetadataRecords, err := other.ContractMetadata().Fetch()
	if err != nil {
		return err
	}

	for _, record := range contractMetadataRecords {
		if err := s.ContractMetadataStore.AddOrUpdate(record); err != nil {
			return err
		}
	}

	envMetadata, err := other.EnvMetadata().Get()
	if err != nil {
		if errors.Is(err, ErrEnvMetadataNotSet) {
			// If the env metadata was not set in `other` data store, Get() will return
			// ErrEnvMetadataNotSet. In this case, we don't need to do anything because
			// since `other` doesn't contain any update to the env metadata, we can just
			// skip the env metadata update.
			return nil
		}
		return err
	}
	// If the env metadata was set, we need to update it in the current
	// data store.
	err = s.EnvMetadataStore.Set(envMetadata)
	if err != nil {
		return err
	}

	return nil
}

// SealedMemoryDataStore is a concrete implementation of the DataStore interface.
// It represents a sealed data store that cannot be modified further.
var _ DataStore[DefaultMetadata, DefaultMetadata] = &sealedMemoryDataStore[DefaultMetadata, DefaultMetadata]{}

type sealedMemoryDataStore[CM Cloneable[CM], EM Cloneable[EM]] struct {
	AddressRefStore       *MemoryAddressRefStore           `json:"addressRefStore"`
	ContractMetadataStore *MemoryContractMetadataStore[CM] `json:"contractMetadataStore"`
	EnvMetadataStore      *MemoryEnvMetadataStore[EM]      `json:"envMetadataStore"`
}

// Addresses returns the AddressRefStore of the sealedMemoryDataStore.
// It implements the BaseDataStore interface.
func (s *sealedMemoryDataStore[CM, EM]) Addresses() AddressRefStore {
	return s.AddressRefStore
}

// ContractMetadata returns the ContractMetadataStore of the sealedMemoryDataStore.
func (s *sealedMemoryDataStore[CM, EM]) ContractMetadata() ContractMetadataStore[CM] {
	return s.ContractMetadataStore
}

// EnvMetadata returns the EnvMetadataStore of the sealedMemoryDataStore.
func (s *sealedMemoryDataStore[CM, EM]) EnvMetadata() EnvMetadataStore[EM] {
	return s.EnvMetadataStore
}
