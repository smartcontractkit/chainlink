package main

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/handler"
)

func TestRevertReason(t *testing.T) {
	errorCodeString := "9fe2f95a000000000000000000000000ec1a1eca1d0f8aceae364ab2306e2368f1c6f99000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000024cbbf11130100000000000000000000000000000000000000000000056bc75e2d6310000000000000000000000000000000000000000000000000000000000000"

	decodedError, err := handler.DecodeErrorStringFromABI(errorCodeString)
	if err != nil {
		fmt.Printf("Error decoding error string: %v\n", err)
		return
	}

	fmt.Println(decodedError)
}
