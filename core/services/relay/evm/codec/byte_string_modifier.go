package codec

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// EVMAddressModifier implements the AddressModifier interface for Ethereum addresses.
// It handles encoding and decoding Ethereum addresses with EIP-55 checksums and hex encoding.
type EVMAddressModifier struct{}

func (e EVMAddressModifier) EncodeAddress(bytes []byte) (string, error) {
	if len(bytes) != e.Length() {
		return "", fmt.Errorf("%w: got length %d, expected 20 for bytes %x", commontypes.ErrInvalidType, len(bytes), bytes)
	}
	// Convert bytes to common.Address type and return the EIP-55 encoded string.

	return common.BytesToAddress(bytes).Hex(), nil
}

// DecodeAddress takes an EIP-55 encoded Ethereum address (e.g., "0x...") and decodes it to a 20-byte array.
func (e EVMAddressModifier) DecodeAddress(str string) ([]byte, error) {
	// Remove the "0x" prefix if present.
	str = strings.TrimPrefix(str, "0x")

	// Validate the address length (40 hex characters for a 20-byte address).
	if len(str) != 40 {
		return nil, fmt.Errorf("%w: got length %d, expected 40 for address %s", commontypes.ErrInvalidType, len(str), str)
	}

	// Use HexToAddress to parse the string into a common.Address type.
	address := common.HexToAddress(str)

	// Ensure the decoded address matches the original string to avoid zero-value addresses.
	if address == (common.Address{}) {
		return nil, fmt.Errorf("%w: address is zero", commontypes.ErrInvalidType)
	}

	return address.Bytes(), nil
}

// Length returns the expected length of an Ethereum address in bytes (20 bytes).
func (e EVMAddressModifier) Length() int {
	return 20
}
