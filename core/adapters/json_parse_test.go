package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestJsonParse_Perform(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		result          string
		path            []string
		wantData        string
		wantStatus      models.RunStatus
		wantResultError bool
	}{
		{"existing path", `{"high":"11850.00","last":"11779.99"}`, []string{"last"},
			`{"result":"11779.99"}`, models.RunStatusCompleted, false},
		{"nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"doesnotexist"},
			`{"result":null}`, models.RunStatusCompleted, false},
		{"double nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"no", "really"},
			`{"result":"{\"high\":\"11850.00\",\"last\":\"11779.99\"}"}`, models.RunStatusErrored, true},
		{"array index path", `{"data":[{"availability":"0.99991"}]}`, []string{"data", "0", "availability"},
			`{"result":"0.99991"}`, models.RunStatusCompleted, false},
		{"float result", `{"availability":0.99991}`, []string{"availability"},
			`{"result":0.99991}`, models.RunStatusCompleted, false},
		{
			"index array",
			`{"data": [0, 1]}`,
			[]string{"data", "0"},
			`{"result":0}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"index array of array",
			`{"data": [[0, 1]]}`,
			[]string{"data", "0", "0"},
			`{"result":0}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"index of negative one",
			`{"data": [0, 1]}`,
			[]string{"data", "-1"},
			`{"result":1}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"index of negative array length",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]}`,
			[]string{"data", "-10"},
			`{"result":0}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"index of negative array length minus one",
			`{"data": [0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55]}`,
			[]string{"data", "-12"},
			`{"result":null}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"maximum index array",
			`{"data": [0, 1]}`,
			[]string{"data", "18446744073709551615"},
			`{"result":null}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"overflow index array",
			`{"data": [0, 1]}`,
			[]string{"data", "18446744073709551616"},
			`{"result":null}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"return array",
			`{"data": [[0, 1]]}`,
			[]string{"data", "0"},
			`{"result":[0,1]}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"return false",
			`{"data": false}`,
			[]string{"data"},
			`{"result":false}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"return true",
			`{"data": true}`,
			[]string{"data"},
			`{"result":true}`,
			models.RunStatusCompleted,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.JSONWithResult(t, test.result)
			adapter := adapters.JSONParse{Path: test.path}
			result := adapter.Perform(input, models.RunResult{}, nil)
			assert.Equal(t, test.wantData, result.Data.String())
			assert.Equal(t, test.wantStatus, result.Status)

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
