package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestMap_Perform(t *testing.T) {
	tests := []struct {
		name            string
		value           string
		path            []string
		want            string
		wantError       bool
		wantResultError bool
	}{
		{"existing path", `{"high":"11850.00","last":"11779.99"}`, []string{"last"},
			`{"high":"11850.00","last":"11779.99","value":"11779.99"}`, false, false},
		{"nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"doesnotexist"},
			`{"high":"11850.00","last":"11779.99","value":null}`, true, false},
		{"double nonexistent path", `{"high":"11850.00","last":"11779.99"}`, []string{"no", "really"},
			`{"high":"11850.00","last":"11779.99","value":"{\"high\":\"11850.00\",\"last\":\"11779.99\"}"}`, true, true},
		{"array index path", `{"data":[{"availability":"0.99991"}]}`, []string{"data", "0", "availability"},
			`{"data":[{"availability":"0.99991"}],"value":"0.99991"}`, false, false},
		{"float value", `{"availability":0.99991}`, []string{"availability"},
			`{"availability":0.99991,"value":"0.99991"}`, false, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.RunResultWithData(test.value)
			log.Print(input)
			adapter := adapters.Map{Path: test.path}
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
