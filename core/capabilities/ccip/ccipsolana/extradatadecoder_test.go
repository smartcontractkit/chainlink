package ccipsolana

import (
	"bytes"
	"encoding/binary"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"

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
		destGasAmount := uint32(10000)
		bitmap := uint64(0)
		extraArgs := ccip_router.SVMExtraArgsV1{
			ComputeUnits:             destGasAmount,
			AccountIsWritableBitmap:  bitmap,
			AllowOutOfOrderExecution: false,
			TokenReceiver:            config.CcipLogicReceiver,
			Accounts: [][32]byte{
				[32]byte(config.CcipLogicReceiver.Bytes()),
				[32]byte(config.ReceiverTargetAccountPDA.Bytes()),
				[32]byte(solana.SystemProgramID.Bytes()),
			},
		}

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err := extraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)
		output, err := DecodeExtraArgsToMap(append(svmExtraArgsV1Tag, buf.Bytes()...))
		require.NoError(t, err)
		require.Len(t, output, 5)

		gasLimit, exist := output["ComputeUnits"]
		require.True(t, exist)
		require.Equal(t, destGasAmount, gasLimit)

		writableBitmap, exist := output["AccountIsWritableBitmap"]
		require.True(t, exist)
		require.Equal(t, bitmap, writableBitmap)

		ooe, exist := output["AllowOutOfOrderExecution"]
		require.True(t, exist)
		require.Equal(t, false, ooe)
	})

	t.Run("decode extra args into map evm", func(t *testing.T) {
		extraArgs := ccip_router.EVMExtraArgsV2{
			GasLimit:                 agbinary.Uint128{Lo: 5000, Hi: 0},
			AllowOutOfOrderExecution: false,
		}

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err := extraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		output, err := DecodeExtraArgsToMap(append(evmExtraArgsV2Tag, buf.Bytes()...))
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
