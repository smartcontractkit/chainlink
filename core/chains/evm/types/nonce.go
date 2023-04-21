package types

import (
	"fmt"
	"strconv"
)

var _ fmt.Stringer = Nonce(0)

// Nonce wraps an EVM nonce into a stringable type
type Nonce int64

func (n Nonce) Int64() int64 {
	return int64(n)
}

func (n Nonce) String() string {
	return strconv.FormatInt(n.Int64(), 10)
}
