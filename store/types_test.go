package store_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
)

func TestLog_UnmarshalEmptyTxHash(t *testing.T) {
	t.Parallel()

	input := `{
		"transactionHash": null,
		"transactionIndex": "0x3",
		"address": "0x1aee7c03606fca5035d204c3818d0660bb230e44",
		"blockNumber": "0x8bf99b",
		"topics": ["0xdeadbeefdeadbeedeadbeedeadbeefffdeadbeefdeadbeedeadbeedeadbeefff"],
		"blockHash": "0xdb777676330c067e3c3a6dbfc2d51282cac5bcc1b7a884dd8d85ba72ca1f147e",
		"data": "0xdeadbeef",
		"logIndex": "0x5",
		"transactionLogIndex": "0x3"
	}`

	var log store.Log
	err := json.Unmarshal([]byte(input), &log)
	assert.NoError(t, err)
}
