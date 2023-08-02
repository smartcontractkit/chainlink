package gateway_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const defaultConfig = `
[nodeServerConfig]
Path = "/node"

[[dons]]
DonId = "my_don_1"
HandlerName = "dummy"

[[dons.members]]
Name = "example_node"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"

[[dons]]
DonId = "my_don_2"
HandlerName = "dummy"

[[dons.members]]
Name = "example_node"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
`

func TestConnectionManager_NewConnectionManager_ValidConfig(t *testing.T) {
	t.Parallel()

	tomlConfig := parseTOMLConfig(t, defaultConfig)

	_, err := gateway.NewConnectionManager(tomlConfig, utils.NewFixedClock(time.Now()), logger.TestLogger(t))
	require.NoError(t, err)
}

func TestConnectionManager_NewConnectionManager_InvalidConfig(t *testing.T) {
	t.Parallel()

	invalidCases := map[string]string{
		"duplicate DON ID": `
[[dons]]
DonId = "my_don"
[[dons]]
DonId = "my_don"
`,
		"duplicate node address": `
[[dons]]
DonId = "my_don"
[[dons.members]]
Name = "node_1"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
[[dons.members]]
Name = "node_2"
Address = "0x68902d681c28119f9b2531473a417088bf008e59"
`,
	}

	for name, config := range invalidCases {
		config := config
		t.Run(name, func(t *testing.T) {
			fullConfig := `
[nodeServerConfig]
Path = "/node"` + config
			_, err := gateway.NewConnectionManager(parseTOMLConfig(t, fullConfig), utils.NewFixedClock(time.Now()), logger.TestLogger(t))
			require.Error(t, err)
		})
	}
}

func TestConnectionManager_StartHandshake_TooShort(t *testing.T) {
	t.Parallel()

	mgr, err := gateway.NewConnectionManager(parseTOMLConfig(t, defaultConfig), utils.NewFixedClock(time.Now()), logger.TestLogger(t))
	require.NoError(t, err)

	_, _, err = mgr.StartHandshake([]byte("ab"))
	require.Error(t, err)
}
