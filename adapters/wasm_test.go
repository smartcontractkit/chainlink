// +build sgx_enclave

package adapters_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

const (
	// CheckEthProgram was compiled then base64ed from internal/fixtures/wasm/checkethf.wat
	// This program compares the input value to 450 using i64.lt_s
	CheckEthProgram = "AGFzbQEAAAABBgFgAXwBfwMCAQAHCwEHcGVyZm9ybQAAChABDgBEAAAAAAAgfEAgAGML"
)

func TestWasm_Perform(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		json      string
		want      string
		errored   bool
		jsonError bool
	}{
		{
			"check eth less than 450",
			fmt.Sprintf(`{"wasm":"%s"}`, CheckEthProgram),
			`{"value": 449.9}`,
			"0",
			false,
			false,
		},
		{
			"check eth equals 450",
			fmt.Sprintf(`{"wasm":"%s"}`, CheckEthProgram),
			`{"value": 450.0}`,
			"0",
			false,
			false,
		},
		{
			"check eth greater than 450",
			fmt.Sprintf(`{"wasm":"%s"}`, CheckEthProgram),
			`{"value": 450.1}`,
			"1",
			false,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.Wasm{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil)

			if test.jsonError {
				assert.Error(t, jsonErr)
			} else if test.errored {
				assert.Error(t, result.GetError())
				assert.NoError(t, jsonErr)
			} else {
				val, err := result.Value()
				assert.NoError(t, err)
				assert.Equal(t, test.want, val)
				assert.NoError(t, result.GetError())
				assert.NoError(t, jsonErr)
			}
		})
	}
}
