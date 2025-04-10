package evmtesting

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/chain_reader_tester"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	_ "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest" // force binding for tx type
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint:revive // dot-imports
)

const (
	triggerWithDynamicTopic        = "TriggeredEventWithDynamicTopic"
	triggerWithAllTopics           = "TriggeredWithFourTopics"
	triggerWithAllTopicsWithHashed = "TriggeredWithFourTopicsWithHashed"
	staticBytesEventName           = "StaticBytes"
)

type EVMChainComponentsInterfaceTesterHelper[T TestingT[T]] interface {
	Client(t T) client.Client
	Commit()
	Backend() bind.ContractBackend
	Database() *sqlx.DB
	ChainID() *big.Int
	Context(t T) context.Context
	NewSqlxDB(t T) *sqlx.DB
	MaxWaitTimeForEvents() time.Duration
	GasPriceBufferPercent() int64
	Accounts(t T) []*bind.TransactOpts
	TXM(T, client.Client) evmtxmgr.TxManager
	// To enable the historical wrappers required for Simulated Backend tests.
	ChainReaderEVMClient(ctx context.Context, t T, ht logpoller.HeadTracker, conf types.ChainReaderConfig) client.Client
	WrappedChainWriter(cw clcommontypes.ContractWriter, client client.Client) clcommontypes.ContractWriter
	LogPoller(t T) logpoller.LogPoller
	HeadTracker(t T) logpoller.HeadTracker
}

type EVMChainComponentsInterfaceTester[T TestingT[T]] struct {
	TestSelectionSupport
	Helper                    EVMChainComponentsInterfaceTesterHelper[T]
	DeployLock                *sync.Mutex
	client                    client.Client
	chainReaderConfigSupplier func(t T) types.ChainReaderConfig
	chainWriterConfigSupplier func(t T) types.ChainWriterConfig
}

func (it *EVMChainComponentsInterfaceTester[T]) GetBindings(t T) []clcommontypes.BoundContract {
	return it.deployNewContracts(t)
}

