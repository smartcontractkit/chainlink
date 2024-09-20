package codec

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestAddressStringDecodeHook(t *testing.T) {
	t.Parallel()

	var nilAddress *common.Address
	hexString := "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF"
	address := common.HexToAddress(hexString)
	addressToString := address.Hex()
	emptyAddress := common.Address{}
	emptyString := ""

	t.Run("Converts from string to common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(""), reflect.TypeOf(common.Address{}), hexString)
		require.NoError(t, err)
		require.IsType(t, common.Address{}, result)
		assert.Equal(t, address, result)
	})

	t.Run("Converts from string to *common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(""), reflect.TypeOf(&common.Address{}), hexString)
		require.NoError(t, err)
		assert.Equal(t, &address, result)
	})

	t.Run("Converts from *string to common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.PointerTo(reflect.TypeOf("")), reflect.TypeOf(common.Address{}), &hexString)
		require.NoError(t, err)
		require.IsType(t, common.Address{}, result)
		assert.Equal(t, address, result)
	})

	t.Run("Converts from *string to *common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.PointerTo(reflect.TypeOf("")), reflect.TypeOf(&common.Address{}), &hexString)
		require.NoError(t, err)
		assert.Equal(t, &address, result)
	})

	t.Run("Converts from common.Address to string", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(common.Address{}), reflect.TypeOf(""), address)
		require.NoError(t, err)
		require.IsType(t, "", result)
		assert.Equal(t, addressToString, result)
	})

	t.Run("Converts from common.Address to *string", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(common.Address{}), reflect.PointerTo(reflect.TypeOf("")), address)
		require.NoError(t, err)
		assert.Equal(t, &addressToString, result)
	})

	t.Run("Converts from *common.Address to string", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(&common.Address{}), reflect.TypeOf(""), &address)
		require.NoError(t, err)
		assert.Equal(t, addressToString, result)
	})

	t.Run("Converts from *common.Address to *string", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(&common.Address{}), reflect.PointerTo(reflect.TypeOf("")), &address)
		require.NoError(t, err)
		assert.Equal(t, &addressToString, result)
	})

	t.Run("Returns error on invalid hex string", func(t *testing.T) {
		_, err := addressStringDecodeHook(reflect.TypeOf(""), reflect.TypeOf(common.Address{}), "NotAHexString")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
		_, err = addressStringDecodeHook(reflect.TypeOf(""), reflect.TypeOf(&common.Address{}), "NotAHexString")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Returns error on empty string and empty *string", func(t *testing.T) {
		_, err := addressStringDecodeHook(reflect.TypeOf(""), reflect.TypeOf(common.Address{}), emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
		_, err = addressStringDecodeHook(reflect.TypeOf(""), reflect.TypeOf(&common.Address{}), emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
		_, err = addressStringDecodeHook(reflect.PointerTo(reflect.TypeOf("")), reflect.TypeOf(common.Address{}), &emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
		_, err = addressStringDecodeHook(reflect.PointerTo(reflect.TypeOf("")), reflect.TypeOf(&common.Address{}), &emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
	})

	t.Run("Returns error for empty common.Address and empty *common.Address", func(t *testing.T) {
		_, err := addressStringDecodeHook(reflect.TypeOf(common.Address{}), reflect.TypeOf(""), emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty common.Address")
		_, err = addressStringDecodeHook(reflect.TypeOf(common.Address{}), reflect.PointerTo(reflect.TypeOf("")), emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty common.Address")
		_, err = addressStringDecodeHook(reflect.TypeOf(&common.Address{}), reflect.TypeOf(""), &emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty *common.Address")
		_, err = addressStringDecodeHook(reflect.TypeOf(&common.Address{}), reflect.PointerTo(reflect.TypeOf("")), &emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty *common.Address")
	})

	t.Run("Returns nil for nil *string", func(t *testing.T) {
		var nilString *string
		result, err := addressStringDecodeHook(reflect.PointerTo(reflect.TypeOf("")), reflect.TypeOf(common.Address{}), nilString)
		require.NoError(t, err)
		assert.Nil(t, result, "Expected nil to be returned for nil *string input")
		result, err = addressStringDecodeHook(reflect.PointerTo(reflect.TypeOf("")), reflect.TypeOf(&common.Address{}), nilString)
		require.NoError(t, err)
		assert.Nil(t, result, "Expected nil to be returned for nil *string input")
	})
	t.Run("Returns nil for nil *common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(reflect.TypeOf(&common.Address{}), reflect.TypeOf(""), nilAddress)
		require.NoError(t, err)
		assert.Nil(t, result, "Expected nil to be returned for nil common.Address input")
		result, err = addressStringDecodeHook(reflect.TypeOf(&common.Address{}), reflect.PointerTo(reflect.TypeOf("")), nilAddress)
		require.NoError(t, err)
		assert.Nil(t, result, "Expected nil to be returned for nil common.Address input")
	})

	t.Run("Returns input unchanged for unsupported conversion", func(t *testing.T) {
		unsupportedCases := []struct {
			fromType reflect.Type
			toType   reflect.Type
			input    interface{}
		}{
			{fromType: reflect.TypeOf(12345), toType: reflect.TypeOf(common.Address{}), input: 12345},
			{fromType: reflect.TypeOf(12345), toType: reflect.TypeOf(""), input: 12345},
			{fromType: reflect.TypeOf([]byte{}), toType: reflect.TypeOf(common.Address{}), input: []byte{0x01, 0x02, 0x03}},
		}

		for _, tc := range unsupportedCases {
			result, err := addressStringDecodeHook(tc.fromType, tc.toType, tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.input, result, "Expected original value to be returned for unsupported conversion")
		}
	})
}
