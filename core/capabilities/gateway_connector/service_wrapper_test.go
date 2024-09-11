package gatewayconnector_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

func generateWrapper(t *testing.T, privateKey *ecdsa.PrivateKey, keystoreKey *ecdsa.PrivateKey) (*gatewayconnector.ServiceWrapper, error) {
	logger := logger.TestLogger(t)
	privateKeyV2 := ethkey.FromPrivateKey(privateKey)
	addr := privateKeyV2.Address
	keystoreKeyV2 := ethkey.FromPrivateKey(keystoreKey)

	config, err := chainlink.GeneralConfigOpts{
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
	ethKeystore := ksmocks.NewEth(t)
	ethKeystore.On("EnabledKeysForChain", mock.Anything, mock.Anything).Return([]ethkey.KeyV2{keystoreKeyV2}, nil)
	gc := config.Capabilities().GatewayConnector()
	wrapper := gatewayconnector.NewGatewayConnectorServiceWrapper(gc, ethKeystore, clockwork.NewFakeClock(), logger)
	require.NoError(t, err)
	return wrapper, err
}

func TestGatewayConnectorServiceWrapper_CleanStartClose(t *testing.T) {
	t.Parallel()

	key, _ := testutils.NewPrivateKeyAndAddress(t)
	wrapper, err := generateWrapper(t, key, key)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	err = wrapper.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, wrapper.Close())
	})
}

func TestGatewayConnectorServiceWrapper_NonexistentKey(t *testing.T) {
	t.Parallel()

	key, _ := testutils.NewPrivateKeyAndAddress(t)
	keystoreKey, _ := testutils.NewPrivateKeyAndAddress(t)
	wrapper, err := generateWrapper(t, key, keystoreKey)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	err = wrapper.Start(ctx)
	require.Error(t, err)
}

func ptr[T any](t T) *T { return &t }
