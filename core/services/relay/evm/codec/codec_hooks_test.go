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

	// Helper vars
	var nilString *string
	var nilAddress *common.Address
	hexString := "0xDEADBEEFDEADBEEFDEADBEEFDEADBEEFDEADBEEF"
	address := common.HexToAddress(hexString)
	addressToString := address.Hex()
	emptyAddress := common.Address{}
	emptyString := ""
	stringType, stringPtrType := reflect.TypeOf(""), reflect.PointerTo(reflect.TypeOf(""))
	addressType, addressPtrType := reflect.TypeOf(common.Address{}), reflect.TypeOf(&common.Address{})

	t.Run("Converts from string to common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(stringType, addressType, hexString)
		require.NoError(t, err)
		require.IsType(t, common.Address{}, result)
		assert.Equal(t, address, result)
	})

	t.Run("Converts from string to *common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(stringType, addressPtrType, hexString)
		require.NoError(t, err)
		assert.Equal(t, &address, result)
	})

	t.Run("Converts from *string to common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(stringPtrType, addressType, &hexString)
		require.NoError(t, err)
		require.IsType(t, common.Address{}, result)
		assert.Equal(t, address, result)
	})

	t.Run("Converts from *string to *common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(stringPtrType, addressPtrType, &hexString)
		require.NoError(t, err)
		assert.Equal(t, &address, result)
	})

	t.Run("Converts from common.Address to string", func(t *testing.T) {
		result, err := addressStringDecodeHook(addressType, stringType, address)
		require.NoError(t, err)
		require.IsType(t, "", result)
		assert.Equal(t, addressToString, result)
	})

	t.Run("Converts from common.Address to *string", func(t *testing.T) {
		result, err := addressStringDecodeHook(addressType, stringPtrType, address)
		require.NoError(t, err)
		assert.Equal(t, &addressToString, result)
	})

	t.Run("Converts from *common.Address to string", func(t *testing.T) {
		result, err := addressStringDecodeHook(addressPtrType, stringType, &address)
		require.NoError(t, err)
		assert.Equal(t, addressToString, result)
	})

	t.Run("Converts from *common.Address to *string", func(t *testing.T) {
		result, err := addressStringDecodeHook(addressPtrType, stringPtrType, &address)
		require.NoError(t, err)
		assert.Equal(t, &addressToString, result)
	})

	t.Run("Returns error on invalid hex string", func(t *testing.T) {
		_, err := addressStringDecodeHook(stringType, addressType, "NotAHexString")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
		_, err = addressStringDecodeHook(stringType, addressPtrType, "NotAHexString")
		assert.True(t, errors.Is(err, types.ErrInvalidType))
	})

	t.Run("Returns error on empty string and empty *string", func(t *testing.T) {
		_, err := addressStringDecodeHook(stringType, addressType, emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
		_, err = addressStringDecodeHook(stringType, addressPtrType, emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
		_, err = addressStringDecodeHook(stringPtrType, addressType, &emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
		_, err = addressStringDecodeHook(stringPtrType, addressPtrType, &emptyString)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected an error for empty string")
	})

	t.Run("Returns error for empty common.Address and empty *common.Address", func(t *testing.T) {
		_, err := addressStringDecodeHook(addressType, stringType, emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty common.Address")
		_, err = addressStringDecodeHook(addressType, stringPtrType, emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty common.Address")
		_, err = addressStringDecodeHook(addressPtrType, stringType, &emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty *common.Address")
		_, err = addressStringDecodeHook(addressPtrType, stringPtrType, &emptyAddress)
		assert.True(t, errors.Is(err, types.ErrInvalidType), "Expected error for empty *common.Address")
	})
	t.Run("Returns error for nil *string", func(t *testing.T) {
		result, err := addressStringDecodeHook(stringPtrType, addressType, nilString)
		require.Error(t, err, "Expected error for nil *string input")
		assert.Contains(t, err.Error(), "nil *string value")
		assert.Nil(t, result, "Expected result to be nil for nil *string input")

		result, err = addressStringDecodeHook(stringPtrType, addressPtrType, nilString)
		require.Error(t, err, "Expected error for nil *string input")
		assert.Contains(t, err.Error(), "nil *string value")
		assert.Nil(t, result, "Expected result to be nil for nil *string input")
	})

	t.Run("Returns error for nil *common.Address", func(t *testing.T) {
		result, err := addressStringDecodeHook(addressPtrType, stringType, nilAddress)
		require.Error(t, err, "Expected error for nil *common.Address input")
		assert.Contains(t, err.Error(), "nil *common.Address value")
		assert.Nil(t, result, "Expected result to be nil for nil *common.Address input")

		result, err = addressStringDecodeHook(addressPtrType, stringPtrType, nilAddress)
		require.Error(t, err, "Expected error for nil *common.Address input")
		assert.Contains(t, err.Error(), "nil *common.Address value")
		assert.Nil(t, result, "Expected result to be nil for nil *common.Address input")
	})

	t.Run("Returns input unchanged for unsupported conversion", func(t *testing.T) {
		unsupportedCases := []struct {
			fromType reflect.Type
			toType   reflect.Type
			input    interface{}
		}{
			{fromType: reflect.TypeOf(12345), toType: addressType, input: 12345},
			{fromType: reflect.TypeOf(12345), toType: stringType, input: 12345},
			{fromType: reflect.TypeOf([]byte{}), toType: addressType, input: []byte{0x01, 0x02, 0x03}},
		}

		for _, tc := range unsupportedCases {
			result, err := addressStringDecodeHook(tc.fromType, tc.toType, tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.input, result, "Expected original value to be returned for unsupported conversion")
		}
	})
}
