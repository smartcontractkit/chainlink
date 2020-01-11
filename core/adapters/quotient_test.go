package adapters_test

import (
	"encoding/json"
	"testing"

	"chainlink/core/adapters"
	"chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			t.Parallel()
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Quotient{}
			json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil)

			require.NoError(t, result.Error())
			assert.Equal(t, test.want, result.Result().String())
		})
	}
}

func TestQuotient_Perform_Error(t *testing.T) {
	tests := []struct {
		name   string
		params string
		json   string
		want   string
	}{
		{"string zero integer", `{"dividend":"1"}`, `{"result":0}`, "0"},
		{"zero string zero float", `{"dividend":"0"}`, `{"result":"0"}`, "0"},
		{"object", `{"dividend":100}`, `{"result":{"foo":"bar"}}`, ""},
		{"string object", `{"dividend":"100"}`, `{"result":{"foo":"bar"}}`, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Quotient{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil)

			require.Error(t, result.Error())
			assert.NoError(t, jsonErr)
		})
	}
}

func TestQuotient_Perform_JSONParseError(t *testing.T) {
	tests := []struct {
		name   string
		params string
		json   string
		want   string
	}{
		{"array string", `{"dividend":[1, 2, 3]}`, `{"result":"1.23"}`, ""},
		{"rubbish string", `{"dividend":"123aaa123"}`, `{"result":"1.23"}`, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Quotient{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			adapter.Perform(input, nil)

			assert.Error(t, jsonErr)
		})
	}
}
