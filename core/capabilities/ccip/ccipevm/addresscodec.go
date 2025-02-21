package ccipevm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	if len(addr) != common.AddressLength {
		return "", fmt.Errorf("invalid EVM address length, expected %v, got %d", common.AddressLength, len(addr))
	}

	return common.BytesToAddress(addr).Hex(), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("invalid EVM address %s", addr)
	}
	return common.HexToAddress(addr).Bytes(), nil
}
