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

	return common.BytesToAddress(bytes).Hex(), nil
}

// DecodeAddress takes an EIP-55 encoded Ethereum address (e.g., "0x...") and decodes it to a 20-byte array.
func (e EVMAddressModifier) DecodeAddress(str string) ([]byte, error) {
	str = strings.TrimPrefix(str, "0x")
	if len(str) != 40 {
		return nil, fmt.Errorf("%w: got length %d, expected 40 for address %s", commontypes.ErrInvalidType, len(str), str)
	}

	address := common.HexToAddress(str)
	if address == (common.Address{}) {
		return nil, fmt.Errorf("%w: address is zero", commontypes.ErrInvalidType)
	}

	return address.Bytes(), nil
}

// Length returns the expected length of an Ethereum address in bytes (20 bytes).
func (e EVMAddressModifier) Length() int {
	return common.AddressLength
}
