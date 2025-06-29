package ccipaptos

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return addressBytesToString(addr)
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return addressStringToBytes(addr)
}

func (a AddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	addr := make([]byte, 32)

	// write oracleID in big endian as done by BCS for addresses
	binary.BigEndian.PutUint32(addr[28:], uint32(oracleID))

	return addr, nil
}

func (a AddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	// Transmitter accounts are ed25519 public keys, and encoded as a hex string without
	// a 0x prefix.
	return hex.EncodeToString(addr), nil
}

func addressBytesToString(addr []byte) (string, error) {
	if len(addr) < 1 || len(addr) > 32 {
		return "", fmt.Errorf("invalid Aptos address length (%d)", len(addr))
	}

	return fmt.Sprintf("0x%064x", addr), nil
}

func addressStringToBytes(addr string) ([]byte, error) {
	a := strings.TrimPrefix(addr, "0x")
	if len(a) == 0 {
		return nil, fmt.Errorf("invalid Aptos address length, expected at least 1 character: %s", addr)
	}
	if len(a) > 64 {
		return nil, fmt.Errorf("invalid Aptos address length, expected at most 64 characters: %s", addr)
	}
	if len(a) < 64 {
		a = strings.Repeat("0", 64-len(a)) + a
	}

	bytes, err := hex.DecodeString(a)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Aptos address '%s': %w", addr, err)
	}
	return bytes, nil
}

func addressBytesToBytes32(addr []byte) ([32]byte, error) {
	if len(addr) > 32 {
		return [32]byte{}, fmt.Errorf("invalid Aptos address length, expected 32, got %d", len(addr))
	}
	var result [32]byte
	// Left pad by copying to the end of the 32 byte array
	copy(result[32-len(addr):], addr)
	return result, nil
}
