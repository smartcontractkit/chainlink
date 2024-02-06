package codec_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestElementExtractor(t *testing.T) {
	first := codec.ElementExtractorLocationFirst
	middle := codec.ElementExtractorLocationMiddle
	last := codec.ElementExtractorLocationLast

	type testStruct struct {
		A string
		B int64
		C int64
		D uint64
	}

	type nestedTestStruct struct {
		A string
		B testStruct
		C []testStruct
		D string
	}

	extractor := codec.NewElementExtractor(map[string]*codec.ElementExtractorLocation{"A": &first, "C": &middle, "D": &last})
	invalidExtractor := codec.NewElementExtractor(map[string]*codec.ElementExtractorLocation{"A": &first, "W": &middle})
	nestedExtractor := codec.NewElementExtractor(map[string]*codec.ElementExtractorLocation{"A": &first, "B.A": &first, "B.C": &middle, "B.D": &last, "C.A": &first, "C.C": &middle, "C.D": &last, "B": &last})
	t.Run("RetypeToOffChain gets non-slice type", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		require.NoError(t, err)

		assertBasicElementExtractTransform(t, inputType)
	})

	t.Run("RetypeToOffChain works on slices", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf([]testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Slice, inputType.Kind())
		assertBasicElementExtractTransform(t, inputType.Elem())
	})

	t.Run("RetypeToOffChain works on pointers", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf(&testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Pointer, inputType.Kind())
		assertBasicElementExtractTransform(t, inputType.Elem())
	})

	t.Run("RetypeToOffChain works on pointers to non structs", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf(&[]testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Pointer, inputType.Kind())
		assert.Equal(t, reflect.Slice, inputType.Elem().Kind())
		assertBasicElementExtractTransform(t, inputType.Elem().Elem())
	})

	t.Run("RetypeToOffChain works on arrays", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf([2]testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Array, inputType.Kind())
		assert.Equal(t, 2, inputType.Len())
		assertBasicElementExtractTransform(t, inputType.Elem())
	})

	t.Run("RetypeToOffChain returns exception if a field is not on the type", func(t *testing.T) {
		_, err := invalidExtractor.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("RetypeToOffChain works on nested fields even if the field itself is also extracted", func(t *testing.T) {
		inputType, err := nestedExtractor.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		fmt.Printf("%+v\n", inputType)
		require.NoError(t, err)
		assert.Equal(t, 4, inputType.NumField())
		f0 := inputType.Field(0)
		assert.Equal(t, "A", f0.Name)
		assert.Equal(t, reflect.TypeOf([]string{}), f0.Type)
		f1 := inputType.Field(1)
		assert.Equal(t, "B", f1.Name)
		require.Equal(t, reflect.Slice, f1.Type.Kind())
		assertBasicElementExtractTransform(t, f1.Type.Elem())
		f2 := inputType.Field(2)
		require.Equal(t, reflect.Slice, f2.Type.Kind())
		assertBasicElementExtractTransform(t, f2.Type.Elem())
		f3 := inputType.Field(3)
		assert.Equal(t, "D", f3.Name)
		assert.Equal(t, reflect.TypeOf(""), f3.Type)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on structs", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		require.NoError(t, err)
		iInput := reflect.Indirect(reflect.New(inputType))
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iInput.FieldByName("B").SetInt(10)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))

		output, err := extractor.TransformToOnChain(iInput.Interface(), "")

		require.NoError(t, err)

		expected := testStruct{
			A: "A",
			B: 10,
			C: 20,
			D: 30,
		}
		assert.Equal(t, expected, output)
		newInput, err := extractor.TransformToOffChain(expected, "")
		require.NoError(t, err)
		// Lossy modification
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		assert.Equal(t, iInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain returns error if input type was not from TransformToOnChain", func(t *testing.T) {
		_, err := invalidExtractor.TransformToOnChain(testStruct{}, "")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
		_, err = invalidExtractor.TransformToOffChain(testStruct{}, "")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("TransformToOnChain and TransformToOffChain works on pointers", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf(&testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.New(inputType.Elem())
		iInput := reflect.Indirect(rInput)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iInput.FieldByName("B").SetInt(10)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))

		output, err := extractor.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := &testStruct{
			A: "A",
			B: 10,
			C: 20,
			D: 30,
		}
		assert.Equal(t, expected, output)
		newInput, err := extractor.TransformToOffChain(expected, "")
		require.NoError(t, err)
		// Lossy modification
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on slices", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf([]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.MakeSlice(inputType, 2, 2)
		iInput := rInput.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iInput.FieldByName("B").SetInt(10)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))
		iInput = rInput.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az", "Bz", "Cz"}))
		iInput.FieldByName("B").SetInt(15)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 25, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 35}))

		output, err := extractor.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := []testStruct{
			{
				A: "A",
				B: 10,
				C: 20,
				D: 30,
			},
			{
				A: "Az",
				B: 15,
				C: 25,
				D: 35,
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := extractor.TransformToOffChain(expected, "")
		require.NoError(t, err)
		// Lossy modification
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{25}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{35}))
		iInput = rInput.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on nested slices", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf([][]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.MakeSlice(inputType, 2, 2)
		rOuter := rInput.Index(0)
		rOuter.Set(reflect.MakeSlice(rOuter.Type(), 2, 2))
		iInput := rOuter.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iInput.FieldByName("B").SetInt(10)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))
		iInput = rOuter.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az", "Bz", "Cz"}))
		iInput.FieldByName("B").SetInt(15)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 25, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 35}))
		rOuter = rInput.Index(1)
		rOuter.Set(reflect.MakeSlice(rOuter.Type(), 2, 2))
		iInput = rOuter.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az", "Bz", "Cz"}))
		iInput.FieldByName("B").SetInt(100)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{150, 200, 350}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{100, 200, 300}))
		iInput = rOuter.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Azz", "Bzz", "Czz"}))
		iInput.FieldByName("B").SetInt(150)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{150, 250, 350}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{100, 200, 350}))

		output, err := extractor.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := [][]testStruct{
			{
				{
					A: "A",
					B: 10,
					C: 20,
					D: 30,
				},
				{
					A: "Az",
					B: 15,
					C: 25,
					D: 35,
				},
			},
			{
				{
					A: "Az",
					B: 100,
					C: 200,
					D: 300,
				},
				{
					A: "Azz",
					B: 150,
					C: 250,
					D: 350,
				},
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := extractor.TransformToOffChain(expected, "")
		require.NoError(t, err)
		// Lossy modification
		rOuter = rInput.Index(0)
		iInput = rOuter.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		iInput = rOuter.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{25}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{35}))
		rOuter = rInput.Index(1)
		iInput = rOuter.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{200}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{300}))
		iInput = rOuter.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Azz"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{250}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{350}))

		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on pointers to non structs", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf(&[]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.New(inputType.Elem())
		rElm := reflect.MakeSlice(inputType.Elem(), 2, 2)
		iInput := rElm.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iInput.FieldByName("B").SetInt(10)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))
		iInput = rElm.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az", "Bz", "Cz"}))
		iInput.FieldByName("B").SetInt(15)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 25, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 35}))
		reflect.Indirect(rInput).Set(rElm)

		output, err := extractor.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := &[]testStruct{
			{
				A: "A",
				B: 10,
				C: 20,
				D: 30,
			},
			{
				A: "Az",
				B: 15,
				C: 25,
				D: 35,
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := extractor.TransformToOffChain(expected, "")
		require.NoError(t, err)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{25}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{35}))
		iInput = rElm.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on arrays", func(t *testing.T) {
		inputType, err := extractor.RetypeToOffChain(reflect.TypeOf([2]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.New(inputType).Elem()
		iInput := rInput.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iInput.FieldByName("B").SetInt(10)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))
		iInput = rInput.Index(1)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az", "Bz", "Cz"}))
		iInput.FieldByName("B").SetInt(15)
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{15, 25, 35}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 35}))

		output, err := extractor.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := [2]testStruct{
			{
				A: "A",
				B: 10,
				C: 20,
				D: 30,
			},
			{
				A: "Az",
				B: 15,
				C: 25,
				D: 35,
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := extractor.TransformToOffChain(expected, "")
		require.NoError(t, err)
		// Lossy modification
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"Az"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{25}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{35}))
		iInput = rInput.Index(0)
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		iInput.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		iInput.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on nested fields even if the field itself is also extracted", func(t *testing.T) {
		inputType, err := nestedExtractor.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NoError(t, err)

		iInput := reflect.Indirect(reflect.New(inputType))
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))

		rB := iInput.FieldByName("B")
		rB.Set(reflect.MakeSlice(rB.Type(), 2, 2))

		rElm := rB.Index(0)
		rElm.FieldByName("A").Set(reflect.ValueOf([]string{"Z", "W", "Z"}))
		rElm.FieldByName("B").SetInt(99)
		rElm.FieldByName("C").Set(reflect.ValueOf([]int64{44, 44, 44}))
		rElm.FieldByName("D").Set(reflect.ValueOf([]uint64{42, 62, 99}))

		rElm = rB.Index(1)
		rElm.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		rElm.FieldByName("B").SetInt(10)
		rElm.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		rElm.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))

		rC := iInput.FieldByName("C")
		rC.Set(reflect.MakeSlice(rC.Type(), 2, 2))
		iElm := rC.Index(0)
		iElm.FieldByName("A").Set(reflect.ValueOf([]string{"A", "B", "C"}))
		iElm.FieldByName("B").SetInt(10)
		iElm.FieldByName("C").Set(reflect.ValueOf([]int64{15, 20, 35}))
		iElm.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 30}))
		iElm = rC.Index(1)
		iElm.FieldByName("A").Set(reflect.ValueOf([]string{"Az", "Bz", "Cz"}))
		iElm.FieldByName("B").SetInt(15)
		iElm.FieldByName("C").Set(reflect.ValueOf([]int64{15, 25, 35}))
		iElm.FieldByName("D").Set(reflect.ValueOf([]uint64{10, 20, 35}))

		iInput.FieldByName("D").SetString("bar")

		output, err := nestedExtractor.TransformToOnChain(iInput.Interface(), "")
		require.NoError(t, err)

		expected := nestedTestStruct{
			A: "A",
			B: testStruct{
				A: "A",
				B: 10,
				C: 20,
				D: 30,
			},
			C: []testStruct{
				{
					A: "A",
					B: 10,
					C: 20,
					D: 30,
				},
				{
					A: "Az",
					B: 15,
					C: 25,
					D: 35,
				},
			},
			D: "bar",
		}

		assert.Equal(t, expected, output)

		newInput, err := nestedExtractor.TransformToOffChain(expected, "")
		require.NoError(t, err)

		// Lossy modification
		iInput.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		rB.Set(rB.Slice(1, 2))
		rElm = rB.Index(0)
		rElm.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		rElm.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		rElm.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))

		rElm = rC.Index(0)
		rElm.FieldByName("A").Set(reflect.ValueOf([]string{"A"}))
		rElm.FieldByName("C").Set(reflect.ValueOf([]int64{20}))
		rElm.FieldByName("D").Set(reflect.ValueOf([]uint64{30}))
		rElm = rC.Index(1)
		rElm.FieldByName("A").Set(reflect.ValueOf([]string{"Az"}))
		rElm.FieldByName("C").Set(reflect.ValueOf([]int64{25}))
		rElm.FieldByName("D").Set(reflect.ValueOf([]uint64{35}))
		assert.Equal(t, iInput.Interface(), newInput)
	})

	for _, test := range []struct {
		location codec.ElementExtractorLocation
	}{
		{location: codec.ElementExtractorLocationFirst},
		{location: codec.ElementExtractorLocationMiddle},
		{location: codec.ElementExtractorLocationLast},
	} {
		t.Run("Json encoding works", func(t *testing.T) {
			b, err := json.Marshal(test.location)
			require.NoError(t, err)
			var actual codec.ElementExtractorLocation
			require.NoError(t, json.Unmarshal(b, &actual))
			assert.Equal(t, test.location, actual)
		})
	}
}

func assertBasicElementExtractTransform(t *testing.T, inputType reflect.Type) {
	require.Equal(t, 4, inputType.NumField())
	f0 := inputType.Field(0)
	assert.Equal(t, "A", f0.Name)
	assert.Equal(t, reflect.TypeOf([]string{}), f0.Type)
	f1 := inputType.Field(1)
	assert.Equal(t, "B", f1.Name)
	assert.Equal(t, reflect.TypeOf(int64(0)), f1.Type)
	f2 := inputType.Field(2)
	assert.Equal(t, "C", f2.Name)
	assert.Equal(t, reflect.TypeOf([]int64{}), f2.Type)
	f3 := inputType.Field(3)
	assert.Equal(t, "D", f3.Name)
	assert.Equal(t, reflect.TypeOf([]uint64{}), f3.Type)
}
