package evm_test

//go:generate ./testfiles/chainlink_reader_test_setup.sh

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	. "github.com/smartcontractkit/chainlink-relay/pkg/types/interfacetests"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/testfiles"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	relaymedian "github.com/smartcontractkit/chainlink/v2/core/services/relay/median"
)

const commonGasLimitOnEvms = uint64(4712388)

var inner = []abi.ArgumentMarshaling{
	{Name: "I", Type: "int64"},
	{Name: "S", Type: "string"},
}

var nested = []abi.ArgumentMarshaling{
	{Name: "FixedBytes", Type: "bytes2"},
	{Name: "Inner", Type: "tuple", Components: inner},
}

var ts = []abi.ArgumentMarshaling{
	{Name: "Field", Type: "int32"},
	{Name: "DifferentField", Type: "string"},
	{Name: "OracleId", Type: "uint8"},
	{Name: "OracleIds", Type: "uint8[32]"},
	{Name: "Account", Type: "bytes32"},
	{Name: "Accounts", Type: "bytes32[]"},
	{Name: "BigField", Type: "int192"},
	{Name: "NestedStruct", Type: "tuple", Components: nested},
}

const sizeItemType = "item for size"

var defs = map[string][]abi.ArgumentMarshaling{
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
	relaymedian.MedianTypeName: {
		{Name: "observationsTimestamp", Type: "uint32"},
		{Name: "rawObservers", Type: "bytes32"},
		{Name: "observations", Type: "int192[]"},
		{Name: "juelsPerFeeCoin", Type: "int192"},
	},
}

func TestChainReader(t *testing.T) {
	RunChainReaderInterfaceTests(t, &interfaceTester{})
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

	t.Run("ReportCodec compatibility tests", func(t *testing.T) {
		it := &interfaceTester{}
		it.Setup(t)

		cr := it.GetChainReader(t)

		pao := median.ParsedAttributedObservation{
			math.MaxUint32,
			median.MaxValue(),
			median.MaxValue(),
			ocr2types.MaxOracles - 1,
		}

		var paos []median.ParsedAttributedObservation

		rc, err := relaymedian.NewReportCodec(cr)
		require.NoError(t, err)
		for n := 1; n <= ocr2types.MaxOracles; n++ {
			paos = append(paos, pao)
			maxReportLen, err := rc.MaxReportLength(n)
			require.NoError(t, err)
			report, err := rc.BuildReport(paos)
			require.NoError(t, err)
			require.Equal(t, len(report), maxReportLen)
		}
	})
}

func FuzzMedianFromReportCompatibility(f *testing.F) {
	it := &interfaceTester{}
	it.setupNoClient(f)

	cr := it.getChainReader(f)

	codec, err := relaymedian.NewReportCodec(cr)
	require.NoError(f, err)
	validReport1, err := codec.BuildReport([]median.ParsedAttributedObservation{{
		12345678,
		big.NewInt(1e12),
		big.NewInt(1e13),
		0,
	}})
	if err != nil {
		f.Fatalf("failed to construct valid report: %s", err)
	}

	validReport2, err := codec.BuildReport([]median.ParsedAttributedObservation{{
		12345678,
		big.NewInt(1e12),
		big.NewInt(1e13),
		0,
	}, {
		12345679,
		big.NewInt(1e13),
		big.NewInt(1e14),
		1,
	}})
	if err != nil {
		f.Fatalf("failed to construct valid report: %s", err)
	}

	f.Add([]byte{})
	f.Add([]byte(validReport1))
	f.Add([]byte(validReport2))
	f.Fuzz(func(t *testing.T, report []byte) {
		_, _ = codec.MedianFromReport(report)
	})
}

type interfaceTester struct {
	chain   *mocks.Chain
	address string
	ropts   *types.RelayOpts
	defs    map[string]abi.Arguments
	auth    *bind.TransactOpts
	sim     *backends.SimulatedBackend
	pk      *ecdsa.PrivateKey
	evmTest *testfiles.Testfiles
}

func (it *interfaceTester) Setup(t *testing.T) {
	if it.chain == nil {
		it.setupNoClient(t)
		it.chain.On("Client").Return(client.NewSimulatedBackendClient(t, it.sim, big.NewInt(1337)))
	}
}

