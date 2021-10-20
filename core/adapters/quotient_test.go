package adapters_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuotient_Marshal(t *testing.T) {
	tests := []struct {
		name string
		obj  adapters.Quotient
		exp  string
	}{
		{
			"w/ value",
			adapters.Quotient{Dividend: big.NewFloat(3.142)},
			`{"dividend":"3.142"}`,
		},
		{
			"w/ value",
			adapters.Quotient{Dividend: big.NewFloat(5)},
			`{"dividend":"5"}`,
		},
		{
			"w/o value",
			adapters.Quotient{Dividend: nil},
			`{}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := json.Marshal(tc.obj)
			require.NoError(t, err)
			require.Equal(t, tc.exp, string(buf))
		})
	}
}

func TestQuotient_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		exp     adapters.Quotient
	}{
		{
			"w/ value",
			`{"dividend": 5}`,
			adapters.Quotient{Dividend: big.NewFloat(5)},
		},
		{
			"w/o value",
			`{}`,
			adapters.Quotient{Dividend: nil},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var m adapters.Quotient
			err := json.Unmarshal([]byte(tc.payload), &m)
			require.NoError(t, err)
			require.Equal(t, tc.exp, m)
		})
	}
}

func TestQuotient_Perform_Success(t *testing.T) {
	tests := []struct {
		name   string
		params string
		json   string
		want   string
	}{
		{"string", `{"dividend":100}`, `{"result":"1.23"}`, "81.30081301"},
		{"integer", `{"dividend":100}`, `{"result":123}`, "0.8130081301"},
		{"float", `{"dividend":100}`, `{"result":1.23}`, "81.30081301"},
		{"zero integer string", `{"dividend":0}`, `{"result":"1.23"}`, "0"},
		{"zero float string", `{"dividend":0.0}`, `{"result":"1.23"}`, "0"},
		{"negative integer string", `{"dividend":-5}`, `{"result":"1.23"}`, "-4.06504065"},
		{"negative integer string", `{"dividend":-5}`, `{"result":"1.23"}`, "-4.06504065"},

		{"no dividend parameter", `{}`, `{"result":"3.14"}`, "3.14"},

		{"string string", `{"dividend":"100"}`, `{"result":"1.23"}`, "81.30081301"},
		{"string integer", `{"dividend":"100"}`, `{"result":123}`, "0.8130081301"},
		{"string float", `{"dividend":"100"}`, `{"result":1.23}`, "81.30081301"},

		{"zero string string", `{"dividend":"0"}`, `{"result":"1.23"}`, "0"},
		{"negative string string", `{"dividend":"-5"}`, `{"result":"1.23"}`, "-4.06504065"},
		{"string", `{"dividend":"1"}`, `{"result":"1.23"}`, "0.8130081301"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Quotient{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil, nil)

			require.NoError(t, result.Error())
			require.NoError(t, jsonErr)
			assert.Equal(t, test.want, result.Result().String())
		})
	}
}

func TestQuotient_Perform_Error(t *testing.T) {
	tests := []struct {
		name   string
		params string
		json   string
	}{
		{"string zero integer", `{"dividend":"1"}`, `{"result":0}`},
		{"zero string zero float", `{"dividend":"0"}`, `{"result":"0"}`},
		{"object", `{"dividend":100}`, `{"result":{"foo":"bar"}}`},
		{"string object", `{"dividend":"100"}`, `{"result":{"foo":"bar"}}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Quotient{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil, nil)

			require.NoError(t, jsonErr)
			assert.Error(t, result.Error())
		})
	}
}

func TestQuotient_Perform_JSONParseError(t *testing.T) {
	tests := []struct {
		name   string
		params string
	}{
		{"array string", `{"dividend":[1, 2, 3]}`},
		{"rubbish string", `{"dividend":"123aaa123"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter := adapters.Quotient{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			assert.Error(t, jsonErr)
		})
	}
}
