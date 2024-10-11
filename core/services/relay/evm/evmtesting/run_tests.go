package evmtesting

import (
	"encoding/binary"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	clcommontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/read"

	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests" //nolint common practice to import test mods with .
)

func RunChainComponentsEvmTests[T TestingT[T]](t T, it *EVMChainComponentsInterfaceTester[T]) {
	RunContractReaderEvmTests[T](t, it)
	// Add ChainWriter tests here
}

func RunChainComponentsInLoopEvmTests[T TestingT[T]](t T, it ChainComponentsInterfaceTester[T]) {
	RunContractReaderInLoopTests[T](t, it)
	// Add ChainWriter tests here
}

func RunContractReaderEvmTests[T TestingT[T]](t T, it *EVMChainComponentsInterfaceTester[T]) {
	RunContractReaderInterfaceTests[T](t, it, false)

	t.Run("Dynamically typed topics can be used to filter and have type correct in return", func(t T) {
		it.Setup(t)

		anyString := "foo"
		ctx := it.Helper.Context(t)

		cr := it.GetContractReader(t)
		bindings := it.GetBindings(t)
		require.NoError(t, cr.Bind(ctx, bindings))

		type DynamicEvent struct {
			Field string
		}
		SubmitTransactionToCW(t, it, "triggerEventWithDynamicTopic", DynamicEvent{Field: anyString}, bindings[0], types.Unconfirmed)

		input := struct{ Field string }{Field: anyString}
		tp := cr.(clcommontypes.ContractTypeProvider)

		readName := types.BoundContract{
			Address: bindings[0].Address,
			Name:    AnyContractName,
		}.ReadIdentifier(triggerWithDynamicTopic)

		output, err := tp.CreateContractType(readName, false)
		require.NoError(t, err)
		rOutput := reflect.Indirect(reflect.ValueOf(output))

		require.Eventually(t, func() bool {
			return cr.GetLatestValue(ctx, readName, primitives.Unconfirmed, input, output) == nil
		}, it.MaxWaitTimeForEvents(), 100*time.Millisecond)

		assert.Equal(t, &anyString, rOutput.FieldByName("Field").Interface())
		topic, err := abi.MakeTopics([]any{anyString})
		require.NoError(t, err)
		assert.Equal(t, &topic[0][0], rOutput.FieldByName("FieldHash").Interface())
	})

	t.Run("Multiple topics can filter together", func(t T) {
		it.Setup(t)
		ctx := it.Helper.Context(t)
		cr := it.GetContractReader(t)
		bindings := it.GetBindings(t)

		require.NoError(t, cr.Bind(ctx, bindings))

		triggerFourTopics(t, it, int32(1), int32(2), int32(3))
		triggerFourTopics(t, it, int32(2), int32(2), int32(3))
		triggerFourTopics(t, it, int32(1), int32(3), int32(3))
		triggerFourTopics(t, it, int32(1), int32(2), int32(4))

		var bound types.BoundContract
		for idx := range bindings {
			if bindings[idx].Name == AnyContractName {
				bound = bindings[idx]
			}
		}

		var latest struct{ Field1, Field2, Field3 int32 }
		params := struct{ Field1, Field2, Field3 int32 }{Field1: 1, Field2: 2, Field3: 3}

		time.Sleep(it.MaxWaitTimeForEvents())

		require.NoError(t, cr.GetLatestValue(ctx, bound.ReadIdentifier(triggerWithAllTopics), primitives.Unconfirmed, params, &latest))
		assert.Equal(t, int32(1), latest.Field1)
		assert.Equal(t, int32(2), latest.Field2)
		assert.Equal(t, int32(3), latest.Field3)
	})

	t.Run("Filtering can be done on indexed topics that get hashed", func(t T) {
		it.Setup(t)

		cr := it.GetContractReader(t)
		ctx := it.Helper.Context(t)
		bindings := it.GetBindings(t)

		require.NoError(t, cr.Bind(ctx, bindings))

		triggerFourTopicsWithHashed(t, it, "1", [32]uint8{2}, [32]byte{5})
		triggerFourTopicsWithHashed(t, it, "2", [32]uint8{2}, [32]byte{3})
		triggerFourTopicsWithHashed(t, it, "1", [32]uint8{3}, [32]byte{3})

		var bound types.BoundContract
		for idx := range bindings {
			if bindings[idx].Name == AnyContractName {
				bound = bindings[idx]
			}
		}

		var latest struct {
			Field3 [32]byte
		}
		params := struct {
			Field1 string
			Field2 [32]uint8
			Field3 [32]byte
		}{Field1: "1", Field2: [32]uint8{2}, Field3: [32]byte{5}}

		time.Sleep(it.MaxWaitTimeForEvents())
		require.NoError(t, cr.GetLatestValue(ctx, bound.ReadIdentifier(triggerWithAllTopicsWithHashed), primitives.Unconfirmed, params, &latest))
		// only checking Field3 topic makes sense since it isn't hashed, to check other fields we'd have to replicate solidity encoding and hashing
		assert.Equal(t, [32]uint8{5}, latest.Field3)
	})

	t.Run("Bind returns error on missing contract at address", func(t T) {
		it.Setup(t)

		addr := common.BigToAddress(big.NewInt(42))
		reader := it.GetContractReader(t)

		ctx := it.Helper.Context(t)
		err := reader.Bind(ctx, []clcommontypes.BoundContract{{Name: AnyContractName, Address: addr.Hex()}})

		require.ErrorIs(t, err, read.NoContractExistsError{Err: clcommontypes.ErrInternal, Address: addr})
	})
}

