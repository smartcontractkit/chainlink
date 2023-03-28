package pipeline_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestJSONParseTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		data              string
		path              string
		separator         string
		lax               string
		vars              pipeline.Vars
		inputs            []pipeline.Result
		wantData          interface{}
		wantErrorCause    error
		wantErrorContains string
	}{
		{
			"array index path",
			"",
			"data,0,availability",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data":[{"availability":"0.99991"}]}`}},
			"0.99991",
			nil,
			"",
		},
		{
			"large int result",
			"",
			"some_id",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"some_id":1564679049192120321}`}},
			int64(1564679049192120321),
			nil,
			"",
		},
		{
			"float result",
			"",
			"availability",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"availability":3.14}`}},
			3.14,
			nil,
			"",
		},
		{
			"index array",
			"",
			"data,0",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			int64(0),
			nil,
			"",
		},
		{
			"index array of array",
			"",
			"data,0,0",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [[0, 1]]}`}},
			int64(0),
			nil,
			"",
		},
		{
			"index of negative one",
			"",
			"data,-1",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			int64(1),
			nil,
			"",
		},
		{
			"index of negative array length",
			"",
			"data,-10",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]}`}},
			int64(0),
			nil,
			"",
		},
		{
			"index of negative array length minus one with lax returns nil",
			"",
			"data,-12",
			"",
			"true",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`}},
			nil,
			nil,
			"",
		},
		{
			"index of negative array length minus one without lax returns error",
			"",
			"data,-12",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"maximum index array with lax returns nil",
			"",
			"data,18446744073709551615",
			"",
			"true",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			nil,
			"",
		},
		{
			"maximum index array without lax returns error",
			"",
			"data,18446744073709551615",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"overflow index array with lax returns nil",
			"",
			"data,18446744073709551616",
			"",
			"true",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			nil,
			"",
		},
		{
			"overflow index array without lax returns error",
			"",
			"data,18446744073709551616",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [0, 1]}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"return array",
			"",
			"data,0",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": [[0, 1]]}`}},
			[]interface{}{int64(0), int64(1)},
			nil,
			"",
		},
		{
			"return false",
			"",
			"data",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": false}`}},
			false,
			nil,
			"",
		},
		{
			"return true",
			"",
			"data",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"data": true}`}},
			true,
			nil,
			"",
		},
		{
			"regression test: keys in the path have dots",
			"",
			"Realtime Currency Exchange Rate,5. Exchange Rate",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
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
			"custom separator: keys in the path have commas",
			"",
			"foo.bar1,bar2,bar3",
			".",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{
                "foo": {
                    "bar1": "LEND",
                    "bar1,bar2": "EthLend",
                    "bar2,bar3": "ETH",
                    "bar1,bar3": "Ethereum",
                    "bar1,bar2,bar3": "0.00058217",
                    "bar1.bar2.bar3": "2020-06-22 19:14:04"
                }
            }`}},
			"0.00058217",
			nil,
			"",
		},
		{
			"custom separator: diabolical keys in the path",
			"",
			"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.,/\\[]{}|<>?_+-=!@#$%^&*()__hacky__separator__foo",
			"__hacky__separator__",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{
                "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.,/\\[]{}|<>?_+-=!@#$%^&*()": {
                    "foo": "LEND",
                    "bar": "EthLend"
                }
            }`}},
			"LEND",
			nil,
			"",
		},
		{
			"missing top-level key with lax=false returns error",
			"",
			"baz",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"foo": 1}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"missing nested key with lax=false returns error",
			"",
			"foo,bar",
			"",
			"false",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"foo": {}}`}},
			nil,
			pipeline.ErrKeypathNotFound,
			"",
		},
		{
			"missing top-level key with lax=true returns nil",
			"",
			"baz",
			"",
			"true",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{}`}},
			nil,
			nil,
			"",
		},
		{
			"missing nested key with lax=true returns nil",
			"",
			"foo,baz",
			"",
			"true",
			pipeline.NewVarsFrom(nil),
			[]pipeline.Result{{Value: `{"foo": {}}`}},
			nil,
			nil,
			"",
		},
		{
			"variable data",
			"$(foo.bar)",
			"data,0,availability",
			"",
			"false",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{"bar": `{"data":[{"availability":"0.99991"}]}`},
			}),
			[]pipeline.Result{},
			"0.99991",
			nil,
			"",
		},
		{
			"empty path",
			"$(foo.bar)",
			"",
			"",
			"false",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo": map[string]interface{}{"bar": `{"data":["stevetoshi sergeymoto"]}`},
			}),
			[]pipeline.Result{},
			map[string]interface{}{"data": []interface{}{"stevetoshi sergeymoto"}},
			nil,
			"",
		},
		{
			"no data or input",
			"",
			"$(chain.link)",
			"",
			"false",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo":   map[string]interface{}{"bar": `{"data":[{"availability":"0.99991"}]}`},
				"chain": map[string]interface{}{"link": "data,0,availability"},
			}),
			[]pipeline.Result{},
			"0.99991",
			pipeline.ErrIndexOutOfRange,
			"data",
		},
		{
			"malformed 'lax' param",
			"$(foo.bar)",
			"$(chain.link)",
			"",
			"sergey",
			pipeline.NewVarsFrom(map[string]interface{}{
				"foo":   map[string]interface{}{"bar": `{"data":[{"availability":"0.99991"}]}`},
				"chain": map[string]interface{}{"link": "data,0,availability"},
			}),
			[]pipeline.Result{},
			"0.99991",
			pipeline.ErrBadInput,
			"lax",
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			task := pipeline.JSONParseTask{
				BaseTask:  pipeline.NewBaseTask(0, "json", nil, nil, 0),
				Path:      test.path,
				Separator: test.separator,
				Data:      test.data,
				Lax:       test.lax,
			}
			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)

			if test.wantErrorCause != nil {
				require.Equal(t, test.wantErrorCause, errors.Cause(result.Error))
				if test.wantErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.wantErrorContains)
				}

				require.Nil(t, result.Value)
				val, err := test.vars.Get("json")
				require.Equal(t, pipeline.ErrKeypathNotFound, errors.Cause(err))
				require.Nil(t, val)
			} else {
				require.NoError(t, result.Error)
				require.Equal(t, test.wantData, result.Value)
			}
		})
	}
}
