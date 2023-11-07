package evm_test

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	. "github.com/smartcontractkit/chainlink-relay/pkg/types/interfacetests"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	chainevm "github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const anyMethodName = "method"

var inner = []abi.ArgumentMarshaling{
	{Name: "I", Type: "int64"},
	{Name: "S", Type: "string"},
}

var nested = []abi.ArgumentMarshaling{
	{Name: "FixedBytes", Type: "bytes2"},
	{Name: "Inner", Type: "tuple", Components: inner},
}

var defs = map[string][]abi.ArgumentMarshaling{
	TestItemType: {
		{Name: "Field", Type: "int32"},
		{Name: "DifferentField", Type: "string"},
		{Name: "OracleId", Type: "uint8"},
		{Name: "OracleIds", Type: "uint8[32]"},
		{Name: "Account", Type: "bytes32"},
		{Name: "Accounts", Type: "bytes32[]"},
		{Name: "BigField", Type: "int192"},
		{Name: "NestedStruct", Type: "tuple", Components: nested},
	},
	TestItemSliceType: {
		{Name: "Field", Type: "int32[]"},
		{Name: "DifferentField", Type: "string[]"},
		{Name: "OracleId", Type: "uint8[]"},
		{Name: "OracleIds", Type: "bytes32[]"},
		{Name: "Account", Type: "bytes32[]"},
		{Name: "Accounts", Type: "bytes32[][]"},
		{Name: "BigField", Type: "int192[]"},
		{Name: "NestedStruct", Type: "tuple[]", Components: nested},
	},
	TestItemArray1Type: {
		{Name: "Field", Type: "int32[1]"},
		{Name: "DifferentField", Type: "string[1]"},
		{Name: "OracleId", Type: "uint8[1]"},
		{Name: "OracleIds", Type: "bytes32[1]"},
		{Name: "Account", Type: "bytes32[1]"},
		{Name: "Accounts", Type: "bytes32[][1]"},
		{Name: "BigField", Type: "int192[1]"},
		{Name: "NestedStruct", Type: "tuple[1]", Components: nested},
	},
	TestItemArray2Type: {
		{Name: "Field", Type: "int32[2]"},
		{Name: "DifferentField", Type: "string[2]"},
		{Name: "OracleId", Type: "uint8[2]"},
		{Name: "OracleIds", Type: "bytes32[2]"},
		{Name: "Account", Type: "bytes32[2]"},
		{Name: "Accounts", Type: "bytes32[][2]"},
		{Name: "BigField", Type: "int192[2]"},
		{Name: "NestedStruct", Type: "tuple[2]", Components: nested},
	},
}

func TestChainReader(t *testing.T) {
	RunChainReaderInterfaceTests(t, &interfaceTester{})
}

type interfaceTester struct {
	chain    chainevm.Chain
	contract relaytypes.BoundContract
	ropts    *types.RelayOpts
	defs     map[string]abi.Arguments
}

func (it *interfaceTester) Setup(t *testing.T) {
	// can re-use the same chain for tests, just make new contract for each test
	if it.chain == nil {
		defBytes, err := json.Marshal(defs)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(defBytes, &it.defs))

		// TODO would like to set up a real EVM here, but tests I see use a mock...
		c := evmclimocks.NewClient(t)
		chainMock := &mocks.Chain{}
		it.chain = chainMock
		chainMock.On("Client").Return(c)
		chainMock.On("LogPoller").Return(lpmocks.NewLogPoller(t))

		relayConfig := types.RelayConfig{
			ChainReader: &types.ChainReaderConfig{
				ChainContractReaders: map[string]types.ChainContractReader{},
				ChainCodecConfigs:    map[string]types.ChainCodedConfig{},
			},
		}

		for k, v := range defs {
			defBytes, err := json.Marshal(v)
			require.NoError(t, err)
			entry := relayConfig.ChainReader.ChainCodecConfigs[k]
			entry.TypeAbi = string(defBytes)
			relayConfig.ChainReader.ChainCodecConfigs[k] = entry
		}

		relayBytes, err := json.Marshal(relayConfig)
		require.NoError(t, err)
		it.ropts = &types.RelayOpts{
			RelayArgs: relaytypes.RelayArgs{RelayConfig: relayBytes},
		}
	}
}

func (it *interfaceTester) Teardown(t *testing.T) {
	it.contract.Address = ""
}

func (it *interfaceTester) Name() string {
	return "EVM"
}

func (it *interfaceTester) EncodeFields(t *testing.T, request *EncodeRequest) ocr2types.Report {
	if request.TestOn == TestItemType {
		return it.encodeFieldsOnItem(t, request)
	}
	return it.encodeFieldsOnSliceOrArray(t, request)
}

func (it *interfaceTester) GetAccountBytes(i int) []byte {
	account := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}
	account[i%32] += byte(i)
	account[(i+3)%32] += byte(i + 3)
	return account
}

