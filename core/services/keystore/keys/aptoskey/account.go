package aptoskey

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// AccountAddress is a 32 byte address on the Aptos blockchain
// It can represent an Object, an Account, and much more.
//
// AccountAddress is copied from the aptos sdk because:
//  1. There are still breaking changes in sdk and we don't want the dependency.
//  2. AccountAddress is just a wrapper and can be easily extracted out.
//
// https://github.com/aptos-labs/aptos-go-sdk/blob/main/internal/types/account.go
type AccountAddress [32]byte

// IsSpecial Returns whether the address is a "special" address. Addresses are considered
// special if the first 63 characters of the hex string are zero. In other words,
// an address is special if the first 31 bytes are zero and the last byte is
// smaller than `0b10000` (16). In other words, special is defined as an address
// that matches the following regex: `^0x0{63}[0-9a-f]$`. In short form this means
// the addresses in the range from `0x0` to `0xf` (inclusive) are special.
// For more details see the v1 address standard defined as part of AIP-40:
// https://github.com/aptos-foundation/AIPs/blob/main/aips/aip-40.md
func (aa *AccountAddress) IsSpecial() bool {
	for _, b := range aa[:31] {
		if b != 0 {
			return false
		}
	}
	return aa[31] < 0x10
}

// String Returns the canonical string representation of the AccountAddress
func (aa *AccountAddress) String() string {
	if aa.IsSpecial() {
		return fmt.Sprintf("0x%x", aa[31])
	}
	return BytesToHex(aa[:])
}

// ParseStringRelaxed parses a string into an AccountAddress
func (aa *AccountAddress) ParseStringRelaxed(x string) error {
	x = strings.TrimPrefix(x, "0x")
	if len(x) < 1 {
		return ErrAddressTooShort
	}
	if len(x) > 64 {
		return ErrAddressTooLong
	}
	if len(x)%2 != 0 {
		x = "0" + x
	}
	bytes, err := hex.DecodeString(x)
	if err != nil {
		return err
	}
	// zero-prefix/right-align what bytes we got
	copy((*aa)[32-len(bytes):], bytes)

	return nil
}

var ErrAddressTooShort = errors.New("AccountAddress too short")
var ErrAddressTooLong = errors.New("AccountAddress too long")

func BytesToHex(bytes []byte) string {
	return "0x" + hex.EncodeToString(bytes)
}
