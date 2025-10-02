package cltest

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// LogFromFixture create ethtypes.log from file path
func LogFromFixture(t *testing.T, path string) types.Log {
	value := gjson.Get(string(MustReadFile(t, path)), "params.result")
	var el types.Log
	require.NoError(t, json.Unmarshal([]byte(value.String()), &el))

	return el
}
