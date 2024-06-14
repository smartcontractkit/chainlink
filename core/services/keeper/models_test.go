package keeper

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestUpkeepIdentifer_String(t *testing.T) {
	for _, test := range []struct {
		name string
		id   string
		hex  string
	}{
		{"small", "10", "UPx000000000000000000000000000000000000000000000000000000000000000a"},
		{"large", "1000000000", "UPx000000000000000000000000000000000000000000000000000000003b9aca00"},
		{"big", "5032485723458348569331745", "UPx0000000000000000000000000000000000000000000429ab990419450db80821"},
	} {
		t.Run(test.name, func(t *testing.T) {
			o, ok := new(big.Int).SetString(test.id, 10)
			if !ok {
				t.Errorf("%s failed to parse test integer", test.name)
				return
			}

			result := NewUpkeepIdentifier(utils.NewBig(o)).String()
			require.Equal(t, test.hex, result)
		})
	}
}
