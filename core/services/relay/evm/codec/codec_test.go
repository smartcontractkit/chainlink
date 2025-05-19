package codec_test

import (
	"encoding/json"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	looptestutils "github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/chain_reader_tester"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint:revive // dot-imports
)

const anyExtraValue = 3

func TestCodec(t *testing.T) {
	tester := &codecInterfaceTester{}
	RunCodecInterfaceTests(t, tester)
	RunCodecInterfaceTests(t, looptestutils.WrapCodecTesterForLoop(tester))

	anyN := 10
	c := tester.GetCodec(t)
	t.Run("Decode works with multiple unnamed return values", func(t *testing.T) {
		encode := &struct {
			F0 int32
			F1 int32
		}{F0: 1, F1: 2}
		codecName := "my_codec"
		evmEncoderConfig := `[{"Name":"","Type":"int32"},{"Name":"","Type":"int32"}]`

		codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
			codecName: {TypeABI: evmEncoderConfig},
		}}
		c, err := codec.NewCodec(codecConfig)
		require.NoError(t, err)

		result, err := c.Encode(testutils.Context(t), encode, codecName)
		require.NoError(t, err)

		decode := &struct {
			F0 int32
			F1 int32
		}{}
		err = c.Decode(testutils.Context(t), result, decode, codecName)
		require.NoError(t, err)
		require.Equal(t, encode.F0, decode.F0)
		require.Equal(t, encode.F1, decode.F1)
	})

	t.Run("GetMaxEncodingSize delegates to GetMaxSize", func(t *testing.T) {
		actual, err := c.GetMaxEncodingSize(testutils.Context(t), anyN, sizeItemType)
		assert.NoError(t, err)

		expected, err := types.GetMaxSize(anyN, parseDefs(t)[sizeItemType])
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("GetMaxDecodingSize delegates to GetMaxSize", func(t *testing.T) {
		actual, err := c.GetMaxDecodingSize(testutils.Context(t), anyN, sizeItemType)
		assert.NoError(t, err)

		expected, err := types.GetMaxSize(anyN, parseDefs(t)[sizeItemType])
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestCodec_SimpleEncode(t *testing.T) {
	codecName := "my_codec"
	input := map[string]any{
		"Report": int32(6),
		"Meta":   "abcdefg",
	}
	evmEncoderConfig := `[{"Name":"Report","Type":"int32"},{"Name":"Meta","Type":"string"}]`

	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		codecName: {TypeABI: evmEncoderConfig},
	}}
	c, err := codec.NewCodec(codecConfig)
	require.NoError(t, err)

	result, err := c.Encode(testutils.Context(t), input, codecName)
	require.NoError(t, err)
	expected :=
		"0000000000000000000000000000000000000000000000000000000000000006" + // int32(6)
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of the beginning of second value (64 bytes)
			"0000000000000000000000000000000000000000000000000000000000000007" + // length of the string (7 chars)
			"6162636465666700000000000000000000000000000000000000000000000000" // actual string

	require.Equal(t, expected, hexutil.Encode(result)[2:])
}

func TestCodec_EncodeTuple(t *testing.T) {
	codecName := "my_codec"
	input := map[string]any{
		"Report": int32(6),
		"Nested": map[string]any{
			"Meta":  "abcdefg",
			"Count": int32(14),
			"Other": "12334",
		},
	}
	evmEncoderConfig := `[{"Name":"Report","Type":"int32"},{"Name":"Nested","Type":"tuple","Components":[{"Name":"Other","Type":"string"},{"Name":"Count","Type":"int32"},{"Name":"Meta","Type":"string"}]}]`

	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		codecName: {TypeABI: evmEncoderConfig},
	}}
	c, err := codec.NewCodec(codecConfig)
	require.NoError(t, err)

	result, err := c.Encode(testutils.Context(t), input, codecName)
	require.NoError(t, err)
	expected :=
		"0000000000000000000000000000000000000000000000000000000000000006" + // Report integer (=6)
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of the first dynamic value (tuple, 64 bytes)
			"0000000000000000000000000000000000000000000000000000000000000060" + // offset of the first nested dynamic value (string, 96 bytes)
			"000000000000000000000000000000000000000000000000000000000000000e" + // "Count" integer (=14)
			"00000000000000000000000000000000000000000000000000000000000000a0" + // offset of the second nested dynamic value (string, 160 bytes)
			"0000000000000000000000000000000000000000000000000000000000000005" + // length of the "Meta" string (5 chars)
			"3132333334000000000000000000000000000000000000000000000000000000" + // "Other" string (="12334")
			"0000000000000000000000000000000000000000000000000000000000000007" + // length of the "Other" string (7 chars)
			"6162636465666700000000000000000000000000000000000000000000000000" // "Meta" string (="abcdefg")

	require.Equal(t, expected, hexutil.Encode(result)[2:])
}

