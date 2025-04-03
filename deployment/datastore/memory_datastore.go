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
// It is parameterized by the type of address reference store it uses.
type BaseDataStore[R AddressRefStore] interface {
	Addresses() R
}

// DataStore is an interface that defines the operations for a read-only data store.
type DataStore interface {
	BaseDataStore[AddressRefStore]
}

// MutableDataStore is an interface that defines the operations for a mutable data store.
type MutableDataStore interface {
	Merger[DataStore]
	Sealer[DataStore]

	BaseDataStore[MutableAddressRefStore]
}

// MemoryDataStore is a concrete implementation of the MutableDataStore interface.
var _ MutableDataStore = &MemoryDataStore{}

type MemoryDataStore struct {
	AddressRefStore *MemoryAddressRefStore `json:"addressRefStore"`
}

// NewMemoryDataStore creates a new instance of MemoryDataStore.
// NOTE: The instance returned is mutable and can be modified.
func NewMemoryDataStore() *MemoryDataStore {
	return &MemoryDataStore{
		AddressRefStore: NewMemoryAddressRefStore(),
	}
}

// Seal seals the MemoryDataStore, by returning a new instance of sealedMemoryDataStore.
func (s *MemoryDataStore) Seal() DataStore {
	return &sealedMemoryDataStore{AddressRefStore: s.AddressRefStore}
}

// Addresses returns the AddressRefStore of the MemoryDataStore.
func (s *MemoryDataStore) Addresses() MutableAddressRefStore {
	return s.AddressRefStore
}

// Merge merges the given mutable data store into the current MemoryDataStore.
func (s *MemoryDataStore) Merge(other DataStore) error {
	addressRefs, err := other.Addresses().Fetch()
	if err != nil {
		return err
	}

	for _, addressRef := range addressRefs {
		if err := s.AddressRefStore.AddOrUpdate(addressRef); err != nil {
			return err
		}
	}

	return nil
}

// SealedMemoryDataStore is a concrete implementation of the DataStore interface.
// It represents a sealed data store that cannot be modified further.
var _ DataStore = &sealedMemoryDataStore{}

type sealedMemoryDataStore struct {
	AddressRefStore *MemoryAddressRefStore `json:"addressRefStore"`
}

// Addresses returns the AddressRefStore of the sealedMemoryDataStore.
// It implements the BaseDataStore interface.
func (s *sealedMemoryDataStore) Addresses() AddressRefStore {
	return s.AddressRefStore
}