func (it *EVMChainComponentsInterfaceTester[T]) getChainReaderConfig(t T) types.ChainReaderConfig {
	testStruct := CreateTestStruct[T](0, it)

	methodTakingLatestParamsReturningTestStructConfig := types.ChainReaderDefinition{
		ChainSpecificName: "getElementAtIndex",
		OutputModifications: codec.ModifiersConfig{
			&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
			&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
			&codec.AddressBytesToStringModifierConfig{
				Fields: []string{"AccountStruct.AccountStr"},
			},
		},
	}

	return types.ChainReaderConfig{
		Contracts: map[string]types.ChainContractReader{
			AnyContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				ContractPollingFilter: types.ContractPollingFilter{
					GenericEventNames: []string{EventName, EventWithFilterName, triggerWithAllTopicsWithHashed, staticBytesEventName},
				},
				Configs: map[string]*types.ChainReaderDefinition{
					MethodTakingLatestParamsReturningTestStruct: &methodTakingLatestParamsReturningTestStructConfig,
					MethodReturningAlterableUint64: {
						ChainSpecificName: "getAlterablePrimitiveValue",
					},
					MethodReturningUint64: {
						ChainSpecificName: "getPrimitiveValue",
					},
					MethodReturningUint64Slice: {
						ChainSpecificName: "getSliceValue",
					},
					EventName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							GenericTopicNames: map[string]string{"field": "Field"},
							GenericDataWordDetails: map[string]types.DataWordDetail{
								"OracleID": {Name: "oracleId"},
								// this is just to illustrate an example, generic names shouldn't really be formatted like this since other chains might not store it in the same way
								"NestedStaticStruct.Inner.IntVal": {Name: "nestedStaticStruct.Inner.IntVal"},
								"NestedDynamicStruct.FixedBytes":  {Name: "nestedDynamicStruct.FixedBytes"},
								"BigField":                        {Name: "bigField"},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
							&codec.AddressBytesToStringModifierConfig{
								Fields: []string{"AccountStruct.AccountStr"},
							},
						},
					},
					staticBytesEventName: {
						ChainSpecificName: staticBytesEventName,
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							GenericDataWordDetails: map[string]types.DataWordDetail{
								"msgTransmitterEvent": {
									Name:  "msgTransmitterEvent",
									Index: testutils.Ptr(2),
									Type:  "bytes32",
								},
							},
						},
					},
					EventWithFilterName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
						OutputModifications: codec.ModifiersConfig{
							&codec.AddressBytesToStringModifierConfig{
								Fields: []string{"AccountStruct.AccountStr"},
							},
						},
					},
					triggerWithDynamicTopic: {
						ChainSpecificName: triggerWithDynamicTopic,
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							// No specific reason for filter being defined here instead of on contract level, this is just for test case variety.
							PollingFilter: &types.PollingFilter{},
						},
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"FieldHash": "Field"}},
						},
					},
					triggerWithAllTopics: {
						ChainSpecificName: triggerWithAllTopics,
						ReadType:          types.Event,
						EventDefinitions: &types.EventDefinitions{
							PollingFilter: &types.PollingFilter{},
						},
						// This doesn't have to be here, since the defalt mapping would work, but is left as an example.
						// Keys which are string float values(confidence levels) are chain agnostic and should be reused across chains.
						// These float values can map to different finality concepts across chains.
						ConfidenceConfirmations: map[string]int{"0.0": int(evmtypes.Unconfirmed), "1.0": int(evmtypes.Finalized)},
					},
					triggerWithAllTopicsWithHashed: {
						ChainSpecificName: triggerWithAllTopicsWithHashed,
						ReadType:          types.Event,
						EventDefinitions:  &types.EventDefinitions{},
					},
					MethodReturningSeenStruct: {
						ChainSpecificName: "returnSeen",
						InputModifications: codec.ModifiersConfig{
							&codec.HardCodeModifierConfig{
								OnChainValues: map[string]any{
									"BigField":              testStruct.BigField.String(),
									"AccountStruct.Account": hexutil.Encode(testStruct.AccountStruct.Account),
								},
							},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
							&codec.AddressBytesToStringModifierConfig{
								Fields: []string{"AccountStruct.AccountStr"},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"ExtraField": AnyExtraValue}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
							&codec.AddressBytesToStringModifierConfig{
								Fields: []string{"AccountStruct.AccountStr"},
							},
						},
					},
				},
			},
			AnySecondContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				Configs: map[string]*types.ChainReaderDefinition{
					MethodTakingLatestParamsReturningTestStruct: &methodTakingLatestParamsReturningTestStructConfig,
					MethodReturningUint64: {
						ChainSpecificName: "getDifferentPrimitiveValue",
					},
				},
			},
		},
	}
}

func (it *EVMChainComponentsInterfaceTester[T]) getChainWriterConfig(t T) types.ChainWriterConfig {
	return types.ChainWriterConfig{
		Contracts: map[string]*types.ContractConfig{
			AnyContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				Configs: map[string]*types.ChainWriterDefinition{
					"addTestStruct": {
						ChainSpecificName: "addTestStruct",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
						},
					},
					"setAlterablePrimitiveValue": {
						ChainSpecificName: "setAlterablePrimitiveValue",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
					},
					"triggerEvent": {
						ChainSpecificName: "triggerEvent",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
						},
					},
					"triggerEventWithDynamicTopic": {
						ChainSpecificName: "triggerEventWithDynamicTopic",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
					},
					"triggerWithFourTopics": {
						ChainSpecificName: "triggerWithFourTopics",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
					},
					"triggerWithFourTopicsWithHashed": {
						ChainSpecificName: "triggerWithFourTopicsWithHashed",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
					},
					"triggerStaticBytes": {
						ChainSpecificName: "triggerStaticBytes",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
					},
				},
			},
			AnySecondContractName: {
				ContractABI: chain_reader_tester.ChainReaderTesterMetaData.ABI,
				Configs: map[string]*types.ChainWriterDefinition{
					"addTestStruct": {
						ChainSpecificName: "addTestStruct",
						FromAddress:       it.Helper.Accounts(t)[1].From,
						GasLimit:          2_000_000,
						Checker:           "simulate",
						InputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedDynamicStruct.Inner.IntVal": "I"}},
							&codec.RenameModifierConfig{Fields: map[string]string{"NestedStaticStruct.Inner.IntVal": "I"}},
						},
					},
				},
			},
		},
		MaxGasPrice: assets.NewWei(big.NewInt(1000000000000000000)),
	}
}

