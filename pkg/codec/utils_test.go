package codec

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestGetMapsFromPath(t *testing.T) {
	type testA struct {
		IntSlice []int
	}
	type testB struct {
		TestASlice []testA
	}

	type testStruct struct {
		A    testA
		B    testB
		C, D int
	}

	testMap := map[string]any{"A": map[string]any{"B": []testStruct{{B: testB{TestASlice: []testA{{IntSlice: []int{3, 2, 0}}, {IntSlice: []int{0, 1, 2}}}}, C: 10, D: 100}, {C: 20, D: 200}}}}
	t.Parallel()
	actual, err := getMapsFromPath(testMap, []string{"A"})
	require.NoError(t, err)
	assert.Equal(t, []map[string]any{{"B": []testStruct{{B: testB{TestASlice: []testA{{IntSlice: []int{3, 2, 0}}, {IntSlice: []int{0, 1, 2}}}}, C: 10, D: 100}, {C: 20, D: 200}}}}, actual)

	actual, err = getMapsFromPath(testMap, []string{"A", "B"})
	require.NoError(t, err)
	assert.Equal(t, []map[string]any{{"A": map[string]any{"IntSlice": []int(nil)}, "B": map[string]any{"TestASlice": []testA{{IntSlice: []int{3, 2, 0}}, {IntSlice: []int{0, 1, 2}}}}, "C": 10, "D": 100}, {"A": map[string]any{"IntSlice": []int(nil)}, "B": map[string]any{"TestASlice": []testA(nil)}, "C": 20, "D": 200}}, actual)

	actual, err = getMapsFromPath(testMap, []string{"A", "B", "B"})
	require.NoError(t, err)
	assert.Equal(t, []map[string]any{{"TestASlice": []testA{{IntSlice: []int{3, 2, 0}}, {IntSlice: []int{0, 1, 2}}}}}, actual)

	actual, err = getMapsFromPath(testMap, []string{"A", "B", "B", "TestASlice"})
	require.NoError(t, err)
	assert.Equal(t, []map[string]any{{"IntSlice": []int{3, 2, 0}}, {"IntSlice": []int{0, 1, 2}}}, actual)
}

func TestFitsInNBitsSigned(t *testing.T) {
	t.Parallel()
	t.Run("Fits", func(t *testing.T) {
		bi := big.NewInt(math.MaxInt16)
		assert.True(t, FitsInNBitsSigned(16, bi))
	})

	t.Run("Too large", func(t *testing.T) {
		bi := big.NewInt(math.MaxInt16 + 1)
		assert.False(t, FitsInNBitsSigned(16, bi))
	})

	t.Run("Too small", func(t *testing.T) {
		bi := big.NewInt(math.MinInt16 - 1)
		assert.False(t, FitsInNBitsSigned(16, bi))
	})
}

