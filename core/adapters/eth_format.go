package adapters

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	input.CompleteWithResult(utils.AddHexPrefix(hex))
	return input
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
		input.SetError(err)
		return input
	}

	input.CompleteWithResult(hexutil.Encode(value))
	return input
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
		input.SetError(err)
		return input
	}

	input.CompleteWithResult(hexutil.Encode(value))
	return input
}
