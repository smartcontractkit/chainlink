package ccipaptos

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_aptos_utils"
)

func Test_decodeExtraData(t *testing.T) {
	gasLimit := big.NewInt(rand.Int63())
	extraDataDecoder := &ExtraDataDecoder{}

	t.Run("decode extra args into map evm v1", func(t *testing.T) {
		encoded, err := aptosUtilsABI.Pack("exposeEVMExtraArgsV1", ccip_aptos_utils.AptosUtilsEVMExtraArgsV1{
			GasLimit: gasLimit,
		})
		encoded = append(evmExtraArgsV1Tag, encoded[4:]...)
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		gl, exist := m["gasLimit"]
		require.True(t, exist)
		require.Equal(t, gl, gasLimit)
	})

	t.Run("decode extra args into map evm v2", func(t *testing.T) {
		encoded, err := aptosUtilsABI.Pack("exposeGenericExtraArgsV2", ccip_aptos_utils.AptosUtilsGenericExtraArgsV2{
			GasLimit:                 gasLimit,
			AllowOutOfOrderExecution: true,
		})
		encoded = append(genericExtraArgsV2Tag, encoded[4:]...)
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
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
		encoded, err := aptosUtilsABI.Pack("exposeSVMExtraArgsV1", ccip_aptos_utils.AptosUtilsSVMExtraArgsV1{
			ComputeUnits:             cu,
			AccountIsWritableBitmap:  bitmap,
			AllowOutOfOrderExecution: ooe,
			TokenReceiver:            tokenReceiver,
			Accounts:                 accounts,
		})
		encoded = append(svmExtraArgsV1Tag, encoded[4:]...)
		require.NoError(t, err)

		m, err := extraDataDecoder.DecodeExtraArgsToMap(encoded)
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
		m, err := extraDataDecoder.DecodeDestExecDataToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)

		decoded, exist := m[aptosDestExecDataKey]
		require.True(t, exist)
		require.Equal(t, destGasAmount, decoded)
	})
}
