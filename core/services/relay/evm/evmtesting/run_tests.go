package evmtesting

import (
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	commontestutils "github.com/smartcontractkit/chainlink-common/pkg/loop/testutils"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .
)

func RunChainReaderTests[T TestingT[T]](t T, it *EvmChainReaderInterfaceTester[T], loopTo bool) {
	RunChainReaderInterfaceTests[T](t, it)
	if loopTo {
		RunChainReaderInterfaceTests(t, commontestutils.WrapChainReaderTesterForLoop[T](it))
	}

	t.Run("Dynamically typed topics can be used to filter and have type correct in return", func(t T) {
		it.Setup(t)

		anyString := "foo"
		tx, err := it.evmTest.LatestValueHolderTransactor.TriggerEventWithDynamicTopic(it.auth, anyString)
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
			return cr.GetLatestValue(ctx, AnyContractName, triggerWithDynamicTopic, input, output) == nil
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)

		assert.Equal(t, &anyString, rOutput.FieldByName("Field").Interface())
		topic, err := abi.MakeTopics([]any{anyString})
		require.NoError(t, err)
		assert.Equal(t, &topic[0][0], rOutput.FieldByName("FieldHash").Interface())
	})

	t.Run("Multiple topics can filter together", func(t T) {
		it.Setup(t)
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

		require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, triggerWithAllTopics, params, &latest))
		assert.Equal(t, int32(1), latest.Field1)
		assert.Equal(t, int32(2), latest.Field2)
		assert.Equal(t, int32(3), latest.Field3)
	})
}

func triggerFourTopics[T TestingT[T]](t T, it *EvmChainReaderInterfaceTester[T], i1, i2, i3 int32) {
	tx, err := it.evmTest.LatestValueHolderTransactor.TriggerWithFourTopics(it.auth, i1, i2, i3)
	require.NoError(t, err)
	require.NoError(t, err)
	it.Helper.Commit()
	it.IncNonce()
	it.AwaitTx(t, tx)
}
