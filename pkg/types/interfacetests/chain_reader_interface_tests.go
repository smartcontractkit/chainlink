package interfacetests

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

type ChainReaderInterfaceTester interface {
	BasicTester
	GetChainReader(t *testing.T) types.ContractReader

	// SetLatestValue is expected to return the same bound contract and method in the same test
	// Any setup required for this should be done in Setup.
	// The contract should take a LatestParams as the params and return the nth TestStruct set
	SetLatestValue(t *testing.T, testStruct *TestStruct)
	TriggerEvent(t *testing.T, testStruct *TestStruct)
	GetBindings(t *testing.T) []types.BoundContract
	MaxWaitTimeForEvents() time.Duration
}

const (
	AnyValueToReadWithoutAnArgument             = uint64(3)
	AnyDifferentValueToReadWithoutAnArgument    = uint64(1990)
	MethodTakingLatestParamsReturningTestStruct = "GetLatestValues"
	MethodReturningUint64                       = "GetPrimitiveValue"
	DifferentMethodReturningUint64              = "GetDifferentPrimitiveValue"
	MethodReturningUint64Slice                  = "GetSliceValue"
	MethodReturningSeenStruct                   = "GetSeenStruct"
	EventName                                   = "SomeEvent"
	EventWithFilterName                         = "SomeEventToFilter"
	AnyContractName                             = "TestContract"
	AnySecondContractName                       = "Not" + AnyContractName
)

var AnySliceToReadWithoutAnArgument = []uint64{3, 4}

const AnyExtraValue = 3

func RunChainReaderInterfaceTests(t *testing.T, tester ChainReaderInterfaceTester) {
	t.Run("GetLatestValue for "+tester.Name(), func(t *testing.T) { runChainReaderGetLatestValueInterfaceTests(t, tester) })
	t.Run("QueryKey for "+tester.Name(), func(t *testing.T) { runQueryKeyInterfaceTests(t, tester) })
}

