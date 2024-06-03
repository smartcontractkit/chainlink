package evm_test

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .

	commontestutils "github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_reader_tester"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const (
	commonGasLimitOnEvms    = uint64(4712388)
	triggerWithDynamicTopic = "TriggeredEventWithDynamicTopic"
	triggerWithAllTopics    = "TriggeredWithFourTopics"
)

func TestChainReaderInterfaceTests(t *testing.T) {
	t.Parallel()
	it := &chainReaderInterfaceTester{}

	RunChainReaderInterfaceTests(t, it)
	RunChainReaderInterfaceTests(t, commontestutils.WrapChainReaderTesterForLoop(it))

	t.Run("Dynamically typed topics can be used to filter and have type correct in return", func(t *testing.T) {
		it.Setup(t)

		// bind event before firing it to avoid log poller race
		ctx := testutils.Context(t)
		cr := it.GetChainReader(t)
		require.NoError(t, cr.Bind(ctx, it.GetBindings(t)))

		anyString := "foo"
		tx, err := it.evmTest.TriggerEventWithDynamicTopic(it.auth, anyString)
		require.NoError(t, err)
		it.sim.Commit()
		it.incNonce()
		it.awaitTx(t, tx)

		input := struct{ Field string }{Field: anyString}
		tp := cr.(clcommontypes.ContractTypeProvider)
		output, err := tp.CreateContractType(AnyContractName, triggerWithDynamicTopic, false)
		require.NoError(t, err)
		rOutput := reflect.Indirect(reflect.ValueOf(output))

		require.Eventually(t, func() bool {
			return cr.GetLatestValue(ctx, AnyContractName, triggerWithDynamicTopic, input, output) == nil
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)

		assert.Equal(t, &anyString, rOutput.FieldByName("Field").Interface())
		topic, err := abi.MakeTopics([]any{anyString})
		require.NoError(t, err)
		assert.Equal(t, &topic[0][0], rOutput.FieldByName("FieldHash").Interface())
	})

	t.Run("Multiple topics can filter together", func(t *testing.T) {
		it.Setup(t)

		// bind event before firing it to avoid log poller race
		ctx := testutils.Context(t)
		cr := it.GetChainReader(t)
		require.NoError(t, cr.Bind(ctx, it.GetBindings(t)))

		triggerFourTopics(t, it, int32(1), int32(2), int32(3))
		triggerFourTopics(t, it, int32(2), int32(2), int32(3))
		triggerFourTopics(t, it, int32(1), int32(3), int32(3))
		triggerFourTopics(t, it, int32(1), int32(2), int32(4))

		var latest struct{ Field1, Field2, Field3 int32 }
		params := struct{ Field1, Field2, Field3 int32 }{Field1: 1, Field2: 2, Field3: 3}

		require.Eventually(t, func() bool {
			return cr.GetLatestValue(ctx, AnyContractName, triggerWithAllTopics, params, &latest) == nil
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)

		assert.Equal(t, int32(1), latest.Field1)
		assert.Equal(t, int32(2), latest.Field2)
		assert.Equal(t, int32(3), latest.Field3)
	})
}

func triggerFourTopics(t *testing.T, it *chainReaderInterfaceTester, i1, i2, i3 int32) {
	tx, err := it.evmTest.ChainReaderTesterTransactor.TriggerWithFourTopics(it.auth, i1, i2, i3)
	require.NoError(t, err)
	require.NoError(t, err)
	it.sim.Commit()
	it.incNonce()
	it.awaitTx(t, tx)
}

type chainReaderInterfaceTester struct {
	client      client.Client
	address     string
	address2    string
	chainConfig types.ChainReaderConfig
	auth        *bind.TransactOpts
	sim         *backends.SimulatedBackend
	pk          *ecdsa.PrivateKey
	evmTest     *chain_reader_tester.ChainReaderTester
	cr          evm.ChainReaderService
}

func (it *chainReaderInterfaceTester) MaxWaitTimeForEvents() time.Duration {
	// From trial and error, when running on CI, sometimes the boxes get slow
	maxWaitTime := time.Second * 20
	maxWaitTimeStr, ok := os.LookupEnv("MAX_WAIT_TIME_FOR_EVENTS_S")
	if ok {
		waitS, err := strconv.ParseInt(maxWaitTimeStr, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing MAX_WAIT_TIME_FOR_EVENTS_S: %v, defaulting to %v\n", err, maxWaitTime)
		}
		maxWaitTime = time.Second * time.Duration(waitS)
	}

	return maxWaitTime
}

func (it *chainReaderInterfaceTester) Setup(t *testing.T) {
	t.Cleanup(func() {
		// DB may be closed by the test already, ignore errors
		if it.cr != nil {
			_ = it.cr.Close()
		}
		it.cr = nil
		it.evmTest = nil
	})

	// can re-use the same chain for tests, just make new contract for each test
	if it.client != nil {
		it.deployNewContracts(t)
		return
	}

	it.setupChainNoClient(t)

	testStruct := CreateTestStruct(0, it)

	it.chainConfig = types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			AnyContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				Configs: map[string]*types.ChainReaderDefinition{
					MethodTakingLatestParamsReturningTestStruct: {
						ChainSpecificName: "getElementAtIndex",
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStruct.Inner.IntVal": "I"}},
						},
					},
					MethodReturningUint64: {
						ChainSpecificName: "getPrimitiveValue",
					},
					DifferentMethodReturningUint64: {
						ChainSpecificName: "getDifferentPrimitiveValue",
					},
					MethodReturningUint64Slice: {
						ChainSpecificName: "getSliceValue",
					},
					EventName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStruct.Inner.IntVal": "I"}},
						},
						ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": -1},
					},
					EventWithFilterName: {
						ChainSpecificName:       "Triggered",
						ReadType:                types.Event,
						EventInputFields:        []string{"Field"},
						ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": -1},
					},
					triggerWithDynamicTopic: {
						ChainSpecificName: triggerWithDynamicTopic,
						ReadType:          types.Event,
						EventInputFields:  []string{"fieldHash"},
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"FieldHash": "Field"}},
						},
						ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": -1},
					},
					triggerWithAllTopics: {
						ChainSpecificName:       triggerWithAllTopics,
						ReadType:                types.Event,
						EventInputFields:        []string{"Field1", "Field2", "Field3"},
						ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": -1},
					},
					MethodReturningSeenStruct: {
						ChainSpecificName: "returnSeen",
						InputModifications: codec.ModifiersConfig{
							&codec.HardCodeModifierConfig{
								OnChainValues: map[string]any{
									"BigField": testStruct.BigField.String(),
									"Account":  hexutil.Encode(testStruct.Account),
								},
							},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStruct.Inner.IntVal": "I"}},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"ExtraField": anyExtraValue}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStruct.Inner.IntVal": "I"}},
						},
					},
				},
			},
			AnySecondContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				Configs: map[string]*types.ChainReaderDefinition{
					MethodReturningUint64: {
						ChainSpecificName: "getDifferentPrimitiveValue",
					},
				},
			},
		},
	}
	it.client = client.NewSimulatedBackendClient(t, it.sim, big.NewInt(1337))
	it.deployNewContracts(t)
}

