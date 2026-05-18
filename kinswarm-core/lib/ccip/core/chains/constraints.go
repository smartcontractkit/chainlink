package chains

// ID types represent unique identifiers within a particular chain type. Using string is recommended.
type ID any

// Node types should be a struct including these default fields:
//
//	ID        int32
//	Name      string
type Node any
