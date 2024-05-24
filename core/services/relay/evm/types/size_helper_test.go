package types_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const anyNumElements = 10

func TestGetMaxSize(t *testing.T) {
	t.Run("Basic types all encode to 32 bytes", func(t *testing.T) {
		args := abi.Arguments{
			{Name: "I8", Type: mustType(t, "int8")},
			{Name: "I80", Type: mustType(t, "int80")},
			{Name: "I256", Type: mustType(t, "int256")},
			{Name: "B3", Type: mustType(t, "bytes3")},
			{Name: "B32", Type: mustType(t, "bytes32")},
			{Name: "TF", Type: mustType(t, "bool")},
		}

		runSizeTest(t, anyNumElements, args, int8(9), big.NewInt(3), big.NewInt(200), [3]byte{1, 3, 4}, make32Bytes(1), true)
	})

	t.Run("Slices of basic types all encode to 32 bytes each + header and footer", func(t *testing.T) {
		args := abi.Arguments{
			{Name: "I8", Type: mustType(t, "int8[]")},
			{Name: "I80", Type: mustType(t, "int80[]")},
			{Name: "I256", Type: mustType(t, "int256[]")},
			{Name: "B3", Type: mustType(t, "bytes3[]")},
			{Name: "B32", Type: mustType(t, "bytes32[]")},
			{Name: "TF", Type: mustType(t, "bool[]")},
		}

		i8 := []int8{9, 2, 1, 3, 5, 6, 2, 1, 2, 3}
		i80 := []*big.Int{big.NewInt(9), big.NewInt(2), big.NewInt(1), big.NewInt(3), big.NewInt(5), big.NewInt(6), big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(3)}
		i256 := []*big.Int{big.NewInt(119), big.NewInt(112), big.NewInt(1), big.NewInt(3), big.NewInt(5), big.NewInt(6), big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(3)}
		b3 := [][3]byte{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}}
		b32 := [][32]byte{make32Bytes(1), make32Bytes(2), make32Bytes(3), make32Bytes(4), make32Bytes(5), make32Bytes(6), make32Bytes(7), make32Bytes(8), make32Bytes(9), make32Bytes(10)}
		tf := []bool{true, false, true, false, true, false, true, false, true, false}
		runSizeTest(t, anyNumElements, args, i8, i80, i256, b3, b32, tf)
	})

	t.Run("Arrays of basic types all encode to 32 bytes each", func(t *testing.T) {
		args := abi.Arguments{
			{Name: "I8", Type: mustType(t, "int8[3]")},
			{Name: "I80", Type: mustType(t, "int80[3]")},
			{Name: "I256", Type: mustType(t, "int256[3]")},
			{Name: "B3", Type: mustType(t, "bytes3[3]")},
			{Name: "B32", Type: mustType(t, "bytes32[3]")},
			{Name: "TF", Type: mustType(t, "bool[3]")},
		}

		i8 := [3]int8{9, 2, 1}
		i80 := [3]*big.Int{big.NewInt(9), big.NewInt(2), big.NewInt(1)}
		i256 := [3]*big.Int{big.NewInt(119), big.NewInt(112), big.NewInt(1)}
		b3 := [3][3]byte{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}
		b32 := [3][32]byte{make32Bytes(1), make32Bytes(2), make32Bytes(3)}
		tf := [3]bool{true, false, true}
		runSizeTest(t, anyNumElements, args, i8, i80, i256, b3, b32, tf)
	})

	t.Run("Tuples are a sum of their elements", func(t *testing.T) {
		tuple1 := []abi.ArgumentMarshaling{
			{Name: "I8", Type: "int8"},
			{Name: "I80", Type: "int80"},
			{Name: "I256", Type: "int256"},
			{Name: "B3", Type: "bytes3"},
			{Name: "B32", Type: "bytes32"},
			{Name: "TF", Type: "bool"},
		}
		t1, err := abi.NewType("tuple", "", tuple1)
		require.NoError(t, err)

		tuple2 := []abi.ArgumentMarshaling{
			{Name: "I80", Type: "int80"},
			{Name: "TF", Type: "bool"},
		}
		t2, err := abi.NewType("tuple", "", tuple2)
		require.NoError(t, err)

		args := abi.Arguments{
			{Name: "t1", Type: t1},
			{Name: "t2", Type: t2},
		}
		arg1 := struct {
			I8   int8
			I80  *big.Int
			I256 *big.Int
			B3   [3]byte
			B32  [32]byte
			TF   bool
		}{
			int8(9), big.NewInt(3), big.NewInt(200), [3]byte{1, 3, 4}, make32Bytes(1), true,
		}

		arg2 := struct {
			I80 *big.Int
			TF  bool
		}{
			big.NewInt(3), true,
		}
		runSizeTest(t, anyNumElements, args, arg1, arg2)
	})

	t.Run("Slices of tuples are a sum of their elements with header and footer", func(t *testing.T) {
		tuple1 := []abi.ArgumentMarshaling{
			{Name: "I80", Type: "int80"},
			{Name: "TF", Type: "bool"},
		}
		t1, err := abi.NewType("tuple[]", "", tuple1)
		require.NoError(t, err)

		args := abi.Arguments{
			{Name: "t1", Type: t1},
		}
		arg1 := []struct {
			I80 *big.Int
			TF  bool
		}{
			{big.NewInt(1), true},
			{big.NewInt(2), true},
			{big.NewInt(3), true},
			{big.NewInt(4), false},
			{big.NewInt(5), true},
			{big.NewInt(6), true},
			{big.NewInt(7), true},
			{big.NewInt(8), false},
			{big.NewInt(9), true},
			{big.NewInt(10), true},
		}
		runSizeTest(t, anyNumElements, args, arg1)
	})

	t.Run("Arrays of tuples are a sum of their elements", func(t *testing.T) {
		tuple1 := []abi.ArgumentMarshaling{
			{Name: "I80", Type: "int80"},
			{Name: "TF", Type: "bool"},
		}
		t1, err := abi.NewType("tuple[3]", "", tuple1)
		require.NoError(t, err)

		args := abi.Arguments{
			{Name: "t1", Type: t1},
		}
		arg1 := []struct {
			I80 *big.Int
			TF  bool
		}{
			{big.NewInt(1), true},
			{big.NewInt(2), true},
			{big.NewInt(3), true},
		}
		runSizeTest(t, anyNumElements, args, arg1)
	})

	t.Run("Bytes pack themselves", func(t *testing.T) {
		args := abi.Arguments{{Name: "B", Type: mustType(t, "bytes")}}
		t.Run("No padding needed", func(t *testing.T) {
			padded := []byte("12345789022345678903234567890412345678905123456789061234")
			runSizeTest(t, 64, args, padded)
		})
		t.Run("Padding needed", func(t *testing.T) {
			needsPadding := []byte("12345789022345678903234567890412345678905123456")
			runSizeTest(t, 56, args, needsPadding)
		})
	})

	t.Run("Strings pack themselves", func(t *testing.T) {
		args := abi.Arguments{{Name: "B", Type: mustType(t, "string")}}
		t.Run("No padding needed", func(t *testing.T) {
			padded := "12345789022345678903234567890412345678905123456789061234"
			runSizeTest(t, 64, args, padded)
		})
		t.Run("Padding needed", func(t *testing.T) {
			needsPadding := "12345789022345678903234567890412345678905123456"
			runSizeTest(t, 56, args, needsPadding)
		})
	})

	t.Run("Nested dynamic types return errors", func(t *testing.T) {
		t.Run("Slice in slice", func(t *testing.T) {
			args := abi.Arguments{{Name: "B", Type: mustType(t, "int32[][]")}}
			_, err := types.GetMaxSize(anyNumElements, args)
			assert.IsType(t, commontypes.ErrInvalidType, err)
		})
		t.Run("Slice in array", func(t *testing.T) {
			args := abi.Arguments{{Name: "B", Type: mustType(t, "int32[][2]")}}
			_, err := types.GetMaxSize(anyNumElements, args)
			assert.IsType(t, commontypes.ErrInvalidType, err)
		})
	})

	t.Run("Slices in a top level tuple works as-if they are the sized element", func(t *testing.T) {
		tuple1 := []abi.ArgumentMarshaling{
			{Name: "I80", Type: "int80[]"},
			{Name: "TF", Type: "bool[]"},
		}
		t1, err := abi.NewType("tuple", "", tuple1)
		require.NoError(t, err)
		args := abi.Arguments{{Name: "tuple", Type: t1}}

		arg1 := struct {
			I80 []*big.Int
			TF  []bool
		}{
			I80: []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5), big.NewInt(6), big.NewInt(7), big.NewInt(8), big.NewInt(9), big.NewInt(10)},
			TF:  []bool{true, true, true, false, true, true, true, false, true, true},
		}

		runSizeTest(t, anyNumElements, args, arg1)
	})

	t.Run("Nested dynamic tuples return errors", func(t *testing.T) {
		tuple1 := []abi.ArgumentMarshaling{
			{Name: "I8", Type: "int8"},
			{Name: "I80", Type: "int80"},
			{Name: "I256", Type: "int256"},
			{Name: "B3", Type: "bytes3"},
			{Name: "B32", Type: "bytes32"},
			{Name: "TF", Type: "bool[]"},
		}

		tuple2 := []abi.ArgumentMarshaling{
			{Name: "I80", Type: "int80"},
			{Name: "T1", Type: "tuple", Components: tuple1},
		}
		t2, err := abi.NewType("tuple", "", tuple2)
		require.NoError(t, err)

		args := abi.Arguments{{Name: "t2", Type: t2}}
		_, err = types.GetMaxSize(anyNumElements, args)
		assert.IsType(t, commontypes.ErrInvalidType, err)
	})
}

func runSizeTest(t *testing.T, n int, args abi.Arguments, params ...any) {
	actual, err := types.GetMaxSize(n, args)
	require.NoError(t, err)

	expected, err := args.Pack(params...)
	require.NoError(t, err)
	assert.Equal(t, len(expected), actual)
}

func mustType(t *testing.T, name string) abi.Type {
	aType, err := abi.NewType(name, "", []abi.ArgumentMarshaling{})
	require.NoError(t, err)
	return aType
}

func make32Bytes(firstByte byte) [32]byte {
	return [32]byte{firstByte, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3}
}
