package types

import "fmt"

// Sequence represents the base type, for any chain's sequence object.
// It should be convertible to a string
type Sequence fmt.Stringer

// Id represents the base type, for any chain's Id.
// It should be convertible to a string, that can uniquely identify this chain
type Id fmt.Stringer
