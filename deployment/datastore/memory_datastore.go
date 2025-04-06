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
type BaseDataStore[
	T Cloneable[T],
	U Cloneable[U],
	R AddressRefStore, CM ContractMetadataStore[T], DM DomainMetadataStore[U],
] interface {
	Addresses() R
	ContractMetadata() CM
	DomainMetadata() DM
}

// DataStore is an interface that defines the operations for a read-only data store.
type DataStore[T Cloneable[T], U Cloneable[U]] interface {
	BaseDataStore[T, U, AddressRefStore, ContractMetadataStore[T], DomainMetadataStore[U]]
}

// MutableDataStore is an interface that defines the operations for a mutable data store.
type MutableDataStore[T Cloneable[T], U Cloneable[U]] interface {
	Merger[DataStore[T, U]]
	Sealer[DataStore[T, U]]

	BaseDataStore[T, U, MutableAddressRefStore, MutableContractMetadataStore[T], DomainMetadataStore[U]]
}

// MemoryDataStore is a concrete implementation of the MutableDataStore interface.
var _ MutableDataStore[DefaultMetadata, DefaultMetadata] = &MemoryDataStore[DefaultMetadata, DefaultMetadata]{}

type MemoryDataStore[CM Cloneable[CM], DM Cloneable[DM]] struct {
	AddressRefStore       *MemoryAddressRefStore           `json:"addressRefStore"`
	ContractMetadataStore *MemoryContractMetadataStore[CM] `json:"contractMetadataStore"`
	DomainMetadataStore   *MemoryDomainMetadataStore[DM]   `json:"domainMetadataStore"`
}

// NewMemoryDataStore creates a new instance of MemoryDataStore.
// NOTE: The instance returned is mutable and can be modified.
func NewMemoryDataStore[CM Cloneable[CM], DM Cloneable[DM]]() *MemoryDataStore[CM, DM] {
	return &MemoryDataStore[CM, DM]{
		AddressRefStore:       NewMemoryAddressRefStore(),
		ContractMetadataStore: NewMemoryContractMetadataStore[CM](),
		DomainMetadataStore:   NewMemoryDomainMetadataStore[DM](),
	}
}

// Seal seals the MemoryDataStore, by returning a new instance of sealedMemoryDataStore.
func (s *MemoryDataStore[CM, DM]) Seal() DataStore[CM, DM] {
	return &sealedMemoryDataStore[CM, DM]{
		AddressRefStore:       s.AddressRefStore,
		ContractMetadataStore: s.ContractMetadataStore,
		DomainMetadataStore:   s.DomainMetadataStore,
	}
}

// Addresses returns the AddressRefStore of the MemoryDataStore.
func (s *MemoryDataStore[CM, DM]) Addresses() MutableAddressRefStore {
	return s.AddressRefStore
}

// ContractMetadata returns the ContractMetadataStore of the MemoryDataStore.
func (s *MemoryDataStore[CM, DM]) ContractMetadata() MutableContractMetadataStore[CM] {
	return s.ContractMetadataStore
}

// DomainMetadata returns the DomainMetadataStore of the MemoryDataStore.
func (s *MemoryDataStore[CM, DM]) DomainMetadata() DomainMetadataStore[DM] {
	return s.DomainMetadataStore
}

// Merge merges the given mutable data store into the current MemoryDataStore.
func (s *MemoryDataStore[CM, DM]) Merge(other DataStore[CM, DM]) error {
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

	domainMetadata, ok, err := other.DomainMetadata().Get()
	if err != nil {
		return err
	}
	// If the domain metadata is not found, we don't need to update it. This
	// means that the domain metadata was not updated inside a changeset
	if ok {
		if err := s.DomainMetadataStore.Update(domainMetadata); err != nil {
			return err
		}
	}

	return nil
}

// SealedMemoryDataStore is a concrete implementation of the DataStore interface.
// It represents a sealed data store that cannot be modified further.
var _ DataStore[DefaultMetadata, DefaultMetadata] = &sealedMemoryDataStore[DefaultMetadata, DefaultMetadata]{}

type sealedMemoryDataStore[CM Cloneable[CM], DM Cloneable[DM]] struct {
	AddressRefStore       *MemoryAddressRefStore           `json:"addressRefStore"`
	ContractMetadataStore *MemoryContractMetadataStore[CM] `json:"contractMetadataStore"`
	DomainMetadataStore   *MemoryDomainMetadataStore[DM]   `json:"domainMetadataStore"`
}

// Addresses returns the AddressRefStore of the sealedMemoryDataStore.
// It implements the BaseDataStore interface.
//
//nolint:revive // this triggers a false positive confusing-naming linter error probably there are two implementations of Addresses() in the same file
func (s *sealedMemoryDataStore[CM, DM]) Addresses() AddressRefStore {
	return s.AddressRefStore
}

// ContractMetadata returns the ContractMetadataStore of the sealedMemoryDataStore.
//
//nolint:revive // this triggers a false positive confusing-naming linter error probably there are two implementations of ContractMetadata() in the same file
func (s *sealedMemoryDataStore[CM, DM]) ContractMetadata() ContractMetadataStore[CM] {
	return s.ContractMetadataStore
}

// DomainMetadata returns the DomainMetadataStore of the sealedMemoryDataStore.
//
//nolint:revive // this triggers a false positive confusing-naming linter error probably there are two implementations of DomainMetadata() in the same file
func (s *sealedMemoryDataStore[CM, DM]) DomainMetadata() DomainMetadataStore[DM] {
	return s.DomainMetadataStore
}
