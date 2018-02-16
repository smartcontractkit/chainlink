package adapters

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// EthBytes32 holds no fields
type EthBytes32 struct{}

const maxBytes32HexLength = 32 * 2

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "16800.00" to hex for
// the blockchain, it would be:
// "31363830302e3030000000000000000000000000000000000000000000000000"
func (eba *EthBytes32) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result, err := input.Get("value")
	if err != nil {
		return models.RunResultWithError(err)
	}

	value := common.RightPadBytes([]byte(result.String()), 32)
	hex := utils.RemoveHexPrefix(common.ToHex(value))

	if len(hex) > maxBytes32HexLength {
		hex = hex[0:maxBytes32HexLength]
	}
	return models.RunResultWithValue(hex)
}
