package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustDecimal(t *testing.T, arg string) *decimal.Decimal {
	ret, err := decimal.NewFromString(arg)
	require.NoError(t, err)
	return &ret
}

func TestMultiply_Marshal(t *testing.T) {
	tests := []struct {
		name string
		obj  adapters.Multiply
		exp  string
	}{
		{
			"w/ value",
			adapters.Multiply{Times: mustDecimal(t, "3.142")},
			`{"times":"3.142"}`,
		},
		{
			"w/ large value",
			adapters.Multiply{Times: mustDecimal(t, "1000000000000000000")},
			`{"times":"1000000000000000000"}`,
		},
		{
			"w/ large float",
			adapters.Multiply{Times: mustDecimal(t, "100000000000000000000.23")},
			`{"times":"100000000000000000000.23"}`,
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
			`{"Times": 1.23}`,
			adapters.Multiply{Times: mustDecimal(t, "1.23")},
		},
		{
			"w/ large value",
			`{"Times": 1000000000000000000}`,
			adapters.Multiply{Times: mustDecimal(t, "1000000000000000000")},
		},
		{
			"w/ large string",
			`{"Times": 100000000000000000000.23}`,
			adapters.Multiply{Times: mustDecimal(t, "100000000000000000000.23")},
		},
		{
			"w/ large float",
			`{"Times": 100000000000000000000.23}`,
			adapters.Multiply{Times: mustDecimal(t, "100000000000000000000.23")},
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
		Times decimal.Decimal
		json  string
		want  string
	}{
		{"by 100", *mustDecimal(t, "100"), `{"result":"1.23"}`, "123"},
		{"float", *mustDecimal(t, "100"), `{"result":1.23}`, "123"},
		{"negative", *mustDecimal(t, "-5"), `{"result":"1.23"}`, "-6.15"},
		{"no times parameter", *mustDecimal(t, "1"), `{"result":"3.14"}`, "3.14"},
		{"zero", *mustDecimal(t, "0"), `{"result":"1.23"}`, "0"},
		{"large value", *mustDecimal(t, "1000000000000000000"), `{"result":"1.23"}`, "1230000000000000000"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Multiply{Times: &test.Times}
			result := adapter.Perform(input, nil, nil)

			require.NoError(t, result.Error())
			assert.Equal(t, test.want, result.Result().String())
		})
	}
}

func TestMultiply_Perform_Failure(t *testing.T) {
	tests := []struct {
		name  string
		Times decimal.Decimal
		json  string
	}{
		{"object", *mustDecimal(t, "100"), `{"result":{"foo":"bar"}}`},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, test.json)
			adapter := adapters.Multiply{Times: &test.Times}
			result := adapter.Perform(input, nil, nil)
			require.Error(t, result.Error())
		})
	}
}
