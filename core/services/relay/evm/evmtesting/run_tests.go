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

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
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