func TestCodec_EncodeTupleWithLists(t *testing.T) {
	codecName := "my_codec"
	input := map[string]any{
		"Elem": map[string]any{
			"Prices":     []any{big.NewInt(234), big.NewInt(456)},
			"Timestamps": []any{int64(111), int64(222)},
		},
	}
	evmEncoderConfig := `[{"Name":"Elem","Type":"tuple","InternalType":"tuple","Components":[{"Name":"Prices","Type":"uint256[]","InternalType":"uint256[]","Components":null,"Indexed":false},{"Name":"Timestamps","Type":"uint32[]","InternalType":"uint32[]","Components":null,"Indexed":false}],"Indexed":false}]`

	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		codecName: {TypeABI: evmEncoderConfig},
	}}
	c, err := codec.NewCodec(codecConfig)
	require.NoError(t, err)

	result, err := c.Encode(testutils.Context(t), input, codecName)
	require.NoError(t, err)
	expected :=
		"0000000000000000000000000000000000000000000000000000000000000020" + // offset of Elem tuple
			"0000000000000000000000000000000000000000000000000000000000000040" + // offset of Prices array
			"00000000000000000000000000000000000000000000000000000000000000a0" + // offset of Timestamps array
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Prices array
			"00000000000000000000000000000000000000000000000000000000000000ea" + // Prices[0] = 234
			"00000000000000000000000000000000000000000000000000000000000001c8" + // Prices[1] = 456
			"0000000000000000000000000000000000000000000000000000000000000002" + // length of Timestamps array
			"000000000000000000000000000000000000000000000000000000000000006f" + // Timestamps[0] = 111
			"00000000000000000000000000000000000000000000000000000000000000de" // Timestamps[1] = 222

	require.Equal(t, expected, hexutil.Encode(result)[2:])
}

type codecInterfaceTester struct {
	TestSelectionSupport
}

func (it *codecInterfaceTester) Setup(_ *testing.T) {}

func (it *codecInterfaceTester) GetAccountBytes(i int) []byte {
	account := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// fuzz tests can make -ve numbers
	if i < 0 {
		i = -i
	}
	account[i%20] += byte(i)
	account[(i+3)%20] += byte(i + 3)
	return account[:]
}

func (it *codecInterfaceTester) GetAccountString(i int) string {
	return common.BytesToAddress(it.GetAccountBytes(i)).Hex()
}

func (it *codecInterfaceTester) EncodeFields(t *testing.T, request *EncodeRequest) []byte {
	if request.TestOn == TestItemType {
		return encodeFieldsOnItem(t, request)
	}

	return encodeFieldsOnSliceOrArray(t, request)
}

func (it *codecInterfaceTester) GetCodec(t *testing.T) commontypes.Codec {
	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{}}
	testStruct := CreateTestStruct[*testing.T](0, it)
	for k, v := range codecDefs {
		defBytes, err := json.Marshal(v)
		require.NoError(t, err)
		entry := codecConfig.Configs[k]
		entry.TypeABI = string(defBytes)

		if k != sizeItemType && k != NilType {
			entry.ModifierConfigs = commoncodec.ModifiersConfig{
				&commoncodec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
				&commoncodec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
			}
		}

		if slices.Contains([]string{TestItemType, TestItemSliceType, TestItemArray1Type, TestItemArray2Type, TestItemWithConfigExtra}, k) {
			addressByteModifier := &commoncodec.AddressBytesToStringModifierConfig{
				Fields:   []string{"AccountStruct.AccountStr"},
				Modifier: codec.EVMAddressModifier{},
			}

			entry.ModifierConfigs = append(entry.ModifierConfigs, addressByteModifier)
		}

		if k == TestItemWithConfigExtra {
			hardCode := &commoncodec.HardCodeModifierConfig{
				OnChainValues: map[string]any{
					"BigField":              testStruct.BigField.String(),
					"AccountStruct.Account": hexutil.Encode(testStruct.AccountStruct.Account),
				},
				OffChainValues: map[string]any{"ExtraField": anyExtraValue},
			}
			entry.ModifierConfigs = append(entry.ModifierConfigs, hardCode)
		}
		codecConfig.Configs[k] = entry
	}

	c, err := codec.NewCodec(codecConfig)
	require.NoError(t, err)
	return c
}

func (it *codecInterfaceTester) IncludeArrayEncodingSizeEnforcement() bool {
	return true
}
func (it *codecInterfaceTester) Name() string {
	return "EVM"
}

