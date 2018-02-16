package adapters

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// EthBytes32 holds no fields.
type EthBytes32 struct{}

const evmWordByteLen = 32
const evmWordHexLen = evmWordByteLen * 2

// Perform returns the hex value of the first 32 bytes of a string
// so that it is in the proper format to be written to the blockchain.
//
// For example, after converting the string "16800.00" to hex for Solidity
// type bytes32, it would be:
// "0x31363830302e3030000000000000000000000000000000000000000000000000"
func (*EthBytes32) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result, err := input.Get("value")
	if err != nil {
		return models.RunResultWithError(err)
	}

	value := common.RightPadBytes([]byte(result.String()), evmWordByteLen)
	hex := utils.RemoveHexPrefix(common.ToHex(value))

	if len(hex) > evmWordHexLen {
		hex = hex[:evmWordHexLen]
	}
	return models.RunResultWithValue(hex)
}

// EthUint256 holds no fields.
type EthUint256 struct{}

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "123.99" to hex for Solidity
// type uint256, it would be:
// "0x000000000000000000000000000000000000000000000000000000000000007b"
func (*EthUint256) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val, err := input.Get("value")
	if err != nil {
		return models.RunResultWithError(err)
	}

	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return models.RunResultWithError(fmt.Errorf("cannot parse into big.Int: %v", val.String()))
	}

	b, err := utils.HexToBytes(bigToUintHex(i))
	if err != nil {
		return models.RunResultWithError(err)
	}
	padded := common.LeftPadBytes(b, evmWordByteLen)
	hex := utils.RemoveHexPrefix(common.ToHex(padded))

	return models.RunResultWithValue(hex)
}

func bigToUintHex(f *big.Float) string {
	i, _ := f.Int(nil)
	if i.Sign() == -1 {
		i.Neg(i)
	}
	hex := fmt.Sprintf("%x", i)
	if len(hex)%2 != 0 {
		hex = "0" + hex
	}
	if len(hex) > evmWordHexLen {
		hex = hex[len(hex)-evmWordHexLen:]
	}
	return hex
}
