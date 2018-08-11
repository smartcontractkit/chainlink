package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestNormalizedJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		object   interface{}
		output   string
		didError bool
	}{
		{"empty object", struct{}{}, "{}", false},
		{"empty array", []string{}, "[]", false},
		{"null", nil, "null", false},
		{"float", 1510599740287532257480015872.0, "1.510600e+27", false},
		{"bool", true, "true", false},
		{"string", "string", "\"string\"", false},
		{"array with one item", []string{"item"}, "[\"item\"]", false},
		{"map with one item", map[string]string{"item": "value"}, "{\"item\":\"value\"}", false},
		// See https://en.wikipedia.org/wiki/Precomposed_character
		{"string with decomposed characters",
			"\u0041\u030a\u0073\u0074\u0072\u006f\u0308\u006d",
			"\"\u00c5\u0073\u0074\u0072\u00f6\u006d\"",
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(test.object)
			assert.NoError(t, err)

			str, err := utils.NormalizedJSON(jsonBytes)

			cltest.ErrorPresence(t, test.didError, err)
			assert.Equal(t, test.output, str)
		})
	}
}
