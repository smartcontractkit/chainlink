package codec_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
)

func TestPropertyExtractor(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		A bool
		B int64
	}

	type nestedTestStruct struct {
		A string
		B testStruct
	}

	onChainType := reflect.TypeOf(nestedTestStruct{})

	extractor := codec.NewPropertyExtractor("A")
	invalidExtractor := codec.NewPropertyExtractor("A.B")
	nestedExtractor := codec.NewPropertyExtractor("B.B")

	t.Run("RetypeToOffChain sets the type for offchain to the onchain property", func(t *testing.T) {
		offChainType, err := extractor.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NoError(t, err)
		require.Equal(t, reflect.TypeOf(""), offChainType)
	})

	t.Run("RetypeToOffChain works on pointers", func(t *testing.T) {
		offChainType, err := extractor.RetypeToOffChain(reflect.TypeOf(&nestedTestStruct{}), "")
		require.NoError(t, err)

		str := ""

		require.Equal(t, reflect.TypeOf(&str), offChainType)
	})

	t.Run("RetypeToOffChain works on slices", func(t *testing.T) {
		offChainType, err := extractor.RetypeToOffChain(reflect.TypeOf([]nestedTestStruct{}), "")
		require.NoError(t, err)

		require.Equal(t, reflect.TypeOf([]string{}), offChainType)
	})

	t.Run("RetypeToOffChain works on arrays", func(t *testing.T) {
		arrayLen := 3
		offChainType, err := extractor.RetypeToOffChain(reflect.ArrayOf(arrayLen, onChainType), "")
		require.NoError(t, err)

		require.Equal(t, reflect.Array, offChainType.Kind())
		require.Equal(t, arrayLen, offChainType.Len())
	})

	t.Run("RetypeToOffChain returns error for missing field", func(t *testing.T) {
		_, err := invalidExtractor.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NotNil(t, err)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on structs", func(t *testing.T) {
		_, err := extractor.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NoError(t, err)

		onChainValue := nestedTestStruct{
			A: "test",
			B: testStruct{
				A: true,
				B: 42,
			},
		}

		offChainValue, err := extractor.TransformToOffChain(onChainValue, "")
		require.NoError(t, err)
		require.Equal(t, "test", offChainValue)

		lossyOnChain, err := extractor.TransformToOnChain(offChainValue, "")
		require.NoError(t, err)

		expected := nestedTestStruct{
			A: "test",
			B: testStruct{
				A: false,
				B: 0,
			},
		}

		assert.Equal(t, expected, lossyOnChain)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on nested structs", func(t *testing.T) {
		_, err := nestedExtractor.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NoError(t, err)

		onChainValue := nestedTestStruct{
			A: "test",
			B: testStruct{
				A: true,
				B: 42,
			},
		}

		offChainValue, err := nestedExtractor.TransformToOffChain(onChainValue, "")
		require.NoError(t, err)
		require.Equal(t, int64(42), offChainValue)

		lossyOnChain, err := nestedExtractor.TransformToOnChain(offChainValue, "")
		require.NoError(t, err)

		expected := nestedTestStruct{
			A: "",
			B: testStruct{
				A: false,
				B: 42,
			},
		}

		assert.Equal(t, expected, lossyOnChain)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on pointers", func(t *testing.T) {
		_, err := nestedExtractor.RetypeToOffChain(reflect.PointerTo(onChainType), "")
		require.NoError(t, err)

		onChainValue := &nestedTestStruct{
			A: "test",
			B: testStruct{
				A: true,
				B: 42,
			},
		}

		offChainValue, err := nestedExtractor.TransformToOffChain(onChainValue, "")
		require.NoError(t, err)

		expectedVal := int64(42)
		require.Equal(t, &expectedVal, offChainValue)

		lossyOnChain, err := nestedExtractor.TransformToOnChain(offChainValue, "")
		require.NoError(t, err)

		expected := &nestedTestStruct{
			A: "",
			B: testStruct{
				A: false,
				B: 42,
			},
		}

		assert.Equal(t, expected, lossyOnChain)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on slices", func(t *testing.T) {
		_, err := nestedExtractor.RetypeToOffChain(reflect.TypeOf([]nestedTestStruct{}), "")
		require.NoError(t, err)

		input := []nestedTestStruct{
			{A: "test0", B: testStruct{A: true, B: 42}},
			{A: "test1", B: testStruct{A: true, B: 43}},
		}

		expected := []int64{42, 43}

		offChainValue, err := nestedExtractor.TransformToOffChain(input, "")
		require.NoError(t, err)
		assert.Equal(t, expected, offChainValue)

		lossyOnChain, err := nestedExtractor.TransformToOnChain(offChainValue, "")
		require.NoError(t, err)

		expectedLossy := []nestedTestStruct{
			{A: "", B: testStruct{A: false, B: 42}},
			{A: "", B: testStruct{A: false, B: 43}},
		}

		assert.Equal(t, expectedLossy, lossyOnChain)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on slices of slices", func(t *testing.T) {
		_, err := nestedExtractor.RetypeToOffChain(reflect.TypeOf([][]nestedTestStruct{}), "")
		require.NoError(t, err)

		input := [][]nestedTestStruct{
			{
				{A: "test00", B: testStruct{A: true, B: 42}},
				{A: "test01", B: testStruct{A: true, B: 43}},
			},
			{
				{A: "test10", B: testStruct{A: true, B: 44}},
				{A: "test11", B: testStruct{A: true, B: 45}},
			},
		}

		expected := [][]int64{{42, 43}, {44, 45}}

		offChainValue, err := nestedExtractor.TransformToOffChain(input, "")
		require.NoError(t, err)
		assert.Equal(t, expected, offChainValue)

		lossyOnChain, err := nestedExtractor.TransformToOnChain(offChainValue, "")
		require.NoError(t, err)

		expectedLossy := [][]nestedTestStruct{
			{
				{A: "", B: testStruct{A: false, B: 42}},
				{A: "", B: testStruct{A: false, B: 43}},
			},
			{
				{A: "", B: testStruct{A: false, B: 44}},
				{A: "", B: testStruct{A: false, B: 45}},
			},
		}

		assert.Equal(t, expectedLossy, lossyOnChain)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on arrays", func(t *testing.T) {
		_, err := nestedExtractor.RetypeToOffChain(reflect.ArrayOf(2, onChainType), "")
		require.NoError(t, err)

		input := [2]nestedTestStruct{
			{A: "test0", B: testStruct{A: true, B: 42}},
			{A: "test1", B: testStruct{A: true, B: 43}},
		}

		expected := [2]int64{42, 43}

		offChainValue, err := nestedExtractor.TransformToOffChain(input, "")
		require.NoError(t, err)
		assert.Equal(t, expected, offChainValue)

		lossyOnChain, err := nestedExtractor.TransformToOnChain(offChainValue, "")
		require.NoError(t, err)

		expectedLossy := [2]nestedTestStruct{
			{A: "", B: testStruct{A: false, B: 42}},
			{A: "", B: testStruct{A: false, B: 43}},
		}

		assert.Equal(t, expectedLossy, lossyOnChain)
	})
}
