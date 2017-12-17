package adapters

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-go/models"
)

type EthBytes32 struct{}

const maxBytes32HexLength = 32 * 2

func (self *EthBytes32) Perform(input models.RunResult) models.RunResult {
	value := common.RightPadBytes([]byte(input.Value()), 32)
	hex := removeHexPrefix(common.ToHex(value))

	if len(hex) > maxBytes32HexLength {
		hex = hex[0:maxBytes32HexLength]
	}
	return models.RunResultWithValue(hex)
}

func removeHexPrefix(hex string) string {
	if hex[0:2] == "0x" {
		return hex[2:len(hex)]
	}
	return hex
}
