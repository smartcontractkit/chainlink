package interfacetests

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

type ChainReaderInterfaceTester interface {
	Setup(t *testing.T)
	Name() string
	GetAccountBytes(i int) []byte
	GetChainReader(t *testing.T) types.ChainReader

	// SetLatestValue is expected to return the same bound contract and method in the same test
	// Any setup required for this should be done in Setup.
	// The contract should take a LatestParams as the params and return the nth TestStruct set
	SetLatestValue(ctx context.Context, t *testing.T, testStruct *TestStruct) types.BoundContract
	GetPrimitiveContract(ctx context.Context, t *testing.T) types.BoundContract
	GetSliceContract(ctx context.Context, t *testing.T) types.BoundContract
}

const (
	MethodTakingLatestParamsReturningTestStruct = "GetLatestValues"
)

var AnySliceToReadWithoutAnArgument = []uint64{3, 4}

// RunChainReaderInterfaceTests uses TestStruct and TestStructWithSpecialFields
func RunChainReaderInterfaceTests(t *testing.T, tester ChainReaderInterfaceTester) {
	ctx := tests.Context(t)
	tests := map[string]func(t *testing.T){
		"Gets the latest value": func(t *testing.T) {
			firstItem := CreateTestStruct(0, tester.GetAccountBytes)
			bc := tester.SetLatestValue(ctx, t, &firstItem)
			secondItem := CreateTestStruct(1, tester.GetAccountBytes)
			tester.SetLatestValue(ctx, t, &secondItem)

			cr := tester.GetChainReader(t)
			actual := &TestStruct{}
			params := &LatestParams{I: 1}

			require.NoError(t, cr.GetLatestValue(ctx, bc, MethodTakingLatestParamsReturningTestStruct, params, actual))
			assert.Equal(t, &firstItem, actual)

			params.I = 2
			actual = &TestStruct{}
			require.NoError(t, cr.GetLatestValue(ctx, bc, MethodTakingLatestParamsReturningTestStruct, params, actual))
			assert.Equal(t, &secondItem, actual)
		},
	}

	runTests(t, tester, tests)
}

func runTests(t *testing.T, tester ChainReaderInterfaceTester, tests map[string]func(t *testing.T)) {
	// Order the tests for consistency
	testNames := make([]string, 0, len(tests))
	for name := range tests {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)

	for i := 0; i < len(testNames); i++ {
		name := testNames[i]
		t.Run(name, func(t *testing.T) {
			tester.Setup(t)
			tests[name](t)
		})
	}
}

type InnerTestStruct struct {
	I int
	S string
}

type MidLevelTestStruct struct {
	FixedBytes [2]byte
	Inner      InnerTestStruct
}

type TestStruct struct {
	Field          int32
	DifferentField string
	OracleID       commontypes.OracleID
	OracleIDs      [32]commontypes.OracleID
	Account        []byte
	Accounts       [][]byte
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
}

type LatestParams struct {
	I int
}

func CreateTestStruct(i int, accGen func(int) []byte) TestStruct {
	s := fmt.Sprintf("field%v", i)
	return TestStruct{
		Field:          int32(i),
		DifferentField: s,
		OracleID:       commontypes.OracleID(i + 1),
		OracleIDs:      [32]commontypes.OracleID{commontypes.OracleID(i + 2), commontypes.OracleID(i + 3)},
		Account:        accGen(i + 3),
		Accounts:       [][]byte{accGen(i + 4), accGen(i + 5)},
		BigField:       big.NewInt(int64((i + 1) * (i + 2))),
		NestedStruct: MidLevelTestStruct{
			FixedBytes: [2]byte{uint8(i), uint8(i + 1)},
			Inner: InnerTestStruct{
				I: i,
				S: s,
			},
		},
	}
}
