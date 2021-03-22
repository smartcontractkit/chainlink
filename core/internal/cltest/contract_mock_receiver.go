package cltest

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewContractMockReceiver(t *testing.T, ethMock *mocks.Client, abi abi.ABI, address common.Address) contractMockReceiver {
	return contractMockReceiver{
		t:       t,
		ethMock: ethMock,
		abi:     abi,
		address: address,
	}
}

type contractMockReceiver struct {
	t       *testing.T
	ethMock *mocks.Client
	abi     abi.ABI
	address common.Address
}

func (receiver contractMockReceiver) MockResponse(funcName string, responseArgs ...interface{}) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != 10 {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	encoded := receiver.mustEncodeResponse(funcName, responseArgs)

	return receiver.ethMock.
		On(
			"CallContract",
			mock.Anything,
			mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
				return *callArgs.To == receiver.address &&
					hexutil.Encode(callArgs.Data)[0:10] == funcSig
			}),
			mock.Anything).
		Return(encoded, nil)
}

func (receiver contractMockReceiver) mustEncodeResponse(funcName string, responseArgs []interface{}) []byte {
	if len(responseArgs) == 0 {
		return []byte{}
	}

	var outputList []interface{}

	firstArg := responseArgs[0]
	isStruct := reflect.TypeOf(firstArg).Kind() == reflect.Struct

	if isStruct && len(responseArgs) > 1 {
		receiver.t.Fatal("cannot encode resonse with struct and multiple return values")
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
