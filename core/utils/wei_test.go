package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWei(t *testing.T) {
	for _, tt := range []struct {
		input string
		exp   string
	}{
		{"0", "0"},
		{"1", "1"},
		{"1000", "1 kwei"},
		{"1000000", "1 mwei"},
		{"1000000000", "1 gwei"},
		{"1000000000000", "1 micro"},
		{"200000000000000", "200 micro"},
		{"200 micro", "200 micro"},
		{"0.2 milli", "200 micro"},
		{"1000000000000000", "1 milli"},
		{"1000000000000000000", "1 ether"},
		{"1000000000000000000000", "1 kether"},
		{"1000000000000000000000000", "1 mether"},
		{"1000000000000000000000000000", "1 gether"},
		{"1000000000000000000000000000000", "1 tether"},
		{"1100000000000000000000000000000", "1.1 tether"},
		//TODO more cases, errors
	} {
		t.Run(tt.input, func(t *testing.T) {
			var w Wei
			err := w.UnmarshalText([]byte(tt.input))
			require.NoError(t, err)
			b, err := w.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.exp, string(b))
		})
	}
}
