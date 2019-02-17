package adapters

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// EthBytes32 holds no fields.
type EthBytes32 struct{}

// Perform returns the hex value of the first 32 bytes of a string
// so that it is in the proper format to be written to the blockchain.
//
// For example, after converting the string "16800.01" to hex encoded Ethereum
// ABI, it would be:
// "0x31363830302e3031000000000000000000000000000000000000000000000000"
func (*EthBytes32) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result := input.Result()
	value := common.RightPadBytes([]byte(result.String()), utils.EVMWordByteLen)
	hex := utils.RemoveHexPrefix(hexutil.Encode(value))

	if len(hex) > utils.EVMWordHexLen {
		hex = hex[:utils.EVMWordHexLen]
	}
	return input.WithResult(utils.AddHexPrefix(hex))
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
	value, err := utils.EVMTranscodeInt256(input.Result())
	if err != nil {
		return input.WithError(err)
	}

	return input.WithResult(hexutil.Encode(value))
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
	value, err := utils.EVMTranscodeUint256(input.Result())
	if err != nil {
		return input.WithError(err)
	}

	return input.WithResult(hexutil.Encode(value))
}

// EthBytesRaw holds no fields.
type EthBytesRaw struct{}

// Perform converts the string to the raw bytes equivalent as if it was hex.
//
// For example, "000aa"
// ABI, it would be:
// "0x00000000000000000000000000000000000000000000000000000000000000aa"
func (*EthBytesRaw) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	result := input.Result()
	h, err := hexutil.Decode(utils.AddHexPrefix(result.String()))
	if err != nil {
		return input.WithError(err)
	}
	value := common.LeftPadBytes(h, utils.EVMWordByteLen)
	hex := hexutil.Encode(value)

	// hex encoded 0x
	if len(hex) > utils.EVMWordHexLen+2 {
		hex = hex[:utils.EVMWordHexLen+2]
	}
	return input.WithResult(hex)
}
