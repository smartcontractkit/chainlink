package evm_test

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/chain_reader_tester"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/bindings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"

	_ "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"       //nolint common practice to import test mods with .
	. "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/evmtesting" //nolint common practice to import test mods with .
)

//go:generate evm-bindings -output core/services/relay/evm/bindings -contracts contracts/src/v0.8/shared/test/helpers
func TestGoBindings(t *testing.T) {
	t.Parallel()

	chainReaderConfig := bindings.NewChainReaderConfig()
	maxGasPrice := assets.NewWei(big.NewInt(1000000000))
	chainWriterConfig := bindings.NewChainWriterConfig(*maxGasPrice, 2_000_000)

	it := &EVMChainComponentsInterfaceTester[*testing.T]{
		Helper: &RpcHelper{
			DeployPrivateKey: "PUT_HERE_PRIVATE_KEY",
			SenderPrivateKey: "PUT_HERE_PRIVATE_KEY",
		},
		ChainReaderConfig: &chainReaderConfig,
		ChainWriterConfig: &chainWriterConfig,
	}

	it.Helper.Init(t)

	t.Run("Deploy contract set value using cw and get value using chain reader", func(t *testing.T) {
		it.Setup(t)

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()
		senderOpts := it.GetAuthWithGasSet(t)
		senderOpts.Value = big.NewInt(0)
		senderOpts.GasLimit = 10_000_000
		senderOpts.Context = ctx
		chainWriterConfig.Contracts["ChainReaderTester"].Configs["SetAlterablePrimitiveValue"].FromAddress = senderOpts.From
		chainWriterConfig.Contracts["ChainReaderTester"].Configs["AddTestStruct"].FromAddress = senderOpts.From

		targetContractAddress, _, _, err := chain_reader_tester.DeployChainReaderTester(it.GetAuthWithGasSet(t), it.Helper.Backend())
		require.NoError(t, err)

		chainReader := it.GetChainReader(t)
		chainWriter := it.GetChainWriter(t)
		testContext := it.Helper.Context(t)

		//TODO properly wait for contract deployment confirmation
		time.Sleep(time.Duration(5) * time.Second)
		require.NoError(t, chainReader.Bind(testContext, []types.BoundContract{{targetContractAddress.String(), "ChainReaderTester"}}))
		it.IncNonce()

		chainReaderTester := bindings.ChainReaderTester{
			ContractReader: chainReader,
			ChainWriter:    chainWriter,
		}

		crt, err := chain_reader_tester.NewChainReaderTester(targetContractAddress, it.Helper.Client(t))
		require.NoError(t, err)

		callOpts := bind.CallOpts{
			Pending: true,
			From:    senderOpts.From,
			Context: ctx,
		}
		testStruct, err := crt.GetElementAtIndex(&callOpts, big.NewInt(5))
		require.NoError(t, err)
		require.Equal(t, testStruct.Field, int32(5))

		txId, err := uuid.NewRandom()
		require.NoError(t, err)
		err = chainReaderTester.SetAlterablePrimitiveValue(testContext, bindings.SetAlterablePrimitiveValueInput{Value: uint64(100)}, txId.String(), targetContractAddress.String(), nil)
		require.NoError(t, err)

		err = interfacetests.WaitForTransactionStatus(t, it, txId.String(), types.Unconfirmed, false)
		require.NoError(t, err)

		value, err := chainReaderTester.GetAlterablePrimitiveValue(testContext, primitives.Unconfirmed)
		require.NoError(t, err)
		assert.Equal(t, value, uint64(100))

		addTestStructTxId, err := uuid.NewRandom()
		require.NoError(t, err)
		err = chainReaderTester.AddTestStruct(ctx, bindings.AddTestStructInput{
			Field:          7,
			DifferentField: "anotherField",
			OracleId:       12,
			OracleIds:      [32]uint8{0, 2, 3},
			Account:        targetContractAddress,
			Accounts:       []common.Address{targetContractAddress, targetContractAddress},
			BigField:       *big.NewInt(100),
			NestedStruct: bindings.MidLevelTestStruct{
				FixedBytes: [2]uint8{7, 6},
				Inner: bindings.InnerTestStruct{
					IntVal: 200,
					S:      "super inner string",
				},
			},
		}, addTestStructTxId.String(), targetContractAddress.String(), nil)
		require.NoError(t, err)

		err = interfacetests.WaitForTransactionStatus(t, it, addTestStructTxId.String(), types.Unconfirmed, false)
		require.NoError(t, err)

		lastElementAtIndexOutput, err := chainReaderTester.GetElementAtIndex(ctx, bindings.GetElementAtIndexInput{I: 7}, primitives.Unconfirmed)
		require.NoError(t, err)

		require.Equal(t, lastElementAtIndexOutput.Field, int32(5))
	})
}

func TestGoBindingsTxWithGeth(t *testing.T) {
	t.Parallel()

	chainReaderConfig := bindings.NewChainReaderConfig()

	it := &EVMChainComponentsInterfaceTester[*testing.T]{Helper: &Helper{},
		ChainReaderConfig: &chainReaderConfig}

	it.Helper.Init(t)
	t.Run("Deploy contract set value using geth and get value using chain reader", func(t *testing.T) {
		it.Setup(t)

		ctx := tests.Context(t)

		address, _, deployedContract, err := chain_reader_tester.DeployChainReaderTester(it.GetAuthWithGasSet(t), it.Helper.Backend())
		require.NoError(t, err)
		it.IncNonce()
		_, err = deployedContract.SetAlterablePrimitiveValue(it.GetAuthWithGasSet(t), 100)
		require.NoError(t, err)

		time.Sleep(time.Duration(5) * time.Second)

		chainReader := it.GetChainReader(t)
		require.NoError(t, chainReader.Bind(ctx, []types.BoundContract{{address.String(), "ChainReaderTester"}}))

		chainReaderTester := bindings.ChainReaderTester{
			ContractReader: chainReader,
		}

		value, err := chainReaderTester.GetAlterablePrimitiveValue(ctx, primitives.Unconfirmed)
		require.NoError(t, err)
		assert.Equal(t, value, uint64(100))

	})
}
