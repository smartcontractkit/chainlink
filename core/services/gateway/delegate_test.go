package gateway_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

func TestDelegate_JobSpecValidator(t *testing.T) {
	t.Parallel()

	var tt = []struct {
		name  string
		toml  string
		valid bool
	}{
		{
			"valid spec",
			`
type = "gateway"
schemaVersion = 1
name = "The Best Gateway Job Ever!"
[gatewayConfig.NodeServerConfig]
Port = 666
`,
			true,
		},
		{
			"parse error",
			`
cantparsethis{{{{
`,
			false,
		},
		{
			"invalid job type",
			`
type = "gatez wayz"
schemaVersion = 1
`,
			false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := gateway.ValidatedGatewaySpec(tc.toml)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
