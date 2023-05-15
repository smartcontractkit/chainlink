package gateway_test

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

func parseTOMLConfig(t *testing.T, tomlConfig string) *gateway.GatewayConfig {
	var cfg gateway.GatewayConfig
	err := toml.Unmarshal([]byte(tomlConfig), &cfg)
	require.NoError(t, err)
	return &cfg
}

func TestGateway_NewGatewayFromConfig_ValidConfig(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[[dons]]
DonId = "my_don_1"
HandlerName = "dummy"

[[dons]]
DonId = "my_don_2"
HandlerName = "dummy"
`

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.NoError(t, err)
}

func TestGateway_NewGatewayFromConfig_DuplicateID(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[[dons]]
DonId = "my_don"
HandlerName = "dummy"

[[dons]]
DonId = "my_don"
HandlerName = "dummy"
`

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_InvalidHandler(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[[dons]]
DonId = "my_don"
HandlerName = "no_such_handler"
`

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_MissingID(t *testing.T) {
	t.Parallel()

	tomlConfig := `
[[dons]]
HandlerName = "dummy"
SomeOtherField = "abcd"
`

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.Error(t, err)
}
