package models_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func TestMarshalReportMetadata(t *testing.T) {
	t.Parallel()

	t.Run("marshal json", func(t *testing.T) {
		rm := models.ReportMetadata{}
		b, err := json.Marshal(rm)
		require.NoError(t, err, "failed to marshal empty ReportMetadata")

		var unmarshalled models.ReportMetadata
		err = json.Unmarshal(b, &unmarshalled)
		require.NoError(t, err, "failed to unmarshal empty ReportMetadata")
		require.Equal(t, rm, unmarshalled, "marshalled and unmarshalled ReportMetadata should be equal")

		rm = models.ReportMetadata{
			Transfers: []models.Transfer{
				models.NewTransfer(1, 2, big.NewInt(3), time.Now().UTC(), []byte{}),
			},
			LiquidityManagerAddress: models.Address(testutils.NewAddress()),
			NetworkID:               1,
			ConfigDigest: models.ConfigDigest{
				ConfigDigest: testutils.Random32Byte(),
			},
		}
		b, err = json.Marshal(rm)
		require.NoError(t, err, "failed to marshal ReportMetadata")

		err = json.Unmarshal(b, &unmarshalled)
		require.NoError(t, err, "failed to unmarshal ReportMetadata")
		require.Equal(t, rm, unmarshalled, "marshalled and unmarshalled ReportMetadata should be equal")
	})

	t.Run("marshal onchain", func(t *testing.T) {
		rm := models.ReportMetadata{
			NetworkID:               1,
			LiquidityManagerAddress: models.Address(testutils.NewAddress()),
			Transfers: []models.Transfer{
				models.NewTransfer(1, 2, big.NewInt(3), time.Now().UTC(), []byte{}), // send from 1 to 2
				models.NewTransfer(3, 1, big.NewInt(3), time.Now().UTC(), []byte{}), // receive from 3 to 1
			},
		}
		instructions, err := rm.ToLiquidityInstructions()
		require.NoError(t, err, "failed to convert ReportMetadata to LiquidityInstructions")

		encoded, err := rm.OnchainEncode()
		require.NoError(t, err, "failed to encode ReportMetadata")

		r, decodedInstructions, err := models.DecodeReport(rm.NetworkID, rm.LiquidityManagerAddress, encoded)
		require.NoError(t, err, "failed to unmarshal ReportMetadata")
		require.Equal(t, instructions, decodedInstructions, "marshalled and unmarshalled instructions should be equal")
		require.Equal(t, rm.NetworkID, r.NetworkID, "marshalled and unmarshalled NetworkID should be equal")
		require.Equal(t, rm.LiquidityManagerAddress, r.LiquidityManagerAddress, "marshalled and unmarshalled LiquidityManagerAddress should be equal")
		require.Equal(t, rm.Transfers[0].Amount, r.Transfers[0].Amount, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[0].From, r.Transfers[0].From, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[0].To, r.Transfers[0].To, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].Amount, r.Transfers[1].Amount, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].From, r.Transfers[1].From, "marshalled and unmarshalled Transfers should be equal")
		require.Equal(t, rm.Transfers[1].To, r.Transfers[1].To, "marshalled and unmarshalled Transfers should be equal")
	})
}
