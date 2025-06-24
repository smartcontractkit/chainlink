package connector

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	gatewaymocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	defaultConfig = `
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
	testMethod1 = "test_method_1"
	testMethod2 = "test_method_2"
)

func parseTOMLConfig(t *testing.T, tomlConfig string) *ConnectorConfig {
	var cfg ConnectorConfig
	err := toml.Unmarshal([]byte(tomlConfig), &cfg)
	require.NoError(t, err)
	return &cfg
}

func newTestConnector(t *testing.T, config *ConnectorConfig) (*gatewayConnector, *gatewaymocks.Signer, *gatewaymocks.GatewayConnectorHandler) {
	signer := gatewaymocks.NewSigner(t)
	handler := gatewaymocks.NewGatewayConnectorHandler(t)
	clock := clockwork.NewFakeClock()
	connector, err := NewGatewayConnector(config, signer, clock, logger.Test(t))
	require.NoError(t, err)
	require.NoError(t, connector.AddHandler(t.Context(), []string{testMethod1}, handler))
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

	signer := gatewaymocks.NewSigner(t)
	clock := clockwork.NewFakeClock()
	for name, config := range invalidCases {
		config := config
		t.Run(name, func(t *testing.T) {
			_, err := NewGatewayConnector(parseTOMLConfig(t, config), signer, clock, logger.Test(t))
			require.Error(t, err)
		})
	}
}

func TestGatewayConnector_CleanStartAndClose(t *testing.T) {
	t.Parallel()

	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))
	signer.On("Sign", mock.AnythingOfType("*context.cancelCtx"), mock.Anything).Return(nil, errors.New("cannot sign"))
	servicetest.Run(t, connector)
}

func TestGatewayConnector_NewAuthHeader_SignerError(t *testing.T) {
	t.Parallel()

	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))
	signer.On("Sign", mock.AnythingOfType("*context.cancelCtx"), mock.Anything).Return(nil, errors.New("cannot sign"))

	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)
	_, err = connector.NewAuthHeader(t.Context(), url)
	require.Error(t, err)
}

func TestGatewayConnector_NewAuthHeader_Success(t *testing.T) {
	t.Parallel()

	testSignature := make([]byte, network.HandshakeSignatureLen)
	testSignature[1] = 0xfa
	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))
	signer.On("Sign", mock.AnythingOfType("*context.cancelCtx"), mock.Anything).Return(testSignature, nil)
	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)

	header, err := connector.NewAuthHeader(t.Context(), url)
	require.NoError(t, err)
	require.Equal(t, testSignature, header[len(header)-65:])
}

func TestGatewayConnector_ChallengeResponse(t *testing.T) {
	t.Parallel()

	testSignature := make([]byte, network.HandshakeSignatureLen)
	testSignature[1] = 0xfa
	now := time.Now()
	connector, signer, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))
	signer.On("Sign", mock.AnythingOfType("*context.cancelCtx"), mock.Anything).Return(testSignature, nil)
	url, err := url.Parse("ws://localhost:8081/node")
	require.NoError(t, err)

	challenge := network.ChallengeElems{
		Timestamp:      uint32(now.Unix()),
		GatewayId:      "example_gateway",
		ChallengeBytes: []byte("1234567890"),
	}

	// valid
	signature, err := connector.ChallengeResponse(t.Context(), url, network.PackChallenge(&challenge))
	require.NoError(t, err)
	require.Equal(t, testSignature, signature)

	// invalid timestamp
	badChallenge := challenge
	badChallenge.Timestamp += 100
	_, err = connector.ChallengeResponse(t.Context(), url, network.PackChallenge(&badChallenge))
	require.Equal(t, network.ErrAuthInvalidTimestamp, err)

	// too short
	badChallenge = challenge
	badChallenge.ChallengeBytes = []byte("aabb")
	_, err = connector.ChallengeResponse(t.Context(), url, network.PackChallenge(&badChallenge))
	require.Equal(t, network.ErrChallengeTooShort, err)

	// invalid GatewayId
	badChallenge = challenge
	badChallenge.GatewayId = "wrong"
	_, err = connector.ChallengeResponse(t.Context(), url, network.PackChallenge(&badChallenge))
	require.Equal(t, network.ErrAuthInvalidGateway, err)
}

func TestGatewayConnector_AddHandler(t *testing.T) {
	t.Parallel()

	connector, _, _ := newTestConnector(t, parseTOMLConfig(t, defaultConfig))
	// testMethod1 already exists
	require.Error(t, connector.AddHandler(t.Context(), []string{testMethod1}, gatewaymocks.NewGatewayConnectorHandler(t)))
	require.NoError(t, connector.AddHandler(t.Context(), []string{testMethod2}, gatewaymocks.NewGatewayConnectorHandler(t)))
}
