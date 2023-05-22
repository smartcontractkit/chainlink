package types

import "fmt"

// ID represents the base type, for any chain's ID.
// It should be convertible to a string, that can uniquely identify this chain
type ID fmt.Stringer
