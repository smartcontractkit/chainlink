package types

import (
	"encoding/json"

	eth "github.com/ethereum/go-ethereum/common"
)

type EthAddress eth.Address

func HexToEthAddress(s string) EthAddress {
	return EthAddress(eth.HexToAddress(s))
}

func BytesToEthAddress(b []byte) EthAddress {
	return EthAddress(eth.BytesToAddress(b))
}

func (a EthAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(eth.Address(a))
}

func (a EthAddress) Bytes() []byte {
	return eth.Address(a).Bytes()
}
