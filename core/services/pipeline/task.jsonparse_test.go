package pipeline_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func TestJSONParseTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		data              string
		path              string
		lax               string
		inputs            []pipeline.Result
		wantData          interface{}
		wantErrorCause    error
		wantErrorContains string
	}{
		{
			"array index path",
			"",
			"data,0,availability",
			"false",
			[]pipeline.Result{{Value: `{"data":[{"availability":"0.99991"}]}`}},
			"0.99991",
			nil,
			"",
		},
		{
			"float result",
			"",
			"availability",
			"false",
			[]pipeline.Result{{Value: `{"availability":0.99991}`}},
			0.99991,
			nil,
			"",
		},
		{
			"index array",
			"",
			"data,0",
			"false",
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			float64(0),
			nil,
			"",
		},
		{
			"index array of array",
			"",
			"data,0,0",
			"false",
			[]pipeline.Result{{Value: `{"data": [[0, 1]]}`}},
			float64(0),
			nil,
			"",
		},
		{
			"index of negative one",
			"",
			"data,-1",
			"false",
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			float64(1),
			nil,
			"",
		},
		{
			"index of negative array length",
			"",
			"data,-10",
			"false",
			[]pipeline.Result{{Value: `{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]}`}},
			float64(0),
			nil,
			"",
		},
		{
			"index of negative array length minus one with lax returns nil",
			"",
			"data,-12",
			"true",
			[]pipeline.Result{{Value: `{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`}},
			nil,
			nil,
			"",
		},
		{
			"index of negative array length minus one without lax returns error",
			"",
			"data,-12",
			"false",
			[]pipeline.Result{{Value: `{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"maximum index array with lax returns nil",
			"",
			"data,18446744073709551615",
			"true",
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			nil,
			"",
		},
		{
			"maximum index array without lax returns error",
			"",
			"data,18446744073709551615",
			"false",
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"overflow index array with lax returns nil",
			"",
			"data,18446744073709551616",
			"true",
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			nil,
			"",
		},
		{
			"overflow index array without lax returns error",
			"",
			"data,18446744073709551616",
			"false",
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"return array",
			"",
			"data,0",
			"false",
			[]pipeline.Result{{Value: `{"data": [[0, 1]]}`}},
			[]interface{}{float64(0), float64(1)},
			nil,
			"",
		},
		{
			"return false",
			"",
			"data",
			"false",
			[]pipeline.Result{{Value: `{"data": false}`}},
			false,
			nil,
			"",
		},
		{
			"return true",
			"",
			"data",
			"false",
			[]pipeline.Result{{Value: `{"data": true}`}},
			true,
			nil,
			"",
		},
		{
			"regression test: keys in the path have dots",
			"",
			"Realtime Currency Exchange Rate,5. Exchange Rate",
			"false",
			[]pipeline.Result{{Value: `{
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
            }`}},
			"0.00058217",
			nil,
			"",
		},
		{
			"missing top-level key with lax=false returns error",
			"",
			"baz",
			"false",
			[]pipeline.Result{{Value: `{"foo": 1}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"missing nested key with lax=false returns error",
			"",
			"foo,bar",
			"false",
			[]pipeline.Result{{Value: `{"foo": {}}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"missing top-level key with lax=true returns nil",
			"",
			"baz",
			"true",
			[]pipeline.Result{{Value: `{}`}},
			nil,
			nil,
			"",
		},
		{
			"missing nested key with lax=true returns nil",
			"",
			"foo,baz",
			"true",
			[]pipeline.Result{{Value: `{"foo": {}}`}},
			nil,
			nil,
			"",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.JSONParseTask{
				BaseTask: pipeline.NewBaseTask("json", nil, 0, 0),
				Path:     test.path,
				Data:     test.data,
				Lax:      test.lax,
			}
			result := task.Run(context.Background(), pipeline.JSONSerializable{}, test.inputs)

			if test.wantErrorCause != nil {
				require.Equal(t, test.wantErrorCause, errors.Cause(result.Error))
				if test.wantErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.wantErrorContains)
				}
				require.Nil(t, result.Value)
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.wantData, result.Value)
			}
		})
	}
}
