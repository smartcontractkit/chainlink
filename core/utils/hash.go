package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Hash is a simplified version of go-ethereum's common.Hash to avoid
// go-ethereum dependency
// It represents a 32 byte fixed size array that marshals/unmarshals assuming a
// 0x prefix
type Hash [32]byte

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return fmt.Sprintf("0x%s", hex.EncodeToString(h[:])) }

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	if !strings.HasPrefix(string(input), "0x") {
		return errors.New("hash: expected a hex string starting with '0x'")
	}
	phex := new(PlainHexBytes)
	if err := phex.UnmarshalText(input[2:]); err != nil {
		return fmt.Errorf("hash: %w", err)
	}
	if len(*phex) != 32 {
		return fmt.Errorf("hash: expected 32-byte sequence, got %d bytes", len(*phex))
	}
	copy((*h)[:], (*phex))
	return nil
}