func (it *EVMChainComponentsInterfaceTester[T]) Name() string {
	return "EVM"
}

func (it *EVMChainComponentsInterfaceTester[T]) GetAccountBytes(i int) []byte {
	account := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	account[i%20] += byte(i)
	account[(i+3)%20] += byte(i + 3)
	return account[:]
}

func (it *EVMChainComponentsInterfaceTester[T]) GetAccountString(i int) string {
	return common.BytesToAddress(it.GetAccountBytes(i)).Hex()
}

func (it *EVMChainComponentsInterfaceTester[T]) GetContractReader(t T) clcommontypes.ContractReader {
	ctx := it.Helper.Context(t)
	lggr := logger.Nop()

	cr, err := evm.NewChainReaderService(ctx, lggr, it.Helper.LogPoller(t), it.Helper.HeadTracker(t), it.client, it.chainReaderConfigSupplier(t))
	require.NoError(t, err)
	servicetest.Run(t, cr)

	return cr
}

func (it *EVMChainComponentsInterfaceTester[T]) GetContractWriter(t T) clcommontypes.ContractWriter {
	cw, err := evm.NewChainWriterService(logger.Nop(), it.client, it.Helper.TXM(t, it.client), nil, it.chainWriterConfigSupplier(t))
	require.NoError(t, err)

	cw = it.Helper.WrappedChainWriter(cw, it.client)

	require.NoError(t, err)

	servicetest.Run(t, cw)

	return cw
}

// This function is no longer necessary for Simulated Backend or Testnet tests.
func (it *EVMChainComponentsInterfaceTester[T]) GenerateBlocksTillConfidenceLevel(t T, contractName, readName string, confidenceLevel primitives.ConfidenceLevel) {
}

func (it *EVMChainComponentsInterfaceTester[T]) DirtyContracts() {
}

func (it *EVMChainComponentsInterfaceTester[T]) GetAuthWithGasSet(t T) *bind.TransactOpts {
	auth := *it.Helper.Accounts(t)[0]
	gasPrice, err := it.Helper.Client(t).SuggestGasPrice(it.Helper.Context(t))
	require.NoError(t, err)
	extra := new(big.Int).Mul(gasPrice, big.NewInt(it.Helper.GasPriceBufferPercent()))
	extra = extra.Div(extra, big.NewInt(100))
	auth.GasPrice = gasPrice.Add(gasPrice, extra)
	auth.GasLimit = 10552800
	nonce, err := it.client.PendingNonceAt(it.Helper.Context(t), auth.From)
	require.NoError(t, err)
	auth.Nonce = new(big.Int).SetUint64(nonce)
	return &auth
}

func (it *EVMChainComponentsInterfaceTester[T]) AwaitTx(t T, tx *gethtypes.Transaction) {
	ctx := it.Helper.Context(t)
	receipt, err := bind.WaitMined(ctx, it.Helper.Client(t), tx)
	require.NoError(t, err)
	require.Equal(t, gethtypes.ReceiptStatusSuccessful, receipt.Status)
}

