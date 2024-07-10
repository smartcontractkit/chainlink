package evmtesting

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jmoiron/sqlx"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_reader_tester"
	_ "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest" // force binding for tx type
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const (
	triggerWithDynamicTopic = "TriggeredEventWithDynamicTopic"
	triggerWithAllTopics    = "TriggeredWithFourTopics"
)

type EVMChainReaderInterfaceTesterHelper[T TestingT[T]] interface {
	SetupAuth(t T) *bind.TransactOpts
	Client(t T) client.Client
	Commit()
	MustGenerateRandomKey(t T) ethkey.KeyV2
	Backend() bind.ContractBackend
	ChainID() *big.Int
	Context(t T) context.Context
	NewSqlxDB(t T) *sqlx.DB
	MaxWaitTimeForEvents() time.Duration
	GasPriceBufferPercent() int64
}

type EVMChainReaderInterfaceTester[T TestingT[T]] struct {
	Helper         EVMChainReaderInterfaceTesterHelper[T]
	client         client.Client
	address        string
	address2       string
	chainConfig    types.ChainReaderConfig
	auth           *bind.TransactOpts
	evmTest        *chain_reader_tester.ChainReaderTester
	cr             evm.ChainReaderService
	dirtyContracts bool
}

func (it *EVMChainReaderInterfaceTester[T]) Setup(t T) {
	t.Cleanup(func() {
		// DB may be closed by the test already, ignore errors
		if it.cr != nil {
			_ = it.cr.Close()
		}
		it.cr = nil

		if it.dirtyContracts {
			it.evmTest = nil
		}
	})

	// can re-use the same chain for tests, just make new contract for each test
	if it.client != nil {
		it.deployNewContracts(t)
		return
	}

	it.auth = it.Helper.SetupAuth(t)

	testStruct := CreateTestStruct[T](0, it)

	it.chainConfig = types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			AnyContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				ContractPollingFilter: types.ContractPollingFilter{
					GenericEventNames: []string{EventName, EventWithFilterName},
				},
				Configs: map[string]*types.ChainReaderDefinition{
					MethodTakingLatestParamsReturningTestStruct: {
						ChainSpecificName: "getElementAtIndex",
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStruct.Inner.IntVal": "I"}},
						},
					},
					// this is supposed to be used for testing confidence levels, but geth simulated backend doesn't support calling past state
					//MethodReturningAlterableUint64: {
					//	ChainSpecificName:       "getAlterablePrimitiveValue",
					//},
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
					},
					EventWithFilterName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
						EventDefinitions:  &types.EventDefinitions{InputFields: []string{"Field"}},
					},
					triggerWithDynamicTopic: {
						ChainSpecificName: triggerWithDynamicTopic,
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							InputFields: []string{"fieldHash"},
							// no specific reason for filter being defined here insted on contract level,
							// this is just for test case variety
							PollingFilter: &types.PollingFilter{},
						},
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"FieldHash": "Field"}},
						},
						ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": -1},
					},
					triggerWithAllTopics: {
						ChainSpecificName: triggerWithAllTopics,
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							InputFields:   []string{"Field1", "Field2", "Field3"},
							PollingFilter: &types.PollingFilter{},
						},
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
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"ExtraField": AnyExtraValue}},
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
	it.client = it.Helper.Client(t)

	it.deployNewContracts(t)
}

func (it *EVMChainReaderInterfaceTester[T]) Name() string {
	return "EVM"
}

func (it *EVMChainReaderInterfaceTester[T]) GetAccountBytes(i int) []byte {
	account := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	account[i%20] += byte(i)
	account[(i+3)%20] += byte(i + 3)
	return account[:]
}

func (it *EVMChainReaderInterfaceTester[T]) GetChainReader(t T) clcommontypes.ContractReader {
	ctx := it.Helper.Context(t)
	if it.cr != nil {
		return it.cr
	}

	lggr := logger.NullLogger
	db := it.Helper.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            4,
		BackfillBatchSize:        1,
		RpcBatchSize:             1,
		KeepFinalizedBlocksDepth: 10000,
	}
	ht := headtracker.NewSimulatedHeadTracker(it.client, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(logpoller.NewORM(it.Helper.ChainID(), db, lggr), it.client, lggr, ht, lpOpts)
	require.NoError(t, lp.Start(ctx))

	// encode and decode the config to ensure the test covers type issues
	confBytes, err := json.Marshal(it.chainConfig)
	require.NoError(t, err)

	conf, err := types.ChainReaderConfigFromBytes(confBytes)
	require.NoError(t, err)

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, ht, it.client, conf)
	require.NoError(t, err)
	require.NoError(t, cr.Start(ctx))
	it.cr = cr
	return cr
}

