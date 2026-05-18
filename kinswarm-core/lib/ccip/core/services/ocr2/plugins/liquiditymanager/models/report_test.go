package models_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestReport_Marshal(t *testing.T) {
	t.Run("marshal json", func(t *testing.T) {
		rm := models.Report{}
		b, err := json.Marshal(rm)
		require.NoError(t, err, "failed to marshal empty ReportMetadata")

		var unmarshalled models.Report
		err = json.Unmarshal(b, &unmarshalled)
		require.NoError(t, err, "failed to unmarshal empty ReportMetadata")
		require.Equal(t, rm, unmarshalled, "marshalled and unmarshalled ReportMetadata should be equal")

		rm = models.Report{
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
}
