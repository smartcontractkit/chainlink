package types

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestCodecEntry(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		type1, err := abi.NewType("uint16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		type2, err := abi.NewType("string", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		type3, err := abi.NewType("uint24", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		type4, err := abi.NewType("int24", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		args := abi.Arguments{
			{Name: "Field1", Type: type1},
			{Name: "Field2", Type: type2},
			{Name: "Field3", Type: type3},
			{Name: "Field4", Type: type4},
		}
		entry := NewCodecEntry(args, nil, nil)
		require.NoError(t, entry.Init())
		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		f1 := uint16(2)
		iChecked.FieldByName("Field1").Set(reflect.ValueOf(&f1))
		f2 := "any string"
		iChecked.FieldByName("Field2").Set(reflect.ValueOf(&f2))

		f3 := big.NewInt( /*2^24 - 1*/ 16777215)
		setAndVerifyLimit(t, (*uint24)(f3), f3, iChecked.FieldByName("Field3"))

		f4 := big.NewInt( /*2^23 - 1*/ 8388607)
		setAndVerifyLimit(t, (*int24)(f4), f4, iChecked.FieldByName("Field4"))

		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		assert.Equal(t, native.Field(0).Interface(), iChecked.Field(0).Interface())
		assert.Equal(t, native.Field(1).Interface(), iChecked.Field(1).Interface())
		assert.Equal(t, native.Field(2).Interface(), f3)
		assert.Equal(t, native.Field(3).Interface(), f4)
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("tuples", func(t *testing.T) {
		type1, err := abi.NewType("uint16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		tupleType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "Field3", Type: "uint24"},
			{Name: "Field4", Type: "int24"},
		})
		require.NoError(t, err)
		args := abi.Arguments{
			{Name: "Field1", Type: type1},
			{Name: "Field2", Type: tupleType},
		}
		entry := NewCodecEntry(args, nil, nil)
		require.NoError(t, entry.Init())

		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		f1 := uint16(2)
		iChecked.FieldByName("Field1").Set(reflect.ValueOf(&f1))
		f2 := iChecked.FieldByName("Field2")
		f2.Set(reflect.New(f2.Type().Elem()))
		f2 = reflect.Indirect(f2)
		f3 := big.NewInt( /*2^24 - 1*/ 16777215)
		setAndVerifyLimit(t, (*uint24)(f3), f3, f2.FieldByName("Field3"))
		f4 := big.NewInt( /*2^23 - 1*/ 8388607)
		setAndVerifyLimit(t, (*int24)(f4), f4, f2.FieldByName("Field4"))

		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		require.Equal(t, native.Field(0).Interface(), iChecked.Field(0).Interface())
		nF2 := reflect.Indirect(native.Field(1))
		assert.Equal(t, nF2.Field(0).Interface(), f3)
		assert.Equal(t, nF2.Field(1).Interface(), f4)
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("nested tuple member names are capitalized", func(t *testing.T) {
		type1, err := abi.NewType("uint16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		tupleType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "field3", Type: "uint24"},
			{Name: "field4", Type: "int24"},
		})
		require.NoError(t, err)
		args := abi.Arguments{
			{Name: "field1", Type: type1},
			{Name: "field2", Type: tupleType},
		}
		entry := NewCodecEntry(args, nil, nil)
		require.NoError(t, entry.Init())

		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		f1 := uint16(2)
		iChecked.FieldByName("Field1").Set(reflect.ValueOf(&f1))
		f2 := iChecked.FieldByName("Field2")
		f2.Set(reflect.New(f2.Type().Elem()))
		f2 = reflect.Indirect(f2)
		f3 := big.NewInt( /*2^24 - 1*/ 16777215)
		setAndVerifyLimit(t, (*uint24)(f3), f3, f2.FieldByName("Field3"))
		f4 := big.NewInt( /*2^23 - 1*/ 8388607)
		setAndVerifyLimit(t, (*int24)(f4), f4, f2.FieldByName("Field4"))

		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		require.Equal(t, native.Field(0).Interface(), iChecked.Field(0).Interface())
		nF2 := reflect.Indirect(native.Field(1))
		assert.Equal(t, nF2.Field(0).Interface(), f3)
		assert.Equal(t, nF2.Field(1).Interface(), f4)
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("unwrapped types", func(t *testing.T) {
		// This exists to allow you to decode single returned values without naming the parameter
		wrappedTuple, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "Field1", Type: "int16"},
		})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "", Type: wrappedTuple}}, nil, nil)
		require.NoError(t, entry.Init())
		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		anyValue := int16(2)
		iChecked.FieldByName("Field1").Set(reflect.ValueOf(&anyValue))
		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		assert.Equal(t, &anyValue, native.FieldByName("Field1").Interface())
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("slice types", func(t *testing.T) {
		type1, err := abi.NewType("int16[]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "Field1", Type: type1}}, nil, nil)

		require.NoError(t, entry.Init())
		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		anySliceValue := &[]int16{2, 3}
		iChecked.FieldByName("Field1").Set(reflect.ValueOf(anySliceValue))
		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		assert.Equal(t, anySliceValue, native.FieldByName("Field1").Interface())
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("array types", func(t *testing.T) {
		type1, err := abi.NewType("int16[3]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "Field1", Type: type1}}, nil, nil)
		require.NoError(t, entry.Init())
		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		anySliceValue := &[3]int16{2, 3, 30}
		iChecked.FieldByName("Field1").Set(reflect.ValueOf(anySliceValue))
		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		assert.Equal(t, anySliceValue, native.FieldByName("Field1").Interface())
	})

	t.Run("Not return values makes struct{}", func(t *testing.T) {
		entry := NewCodecEntry(abi.Arguments{}, nil, nil)
		require.NoError(t, entry.Init())
		assert.Equal(t, reflect.TypeOf(struct{}{}), entry.CheckedType())
		native, err := entry.ToNative(reflect.ValueOf(&struct{}{}))
		require.NoError(t, err)
		assert.Equal(t, struct{}{}, native.Interface())
	})

	t.Run("Address works", func(t *testing.T) {
		address, err := abi.NewType("address", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "foo", Type: address}}, nil, nil)
		require.NoError(t, entry.Init())

		checked := reflect.New(entry.CheckedType())
		iChecked := reflect.Indirect(checked)
		anyAddr := &common.Address{1, 2, 3}
		iChecked.FieldByName("Foo").Set(reflect.ValueOf(anyAddr))

		native, err := entry.ToNative(checked)
		require.NoError(t, err)
		assert.Equal(t, anyAddr, native.FieldByName("Foo").Interface())
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("Unnamed parameters are named after their locations", func(t *testing.T) {
		// use different types to make sure that the right fields have the right types.
		anyType1, err := abi.NewType("int64", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		anyType2, err := abi.NewType("int32", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "", Type: anyType1}, {Name: "", Type: anyType2}}, nil, nil)
		assert.NoError(t, entry.Init())
		ct := entry.CheckedType()
		require.Equal(t, 2, ct.NumField())
		args := entry.Args()
		require.Equal(t, 2, len(args))
		f0 := ct.Field(0)
		assert.Equal(t, "F0", f0.Name)
		assert.Equal(t, "F0", args[0].Name)
		assert.Equal(t, reflect.TypeOf((*int64)(nil)), f0.Type)
		f1 := ct.Field(1)
		assert.Equal(t, "F1", f1.Name)
		assert.Equal(t, "F1", args[1].Name)
		assert.Equal(t, reflect.TypeOf((*int32)(nil)), f1.Type)
	})

	t.Run("Unnamed parameters adds _Xes at the end if their location name is taken", func(t *testing.T) {
		// use different types to make sure that the right fields have the right types.
		anyType1, err := abi.NewType("int64", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		anyType2, err := abi.NewType("int32", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "F1", Type: anyType1}, {Name: "", Type: anyType2}}, nil, nil)
		assert.NoError(t, entry.Init())
		ct := entry.CheckedType()
		require.Equal(t, 2, ct.NumField())
		args := entry.Args()
		require.Equal(t, 2, len(args))
		f0 := ct.Field(0)
		assert.Equal(t, "F1", f0.Name)
		assert.Equal(t, reflect.TypeOf((*int64)(nil)), f0.Type)
		f1 := ct.Field(1)
		assert.Equal(t, "F1_X", f1.Name)
		assert.Equal(t, "F1_X", args[1].Name)
		assert.Equal(t, reflect.TypeOf((*int32)(nil)), f1.Type)
	})

	t.Run("Multiple abi arguments with the same name returns an error", func(t *testing.T) {
		anyType, err := abi.NewType("int16[3]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "Name", Type: anyType}, {Name: "Name", Type: anyType}}, nil, nil)
		assert.True(t, errors.Is(entry.Init(), commontypes.ErrInvalidConfig))
	})

	t.Run("Indexed basic types leave their native and checked types as-is", func(t *testing.T) {
		anyType, err := abi.NewType("int16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "Name", Type: anyType, Indexed: true}}, nil, nil)
		require.NoError(t, entry.Init())
		checkedField, ok := entry.CheckedType().FieldByName("Name")
		require.True(t, ok)
		assert.Equal(t, reflect.TypeOf((*int16)(nil)), checkedField.Type)
		native, err := entry.ToNative(reflect.New(entry.CheckedType()))
		require.NoError(t, err)
		assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
	})

	t.Run("Indexed string and bytes array change to hash", func(t *testing.T) {
		stringType, err := abi.NewType("string", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		arrayType, err := abi.NewType("uint8[32]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)

		abiArgs := abi.Arguments{
			{Name: "String", Type: stringType, Indexed: true},
			{Name: "Array", Type: arrayType, Indexed: true},
		}

		for i := 0; i < len(abiArgs); i++ {
			entry := NewCodecEntry(abi.Arguments{abiArgs[i]}, nil, nil)
			require.NoError(t, entry.Init())
			nativeField, ok := entry.CheckedType().FieldByName(abiArgs[i].Name)
			require.True(t, ok)
			assert.Equal(t, reflect.TypeOf(&common.Hash{}), nativeField.Type)
			native, err := entry.ToNative(reflect.New(entry.CheckedType()))
			require.NoError(t, err)
			assertHaveSameStructureAndNames(t, native.Type(), entry.CheckedType())
		}
	})

	t.Run("Too many indexed items returns an error", func(t *testing.T) {
		anyType, err := abi.NewType("int16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(
			abi.Arguments{
				{Name: "Name1", Type: anyType, Indexed: true},
				{Name: "Name2", Type: anyType, Indexed: true},
				{Name: "Name3", Type: anyType, Indexed: true},
				{Name: "Name4", Type: anyType, Indexed: true},
			}, nil, nil)
		require.True(t, errors.Is(entry.Init(), commontypes.ErrInvalidConfig))
	})

	// TODO: when the TODO on
	// https://github.com/ethereum/go-ethereum/blob/release/1.12/accounts/abi/topics.go#L78
	// is removed, remove this test.
	t.Run("Using unsupported types by go-ethereum returns an error", func(t *testing.T) {
		anyType, err := abi.NewType("int256[2]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := NewCodecEntry(abi.Arguments{{Name: "Name", Type: anyType, Indexed: true}}, nil, nil)
		assert.True(t, errors.Is(entry.Init(), commontypes.ErrInvalidConfig))
	})

	t.Run("Modifier returns provided modifier", func(t *testing.T) {
		anyType, err := abi.NewType("int16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		mod := codec.NewRenamer(map[string]string{"Name": "RenamedName"})
		entry := NewCodecEntry(abi.Arguments{{Name: "Name", Type: anyType, Indexed: true}}, nil, mod)
		assert.Equal(t, mod, entry.Modifier())
	})

	t.Run("EncodingPrefix returns provided prefix", func(t *testing.T) {
		anyType, err := abi.NewType("int16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		prefix := []byte{1, 2, 3}
		entry := NewCodecEntry(abi.Arguments{{Name: "Name", Type: anyType, Indexed: true}}, prefix, nil)
		assert.Equal(t, prefix, entry.EncodingPrefix())
	})
}

// sized and bi must be the same pointer.
func setAndVerifyLimit(t *testing.T, sbi SizedBigInt, bi *big.Int, field reflect.Value) {
	require.Same(t, reflect.NewAt(reflect.TypeOf(big.Int{}), reflect.ValueOf(sbi).UnsafePointer()).Interface(), bi)
	field.Set(reflect.ValueOf(sbi))
	assert.NoError(t, sbi.Verify())
	bi.Add(bi, big.NewInt(1))
	assert.IsType(t, commontypes.ErrInvalidType, sbi.Verify())
}

// verifying the same structure allows us to use unsafe pointers to cast between them.
// This is done for perf and simplicity in mapping the two structures.
// [reflect.NewAt]'s use is the same as (*native)(unsafe.Pointer(checked))
// See the safe usecase 1 from [unsafe.Pointer], as this is a subset of that.
// This also verifies field names are the same.
func assertHaveSameStructureAndNames(t *testing.T, t1, t2 reflect.Type) {
	require.Equal(t, t1.Kind(), t2.Kind())

	switch t1.Kind() {
	case reflect.Array:
		require.Equal(t, t1.Len(), t2.Len())
		assertHaveSameStructureAndNames(t, t1.Elem(), t2.Elem())
	case reflect.Slice, reflect.Pointer:
		assertHaveSameStructureAndNames(t, t1.Elem(), t2.Elem())
	case reflect.Struct:
		numFields := t1.NumField()
		require.Equal(t, numFields, t2.NumField())
		for i := 0; i < numFields; i++ {
			require.Equal(t, t1.Field(i).Name, t2.Field(i).Name)
			assertHaveSameStructureAndNames(t, t1.Field(i).Type, t2.Field(i).Type)
		}
	default:
		require.Equal(t, t1, t2)
	}
}
