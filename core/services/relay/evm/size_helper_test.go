package evm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

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

		runSizeTest(t, args, int8(9), big.NewInt(3), big.NewInt(200), [3]byte{1, 3, 4}, make32Bytes(1), true)
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
		runSizeTest(t, args, i8, i80, i256, b3, b32, tf)
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
		runSizeTest(t, args, i8, i80, i256, b3, b32, tf)
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
		runSizeTest(t, args, arg1, arg2)
	})

	t.Run("Slices of tuples are a sum of their elements with header and footer", func(t *testing.T) {
		assert.Fail(t, "not written yet")
	})

	t.Run("Arrays of tuples are a sum of their elements", func(t *testing.T) {
		assert.Fail(t, "not written yet")
	})

	t.Run("Bytes pack themselves", func(t *testing.T) {
		t.Run("No padding needed", func(t *testing.T) {
			assert.Fail(t, "not written yet")
		})
		t.Run("Padding needed", func(t *testing.T) {
			assert.Fail(t, "not written yet")
		})
	})

	t.Run("Nested dynamic types return errors", func(t *testing.T) {
		t.Run("Slice in slice", func(t *testing.T) {
			assert.Fail(t, "not written yet")
		})
		t.Run("Slice in array", func(t *testing.T) {
			assert.Fail(t, "not written yet")
		})
		t.Run("Slice in tuple", func(t *testing.T) {
			assert.Fail(t, "not written yet")
		})
	})

	t.Run("Dynamic tuples return errors", func(t *testing.T) {
		assert.Fail(t, "not written yet")
	})
}

func runSizeTest(t *testing.T, args abi.Arguments, params ...any) {
	anyNumElements := 10

	actual, err := evm.GetMaxSize(anyNumElements, args)
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
