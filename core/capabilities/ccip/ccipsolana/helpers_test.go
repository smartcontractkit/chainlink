package ccipsolana

import (
	"bytes"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
)

func Test_decodeExtraArgs(t *testing.T) {
	t.Run("decode extra args into map svm", func(t *testing.T) {
		extraArgs := ccip_router.AnyExtraArgs{
			// TODO wait for onchain team fix this AnyExtraArgs
			//ComputeUnits:     1000,
			//IsWritableBitmap: 2,
			//Accounts: []solana.PublicKey{
			//	config.ReceiverExternalExecutionConfigPDA,
			//	config.ReceiverTargetAccountPDA,
			//	solana.SystemProgramID,
			//},
		}

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err := extraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)
		output, err := DecodeExtraArgsToMap(buf.Bytes())
		require.NoError(t, err)
		require.Len(t, output, 3)
	})
}
