package adapters_test

import (
	"encoding/json"
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
		wantResultError bool
	}{
		{
			"existing path",
			`{"high":"11850.00","last":"11779.99"}`,
			[]string{"last"},
			`{"result":"11779.99"}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"nonexistent path",
			`{"high":"11850.00","last":"11779.99"}`,
			[]string{"doesnotexist"},
			`{"result":null}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"double nonexistent path",
			`{"high":"11850.00","last":"11779.99"}`,
			[]string{"no", "really"},
			``,
			models.RunStatusErrored,
			true,
		},
		{
			"array index path",
			`{"data":[{"availability":"0.99991"}]}`,
			[]string{"data", "0", "availability"},
			`{"result":"0.99991"}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"float result",
			`{"availability":0.99991}`,
			[]string{"availability"},
			`{"result":0.99991}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"result with quotes",
			`{"availability":"\""}`,
			[]string{`"`},
			`{"result":null}`,
			models.RunStatusCompleted,
			false,
		},
		{
			"index array of array",
			`{"data":[[0,1]]}`,
			[]string{"data", "0", "0"},
			`{"result":0}`,
			models.RunStatusCompleted,
			false,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.NewRunInput(cltest.JSONFromString(t, test.input))
			adapter := adapters.Copy{CopyPath: test.copyPath}
			result := adapter.Perform(input, nil)
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
