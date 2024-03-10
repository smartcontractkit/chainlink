package types

import "fmt"

// A chain-agnostic generic interface to represent the following native types on various chains:
// PublicKey, Address, Account, BlockHash, TxHash
type Hashable interface {
	fmt.Stringer
	comparable

	Bytes() []byte
}
