package client

import (
	"strings"
	"testing"
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
