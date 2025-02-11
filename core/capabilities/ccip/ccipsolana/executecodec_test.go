package ccipsolana

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"math/rand"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

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

			extraArgs := ccip_offramp.Any2SVMRampExtraArgs{
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

func Test_DecodingExecuteReport(t *testing.T) {
	t.Run("decode on-chain execute report", func(t *testing.T) {
		chainSel := cciptypes.ChainSelector(rand.Uint64())
		onRampAddr, err := solanago.NewRandomPrivateKey()
		require.NoError(t, err)

		destGasAmount := uint32(10)
		tokenAmount := big.NewInt(rand.Int63())
		tokenReceiver := solanago.MustPublicKeyFromBase58("C8WSPj3yyus1YN3yNB6YA5zStYtbjQWtpmKadmvyUXq8")
		extraArgs := ccip_offramp.Any2SVMRampExtraArgs{
			ComputeUnits:     1000,
			IsWritableBitmap: 2,
		}

		onChainReport := ccip_offramp.ExecutionReportSingleChain{
			SourceChainSelector: uint64(chainSel),
			Message: ccip_offramp.Any2SVMRampMessage{
				Header: ccip_offramp.RampMessageHeader{
					SourceChainSelector: uint64(chainSel),
				},
				TokenReceiver: tokenReceiver,
				ExtraArgs:     extraArgs,
				TokenAmounts: []ccip_offramp.Any2SVMTokenTransfer{
					{
						Amount:        ccip_offramp.CrossChainAmount{LeBytes: [32]uint8(encodeBigIntToFixedLengthLE(tokenAmount, 32))},
						DestGasAmount: destGasAmount,
					},
				},
				OnRampAddress: onRampAddr.PublicKey().Bytes(),
			},
		}

		var extraArgsBuf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&extraArgsBuf)
		err = extraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		var buf bytes.Buffer
		encoder = agbinary.NewBorshEncoder(&buf)
		err = onChainReport.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		executeCodec := NewExecutePluginCodecV1()
		decode, err := executeCodec.Decode(testutils.Context(t), buf.Bytes())
		require.NoError(t, err)

		report := decode.ChainReports[0]
		require.Equal(t, chainSel, report.SourceChainSelector)

		msg := report.Messages[0]
		require.Equal(t, cciptypes.UnknownAddress(tokenReceiver.Bytes()), msg.Receiver)
		require.Equal(t, cciptypes.Bytes(extraArgsBuf.Bytes()), msg.ExtraArgs)
		require.Equal(t, tokenAmount, msg.TokenAmounts[0].Amount.Int)
		require.Equal(t, destGasAmount, bytesToUint32LE(msg.TokenAmounts[0].DestExecData))
	})

	t.Run("decode Borsh encoded execute report", func(t *testing.T) {
		ocrReport := randomExecuteReport(t, 124615329519749607)
		cd := NewExecutePluginCodecV1()
		encodedReport, err := cd.Encode(testutils.Context(t), ocrReport)
		require.NoError(t, err)

		decoder := agbinary.NewBorshDecoder(encodedReport)
		executeReport := ccip_offramp.ExecutionReportSingleChain{}
		err = executeReport.UnmarshalWithDecoder(decoder)
		require.NoError(t, err)

		originReport := ocrReport.ChainReports[0]
		require.Equal(t, originReport.SourceChainSelector, cciptypes.ChainSelector(executeReport.SourceChainSelector))

		originMsg := originReport.Messages[0]
		require.Equal(t, originMsg.Header.MessageID, cciptypes.Bytes32(executeReport.Message.Header.MessageId))
		require.Equal(t, originMsg.Header.DestChainSelector, cciptypes.ChainSelector(executeReport.Message.Header.DestChainSelector))
		require.Equal(t, originMsg.Header.SourceChainSelector, cciptypes.ChainSelector(executeReport.Message.Header.SourceChainSelector))

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err = executeReport.Message.ExtraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)
		require.Equal(t, originMsg.ExtraArgs, cciptypes.Bytes(buf.Bytes()))

		originTokenAmount := originMsg.TokenAmounts[0]
		require.Equal(t, originTokenAmount.Amount, decodeLEToBigInt(executeReport.Message.TokenAmounts[0].Amount.LeBytes[:]))
		require.Equal(t, originTokenAmount.DestTokenAddress, cciptypes.UnknownAddress(executeReport.Message.TokenAmounts[0].DestTokenAddress.Bytes()))
		require.Equal(t, bytesToUint32LE(originTokenAmount.DestExecData), executeReport.Message.TokenAmounts[0].DestGasAmount)
		require.Equal(t, originMsg.Sender, cciptypes.UnknownAddress(executeReport.Message.Sender))
	})
}
