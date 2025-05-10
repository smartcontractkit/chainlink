package evm

// LogIterator defines the interface for iterating over events in Generated GETH Client `Filter_*` methods
// The GETH Client generates a concrete iterator type for each event. This interface allows us to mock
// the iterator in tests without needing to depend on the generated code.
type LogIterator[Event any] interface {
	// Next advances the iterator to the next event, returning false when no more events
	Next() bool

	// Error returns any error encountered during iteration
	Error() error

	// Close releases any resources associated with the iterator
	Close() error

	// GetEvent returns the current event
	GetEvent() *Event
}
