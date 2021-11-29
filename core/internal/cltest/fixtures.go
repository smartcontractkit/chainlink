package cltest

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// JSONFromFixture create models.JSON from file path
func JSONFromFixture(t *testing.T, path string) models.JSON {
	return JSONFromBytes(t, MustReadFile(t, path))
}

// LogFromFixture create ethtypes.log from file path
func LogFromFixture(t *testing.T, path string) types.Log {
	value := gjson.Get(string(MustReadFile(t, path)), "params.result")
	var el types.Log
	require.NoError(t, json.Unmarshal([]byte(value.String()), &el))

	return el
}

// TxReceiptFromFixture create ethtypes.log from file path
func TxReceiptFromFixture(t *testing.T, path string) *types.Receipt {
	jsonStr := JSONFromFixture(t, path).Get("result").String()

	var receipt types.Receipt
	err := json.Unmarshal([]byte(jsonStr), &receipt)
	require.NoError(t, err)

	return &receipt
}