func TestBigIntHook(t *testing.T) {
	intTypes := []struct {
		Type reflect.Type
		Max  *big.Int
		Min  *big.Int
	}{
		{Type: reflect.TypeOf(0), Min: big.NewInt(math.MinInt), Max: big.NewInt(math.MaxInt)},
		{Type: reflect.TypeOf(uint(0)), Min: big.NewInt(0), Max: new(big.Int).SetUint64(math.MaxUint)},
		{Type: reflect.TypeOf(int8(0)), Min: big.NewInt(math.MinInt8), Max: big.NewInt(math.MaxInt8)},
		{Type: reflect.TypeOf(uint8(0)), Min: big.NewInt(0), Max: new(big.Int).SetUint64(math.MaxUint8)},
		{Type: reflect.TypeOf(int16(0)), Min: big.NewInt(math.MinInt16), Max: big.NewInt(math.MaxInt16)},
		{Type: reflect.TypeOf(uint16(0)), Min: big.NewInt(0), Max: new(big.Int).SetUint64(math.MaxUint16)},
		{Type: reflect.TypeOf(int32(0)), Min: big.NewInt(math.MinInt32), Max: big.NewInt(math.MaxInt32)},
		{Type: reflect.TypeOf(uint32(0)), Min: big.NewInt(0), Max: new(big.Int).SetUint64(math.MaxUint32)},
		{Type: reflect.TypeOf(int64(0)), Min: big.NewInt(math.MinInt64), Max: big.NewInt(math.MaxInt64)},
		{Type: reflect.TypeOf(uint64(0)), Min: big.NewInt(0), Max: new(big.Int).SetUint64(math.MaxUint64)},
	}
	for _, intType := range intTypes {
		t.Run(fmt.Sprintf("Fits conversion %v", intType.Type), func(t *testing.T) {
			anyValidNumber := big.NewInt(5)
			result, err := BigIntHook(reflect.TypeOf((*big.Int)(nil)), intType.Type, anyValidNumber)
			require.NoError(t, err)
			require.IsType(t, reflect.New(intType.Type).Elem().Interface(), result)
			if intType.Min.Cmp(big.NewInt(0)) == 0 {
				u64 := reflect.ValueOf(result).Convert(reflect.TypeOf(uint64(0))).Interface().(uint64)
				actual := new(big.Int).SetUint64(u64)
				require.Equal(t, anyValidNumber, actual)
			} else {
				i64 := reflect.ValueOf(result).Convert(reflect.TypeOf(int64(0))).Interface().(int64)
				actual := big.NewInt(i64)
				require.Equal(t, 0, anyValidNumber.Cmp(actual))
			}
		})

		t.Run("Overflow return an error "+intType.Type.String(), func(t *testing.T) {
			bigger := new(big.Int).Add(intType.Max, big.NewInt(1))
			_, err := BigIntHook(reflect.TypeOf((*big.Int)(nil)), intType.Type, bigger)
			assert.True(t, errors.Is(err, types.ErrInvalidType))
		})

		t.Run("Underflow return an error "+intType.Type.String(), func(t *testing.T) {
			smaller := new(big.Int).Sub(intType.Min, big.NewInt(1))
			_, err := BigIntHook(reflect.TypeOf((*big.Int)(nil)), intType.Type, smaller)
			assert.True(t, errors.Is(err, types.ErrInvalidType))
		})

		t.Run("Converts from "+intType.Type.String(), func(t *testing.T) {
			anyValidNumber := int64(5)
			asType := reflect.ValueOf(anyValidNumber).Convert(intType.Type).Interface()
			result, err := BigIntHook(intType.Type, reflect.TypeOf((*big.Int)(nil)), asType)
			require.NoError(t, err)
			bi, ok := result.(*big.Int)
			require.True(t, ok)
			assert.Equal(t, anyValidNumber, bi.Int64())
		})
	}

	t.Run("Converts from string", func(t *testing.T) {
		anyNumber := int64(5)
		result, err := BigIntHook(reflect.TypeOf(""), reflect.TypeOf((*big.Int)(nil)), strconv.FormatInt(anyNumber, 10))
		require.NoError(t, err)
		bi, ok := result.(*big.Int)
		require.True(t, ok)
		assert.Equal(t, anyNumber, bi.Int64())
	})

	t.Run("Converts to string", func(t *testing.T) {
		anyNumber := int64(5)
		result, err := BigIntHook(reflect.TypeOf((*big.Int)(nil)), reflect.TypeOf(""), big.NewInt(anyNumber))
		require.NoError(t, err)
		assert.Equal(t, strconv.FormatInt(anyNumber, 10), result)
	})

	t.Run("Errors for invalid string", func(t *testing.T) {
		_, err := BigIntHook(reflect.TypeOf(""), reflect.TypeOf((*big.Int)(nil)), "Not a number :(")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Not a big int returns the input data", func(t *testing.T) {
		input := "foo"
		result, err := BigIntHook(reflect.TypeOf(""), reflect.TypeOf(10), input)
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})
}

func TestSliceToArrayVerifySizeHook(t *testing.T) {
	t.Run("correct size slice converts", func(t *testing.T) {
		to := reflect.TypeOf([2]int64{})
		data := []int64{1, 2}
		res, err := SliceToArrayVerifySizeHook(reflect.TypeOf(data), to, data)
		assert.NoError(t, err)

		// Mapstructure will convert slices to arrays, all we need in this hook is to pass it along
		assert.Equal(t, data, res)
	})

	t.Run("Too large slice returns error", func(t *testing.T) {
		to := reflect.TypeOf([2]int64{})
		data := []int64{1, 2, 3}
		_, err := SliceToArrayVerifySizeHook(reflect.TypeOf(data), to, data)
		assert.True(t, errors.Is(err, types.ErrSliceWrongLen))
	})

	t.Run("Too small slice returns error", func(t *testing.T) {
		to := reflect.TypeOf([2]int64{})
		data := []int64{1}
		_, err := SliceToArrayVerifySizeHook(reflect.TypeOf(data), to, data)
		assert.True(t, errors.Is(err, types.ErrSliceWrongLen))
	})

	t.Run("Empty slices are treated as ok to allow unset values", func(t *testing.T) {
		to := reflect.TypeOf([2]int64{})
		var data []int64
		res, err := SliceToArrayVerifySizeHook(reflect.TypeOf(data), to, data)
		assert.NoError(t, err)

		// Mapstructure will convert slices to arrays, all we need in this hook is to pass it along
		assert.Equal(t, []int64{0, 0}, res)
	})

	t.Run("Not a slice returns the input data", func(t *testing.T) {
		input := "foo"
		result, err := BigIntHook(reflect.TypeOf(""), reflect.TypeOf(10), input)
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})
}

