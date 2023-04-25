package types

import "fmt"

// Sequence represents the base type, for any chain's sequence object.
// It should be convertible to a string
type Sequence fmt.Stringer

// ID represents the base type, for any chain's ID.
// It should be convertible to a string, that can uniquely identify this chain
type ID fmt.Stringer
