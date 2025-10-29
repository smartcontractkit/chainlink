package ccipsolana

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"math/rand"
	"testing"

	agbinary "github.com/gagliardetto/binary"
	solanago "github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/latest/ccip_offramp"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var randomExecuteReport = func(t *testing.T, sourceChainSelector uint64) ccipocr3.ExecutePluginReport {
	const numChainReports = 1
	const msgsPerReport = 1
	const numTokensPerMsg = 1

	chainReports := make([]ccipocr3.ExecutePluginReportSingleChain, numChainReports)
	for i := range numChainReports {
		reportMessages := make([]ccipocr3.Message, msgsPerReport)
		for j := range msgsPerReport {
			key, err := solanago.NewRandomPrivateKey()
			if err != nil {
				panic(err)
			}
			extraData, err := ccipocr3.NewBytesFromString("0x1234")
			require.NoError(t, err)
			tokenReceiver := solanago.MustPublicKeyFromBase58("42Gia5bGsh8R2S44e37t9fsucap1qsgjr6GjBmWotgdF")
			destGasAmount := uint32(10)
			destExecData := make([]byte, 4)
			binary.LittleEndian.PutUint32(destExecData, destGasAmount)

			tokenAmounts := make([]ccipocr3.RampTokenAmount, numTokensPerMsg)
			for z := range numTokensPerMsg {
				tokenAmounts[z] = ccipocr3.RampTokenAmount{
					SourcePoolAddress: ccipocr3.UnknownAddress(key.PublicKey().String()),
					DestTokenAddress:  key.PublicKey().Bytes(),
					ExtraData:         extraData,
					Amount:            ccipocr3.NewBigInt(big.NewInt(rand.Int63())),
					DestExecData:      destExecData,
				}
			}

			extraArgs := ccip_offramp.Any2SVMRampExtraArgs{
				ComputeUnits:     1000,
				IsWritableBitmap: 2,
			}

			var buf bytes.Buffer
			encoder := agbinary.NewBorshEncoder(&buf)
			err = extraArgs.MarshalWithEncoder(encoder)
			require.NoError(t, err)

			reportMessages[j] = ccipocr3.Message{
				Header: ccipocr3.RampMessageHeader{
					MessageID:           utils.RandomBytes32(),
					SourceChainSelector: ccipocr3.ChainSelector(sourceChainSelector),
					DestChainSelector:   ccipocr3.ChainSelector(rand.Uint64()),
					SequenceNumber:      ccipocr3.SeqNum(rand.Uint64()),
					Nonce:               rand.Uint64(),
					MsgHash:             utils.RandomBytes32(),
					OnRamp:              ccipocr3.UnknownAddress(key.PublicKey().String()),
				},
				Sender:         ccipocr3.UnknownAddress(key.PublicKey().String()),
				Data:           extraData,
				Receiver:       tokenReceiver.Bytes(),
				ExtraArgs:      buf.Bytes(),
				FeeToken:       ccipocr3.UnknownAddress(key.PublicKey().String()),
				FeeTokenAmount: ccipocr3.NewBigInt(big.NewInt(rand.Int63())),
				TokenAmounts:   tokenAmounts,
			}
		}

		tokenData := make([][][]byte, numTokensPerMsg)
		for j := range numTokensPerMsg {
			tokenData[j] = [][]byte{{0x1}, {0x2, 0x3}}
		}

		chainReports[i] = ccipocr3.ExecutePluginReportSingleChain{
			SourceChainSelector: ccipocr3.ChainSelector(sourceChainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   tokenData,
			Proofs:              []ccipocr3.Bytes32{utils.RandomBytes32(), utils.RandomBytes32()},
		}
	}

	return ccipocr3.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1(t *testing.T) {
	testCases := []struct {
		name          string
		report        func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport
		expErr        bool
		chainSelector uint64
	}{
		{
			name:          "base report with Solana as source chain",
			report:        func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
		{
			name:          "base report with EVM as source chain",
			report:        func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
		},
		// TODO: check if empty msg if necessary since there is only single msg in solana execute report
		// {
		//	 name: "reports have empty msgs",
		//	 report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
		//		 report.ChainReports[0].Messages = []ccipocr3.Message{}
		//		 return report
		//	 },
		//	 expErr: false,
		// },
		{
			name: "reports have empty offchain token data",
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
				report.ChainReports[0].OffchainTokenData = [][][]byte{}
				return report
			},
			expErr:        false,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
		{
			name: "reports have invalid DestTokenAddress",
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
				report.ChainReports[0].Messages[0].TokenAmounts[0].DestTokenAddress = []byte{0, 0}
				return report
			},
			expErr:        true,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
		{
			name: "reports have invalid receiver",
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
				report.ChainReports[0].Messages[0].Receiver = []byte{0, 0}
				return report
			},
			expErr:        true,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
		{
			name: "reports have negative token amount",
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
				report.ChainReports[0].Messages[0].TokenAmounts[0].Amount = ccipocr3.NewBigInt(big.NewInt(-1))
				return report
			},
			expErr:        true,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
		},
	}

	ctx := testutils.Context(t)
	mockExtraDataCodec := mocks.NewSourceChainExtraDataCodec(t)
	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything).Return(map[string]any{
		"destGasAmount": uint32(10),
	}, nil).Maybe()
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything).Return(map[string]any{
		"ComputeUnits":            uint32(1000),
		"accountIsWritableBitmap": uint64(2),
		"TokenReceiver":           [32]byte(solanago.MustPublicKeyFromBase58("42Gia5bGsh8R2S44e37t9fsucap1qsgjr6GjBmWotgdF").Bytes()),
	}, nil).Maybe()
	registeredMockExtraDataCodecMap := map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
	}

	edc := ccipocr3.ExtraDataCodecMap(registeredMockExtraDataCodecMap)
	cd := NewExecutePluginCodecV1(edc)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
					report.ChainReports[i].Messages[j].Header.MsgHash = ccipocr3.Bytes32{}
					report.ChainReports[i].Messages[j].Header.OnRamp = ccipocr3.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeToken = ccipocr3.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeTokenAmount = ccipocr3.BigInt{}
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
	mockExtraDataCodec := mocks.NewSourceChainExtraDataCodec(t)
	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything, mock.Anything).Return(map[string]any{
		"destGasAmount": uint32(10),
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything, mock.Anything).Return(map[string]any{
		"ComputeUnits":            uint32(1000),
		"accountIsWritableBitmap": uint64(2),
	}, nil)
	registeredMockExtraDataCodecMap := map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
	}

	t.Run("decode on-chain execute report", func(t *testing.T) {
		chainSel := ccipocr3.ChainSelector(rand.Uint64())

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
			},
		}

		var extraArgsBuf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&extraArgsBuf)
		err := extraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		var buf bytes.Buffer
		encoder = agbinary.NewBorshEncoder(&buf)
		err = onChainReport.MarshalWithEncoder(encoder)
		require.NoError(t, err)

		edc := ccipocr3.ExtraDataCodecMap(registeredMockExtraDataCodecMap)
		executeCodec := NewExecutePluginCodecV1(edc)
		decode, err := executeCodec.Decode(testutils.Context(t), buf.Bytes())
		require.NoError(t, err)

		report := decode.ChainReports[0]
		require.Equal(t, chainSel, report.SourceChainSelector)

		msg := report.Messages[0]
		require.Equal(t, ccipocr3.UnknownAddress(tokenReceiver.Bytes()), msg.Receiver)
		require.Equal(t, ccipocr3.Bytes(extraArgsBuf.Bytes()), msg.ExtraArgs)
		require.Equal(t, tokenAmount, msg.TokenAmounts[0].Amount.Int)
		require.Equal(t, destGasAmount, binary.LittleEndian.Uint32(msg.TokenAmounts[0].DestExecData))
	})

	t.Run("decode Borsh encoded execute report", func(t *testing.T) {
		ocrReport := randomExecuteReport(t, 124615329519749607)
		edc := ccipocr3.ExtraDataCodecMap(registeredMockExtraDataCodecMap)
		cd := NewExecutePluginCodecV1(edc)
		encodedReport, err := cd.Encode(testutils.Context(t), ocrReport)
		require.NoError(t, err)

		decoder := agbinary.NewBorshDecoder(encodedReport)
		executeReport := ccip_offramp.ExecutionReportSingleChain{}
		err = executeReport.UnmarshalWithDecoder(decoder)
		require.NoError(t, err)

		originReport := ocrReport.ChainReports[0]
		require.Equal(t, originReport.SourceChainSelector, ccipocr3.ChainSelector(executeReport.SourceChainSelector))

		originMsg := originReport.Messages[0]
		require.Equal(t, originMsg.Header.MessageID, ccipocr3.Bytes32(executeReport.Message.Header.MessageId))
		require.Equal(t, originMsg.Header.DestChainSelector, ccipocr3.ChainSelector(executeReport.Message.Header.DestChainSelector))
		require.Equal(t, originMsg.Header.SourceChainSelector, ccipocr3.ChainSelector(executeReport.Message.Header.SourceChainSelector))

		var buf bytes.Buffer
		encoder := agbinary.NewBorshEncoder(&buf)
		err = executeReport.Message.ExtraArgs.MarshalWithEncoder(encoder)
		require.NoError(t, err)
		require.Equal(t, originMsg.ExtraArgs, ccipocr3.Bytes(buf.Bytes()))

		originTokenAmount := originMsg.TokenAmounts[0]
		require.Equal(t, originTokenAmount.Amount, decodeLEToBigInt(executeReport.Message.TokenAmounts[0].Amount.LeBytes[:]))
		require.Equal(t, originTokenAmount.DestTokenAddress, ccipocr3.UnknownAddress(executeReport.Message.TokenAmounts[0].DestTokenAddress.Bytes()))
		require.Equal(t, binary.LittleEndian.Uint32(originTokenAmount.DestExecData), executeReport.Message.TokenAmounts[0].DestGasAmount)
		require.Equal(t, originMsg.Sender, ccipocr3.UnknownAddress(executeReport.Message.Sender))
	})
}
