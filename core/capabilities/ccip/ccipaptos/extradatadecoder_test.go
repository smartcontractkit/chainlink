package ccipaptos

import (
	"math/big"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func Test_decodeExtraData(t *testing.T) {
	extraDataDecoder := &ExtraDataDecoder{}

	t.Run("decode extra args into map evm v1", func(t *testing.T) {
		encodedGasLimit := uint256.MustFromDecimal("500000")

		encoded := hexutil.MustDecode("0x97a657c920a1070000000000000000000000000000000000000000000000000000000000")
		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, encodedGasLimit.Uint64(), gl.(*big.Int).Uint64(), "Expected %s, got %s", encodedGasLimit.String(), gl.(*big.Int).String())
	})

	t.Run("decode extra args into map evm v2", func(t *testing.T) {
		encodedGasLimit := uint256.MustFromDecimal("500000")
		encodedAllowOOO := true

		// Value generated using aptos contracts
		encoded := hexutil.MustDecode("0x181dcf1020a107000000000000000000000000000000000000000000000000000000000001")
		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 2)

		gl, exist := m["gasLimit"]
		require.True(t, exist)

		require.Equal(t, encodedGasLimit.Uint64(), gl.(*big.Int).Uint64(), "Expected %s, got %s", encodedGasLimit.String(), gl.(*big.Int).String())

		ooe, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, encodedAllowOOO, ooe)
	})

	t.Run("decode extra args into map svm", func(t *testing.T) {
		encodedComputeUnits := uint32(100000)
		encodedBitmap := uint64(255)
		encodedOOO := true
		encodedTokenReceiver := hexutil.MustDecode("0x1234567890123456789012345678901234567890123456789012345678901234")
		encodedAccounts := [][]byte{hexutil.MustDecode("0x1234567890123456789012345678901212345678901234567890123456789012")}

		encoded := hexutil.MustDecode("0x1f3b3abaa0860100ff000000000000000120123456789012345678901234567890123456789012345678901234567890123401201234567890123456789012345678901212345678901234567890123456789012")
		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 5)

		cu, exist := m["computeUnits"]
		require.True(t, exist)
		require.Equal(t, encodedComputeUnits, cu.(uint32), "Expected %d, got %d", encodedComputeUnits, cu.(uint32))

		bitmap, exist := m["accountIsWritableBitmap"]
		require.True(t, exist)
		require.Equal(t, encodedBitmap, bitmap.(uint64), "Expected %d, got %d", encodedBitmap, bitmap.(uint64))

		ooe, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, encodedOOO, ooe)

		tokenReceiver, exist := m["tokenReceiver"]
		require.True(t, exist)
		require.Equal(t, encodedTokenReceiver, tokenReceiver, "Expected %s, got %s", hexutil.Encode(encodedTokenReceiver), hexutil.Encode(tokenReceiver.([]byte)))

		accounts, exist := m["accounts"]
		require.True(t, exist)
		require.Equal(t, encodedAccounts, accounts)
	})

	t.Run("decode extra args into map svm with multiple accounts", func(t *testing.T) {
		encodedComputeUnits := uint32(100000)
		encodedBitmap := uint64(255)
		encodedOOO := true
		encodedTokenReceiver := hexutil.MustDecode("0x1234567890123456789012345678901234567890123456789012345678901234")
		encodedAccounts := [][]byte{hexutil.MustDecode("0x1234567890123456789012345678901212345678901234567890123456789012"), hexutil.MustDecode("0x9ab25d7fff22ac56789012345678901212345678901234567890123456789012")}

		encoded := hexutil.MustDecode("0x1f3b3abaa0860100ff000000000000000120123456789012345678901234567890123456789012345678901234567890123402201234567890123456789012345678901212345678901234567890123456789012209ab25d7fff22ac56789012345678901212345678901234567890123456789012")
		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 5)

		cu, exist := m["computeUnits"]
		require.True(t, exist)
		require.Equal(t, encodedComputeUnits, cu.(uint32), "Expected %d, got %d", encodedComputeUnits, cu.(uint32))

		bitmap, exist := m["accountIsWritableBitmap"]
		require.True(t, exist)
		require.Equal(t, encodedBitmap, bitmap.(uint64), "Expected %d, got %d", encodedBitmap, bitmap.(uint64))

		ooe, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, encodedOOO, ooe)

		tokenReceiver, exist := m["tokenReceiver"]
		require.True(t, exist)
		require.Equal(t, encodedTokenReceiver, tokenReceiver, "Expected %s, got %s", hexutil.Encode(encodedTokenReceiver), hexutil.Encode(tokenReceiver.([]byte)))

		accounts, exist := m["accounts"]
		require.True(t, exist)
		require.Len(t, accounts.([][]byte), 2)
		require.Equal(t, encodedAccounts[0], accounts.([][]byte)[0])
		require.Equal(t, encodedAccounts[1], accounts.([][]byte)[1])
	})

	t.Run("decode dest exec data into map", func(t *testing.T) {
		destGasAmount := uint32(10000)
		encoded, err := bcs.SerializeU32(destGasAmount)
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeDestExecDataToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		decoded, exist := m[aptosDestExecDataKey]
		require.True(t, exist)
		require.Equal(t, destGasAmount, decoded.(uint32)) // Type assert and compare uint32
	})

	t.Run("error on short extra args", func(t *testing.T) {
		shortData := evmExtraArgsV1Tag[:2] // Less than 4 bytes
		_, err := extraDataDecoder.DecodeExtraArgsToMap(shortData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "extra args too short")
	})

	t.Run("error on unknown tag", func(t *testing.T) {
		dataWithUnknownTag := []byte{0xde, 0xad, 0xbe, 0xef}
		dummyData, err := bcs.SerializeU256(*big.NewInt(1))
		require.NoError(t, err)
		dataWithUnknownTag = append(dataWithUnknownTag, dummyData...)
		_, err = extraDataDecoder.DecodeExtraArgsToMap(dataWithUnknownTag)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown extra args tag")
	})

	t.Run("error on malformed evm v1 data", func(t *testing.T) {
		malformedData, err := bcs.SerializeU256(*big.NewInt(1))
		require.NoError(t, err)
		encoded := append([]byte{}, evmExtraArgsV1Tag...)
		encoded = append(encoded, malformedData[:4]...)
		_, err = extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not enough bytes remaining to deserialize u256")
	})

	t.Run("error on malformed dest exec data", func(t *testing.T) {
		malformedData := []byte{0x01, 0x02, 0x03} // Too short for uint32 (expects 4 bytes)
		_, err := extraDataDecoder.DecodeDestExecDataToMap(malformedData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "dest exec data invalid length")
	})

	t.Run("error on dest exec data exceeding uint32 max", func(t *testing.T) {
		tooLargeValue := new(big.Int).Lsh(big.NewInt(1), 32)
		encodedTooLarge, err := bcs.SerializeU256(*tooLargeValue) // Pack as uint256
		require.NoError(t, err)

		_, err = extraDataDecoder.DecodeDestExecDataToMap(encodedTooLarge)
		require.Error(t, err)
		require.Contains(t, err.Error(), "dest exec data invalid length")
	})
}
