package ccipevm

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/message_hasher"

	"github.com/stretchr/testify/require"
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
}
