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

// Perform returns the hex value of the first 32 bytes of a string
// so that it is in the proper format to be written to the blockchain.
//
// For example, after converting the string "123.99" to hex for
// the blockchain, it would be:
// "0x000000000000000000000000000000000000000000000000000000000000007b"
func (*EthBytes32) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result, err := input.Get("value")
	if err != nil {
		return input.WithError(err)
	}

	value := common.RightPadBytes([]byte(result.String()), utils.EVMWordByteLen)
	hex := utils.RemoveHexPrefix(common.ToHex(value))

	if len(hex) > utils.EVMWordHexLen {
		hex = hex[:utils.EVMWordHexLen]
	}
	return input.WithValue(utils.AddHexPrefix(hex))
}

// EthUint256 holds no fields.
type EthUint256 struct{}

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "16800.00" to hex for
// the blockchain, it would be:
// "0x31363830302e3030000000000000000000000000000000000000000000000000"
func (*EthUint256) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val, err := input.Get("value")
	if err != nil {
		return input.WithError(err)
	}

	i, ok := (&big.Float{}).SetString(val.String())
	if !ok {
		return input.WithError(fmt.Errorf("cannot parse into big.Float: %v", val.String()))
	}

	b, err := utils.HexToBytes(bigToUintHex(i))
	if err != nil {
		return input.WithError(err)
	}
	padded := common.LeftPadBytes(b, utils.EVMWordByteLen)

	return input.WithValue(common.ToHex(padded))
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
	if len(hex) > utils.EVMWordHexLen {
		hex = hex[len(hex)-utils.EVMWordHexLen:]
	}
	return hex
}
