package evm_test

//go:generate ./chainlink_reader_test_setup.sh

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/types/interfacetests"
	"github.com/smartcontractkit/libocr/commontypes"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const anyMethodName = "method"

const commonGasLimitOnEvms = uint64(4712388)

var inner = []abi.ArgumentMarshaling{
	{Name: "I", Type: "int64"},
	{Name: "S", Type: "string"},
}

var nested = []abi.ArgumentMarshaling{
	{Name: "FixedBytes", Type: "bytes2"},
	{Name: "Inner", Type: "tuple", Components: inner},
}

const sizeItemType = "item for size"

var defs = map[string][]abi.ArgumentMarshaling{
	interfacetests.TestItemType: {
		{Name: "Field", Type: "int32"},
		{Name: "DifferentField", Type: "string"},
		{Name: "OracleId", Type: "uint8"},
		{Name: "OracleIds", Type: "uint8[32]"},
		{Name: "Account", Type: "bytes32"},
		{Name: "Accounts", Type: "bytes32[]"},
		{Name: "BigField", Type: "int192"},
		{Name: "NestedStruct", Type: "tuple", Components: nested},
	},
	interfacetests.TestItemSliceType: {
		{Name: "Field", Type: "int32[]"},
		{Name: "DifferentField", Type: "string[]"},
		{Name: "OracleId", Type: "uint8[]"},
		{Name: "OracleIds", Type: "bytes32[]"},
		{Name: "Account", Type: "bytes32[]"},
		{Name: "Accounts", Type: "bytes32[][]"},
		{Name: "BigField", Type: "int192[]"},
		{Name: "NestedStruct", Type: "tuple[]", Components: nested},
	},
	interfacetests.TestItemArray1Type: {
		{Name: "Field", Type: "int32[1]"},
		{Name: "DifferentField", Type: "string[1]"},
		{Name: "OracleId", Type: "uint8[1]"},
		{Name: "OracleIds", Type: "bytes32[1]"},
		{Name: "Account", Type: "bytes32[1]"},
		{Name: "Accounts", Type: "bytes32[][1]"},
		{Name: "BigField", Type: "int192[1]"},
		{Name: "NestedStruct", Type: "tuple[1]", Components: nested},
	},
	interfacetests.TestItemArray2Type: {
		{Name: "Field", Type: "int32[2]"},
		{Name: "DifferentField", Type: "string[2]"},
		{Name: "OracleId", Type: "uint8[2]"},
		{Name: "OracleIds", Type: "bytes32[2]"},
		{Name: "Account", Type: "bytes32[2]"},
		{Name: "Accounts", Type: "bytes32[][2]"},
		{Name: "BigField", Type: "int192[2]"},
		{Name: "NestedStruct", Type: "tuple[2]", Components: nested},
	},
	sizeItemType: {
		{Name: "Stuff", Type: "int256[]"},
		{Name: "OtherStuff", Type: "int256"},
	},
}

func TestChainReader(t *testing.T) {
	interfacetests.RunChainReaderInterfaceTests(t, &interfaceTester{})
	t.Run("GetMaxEncodingSize delegates to GetMaxSize", func(t *testing.T) {
		runSizeDelegationTest(t, func(reader relaytypes.ChainReader, ctx context.Context, i int, s string) (int, error) {
			return reader.GetMaxEncodingSize(ctx, i, s)
		})
	})

	t.Run("GetMaxDecodingSize delegates to GetMaxSize", func(t *testing.T) {
		runSizeDelegationTest(t, func(reader relaytypes.ChainReader, ctx context.Context, i int, s string) (int, error) {
			return reader.GetMaxDecodingSize(ctx, i, s)
		})
	})
}

type interfaceTester struct {
	chain    *mocks.Chain
	contract relaytypes.BoundContract
	ropts    *types.RelayOpts
	defs     map[string]abi.Arguments
	auth     *bind.TransactOpts
	sim      *backends.SimulatedBackend
	pk       *ecdsa.PrivateKey
	evmTest  *EvmTest
}

