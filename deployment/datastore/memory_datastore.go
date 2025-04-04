package datastore

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
type BaseDataStore[T Cloneable[T], R AddressRefStore, M ContractMetadataStore[T]] interface {
	Addresses() R
	Metadata() M
}

// DataStore is an interface that defines the operations for a read-only data store.
type DataStore[T Cloneable[T]] interface {
	BaseDataStore[T, AddressRefStore, ContractMetadataStore[T]]
}

// MutableDataStore is an interface that defines the operations for a mutable data store.
type MutableDataStore[T Cloneable[T]] interface {
	Merger[DataStore[T]]
	Sealer[DataStore[T]]

	BaseDataStore[T, MutableAddressRefStore, MutableContractMetadataStore[T]]
}

// MemoryDataStore is a concrete implementation of the MutableDataStore interface.
var _ MutableDataStore[DefaultMetadata] = &MemoryDataStore[DefaultMetadata]{}

type MemoryDataStore[M Cloneable[M]] struct {
	AddressRefStore *MemoryAddressRefStore          `json:"addressRefStore"`
	MetadataStore   *MemoryContractMetadataStore[M] `json:"metadataStore"`
}

// NewMemoryDataStore creates a new instance of MemoryDataStore.
// NOTE: The instance returned is mutable and can be modified.
func NewMemoryDataStore[M Cloneable[M]]() *MemoryDataStore[M] {
	return &MemoryDataStore[M]{
		AddressRefStore: NewMemoryAddressRefStore(),
		MetadataStore:   NewMemoryContractMetadataStore[M](),
	}
}

// Seal seals the MemoryDataStore, by returning a new instance of sealedMemoryDataStore.
func (s *MemoryDataStore[M]) Seal() DataStore[M] {
	return &sealedMemoryDataStore[M]{
		AddressRefStore: s.AddressRefStore,
		MetadataStore:   s.MetadataStore,
	}
}

// Addresses returns the AddressRefStore of the MemoryDataStore.
func (s *MemoryDataStore[M]) Addresses() MutableAddressRefStore {
	return s.AddressRefStore
}

// Metadata returns the MetadataStore of the MemoryDataStore.
func (s *MemoryDataStore[M]) Metadata() MutableContractMetadataStore[M] {
	return s.MetadataStore
}

// Merge merges the given mutable data store into the current MemoryDataStore.
func (s *MemoryDataStore[M]) Merge(other DataStore[M]) error {
	addressRefs, err := other.Addresses().Fetch()
	if err != nil {
		return err
	}

	for _, addressRef := range addressRefs {
		if err := s.AddressRefStore.AddOrUpdate(addressRef); err != nil {
			return err
		}
	}

	metadataRecords, err := other.Metadata().Fetch()
	if err != nil {
		return err
	}

	for _, record := range metadataRecords {
		if err := s.MetadataStore.AddOrUpdate(record); err != nil {
			return err
		}
	}

	return nil
}

// SealedMemoryDataStore is a concrete implementation of the DataStore interface.
// It represents a sealed data store that cannot be modified further.
var _ DataStore[DefaultMetadata] = &sealedMemoryDataStore[DefaultMetadata]{}

type sealedMemoryDataStore[M Cloneable[M]] struct {
	AddressRefStore *MemoryAddressRefStore          `json:"addressRefStore"`
	MetadataStore   *MemoryContractMetadataStore[M] `json:"metadataStore"`
}

// Addresses returns the AddressRefStore of the sealedMemoryDataStore.
// It implements the BaseDataStore interface.
func (s *sealedMemoryDataStore[M]) Addresses() AddressRefStore {
	return s.AddressRefStore
}

// Metadata returns the MetadataStore of the sealedMemoryDataStore.
func (s *sealedMemoryDataStore[M]) Metadata() ContractMetadataStore[M] {
	return s.MetadataStore
}
