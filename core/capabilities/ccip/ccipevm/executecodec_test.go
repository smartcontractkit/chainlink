package ccipevm

import (
	"encoding/base64"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/report_codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common/mocks"
)

var randomExecuteReport = func(t *testing.T, d *testSetupData, chainSelector uint64, gasLimit *big.Int, destGasAmount uint32) ccipocr3.ExecutePluginReport {
	const numChainReports = 10
	const msgsPerReport = 10
	const numTokensPerMsg = 3

	chainReports := make([]ccipocr3.ExecutePluginReportSingleChain, numChainReports)
	for i := range numChainReports {
		reportMessages := make([]ccipocr3.Message, msgsPerReport)
		for j := range msgsPerReport {
			data, err := ccipocr3.NewBytesFromString(utils.RandomAddress().String())
			require.NoError(t, err)

			tokenAmounts := make([]ccipocr3.RampTokenAmount, numTokensPerMsg)
			for z := range numTokensPerMsg {
				encodedDestExecData, err2 := abiEncodeUint32(destGasAmount)
				require.NoError(t, err2)

				tokenAmounts[z] = ccipocr3.RampTokenAmount{
					SourcePoolAddress: utils.RandomAddress().Bytes(),
					DestTokenAddress:  utils.RandomAddress().Bytes(),
					ExtraData:         data,
					Amount:            ccipocr3.NewBigInt(utils.RandUint256()),
					DestExecData:      encodedDestExecData,
				}
			}

			extraArgs, err := d.contract.EncodeEVMExtraArgsV1(nil, message_hasher.ClientEVMExtraArgsV1{
				GasLimit: gasLimit,
			})
			require.NoError(t, err)

			reportMessages[j] = ccipocr3.Message{
				Header: ccipocr3.RampMessageHeader{
					MessageID:           utils.RandomBytes32(),
					SourceChainSelector: ccipocr3.ChainSelector(rand.Uint64()),
					DestChainSelector:   ccipocr3.ChainSelector(rand.Uint64()),
					SequenceNumber:      ccipocr3.SeqNum(rand.Uint64()),
					Nonce:               rand.Uint64(),
					MsgHash:             utils.RandomBytes32(),
					OnRamp:              utils.RandomAddress().Bytes(),
				},
				Sender:         common.LeftPadBytes(utils.RandomAddress().Bytes(), 32),
				Data:           data,
				Receiver:       utils.RandomAddress().Bytes(),
				ExtraArgs:      extraArgs,
				FeeToken:       utils.RandomAddress().Bytes(),
				FeeTokenAmount: ccipocr3.NewBigInt(utils.RandUint256()),
				TokenAmounts:   tokenAmounts,
			}
		}

		tokenData := make([][][]byte, numTokensPerMsg)
		for j := range numTokensPerMsg {
			tokenData[j] = [][]byte{{0x1}, {0x2, 0x3}}
		}

		chainReports[i] = ccipocr3.ExecutePluginReportSingleChain{
			SourceChainSelector: ccipocr3.ChainSelector(chainSelector),
			Messages:            reportMessages,
			OffchainTokenData:   tokenData,
			Proofs:              []ccipocr3.Bytes32{utils.RandomBytes32(), utils.RandomBytes32()},
			ProofFlagBits:       ccipocr3.NewBigInt(utils.RandUint256()),
		}
	}

	return ccipocr3.ExecutePluginReport{ChainReports: chainReports}
}

