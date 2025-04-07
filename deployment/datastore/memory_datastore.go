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
type BaseDataStore[T Cloneable[T], R AddressRefStore, CM ContractMetadataStore[T]] interface {
	Addresses() R
	ContractMetadata() CM
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

type MemoryDataStore[CM Cloneable[CM]] struct {
	AddressRefStore       *MemoryAddressRefStore           `json:"addressRefStore"`
	ContractMetadataStore *MemoryContractMetadataStore[CM] `json:"contractMetadataStore"`
}

// NewMemoryDataStore creates a new instance of MemoryDataStore.
// NOTE: The instance returned is mutable and can be modified.
func NewMemoryDataStore[CM Cloneable[CM]]() *MemoryDataStore[CM] {
	return &MemoryDataStore[CM]{
		AddressRefStore:       NewMemoryAddressRefStore(),
		ContractMetadataStore: NewMemoryContractMetadataStore[CM](),
	}
}

// Seal seals the MemoryDataStore, by returning a new instance of sealedMemoryDataStore.
func (s *MemoryDataStore[CM]) Seal() DataStore[CM] {
	return &sealedMemoryDataStore[CM]{
		AddressRefStore:       s.AddressRefStore,
		ContractMetadataStore: s.ContractMetadataStore,
	}
}

// Addresses returns the AddressRefStore of the MemoryDataStore.
func (s *MemoryDataStore[CM]) Addresses() MutableAddressRefStore {
	return s.AddressRefStore
}

// Metadata returns the MetadataStore of the MemoryDataStore.
func (s *MemoryDataStore[CM]) ContractMetadata() MutableContractMetadataStore[CM] {
	return s.ContractMetadataStore
}

// Merge merges the given mutable data store into the current MemoryDataStore.
func (s *MemoryDataStore[CM]) Merge(other DataStore[CM]) error {
	addressRefs, err := other.Addresses().Fetch()
	if err != nil {
		return err
	}

	for _, addressRef := range addressRefs {
		if err := s.AddressRefStore.AddOrUpdate(addressRef); err != nil {
			return err
		}
	}

	metadataRecords, err := other.ContractMetadata().Fetch()
	if err != nil {
		return err
	}

	for _, record := range metadataRecords {
		if err := s.ContractMetadataStore.AddOrUpdate(record); err != nil {
			return err
		}
	}

	return nil
}

// SealedMemoryDataStore is a concrete implementation of the DataStore interface.
// It represents a sealed data store that cannot be modified further.
var _ DataStore[DefaultMetadata] = &sealedMemoryDataStore[DefaultMetadata]{}

type sealedMemoryDataStore[CM Cloneable[CM]] struct {
	AddressRefStore       *MemoryAddressRefStore           `json:"addressRefStore"`
	ContractMetadataStore *MemoryContractMetadataStore[CM] `json:"contractMetadataStore"`
}

// Addresses returns the AddressRefStore of the sealedMemoryDataStore.
// It implements the BaseDataStore interface.
func (s *sealedMemoryDataStore[CM]) Addresses() AddressRefStore {
	return s.AddressRefStore
}

// Metadata returns the MetadataStore of the sealedMemoryDataStore.
func (s *sealedMemoryDataStore[CM]) ContractMetadata() ContractMetadataStore[CM] {
	return s.ContractMetadataStore
}
