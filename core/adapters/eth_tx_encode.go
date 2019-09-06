package adapters

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/tidwall/gjson"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// EthTx holds the Address to send the result to and the FunctionSelector
// to execute.
type EthTxEncode struct {
	// Ethereum address of the contract this task calls
	Address common.Address `json:"address"`
	// Name of the contract method this task calls
	MethodName string `json:"methodName"`
	// Solidity types of the arguments to this method. (Must be primitive types.)
	// Keys are the argument names in `Order` field
	Types map[string]string `json:"types"`
	// Names of the arguments to the method, in appropriate order
	Order    []string    `json:"order"`
	GasPrice *models.Big `json:"gasPrice" gorm:"type:numeric"`
	GasLimit uint64      `json:"gasLimit"`
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (etx *EthTxEncode) Perform(
	input models.RunResult, store *strpkg.Store) models.RunResult {
	if !store.TxManager.Connected() {
		input.MarkPendingConnection()
		return input
	}
	if !input.Status.PendingConfirmations() {
		data, err := getTxEncodeData(etx, &input)
		if err != nil {
			input.SetError(errors.Wrap(err, "while constructing EthTx data"))
			return input
		}
		createTxRunResult(etx.Address, etx.GasPrice, etx.GasLimit, data, &input, store)
		return input
	}
	ensureTxRunResult(&input, store)
	return input
}

// ABI presents the method specified in etx as required by the abi package
func (etx EthTxEncode) ABI() (*abi.ABI, error) {
	rv := abi.ABI{}
	method := abi.Method{}
	method.Name = etx.MethodName
	method.RawName = etx.MethodName
	method.Inputs = make([]abi.Argument, len(etx.Order))
	for idx, argName := range etx.Order {
		typeName := etx.Types[argName]
		// TODO(alx): Enable composite types with a parser for the input JSON into
		// the correct golang types. Then `nil` must be replaced with the components
		typ, err := abi.NewType(typeName, nil)
		if err != nil {
			return nil, errors.Wrapf(
				err, `bad type for argument %s: %s`, argName, typeName)
		}
		method.Inputs[idx] = abi.Argument{Name: argName, Type: typ}
	}
	rv.Methods = make(map[string]abi.Method)
	rv.Methods[method.Name] = method
	return &rv, nil
}

// getTxData returns the data to save against the callback encoded according to
// the `types` and `order` fields of the job.
//
// At the time of writing it only uint256 arguments, so the use of the abi
// package is kind of overkill. Should make more complex method signatures
// easier in future, though.
func getTxEncodeData(e *EthTxEncode, input *models.RunResult) ([]byte, error) {
	abi, err := e.ABI()
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse method ABI from %+v", e)
	}
	// Extract values in the order specified by e.Order
	unorderedValues := input.Result()
	values := make([]interface{}, len(e.Order))
	for idx, name := range e.Order {
		switch e.Types[name] {
		case "uint256":
			var rawNum *big.Int
			rawVal := unorderedValues.Get(name)
			switch rawVal.Type {
			case gjson.String:
				rawNum, err = utils.HexToUint256(rawVal.Str)
				if err != nil {
					return nil, errors.Wrapf(
						err, "while casting argument %s for EthTxEncode", name)
				}
			case gjson.Number:
				var accuracy big.Accuracy
				rawNum, accuracy = big.NewFloat(rawVal.Num).Int(big.NewInt(0))
				if accuracy != big.Exact {
					return nil, errors.Errorf(
						"argument %s is not a whole number, as required for uint256 type",
						name)
				}
			default:
				return nil, errors.Errorf(
					"argument %s, which is of uint256 type, is not a number", name)
			}
			if rawNum.Cmp(big.NewInt(0)) == -1 {
				return nil, errors.Errorf(
					"cannot use negative number for uint256 argument %s", name)
			}
			values[idx] = rawNum
		default:
			return nil, errors.Errorf(
				"unimplelmented type for argument %s", name)
		}
	}
	data, err := abi.Pack(e.MethodName, values...)
	if err != nil {
		err = errors.Wrapf(err, "while packing %+v into %+v", values, abi)
	}
	return data, err
}
