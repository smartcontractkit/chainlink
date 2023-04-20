package types

// A chain-agnostic generic interface to represent the following native types on various chains:
// PublicKey, Address, Account, BlockHash, TxHash
//
//go:generate mockery --quiet --name Hashable --output ./mocks/ --case=underscore
type Hashable interface {
	MarshalText() (text []byte, err error)
	String() string
	comparable
}
