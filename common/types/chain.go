package types

import "fmt"

// TODO : Remove ID from txmgr/types to move to common/types
// ID represents the base type, for any chain's ID.
// It should be convertible to a string, that can uniquely identify this chain
type ID fmt.Stringer
