package chain

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// funcSigLength is the length of the function signature (including the 0x)
// ex: 0x1234ABCD
const funcSigLength = 10

type ContractMockReceiver struct {
	t       *testing.T
	ethMock *types.MockEVMClient
	abi     abi.ABI
}

func NewContractMockReceiver(t *testing.T, ethMock *types.MockEVMClient, abi abi.ABI) ContractMockReceiver {
	return ContractMockReceiver{
		t:       t,
		ethMock: ethMock,
		abi:     abi,
	}
}

func (receiver ContractMockReceiver) MockResponse(funcName string, responseArgs ...interface{}) *mock.Call {
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
				return hexutil.Encode(callArgs.Data)[0:funcSigLength] == funcSig
			}),
			mock.Anything).
		Return(encoded, nil)
}

func (receiver ContractMockReceiver) MockRevertResponse(funcName string, msg string) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != funcSigLength {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	return receiver.ethMock.
		On(
			"CallContract",
			mock.Anything,
			mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
				return hexutil.Encode(callArgs.Data)[0:funcSigLength] == funcSig
			}),
			mock.Anything).
		Return(nil, fmt.Errorf("revert%s", msg))
}

func (receiver ContractMockReceiver) MockNonRevertError(funcName string, err error, after time.Duration) *mock.Call {
	funcSig := hexutil.Encode(receiver.abi.Methods[funcName].ID)
	if len(funcSig) != funcSigLength {
		receiver.t.Fatalf("Unable to find Registry contract function with name %s", funcName)
	}

	return receiver.ethMock.
		On(
			"CallContract",
			mock.Anything,
			mock.MatchedBy(func(callArgs ethereum.CallMsg) bool {
				return hexutil.Encode(callArgs.Data)[0:funcSigLength] == funcSig
			}),
			mock.Anything).
		Return(nil, err).After(after)
}

func (receiver ContractMockReceiver) mustEncodeResponse(funcName string, responseArgs ...interface{}) []byte {
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
