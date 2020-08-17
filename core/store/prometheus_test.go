package store

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"

	"github.com/stretchr/testify/require"
)

func Test_approximateFloat64(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      float64
		wantError bool
	}{
		{"zero", "0", 0, false},
		{"small", "1", 0.000000000000000001, false},
		{"rounding", "12345678901234567890", 12.345678901234567, false},
		{"large", "123456789012345678901234567890", 123456789012.34567, false},
		{"extreme", "1234567890123456789012345678901234567890123456789012345678901234567890", 1.2345678901234568e+51, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eth := assets.NewEth(0)
			eth.SetString(test.input, 10)
			float, err := approximateFloat64(eth)
			require.NoError(t, err)
			require.Equal(t, test.want, float)
		})
	}
}