func (it *EVMChainComponentsInterfaceTester[T]) deployNewContracts(t T) []clcommontypes.BoundContract {
	it.DeployLock.Lock()
	defer it.DeployLock.Unlock()
	address := it.deployNewContract(t)
	address2 := it.deployNewContract(t)
	return []clcommontypes.BoundContract{
		{Name: AnyContractName, Address: address},
		{Name: AnySecondContractName, Address: address2},
	}
}

func (it *EVMChainComponentsInterfaceTester[T]) deployNewContract(t T) string {
	address, tx, _, err := chain_reader_tester.DeployChainReaderTester(it.GetAuthWithGasSet(t), it.Helper.Backend())
	require.NoError(t, err)

	it.AwaitTx(t, tx)
	return address.String()
}

func (it *EVMChainComponentsInterfaceTester[T]) MaxWaitTimeForEvents() time.Duration {
	return it.Helper.MaxWaitTimeForEvents()
}

func (it *EVMChainComponentsInterfaceTester[T]) Setup(t T) {
	if it.chainReaderConfigSupplier == nil {
		it.chainReaderConfigSupplier = func(t T) types.ChainReaderConfig { return it.getChainReaderConfig(t) }
	}
	if it.chainWriterConfigSupplier == nil {
		it.chainWriterConfigSupplier = func(t T) types.ChainWriterConfig { return it.getChainWriterConfig(t) }
	}
	it.client = it.Helper.ChainReaderEVMClient(it.Helper.Context(t), t, it.Helper.HeadTracker(t), it.chainReaderConfigSupplier(t))
}

func (it *EVMChainComponentsInterfaceTester[T]) SetChainReaderConfigSupplier(chainReaderConfigSupplier func(t T) types.ChainReaderConfig) {
	it.chainReaderConfigSupplier = chainReaderConfigSupplier
}

func (it *EVMChainComponentsInterfaceTester[T]) SetChainWriterConfigSupplier(chainWriterConfigSupplier func(t T) types.ChainWriterConfig) {
	it.chainWriterConfigSupplier = chainWriterConfigSupplier
}

func OracleIDsToBytes(oracleIDs [32]commontypes.OracleID) [32]byte {
	convertedIDs := [32]byte{}
	for i, id := range oracleIDs {
		convertedIDs[i] = byte(id)
	}
	return convertedIDs
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
		Field:               *testStruct.Field,
		DifferentField:      testStruct.DifferentField,
		OracleId:            byte(testStruct.OracleID),
		OracleIds:           OracleIDsToBytes(testStruct.OracleIDs),
		AccountStruct:       AccountStructToInternalType(testStruct.AccountStruct),
		Accounts:            ConvertAccounts(testStruct.Accounts),
		BigField:            testStruct.BigField,
		NestedDynamicStruct: MidDynamicToInternalType(testStruct.NestedDynamicStruct),
		NestedStaticStruct:  MidStaticToInternalType(testStruct.NestedStaticStruct),
	}
}

func AccountStructToInternalType(a AccountStruct) chain_reader_tester.AccountStruct {
	return chain_reader_tester.AccountStruct{
		Account:    common.Address(a.Account),
		AccountStr: common.HexToAddress(a.AccountStr),
	}
}

func MidDynamicToInternalType(m MidLevelDynamicTestStruct) chain_reader_tester.MidLevelDynamicTestStruct {
	return chain_reader_tester.MidLevelDynamicTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: chain_reader_tester.InnerDynamicTestStruct{
			IntVal: int64(m.Inner.I),
			S:      m.Inner.S,
		},
	}
}

func MidStaticToInternalType(m MidLevelStaticTestStruct) chain_reader_tester.MidLevelStaticTestStruct {
	return chain_reader_tester.MidLevelStaticTestStruct{
		FixedBytes: m.FixedBytes,
		Inner: chain_reader_tester.InnerStaticTestStruct{
			IntVal: int64(m.Inner.I),
			A:      common.BytesToAddress(m.Inner.A),
		},
	}
}
