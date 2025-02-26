package ccipaptos

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_aptos_utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var randomExecuteReport = func(t *testing.T, chainSelector uint64, gasLimit *big.Int, destGasAmount uint32) cciptypes.ExecutePluginReport {
	const numChainReports = 1
	const msgsPerReport = 1
	const numTokensPerMsg = 3

	chainReports := make([]cciptypes.ExecutePluginReportSingleChain, numChainReports)
	for i := 0; i < numChainReports; i++ {
		reportMessages := make([]cciptypes.Message, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			data, err := cciptypes.NewBytesFromString(utils.RandomAddress().String())
			require.NoError(t, err)

			tokenAmounts := make([]cciptypes.RampTokenAmount, numTokensPerMsg)
			for z := 0; z < numTokensPerMsg; z++ {
				encodedDestExecData, err2 := abiEncodeUint32(destGasAmount)
				require.NoError(t, err2)

				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: utils.RandomAddress().Bytes(),
					DestTokenAddress:  generateAddressBytes(),
					ExtraData:         data,
					Amount:            cciptypes.NewBigInt(utils.RandUint256()),
					DestExecData:      encodedDestExecData,
				}
			}

			extraArgs, err := aptosUtilsABI.Pack("exposeEVMExtraArgsV1", ccip_aptos_utils.AptosUtilsEVMExtraArgsV1{
				GasLimit: gasLimit,
			})
			if err != nil {
				t.Fatalf("failed to pack extra args: %v", err)
			}
			extraArgs = append(evmExtraArgsV1Tag, extraArgs[4:]...)

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
				Sender:         common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				Data:           data,
				Receiver:       generateAddressBytes(),
				ExtraArgs:      extraArgs,
				FeeToken:       generateAddressBytes(),
				FeeTokenAmount: cciptypes.NewBigInt(utils.RandUint256()),
				TokenAmounts:   tokenAmounts,
			}
		}

		tokenData := make([][][]byte, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			tokenData[j] = [][]byte{{0x1}, {0x2, 0x3}}
		}

		chainReports[i] = cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(chainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   tokenData,
			Proofs:              []cciptypes.Bytes32{utils.RandomBytes32(), utils.RandomBytes32()},
			ProofFlagBits:       cciptypes.NewBigInt(big.NewInt(0)),
		}
	}

	return cciptypes.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1(t *testing.T) {
	ctx := testutils.Context(t)
	mockExtraDataCodec := &mocks.ExtraDataCodec{}
	destGasAmount := rand.Uint32()
	gasLimit := utils.RandUint256()
	mockExtraDataCodec.On("DecodeTokenAmountDestExecData", mock.Anything, mock.Anything).Return(map[string]any{
		"destgasamount": destGasAmount,
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgs", mock.Anything, mock.Anything).Return(map[string]any{
		"gasLimit":                utils.RandUint256(),
		"accountIsWritableBitmap": gasLimit,
	}, nil)

	testCases := []struct {
		name          string
		report        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport
		expErr        bool
		chainSelector uint64
		destGasAmount uint32
		gasLimit      *big.Int
	}{
		{
			name:          "base report",
			report:        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name:          "base report",
			report:        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name: "reports have empty msgs",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].Messages = []cciptypes.Message{}
				return report
			},
			expErr:        true,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name: "reports have empty offchain token data",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].OffchainTokenData = [][][]byte{}
				return report
			},
			expErr:        true,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewExecutePluginCodecV1(mockExtraDataCodec)
			report := tc.report(randomExecuteReport(t, tc.chainSelector, tc.gasLimit, tc.destGasAmount))
			bytes, err := codec.Encode(ctx, report)
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// ignore unavailable fields in comparison
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
			codecDecoded, err := codec.Decode(ctx, bytes)
			require.NoError(t, err)
			require.Equal(t, report, codecDecoded)
		})
	}
}
