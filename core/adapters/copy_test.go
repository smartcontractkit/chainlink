package adapters_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/stretchr/testify/assert"
)

func TestCopy_Perform(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		copyPath        []string
		wantData        string
		wantStatus      models.RunStatus
		wantResultError error
	}{
		{
			"existing path",
			`{"high":"11850.00","last":"11779.99"}`,
			[]string{"last"},
			`{"result":"11779.99"}`,
			models.RunStatusCompleted,
			nil,
		},
		{
			"nonexistent path",
			`{"high":"11850.00","last":"11779.99"}`,
			[]string{"doesnotexist"},
			`{"result":null}`,
			models.RunStatusCompleted,
			nil,
		},
		{
			"array index path",
			`{"data":[{"availability":"0.99991"}]}`,
			[]string{"data", "0", "availability"},
			`{"result":"0.99991"}`,
			models.RunStatusCompleted,
			nil,
		},
		{
			"float result",
			`{"availability":0.99991}`,
			[]string{"availability"},
			`{"result":0.99991}`,
			models.RunStatusCompleted,
			nil,
		},
		{
			"result with quotes",
			`{"availability":"\""}`,
			[]string{`"`},
			`{"result":null}`,
			models.RunStatusCompleted,
			nil,
		},
		{
			"index array of array",
			`{"data":[[0,1]]}`,
			[]string{"data", "0", "0"},
			`{"result":0}`,
			models.RunStatusCompleted,
			nil,
		},
		{
			"double nonexistent path",
			`{"high":"11850.00","last":"11779.99"}`,
			[]string{"no", "really"},
			``,
			models.RunStatusErrored,
			errors.New("No value could be found for the key 'no'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := cltest.NewRunInputWithString(t, tt.input)
			adapter := adapters.Copy{CopyPath: tt.copyPath}
			result := adapter.Perform(input, nil, nil)
			assert.Equal(t, tt.wantData, result.Data().String())
			assert.Equal(t, tt.wantStatus, result.Status())

			assert.Equal(t, tt.wantResultError, result.Error())
		})
	}
}

func TestCopy_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		want      []string
		wantError bool
	}{
		{"array", `{"copyPath":["1","b"]}`, []string{"1", "b"}, false},
		{"array with dots", `{"copyPath":["1",".","b"]}`, []string{"1", ".", "b"}, false},
		{"string", `{"copyPath":"first"}`, []string{"first"}, false},
		{"dot delimited", `{"copyPath":"1.b"}`, []string{"1", "b"}, false},
		{"dot delimited empty string", `{"copyPath":"1...b"}`, []string{"1", "", "", "b"}, false},
		{"unclosed array errors", `{"copyPath":["1"}`, []string{}, true},
		{"unclosed string errors", `{"copyPath":"1.2}`, []string{}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := adapters.Copy{}
			err := json.Unmarshal([]byte(test.input), &a)
			cltest.AssertError(t, test.wantError, err)
			if !test.wantError {
				assert.Equal(t, test.want, []string(a.CopyPath))
			}
		})
	}
}
