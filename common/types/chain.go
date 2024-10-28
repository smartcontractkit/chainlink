package types

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// Sequence represents the base type, for any chain's sequence object.
// It should be convertible to a string
type Sequence interface {
	fmt.Stringer
	Int64() int64 // needed for numeric sequence confirmation - to be removed with confirmation logic generalization: https://smartcontract-it.atlassian.net/browse/BCI-860
}

// ID represents the base type, for any chain's ID.
// It should be convertible to a string, that can uniquely identify this chain
type ID fmt.Stringer

// ChainStatusWithID compose of ChainStatus and RelayID. This is useful for
// storing the Network associated with the ChainStatus.
type ChainStatusWithID struct {
	types.ChainStatus
	types.RelayID
}