func (it *interfaceTester) Setup(t *testing.T) {
	// can re-use the same chain for tests, just make new contract for each test
	if it.chain == nil {
		defBytes, err := json.Marshal(defs)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(defBytes, &it.defs))
		it.chain = &mocks.Chain{}
		it.setupChain(t)
		it.chain.On("LogPoller").Return(lpmocks.NewLogPoller(t))

		relayConfig := types.RelayConfig{
			ChainReader: &types.ChainReaderConfig{
				ChainContractReaders: map[string]types.ChainContractReader{
					"LatestValueHolder": {
						ContractABI: EvmTestMetaData.ABI,
						ChainReaderDefinitions: map[string]types.ChainReaderDefinition{
							anyMethodName: {
								ChainSpecificName: "GetElementAtIndex",
								ReturnValues: []string{
									"Field",
									"DifferentField",
									"OracleId",
									"OracleIds",
									"Account",
									"Accounts",
									"BigField",
									"NestedStruct",
								},
							},
						},
					},
				},
				ChainCodecConfigs: map[string]types.ChainCodedConfig{},
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

func (it *interfaceTester) Teardown(_ *testing.T) {
	it.contract.Address = ""
}

func (it *interfaceTester) Name() string {
	return "EVM"
}

func (it *interfaceTester) EncodeFields(t *testing.T, request *interfacetests.EncodeRequest) ocr2types.Report {
	if request.TestOn == interfacetests.TestItemType {
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

func (it *interfaceTester) SetLatestValue(t *testing.T, ctx context.Context, testStruct *interfacetests.TestStruct) (relaytypes.BoundContract, string) {
	// Since most tests don't use the contract, it's set up lazily to save time
	if it.contract.Address == "" {
		it.contract.Address = it.deployNewContract(t, ctx)
		it.contract.Name = anyMethodName
		it.contract.Pending = false
	}

	tx, err := it.evmTest.AddTestStruct(
		it.auth,
		testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleId),
		convertOracleIds(testStruct.OracleIds),
		[32]byte(testStruct.Account),
		convertAccounts(testStruct.Accounts),
		testStruct.BigField,
		toInternalType(testStruct.NestedStruct),
	)
	require.NoError(t, err)
	it.sim.Commit()
	it.incNonce()
	it.awaitTx(t, ctx, tx)
	return it.contract, anyMethodName
}

func (it *interfaceTester) encodeFieldsOnItem(t *testing.T, request *interfacetests.EncodeRequest) ocr2types.Report {
	return packArgs(t, argsFromTestStruct(request.TestStructs[0]), it.defs[interfacetests.TestItemType], request)
}

func (it *interfaceTester) encodeFieldsOnSliceOrArray(t *testing.T, request *interfacetests.EncodeRequest) []byte {
	oargs := it.defs[request.TestOn]

	var field []int32
	var differentField []string
	var oracleId []byte
	var oracleIds [][32]byte
	var account [][32]byte
	var accounts [][][32]byte
	var bigField []*big.Int
	var nested []MidLevelTestStruct

	for _, testStruct := range request.TestStructs {
		field = append(field, testStruct.Field)
		differentField = append(differentField, testStruct.DifferentField)
		oracleId = append(oracleId, byte(testStruct.OracleId))
		convertedIds := convertOracleIds(testStruct.OracleIds)

		convertedAccounts := convertAccounts(testStruct.Accounts)

		oracleIds = append(oracleIds, convertedIds)
		account = append(account, [32]byte(testStruct.Account))
		accounts = append(accounts, convertedAccounts)
		bigField = append(bigField, testStruct.BigField)
		nested = append(nested, toInternalType(testStruct.NestedStruct))
	}

	allArgs := []any{field, differentField, oracleId, oracleIds, account, accounts, bigField, nested}

	switch request.TestOn {
	case interfacetests.TestItemArray1Type:
		for i, arg := range allArgs {
			allArgs[i] = toFixedSized(arg, 1)
		}
	case interfacetests.TestItemArray2Type:
		for i, arg := range allArgs {
			allArgs[i] = toFixedSized(arg, 2)
		}
	}

	return packArgs(t, allArgs, oargs, request)
}

func convertOracleIds(oracleIds [32]commontypes.OracleID) [32]byte {
	convertedIds := [32]byte{}
	for i, id := range oracleIds {
		convertedIds[i] = byte(id)
	}
	return convertedIds
}

func convertAccounts(accounts [][]byte) [][32]byte {
	convertedAccounts := make([][32]byte, len(accounts))
	for i, a := range accounts {
		convertedAccounts[i] = [32]byte(a)
	}
	return convertedAccounts
}

func (it *interfaceTester) setupChain(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	it.pk = privateKey

	it.auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	it.sim = backends.NewSimulatedBackend(core.GenesisAlloc{it.auth.From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
	it.sim.Commit()
	it.chain.On("Client").Return(client.NewSimulatedBackendClient(t, it.sim, big.NewInt(1337)))
}

func (it *interfaceTester) deployNewContract(t *testing.T, ctx context.Context) string {
	gasPrice, err := it.sim.SuggestGasPrice(ctx)
	require.NoError(t, err)
	it.auth.GasPrice = gasPrice

	// 105528 was in the error: gas too low: have 0, want 105528
	// Not sure if there's a better way to get it.
	it.auth.GasLimit = 1055280000

	address, tx, ts, err := DeployEvmTest(it.auth, it.sim)

	require.NoError(t, err)
	it.sim.Commit()
	it.evmTest = ts
	it.incNonce()
	it.awaitTx(t, ctx, tx)
	return address.String()
}

func (it *interfaceTester) awaitTx(t *testing.T, ctx context.Context, tx *evmtypes.Transaction) {
	receipt, err := it.sim.TransactionReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, evmtypes.ReceiptStatusSuccessful, receipt.Status)
}

func (it *interfaceTester) incNonce() {
	if it.auth.Nonce == nil {
		it.auth.Nonce = big.NewInt(1)
	} else {
		it.auth.Nonce = it.auth.Nonce.Add(it.auth.Nonce, big.NewInt(1))
	}
}

func packArgs(t *testing.T, allArgs []any, oargs abi.Arguments, request *interfacetests.EncodeRequest) []byte {
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

func getAccounts(first interfacetests.TestStruct) [][32]byte {
	accountBytes := make([][32]byte, len(first.Accounts))
	for i, account := range first.Accounts {
		accountBytes[i] = [32]byte(account)
	}
	return accountBytes
}

func argsFromTestStruct(ts interfacetests.TestStruct) []any {
	return []any{
		ts.Field,
		ts.DifferentField,
		uint8(ts.OracleId),
		getOracleIds(ts),
		[32]byte(ts.Account),
		getAccounts(ts),
		ts.BigField,
		toInternalType(ts.NestedStruct),
	}
}

func getOracleIds(first interfacetests.TestStruct) [32]byte {
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

func toInternalType(m interfacetests.MidLevelTestStruct) MidLevelTestStruct {
	return MidLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: InnerTestStruct{
			I: int64(m.Inner.I),
			S: m.Inner.S,
		},
	}
}

func runSizeDelegationTest(t *testing.T, run func(relaytypes.ChainReader, context.Context, int, string) (int, error)) {
	it := &interfaceTester{}
	it.Setup(t)

	cr := it.GetChainReader(t)

	ctx := context.Background()
	actual, err := run(cr, ctx, 10, sizeItemType)
	require.NoError(t, err)

	expected, _ := evm.GetMaxSize(10, it.defs[sizeItemType])
	assert.Equal(t, expected, actual)
}
