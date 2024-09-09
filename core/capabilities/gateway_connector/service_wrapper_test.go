package gatewayconnector

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	chainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/mock"
)

func generateConfig(addr common.Address) (chainlink.GeneralConfig, error) {
	return chainlink.GeneralConfigOpts{
		Config: chainlink.Config{
			Core: toml.Core{
				Capabilities: toml.Capabilities{
					GatewayConnector: toml.GatewayConnector{
						ChainIDForNodeKey:         ptr("1"),
						NodeAddress:               ptr(addr.Hex()),
						DonID:                     ptr("5"),
						WSHandshakeTimeoutMillis:  ptr[uint32](100),
						AuthMinChallengeLen:       ptr[int](0),
						AuthTimestampToleranceSec: ptr[uint32](10),
						Gateways:                  []toml.ConnectorGateway{{ID: ptr("example_gateway"), URL: ptr("wss://localhost:8081/node")}},
					},
				},
			},
		},
	}.New()
}

// Unit test that creates the ServiceWrapper object and then calls Start() can Close() on it.
// Take inspiration from functions/plugin_test.go and functions/connector_handler_test.go on how to mock the dependencies.
//
// Test valid NodeAddress and an invalid one (i.e. key doesn't exit).

func TestGatewayConnectorServiceWrapper(t *testing.T) {
	t.Parallel()

	logger := logger.TestLogger(t)
	_, addr := testutils.NewPrivateKeyAndAddress(t)

	config, err := generateConfig(addr)
	ethKeystore := ksmocks.NewEth(t)
	ethKeystore.On("EnabledKeysForChain", mock.Anything).Return([]ethkey.KeyV2{{Address: addr}})

	gc := config.Capabilities().GatewayConnector()
	handler := NewGatewayConnectorServiceWrapper(&gc, ethKeystore, logger)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, handler.Close())
	})

	t.Run("Start & Stop Success", func(t *testing.T) {
		ctx := testutils.Context(t)

		err := handler.Start(ctx)
		require.NoError(t, err)
		err = handler.Close()
		require.NoError(t, err)
	})
}

func TestGatewayConnectorServiceWrapperConfigError(t *testing.T) {
	t.Parallel()

	logger := logger.TestLogger(t)
	_, addr := testutils.NewPrivateKeyAndAddress(t)

	config, err := generateConfig(addr)
	ethKeystore := ksmocks.NewEth(t)
	_, addr2 := testutils.NewPrivateKeyAndAddress(t)

	ethKeystore.On("EnabledKeysForChain", mock.Anything).Return([]ethkey.KeyV2{{Address: addr2}})

	gc := config.Capabilities().GatewayConnector()
	handler := NewGatewayConnectorServiceWrapper(&gc, ethKeystore, logger)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, handler.Close())
	})

	t.Run("Start Error", func(t *testing.T) {
		ctx := testutils.Context(t)
		err := handler.Start(ctx)
		require.Error(t, err)
	})
}

func ptr[T any](t T) *T { return &t }
