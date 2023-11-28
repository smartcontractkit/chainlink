package evm

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
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
		entry := CodecEntry{
			Args: abi.Arguments{
				{Name: "Field1", Type: type1},
				{Name: "Field2", Type: type2},
				{Name: "Field3", Type: type3},
				{Name: "Field4", Type: type4},
			},
		}
		require.NoError(t, entry.Init())
		native := reflect.New(entry.nativeType)
		iNative := reflect.Indirect(native)
		iNative.FieldByName("Field1").Set(reflect.ValueOf(uint16(2)))
		iNative.FieldByName("Field2").Set(reflect.ValueOf("any string"))
		iNative.FieldByName("Field3").Set(reflect.ValueOf(big.NewInt( /*2^24 - 1*/ 16777215)))
		iNative.FieldByName("Field4").Set(reflect.ValueOf(big.NewInt( /*2^23 - 1*/ 8388607)))
		// native and checked point to the same item, even though they have different "types"
		// they have the same memory layout so this is safe per unsafe casting rules, see unsafe.Pointer for details
		checked := reflect.NewAt(entry.checkedType, native.UnsafePointer())
		iChecked := reflect.Indirect(checked)
		checkedField := iChecked.FieldByName("Field3").Interface()

		sbi, ok := checkedField.(types.SizedBigInt)
		require.True(t, ok)
		assert.NoError(t, sbi.Verify())
		bi, ok := iNative.FieldByName("Field3").Interface().(*big.Int)
		require.True(t, ok)
		bi.Add(bi, big.NewInt(1))
		assert.IsType(t, commontypes.ErrInvalidType, sbi.Verify())
		bi, ok = iNative.FieldByName("Field4").Interface().(*big.Int)
		require.True(t, ok)
		bi.Add(bi, big.NewInt(1))
		assert.IsType(t, commontypes.ErrInvalidType, sbi.Verify())
	})

	t.Run("tuples", func(t *testing.T) {
		type1, err := abi.NewType("uint16", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		tupleType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "Field3", Type: "uint24"},
			{Name: "Field4", Type: "int24"},
		})
		require.NoError(t, err)
		entry := CodecEntry{
			Args: abi.Arguments{
				{Name: "Field1", Type: type1},
				{Name: "Field2", Type: tupleType},
			},
		}
		require.NoError(t, entry.Init())
		native := reflect.New(entry.nativeType)
		iNative := reflect.Indirect(native)
		iNative.FieldByName("Field1").Set(reflect.ValueOf(uint16(2)))
		f2 := iNative.FieldByName("Field2")
		f2.FieldByName("Field3").Set(reflect.ValueOf(big.NewInt( /*2^24 - 1*/ 16777215)))
		f2.FieldByName("Field4").Set(reflect.ValueOf(big.NewInt( /*2^23 - 1*/ 8388607)))
		// native and checked point to the same item, even though they have different "types"
		// they have the same memory layout so this is safe per unsafe casting rules, see unsafe.Pointer for details
		checked := reflect.NewAt(entry.checkedType, native.UnsafePointer())
		tuple := reflect.Indirect(checked).FieldByName("Field2")
		checkedField := tuple.FieldByName("Field3").Interface()

		sbi, ok := checkedField.(types.SizedBigInt)
		require.True(t, ok)
		assert.NoError(t, sbi.Verify())
		bi, ok := f2.FieldByName("Field3").Interface().(*big.Int)
		require.True(t, ok)
		bi.Add(bi, big.NewInt(1))
		assert.IsType(t, commontypes.ErrInvalidType, sbi.Verify())
		bi, ok = f2.FieldByName("Field4").Interface().(*big.Int)
		require.True(t, ok)
		bi.Add(bi, big.NewInt(1))
		assert.IsType(t, commontypes.ErrInvalidType, sbi.Verify())
	})

	t.Run("unwrapped types", func(t *testing.T) {
		// This exists to allow you to decode single returned values without naming the parameter
		wrappedTuple, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
			{Name: "Field1", Type: "int16"},
		})
		require.NoError(t, err)
		entry := CodecEntry{
			Args: abi.Arguments{{Name: "", Type: wrappedTuple}},
		}
		require.NoError(t, entry.Init())
		native := reflect.New(entry.nativeType)
		iNative := reflect.Indirect(native)
		iNative.FieldByName("Field1").Set(reflect.ValueOf(int16(2)))
	})

	t.Run("slice types", func(t *testing.T) {
		type1, err := abi.NewType("int16[]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := CodecEntry{
			Args: abi.Arguments{{Name: "Field1", Type: type1}},
		}
		require.NoError(t, entry.Init())
		native := reflect.New(entry.nativeType)
		iNative := reflect.Indirect(native)
		iNative.FieldByName("Field1").Set(reflect.ValueOf([]int16{2, 3}))
	})

	t.Run("array types", func(t *testing.T) {
		type1, err := abi.NewType("int16[3]", "", []abi.ArgumentMarshaling{})
		require.NoError(t, err)
		entry := CodecEntry{
			Args: abi.Arguments{{Name: "Field1", Type: type1}},
		}
		require.NoError(t, entry.Init())
		native := reflect.New(entry.nativeType)
		iNative := reflect.Indirect(native)
		iNative.FieldByName("Field1").Set(reflect.ValueOf([3]int16{2, 3, 30}))
	})
}
