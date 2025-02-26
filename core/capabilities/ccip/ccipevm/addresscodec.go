package ccipevm

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	// TODO support EIP-55 checksum, https://smartcontract-it.atlassian.net/browse/CCIP-5340
	if len(addr) != common.AddressLength {
		return "", fmt.Errorf("invalid EVM address length, expected %v, got %d", common.AddressLength, len(addr))
	}

	return "0x" + hex.EncodeToString(addr), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	addrBytes, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(addr, "0x")))
	if err != nil {
		return nil, fmt.Errorf("failed to decode EVM address '%s': %w", addr, err)
	}

	return addrBytes, nil
}
