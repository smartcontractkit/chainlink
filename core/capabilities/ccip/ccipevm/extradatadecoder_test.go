package ccipevm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/message_hasher"
)

func Test_decodeExtraData(t *testing.T) {
	d := testSetup(t)
	gasLimit := big.NewInt(rand.Int63())

	t.Run("decode extra args into map evm v1", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
			GasLimit: gasLimit,
		})
		require.NoError(t, err)

		m, err := DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, gl, gasLimit)
	})

	t.Run("decode extra args into map evm v2", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV2(nil, message_hasher.ClientEVMExtraArgsV2{
			GasLimit:                 gasLimit,
			AllowOutOfOrderExecution: true,
		})
		require.NoError(t, err)

		m, err := DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 2)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, gl, gasLimit)

		ooe, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, true, ooe)
	})

	t.Run("decode extra args into map svm", func(t *testing.T) {
		key, err := solana.NewRandomPrivateKey()
		require.NoError(t, err)
		cu := uint32(10000)
		bitmap := uint64(4)
		ooe := false
		tokenReceiver := [32]byte(key.PublicKey().Bytes())
		accounts := [][32]byte{[32]byte(key.PublicKey().Bytes())}
		decoded, err := d.contract.DecodeSVMExtraArgsV1(nil, cu, bitmap, ooe, tokenReceiver, accounts)
		if err != nil {
			return
		}
		encoded, err := d.contract.EncodeSVMExtraArgsV1(nil, decoded)
		require.NoError(t, err)

		m, err := DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 5)

		cuDecoded, exist := m["computeUnits"]
		require.True(t, exist)
		require.Equal(t, cuDecoded, cu)

		bitmapDecoded, exist := m["accountIsWritableBitmap"]
		require.True(t, exist)
		require.Equal(t, bitmapDecoded, bitmap)

		ooeDecoded, exist := m["allowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, ooeDecoded, ooe)

		tokenReceiverDecoded, exist := m["tokenReceiver"]
		require.True(t, exist)
		require.Equal(t, tokenReceiverDecoded, tokenReceiver)

		accountsDecoded, exist := m["accounts"]
		require.True(t, exist)
		require.Equal(t, accountsDecoded, accounts)
	})

	t.Run("decode dest exec data into map", func(t *testing.T) {
		destGasAmount := uint32(10000)
		encoded, err := abiEncodeUint32(destGasAmount)
		require.NoError(t, err)
		m, err := DecodeDestExecDataToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		decoded, exist := m[evmDestExecDataKey]
		require.True(t, exist)
		require.Equal(t, destGasAmount, decoded)
	})
}
