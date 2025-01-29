package ccipsolana

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"math/rand"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/evm/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var randomExecuteReport = func(t *testing.T, sourceChainSelector uint64) cciptypes.ExecutePluginReport {
	const numChainReports = 1
	const msgsPerReport = 1
	const numTokensPerMsg = 1

	chainReports := make([]cciptypes.ExecutePluginReportSingleChain, numChainReports)
	for i := 0; i < numChainReports; i++ {
		reportMessages := make([]cciptypes.Message, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			key, err := solanago.NewRandomPrivateKey()
			if err != nil {
				panic(err)
			}
			extraData, err := cciptypes.NewBytesFromString("0x1234")
			require.NoError(t, err)

			destGasAmount := uint32(10)
			destExecData := make([]byte, 4)
			binary.LittleEndian.PutUint32(destExecData, destGasAmount)

			tokenAmounts := make([]cciptypes.RampTokenAmount, numTokensPerMsg)
			for z := 0; z < numTokensPerMsg; z++ {
				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: cciptypes.UnknownAddress(key.PublicKey().String()),
					DestTokenAddress:  key.PublicKey().Bytes(),
					ExtraData:         extraData,
					Amount:            cciptypes.NewBigInt(big.NewInt(rand.Int63())),
					DestExecData:      destExecData,
					DestExecDataDecoded: map[string]any{
						"destGasAmount": uint32(10),
					},
				}
			}

			extraArgs := ccip_router.SVMExtraArgs{
				ComputeUnits:     1000,
				IsWritableBitmap: 2,
			}

			extraArgsMap := map[string]any{
				"ComputeUnits":            uint32(1000),
				"accountIsWritableBitmap": uint64(2),
			}

			var buf bytes.Buffer
			encoder := agbinary.NewBorshEncoder(&buf)
			err = extraArgs.MarshalWithEncoder(encoder)
			require.NoError(t, err)

			reportMessages[j] = cciptypes.Message{
				Header: cciptypes.RampMessageHeader{
					MessageID:           utils.RandomBytes32(),
					SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
					DestChainSelector:   cciptypes.ChainSelector(rand.Uint64()),
					SequenceNumber:      cciptypes.SeqNum(rand.Uint64()),
					Nonce:               rand.Uint64(),
					MsgHash:             utils.RandomBytes32(),
					OnRamp:              cciptypes.UnknownAddress(key.PublicKey().String()),
				},
				Sender:           cciptypes.UnknownAddress(key.PublicKey().String()),
				Data:             extraData,
				Receiver:         key.PublicKey().Bytes(),
				ExtraArgs:        buf.Bytes(),
				FeeToken:         cciptypes.UnknownAddress(key.PublicKey().String()),
				FeeTokenAmount:   cciptypes.NewBigInt(big.NewInt(rand.Int63())),
				TokenAmounts:     tokenAmounts,
				ExtraArgsDecoded: extraArgsMap,
			}
		}

		tokenData := make([][][]byte, numTokensPerMsg)
		for j := 0; j < numTokensPerMsg; j++ {
			tokenData[j] = [][]byte{{0x1}, {0x2, 0x3}}
		}

		chainReports[i] = cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(sourceChainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   tokenData,
			Proofs:              []cciptypes.Bytes32{utils.RandomBytes32(), utils.RandomBytes32()},
		}
	}

	return cciptypes.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1(t *testing.T) {
	testCases := []struct {
		name          string
		report        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport
		expErr        bool
		chainSelector uint64
	}{
		{
			name:          "base report with Solana as source chain",
			report:        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
		{
			name:          "base report with EVM as source chain",
			report:        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
		},
		// TODO: check if empty msg if necessary since there is only single msg in solana execute report
		// {
		//	 name: "reports have empty msgs",
		//	 report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
		//		 report.ChainReports[0].Messages = []cciptypes.Message{}
		//		 return report
		//	 },
		//	 expErr: false,
		// },
		{
			name: "reports have empty offchain token data",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].OffchainTokenData = [][][]byte{}
				return report
			},
			expErr:        false,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
	}

	ctx := testutils.Context(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cd := NewExecutePluginCodecV1()
			report := tc.report(randomExecuteReport(t, tc.chainSelector))
			bytes, err := cd.Encode(ctx, report)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// ignore msg hash, extraArgsDecoded map and DestExecDataDecoded map in comparison
			for i := range report.ChainReports {
				for j := range report.ChainReports[i].Messages {
					report.ChainReports[i].Messages[j].Header.MsgHash = cciptypes.Bytes32{}
					report.ChainReports[i].Messages[j].Header.OnRamp = cciptypes.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeToken = cciptypes.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeTokenAmount = cciptypes.BigInt{}
					report.ChainReports[i].Messages[j].ExtraArgsDecoded = nil
					for k := range report.ChainReports[i].Messages[j].TokenAmounts {
						report.ChainReports[i].Messages[j].TokenAmounts[k].DestExecDataDecoded = nil
					}
				}
			}

			// decode using the codec
			codecDecoded, err := cd.Decode(ctx, bytes)
			require.NoError(t, err)
			assert.Equal(t, report, codecDecoded)
		})
	}
}
