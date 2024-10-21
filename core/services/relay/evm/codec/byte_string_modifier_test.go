package codec_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
)

func TestEVMAddressModifier(t *testing.T) {
	modifier := codec.EVMAddressModifier{}
	validAddressStr := "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"
	validAddressBytes := common.HexToAddress(validAddressStr).Bytes()
	invalidLengthAddressStr := "0xabcdef1234567890abcdef"

	t.Run("EncodeAddress encodes valid Ethereum address bytes", func(t *testing.T) {
		encoded, err := modifier.EncodeAddress(validAddressBytes)
		require.NoError(t, err)
		assert.Equal(t, validAddressStr, encoded)
	})

	t.Run("EncodeAddress returns error for invalid byte length", func(t *testing.T) {
		invalidBytes := []byte(invalidLengthAddressStr)
		_, err := modifier.EncodeAddress(invalidBytes)
		require.Error(t, err)
		assert.Contains(t, err.Error(), commontypes.ErrInvalidType)
	})

	t.Run("DecodeAddress decodes valid Ethereum address", func(t *testing.T) {
		decodedBytes, err := modifier.DecodeAddress(validAddressStr)
		require.NoError(t, err)
		assert.Equal(t, validAddressBytes, decodedBytes)
	})

	t.Run("DecodeAddress returns error for invalid address length", func(t *testing.T) {
		_, err := modifier.DecodeAddress(invalidLengthAddressStr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), commontypes.ErrInvalidType)
	})

	t.Run("DecodeAddress returns error for zero-value address", func(t *testing.T) {
		_, err := modifier.DecodeAddress(common.Address{}.Hex())
		require.Error(t, err)
		assert.Contains(t, err.Error(), commontypes.ErrInvalidType)
	})

	t.Run("Length returns 20 for Ethereum addresses", func(t *testing.T) {
		assert.Equal(t, common.AddressLength, modifier.Length())
	})
}
