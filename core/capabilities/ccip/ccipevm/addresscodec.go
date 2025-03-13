package ccipevm

import (
	"encoding/hex"
	"fmt"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return gethcommon.BytesToAddress(addr).Hex(), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	addrBytes, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(addr, "0x")))
	if err != nil {
		return nil, fmt.Errorf("failed to decode EVM address '%s': %w", addr, err)
	}

	return addrBytes, nil
}
