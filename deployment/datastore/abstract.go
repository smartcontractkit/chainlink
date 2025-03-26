package datastore

// Cloneable provides a Clone() method which returns a semi-deep copy of the type.
type Cloneable[R any] interface {
	// Clone() returns a semi-deep copy of the type. The implementation should call Clone() on
	// any nested Cloneable fields, and perform a shallow copy of any slices or maps.
	// NOTE: Handling non-Cloneable references to Cloneable types is beyond the intended scope of this interface.
	Clone() R
}

// Comparable provides an Equals() method which returns true if the two instances are equal, false otherwise.
type Comparable[T any] interface {
	// Equals()	returns true if the two instances are equal, false otherwise.
	Equals(T) bool
}

// Fetcher provides a Fetch() method which is used to complete a read query from a Store.
type Fetcher[R any] interface {
	// Fetch() returns a slice of records representing the entire data set. The returned slice
	// will be a newly allocated slice (not a reference to an existing one), and each record should
	// be a copy of the corresponding stored data. Modifying the returned slice or its records must
	// not affect the underlying data.
	Fetch() ([]R, error)
}

// Getter provides a Get() method which is used to complete a read by key query from a Store.
type Getter[K Comparable[K], R Record[K, R]] interface {
	// Get() returns the record with the given key, or an error if no such record exists.
	Get(K) (R, error)
}

// PrimaryKeyHolder is an interface for types that can provide a unique identifier key for themselves.
type PrimaryKeyHolder[K Comparable[K]] interface {
	// Key() returns the primary key for the implementing type.
	Key() K
}

// Record represents a data entry that is both Cloneable and uniquely identifiable by its primary key.
type Record[K Comparable[K], R PrimaryKeyHolder[K]] interface {
	Cloneable[R]
	PrimaryKeyHolder[K]
}

// FilterFunc is a function that filters a slice of records.
type FilterFunc[K Comparable[K], R Record[K, R]] func([]R) []R

// Filterable provides a Filter() method which is used to complete a filtered query with from a Store.
type Filterable[K Comparable[K], R Record[K, R]] interface {
	Filter(filters ...FilterFunc[K, R]) []R
}

// Store is an interface that represents an immutable set of records.
type Store[K Comparable[K], R Record[K, R]] interface {
	Fetcher[R]
	Getter[K, R]
	Filterable[K, R]
}

// MutableStore is an interface that represents a mutable set of records.
type MutableStore[K Comparable[K], R Record[K, R]] interface {
	Store[K, R]

	// Add inserts a new record into the MutableStore.
	Add(record R) error

	// AddOrUpdate behaves like Add where there is not already a record with the same composite primary key as the
	// supplied record, otherwise it behaves like an update.
	AddOrUpdate(record R) error

	// Update edits an existing record whose fields match the primary key elements of the supplied AddressRecord, with
	// the non-primary-key values of the supplied AddressRecord.
	Update(record R) error

	// Delete deletes record whose primary key elements match the supplied key, returning an error if no
	// such record exists to be deleted
	Delete(key K) error
}
