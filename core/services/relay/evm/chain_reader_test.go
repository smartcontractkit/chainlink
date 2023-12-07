package evm_test

//go:generate ./testfiles/chainlink_reader_test_setup.sh

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
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
	mocklogpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm/mocks"
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
}

func (it *chainReaderInterfaceTester) Setup(t *testing.T) {
	// can re-use the same chain for tests, just make new contract for each test
	if it.chain != nil {
		return
	}

	t.Cleanup(func() { it.address = "" })
	it.chain = &mocks.Chain{}
	it.setupChainNoClient(t)
	it.chain.On("LogPoller").Return(logger.NullLogger)

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
	cr, err := evm.NewChainReaderService(logger.NullLogger, mocklogpoller.NewLogPoller(t), it.chain, it.chainConfig)
	require.NoError(t, err)
	return cr
}

func (it *chainReaderInterfaceTester) GetPrimitiveContract(ctx context.Context, t *testing.T) clcommontypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(ctx, t)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    MethodReturningUint64,
	}
}

func (it *chainReaderInterfaceTester) GetReturnSeenContract(ctx context.Context, t *testing.T) clcommontypes.BoundContract {
	it.deployNewContract(ctx, t)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    returnSeenName,
	}
}
func (it *chainReaderInterfaceTester) GetSliceContract(ctx context.Context, t *testing.T) clcommontypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(ctx, t)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    MethodReturningUint64Slice,
	}
}

func (it *chainReaderInterfaceTester) SetLatestValue(ctx context.Context, t *testing.T, testStruct *TestStruct) clcommontypes.BoundContract {
	// Since most tests don't use the contract, it's set up lazily to save time
	it.deployNewContract(ctx, t)

	tx, err := it.evmTest.AddTestStruct(
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
	it.awaitTx(ctx, t, tx)
	return clcommontypes.BoundContract{
		Address: it.address,
		Name:    MethodTakingLatestParamsReturningTestStruct,
	}
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

func (it *chainReaderInterfaceTester) deployNewContract(ctx context.Context, t *testing.T) {
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

func (it *chainReaderInterfaceTester) awaitTx(ctx context.Context, t *testing.T, tx *evmtypes.Transaction) {
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
