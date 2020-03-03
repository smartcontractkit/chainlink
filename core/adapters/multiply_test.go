package adapters_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"chainlink/core/adapters"
	"chainlink/core/internal/cltest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiply_Marshal(t *testing.T) {
	tests := []struct {
		name string
		obj  adapters.Multiply
		exp  string
	}{
		{
			"w/ value",
			adapters.Multiply{Times: big.NewFloat(3.142)},
			`{"times":"3.142"}`,
		},
		{
			"w/ value",
			adapters.Multiply{Times: big.NewFloat(5)},
			`{"times":"5"}`,
		},
		{
			"w/o value",
			adapters.Multiply{},
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

func TestMultiply_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		exp     adapters.Multiply
	}{
		{
			"w/ value",
			`{"Times": 5}`,
			adapters.Multiply{Times: big.NewFloat(5)},
		},
		{
			"w/o value",
			`{}`,
			adapters.Multiply{Times: nil},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var m adapters.Multiply
			err := json.Unmarshal([]byte(tc.payload), &m)
			require.NoError(t, err)
			require.Equal(t, tc.exp, m)
		})
	}
}

func TestMultiply_Perform(t *testing.T) {
	tests := []struct {
		name  string
		Times *big.Float
		json  string
		want  string
	}{
		{"by 100", big.NewFloat(100), `{"result":"1.23"}`, "123"},
		{"float", big.NewFloat(100), `{"result":1.23}`, "123"},
		{"negative", big.NewFloat(-5), `{"result":"1.23"}`, "-6.15"},
		{"no times parameter", nil, `{"result":"3.14"}`, "3.14"},
		{"zero", big.NewFloat(0), `{"result":"1.23"}`, "0"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Multiply{Times: test.Times}
			result := adapter.Perform(input, nil)

			require.NoError(t, result.Error())
			assert.Equal(t, test.want, result.Result().String())
		})
	}
}

func TestMultiply_Perform_Failure(t *testing.T) {
	tests := []struct {
		name  string
		Times *big.Float
		json  string
		want  string
	}{
		{"object", big.NewFloat(100), `{"result":{"foo":"bar"}}`, ""},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Multiply{Times: test.Times}
			result := adapter.Perform(input, nil)
			require.Error(t, result.Error())
		})
	}
}
