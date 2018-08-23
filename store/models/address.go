package models

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// EIP55Address is a newtype for string which persists an ethereum address in
// its original string representation which includes a leading 0x, and EIP55
// checksum which is represented by the case of digits A-F.
type EIP55Address string

// NewEIP55Address creates an EIP55Address from a string, an error is returned if:
//
// 1) There is no leading 0x
// 2) The length is wrong
// 3) There are any non hexadecimal characters
// 4) The checksum fails
//
func NewEIP55Address(s string) (EIP55Address, error) {
	address := common.HexToAddress(s)
	if s != address.Hex() {
		return EIP55Address(""), fmt.Errorf(`"%s" is not a valid EIP55 formatted address`, s)
	}
	return EIP55Address(s), nil
}

// Bytes returns the raw bytes
func (a EIP55Address) Bytes() []byte { return a.Address().Bytes() }

// Big returns a big.Int representation
func (a EIP55Address) Big() *big.Int { return a.Address().Big() }

// Hash returns the Hash
func (a EIP55Address) Hash() common.Hash { return a.Address().Hash() }

// Address returns EIP55Address as a go-ethereum Address type
func (a EIP55Address) Address() common.Address { return common.HexToAddress(a.String()) }

// String implements the stringer interface and is used also by the logger.
func (a EIP55Address) String() string {
	return string(a)
}

// Format implements fmt.Formatter
func (a EIP55Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, a.String())
}

// UnmarshalText parses a hash from plain text
func (a *EIP55Address) UnmarshalText(input []byte) error {
	var err error
	*a, err = NewEIP55Address(string(input))
	return err
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// UnmarshalJSON parses a hash from a JSON string
func (a *EIP55Address) UnmarshalJSON(input []byte) error {
	if isString(input) {
		input = input[1 : len(input)-1]
	}

	return a.UnmarshalText([]byte(input))
}
