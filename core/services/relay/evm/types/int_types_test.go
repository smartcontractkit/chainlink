package types

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestIntTypes(t *testing.T) {
	t.Parallel()
	for i := 24; i <= 256; i += 8 {
		if i == 64 || i == 32 {
			continue
		}
		t.Run(fmt.Sprintf("int%v", i), func(t *testing.T) {
			tpe, ok := GetAbiEncodingType(fmt.Sprintf("int%v", i))
			require.True(t, ok)
			minVal := new(big.Int).Neg(new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i-1)), nil))
			maxVal := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i-1)), nil), big.NewInt(1))
			assertBigIntBounds(t, tpe, minVal, maxVal)
		})

		t.Run(fmt.Sprintf("uint%v", i), func(t *testing.T) {
			tep, ok := GetAbiEncodingType(fmt.Sprintf("uint%v", i))
			require.True(t, ok)
			minVal := big.NewInt(0)
			maxVal := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil), big.NewInt(1))
			assertBigIntBounds(t, tep, minVal, maxVal)
		})
	}
}

func assertBigIntBounds(t *testing.T, tpe *ABIEncodingType, min, max *big.Int) {
	t.Helper()
	assert.Equal(t, reflect.TypeOf(min), tpe.native)
	assert.True(t, tpe.checked.ConvertibleTo(reflect.TypeOf(min)))
	minMinusOne := new(big.Int).Sub(min, big.NewInt(1))
	maxPlusOne := new(big.Int).Add(max, big.NewInt(1))
	sbi := reflect.ValueOf(min).Convert(tpe.checked).Interface().(SizedBigInt)
	assert.NoError(t, sbi.Verify())
	sbi = reflect.ValueOf(max).Convert(tpe.checked).Interface().(SizedBigInt)
	assert.NoError(t, sbi.Verify())
	sbi = reflect.ValueOf(minMinusOne).Convert(tpe.checked).Interface().(SizedBigInt)
	assert.True(t, errors.Is(types.ErrInvalidType, sbi.Verify()))
	sbi = reflect.ValueOf(maxPlusOne).Convert(tpe.checked).Interface().(SizedBigInt)
	assert.True(t, errors.Is(types.ErrInvalidType, sbi.Verify()))
}
