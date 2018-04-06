package adapters

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

var evmUint256Max *big.Int

func init() {
	var ok bool
	evmUint256Max, ok = (&big.Int{}).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	if !ok {
		panic("could not parse evmUint256Max")
	}
}

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

	parts := strings.Split(val.String(), ".")
	i, ok := (&big.Int{}).SetString(parts[0], 10)
	if !ok {
		return input.WithError(fmt.Errorf("cannot parse into big.Int: %v", val.String()))
	}

	if err = validateRange(i); err != nil {
		return input.WithError(err)
	}

	return input.WithValue(utils.EVMHexNumber(i))
}

func validateRange(i *big.Int) error {
	if i.Sign() == -1 {
		return fmt.Errorf("ethUint256: value %v is negative", i.String())
	}

	if evmUint256Max.Cmp(i) == -1 {
		return fmt.Errorf("ethUint256: value %v too large", i.String())
	}
	return nil
}
