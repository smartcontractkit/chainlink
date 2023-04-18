package types

import (
	"fmt"
	"strconv"
)

var _ fmt.Stringer = &Nonce{}

// Nonce wraps an EVM nonce into a stringable type
type Nonce struct {
	N int64
}

func (n Nonce) String() string {
	return strconv.FormatInt(n.N, 10)
}