func (it *EVMChainReaderInterfaceTester[T]) SetTestStructLatestValue(t T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*chain_reader_tester.ChainReaderTesterTransactor).AddTestStruct)
}

// SetUintLatestValue is supposed to be used for testing confidence levels, but geth simulated backend doesn't support calling past state
func (it *EVMChainReaderInterfaceTester[T]) SetUintLatestValue(t T, val uint64) {
	it.sendTxWithUintVal(t, val, (*chain_reader_tester.ChainReaderTesterTransactor).SetAlterablePrimitiveValue)
}

func (it *EVMChainReaderInterfaceTester[T]) TriggerEvent(t T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*chain_reader_tester.ChainReaderTesterTransactor).TriggerEvent)
}

// GenerateBlocksTillConfidenceLevel is supposed to be used for testing confidence levels, but geth simulated backend doesn't support calling past state
func (it *EVMChainReaderInterfaceTester[T]) GenerateBlocksTillConfidenceLevel(t T, contractName, readName string, confidenceLevel primitives.ConfidenceLevel) {
	contractCfg, ok := it.chainConfig.Contracts[contractName]
	if !ok {
		t.Errorf("contract %s not found", contractName)
		return
	}

	readCfg, ok := contractCfg.Configs[readName]
	require.True(t, ok, fmt.Sprintf("readName: %s not found for contract: %s", readName, contractName))

	toEvmConf, err := evm.ConfirmationsFromConfig(readCfg.ConfidenceConfirmations)
	require.True(t, ok, fmt.Errorf("failed to parse confidence level mapping:%s not found for contract: %s readName: %s, err:%w", confidenceLevel, readName, contractName, err))

	confirmations, ok := toEvmConf[confidenceLevel]
	require.True(t, ok, fmt.Sprintf("confidence level mapping:%s not found for contract: %s readName: %s", confidenceLevel, readName, contractName))

	key := it.Helper.MustGenerateRandomKey(t)
	pk := key.ToEcdsaPrivKey()
	toAddress := common.HexToAddress("0x0")

	// confirmations are in form of negative values that signify how many blocks are needed for a specific confidence level
	for i := confirmations; i == 0; i++ {
		nonce, err := it.client.PendingNonceAt(it.Helper.Context(t), key.Address)
		require.NoError(t, err)

		tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{ChainID: it.client.ConfiguredChainID(), Nonce: nonce, Gas: 21000, To: &toAddress})
		signedTx, err := gethtypes.SignTx(tx, gethtypes.NewEIP155Signer(it.client.ConfiguredChainID()), pk)
		require.NoError(t, err)

		require.NoError(t, it.client.SendTransaction(it.Helper.Context(t), signedTx))
		it.AwaitTx(t, &gethtypes.Transaction{})
	}
}

func (it *EVMChainReaderInterfaceTester[T]) GetBindings(_ T) []clcommontypes.BoundContract {
	return []clcommontypes.BoundContract{
		{Name: AnyContractName, Address: it.address},
		{Name: AnySecondContractName, Address: it.address2},
	}
}

type uintFn = func(*chain_reader_tester.ChainReaderTesterTransactor, *bind.TransactOpts, uint64) (*gethtypes.Transaction, error)

// sendTxWithUintVal is supposed to be used for testing confidence levels, but geth simulated backend doesn't support calling past state
func (it *EVMChainReaderInterfaceTester[T]) sendTxWithUintVal(t T, val uint64, fn uintFn) {
	tx, err := fn(
		&it.evmTest.ChainReaderTesterTransactor,
		it.GetAuthWithGasSet(t),
		val,
	)

	require.NoError(t, err)
	it.Helper.Commit()
	it.IncNonce()
	it.AwaitTx(t, tx)
	it.dirtyContracts = true
}

type testStructFn = func(*chain_reader_tester.ChainReaderTesterTransactor, *bind.TransactOpts, int32, string, uint8, [32]uint8, common.Address, []common.Address, *big.Int, chain_reader_tester.MidLevelTestStruct) (*gethtypes.Transaction, error)

