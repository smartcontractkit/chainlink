package ccipsolana

import (
	"bytes"
	"encoding/binary"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
)

func Test_decodeExtraArgs(t *testing.T) {
	t.Run("decode dest exec data into map svm", func(t *testing.T) {
		destGasAmount := uint32(10000)
		encoded := make([]byte, 4)
		binary.LittleEndian.PutUint32(encoded, destGasAmount)
		output, err := DecodeDestExecDataToMap(encoded)
		require.NoError(t, err)

		decoded, exist := output[svmDestExecDataKey]
		require.True(t, exist)
		require.Equal(t, destGasAmount, decoded)
	})

	t.Run("decode extra args into map svm", func(t *testing.T) {
		extraArgs := ccip_router.AnyExtraArgs{
			GasLimit:                 agbinary.Uint128{Lo: 5000, Hi: 0},
			AllowOutOfOrderExecution: false,
		}

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err := extraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)
		output, err := DecodeExtraArgsToMap(buf.Bytes())
		require.NoError(t, err)
		require.Len(t, output, 2)

		gasLimit, exist := output["GasLimit"]
		require.True(t, exist)
		require.Equal(t, agbinary.Uint128{Lo: 5000, Hi: 0}, gasLimit)

		ooe, exist := output["AllowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, false, ooe)
	})
}
