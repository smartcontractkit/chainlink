package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store/models"
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
		{"existing path", `{"high": "11850.00", "last": "11779.99"}`, []string{"last"}, "11779.99", false, false},
		{"nonexistent path", `{"high": "11850.00", "last": "11779.99"}`, []string{"doesnotexist"}, "", true, false},
		{"double nonexistent path", `{"high": "11850.00", "last": "11779.99"}`, []string{"no", "really"}, "", true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := models.RunResultWithValue(test.value)
			adapter := adapters.JsonParse{Path: test.path}
			result := adapter.Perform(input, nil)
			val, err := result.Value()
			assert.Equal(t, test.want, val)
			if test.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			if test.wantResultError {
				assert.NotNil(t, result.GetError())
			} else {
				assert.Nil(t, result.GetError())
			}
		})
	}
}