func TestEpochToTimeHook(t *testing.T) {
	anyTime := int64(math.MaxInt8 - 40)
	testTime := time.Unix(anyTime, 0).UTC()
	testValues := []any{
		int(anyTime),
		uint(anyTime),
		int8(anyTime),
		uint8(anyTime),
		int16(anyTime),
		uint16(anyTime),
		int32(anyTime),
		uint32(anyTime),
		anyTime,
		uint64(anyTime),
	}

	t.Run("converts epoch to time", func(t *testing.T) {
		for _, testValue := range testValues {
			t.Run(fmt.Sprintf("%T", testValue), func(t *testing.T) {
				actual, err := EpochToTimeHook(reflect.TypeOf(testValue), reflect.TypeOf(testTime), testValue)
				require.NoError(t, err)
				assert.Equal(t, testTime, actual)
			})
		}
	})

	t.Run("Converts timestamps to integer type", func(t *testing.T) {
		for _, testValue := range testValues {
			t.Run(fmt.Sprintf("%T", testValue), func(t *testing.T) {
				actual, err := EpochToTimeHook(reflect.TypeOf(testTime), reflect.TypeOf(testValue), testTime)
				require.NoError(t, err)
				assert.Equal(t, testValue, actual)
			})
		}
	})

	t.Run("returns data for non time types", func(t *testing.T) {
		actual, err := EpochToTimeHook(reflect.TypeOf(""), reflect.TypeOf(0), "foo")
		require.NoError(t, err)
		assert.Equal(t, "foo", actual)
	})

	t.Run("pointers are maintained in non-conversion scenarios", func(t *testing.T) {
		t.Run("*time.Time to *time.Time", func(t *testing.T) {
			tp := reflect.PointerTo(reflect.TypeOf(testTime))
			output, err := EpochToTimeHook(tp, tp, &testTime)

			require.NoError(t, err)

			value, ok := output.(*time.Time)

			require.True(t, ok)
			require.Equal(t, testTime.Unix(), value.Unix())
		})

		t.Run("*time.Time to time.Time", func(t *testing.T) {
			output, err := EpochToTimeHook(reflect.PointerTo(reflect.TypeOf(testTime)), reflect.TypeOf(testTime), &testTime)

			require.NoError(t, err)

			value, ok := output.(time.Time)

			require.True(t, ok)
			require.Equal(t, testTime.Unix(), value.Unix())
		})

		t.Run("time.Time to *time.Time", func(t *testing.T) {
			output, err := EpochToTimeHook(reflect.TypeOf(testTime), reflect.PointerTo(reflect.TypeOf(testTime)), testTime)

			require.NoError(t, err)

			value, ok := output.(*time.Time)

			require.True(t, ok)
			require.Equal(t, testTime.Unix(), value.Unix())
		})
	})

	t.Run("Converts timestamps to integer type using mapstructure", func(t *testing.T) {
		type A struct {
			Val1 time.Time
			Val2 time.Time
		}

		type B struct {
			Val1 int64
			Val2 *big.Int
		}

		input := A{
			Val1: testTime,
			Val2: testTime.Add(time.Hour),
		}

		var output B

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: EpochToTimeHook,
			Result:     &output,
		})

		require.NoError(t, err)
		require.NoError(t, decoder.Decode(input))

		expected := B{
			Val1: testTime.Unix(),
			Val2: new(big.Int).SetInt64(testTime.Add(time.Hour).Unix()),
		}

		require.Equal(t, expected, output)
	})
}
