package evmliquiditymanager

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/report_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/testonlybridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestEvmReportCodec(t *testing.T) {
	t.Run("marshal onchain", func(t *testing.T) {
		bridgeData1 := testutils.Random32Byte()
		bridgeData2 := testutils.Random32Byte()
		rm := models.Report{
			NetworkID:               1,
			LiquidityManagerAddress: models.Address(testutils.NewAddress()),
			Transfers: []models.Transfer{
				// send instruction
				{
					From:               1,
					To:                 2,
					Amount:             ubig.NewI(3),
					Sender:             models.Address(testutils.NewAddress()),
					Receiver:           models.Address(testutils.NewAddress()),
					LocalTokenAddress:  models.Address(testutils.NewAddress()),
					RemoteTokenAddress: models.Address(testutils.NewAddress()),
					Date:               time.Now().UTC(),
					BridgeData:         bridgeData1[:],
					NativeBridgeFee:    ubig.NewI(4),
				},
				// receive instruction
				{
					From:               3,
					To:                 1,
					Amount:             ubig.NewI(5),
					Sender:             models.Address(testutils.NewAddress()),
					Receiver:           models.Address(testutils.NewAddress()),
					LocalTokenAddress:  models.Address(testutils.NewAddress()),
					RemoteTokenAddress: models.Address(testutils.NewAddress()),
					Date:               time.Now().UTC(),
					BridgeData:         bridgeData2[:],
					NativeBridgeFee:    ubig.NewI(6),
				},
			},
		}
		instructions, err := rm.ToLiquidityInstructions()
		require.NoError(t, err, "failed to convert ReportMetadata to LiquidityInstructions")

		var evmEncoder = NewEvmReportCodec()
		encoded, err := evmEncoder.Encode(rm)
		require.NoError(t, err, "failed to encode ReportMetadata")

		r, decodedInstructions, err := evmEncoder.Decode(rm.NetworkID, rm.LiquidityManagerAddress, encoded)
		require.NoError(t, err, "failed to unmarshal ReportMetadata")
		require.Equal(t, instructions, decodedInstructions, "marshalled and unmarshalled instructions should be equal")
		require.Equal(t, rm.NetworkID, r.NetworkID, "marshalled and unmarshalled NetworkID should be equal")
		require.Equal(t, rm.LiquidityManagerAddress, r.LiquidityManagerAddress, "marshalled and unmarshalled LiquidityManagerAddress should be equal")
		require.Equal(t, rm.Transfers[0].Amount, r.Transfers[0].Amount, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[0].From, r.Transfers[0].From, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[0].To, r.Transfers[0].To, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[0].BridgeData, r.Transfers[0].BridgeData, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].Amount, r.Transfers[1].Amount, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].From, r.Transfers[1].From, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].To, r.Transfers[1].To, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].BridgeData, r.Transfers[1].BridgeData, "marshalled and unmarshalled Transfers should be equal")
	})

	t.Run("unmarshal onchain", func(t *testing.T) {
		evmEncoder := NewEvmReportCodec()
		// an actual report from one integration test run
		// should consist of 1 send operation from 909606746561742123 to 3379446385462418246
		packedBridgeData, err := testonlybridge.PackBridgeSendReturnData(big.NewInt(1))
		require.NoError(t, err)
		encodedReport, err := evmEncoder.onchainReportArguments.Pack(report_encoder.ILiquidityManagerLiquidityInstructions{
			SendLiquidityParams: []report_encoder.ILiquidityManagerSendLiquidityParams{
				{
					Amount:              assets.Ether(5).ToInt(),
					NativeBridgeFee:     big.NewInt(0),
					RemoteChainSelector: 3379446385462418246,
					BridgeData:          packedBridgeData,
				},
			},
			ReceiveLiquidityParams: []report_encoder.ILiquidityManagerReceiveLiquidityParams{},
		})
		require.NoError(t, err)
		_, instructions, err := evmEncoder.Decode(
			909606746561742123,
			models.Address(common.HexToAddress("0x2033C546BC60900f8B765F0a8e7E2376a17cba5d")),
			encodedReport)
		require.NoError(t, err)
		require.Len(t, instructions.SendLiquidityParams, 1)
		require.Len(t, instructions.ReceiveLiquidityParams, 0)
		require.Equal(t, uint64(3379446385462418246), instructions.SendLiquidityParams[0].RemoteChainSelector)
		require.Equal(t, assets.Ether(5).ToInt(), instructions.SendLiquidityParams[0].Amount)
		require.Equal(t, big.NewInt(0).String(), instructions.SendLiquidityParams[0].NativeBridgeFee.String())
		amount, err := testonlybridge.UnpackBridgeSendReturnData(instructions.SendLiquidityParams[0].BridgeData)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(1).String(), amount.String())

		// should consist of 1 receive instruction from 909606746561742123 to 3379446385462418246
		packedBridgeData, err = testonlybridge.PackFinalizeBridgePayload(assets.Ether(5).ToInt(), big.NewInt(1))
		require.NoError(t, err)
		encodedReport, err = evmEncoder.onchainReportArguments.Pack(report_encoder.ILiquidityManagerLiquidityInstructions{
			SendLiquidityParams: []report_encoder.ILiquidityManagerSendLiquidityParams{},
			ReceiveLiquidityParams: []report_encoder.ILiquidityManagerReceiveLiquidityParams{
				{
					Amount:              assets.Ether(5).ToInt(),
					RemoteChainSelector: 909606746561742123,
					BridgeData:          packedBridgeData,
				},
			},
		})
		require.NoError(t, err)
		_, instructions, err = evmEncoder.Decode(
			3379446385462418246,
			models.Address(common.HexToAddress("0x2033C546BC60900f8B765F0a8e7E2376a17cba5d")),
			encodedReport)
		require.NoError(t, err)
		require.Len(t, instructions.SendLiquidityParams, 0)
		require.Len(t, instructions.ReceiveLiquidityParams, 1)
		require.Equal(t, uint64(909606746561742123), instructions.ReceiveLiquidityParams[0].RemoteChainSelector)
		require.Equal(t, assets.Ether(5).ToInt(), instructions.ReceiveLiquidityParams[0].Amount)
		amount, nonce, err := testonlybridge.UnpackFinalizeBridgePayload(instructions.ReceiveLiquidityParams[0].BridgeData)
		require.NoError(t, err)
		require.Equal(t, assets.Ether(5).ToInt().String(), amount.String())
		require.Equal(t, big.NewInt(1).String(), nonce.String())
	})
}
