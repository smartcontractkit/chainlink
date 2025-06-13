package ccipaptos

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
				// Use BCS to pack destGasAmount
				encodedDestExecData, err2 := bcs.SerializeU32(destGasAmount)
				require.NoError(t, err2)

				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: utils.RandomAddress().Bytes(),
					DestTokenAddress:  generateAddressBytes(),
					ExtraData:         data,
					Amount:            cciptypes.NewBigInt(utils.RandUint256()),
					DestExecData:      encodedDestExecData,
				}
			}

			// Use BCS to pack EVM V1 fields
			encodedExtraArgsFields, err := bcs.SerializeU256(*gasLimit)
			require.NoError(t, err, "failed to pack extra args fields")

			// Prepend the tag
			extraArgs := append(evmExtraArgsV1Tag, encodedExtraArgsFields...)

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
	mockExtraDataCodec := &mocks.SourceChainExtraDataCodec{}
	destGasAmount := rand.Uint32()
	gasLimit := utils.RandUint256()

	// Update mock return values to use the correct keys expected by the codec
	// The codec uses the ExtraDataDecoder internally, which returns maps like these.
	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything, mock.Anything).Return(map[string]any{
		aptosDestExecDataKey: destGasAmount, // Use the constant defined in the decoder
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything, mock.Anything).Return(map[string]any{
		"gasLimit": gasLimit, // Match the key used in the decoder for EVM V1/V2 gasLimit
		// "allowOutOfOrderExecution": false, // Optionally mock other fields if needed by codec logic
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
			name:          "base report EVM chain",
			report:        func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name:          "base report non-EVM chain", // Name updated for clarity
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
			chainSelector: 5009297550715157269,
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
			chainSelector: 5009297550715157269,
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
	}

	registeredMockExtraDataCodecMap := map[string]ccipcommon.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
		chainsel.FamilyAptos:  mockExtraDataCodec,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewExecutePluginCodecV1(registeredMockExtraDataCodecMap)
			// randomExecuteReport now uses the new encoding internally
			report := tc.report(randomExecuteReport(t, tc.chainSelector, tc.gasLimit, tc.destGasAmount))
			bytes, err := codec.Encode(ctx, report)
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// ignore unavailable fields in comparison - This part remains the same
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
			require.Equal(t, report, codecDecoded) // Comparison should still work
		})
	}
}
