package ccipsolana

import (
	"math/rand"
	"testing"

	solanacommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var randomExecuteReport = func(t *testing.T) cciptypes.ExecutePluginReport {
	const numChainReports = 1
	const msgsPerReport = 1
	const numTokensPerMsg = 3

	chainReports := make([]cciptypes.ExecutePluginReportSingleChain, numChainReports)
	for i := 0; i < numChainReports; i++ {
		reportMessages := make([]cciptypes.Message, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			data, err := cciptypes.NewBytesFromString(utils.RandomAddress().String())
			assert.NoError(t, err)

			tokenAmounts := make([]cciptypes.RampTokenAmount, numTokensPerMsg)
			for z := 0; z < numTokensPerMsg; z++ {
				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: utils.RandomAddress().Bytes(),
					DestTokenAddress:  utils.RandomAddress().Bytes(),
					ExtraData:         data,
					Amount:            cciptypes.NewBigInt(utils.RandUint256()),
				}
			}

			// TODO enable extraArgs ?
			//extraArgs := ccip_router.SolanaExtraArgs{
			//	ComputeUnits: 1000,
			//	Accounts: []ccip_router.SolanaAccountMeta{
			//		{Pubkey: config.CcipReceiverProgram},
			//		{Pubkey: config.ReceiverTargetAccountPDA, IsWritable: true},
			//		{Pubkey: solana.SystemProgramID, IsWritable: false},
			//	},
			//}
			assert.NoError(t, err)
			senderAddr := solanacommon.MakeRandom32ByteArray()
			receiverAddr := solanacommon.MakeRandom32ByteArray()
			feeTokenAddr := solanacommon.MakeRandom32ByteArray()

			reportMessages[j] = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           utils.RandomBytes32(),
					SourceChainSelector: cciptypes.ChainSelector(rand.Uint64()),
					DestChainSelector:   cciptypes.ChainSelector(rand.Uint64()),
					SequenceNumber:      cciptypes.SeqNum(rand.Uint64()),
					Nonce:               rand.Uint64(),
					MsgHash:             utils.RandomBytes32(),
					OnRamp:              utils.RandomAddress().Bytes(),
				},
				Sender:   senderAddr[:],
				Data:     data,
				Receiver: receiverAddr[:],
				//ExtraArgs:      extraArgs,
				FeeToken:       feeTokenAddr[:],
				FeeTokenAmount: cciptypes.NewBigInt(utils.RandUint256()),
				TokenAmounts:   tokenAmounts,
			}
		}

		tokenData := make([][][]byte, numTokensPerMsg)
		for j := 0; j < numTokensPerMsg; j++ {
			tokenData[j] = [][]byte{{0x1}, {0x2, 0x3}}
		}

		chainReports[i] = cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(rand.Uint64()),
			Messages:            reportMessages,
			OffchainTokenData:   tokenData,
			Proofs:              []cciptypes.Bytes32{utils.RandomBytes32(), utils.RandomBytes32()},
			ProofFlagBits:       cciptypes.NewBigInt(utils.RandUint256()),
		}
	}

	return cciptypes.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1(t *testing.T) {
	testCases := []struct {
		name   string
		report func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport
		expErr bool
	}{
		{
			name:   "base report",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr: false,
		},
		{
			name: "reports have empty msgs",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].Messages = []cciptypes.Message{}
				report.ChainReports[4].Messages = []cciptypes.Message{}
				return report
			},
			expErr: false,
		},
		{
			name: "reports have empty offchain token data",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].OffchainTokenData = [][][]byte{}
				report.ChainReports[4].OffchainTokenData[1] = [][]byte{}
				return report
			},
			expErr: false,
		},
	}

	ctx := testutils.Context(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cd := NewExecutePluginCodecV1()
			report := tc.report(randomExecuteReport(t))
			bytes, err := cd.Encode(ctx, report)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// ignore msg hash in comparison
			for i := range report.ChainReports {
				for j := range report.ChainReports[i].Messages {
					report.ChainReports[i].Messages[j].Header.MsgHash = cciptypes.Bytes32{}
					report.ChainReports[i].Messages[j].Header.OnRamp = cciptypes.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeToken = cciptypes.UnknownAddress{}
					report.ChainReports[i].Messages[j].ExtraArgs = cciptypes.Bytes{}
					report.ChainReports[i].Messages[j].FeeTokenAmount = cciptypes.BigInt{}
				}
			}

			// decode using the codec
			codecDecoded, err := cd.Decode(ctx, bytes)
			require.NoError(t, err)
			assert.Equal(t, report, codecDecoded)
		})
	}
}