func (it *interfaceTester) setupNoClient(t require.TestingT) {
	// can re-use the same chain for tests, just make new contract for each test
	if it.chain != nil {
		return
	}

	defBytes, err := json.Marshal(defs)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(defBytes, &it.defs))
	it.chain = &mocks.Chain{}
	it.setupChainNoClient(t)
	it.chain.On("LogPoller").Return(logger.NullLogger)

	relayConfig := types.RelayConfig{
		ChainReader: &types.ChainReaderConfig{
			ChainContractReaders: map[string]types.ChainContractReader{
				"LatestValueHolder": {
					ContractABI: testfiles.TestfilesMetaData.ABI,
					ChainReaderDefinitions: map[string]types.ChainReaderDefinition{
						MethodTakingLatestParamsReturningTestStruct: {
							ChainSpecificName: "GetElementAtIndex",
						},
						MethodReturningUint64: {
							ChainSpecificName: "GetPrimitiveValue",
						},
						MethodReturningUint64Slice: {
							ChainSpecificName: "GetSliceValue",
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

func (it *interfaceTester) Teardown(_ *testing.T) {
	it.address = ""
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
	return it.getChainReader(t)
}

func (it *interfaceTester) getChainReader(t require.TestingT) relaytypes.ChainReader {
	cr, err := evm.NewChainReader(logger.NullLogger, it.chain, it.ropts)
	require.NoError(t, err)
	return cr
}

func (it *interfaceTester) IncludeArrayEncodingSizeEnforcement() bool {
	return true
}

func (it *interfaceTester) GetPrimitiveContract(ctx context.Context, t *testing.T) relaytypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(ctx, t)
	return relaytypes.BoundContract{
		Address: it.address,
		Name:    MethodReturningUint64,
	}
}

func (it *interfaceTester) GetSliceContract(ctx context.Context, t *testing.T) relaytypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(ctx, t)
	return relaytypes.BoundContract{
		Address: it.address,
		Name:    MethodReturningUint64Slice,
	}
}

func (it *interfaceTester) SetLatestValue(ctx context.Context, t *testing.T, testStruct *TestStruct) relaytypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(ctx, t)

	tx, err := it.evmTest.AddTestStruct(
		it.auth,
		testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleId),
		convertOracleIds(testStruct.OracleIds),
		[32]byte(testStruct.Account),
		convertAccounts(testStruct.Accounts),
		testStruct.BigField,
		midToInternalType(testStruct.NestedStruct),
	)
	require.NoError(t, err)
	it.sim.Commit()
	it.incNonce()
	it.awaitTx(ctx, t, tx)
	return relaytypes.BoundContract{
		Address: it.address,
		Name:    MethodTakingLatestParamsReturningTestStruct,
	}
}

func (it *interfaceTester) encodeFieldsOnItem(t *testing.T, request *EncodeRequest) ocr2types.Report {
	return packArgs(t, argsFromTestStruct(request.TestStructs[0]), it.defs[TestItemType], request)
}

func (it *interfaceTester) encodeFieldsOnSliceOrArray(t *testing.T, request *EncodeRequest) []byte {
	oargs := it.defs[request.TestOn]
	args := make([]any, 1)

	switch request.TestOn {
	case TestItemArray1Type:
		args[0] = [1]testfiles.TestStruct{toInternalType(request.TestStructs[0])}
	case TestItemArray2Type:
		args[0] = [2]testfiles.TestStruct{toInternalType(request.TestStructs[0]), toInternalType(request.TestStructs[1])}
	default:
		tmp := make([]testfiles.TestStruct, len(request.TestStructs))
		for i, ts := range request.TestStructs {
			tmp[i] = toInternalType(ts)
		}
		args[0] = tmp
	}

	return packArgs(t, args, oargs, request)
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

func (it *interfaceTester) setupChainNoClient(t require.TestingT) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	it.pk = privateKey

	it.auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	it.sim = backends.NewSimulatedBackend(core.GenesisAlloc{it.auth.From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
	it.sim.Commit()
}

func (it *interfaceTester) deployNewContract(ctx context.Context, t *testing.T) {
	if it.address != "" {
		return
	}
	gasPrice, err := it.sim.SuggestGasPrice(ctx)
	require.NoError(t, err)
	it.auth.GasPrice = gasPrice

	// 105528 was in the error: gas too low: have 0, want 105528
	// Not sure if there's a better way to get it.
	it.auth.GasLimit = 1055280000

	address, tx, ts, err := testfiles.DeployTestfiles(it.auth, it.sim)

	require.NoError(t, err)
	it.sim.Commit()
	it.evmTest = ts
	it.incNonce()
	it.awaitTx(ctx, t, tx)
	it.address = address.String()
}

func (it *interfaceTester) awaitTx(ctx context.Context, t *testing.T, tx *evmtypes.Transaction) {
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

func argsFromTestStruct(ts TestStruct) []any {
	return []any{
		ts.Field,
		ts.DifferentField,
		uint8(ts.OracleId),
		getOracleIds(ts),
		[32]byte(ts.Account),
		getAccounts(ts),
		ts.BigField,
		midToInternalType(ts.NestedStruct),
	}
}

func getOracleIds(first TestStruct) [32]byte {
	oracleIds := [32]byte{}
	for i, oracleId := range first.OracleIds {
		oracleIds[i] = byte(oracleId)
	}
	return oracleIds
}

func toInternalType(testStruct TestStruct) testfiles.TestStruct {
	return testfiles.TestStruct{
		Field:          testStruct.Field,
		DifferentField: testStruct.DifferentField,
		OracleId:       byte(testStruct.OracleId),
		OracleIds:      convertOracleIds(testStruct.OracleIds),
		Account:        [32]byte(testStruct.Account),
		Accounts:       convertAccounts(testStruct.Accounts),
		BigField:       testStruct.BigField,
		NestedStruct:   midToInternalType(testStruct.NestedStruct),
	}
}

func midToInternalType(m MidLevelTestStruct) testfiles.MidLevelTestStruct {
	return testfiles.MidLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: testfiles.InnerTestStruct{
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