func (it *chainReaderInterfaceTester) Name() string {
	return "EVM"
}

func (it *chainReaderInterfaceTester) GetAccountBytes(i int) []byte {
	account := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	account[i%20] += byte(i)
	account[(i+3)%20] += byte(i + 3)
	return account[:]
}

func (it *chainReaderInterfaceTester) GetChainReader(t *testing.T) clcommontypes.ContractReader {
	ctx := testutils.Context(t)
	if it.cr != nil {
		return it.cr
	}

	lggr := logger.NullLogger
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            4,
		BackfillBatchSize:        1,
		RpcBatchSize:             1,
		KeepFinalizedBlocksDepth: 10000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, db, lggr), it.client, lggr, lpOpts)
	require.NoError(t, lp.Start(ctx))

	// TODO  uncomment this after this is fixed BCF-3242
	//chain := mocks.NewChain(t)
	//chain.Mock.On("LogPoller").Return(lp)
	//chain.Mock.On("ID").Return(it.client.ConfiguredChainID())
	//
	//keyStore := cltest.NewKeyStore(t, db)
	//relayer, err := evm.NewRelayer(lggr, chain, evm.RelayerOpts{DS: db, CSAETHKeystore: keyStore, CapabilitiesRegistry: capabilities.NewRegistry(lggr)})
	//require.NoError(t, err)
	//
	//cfgBytes, err := cbor.Marshal(it.chainConfig)
	//require.NoError(t, err)
	//cr, err := relayer.NewContractReader(cfgBytes)

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, it.client, it.chainConfig)
	require.NoError(t, err)
	require.NoError(t, cr.Start(ctx))
	it.cr = cr
	return cr
}

func (it *chainReaderInterfaceTester) SetLatestValue(t *testing.T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*chain_reader_tester.ChainReaderTesterTransactor).AddTestStruct)
}

func (it *chainReaderInterfaceTester) TriggerEvent(t *testing.T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*chain_reader_tester.ChainReaderTesterTransactor).TriggerEvent)
}

func (it *chainReaderInterfaceTester) GetBindings(_ *testing.T) []clcommontypes.BoundContract {
	return []clcommontypes.BoundContract{
		{Name: AnyContractName, Address: it.address, Pending: true},
		{Name: AnySecondContractName, Address: it.address2, Pending: true},
	}
}

