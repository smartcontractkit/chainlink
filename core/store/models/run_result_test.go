package models_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRunResult_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		json        string
		want        string
		wantErrored bool
	}{
		{"string", `{"result": "100", "other": "101"}`, "100", false},
		{"integer", `{"result": 100}`, "", true},
		{"float", `{"result": 100.01}`, "", true},
		{"boolean", `{"result": true}`, "", true},
		{"null", `{"result": null}`, "", true},
		{"no key", `{"other": 100}`, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data models.JSON
			assert.NoError(t, json.Unmarshal([]byte(test.json), &data))
			rr := models.RunResult{Data: data}

			val, err := rr.ResultString()
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}
