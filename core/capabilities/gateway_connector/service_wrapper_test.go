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
	ethKeystore.On("EnabledKeysForChain", mock.Anything, mock.Anything).Return([]ethkey.KeyV2{{Address: keystoreAddr}}, nil)

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
	// panic: runtime error: invalid memory address or nil pointer dereference
	// [signal SIGSEGV: segmentation violation code=0x1 addr=0x0 pc=0x1008fbb6b]

	// goroutine 280 [running]:
	// github.com/ethereum/go-ethereum/crypto.Sign({0xc001a1b880, 0x20, 0x20}, 0x0)
	// 	/Users/davidorchard/go/pkg/mod/github.com/ethereum/go-ethereum@v1.13.8/crypto/signature_cgo.go:60 +0x10b
	// github.com/smartcontractkit/chainlink/v2/core/utils.GenerateEthSignature(0x0, {0xc000d128f0, 0xc4, 0xd0})
	// 	/Users/davidorchard/code/chainlink/core/utils/eth_signatures.go:46 +0x6e
	// github.com/smartcontractkit/chainlink/v2/core/services/gateway/common.SignData(0x0, {0xc000bc1ef0?, 0x102412600?, 0xc000d12801?})
	// 	/Users/davidorchard/code/chainlink/core/services/gateway/common/utils.go:45 +0xdb
	// github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector.(*connectorSigner).Sign(0xc001a6c280?, {0xc000bc1ef0?, 0xc001a3d140?, 0x101bdfaae?})
	// 	/Users/davidorchard/code/chainlink/core/capabilities/gateway_connector/service_wrapper.go:53 +0x1c
	// github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector.(*gatewayConnector).NewAuthHeader(0xc00194bc20, 0x0?)
	// 	/Users/davidorchard/code/chainlink/core/services/gateway/connector/connector.go:256 +0x28f
	// github.com/smartcontractkit/chainlink/v2/core/services/gateway/network.(*webSocketClient).Connect(0xc00194d530, {0x1027f9260, 0xc001a6e140}, 0xc001950d80)
	// 	/Users/davidorchard/code/chainlink/core/services/gateway/network/wsclient.go:42 +0x50
	// github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector.(*gatewayConnector).reconnectLoop(0xc00194bc20, 0xc00195e780)
	// 	/Users/davidorchard/code/chainlink/core/services/gateway/connector/connector.go:195 +0x14a
	// created by github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector.(*gatewayConnector).Start.func1 in goroutine 277
	// 	/Users/davidorchard/code/chainlink/core/services/gateway/connector/connector.go:226 +0xd4
	// FAIL	github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector	3.241s
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
