package evmtesting

import (
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .
)

func RunChainReaderEvmTests[T TestingT[T]](t T, it *EVMChainReaderInterfaceTester[T]) {
	RunChainReaderInterfaceTests[T](t, it)

	t.Run("Dynamically typed topics can be used to filter and have type correct in return", func(t T) {
		it.Setup(t)

		anyString := "foo"
		it.dirtyContracts = true
		tx, err := it.contractTesters[it.address].ChainReaderTesterTransactor.TriggerEventWithDynamicTopic(it.GetAuthWithGasSet(t), anyString)
		require.NoError(t, err)
		it.Helper.Commit()
		it.IncNonce()
		it.AwaitTx(t, tx)
		ctx := it.Helper.Context(t)

		cr := it.GetChainReader(t)
		require.NoError(t, cr.Bind(ctx, it.GetBindings(t)))

		input := struct{ Field string }{Field: anyString}
		tp := cr.(clcommontypes.ContractTypeProvider)
		output, err := tp.CreateContractType(AnyContractName, triggerWithDynamicTopic, false)
		require.NoError(t, err)
		rOutput := reflect.Indirect(reflect.ValueOf(output))

		require.Eventually(t, func() bool {
			return cr.GetLatestValue(ctx, AnyContractName, triggerWithDynamicTopic, primitives.Unconfirmed, input, output) == nil
		}, it.MaxWaitTimeForEvents(), 100*time.Millisecond)

		assert.Equal(t, &anyString, rOutput.FieldByName("Field").Interface())
		topic, err := abi.MakeTopics([]any{anyString})
		require.NoError(t, err)
		assert.Equal(t, &topic[0][0], rOutput.FieldByName("FieldHash").Interface())
	})

	t.Run("Multiple topics can filter together", func(t T) {
		it.Setup(t)
		it.dirtyContracts = true
		triggerFourTopics(t, it, int32(1), int32(2), int32(3))
		triggerFourTopics(t, it, int32(2), int32(2), int32(3))
		triggerFourTopics(t, it, int32(1), int32(3), int32(3))
		triggerFourTopics(t, it, int32(1), int32(2), int32(4))

		ctx := it.Helper.Context(t)
		cr := it.GetChainReader(t)
		require.NoError(t, cr.Bind(ctx, it.GetBindings(t)))
		var latest struct{ Field1, Field2, Field3 int32 }
		params := struct{ Field1, Field2, Field3 int32 }{Field1: 1, Field2: 2, Field3: 3}

		time.Sleep(it.MaxWaitTimeForEvents())

		require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, triggerWithAllTopics, primitives.Unconfirmed, params, &latest))
		assert.Equal(t, int32(1), latest.Field1)
		assert.Equal(t, int32(2), latest.Field2)
		assert.Equal(t, int32(3), latest.Field3)
	})

	t.Run("Filtering can be done on indexed topics that get hashed", func(t T) {
		it.Setup(t)
		it.dirtyContracts = true
		triggerFourTopicsWithHashed(t, it, "1", [32]uint8{2}, [32]byte{5})
		triggerFourTopicsWithHashed(t, it, "2", [32]uint8{2}, [32]byte{3})
		triggerFourTopicsWithHashed(t, it, "1", [32]uint8{3}, [32]byte{3})

		ctx := it.Helper.Context(t)
		cr := it.GetChainReader(t)
		require.NoError(t, cr.Bind(ctx, it.GetBindings(t)))
		var latest struct {
			Field3 [32]byte
		}
		params := struct {
			Field1 string
			Field2 [32]uint8
			Field3 [32]byte
		}{Field1: "1", Field2: [32]uint8{2}, Field3: [32]byte{5}}

		time.Sleep(it.MaxWaitTimeForEvents())
		require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, triggerWithAllTopicsWithHashed, primitives.Unconfirmed, params, &latest))
		// only checking Field3 topic makes sense since it isn't hashed, to check other fields we'd have to replicate solidity encoding and hashing
		assert.Equal(t, [32]uint8{5}, latest.Field3)
	})

	t.Run("Bind returns error on missing contract at address", func(t T) {
		it.Setup(t)

		addr := common.BigToAddress(big.NewInt(42))
		reader := it.GetChainReader(t)

		ctx := it.Helper.Context(t)
		err := reader.Bind(ctx, []clcommontypes.BoundContract{{Name: AnyContractName, Address: addr.Hex()}})

		require.ErrorIs(t, err, evm.NoContractExistsError{Address: addr})
	})
}

func triggerFourTopics[T TestingT[T]](t T, it *EVMChainReaderInterfaceTester[T], i1, i2, i3 int32) {
	tx, err := it.contractTesters[it.address].ChainReaderTesterTransactor.TriggerWithFourTopics(it.GetAuthWithGasSet(t), i1, i2, i3)
	require.NoError(t, err)
	require.NoError(t, err)
	it.Helper.Commit()
	it.IncNonce()
	it.AwaitTx(t, tx)
}

func triggerFourTopicsWithHashed[T TestingT[T]](t T, it *EVMChainReaderInterfaceTester[T], i1 string, i2 [32]uint8, i3 [32]byte) {
	tx, err := it.contractTesters[it.address].ChainReaderTesterTransactor.TriggerWithFourTopicsWithHashed(it.GetAuthWithGasSet(t), i1, i2, i3)
	require.NoError(t, err)
	require.NoError(t, err)
	it.Helper.Commit()
	it.IncNonce()
	it.AwaitTx(t, tx)
}
