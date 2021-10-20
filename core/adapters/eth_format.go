package adapters

import (
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthBytes32 holds no fields.
type EthBytes32 struct{}

// TaskType returns the type of Adapter.
func (e *EthBytes32) TaskType() models.TaskType {
	return TaskTypeEthBytes32
}

// Perform returns the hex value of the first 32 bytes of a string
// so that it is in the proper format to be written to the blockchain.
//
// For example, after converting the string "16800.01" to hex encoded Ethereum
// ABI, it would be:
// "0x31363830302e3031000000000000000000000000000000000000000000000000"
func (*EthBytes32) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	result := input.Result()
	value := common.RightPadBytes([]byte(result.String()), utils.EVMWordByteLen)
	hex := utils.RemoveHexPrefix(hexutil.Encode(value))

	if len(hex) > utils.EVMWordHexLen {
		hex = hex[:utils.EVMWordHexLen]
	}

	return models.NewRunOutputCompleteWithResult(utils.AddHexPrefix(hex), input.ResultCollection())
}

// EthInt256 holds no fields
type EthInt256 struct{}

// TaskType returns the type of Adapter.
func (e *EthInt256) TaskType() models.TaskType {
	return TaskTypeEthInt256
}

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "-123.99" to hex encoded Ethereum
// ABI, it would be:
// "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff85"
func (*EthInt256) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	value, err := utils.EVMTranscodeInt256(input.Result())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	return models.NewRunOutputCompleteWithResult(hexutil.Encode(value), input.ResultCollection())
}

// EthUint256 holds no fields.
type EthUint256 struct{}

// TaskType returns the type of Adapter.
func (e *EthUint256) TaskType() models.TaskType {
	return TaskTypeEthUint256
}

// Perform returns the hex value of a given string so that it
// is in the proper format to be written to the blockchain.
//
// For example, after converting the string "123.99" to hex encoded Ethereum
// ABI, it would be:
// "0x000000000000000000000000000000000000000000000000000000000000007b"
func (*EthUint256) Perform(input models.RunInput, _ *store.Store, _ *keystore.Master) models.RunOutput {
	value, err := utils.EVMTranscodeUint256(input.Result())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	return models.NewRunOutputCompleteWithResult(hexutil.Encode(value), input.ResultCollection())
}
