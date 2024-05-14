package evmtesting

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jmoiron/sqlx"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_reader_example"
	_ "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest" // force binding for tx type
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .
)

const (
	triggerWithDynamicTopic = "TriggeredEventWithDynamicTopic"
	triggerWithAllTopics    = "TriggeredWithFourTopics"
)

type EvmChainReaderInterfaceTesterHelper[T TestingT[T]] interface {
	SetupAuth(t T) *bind.TransactOpts
	Client(t T) client.Client
	Commit()
	Backend() bind.ContractBackend
	ChainID() *big.Int
	Context(t T) context.Context
	NewSqlxDB(t T) *sqlx.DB
	MaxWaitTimeForEvents() time.Duration
}

type EvmChainReaderInterfaceTester[T TestingT[T]] struct {
	Helper      EvmChainReaderInterfaceTesterHelper[T]
	client      client.Client
	address     string
	address2    string
	chainConfig types.ChainReaderConfig
	auth        *bind.TransactOpts
	evmTest     *chain_reader_example.LatestValueHolder
	cr          evm.ChainReaderService
}

func (it *EvmChainReaderInterfaceTester[T]) Setup(t T) {
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

	it.auth = it.Helper.SetupAuth(t)

	testStruct := CreateTestStruct[T](0, it)

	it.chainConfig = types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			AnyContractName: {
				ContractABI: chain_reader_example.LatestValueHolderMetaData.ABI,
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
					},
					EventWithFilterName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
						EventInputFields:  []string{"Field"},
					},
					triggerWithDynamicTopic: {
						ChainSpecificName: triggerWithDynamicTopic,
						ReadType:          types.Event,
						EventInputFields:  []string{"fieldHash"},
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"FieldHash": "Field"}},
						},
					},
					triggerWithAllTopics: {
						ChainSpecificName: triggerWithAllTopics,
						ReadType:          types.Event,
						EventInputFields:  []string{"Field1", "Field2", "Field3"},
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
				ContractABI: chain_reader_example.LatestValueHolderMetaData.ABI,
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

func (it *EvmChainReaderInterfaceTester[T]) Name() string {
	return "EVM"
}

func (it *EvmChainReaderInterfaceTester[T]) GetAccountBytes(i int) []byte {
	account := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	account[i%20] += byte(i)
	account[(i+3)%20] += byte(i + 3)
	return account[:]
}

func (it *EvmChainReaderInterfaceTester[T]) GetChainReader(t T) clcommontypes.ChainReader {
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
	lp := logpoller.NewLogPoller(logpoller.NewORM(it.Helper.ChainID(), db, lggr), it.client, lggr, lpOpts)
	require.NoError(t, lp.Start(ctx))

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, it.client, it.chainConfig)
	require.NoError(t, err)
	require.NoError(t, cr.Start(ctx))
	it.cr = cr
	return cr
}

func (it *EvmChainReaderInterfaceTester[T]) SetLatestValue(t T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*chain_reader_example.LatestValueHolderTransactor).AddTestStruct)
}

func (it *EvmChainReaderInterfaceTester[T]) TriggerEvent(t T, testStruct *TestStruct) {
	it.sendTxWithTestStruct(t, testStruct, (*chain_reader_example.LatestValueHolderTransactor).TriggerEvent)
}

func (it *EvmChainReaderInterfaceTester[T]) GetBindings(_ T) []clcommontypes.BoundContract {
	return []clcommontypes.BoundContract{
		{Name: AnyContractName, Address: it.address, Pending: true},
		{Name: AnySecondContractName, Address: it.address2, Pending: true},
	}
}

type testStructFn = func(*chain_reader_example.LatestValueHolderTransactor, *bind.TransactOpts, int32, string, uint8, [32]uint8, common.Address, []common.Address, *big.Int, chain_reader_example.MidLevelTestStruct) (*evmtypes.Transaction, error)

func (it *EvmChainReaderInterfaceTester[T]) sendTxWithTestStruct(t T, testStruct *TestStruct, fn testStructFn) {
	tx, err := fn(
		&it.evmTest.LatestValueHolderTransactor,
		it.auth,
		*testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleID),
		ConvertOracleIDs(testStruct.OracleIDs),
		common.Address(testStruct.Account),
		ConvertAccounts(testStruct.Accounts),
		testStruct.BigField,
		MidToInternalType(testStruct.NestedStruct),
	)
	require.NoError(t, err)
	it.Helper.Commit()
	it.IncNonce()
	it.AwaitTx(t, tx)
}

func (it *EvmChainReaderInterfaceTester[T]) IncNonce() {
	if it.auth.Nonce == nil {
		it.auth.Nonce = big.NewInt(1)
	} else {
		it.auth.Nonce = it.auth.Nonce.Add(it.auth.Nonce, big.NewInt(1))
	}
}

func (it *EvmChainReaderInterfaceTester[T]) AwaitTx(t T, tx *evmtypes.Transaction) {
	ctx := it.Helper.Context(t)
	receipt, err := bind.WaitMined(ctx, it.client, tx)
	require.NoError(t, err)
	require.Equal(t, evmtypes.ReceiptStatusSuccessful, receipt.Status)
}

func (it *EvmChainReaderInterfaceTester[T]) deployNewContracts(t T) {
	it.address = it.deployNewContract(t)
	it.address2 = it.deployNewContract(t)
}

func (it *EvmChainReaderInterfaceTester[T]) deployNewContract(t T) string {
	ctx := it.Helper.Context(t)
	gasPrice, err := it.client.SuggestGasPrice(ctx)
	require.NoError(t, err)
	it.auth.GasPrice = gasPrice

	// 105528 was in the error: gas too low: have 0, want 105528
	// Not sure if there's a better way to get it.
	it.auth.GasLimit = 10552800

	address, tx, ts, err := chain_reader_example.DeployLatestValueHolder(it.auth, it.Helper.Backend())
	require.NoError(t, err)
	it.Helper.Commit()
	if it.evmTest == nil {
		it.evmTest = ts
	}

	it.IncNonce()
	it.AwaitTx(t, tx)
	return address.String()
}

func (it *EvmChainReaderInterfaceTester[T]) MaxWaitTimeForEvents() time.Duration {
	return it.Helper.MaxWaitTimeForEvents()
}

func ConvertOracleIDs(oracleIDs [32]commontypes.OracleID) [32]byte {
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

func ToInternalType(testStruct TestStruct) chain_reader_example.TestStruct {
	return chain_reader_example.TestStruct{
		Field:          *testStruct.Field,
		DifferentField: testStruct.DifferentField,
		OracleId:       byte(testStruct.OracleID),
		OracleIds:      ConvertOracleIDs(testStruct.OracleIDs),
		Account:        common.Address(testStruct.Account),
		Accounts:       ConvertAccounts(testStruct.Accounts),
		BigField:       testStruct.BigField,
		NestedStruct:   MidToInternalType(testStruct.NestedStruct),
	}
}

func MidToInternalType(m MidLevelTestStruct) chain_reader_example.MidLevelTestStruct {
	return chain_reader_example.MidLevelTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: chain_reader_example.InnerTestStruct{
			IntVal: int64(m.Inner.I),
			S:      m.Inner.S,
		},
	}
}