func (it *interfaceTester) GetChainReader(t *testing.T) relaytypes.ChainReader {
	cr, err := evm.NewChainReader(logger.TestLogger(t), it.chain, it.ropts)
	require.NoError(t, err)
	return cr
}

func (it *interfaceTester) IncludeArrayEncodingSizeEnforcement() bool {
	return true
}

func (it *interfaceTester) SetLatestValue(t *testing.T, testStruct *TestStruct) (relaytypes.BoundContract, string) {
	// Since most tests don't use the contract, it's set up lazily to save time
	if it.contract.Address == "" {
		// TODO set up the contract and get the address back
		// it.contract.Address =
		it.contract.Name = anyMethodName
		// TODO tests for pending?
		it.contract.Pending = false
	}
	return it.contract, anyMethodName
}

func (it *interfaceTester) encodeFieldsOnItem(t *testing.T, request *EncodeRequest) ocr2types.Report {
	first := request.TestStructs[0]

	allArgs := []any{
		first.Field,
		first.DifferentField,
		uint8(first.OracleId),
		getOracleIds(first),
		[32]byte(first.Account),
		getAccounts(first),
		first.BigField,
		toInternalType(first.NestedStruct),
	}
	return packArgs(t, allArgs, it.defs[TestItemType], request)
}

func (it *interfaceTester) encodeFieldsOnSliceOrArray(t *testing.T, request *EncodeRequest) []byte {
	oargs := it.defs[request.TestOn]

	var field []int32
	var differentField []string
	var oracleId []byte
	var oracleIds [][32]byte
	var account [][32]byte
	var accounts [][][32]byte
	var bigField []*big.Int
	var nested []midLevelTestStruct

	for _, testStruct := range request.TestStructs {
		field = append(field, testStruct.Field)
		differentField = append(differentField, testStruct.DifferentField)
		oracleId = append(oracleId, byte(testStruct.OracleId))
		convertedIds := [32]byte{}
		for i, id := range testStruct.OracleIds {
			convertedIds[i] = byte(id)
		}
		convertedAccount := [32]byte(testStruct.Account)

		convertedAccounts := make([][32]byte, len(testStruct.Accounts))
		for i, a := range testStruct.Accounts {
			convertedAccounts[i] = [32]byte(a)
		}

		oracleIds = append(oracleIds, convertedIds)
		account = append(account, convertedAccount)
		accounts = append(accounts, convertedAccounts)
		bigField = append(bigField, testStruct.BigField)
		nested = append(nested, toInternalType(testStruct.NestedStruct))
	}

	allArgs := []any{field, differentField, oracleId, oracleIds, account, accounts, bigField, nested}

	switch request.TestOn {
	case TestItemArray1Type:
		for i, arg := range allArgs {
			allArgs[i] = toFixedSized(arg, 1)
		}
	case TestItemArray2Type:
		for i, arg := range allArgs {
			allArgs[i] = toFixedSized(arg, 2)
		}
	}

	return packArgs(t, allArgs, oargs, request)
}

func packArgs(t *testing.T, allArgs []any, oargs abi.Arguments, request *EncodeRequest) []byte {
	// extra capacity in case we add an argument
	args := make(abi.Arguments, len(oargs), len(oargs)+1)
	copy(args, oargs)
	// decoding has extra field to decode
	if request.ExtraField {
		fakeType, err := abi.NewType("int32", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		args = append(args, abi.Argument{Name: "FakeField", Type: fakeType})
		allArgs = append(allArgs)
	}

	if request.MissingField {
		args = args[1:]
		allArgs = allArgs[1:]
	}

	bytes, err := args.Pack(allArgs...)
	require.NoError(t, err)
	return bytes
}

func getAccounts(first TestStruct) [][32]byte {
	accountBytes := make([][32]byte, len(first.Accounts))
	for i, account := range first.Accounts {
		accountBytes[i] = [32]byte(account)
	}
	return accountBytes
}

func getOracleIds(first TestStruct) [32]byte {
	oracleIds := [32]byte{}
	for i, oracleId := range first.OracleIds {
		oracleIds[i] = byte(oracleId)
	}
	return oracleIds
}

func toFixedSized(item any, size int) any {
	rItem := reflect.ValueOf(item)
	arrayType := reflect.ArrayOf(size, rItem.Type().Elem())
	return rItem.Convert(arrayType).Interface()
}

func toInternalType(m MidLevelTestStruct) midLevelTestStruct {
	return midLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: innerTestStruct{
			I: int64(m.Inner.I),
			S: m.Inner.S,
		},
	}
}

type innerTestStruct struct {
	I int64
	S string
}

type midLevelTestStruct struct {
	FixedBytes [2]byte
	Inner      innerTestStruct
}
