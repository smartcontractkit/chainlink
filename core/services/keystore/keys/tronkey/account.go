package tronkey

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

// Extracted from go-tron sdk: https://github.com/fbsobreira/gotron-sdk

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLengthBase58 is the expected length of the address in base58format
	AddressLengthBase58 = 34
	// Tron Address Prefix
	prefixMainnet = 0x41
	// TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(prefixMainnet)
	// Tron address should have 21 bytes (20 bytes + 1 byte prefix)
	AddressLength = 21
)

// Address represents the 21 byte address of an Tron account.
type Address [AddressLength]byte

// Bytes get bytes from address
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex get bytes from address in string
func (a Address) Hex() string {
	return BytesToHexString(a[:])
}

// HexToAddress returns Address with byte values of s.
func HexToAddress(s string) (Address, error) {
	addr, err := FromHex(s)
	if err != nil {
		return Address{}, err
	}
	// Check if the address starts with '41' and is 21 characters long
	if len(addr) != AddressLength || addr[0] != prefixMainnet {
		return Address{}, errors.New("invalid Tron address")
	}
	return Address(addr), nil
}

// Base58ToAddress returns Address with byte values of s.
func Base58ToAddress(s string) (Address, error) {
	addr, err := DecodeCheck(s)
	if err != nil {
		return Address{}, err
	}
	return Address(addr), nil
}

// String implements fmt.Stringer.
// Returns the address as a base58 encoded string.
func (a Address) String() string {
	if len(a) == 0 {
		return ""
	}

	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}
	return EncodeCheck(a.Bytes())
}

// PubkeyToAddress returns address from ecdsa public key
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	address := crypto.PubkeyToAddress(p)

	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return Address(addressTron)
}

// BytesToHexString encodes bytes as a hex string.
func BytesToHexString(bytes []byte) string {
	encode := make([]byte, len(bytes)*2)
	hex.Encode(encode, bytes)
	return "0x" + string(encode)
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) ([]byte, error) {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return HexToBytes(s)
}

// Has0xPrefix validates str begins with '0x' or '0X'.
func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// HexToBytes returns the bytes represented by the hexadecimal string str.
func HexToBytes(str string) ([]byte, error) {
	return hex.DecodeString(str)
}

func Encode(input []byte) string {
	return base58.Encode(input)
}

func EncodeCheck(input []byte) string {
	h256h0 := sha256.New()
	h256h0.Write(input)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	inputCheck := input
	inputCheck = append(inputCheck, h1[:4]...)

	return Encode(inputCheck)
}

func DecodeCheck(input string) ([]byte, error) {
	decodeCheck, err := Decode(input)
	if err != nil {
		return nil, err
	}

	if len(decodeCheck) < 4 {
		return nil, errors.New("base58 check error")
	}

	// tron address should should have 21 bytes (including prefix) + 4 checksum
	if len(decodeCheck) != AddressLength+4 {
		return nil, fmt.Errorf("invalid address length: %d", len(decodeCheck))
	}

	// check prefix
	if decodeCheck[0] != prefixMainnet {
		return nil, errors.New("invalid prefix")
	}

	decodeData := decodeCheck[:len(decodeCheck)-4]

	h256h0 := sha256.New()
	h256h0.Write(decodeData)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	if h1[0] == decodeCheck[len(decodeData)] &&
		h1[1] == decodeCheck[len(decodeData)+1] &&
		h1[2] == decodeCheck[len(decodeData)+2] &&
		h1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}

	return nil, errors.New("base58 check error")
}

func Decode(input string) ([]byte, error) {
	return base58.Decode(input)
}
