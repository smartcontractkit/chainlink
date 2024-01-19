package codec_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestRenamer(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		A string
		B int64
		C int64
	}

	type nestedTestStruct struct {
		A string
		B testStruct
		C []testStruct
		D string
	}

	renamer := codec.NewRenamer(map[string]string{"A": "X", "C": "Z"})
	invalidRenamer := codec.NewRenamer(map[string]string{"W": "X", "C": "Z"})
	nestedRenamer := codec.NewRenamer(map[string]string{"A": "X", "B.A": "X", "B.C": "Z", "C.A": "X", "C.C": "Z", "B": "Y"})
	t.Run("RetypeToOffChain renames fields keeping structure", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		require.NoError(t, err)

		assertBasicRenameTransform(t, offChainType)
	})

	t.Run("RetypeToOffChain works on slices", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf([]testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Slice, offChainType.Kind())
		assertBasicRenameTransform(t, offChainType.Elem())
	})

	t.Run("RetypeToOffChain works on pointers", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf(&testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Pointer, offChainType.Kind())
		assertBasicRenameTransform(t, offChainType.Elem())
	})

	t.Run("RetypeToOffChain works on pointers to non structs", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf(&[]testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Pointer, offChainType.Kind())
		assert.Equal(t, reflect.Slice, offChainType.Elem().Kind())
		assertBasicRenameTransform(t, offChainType.Elem().Elem())
	})

	t.Run("RetypeToOffChain works on arrays", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf([2]testStruct{}), "")
		require.NoError(t, err)

		assert.Equal(t, reflect.Array, offChainType.Kind())
		assert.Equal(t, 2, offChainType.Len())
		assertBasicRenameTransform(t, offChainType.Elem())
	})

	t.Run("RetypeToOffChain returns exception if a field is not on the type", func(t *testing.T) {
		_, err := invalidRenamer.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("RetypeToOffChain works on nested fields even if the field itself is renamed", func(t *testing.T) {
		offChainType, err := nestedRenamer.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NoError(t, err)
		assert.Equal(t, 4, offChainType.NumField())
		f0 := offChainType.Field(0)
		assert.Equal(t, "X", f0.Name)
		assert.Equal(t, reflect.TypeOf(""), f0.Type)
		f1 := offChainType.Field(1)
		assert.Equal(t, "Y", f1.Name)
		assertBasicRenameTransform(t, f1.Type)
		f2 := offChainType.Field(2)
		assert.Equal(t, "C", f2.Name)
		assert.Equal(t, reflect.Slice, f2.Type.Kind())
		assertBasicRenameTransform(t, f2.Type.Elem())
		f3 := offChainType.Field(3)
		assert.Equal(t, "D", f3.Name)
		assert.Equal(t, reflect.TypeOf(""), f3.Type)
	})

	t.Run("RetypeToOffChain returns an error if the name is already in use", func(t *testing.T) {
		dup := codec.NewRenamer(map[string]string{"A": "B"})
		_, err := dup.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		require.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("TransformToOnChain and TransformToOffChain works on structs", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf(testStruct{}), "")
		require.NoError(t, err)
		iOffchain := reflect.Indirect(reflect.New(offChainType))
		iOffchain.FieldByName("X").SetString("foo")
		iOffchain.FieldByName("B").SetInt(10)
		iOffchain.FieldByName("Z").SetInt(20)

		output, err := renamer.TransformToOnChain(iOffchain.Interface(), "")

		require.NoError(t, err)

		expected := testStruct{
			A: "foo",
			B: 10,
			C: 20,
		}
		assert.Equal(t, expected, output)
		newInput, err := renamer.TransformToOffChain(expected, "")
		require.NoError(t, err)
		assert.Equal(t, iOffchain.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain returns error if input type was not from TransformToOnChain", func(t *testing.T) {
		_, err := invalidRenamer.TransformToOnChain(testStruct{}, "")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("TransformToOnChain and TransformToOffChain works on pointers", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf(&testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.New(offChainType.Elem())
		iOffchain := reflect.Indirect(rInput)
		iOffchain.FieldByName("X").SetString("foo")
		iOffchain.FieldByName("B").SetInt(10)
		iOffchain.FieldByName("Z").SetInt(20)

		output, err := renamer.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := &testStruct{
			A: "foo",
			B: 10,
			C: 20,
		}
		assert.Equal(t, expected, output)

		// Optimization to avoid creating objects unnecessarily
		iOffchain.FieldByName("X").SetString("Z")
		expected.A = "Z"
		assert.Equal(t, expected, output)
		newInput, err := renamer.TransformToOffChain(output, "")
		require.NoError(t, err)
		assert.Same(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on slices", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf([]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.MakeSlice(offChainType, 2, 2)
		iOffchain := rInput.Index(0)
		iOffchain.FieldByName("X").SetString("foo")
		iOffchain.FieldByName("B").SetInt(10)
		iOffchain.FieldByName("Z").SetInt(20)
		iOffchain = rInput.Index(1)
		iOffchain.FieldByName("X").SetString("baz")
		iOffchain.FieldByName("B").SetInt(15)
		iOffchain.FieldByName("Z").SetInt(25)

		output, err := renamer.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := []testStruct{
			{
				A: "foo",
				B: 10,
				C: 20,
			},
			{
				A: "baz",
				B: 15,
				C: 25,
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := renamer.TransformToOffChain(expected, "")
		require.NoError(t, err)
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on nested slices", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf([][]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.MakeSlice(offChainType, 2, 2)
		rOuter := rInput.Index(0)
		rOuter.Set(reflect.MakeSlice(rOuter.Type(), 2, 2))
		iOffchain := rOuter.Index(0)
		iOffchain.FieldByName("X").SetString("foo")
		iOffchain.FieldByName("B").SetInt(10)
		iOffchain.FieldByName("Z").SetInt(20)
		iOffchain = rOuter.Index(1)
		iOffchain.FieldByName("X").SetString("baz")
		iOffchain.FieldByName("B").SetInt(15)
		iOffchain.FieldByName("Z").SetInt(25)
		rOuter = rInput.Index(1)
		rOuter.Set(reflect.MakeSlice(rOuter.Type(), 2, 2))
		iOffchain = rOuter.Index(0)
		iOffchain.FieldByName("X").SetString("fooz")
		iOffchain.FieldByName("B").SetInt(100)
		iOffchain.FieldByName("Z").SetInt(200)
		iOffchain = rOuter.Index(1)
		iOffchain.FieldByName("X").SetString("bazz")
		iOffchain.FieldByName("B").SetInt(150)
		iOffchain.FieldByName("Z").SetInt(250)

		output, err := renamer.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := [][]testStruct{
			{
				{
					A: "foo",
					B: 10,
					C: 20,
				},
				{
					A: "baz",
					B: 15,
					C: 25,
				},
			},
			{
				{
					A: "fooz",
					B: 100,
					C: 200,
				},
				{
					A: "bazz",
					B: 150,
					C: 250,
				},
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := renamer.TransformToOffChain(expected, "")
		require.NoError(t, err)
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on pointers to non structs", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf(&[]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.New(offChainType.Elem())
		rElm := reflect.MakeSlice(offChainType.Elem(), 2, 2)
		iElm := rElm.Index(0)
		iElm.FieldByName("X").SetString("foo")
		iElm.FieldByName("B").SetInt(10)
		iElm.FieldByName("Z").SetInt(20)
		iElm = rElm.Index(1)
		iElm.FieldByName("X").SetString("baz")
		iElm.FieldByName("B").SetInt(15)
		iElm.FieldByName("Z").SetInt(25)
		reflect.Indirect(rInput).Set(rElm)

		output, err := renamer.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := &[]testStruct{
			{
				A: "foo",
				B: 10,
				C: 20,
			},
			{
				A: "baz",
				B: 15,
				C: 25,
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := renamer.TransformToOffChain(expected, "")
		require.NoError(t, err)
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on arrays", func(t *testing.T) {
		offChainType, err := renamer.RetypeToOffChain(reflect.TypeOf([2]testStruct{}), "")
		require.NoError(t, err)
		rInput := reflect.New(offChainType).Elem()
		iOffchain := rInput.Index(0)
		iOffchain.FieldByName("X").SetString("foo")
		iOffchain.FieldByName("B").SetInt(10)
		iOffchain.FieldByName("Z").SetInt(20)
		iOffchain = rInput.Index(1)
		iOffchain.FieldByName("X").SetString("baz")
		iOffchain.FieldByName("B").SetInt(15)
		iOffchain.FieldByName("Z").SetInt(25)

		output, err := renamer.TransformToOnChain(rInput.Interface(), "")

		require.NoError(t, err)

		expected := [2]testStruct{
			{
				A: "foo",
				B: 10,
				C: 20,
			},
			{
				A: "baz",
				B: 15,
				C: 25,
			},
		}
		assert.Equal(t, expected, output)

		newInput, err := renamer.TransformToOffChain(expected, "")
		require.NoError(t, err)
		assert.Equal(t, rInput.Interface(), newInput)
	})

	t.Run("TransformToOnChain and TransformToOffChain works on nested fields even if the field itself is renamed", func(t *testing.T) {
		offChainType, err := nestedRenamer.RetypeToOffChain(reflect.TypeOf(nestedTestStruct{}), "")
		require.NoError(t, err)
		iOffchain := reflect.Indirect(reflect.New(offChainType))

		iOffchain.FieldByName("X").SetString("foo")
		rY := iOffchain.FieldByName("Y")
		rY.FieldByName("X").SetString("foo")
		rY.FieldByName("B").SetInt(10)
		rY.FieldByName("Z").SetInt(20)

		rC := iOffchain.FieldByName("C")
		rC.Set(reflect.MakeSlice(rC.Type(), 2, 2))
		iElm := rC.Index(0)
		iElm.FieldByName("X").SetString("foo")
		iElm.FieldByName("B").SetInt(10)
		iElm.FieldByName("Z").SetInt(20)
		iElm = rC.Index(1)
		iElm.FieldByName("X").SetString("baz")
		iElm.FieldByName("B").SetInt(15)
		iElm.FieldByName("Z").SetInt(25)

		iOffchain.FieldByName("D").SetString("bar")

		output, err := nestedRenamer.TransformToOnChain(iOffchain.Interface(), "")

		require.NoError(t, err)

		expected := nestedTestStruct{
			A: "foo",
			B: testStruct{
				A: "foo",
				B: 10,
				C: 20,
			},
			C: []testStruct{
				{
					A: "foo",
					B: 10,
					C: 20,
				},
				{
					A: "baz",
					B: 15,
					C: 25,
				},
			},
			D: "bar",
		}
		assert.Equal(t, expected, output)
		newInput, err := nestedRenamer.TransformToOffChain(expected, "")
		require.NoError(t, err)
		assert.Equal(t, iOffchain.Interface(), newInput)
	})
}

func assertBasicRenameTransform(t *testing.T, offChainType reflect.Type) {
	require.Equal(t, 3, offChainType.NumField())
	f0 := offChainType.Field(0)
	assert.Equal(t, "X", f0.Name)
	assert.Equal(t, reflect.TypeOf(""), f0.Type)
	f1 := offChainType.Field(1)
	assert.Equal(t, "B", f1.Name)
	assert.Equal(t, reflect.TypeOf(int64(0)), f1.Type)
	f2 := offChainType.Field(2)
	assert.Equal(t, "Z", f2.Name)
	assert.Equal(t, reflect.TypeOf(int64(0)), f2.Type)
}