type testStructFn = func(*chain_reader_tester.ChainReaderTesterTransactor, *bind.TransactOpts, int32, string, uint8, [32]uint8, common.Address, []common.Address, *big.Int, chain_reader_tester.MidLevelTestStruct) (*evmtypes.Transaction, error)

func (it *chainReaderInterfaceTester) sendTxWithTestStruct(t *testing.T, testStruct *TestStruct, fn testStructFn) {
	tx, err := fn(
		&it.evmTest.ChainReaderTesterTransactor,
		it.auth,
		*testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleID),
		convertOracleIDs(testStruct.OracleIDs),
		common.Address(testStruct.Account),
		convertAccounts(testStruct.Accounts),
		testStruct.BigField,
		midToInternalType(testStruct.NestedStruct),
	)
	require.NoError(t, err)
	it.sim.Commit()
	it.incNonce()
	it.awaitTx(t, tx)
}

func convertOracleIDs(oracleIDs [32]commontypes.OracleID) [32]byte {
	convertedIds := [32]byte{}
	for i, id := range oracleIDs {
		convertedIds[i] = byte(id)
	}
	return convertedIds
}

func convertAccounts(accounts [][]byte) []common.Address {
	convertedAccounts := make([]common.Address, len(accounts))
	for i, a := range accounts {
		convertedAccounts[i] = common.Address(a)
	}
	return convertedAccounts
}

func (it *chainReaderInterfaceTester) setupChainNoClient(t require.TestingT) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	it.pk = privateKey

	it.auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	it.sim = backends.NewSimulatedBackend(core.GenesisAlloc{it.auth.From: {Balance: big.NewInt(math.MaxInt64)}}, commonGasLimitOnEvms*5000)
	it.sim.Commit()
}

func (it *chainReaderInterfaceTester) deployNewContracts(t *testing.T) {
	it.address = it.deployNewContract(t)
	it.address2 = it.deployNewContract(t)
}

func (it *chainReaderInterfaceTester) deployNewContract(t *testing.T) string {
	ctx := testutils.Context(t)
	gasPrice, err := it.sim.SuggestGasPrice(ctx)
	require.NoError(t, err)
	it.auth.GasPrice = gasPrice

	// 105528 was in the error: gas too low: have 0, want 105528
	// Not sure if there's a better way to get it.
	it.auth.GasLimit = 10552800

	address, tx, ts, err := chain_reader_tester.DeployChainReaderTester(it.auth, it.sim)

	require.NoError(t, err)
	it.sim.Commit()
	if it.evmTest == nil {
		it.evmTest = ts
	}
	it.incNonce()
	it.awaitTx(t, tx)
	return address.String()
}

func (it *chainReaderInterfaceTester) awaitTx(t *testing.T, tx *evmtypes.Transaction) {
	ctx := testutils.Context(t)
	receipt, err := it.sim.TransactionReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	require.Equal(t, evmtypes.ReceiptStatusSuccessful, receipt.Status)
}

func (it *chainReaderInterfaceTester) incNonce() {
	if it.auth.Nonce == nil {
		it.auth.Nonce = big.NewInt(1)
	} else {
		it.auth.Nonce = it.auth.Nonce.Add(it.auth.Nonce, big.NewInt(1))
	}
}

func getAccounts(first TestStruct) []common.Address {
	accountBytes := make([]common.Address, len(first.Accounts))
	for i, account := range first.Accounts {
		accountBytes[i] = common.Address(account)
	}
	return accountBytes
}

func argsFromTestStruct(ts TestStruct) []any {
	return []any{
		ts.Field,
		ts.DifferentField,
		uint8(ts.OracleID),
		getOracleIDs(ts),
		common.Address(ts.Account),
		getAccounts(ts),
		ts.BigField,
		midToInternalType(ts.NestedStruct),
	}
}

func getOracleIDs(first TestStruct) [32]byte {
	oracleIDs := [32]byte{}
	for i, oracleID := range first.OracleIDs {
		oracleIDs[i] = byte(oracleID)
	}
	return oracleIDs
}

func toInternalType(testStruct TestStruct) chain_reader_tester.TestStruct {
	return chain_reader_tester.TestStruct{
		Field:          *testStruct.Field,
		DifferentField: testStruct.DifferentField,
		OracleId:       byte(testStruct.OracleID),
		OracleIds:      convertOracleIDs(testStruct.OracleIDs),
		Account:        common.Address(testStruct.Account),
		Accounts:       convertAccounts(testStruct.Accounts),
		BigField:       testStruct.BigField,
		NestedStruct:   midToInternalType(testStruct.NestedStruct),
	}
}

func midToInternalType(m MidLevelTestStruct) chain_reader_tester.MidLevelTestStruct {
	return chain_reader_tester.MidLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: chain_reader_tester.InnerTestStruct{
			IntVal: int64(m.Inner.I),
			S:      m.Inner.S,
		},
	}
}
