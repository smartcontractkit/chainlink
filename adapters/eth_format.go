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
var evmInt256Max *big.Int
var evmInt256Min *big.Int

func init() {
	var ok bool
	evmInt256Max, ok = (&big.Int{}).SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	if !ok {
		panic("could not parse evmInt256Max")
	}
	evmInt256Min, ok = (&big.Int{}).SetString("-8fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	if !ok {
		panic("could not parse evmInt256Min")
	}
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
// For example, after converting the string "16800.01" to hex encoded Ethereum
// ABI, it would be:
// "0x31363830302e3031000000000000000000000000000000000000000000000000"
func (*EthBytes32) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result := input.Get("value")
	value := common.RightPadBytes([]byte(result.String()), utils.EVMWordByteLen)
	hex := utils.RemoveHexPrefix(common.ToHex(value))

	if len(hex) > utils.EVMWordHexLen {
		hex = hex[:utils.EVMWordHexLen]
	}
	return input.WithValue(utils.AddHexPrefix(hex))
}

// EthInt256 holds no fields
type EthInt256 struct{}

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "-123.99" to hex encoded Ethereum
// ABI, it would be:
// "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff85"
func (*EthInt256) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	i, err := parseBigInt(input)
	if err != nil {
		return input.WithError(err)
	}

	if err = validateSignedRange(i); err != nil {
		return input.WithError(err)
	}

	sh, err := utils.EVMSignedHexNumber(i)
	if err != nil {
		return input.WithError(err)
	}

	return input.WithValue(sh)
}

// EthUint256 holds no fields.
type EthUint256 struct{}

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "123.99" to hex encoded Ethereum
// ABI, it would be:
// "0x000000000000000000000000000000000000000000000000000000000000007b"
func (*EthUint256) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	i, err := parseBigInt(input)
	if err != nil {
		return input.WithError(err)
	}

	if err = validateUnsignedRange(i); err != nil {
		return input.WithError(err)
	}

	return input.WithValue(utils.EVMHexNumber(i))
}

func parseBigInt(input models.RunResult) (*big.Int, error) {
	val := input.Get("value")
	parts := strings.Split(val.String(), ".")
	i, ok := (&big.Int{}).SetString(parts[0], 10)
	if !ok {
		return nil, fmt.Errorf("cannot parse into big.Int: %v", val)
	}
	return i, nil
}

func validateSignedRange(i *big.Int) error {
	if evmInt256Max.Cmp(i) == -1 {
		return fmt.Errorf("ethInt256: value %v too large", i.String())
	}

	if evmInt256Min.Cmp(i) == 1 {
		return fmt.Errorf("ethInt256: value %v too small", i.String())
	}
	return nil
}

func validateUnsignedRange(i *big.Int) error {
	if i.Sign() == -1 {
		return fmt.Errorf("ethUint256: value %v is negative", i.String())
	}

	if evmUint256Max.Cmp(i) == -1 {
		return fmt.Errorf("ethUint256: value %v too large", i.String())
	}
	return nil
}
