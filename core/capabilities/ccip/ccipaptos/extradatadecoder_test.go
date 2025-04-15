package ccipaptos

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
)

func Test_decodeExtraData(t *testing.T) {
	gasLimit := big.NewInt(rand.Int63())
	extraDataDecoder := &ExtraDataDecoder{}

	t.Run("decode extra args into map evm v1", func(t *testing.T) {
		encodedFields, err := evmExtraArgsV1Fields.Pack(gasLimit)
		require.NoError(t, err)

		encoded := append(evmExtraArgsV1Tag, encodedFields...)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, 0, gl.(*big.Int).Cmp(gasLimit), "Expected %s, got %s", gasLimit.String(), gl.(*big.Int).String()) // Use Cmp for big.Int comparison
	})

	t.Run("decode extra args into map evm v2", func(t *testing.T) {
		allowOoe := true
		encodedFields, err := genericExtraArgsV2Fields.Pack(gasLimit, allowOoe)
		require.NoError(t, err)

		encoded := append(genericExtraArgsV2Tag, encodedFields...)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 2)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, 0, gl.(*big.Int).Cmp(gasLimit), "Expected %s, got %s", gasLimit.String(), gl.(*big.Int).String()) // Use Cmp

		ooeDecoded, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, allowOoe, ooeDecoded) // Check boolean directly
	})

	t.Run("decode extra args into map svm", func(t *testing.T) {
		key, err := solana.NewRandomPrivateKey()
		require.NoError(t, err)
		cu := uint32(10000)
		bitmap := uint64(4)
		ooe := false
		tokenReceiver := [32]byte(key.PublicKey().Bytes())
		accounts := [][32]byte{[32]byte(key.PublicKey().Bytes())} // Example with one account

		encodedFields, err := svmExtraArgsV1Fields.Pack(
			cu,
			bitmap,
			ooe,
			tokenReceiver,
			accounts,
		)
		require.NoError(t, err)

		encoded := append(svmExtraArgsV1Tag, encodedFields...)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 5)

		cuDecoded, exist := m["computeUnits"]
		require.True(t, exist)
		require.Equal(t, cu, cuDecoded)

		bitmapDecoded, exist := m["accountIsWritableBitmap"]
		require.True(t, exist)
		require.Equal(t, bitmap, bitmapDecoded)

		ooeDecoded, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, ooe, ooeDecoded)

		tokenReceiverDecoded, exist := m["tokenReceiver"]
		require.True(t, exist)
		require.Equal(t, tokenReceiver, tokenReceiverDecoded.([32]byte))

		accountsDecoded, exist := m["accounts"]
		require.True(t, exist)
		require.Equal(t, accounts, accountsDecoded.([][32]byte))
	})

	t.Run("decode dest exec data into map", func(t *testing.T) {
		destGasAmount := uint32(10000)
		encoded, err := destGasAmountArguments.Pack(destGasAmount)
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
		unknownTag := []byte{0xde, 0xad, 0xbe, 0xef}
		dummyData, err := evmExtraArgsV1Fields.Pack(big.NewInt(1))
		require.NoError(t, err)
		dataWithUnknownTag := append(unknownTag, dummyData...)
		_, err = extraDataDecoder.DecodeExtraArgsToMap(dataWithUnknownTag)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown extra args tag")
	})

	t.Run("error on malformed evm v1 data", func(t *testing.T) {
		malformedData := []byte{0x01, 0x02, 0x03} // Too short for uint256
		encoded := append(evmExtraArgsV1Tag, malformedData...)
		_, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.Error(t, err)
		require.Contains(t, err.Error(), "abi decode extra args")
	})

	t.Run("error on malformed dest exec data", func(t *testing.T) {
		malformedData := []byte{0x01, 0x02, 0x03} // Too short for uint32 (expects 32 bytes)
		_, err := extraDataDecoder.DecodeDestExecDataToMap(malformedData)
		require.Error(t, err)
		require.Contains(t, err.Error(), "decode dest gas amount")
		require.Contains(t, err.Error(), "expected 32 bytes for uint32")
	})

	t.Run("error on dest exec data exceeding uint32 max", func(t *testing.T) {
		tooLargeValue := new(big.Int).Lsh(big.NewInt(1), 32)
		encodedTooLarge, err := abi.Arguments{{Type: uint256Type}}.Pack(tooLargeValue) // Pack as uint256
		require.NoError(t, err)

		_, err = extraDataDecoder.DecodeDestExecDataToMap(encodedTooLarge)
		require.Error(t, err)
		require.Contains(t, err.Error(), "decode dest gas amount")
		require.Contains(t, err.Error(), "exceeds uint32 max")
	})
}
