package ccipevm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return hexutil.Encode(addr), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return common.HexToAddress(addr).Bytes(), nil
}
