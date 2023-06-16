package connector_test

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const defaultConfig = `
NodeAddress = "0x68902d681c28119f9b2531473a417088bf008e59"
DonId = "example_don"

[[Gateways]]
Id = "example_gateway"
URL = "ws://localhost:8081/node"

[[Gateways]]
Id = "another_one"
URL = "wss://example.com:8090/node_endpoint"
`

func parseTOMLConfig(t *testing.T, tomlConfig string) *connector.ConnectorConfig {
	var cfg connector.ConnectorConfig
	err := toml.Unmarshal([]byte(tomlConfig), &cfg)
	require.NoError(t, err)
	return &cfg
}

func newTestConnector(t *testing.T, config *connector.ConnectorConfig) (connector.GatewayConnector, *mocks.Signer, *mocks.GatewayConnectorHandler) {
	signer := mocks.NewSigner(t)
	handler := mocks.NewGatewayConnectorHandler(t)
	clock := utils.NewFixedClock(time.Now())
	connector, err := connector.NewGatewayConnector(config, signer, handler, clock, logger.TestLogger(t))
	require.NoError(t, err)
	return connector, signer, handler
}

func TestGatewayConnector_NewGatewayConnector_ValidConfig(t *testing.T) {
	t.Parallel()

	tomlConfig := parseTOMLConfig(t, `
NodeAddress = "0x68902d681c28119f9b2531473a417088bf008e59"
DonId = "example_don"

[[Gateways]]
Id = "example_gateway"
URL = "ws://localhost:8081/node"
`)

	newTestConnector(t, tomlConfig)
}

func TestGatewayConnector_NewGatewayConnector_InvalidConfig(t *testing.T) {
	t.Parallel()

	invalidCases := map[string]string{
		"invalid DON ID": `
NodeAddress = "0x68902d681c28119f9b2531473a417088bf008e59"
DonId = ""
`,
		"invalid node address": `
NodeAddress = "2531473a417088bf008e59"
DonId = "example_don"
`,
		"duplicate gateway ID": `
NodeAddress = "0x68902d681c28119f9b2531473a417088bf008e59"
DonId = "example_don"

[[Gateways]]
Id = "example_gateway"
URL = "ws://localhost:8081/node"

[[Gateways]]
Id = "example_gateway"
URL = "ws://localhost:8081/node"
`,
	}

	signer := mocks.NewSigner(t)
	handler := mocks.NewGatewayConnectorHandler(t)
	clock := utils.NewFixedClock(time.Now())
	for name, config := range invalidCases {
		config := config
		t.Run(name, func(t *testing.T) {
			_, err := connector.NewGatewayConnector(parseTOMLConfig(t, config), signer, handler, clock, logger.TestLogger(t))
			require.Error(t, err)
		})
	}
}

func TestGatewayConnector_NewAuthHeader_SignerError(t *testing.T) {
	t.Parallel()

	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))
	signer.On("Sign", mock.Anything).Return(nil, errors.New("cannot sign"))

	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)
	_, err = connector.NewAuthHeader(url)
	require.Error(t, err)
}
