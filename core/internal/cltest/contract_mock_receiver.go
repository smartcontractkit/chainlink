package cltest

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
)

// funcSigLength is the length of the function signature (including the 0x)
// ex: 0x1234ABCD
const funcSigLength = 10

func NewContractMockReceiver(t *testing.T, ethMock *evmclimocks.Client, abi abi.ABI, address common.Address) contractMockReceiver {
	return contractMockReceiver{
		t:       t,
		ethMock: ethMock,
		abi:     abi,
		address: address,
	}
}

type contractMockReceiver struct {
	t       *testing.T
	ethMock *evmclimocks.Client
	abi     abi.ABI
	address common.Address
}

func (receiver contractMockReceiver) MockCallContractResponse(funcName string, responseArgs ...interface{}) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != funcSigLength {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	encoded := receiver.mustEncodeResponse(funcName, responseArgs...)

	return receiver.ethMock.
		On(
			"CallContract",
			mock.Anything,
			mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
				return *callArgs.To == receiver.address &&
					hexutil.Encode(callArgs.Data)[0:funcSigLength] == funcSig
			}),
			mock.Anything).
		Return(encoded, nil)
}

func (receiver contractMockReceiver) MockCallContextResponse(funcName string, responseArgs ...interface{}) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != funcSigLength {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	encoded := receiver.mustEncodeResponse(funcName, responseArgs...)

	return receiver.ethMock.
		On(
			"CallContext",
			mock.Anything,
			mock.Anything,
			"eth_call",
			mock.MatchedBy(func(args map[string]interface{}) bool {
				to := args["to"].(*common.Address)
				data := args["input"].(hexutil.Bytes)
				return *to == receiver.address &&
					hexutil.Encode(data)[0:funcSigLength] == funcSig
			}),
			mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(1).(*hexutil.Bytes)
		*resp = encoded
	})

}

func (receiver contractMockReceiver) MockCallContextMatchedResponse(funcName string, matcher func(args map[string]interface{}) bool, responseArgs ...interface{}) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != funcSigLength {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	encoded := receiver.mustEncodeResponse(funcName, responseArgs...)

	// TODO: ALL CALLER MATCHER FUNCTIONS SHOULD BE CHANGED

	return receiver.ethMock.
		On(
			"CallContext",
			mock.Anything,
			mock.Anything,
			"eth_call",
			mock.MatchedBy(func(args map[string]interface{}) bool {
				to := args["to"].(*common.Address)
				data := args["input"].(hexutil.Bytes)
				return *to == receiver.address &&
					hexutil.Encode(data)[0:funcSigLength] == funcSig &&
					matcher(args)
			}),
			mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(1).(*hexutil.Bytes)
		*resp = encoded
	})
}

func (receiver contractMockReceiver) MockCallContextRevertResponse(funcName string) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != funcSigLength {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	return receiver.ethMock.
		On(
			"CallContext",
			mock.Anything,
			mock.Anything,
			"eth_call",
			mock.MatchedBy(func(args map[string]interface{}) bool {
				to := args["to"].(*common.Address)
				data := args["input"].(hexutil.Bytes)
				return *to == receiver.address &&
					hexutil.Encode(data)[0:funcSigLength] == funcSig
			}),
			mock.Anything).
		Return(errors.New("revert"))

}

func (receiver contractMockReceiver) mustEncodeResponse(funcName string, responseArgs ...interface{}) []byte {
	if len(responseArgs) == 0 {
		return []byte{}
	}

	var outputList []interface{}

	firstArg := responseArgs[0]
	isStruct := reflect.TypeOf(firstArg).Kind() == reflect.Struct

	if isStruct && len(responseArgs) > 1 {
		receiver.t.Fatal("cannot encode response with struct and multiple return values")
	} else if isStruct {
		outputList = structToInterfaceSlice(firstArg)
	} else {
		outputList = responseArgs
	}

	encoded, err := receiver.abi.Methods[funcName].Outputs.PackValues(outputList)
	require.NoError(receiver.t, err)
	return encoded
}

func structToInterfaceSlice(structArg interface{}) []interface{} {
	v := reflect.ValueOf(structArg)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	return values
}
