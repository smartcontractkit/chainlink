package types

import "fmt"

// SEQUENCE represents constraint on a chain's sequence object.
// It should be convertible to a string
type SEQUENCE fmt.Stringer

// ID represents constraint on a chain's Id.
// It should be convertible to a string, that can uniquely identify this chain
type ID fmt.Stringer