func runChainReaderGetLatestValueInterfaceTests(t *testing.T, tester ChainReaderInterfaceTester) {
	tests := []testcase{
		{
			name: "Gets the latest value",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				firstItem := CreateTestStruct(0, tester)
				tester.SetLatestValue(t, &firstItem)
				secondItem := CreateTestStruct(1, tester)
				tester.SetLatestValue(t, &secondItem)

				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				actual := &TestStruct{}
				params := &LatestParams{I: 1}
				require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, MethodTakingLatestParamsReturningTestStruct, params, actual))
				assert.Equal(t, &firstItem, actual)

				params.I = 2
				actual = &TestStruct{}
				require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, MethodTakingLatestParamsReturningTestStruct, params, actual))
				assert.Equal(t, &secondItem, actual)
			},
		},
		{
			name: "Get latest value without arguments and with primitive return",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				var prim uint64
				require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, MethodReturningUint64, nil, &prim))

				assert.Equal(t, AnyValueToReadWithoutAnArgument, prim)
			},
		},
		{
			name: "Get latest value allows a contract name to resolve different contracts internally",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				var prim uint64
				require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, DifferentMethodReturningUint64, nil, &prim))

				assert.Equal(t, AnyDifferentValueToReadWithoutAnArgument, prim)
			},
		},
		{
			name: "Get latest value allows multiple constract names to have the same function name",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				bindings := tester.GetBindings(t)
				seenAddrs := map[string]bool{}
				for _, binding := range bindings {
					assert.False(t, seenAddrs[binding.Address])
					seenAddrs[binding.Address] = true
				}

				require.NoError(t, cr.Bind(ctx, bindings))

				var prim uint64
				require.NoError(t, cr.GetLatestValue(ctx, AnySecondContractName, MethodReturningUint64, nil, &prim))

				assert.Equal(t, AnyDifferentValueToReadWithoutAnArgument, prim)
			},
		},
		{
			name: "Get latest value without arguments and with slice return",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				var slice []uint64
				require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, MethodReturningUint64Slice, nil, &slice))

				assert.Equal(t, AnySliceToReadWithoutAnArgument, slice)
			},
		},
		{
			name: "Get latest value wraps config with modifiers using its own mapstructure overrides",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				testStruct := CreateTestStruct(0, tester)
				testStruct.BigField = nil
				testStruct.Account = nil
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				actual := &TestStructWithExtraField{}
				require.NoError(t, cr.GetLatestValue(ctx, AnyContractName, MethodReturningSeenStruct, testStruct, actual))

				expected := &TestStructWithExtraField{
					ExtraField: AnyExtraValue,
					TestStruct: CreateTestStruct(0, tester),
				}

				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Get latest value gets latest event",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))
				ts := CreateTestStruct(0, tester)
				tester.TriggerEvent(t, &ts)
				ts = CreateTestStruct(1, tester)
				tester.TriggerEvent(t, &ts)

				result := &TestStruct{}
				assert.Eventually(t, func() bool {
					err := cr.GetLatestValue(ctx, AnyContractName, EventName, nil, &result)
					return err == nil && reflect.DeepEqual(result, &ts)
				}, tester.MaxWaitTimeForEvents(), time.Millisecond*10)
			},
		},
		{
			name: "Get latest value returns not found if event was never triggered",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				result := &TestStruct{}
				err := cr.GetLatestValue(ctx, AnyContractName, EventName, nil, &result)
				assert.True(t, errors.Is(err, types.ErrNotFound))
			},
		},
		{
			name: "Get latest value gets latest event with filtering",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))
				ts0 := CreateTestStruct(0, tester)
				tester.TriggerEvent(t, &ts0)
				ts1 := CreateTestStruct(1, tester)
				tester.TriggerEvent(t, &ts1)

				filterParams := &FilterEventParams{Field: *ts0.Field}
				assert.Never(t, func() bool {
					result := &TestStruct{}
					err := cr.GetLatestValue(ctx, AnyContractName, EventWithFilterName, filterParams, &result)
					return err == nil && reflect.DeepEqual(result, &ts1)
				}, tester.MaxWaitTimeForEvents(), time.Millisecond*10)
				// get the result one more time to verify it.
				// Using the result from the Never statement by creating result outside the block is a data race
				result := &TestStruct{}
				err := cr.GetLatestValue(ctx, AnyContractName, EventWithFilterName, filterParams, &result)
				require.NoError(t, err)
				assert.Equal(t, &ts0, result)
			},
		},
	}
	runTests(t, tester, tests)
}

func runQueryKeyInterfaceTests(t *testing.T, tester ChainReaderInterfaceTester) {
	tests := []testcase{
		{
			name: "QueryKey returns not found if sequence never happened",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)

				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))

				logs, err := cr.QueryKey(ctx, AnyContractName, query.KeyFilter{Key: EventName}, query.LimitAndSort{}, &TestStruct{})

				require.NoError(t, err)
				assert.Len(t, logs, 0)
			},
		},
		{
			name: "QueryKey returns sequence data properly",
			test: func(t *testing.T) {
				ctx := tests.Context(t)
				cr := tester.GetChainReader(t)
				require.NoError(t, cr.Bind(ctx, tester.GetBindings(t)))
				ts1 := CreateTestStruct(0, tester)
				tester.TriggerEvent(t, &ts1)
				ts2 := CreateTestStruct(1, tester)
				tester.TriggerEvent(t, &ts2)

				ts := &TestStruct{}
				assert.Eventually(t, func() bool {
					sequences, err := cr.QueryKey(ctx, AnyContractName, query.KeyFilter{Key: EventName}, query.LimitAndSort{}, ts)
					return err == nil && len(sequences) == 2 && reflect.DeepEqual(&ts1, sequences[0].Data) && reflect.DeepEqual(&ts2, sequences[1].Data)
				}, tester.MaxWaitTimeForEvents(), time.Millisecond*10)
			},
		},
	}

	runTests(t, tester, tests)
}
