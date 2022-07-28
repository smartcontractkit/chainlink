package keeper

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestUpkeepIdentifer_String(t *testing.T) {
	for _, test := range []struct {
		name string
		id   string
		hex  string
	}{
		{"small", "10", "a"},
		{"large", "1000000000", "3b9aca00"},
		{"big", "5032485723458348569331745", "429ab990419450db80821"},
	} {
		t.Run(test.name, func(t *testing.T) {
			o, ok := new(big.Int).SetString(test.id, 10)
			if !ok {
				t.Errorf("%s failed to parse test integer", test.name)
				return
			}

			result := NewUpkeepIdentifier(utils.NewBig(o)).String()
			require.Equal(t, fmt.Sprintf("UPx%064s", test.hex), result)
		})
	}
}

func TestUpkeepIdentifer_Scan(t *testing.T) {
	for _, test := range []struct {
		name string
		id   string
		hex  string
	}{
		{"small", "10", "a"},
		{"large", "1000000000", "3b9aca00"},
		{"big", "5032485723458348569331745", "429ab990419450db80821"},
	} {
		t.Run(test.name, func(t *testing.T) {
			id := NewUpkeepIdentifier(utils.NewBigI(0))

			err := id.Scan(test.id)
			require.NoError(t, err)

			result := id.String()
			require.Equal(t, fmt.Sprintf("UPx%064s", test.hex), result)
		})
	}
}
