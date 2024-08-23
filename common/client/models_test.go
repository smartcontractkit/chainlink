package client

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendTxReturnCode_String(t *testing.T) {
	// ensure all the SendTxReturnCodes have proper name
	for c := 1; c < int(sendTxReturnCodeLen); c++ {
		strC := SendTxReturnCode(c).String()
		if strings.Contains(strC, "SendTxReturnCode(") {
			t.Errorf("Expected %s to have a proper string representation", strC)
		}
	}
}

func TestSyncStatus_String(t *testing.T) {
	t.Run("All of the statuses have proper string representation", func(t *testing.T) {
		for i := syncStatusNotInSyncWithPool; i < syncStatusLen; i <<= 1 {
			// ensure that i's string representation is not equal to `syncStatus(%d)`
			assert.NotContains(t, i.String(), "syncStatus(")
		}
	})
	t.Run("Unwraps mask", func(t *testing.T) {
		testCases := []struct {
			Mask        syncStatus
			ExpectedStr string
		}{
			{
				ExpectedStr: "Synced",
			},
			{
				Mask:        syncStatusNotInSyncWithPool | syncStatusNoNewHead,
				ExpectedStr: "NotInSyncWithRPCPool,NoNewHead",
			},
			{
				Mask:        syncStatusNotInSyncWithPool | syncStatusNoNewHead | syncStatusNoNewFinalizedHead,
				ExpectedStr: "NotInSyncWithRPCPool,NoNewHead,NoNewFinalizedHead",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.ExpectedStr, func(t *testing.T) {
				assert.Equal(t, testCase.ExpectedStr, testCase.Mask.String())
			})
		}
	})
}
