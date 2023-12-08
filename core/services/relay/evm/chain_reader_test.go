package evm_test

//go:generate ./testfiles/chainlink_reader_test_setup.sh

import (
	"crypto/ecdsa"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/testfiles"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const commonGasLimitOnEvms = uint64(4712388)
const chainReaderContractName = "LatestValueHolder"
const returnSeenName = "ReturnSeen"

func TestChainReader(t *testing.T) {
	RunChainReaderInterfaceTests(t, &chainReaderInterfaceTester{})
}

type chainReaderInterfaceTester struct {
	chain       *mocks.Chain
	address     string
	chainConfig types.ChainReaderConfig
	auth        *bind.TransactOpts
	sim         *backends.SimulatedBackend
	pk          *ecdsa.PrivateKey
	evmTest     *testfiles.Testfiles
	cr          evm.ChainReaderService
}

func (it *chainReaderInterfaceTester) Setup(t *testing.T) {
	t.Cleanup(func() {
		it.address = ""
		require.NoError(t, it.cr.Close())
		it.cr = nil
	})

	// can re-use the same chain for tests, just make new contract for each test
	if it.chain != nil {
		it.deployNewContract(t)
		return
	}

	it.chain = &mocks.Chain{}
	it.setupChainNoClient(t)

	testStruct := CreateTestStruct(0, it)

	it.chainConfig = types.ChainReaderConfig{
		ChainContractReaders: map[string]types.ChainContractReader{
			chainReaderContractName: {
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
					EventName: {
						ChainSpecificName: "Triggered",
						ReadType:          types.Event,
					},
					MethodReturningSeenStruct: {
						ChainSpecificName: returnSeenName,
						InputModifications: codec.ModifiersConfig{
							&codec.HardCodeConfig{
								OnChainValues: map[string]any{
									"BigField": testStruct.BigField.String(),
									"Account":  hexutil.Encode(testStruct.Account),
								},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.HardCodeConfig{
								OffChainValues: map[string]any{"ExtraField": anyExtraValue}},
						},
					},
				},
			},
		},
	}
	it.chain.On("Client").Return(client.NewSimulatedBackendClient(t, it.sim, big.NewInt(1337)))
	it.deployNewContract(t)
}

func (it *chainReaderInterfaceTester) Name() string {
	return "EVM"
}

func (it *chainReaderInterfaceTester) GetAccountBytes(i int) []byte {
	account := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2}
	account[i%32] += byte(i)
	account[(i+3)%32] += byte(i + 3)
	return account
}

func (it *chainReaderInterfaceTester) GetChainReader(t *testing.T) clcommontypes.ChainReader {
	ctx := testutils.Context(t)
	if it.cr != nil {
		return it.cr
	}

	addr := common.HexToAddress(it.address)
	lggr := logger.NullLogger
	db := pgtest.NewSqlxDB(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.SimulatedChainID, db, lggr, pgtest.NewQConfig(true)), it.chain.Client(), lggr, time.Millisecond, false, 0, 1, 1, 10000)
	require.NoError(t, lp.Start(ctx))
	it.chain.On("LogPoller").Return(lp)
	cr, err := evm.NewChainReaderService(lggr, lp, addr, it.chain, it.chainConfig)
	require.NoError(t, err)
	require.NoError(t, cr.Start(ctx))
	it.cr = cr
	return cr
}

func (it *chainReaderInterfaceTester) GetPrimitiveContract(_ *testing.T) clcommontypes.BoundContract {
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    MethodReturningUint64,
	}
}

func (it *chainReaderInterfaceTester) GetReturnSeenContract(t *testing.T) clcommontypes.BoundContract {
	it.deployNewContract(t)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    returnSeenName,
	}
}
func (it *chainReaderInterfaceTester) GetSliceContract(t *testing.T) clcommontypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(t)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    MethodReturningUint64Slice,
	}
}

func (it *chainReaderInterfaceTester) SetLatestValue(t *testing.T, testStruct *TestStruct) clcommontypes.BoundContract {
	it.sendTxWithTestStruct(t, testStruct, (*testfiles.TestfilesTransactor).AddTestStruct)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    MethodTakingLatestParamsReturningTestStruct,
	}
}

func (it *chainReaderInterfaceTester) TriggerEvent(t *testing.T, testStruct *TestStruct) clcommontypes.BoundContract {
	it.sendTxWithTestStruct(t, testStruct, (*testfiles.TestfilesTransactor).TriggerEvent)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    EventName,
	}
}

type testStructFn = func(*testfiles.TestfilesTransactor, *bind.TransactOpts, int32, string, uint8, [32]uint8, [32]byte, [][32]byte, *big.Int, testfiles.MidLevelTestStruct) (*evmtypes.Transaction, error)

func (it *chainReaderInterfaceTester) sendTxWithTestStruct(t *testing.T, testStruct *TestStruct, fn testStructFn) {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(t)

	tx, err := fn(
		&it.evmTest.TestfilesTransactor,
		it.auth,
		testStruct.Field,
		testStruct.DifferentField,
		uint8(testStruct.OracleID),
		convertOracleIDs(testStruct.OracleIDs),
		[32]byte(testStruct.Account),
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

func convertAccounts(accounts [][]byte) [][32]byte {
	convertedAccounts := make([][32]byte, len(accounts))
	for i, a := range accounts {
		convertedAccounts[i] = [32]byte(a)
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

func (it *chainReaderInterfaceTester) deployNewContract(t *testing.T) {
	ctx := testutils.Context(t)
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
	it.awaitTx(t, tx)
	it.address = address.String()
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
		uint8(ts.OracleID),
		getOracleIDs(ts),
		[32]byte(ts.Account),
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

func toInternalType(testStruct TestStruct) testfiles.TestStruct {
	return testfiles.TestStruct{
		Field:          testStruct.Field,
		DifferentField: testStruct.DifferentField,
		OracleId:       byte(testStruct.OracleID),
		OracleIds:      convertOracleIDs(testStruct.OracleIDs),
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
