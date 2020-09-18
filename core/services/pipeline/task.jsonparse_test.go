package pipeline

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONParseTask(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		path            []string
		wantData        interface{}
		wantResultError bool
	}{
		{"existing path", `{"high":"11850.00","last":"11779.99"}`, []string{"last"},
			"11779.99", false},
		{"nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"doesnotexist"},
			nil, false},
		{"double nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"no", "really"},
			nil, true},
		{"array index path", `{"data":[{"availability":"0.99991"}]}`, []string{"data", "0", "availability"},
			"0.99991", false},
		{"float result", `{"availability":0.99991}`, []string{"availability"},
			0.99991, false},
		{
			"index array",
			`{"data": [0, 1]}`,
			[]string{"data", "0"},
			float64(0),
			false,
		},
		{
			"index array of array",
			`{"data": [[0, 1]]}`,
			[]string{"data", "0", "0"},
			float64(0),
			false,
		},
		{
			"index of negative one",
			`{"data": [0, 1]}`,
			[]string{"data", "-1"},
			float64(1),
			false,
		},
		{
			"index of negative array length",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]}`,
			[]string{"data", "-10"},
			float64(0),
			false,
		},
		{
			"index of negative array length minus one",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`,
			[]string{"data", "-12"},
			nil,
			false,
		},
		{
			"maximum index array",
			`{"data": [0, 1]}`,
			[]string{"data", "18446744073709551615"},
			nil,
			false,
		},
		{
			"overflow index array",
			`{"data": [0, 1]}`,
			[]string{"data", "18446744073709551616"},
			nil,
			false,
		},
		{
			"return array",
			`{"data": [[0, 1]]}`,
			[]string{"data", "0"},
			[]interface{}{float64(0), float64(1)},
			false,
		},
		{
			"return false",
			`{"data": false}`,
			[]string{"data"},
			false,
			false,
		},
		{
			"return true",
			`{"data": true}`,
			[]string{"data"},
			true,
			false,
		},
		{
			"regression test: keys in the path have dots",
			`{
                "Realtime Currency Exchange Rate": {
                    "1. From_Currency Code": "LEND",
                    "2. From_Currency Name": "EthLend",
                    "3. To_Currency Code": "ETH",
                    "4. To_Currency Name": "Ethereum",
                    "5. Exchange Rate": "0.00058217",
                    "6. Last Refreshed": "2020-06-22 19:14:04",
                    "7. Time Zone": "UTC",
                    "8. Bid Price": "0.00058217",
                    "9. Ask Price": "0.00058217"
                }
            }`,
			[]string{
				"Realtime Currency Exchange Rate",
				"5. Exchange Rate",
			},
			"0.00058217",
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			task := JSONParseTask{Path: test.path}
			result := task.Run([]Result{{Value: test.input}})

			if test.wantResultError {
				require.Error(t, result.Error)
				require.Nil(t, result.Value)
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.wantData, result.Value)
			}
		})
	}
}
