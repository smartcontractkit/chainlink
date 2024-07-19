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

func TestSyncIssue(t *testing.T) {
	t.Run("All of the issues have proper string representation", func(t *testing.T) {
		for i := syncIssueNotInSyncWithPool; i < syncIssueLen; i <<= 1 {
			// ensure that i's string representation is not equal to `syncIssue(%d)`
			assert.NotContains(t, i.String(), "syncIssue(")
		}
	})
	t.Run("Unwraps mask", func(t *testing.T) {
		testCases := []struct {
			Mask        syncIssue
			ExpectedStr string
		}{
			{
				ExpectedStr: "synced",
			},
			{
				Mask:        syncIssueNotInSyncWithPool | syncIssueHeadIsNotIncreasing,
				ExpectedStr: "NotInSyncWithRPCPool,HeadIsNotIncreasing",
			},
			{
				Mask:        syncIssueNotInSyncWithPool | syncIssueHeadIsNotIncreasing | syncIssueFinalizedHeadIsNotIncreasing,
				ExpectedStr: "NotInSyncWithRPCPool,HeadIsNotIncreasing,FinalizedHeadIsNotIncreasing",
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.ExpectedStr, func(t *testing.T) {
				assert.Equal(t, testCase.ExpectedStr, testCase.Mask.String())
			})
		}
	})
}
