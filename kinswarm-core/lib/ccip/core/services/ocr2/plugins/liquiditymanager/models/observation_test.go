package models_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestMarshalObservation(t *testing.T) {
	observation := models.Observation{
		LiquidityPerChain: []models.NetworkLiquidity{
			{
				Network:   models.NetworkSelector(1),
				Liquidity: ubig.NewI(100),
			},
			{
				Network:   models.NetworkSelector(2),
				Liquidity: ubig.NewI(200),
			},
		},
		ResolvedTransfers: []models.Transfer{
			{
				From:               models.NetworkSelector(1),
				To:                 models.NetworkSelector(2),
				Amount:             ubig.NewI(3),
				Sender:             models.Address(common.HexToAddress("0x1")),
				Receiver:           models.Address(common.HexToAddress("0x2")),
				LocalTokenAddress:  models.Address(common.HexToAddress("0x3")),
				RemoteTokenAddress: models.Address(common.HexToAddress("0x4")),
				BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3},
				NativeBridgeFee:    ubig.NewI(4),
			},
			{
				From:               models.NetworkSelector(3),
				To:                 models.NetworkSelector(4),
				Amount:             ubig.NewI(4),
				Sender:             models.Address(common.HexToAddress("0x5")),
				Receiver:           models.Address(common.HexToAddress("0x6")),
				LocalTokenAddress:  models.Address(common.HexToAddress("0x7")),
				RemoteTokenAddress: models.Address(common.HexToAddress("0x8")),
				BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3, 0x4},
				NativeBridgeFee:    ubig.NewI(5),
			},
		},
		PendingTransfers: []models.PendingTransfer{
			{
				Transfer: models.Transfer{
					From:               models.NetworkSelector(1),
					To:                 models.NetworkSelector(2),
					Amount:             ubig.NewI(3),
					Sender:             models.Address(common.HexToAddress("0x1")),
					Receiver:           models.Address(common.HexToAddress("0x2")),
					LocalTokenAddress:  models.Address(common.HexToAddress("0x3")),
					RemoteTokenAddress: models.Address(common.HexToAddress("0x4")),
					BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3},
					NativeBridgeFee:    ubig.NewI(4),
				},
				Status: models.TransferStatusReady,
			},
			{
				Transfer: models.Transfer{
					From:               models.NetworkSelector(3),
					To:                 models.NetworkSelector(4),
					Amount:             ubig.NewI(4),
					Sender:             models.Address(common.HexToAddress("0x5")),
					Receiver:           models.Address(common.HexToAddress("0x6")),
					LocalTokenAddress:  models.Address(common.HexToAddress("0x7")),
					RemoteTokenAddress: models.Address(common.HexToAddress("0x8")),
					BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3, 0x4},
					NativeBridgeFee:    ubig.NewI(5),
				},
				Status: models.TransferStatusReady,
			},
		},
		Edges: []models.Edge{
			{
				Source: models.NetworkSelector(1),
				Dest:   models.NetworkSelector(2),
			},
			{
				Source: models.NetworkSelector(3),
				Dest:   models.NetworkSelector(4),
			},
		},
	}
	jsonBytes, err := json.Marshal(observation)
	require.NoError(t, err, "failed to marshal observation to json")
	t.Log(string(jsonBytes))
	var unmarshaled models.Observation
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	require.NoError(t, err, "failed to unmarshal observation from json")
	require.Equal(t, observation, unmarshaled, "marshalled and unmarshalled observation should be equal")
}

func TestMarshalOutcome(t *testing.T) {
	outcome := models.Outcome{
		ProposedTransfers: []models.ProposedTransfer{
			{
				From:   models.NetworkSelector(1),
				To:     models.NetworkSelector(2),
				Amount: ubig.NewI(3),
			},
		},
		ResolvedTransfers: []models.Transfer{
			{
				From:               models.NetworkSelector(1),
				To:                 models.NetworkSelector(2),
				Amount:             ubig.NewI(3),
				Sender:             models.Address(common.HexToAddress("0x1")),
				Receiver:           models.Address(common.HexToAddress("0x2")),
				LocalTokenAddress:  models.Address(common.HexToAddress("0x3")),
				RemoteTokenAddress: models.Address(common.HexToAddress("0x4")),
				BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3},
				NativeBridgeFee:    ubig.NewI(4),
			},
			{
				From:               models.NetworkSelector(3),
				To:                 models.NetworkSelector(4),
				Amount:             ubig.NewI(4),
				Sender:             models.Address(common.HexToAddress("0x5")),
				Receiver:           models.Address(common.HexToAddress("0x6")),
				LocalTokenAddress:  models.Address(common.HexToAddress("0x7")),
				RemoteTokenAddress: models.Address(common.HexToAddress("0x8")),
				BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3, 0x4},
				NativeBridgeFee:    ubig.NewI(5),
			},
		},
		PendingTransfers: []models.PendingTransfer{
			{
				Transfer: models.Transfer{
					From:               models.NetworkSelector(1),
					To:                 models.NetworkSelector(2),
					Amount:             ubig.NewI(3),
					Sender:             models.Address(common.HexToAddress("0x1")),
					Receiver:           models.Address(common.HexToAddress("0x2")),
					LocalTokenAddress:  models.Address(common.HexToAddress("0x3")),
					RemoteTokenAddress: models.Address(common.HexToAddress("0x4")),
					BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3},
					NativeBridgeFee:    ubig.NewI(4),
				},
				Status: models.TransferStatusReady,
			},
			{
				Transfer: models.Transfer{
					From:               models.NetworkSelector(3),
					To:                 models.NetworkSelector(4),
					Amount:             ubig.NewI(4),
					Sender:             models.Address(common.HexToAddress("0x5")),
					Receiver:           models.Address(common.HexToAddress("0x6")),
					LocalTokenAddress:  models.Address(common.HexToAddress("0x7")),
					RemoteTokenAddress: models.Address(common.HexToAddress("0x8")),
					BridgeData:         hexutil.Bytes{0x1, 0x2, 0x3, 0x4},
					NativeBridgeFee:    ubig.NewI(5),
				},
				Status: models.TransferStatusReady,
			},
		},
	}
	jsonBytes, err := json.Marshal(outcome)
	require.NoError(t, err, "failed to marshal outcome to json")
	t.Log(string(jsonBytes))
	var unmarshaled models.Outcome
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	require.NoError(t, err, "failed to unmarshal outcome from json")
}
