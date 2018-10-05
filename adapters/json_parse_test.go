package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestJsonParse_Perform(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		value           string
		path            []string
		want            string
		wantError       bool
		wantResultError bool
	}{
		{"existing path", `{"high":"11850.00","last":"11779.99"}`, []string{"last"},
			`{"value":"11779.99"}`, false, false},
		{"nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"doesnotexist"},
			`{"value":null}`, true, false},
		{"double nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"no", "really"},
			`{"value":"{\"high\":\"11850.00\",\"last\":\"11779.99\"}"}`, true, true},
		{"array index path", `{"data":[{"availability":"0.99991"}]}`, []string{"data", "0", "availability"},
			`{"value":"0.99991"}`, false, false},
		{"float value", `{"availability":0.99991}`, []string{"availability"},
			`{"value":"0.99991"}`, false, false},
		{
			"index array",
			`{"data": [0, 1]}`,
			[]string{"data", "0"},
			`{"value":"0"}`,
			false,
			false,
		},
		{
			"index array of array",
			`{"data": [[0, 1]]}`,
			[]string{"data", "0", "0"},
			`{"value":"0"}`,
			false,
			false,
		},
		{
			"negative index array",
			`{"data": [0, 1]}`,
			[]string{"data", "-1"},
			`{"value":null}`,
			false,
			false,
		},
		{
			"maximum index array",
			`{"data": [0, 1]}`,
			[]string{"data", "18446744073709551615"},
			`{"value":null}`,
			false,
			false,
		},
		{
			"overflow index array",
			`{"data": [0, 1]}`,
			[]string{"data", "18446744073709551616"},
			`{"value":null}`,
			false,
			false,
		},
		{
			"return array",
			`{"data": [[0, 1]]}`,
			[]string{"data", "0"},
			`{"value":"[0,1]"}`,
			false,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.RunResultWithValue(test.value)
			adapter := adapters.JSONParse{Path: test.path}
			result := adapter.Perform(input, nil)
			assert.Equal(t, test.want, result.Data.String())

			if test.wantResultError {
				assert.Error(t, result.GetError())
			} else {
				assert.NoError(t, result.GetError())
			}
		})
	}
}

func TestJSON_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		want      []string
		wantError bool
	}{
		{"array", `{"path":["1","b"]}`, []string{"1", "b"}, false},
		{"array with dots", `{"path":["1",".","b"]}`, []string{"1", ".", "b"}, false},
		{"string", `{"path":"first"}`, []string{"first"}, false},
		{"dot delimited", `{"path":"1.b"}`, []string{"1", "b"}, false},
		{"dot delimited empty string", `{"path":"1...b"}`, []string{"1", "", "", "b"}, false},
		{"unclosed array errors", `{"path":["1"}`, []string{}, true},
		{"unclosed string errors", `{"path":"1.2}`, []string{}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := adapters.JSONParse{}
			err := json.Unmarshal([]byte(test.input), &a)
			cltest.AssertError(t, test.wantError, err)
			if !test.wantError {
				assert.Equal(t, test.want, []string(a.Path))
			}
		})
	}
}
