package types

import "fmt"

// Sequence represents the base type, for any chain's sequence object.
// It should be convertible to a string
type Sequence interface {
	fmt.Stringer
	Int64() int64 // needed for numeric sequence confirmation - to be removed with confirmation logic generalization: https://smartcontract-it.atlassian.net/browse/BCI-860
}

type Unit fmt.Stringer