func RunContractReaderInLoopTests[T TestingT[T]](t T, it ChainComponentsInterfaceTester[T]) {
	RunContractReaderInterfaceTests[T](t, it, false)

	it.Setup(t)
	ctx := tests.Context(t)
	cr := it.GetContractReader(t)
	require.NoError(t, cr.Bind(ctx, it.GetBindings(t)))
	bindings := it.GetBindings(t)
	boundContract := BindingsByName(bindings, AnyContractName)[0]
	require.NoError(t, cr.Bind(ctx, bindings))

	ts1 := CreateTestStruct[T](0, it)
	_ = SubmitTransactionToCW(t, it, MethodTriggeringEvent, ts1, boundContract, types.Unconfirmed)
	ts2 := CreateTestStruct[T](15, it)
	_ = SubmitTransactionToCW(t, it, MethodTriggeringEvent, ts2, boundContract, types.Unconfirmed)
	ts3 := CreateTestStruct[T](35, it)
	_ = SubmitTransactionToCW(t, it, MethodTriggeringEvent, ts3, boundContract, types.Unconfirmed)

	t.Run("Filtering can be done on data words using value comparator", func(t T) {
		ts := &TestStruct{}
		assert.Eventually(t, func() bool {
			sequences, err := cr.QueryKey(ctx, boundContract, query.KeyFilter{Key: EventName, Expressions: []query.Expression{
				query.Comparator("OracleID",
					primitives.ValueComparator{
						Value:    uint8(ts2.OracleID),
						Operator: primitives.Eq,
					}),
			},
			}, query.LimitAndSort{}, ts)
			return err == nil && len(sequences) == 1 && reflect.DeepEqual(&ts2, sequences[0].Data)
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)
	})

	t.Run("Filtering can be done on data words using value comparator on a nested field", func(t T) {
		ts := &TestStruct{}
		assert.Eventually(t, func() bool {
			sequences, err := cr.QueryKey(ctx, boundContract, query.KeyFilter{Key: EventName, Expressions: []query.Expression{
				query.Comparator("OracleID",
					primitives.ValueComparator{
						Value:    uint8(ts2.OracleID),
						Operator: primitives.Eq,
					}),
				query.Comparator("NestedStaticStruct.Inner.IntVal",
					primitives.ValueComparator{
						Value:    ts2.NestedStaticStruct.Inner.I,
						Operator: primitives.Eq,
					}),
			},
			}, query.LimitAndSort{}, ts)
			return err == nil && len(sequences) == 1 && reflect.DeepEqual(&ts2, sequences[0].Data)
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)
	})

	t.Run("Filtering can be done on data words using value comparator on field that follows a dynamic field", func(t T) {
		ts := &TestStruct{}
		assert.Eventually(t, func() bool {
			sequences, err := cr.QueryKey(ctx, boundContract, query.KeyFilter{Key: EventName, Expressions: []query.Expression{
				query.Comparator("OracleID",
					primitives.ValueComparator{
						Value:    uint8(ts2.OracleID),
						Operator: primitives.Eq,
					}),
				query.Comparator("BigField",
					primitives.ValueComparator{
						Value:    ts2.BigField,
						Operator: primitives.Eq,
					}),
			},
			}, query.LimitAndSort{}, ts)
			return err == nil && len(sequences) == 1 && reflect.DeepEqual(&ts2, sequences[0].Data)
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)
	})

	t.Run("Filtering can be done on data words using value comparators on fields that require manual index input", func(t T) {
		empty12Bytes := [12]byte{}
		val1, val2, val3, val4 := uint32(1), uint32(2), uint32(3), uint64(4)
		val5, val6, val7 := [32]byte{}, [32]byte{6}, [32]byte{7}
		copy(val5[:], append(empty12Bytes[:], 5))
		raw := []byte{9, 8}

		var buf []byte
		buf = binary.BigEndian.AppendUint32(buf, val1)
		buf = binary.BigEndian.AppendUint32(buf, val2)
		buf = binary.BigEndian.AppendUint32(buf, val3)
		buf = binary.BigEndian.AppendUint64(buf, val4)
		dataWordOnChainValueToQuery := buf[:]

		resExpected := append(buf, common.LeftPadBytes(val5[:], 32)...)
		resExpected = append(resExpected, common.LeftPadBytes(val6[:], 32)...)
		resExpected = append(resExpected, common.LeftPadBytes(val7[:], 32)...)
		resExpected = append(resExpected, raw...)

		type eventResAsStruct struct {
			Message *[]uint8
		}
		wrapExpectedRes := eventResAsStruct{Message: &resExpected}

		// emit the one we want to search for and a couple of random ones to confirm that filtering works
		triggerStaticBytes(t, it, val1, val2, val3, val4, val5, val6, val7, raw)
		triggerStaticBytes(t, it, 1337, 7331, 4747, val4, val5, val6, val7, raw)
		triggerStaticBytes(t, it, 7331, 4747, 1337, val4, val5, val6, val7, raw)
		triggerStaticBytes(t, it, 4747, 1337, 7331, val4, val5, val6, val7, raw)

		assert.Eventually(t, func() bool {
			sequences, err := cr.QueryKey(ctx, boundContract, query.KeyFilter{Key: staticBytesEventName, Expressions: []query.Expression{
				query.Comparator("msgTransmitterEvent",
					primitives.ValueComparator{
						Value:    dataWordOnChainValueToQuery,
						Operator: primitives.Eq,
					}),
			},
			}, query.LimitAndSort{}, eventResAsStruct{})
			return err == nil && len(sequences) == 1 && reflect.DeepEqual(wrapExpectedRes, sequences[0].Data)
		}, it.MaxWaitTimeForEvents(), time.Millisecond*10)
	})
}

func triggerFourTopics[T TestingT[T]](t T, it *EVMChainComponentsInterfaceTester[T], i1, i2, i3 int32) {
	type DynamicEvent struct {
		Field1 int32
		Field2 int32
		Field3 int32
	}
	contracts := it.GetBindings(t)
	SubmitTransactionToCW(t, it, "triggerWithFourTopics", DynamicEvent{Field1: i1, Field2: i2, Field3: i3}, contracts[0], types.Unconfirmed)
}

func triggerFourTopicsWithHashed[T TestingT[T]](t T, it *EVMChainComponentsInterfaceTester[T], i1 string, i2 [32]uint8, i3 [32]byte) {
	type DynamicEvent struct {
		Field1 string
		Field2 [32]uint8
		Field3 [32]byte
	}
	contracts := it.GetBindings(t)
	SubmitTransactionToCW(t, it, "triggerWithFourTopicsWithHashed", DynamicEvent{Field1: i1, Field2: i2, Field3: i3}, contracts[0], types.Unconfirmed)
}

// triggerStaticBytes emits a staticBytes events and returns the expected event bytes.
func triggerStaticBytes[T TestingT[T]](t T, it ChainComponentsInterfaceTester[T], val1, val2, val3 uint32, val4 uint64, val5, val6, val7 [32]byte, raw []byte) {
	type StaticBytesEvent struct {
		Val1 uint32
		Val2 uint32
		Val3 uint32
		Val4 uint64
		Val5 [32]byte
		Val6 [32]byte
		Val7 [32]byte
		Raw  []byte
	}

	contracts := it.GetBindings(t)
	SubmitTransactionToCW(t, it, "triggerStaticBytes",
		StaticBytesEvent{
			Val1: val1,
			Val2: val2,
			Val3: val3,
			Val4: val4,
			Val5: val5,
			Val6: val6,
			Val7: val7,
			Raw:  raw,
		},
		contracts[0], types.Unconfirmed)
}
