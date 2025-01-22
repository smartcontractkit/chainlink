package ccipevm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"
)

func Test_decodeExtraArgs(t *testing.T) {
	d := testSetup(t)
	gasLimit := big.NewInt(rand.Int63())

	t.Run("v1", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
			GasLimit: gasLimit,
		})
		require.NoError(t, err)

		decodedGasLimit, err := decodeExtraArgsV1V2(encoded)
		require.NoError(t, err)

		require.Equal(t, gasLimit, decodedGasLimit)
	})

	t.Run("v2", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV2(nil, message_hasher.ClientEVMExtraArgsV2{
			GasLimit:                 gasLimit,
			AllowOutOfOrderExecution: true,
		})
		require.NoError(t, err)

		decodedGasLimit, err := decodeExtraArgsV1V2(encoded)
		require.NoError(t, err)

		require.Equal(t, gasLimit, decodedGasLimit)
	})

	t.Run("decode extra args into map evm v1", func(t *testing.T) {
		encoded, err := d.contract.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
			GasLimit: gasLimit,
		})
		require.NoError(t, err)

		m, err := DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
		require.Len(t, m, 1)
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
	})

	t.Run("decode extra args into map svm", func(t *testing.T) {
		key, err := solana.NewRandomPrivateKey()
		require.NoError(t, err)
		encoded, err := d.contract.EncodeSVMExtraArgsV1(nil, message_hasher.ClientSVMExtraArgsV1{
			ComputeUnits:             10000,
			AccountIsWritableBitmap:  4,
			AllowOutOfOrderExecution: false,
			TokenReceiver:            [32]byte(key.PublicKey().Bytes()),
			Accounts: [][32]byte{
				[32]byte(key.PublicKey().Bytes()),
			},
		})
		require.NoError(t, err)

		_, err = DecodeExtraArgsToMap(encoded)
		require.NoError(t, err)
	})
}
