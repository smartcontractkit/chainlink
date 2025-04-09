package ccipevm

import (
	"encoding/base64"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/report_codec"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var randomExecuteReport = func(t *testing.T, d *testSetupData, chainSelector uint64, gasLimit *big.Int, destGasAmount uint32) cciptypes.ExecutePluginReport {
	const numChainReports = 10
	const msgsPerReport = 10
	const numTokensPerMsg = 3

	chainReports := make([]cciptypes.ExecutePluginReportSingleChain, numChainReports)
	for i := 0; i < numChainReports; i++ {
		reportMessages := make([]cciptypes.Message, msgsPerReport)
		for j := 0; j < msgsPerReport; j++ {
			data, err := cciptypes.NewBytesFromString(utils.RandomAddress().String())
			assert.NoError(t, err)

			tokenAmounts := make([]cciptypes.RampTokenAmount, numTokensPerMsg)
			for z := 0; z < numTokensPerMsg; z++ {
				encodedDestExecData, err2 := abiEncodeUint32(destGasAmount)
				require.NoError(t, err2)

				tokenAmounts[z] = cciptypes.RampTokenAmount{
					SourcePoolAddress: utils.RandomAddress().Bytes(),
					DestTokenAddress:  utils.RandomAddress().Bytes(),
					ExtraData:         data,
					Amount:            cciptypes.NewBigInt(utils.RandUint256()),
					DestExecData:      encodedDestExecData,
				}
			}

			extraArgs, err := d.contract.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
				GasLimit: gasLimit,
			})
			assert.NoError(t, err)

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
				Receiver:       utils.RandomAddress().Bytes(),
				ExtraArgs:      extraArgs,
				FeeToken:       utils.RandomAddress().Bytes(),
				FeeTokenAmount: cciptypes.NewBigInt(utils.RandUint256()),
				TokenAmounts:   tokenAmounts,
			}
		}

		tokenData := make([][][]byte, numTokensPerMsg)
		for j := 0; j < numTokensPerMsg; j++ {
			tokenData[j] = [][]byte{{0x1}, {0x2, 0x3}}
		}

		chainReports[i] = cciptypes.ExecutePluginReportSingleChain{
			SourceChainSelector: cciptypes.ChainSelector(chainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   tokenData,
			Proofs:              []cciptypes.Bytes32{utils.RandomBytes32(), utils.RandomBytes32()},
			ProofFlagBits:       cciptypes.NewBigInt(utils.RandUint256()),
		}
	}

	return cciptypes.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1(t *testing.T) {
	d := testSetup(t)
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
				report.ChainReports[4].Messages = []cciptypes.Message{}
				return report
			},
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name: "reports have empty offchain token data",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].OffchainTokenData = [][][]byte{}
				report.ChainReports[4].OffchainTokenData[1] = [][]byte{}
				return report
			},
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name: "reports have negative token amounts",
			report: func(report cciptypes.ExecutePluginReport) cciptypes.ExecutePluginReport {
				report.ChainReports[0].Messages[0].TokenAmounts[0].Amount = cciptypes.NewBigInt(big.NewInt(-1))
				return report
			},
			expErr:        true,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
	}

	// Deploy the contract
	transactor := evmtestutils.MustNewSimTransactor(t)
	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)
	address, _, _, err := report_codec.DeployReportCodec(transactor, simulatedBackend)
	require.NoError(t, err)
	simulatedBackend.Commit()
	contract, err := report_codec.NewReportCodec(address, simulatedBackend)
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			codec := NewExecutePluginCodecV1(mockExtraDataCodec)
			report := tc.report(randomExecuteReport(t, d, tc.chainSelector, tc.gasLimit, tc.destGasAmount))
			bytes, err := codec.Encode(ctx, report)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			testSetup(t)

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

			// decode using the contract
			contractDecodedReport, err := contract.DecodeExecuteReport(&bind.CallOpts{Context: ctx}, bytes)
			assert.NoError(t, err)
			assert.Equal(t, len(report.ChainReports), len(contractDecodedReport))
			for i, expReport := range report.ChainReports {
				actReport := contractDecodedReport[i]
				assert.Equal(t, expReport.OffchainTokenData, actReport.OffchainTokenData)
				assert.Equal(t, len(expReport.Messages), len(actReport.Messages))
				assert.Equal(t, uint64(expReport.SourceChainSelector), actReport.SourceChainSelector)
			}

			// decode using the codec
			codecDecoded, err := codec.Decode(ctx, bytes)
			assert.NoError(t, err)
			assert.Equal(t, report, codecDecoded)
		})
	}
}

func Test_DecodeReport(t *testing.T) {
	ExtraDataCodec := ccipcommon.NewExtraDataCodec(ccipcommon.NewExtraDataCodecParams(ExtraDataDecoder{}, ccipsolana.ExtraDataDecoder{}))
	offRampABI, err := offramp.OffRampMetaData.GetAbi()
	require.NoError(t, err)

	reportBase64 := "9Y4D/AAKbBOGy6NAcrBVaUmVfONiiic4CO6GbHwProoOgQEuAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATQECAwyzhPUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAPgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAx7Gv0ioZNqUJw0gfsLHFIQQr3lR0XlsvWXBIfQJNfOgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE0BAgMMs4T1AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADJ+ShEYchSsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABgAAAAAAAAAAAAAAAANczLkN32zc8KqgZj7Tm8KueHiruAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADDUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABoAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAZAUE4xFSzvn6m/FNtkX62lAsPigAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAAAAAAAAAAAAAIn5uKZ/XwqWg+aAsi9+uE03oJuIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB6EgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA3gtrOnZAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABTDIjV28BIw1+6agsDfXqqEiuKGowAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	reportBytes, err := base64.StdEncoding.DecodeString(reportBase64)
	require.NoError(t, err)

	executeInputs, err := offRampABI.Methods["execute"].Inputs.Unpack(reportBytes[4:])
	require.NoError(t, err)
	require.Len(t, executeInputs, 2)

	// first param is report ctx, which is bytes32[2], so cast to that using
	// abi.ConvertType
	reportCtx := *abi.ConvertType(executeInputs[0], new([2][32]byte)).(*[2][32]byte)
	t.Logf("reportCtx[0]: %x, reportCtx[1]: %x", reportCtx[0], reportCtx[1])

	rawReport := *abi.ConvertType(executeInputs[1], new([]byte)).(*[]byte)

	codec := NewExecutePluginCodecV1(ExtraDataCodec)
	decoded, err := codec.Decode(t.Context(), rawReport)
	require.NoError(t, err)

	t.Logf("decoded: %+v", decoded)
}
