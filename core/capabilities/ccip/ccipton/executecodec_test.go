package ccipton

import (
	"context"
	"math/big"
	"math/rand"
	"testing"

	"github.com/gagliardetto/solana-go"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomTONExecuteReport(t *testing.T, sourceChainSelector uint64) cciptypes.ExecutePluginReport {
	const numChainReports = 1
	const msgsPerReport = 1
	const numTokensPerMsg = 1

	chainReports := make([]cciptypes.ExecutePluginReportSingleChain, numChainReports)
	for i := 0; i < numChainReports; i++ {
		reportMessages := make([]cciptypes.Message, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			key, err := solana.NewRandomPrivateKey()
			require.NoError(t, err)
			extraData := []byte{0x12, 0x34}
			tokenReceiver := make([]byte, 36)
			copy(tokenReceiver, key.PublicKey().Bytes())
			tokenAmounts := make([]cciptypes.RampTokenAmount, numTokensPerMsg)
			for z := 0; z < numTokensPerMsg; z++ {
				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: cciptypes.UnknownAddress(key.PublicKey().String()),
					DestTokenAddress:  append(tokenReceiver, 0, 0, 0, 0), // pad to 36 bytes
					ExtraData:         extraData,
					Amount:            cciptypes.NewBigInt(big.NewInt(rand.Int63())),
					DestExecData:      []byte{0, 0, 0, 0},
				}
			}
			reportMessages[j] = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           [32]byte{},
					SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
					DestChainSelector:   cciptypes.ChainSelector(rand.Uint64()),
					SequenceNumber:      cciptypes.SeqNum(rand.Uint64()),
					Nonce:               rand.Uint64(),
				},
				Sender:       cciptypes.UnknownAddress(key.PublicKey().String()),
				Data:         extraData,
				Receiver:     tokenReceiver,
				ExtraArgs:    []byte{},
				TokenAmounts: tokenAmounts,
			}
		}
		chainReports[i] = cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   [][][]byte{{{0x1}, {0x2, 0x3}}},
			Proofs:              []cciptypes.Bytes32{},
		}
	}
	return cciptypes.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1_TON(t *testing.T) {
	ctx := context.Background()
	mockExtraDataCodec := new(mocks.SourceChainExtraDataCodec)
	edc := common.ExtraDataCodec(map[string]common.SourceChainExtraDataCodec{
		"TON": mockExtraDataCodec,
	})
	codec := NewExecutePluginCodecV1(edc)

	//t.Run("encode/decode roundtrip", func(t *testing.T) {
	//	report := randomTONExecuteReport(t, 123456)
	//	encoded, err := codec.Encode(ctx, report)
	//	require.NoError(t, err)
	//	decoded, err := codec.Decode(ctx, encoded)
	//	require.NoError(t, err)
	//	assert.Equal(t, report.ChainReports[0].SourceChainSelector, decoded.ChainReports[0].SourceChainSelector)
	//	assert.Equal(t, report.ChainReports[0].Messages[0].TokenAmounts[0].Amount, decoded.ChainReports[0].Messages[0].TokenAmounts[0].Amount)
	//})

	t.Run("empty report", func(t *testing.T) {
		encoded, err := codec.Encode(ctx, cciptypes.ExecutePluginReport{})
		require.NoError(t, err)
		assert.Nil(t, encoded)
	})
}