func TestExecutePluginCodecV1(t *testing.T) {
	d := testSetup(t)
	ctx := t.Context()
	mockExtraDataCodec := mocks.NewSourceChainExtraDataCodec(t)
	destGasAmount := rand.Uint32()
	gasLimit := utils.RandUint256()
	mockExtraDataCodec.On("DecodeDestExecDataToMap", mock.Anything).Return(map[string]any{
		"destgasamount": destGasAmount,
	}, nil)
	mockExtraDataCodec.On("DecodeExtraArgsToMap", mock.Anything).Return(map[string]any{
		"gasLimit":                utils.RandUint256(),
		"accountIsWritableBitmap": gasLimit,
	}, nil)

	testCases := []struct {
		name          string
		report        func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport
		expErr        bool
		chainSelector uint64
		destGasAmount uint32
		gasLimit      *big.Int
	}{
		{
			name:          "base report",
			report:        func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name:          "base report",
			report:        func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport { return report },
			expErr:        false,
			chainSelector: 124615329519749607, // Solana mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name: "reports have empty msgs",
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
				report.ChainReports[0].Messages = []ccipocr3.Message{}
				report.ChainReports[4].Messages = []ccipocr3.Message{}
				return report
			},
			expErr:        false,
			chainSelector: 5009297550715157269, // ETH mainnet chain selector
			gasLimit:      gasLimit,
			destGasAmount: destGasAmount,
		},
		{
			name: "reports have empty offchain token data",
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
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
			report: func(report ccipocr3.ExecutePluginReport) ccipocr3.ExecutePluginReport {
				report.ChainReports[0].Messages[0].TokenAmounts[0].Amount = ccipocr3.NewBigInt(big.NewInt(-1))
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
	b := simulated.NewBackend(types.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, simulated.WithBlockGasLimit(30e6), func(_ *node.Config, ethCfg *ethconfig.Config) {
		ethCfg.RPCEVMTimeout = 60 * time.Second
	})
	simulatedBackend := &backends.SimulatedBackend{Backend: b, Client: b.Client()}
	address, _, _, err := report_codec.DeployReportCodec(transactor, simulatedBackend)
	require.NoError(t, err)
	simulatedBackend.Commit()
	contract, err := report_codec.NewReportCodec(address, simulatedBackend)
	require.NoError(t, err)
	registeredMockExtraDataCodecMap := map[string]ccipocr3.SourceChainExtraDataCodec{
		chainsel.FamilyEVM:    mockExtraDataCodec,
		chainsel.FamilySolana: mockExtraDataCodec,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			edc := ccipocr3.ExtraDataCodecMap(registeredMockExtraDataCodecMap)
			codec := NewExecutePluginCodecV1(edc)
			report := tc.report(randomExecuteReport(t, d, tc.chainSelector, tc.gasLimit, tc.destGasAmount))
			bytes, err := codec.Encode(ctx, report)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			testSetup(t)

			// ignore msg hash in comparison
			for i := range report.ChainReports {
				for j := range report.ChainReports[i].Messages {
					report.ChainReports[i].Messages[j].Header.MsgHash = ccipocr3.Bytes32{}
					report.ChainReports[i].Messages[j].Header.OnRamp = ccipocr3.UnknownAddress{}
					report.ChainReports[i].Messages[j].FeeToken = ccipocr3.UnknownAddress{}
					report.ChainReports[i].Messages[j].ExtraArgs = ccipocr3.Bytes{}
					report.ChainReports[i].Messages[j].FeeTokenAmount = ccipocr3.BigInt{}
				}
			}

			// decode using the contract
			contractDecodedReport, err := contract.DecodeExecuteReport(&bind.CallOpts{Context: ctx}, bytes)
			require.NoError(t, err)
			require.Len(t, contractDecodedReport, len(report.ChainReports))
			for i, expReport := range report.ChainReports {
				actReport := contractDecodedReport[i]
				assert.Equal(t, expReport.OffchainTokenData, actReport.OffchainTokenData)
				assert.Len(t, actReport.Messages, len(expReport.Messages))
				assert.Equal(t, uint64(expReport.SourceChainSelector), actReport.SourceChainSelector)
			}

			// decode using the codec
			codecDecoded, err := codec.Decode(ctx, bytes)
			require.NoError(t, err)
			assert.Equal(t, report, codecDecoded)
		})
	}
}

func Test_DecodeReport(t *testing.T) {
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
	codec := NewExecutePluginCodecV1(extraDataCodec)
	decoded, err := codec.Decode(t.Context(), rawReport)
	require.NoError(t, err)

	t.Logf("decoded: %+v", decoded)
}
