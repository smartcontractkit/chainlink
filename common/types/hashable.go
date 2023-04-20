package types

// A chain-agnostic generic interface to represent the following native types on various chains:
// PublicKey, Address, Account, BlockHash, TxHash
type Hashable interface {
	MarshalText() (text []byte, err error)
	String() string
	comparable
}
