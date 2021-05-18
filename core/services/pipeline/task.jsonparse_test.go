package pipeline

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONParseTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		input           string
		path            string
		lax             bool
		wantData        interface{}
		wantResultError bool
	}{
		{"array index path", `{"data":[{"availability":"0.99991"}]}`, "[data,0,availability]", false, "0.99991", false},
		{"float result", `{"availability":0.99991}`, "[availability]", false, 0.99991, false},
		{
			"index array",
			`{"data": [0, 1]}`,
			"[data,0]",
			false,
			float64(0),
			false,
		},
		{
			"index array of array",
			`{"data": [[0, 1]]}`,
			"[data,0,0]",
			false,
			float64(0),
			false,
		},
		{
			"index of negative one",
			`{"data": [0, 1]}`,
			"[data,-1]",
			false,
			float64(1),
			false,
		},
		{
			"index of negative array length",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]}`,
			"[data,-10]",
			false,
			float64(0),
			false,
		},
		{
			"index of negative array length minus one with lax returns nil",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`,
			"[data,-12]",
			true,
			nil,
			false,
		},
		{
			"index of negative array length minus one without lax returns error",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`,
			"[data,-12]",
			false,
			nil,
			true,
		},
		{
			"maximum index array with lax returns nil",
			`{"data": [0, 1]}`,
			"[data,18446744073709551615]",
			true,
			nil,
			false,
		},
		{
			"maximum index array without lax returns error",
			`{"data": [0, 1]}`,
			"[data,18446744073709551615]",
			false,
			nil,
			true,
		},
		{
			"overflow index array with lax returns nil",
			`{"data": [0, 1]}`,
			"[data,18446744073709551616]",
			true,
			nil,
			false,
		},
		{
			"overflow index array without lax returns error",
			`{"data": [0, 1]}`,
			"[data,18446744073709551616]",
			false,
			nil,
			true,
		},
		{
			"return array",
			`{"data": [[0, 1]]}`,
			"[data,0]",
			false,
			[]interface{}{float64(0), float64(1)},
			false,
		},
		{
			"return false",
			`{"data": false}`,
			"[data]",
			false,
			false,
			false,
		},
		{
			"return true",
			`{"data": true}`,
			"[data]",
			false,
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
			"[Realtime Currency Exchange Rate,5. Exchange Rate]",
			false,
			"0.00058217",
			false,
		},
		{
			"missing top-level key with lax=false returns error",
			`{"foo": 1}`,
			"[baz]",
			false,
			nil,
			true,
		},
		{
			"missing nested key with lax=false returns error",
			`{"foo": {}}`,
			"[foo,bar]",
			false,
			nil,
			true,
		},
		{
			"missing top-level key with lax=true returns nil",
			`{}`,
			"[baz]",
			true,
			nil,
			false,
		},
		{
			"missing nested key with lax=true returns nil",
			`{"foo": {}}`,
			"[foo,baz]",
			true,
			nil,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			task := JSONParseTask{Path: test.path, Lax: fmt.Sprintf("%v", test.lax)}
			result := task.Run(context.Background(), nil, JSONSerializable{}, []Result{{Value: test.input}})

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
