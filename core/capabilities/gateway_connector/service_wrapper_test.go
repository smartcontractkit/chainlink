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

func generateWrapper(t *testing.T, addr common.Address, keystoreAddr common.Address) (*ServiceWrapper, error) {
	logger := logger.TestLogger(t)

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
	// github.com/stretchr/testify/mock.Arguments.Get(...)
	// /Users/davidorchard/go/pkg/mod/github.com/stretchr/testify@v1.9.0/mock/mock.go:900
	// github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks.(*Eth).EnabledKeysForChain(0xc0017a0f00, {0x1071e02d0, 0xc000776770}, 0xc0009148c0)
	ethKeystore.On("EnabledKeysForChain", mock.Anything, mock.Anything).Return([]ethkey.KeyV2{{Address: keystoreAddr}})

	gc := config.Capabilities().GatewayConnector()
	wrapper := NewGatewayConnectorServiceWrapper(gc, ethKeystore, logger)
	require.NoError(t, err)
	return wrapper, err

}

func TestGatewayConnectorServiceWrapper_CleanStartClose(t *testing.T) {
	t.Parallel()

	_, addr := testutils.NewPrivateKeyAndAddress(t)
	wrapper, err := generateWrapper(t, addr, addr)

	require.NoError(t, err)

	ctx := testutils.Context(t)

	err = wrapper.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, wrapper.Close())
	})
}

func TestGatewayConnectorServiceWrapper_NonexistantKey(t *testing.T) {
	t.Parallel()

	_, addr := testutils.NewPrivateKeyAndAddress(t)
	_, keystoreAddr := testutils.NewPrivateKeyAndAddress(t)
	wrapper, err := generateWrapper(t, addr, keystoreAddr)

	require.NoError(t, err)
	ctx := testutils.Context(t)
	err = wrapper.Start(ctx)
	require.Error(t, err)
}

func ptr[T any](t T) *T { return &t }
