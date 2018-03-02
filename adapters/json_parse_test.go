package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestJsonParse_Perform(t *testing.T) {
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
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.RunResultWithValue(test.value)
			adapter := adapters.JSONParse{Path: test.path}
			result := adapter.Perform(input, nil)
			assert.Equal(t, test.want, result.Data.String())

			if test.wantResultError {
				assert.NotNil(t, result.GetError())
			} else {
				assert.Nil(t, result.GetError())
			}
		})
	}
}
