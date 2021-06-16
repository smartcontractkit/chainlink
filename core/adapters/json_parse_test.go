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
			``, models.RunStatusErrored, true},
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
			`{"result":"0.00058217"}`,
			models.RunStatusCompleted,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			input := cltest.NewRunInputWithResult(test.result)
			adapter := adapters.JSONParse{Path: test.path}
			result := adapter.Perform(input, nil, nil)
			assert.Equal(t, test.wantData, result.Data().String())
			assert.Equal(t, test.wantStatus, result.Status())

			if test.wantResultError {
				assert.Error(t, result.Error())
			} else {
				assert.NoError(t, result.Error())
			}
		})
	}
}

func TestJsonParse_Perform_WithPreParsedJSON(t *testing.T) {
	var parsed models.JSON
	err := json.Unmarshal([]byte(`{"high":"11850.00","last":"11779.99"}`), &parsed)
	assert.NoError(t, err)

	input := cltest.NewRunInputWithResult(parsed)

	adapter := adapters.JSONParse{Path: []string{"last"}}
	result := adapter.Perform(input, nil, nil)
	assert.Equal(t, `{"result":"11779.99"}`, result.Data().String())
	assert.Equal(t, models.RunStatusCompleted, result.Status())
	assert.NoError(t, result.Error())
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
