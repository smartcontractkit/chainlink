package connector_test

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const defaultConfig = `
NodeAddress = "0x68902d681c28119f9b2531473a417088bf008e59"
DonId = "example_don"
AuthMinChallengeLen = 10
AuthTimestampToleranceSec = 5

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

func newTestConnector(t *testing.T, config *connector.ConnectorConfig, now time.Time) (connector.GatewayConnector, *mocks.Signer, *mocks.GatewayConnectorHandler) {
	signer := mocks.NewSigner(t)
	handler := mocks.NewGatewayConnectorHandler(t)
	clock := utils.NewFixedClock(now)
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

	newTestConnector(t, tomlConfig, time.Now())
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
URL = "ws://localhost:8081/a"

[[Gateways]]
Id = "example_gateway"
URL = "ws://localhost:8081/b"
`,
		"duplicate gateway URL": `
NodeAddress = "0x68902d681c28119f9b2531473a417088bf008e59"
DonId = "example_don"

[[Gateways]]
Id = "gateway_A"
URL = "ws://localhost:8081/node"

[[Gateways]]
Id = "gateway_B"
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

func TestGatewayConnector_CleanStartAndClose(t *testing.T) {
	t.Parallel()

	connector, signer, handler := newTestConnector(t, parseTOMLConfig(t, defaultConfig), time.Now())
	handler.On("Start", mock.Anything).Return(nil)
	handler.On("Close").Return(nil)
	signer.On("Sign", mock.Anything).Return(nil, errors.New("cannot sign"))
	require.NoError(t, connector.Start(testutils.Context(t)))
	require.NoError(t, connector.Close())
}

func TestGatewayConnector_NewAuthHeader_SignerError(t *testing.T) {
	t.Parallel()

	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig), time.Now())
	signer.On("Sign", mock.Anything).Return(nil, errors.New("cannot sign"))

	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)
	_, err = connector.NewAuthHeader(url)
	require.Error(t, err)
}

func TestGatewayConnector_NewAuthHeader_Success(t *testing.T) {
	t.Parallel()

	testSignature := make([]byte, network.HandshakeSignatureLen)
	testSignature[1] = 0xfa
	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig), time.Now())
	signer.On("Sign", mock.Anything).Return(testSignature, nil)
	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)

	header, err := connector.NewAuthHeader(url)
	require.NoError(t, err)
	require.Equal(t, testSignature, header[len(header)-65:])
}

func TestGatewayConnector_ChallengeResponse(t *testing.T) {
	t.Parallel()

	testSignature := make([]byte, network.HandshakeSignatureLen)
	testSignature[1] = 0xfa
	now := time.Now()
	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig), now)
	signer.On("Sign", mock.Anything).Return(testSignature, nil)
	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)

	challenge := network.ChallengeElems{
		Timestamp:      uint32(now.Unix()),
		GatewayId:      "example_gateway",
		ChallengeBytes: []byte("1234567890"),
	}

	// valid
	signature, err := connector.ChallengeResponse(url, network.PackChallenge(&challenge))
	require.NoError(t, err)
	require.Equal(t, testSignature, signature)

	// invalid timestamp
	badChallenge := challenge
	badChallenge.Timestamp += 100
	_, err = connector.ChallengeResponse(url, network.PackChallenge(&badChallenge))
	require.Equal(t, network.ErrAuthInvalidTimestamp, err)

	// too short
	badChallenge = challenge
	badChallenge.ChallengeBytes = []byte("aabb")
	_, err = connector.ChallengeResponse(url, network.PackChallenge(&badChallenge))
	require.Equal(t, network.ErrChallengeTooShort, err)

	// invalid GatewayId
	badChallenge = challenge
	badChallenge.GatewayId = "wrong"
	_, err = connector.ChallengeResponse(url, network.PackChallenge(&badChallenge))
	require.Equal(t, network.ErrAuthInvalidGateway, err)
}