func (it *EVMChainReaderInterfaceTester[T]) sendTxWithTestStruct(t T, testStruct *TestStruct, fn testStructFn) {
	tx, err := fn(
		&it.evmTest.ChainReaderTesterTransactor,
		it.GetAuthWithGasSet(t),
		*testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleID),
		OracleIdsToBytes(testStruct.OracleIDs),
		common.Address(testStruct.Account),
		ConvertAccounts(testStruct.Accounts),
		testStruct.BigField,
		MidToInternalType(testStruct.NestedStruct),
	)
	require.NoError(t, err)
	it.Helper.Commit()
	it.IncNonce()
	it.AwaitTx(t, tx)
	it.dirtyContracts = true
}

func (it *EVMChainReaderInterfaceTester[T]) GetAuthWithGasSet(t T) *bind.TransactOpts {
	gasPrice, err := it.client.SuggestGasPrice(it.Helper.Context(t))
	require.NoError(t, err)
	extra := new(big.Int).Mul(gasPrice, big.NewInt(it.Helper.GasPriceBufferPercent()))
	extra = extra.Div(extra, big.NewInt(100))
	it.auth.GasPrice = gasPrice.Add(gasPrice, extra)
	return it.auth
}

func (it *EVMChainReaderInterfaceTester[T]) IncNonce() {
	if it.auth.Nonce == nil {
		it.auth.Nonce = big.NewInt(1)
	} else {
		it.auth.Nonce = it.auth.Nonce.Add(it.auth.Nonce, big.NewInt(1))
	}
}

func (it *EVMChainReaderInterfaceTester[T]) AwaitTx(t T, tx *gethtypes.Transaction) {
	ctx := it.Helper.Context(t)
	receipt, err := bind.WaitMined(ctx, it.client, tx)
	require.NoError(t, err)
	require.Equal(t, gethtypes.ReceiptStatusSuccessful, receipt.Status)
}

func (it *EVMChainReaderInterfaceTester[T]) deployNewContracts(t T) {
	// First test deploy both contracts, otherwise only deploy contracts if cleanup decides that we need to.
	if it.address == "" {
		it.address = it.deployNewContract(t)
		it.address2 = it.deployNewContract(t)
	} else if it.evmTest == nil {
		it.address = it.deployNewContract(t)
		it.dirtyContracts = false
	}
}

func (it *EVMChainReaderInterfaceTester[T]) deployNewContract(t T) string {
	// 105528 was in the error: gas too low: have 0, want 105528
	// Not sure if there's a better way to get it.
	it.auth.GasLimit = 10552800

	address, tx, ts, err := chain_reader_tester.DeployChainReaderTester(it.GetAuthWithGasSet(t), it.Helper.Backend())
	require.NoError(t, err)
	it.Helper.Commit()
	if it.evmTest == nil {
		it.evmTest = ts
	}

	it.IncNonce()
	it.AwaitTx(t, tx)
	return address.String()
}

func (it *EVMChainReaderInterfaceTester[T]) MaxWaitTimeForEvents() time.Duration {
	return it.Helper.MaxWaitTimeForEvents()
}

func OracleIdsToBytes(oracleIDs [32]commontypes.OracleID) [32]byte {
	convertedIds := [32]byte{}
	for i, id := range oracleIDs {
		convertedIds[i] = byte(id)
	}
	return convertedIds
}

func ConvertAccounts(accounts [][]byte) []common.Address {
	convertedAccounts := make([]common.Address, len(accounts))
	for i, a := range accounts {
		convertedAccounts[i] = common.Address(a)
	}
	return convertedAccounts
}

func ToInternalType(testStruct TestStruct) chain_reader_tester.TestStruct {
	return chain_reader_tester.TestStruct{
		Field:          *testStruct.Field,
		DifferentField: testStruct.DifferentField,
		OracleId:       byte(testStruct.OracleID),
		OracleIds:      OracleIdsToBytes(testStruct.OracleIDs),
		Account:        common.Address(testStruct.Account),
		Accounts:       ConvertAccounts(testStruct.Accounts),
		BigField:       testStruct.BigField,
		NestedStruct:   MidToInternalType(testStruct.NestedStruct),
	}
}

func MidToInternalType(m MidLevelTestStruct) chain_reader_tester.MidLevelTestStruct {
	return chain_reader_tester.MidLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: chain_reader_tester.InnerTestStruct{
			IntVal: int64(m.Inner.I),
			S:      m.Inner.S,
		},
	}
}
