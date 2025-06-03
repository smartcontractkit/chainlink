package ccipevm

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return common.BytesToAddress(addr).Hex(), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("invalid EVM address: %s", addr)
	}
	return common.HexToAddress(addr).Bytes(), nil
}

func (a AddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	addr := make([]byte, 20)

	// write oracleID into addr in big endian
	binary.BigEndian.PutUint32(addr, uint32(oracleID))

	return common.BytesToAddress(addr).Bytes(), nil
}
