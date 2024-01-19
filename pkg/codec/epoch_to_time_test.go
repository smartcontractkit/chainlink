package codec_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestTimeToUnix(t *testing.T) {
	t.Parallel()
	type testStruct struct {
		A string
		T int64
	}
	tst := reflect.TypeOf(&testStruct{})

	type testSliceStruct struct{ T []int64 }
	tsst := reflect.TypeOf(&testSliceStruct{})

	type testArrayStruct struct{ T [2]int64 }
	tast := reflect.TypeOf(&testArrayStruct{})

	type otherIntegerType struct {
		A string
		T uint32
	}
	oit := reflect.TypeOf(&otherIntegerType{})

	type bigIntType struct {
		A string
		T *big.Int
	}
	bit := reflect.TypeOf(&bigIntType{})

	type bigIntAlias big.Int
	type bigIntAliasType struct {
		A string
		T *bigIntAlias
	}
	biat := reflect.TypeOf(&bigIntAliasType{})

	type testInvalidStruct struct{ T string }

	anyTimeEpoch := int64(631515600)
	testTime := time.Unix(anyTimeEpoch, 0).UTC()
	anyTimeEpoch2 := int64(631515601)
	testTime2 := time.Unix(anyTimeEpoch2, 0).UTC()

	t.Run("RetypeToOffChain returns error if type is not an integer type", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		_, err := converter.RetypeToOffChain(reflect.TypeOf(testInvalidStruct{}), "")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("RetypeToOffChain converts integer types", func(t *testing.T) {
		for _, test := range []struct {
			name string
			t    reflect.Type
		}{
			{"int64", tst},
			{"other integer types", oit},
			{"big.Int", bit},
			{"big.Int alias", biat},
		} {
			t.Run(test.name, func(t *testing.T) {
				converter := codec.NewEpochToTimeModifier([]string{"T"})
				convertedType, err := converter.RetypeToOffChain(test.t, "")

				require.NoError(t, err)
				assert.Equal(t, reflect.Pointer, convertedType.Kind())
				convertedType = convertedType.Elem()

				require.Equal(t, 2, convertedType.NumField())
				assert.Equal(t, tst.Elem().Field(0), convertedType.Field(0))
				assert.Equal(t, tst.Elem().Field(1).Name, convertedType.Field(1).Name)
				assert.Equal(t, reflect.TypeOf(&time.Time{}), convertedType.Field(1).Type)
			})
		}
	})

	t.Run("RetypeToOffChain converts slices", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		convertedType, err := converter.RetypeToOffChain(tsst, "")

		require.NoError(t, err)
		assert.Equal(t, reflect.Pointer, convertedType.Kind())
		convertedType = convertedType.Elem()

		require.Equal(t, 1, convertedType.NumField())
		assert.Equal(t, tsst.Elem().Field(0).Name, convertedType.Field(0).Name)
		assert.Equal(t, reflect.TypeOf([]*time.Time{}), convertedType.Field(0).Type)
	})

	t.Run("RetypeToOffChain converts arrays", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		convertedType, err := converter.RetypeToOffChain(tast, "")

		require.NoError(t, err)
		assert.Equal(t, reflect.Pointer, convertedType.Kind())
		convertedType = convertedType.Elem()

		require.Equal(t, 1, convertedType.NumField())
		assert.Equal(t, tast.Elem().Field(0).Name, convertedType.Field(0).Name)
		assert.Equal(t, reflect.TypeOf([]*time.Time{}), convertedType.Field(0).Type)
	})

	t.Run("TransformToOnChain converts time to integer types", func(t *testing.T) {
		anyString := "test"
		for _, test := range []struct {
			name     string
			t        reflect.Type
			expected any
		}{
			{"int64", tst, &testStruct{A: anyString, T: anyTimeEpoch}},
			{"other integer types", oit, &otherIntegerType{A: anyString, T: uint32(anyTimeEpoch)}},
			{"big.Int", bit, &bigIntType{A: anyString, T: big.NewInt(anyTimeEpoch)}},
			{"big.Int alias", biat, &bigIntAliasType{A: anyString, T: (*bigIntAlias)(big.NewInt(anyTimeEpoch))}},
		} {
			t.Run(test.name, func(t *testing.T) {
				converter := codec.NewEpochToTimeModifier([]string{"T"})
				convertedType, err := converter.RetypeToOffChain(test.t, "")
				require.NoError(t, err)

				rOffchain := reflect.New(convertedType.Elem())
				iOffChain := reflect.Indirect(rOffchain)
				iOffChain.FieldByName("A").SetString(anyString)
				iOffChain.FieldByName("T").Set(reflect.ValueOf(&testTime))

				actual, err := converter.TransformToOnChain(rOffchain.Interface(), "")
				require.NoError(t, err)

				assert.Equal(t, test.expected, actual)
			})
		}
	})

	t.Run("TransformToOnChain converts times to integer array", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		convertedType, err := converter.RetypeToOffChain(tast, "")
		require.NoError(t, err)

		rOffchain := reflect.New(convertedType.Elem())
		iOffChain := reflect.Indirect(rOffchain)
		iOffChain.FieldByName("T").Set(reflect.ValueOf([]*time.Time{&testTime, &testTime2}))

		actual, err := converter.TransformToOnChain(rOffchain.Interface(), "")
		require.NoError(t, err)

		expected := &testArrayStruct{T: [2]int64{anyTimeEpoch, anyTimeEpoch2}}
		assert.Equal(t, expected, actual)
	})

	t.Run("TransformToOnChain converts times to integer slice", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		convertedType, err := converter.RetypeToOffChain(tsst, "")
		require.NoError(t, err)

		rOffchain := reflect.New(convertedType.Elem())
		iOffChain := reflect.Indirect(rOffchain)
		iOffChain.FieldByName("T").Set(reflect.ValueOf([]*time.Time{&testTime, &testTime2}))

		actual, err := converter.TransformToOnChain(rOffchain.Interface(), "")
		require.NoError(t, err)

		expected := &testSliceStruct{T: []int64{anyTimeEpoch, anyTimeEpoch2}}
		assert.Equal(t, expected, actual)
	})

	t.Run("TransformToOffChain converts integer to *time.Time", func(t *testing.T) {
		anyString := "test"
		for _, test := range []struct {
			name     string
			t        reflect.Type
			offChain any
		}{
			{"int64", tst, &testStruct{A: anyString, T: anyTimeEpoch}},
			{"other integer types", oit, &otherIntegerType{A: anyString, T: uint32(anyTimeEpoch)}},
			{"big.Int", bit, &bigIntType{A: anyString, T: big.NewInt(anyTimeEpoch)}},
			{"big.Int alias", biat, &bigIntAliasType{A: anyString, T: (*bigIntAlias)(big.NewInt(anyTimeEpoch))}},
		} {
			t.Run(test.name, func(t *testing.T) {
				converter := codec.NewEpochToTimeModifier([]string{"T"})
				convertedType, err := converter.RetypeToOffChain(test.t, "")
				require.NoError(t, err)

				actual, err := converter.TransformToOffChain(test.offChain, "")
				require.NoError(t, err)

				expected := reflect.New(convertedType.Elem())
				iOffChain := reflect.Indirect(expected)
				iOffChain.FieldByName("A").SetString(anyString)
				iOffChain.FieldByName("T").Set(reflect.ValueOf(&testTime))
				assert.Equal(t, expected.Interface(), actual)
			})
		}
	})

	t.Run("TransformToOffChain converts times to integer array", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		convertedType, err := converter.RetypeToOffChain(tast, "")
		require.NoError(t, err)

		actual, err := converter.TransformToOffChain(&testArrayStruct{T: [2]int64{anyTimeEpoch, anyTimeEpoch2}}, "")
		require.NoError(t, err)

		expected := reflect.New(convertedType.Elem())
		iOffChain := reflect.Indirect(expected)
		iOffChain.FieldByName("T").Set(reflect.ValueOf([]*time.Time{&testTime, &testTime2}))
		assert.Equal(t, expected.Interface(), actual)
	})

	t.Run("TransformToOffChain converts times to integer slice", func(t *testing.T) {
		converter := codec.NewEpochToTimeModifier([]string{"T"})
		convertedType, err := converter.RetypeToOffChain(tsst, "")
		require.NoError(t, err)

		actual, err := converter.TransformToOffChain(&testSliceStruct{T: []int64{anyTimeEpoch, anyTimeEpoch2}}, "")
		require.NoError(t, err)

		expected := reflect.New(convertedType.Elem())
		iOffChain := reflect.Indirect(expected)
		iOffChain.FieldByName("T").Set(reflect.ValueOf([]*time.Time{&testTime, &testTime2}))
		assert.Equal(t, expected.Interface(), actual)
	})
}