func encodeFieldsOnItem(t *testing.T, request *EncodeRequest) ocr2types.Report {
	return packArgs(t, argsFromTestStruct(request.TestStructs[0]), parseDefs(t)[TestItemType], request)
}

func encodeFieldsOnSliceOrArray(t *testing.T, request *EncodeRequest) []byte {
	oargs := parseDefs(t)[request.TestOn]
	args := make([]any, 1)

	switch request.TestOn {
	case TestItemArray1Type:
		args[0] = [1]chain_reader_tester.TestStruct{evmtesting.ToInternalType(request.TestStructs[0])}
	case TestItemArray2Type:
		args[0] = [2]chain_reader_tester.TestStruct{evmtesting.ToInternalType(request.TestStructs[0]), evmtesting.ToInternalType(request.TestStructs[1])}
	default:
		tmp := make([]chain_reader_tester.TestStruct, len(request.TestStructs))
		for i, ts := range request.TestStructs {
			tmp[i] = evmtesting.ToInternalType(ts)
		}
		args[0] = tmp
	}

	return packArgs(t, args, oargs, request)
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
		allArgs = append(allArgs, 11)
	}

	if request.MissingField {
		args = args[1:]
		allArgs = allArgs[1:]
	}

	bytes, err := args.Pack(allArgs...)
	require.NoError(t, err)
	return bytes
}

var innerDynamic = []abi.ArgumentMarshaling{
	{Name: "IntVal", Type: "int64"},
	{Name: "S", Type: "string"},
}

var nestedDynamic = []abi.ArgumentMarshaling{
	{Name: "FixedBytes", Type: "bytes2"},
	{Name: "Inner", Type: "tuple", Components: innerDynamic},
}

var innerStatic = []abi.ArgumentMarshaling{
	{Name: "IntVal", Type: "int64"},
	{Name: "A", Type: "address"},
}

var nestedStatic = []abi.ArgumentMarshaling{
	{Name: "FixedBytes", Type: "bytes2"},
	{Name: "Inner", Type: "tuple", Components: innerStatic},
}

var accountStruct = []abi.ArgumentMarshaling{
	{Name: "Account", Type: "address"},
	{Name: "AccountStr", Type: "address"},
}

var ts = []abi.ArgumentMarshaling{
	{Name: "Field", Type: "int32"},
	{Name: "DifferentField", Type: "string"},
	{Name: "OracleId", Type: "uint8"},
	{Name: "OracleIds", Type: "uint8[32]"},
	{Name: "AccountStruct", Type: "tuple", Components: accountStruct},
	{Name: "Accounts", Type: "address[]"},
	{Name: "BigField", Type: "int192"},
	{Name: "NestedDynamicStruct", Type: "tuple", Components: nestedDynamic},
	{Name: "NestedStaticStruct", Type: "tuple", Components: nestedStatic},
}

const sizeItemType = "item for size"

var codecDefs = map[string][]abi.ArgumentMarshaling{
	TestItemType: ts,
	TestItemSliceType: {
		{Name: "", Type: "tuple[]", Components: ts},
	},
	TestItemArray1Type: {
		{Name: "", Type: "tuple[1]", Components: ts},
	},
	TestItemArray2Type: {
		{Name: "", Type: "tuple[2]", Components: ts},
	},
	sizeItemType: {
		{Name: "Stuff", Type: "int256[]"},
		{Name: "OtherStuff", Type: "int256"},
	},
	TestItemWithConfigExtra: ts,
	NilType:                 {},
}

func parseDefs(t *testing.T) map[string]abi.Arguments {
	bytes, err := json.Marshal(codecDefs)
	require.NoError(t, err)
	var results map[string]abi.Arguments
	require.NoError(t, json.Unmarshal(bytes, &results))
	return results
}

func getAccounts(first TestStruct) []common.Address {
	accountBytes := make([]common.Address, len(first.Accounts))
	for i, account := range first.Accounts {
		accountBytes[i] = common.Address(account)
	}
	return accountBytes
}

func getOracleIDs(first TestStruct) [32]byte {
	oracleIDs := [32]byte{}
	for i, oracleID := range first.OracleIDs {
		oracleIDs[i] = byte(oracleID)
	}
	return oracleIDs
}

func argsFromTestStruct(ts TestStruct) []any {
	return []any{
		ts.Field,
		ts.DifferentField,
		uint8(ts.OracleID),
		getOracleIDs(ts),
		evmtesting.AccountStructToInternalType(ts.AccountStruct),
		getAccounts(ts),
		ts.BigField,
		evmtesting.MidDynamicToInternalType(ts.NestedDynamicStruct),
		evmtesting.MidStaticToInternalType(ts.NestedStaticStruct),
	}
}
