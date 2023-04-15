package types

import "fmt"

var _ fmt.Stringer = &Nonce{}

// Nonce wraps an EVM nonce into a stringable type
type Nonce struct {
	N int64
}

func (n Nonce) String() string {
	return string(n.N)
}
