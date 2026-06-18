package gatewayconnector_test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/ethkey"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	gatewayconnector "github.com/smartcontractkit/chainlink/v2/core/capabilities/gateway_connector"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	evmtestutils "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestGatewayConnectorConfigFromTOML(t *testing.T) {
	t.Parallel()

	cfg := newTestGatewayConnectorConfig(t)
	gc := cfg.Capabilities().GatewayConnector()
	assert.Equal(t, "0x68902d681c28119f9b2531473a417088bf008e59", gc.NodeAddress())
	assert.Equal(t, "example_don", gc.DonID())
	assert.Equal(t, uint32(100), gc.WSHandshakeTimeoutMillis())
	assert.Equal(t, 10, gc.AuthMinChallengeLen())
	assert.Equal(t, uint32(5), gc.AuthTimestampToleranceSec())

	gateways := gc.Gateways()
	require.Len(t, gateways, 2)
	assert.Equal(t, "example_gateway", gateways[0].ID())
	assert.Equal(t, "example_gateway_don", gateways[0].DonID())
	assert.Equal(t, "wss://localhost:8081/node", gateways[0].URL())
	assert.Equal(t, "another_gateway", gateways[1].ID())
	assert.Equal(t, "another_gateway_don", gateways[1].DonID())
	assert.Equal(t, "wss://example.com:8090/node", gateways[1].URL())
}

func newTestGatewayConnectorConfig(t *testing.T) chainlink.GeneralConfig {
	t.Helper()

	cfg, err := chainlink.GeneralConfigOpts{
		Config: chainlink.Config{
			Core: toml.Core{
				Capabilities: toml.Capabilities{
					GatewayConnector: toml.GatewayConnector{
						ChainIDForNodeKey:         new("1"),
						NodeAddress:               new("0x68902d681c28119f9b2531473a417088bf008e59"),
						DonID:                     new("example_don"),
						WSHandshakeTimeoutMillis:  new(uint32(100)),
						AuthMinChallengeLen:       new(10),
						AuthTimestampToleranceSec: new(uint32(5)),
						Gateways: []toml.ConnectorGateway{
							{
								ID:    new("example_gateway"),
								DonID: new("example_gateway_don"),
								URL:   new("wss://localhost:8081/node"),
							},
							{
								ID:    new("another_gateway"),
								DonID: new("another_gateway_don"),
								URL:   new("wss://example.com:8090/node"),
							},
						},
					},
				},
			},
		},
	}.New()
	require.NoError(t, err)
	return cfg
}

// fakeOrderedKeyProvider is a simple fake implementation for testing
type fakeOrderedKeyProvider struct {
	keys    []ethkey.KeyV2
	err     error
	chainID *big.Int
}

func (f *fakeOrderedKeyProvider) ListKeys(ctx context.Context, chainID *big.Int, opts *keystore.ListKeysOptions) ([]ethkey.KeyV2, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.chainID != nil && f.chainID.Cmp(chainID) != 0 {
		return nil, assert.AnError
	}
	return f.keys, nil
}

func generateWrapper(t *testing.T, privateKey *ecdsa.PrivateKey, keystoreKey *ecdsa.PrivateKey) (*gatewayconnector.ServiceWrapper, error) {
	lggr := logger.Test(t)
	privateKeyV2 := ethkey.FromPrivateKey(privateKey)
	addr := privateKeyV2.Address
	keystoreKeyV2 := ethkey.FromPrivateKey(keystoreKey)

	config, err := chainlink.GeneralConfigOpts{
		Config: chainlink.Config{
			Core: toml.Core{
				Capabilities: toml.Capabilities{
					GatewayConnector: toml.GatewayConnector{
						ChainIDForNodeKey:         new("1"),
						NodeAddress:               new(addr.Hex()),
						DonID:                     new("5"),
						WSHandshakeTimeoutMillis:  new(uint32(100)),
						AuthMinChallengeLen:       new(0),
						AuthTimestampToleranceSec: new(uint32(10)),
						Gateways:                  []toml.ConnectorGateway{{ID: new("example_gateway"), URL: new("wss://localhost:8081/node")}},
					},
				},
			},
		},
	}.New()
	require.NoError(t, err)

	ethKeystore := &keystest.FakeChainStore{Addresses: keystest.Addresses{keystoreKeyV2.Address}}
	gc := config.Capabilities().GatewayConnector()
	wrapper := gatewayconnector.NewGatewayConnectorServiceWrapper(gc, ethKeystore, nil, big.NewInt(1), clockwork.NewFakeClock(), lggr)
	return wrapper, nil
}

func setupAutoDiscoverTest(
	t *testing.T,
	nodeAddress *string,
	orderedKeyProvider gatewayconnector.OrderedKeyProvider,
	keystoreAddresses []ethkey.KeyV2,
	addressToPrivateKey map[string]*ecdsa.PrivateKey,
) (*gatewayconnector.ServiceWrapper, error) {
	lggr := logger.Test(t)

	config, err := chainlink.GeneralConfigOpts{
		Config: chainlink.Config{
			Core: toml.Core{
				Capabilities: toml.Capabilities{
					GatewayConnector: toml.GatewayConnector{
						ChainIDForNodeKey:         new("1"),
						NodeAddress:               nodeAddress,
						DonID:                     new("5"),
						WSHandshakeTimeoutMillis:  new(uint32(100)),
						AuthMinChallengeLen:       new(0),
						AuthTimestampToleranceSec: new(uint32(10)),
						Gateways:                  []toml.ConnectorGateway{{ID: new("example_gateway"), URL: new("wss://localhost:8081/node")}},
					},
				},
			},
		},
	}.New()
	require.NoError(t, err)

	var ethKeystore keys.Store
	if addressToPrivateKey != nil {
		ethKeystore = evmtestutils.NewSigningKeystore(addressToPrivateKey, keystoreAddresses)
	} else {
		addresses := make(keystest.Addresses, len(keystoreAddresses))
		for i, key := range keystoreAddresses {
			addresses[i] = key.Address
		}
		ethKeystore = &keystest.FakeChainStore{Addresses: addresses}
	}

	return gatewayconnector.NewGatewayConnectorServiceWrapper(
		config.Capabilities().GatewayConnector(),
		ethKeystore,
		orderedKeyProvider,
		big.NewInt(1),
		clockwork.NewFakeClock(),
		lggr,
	), nil
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

func TestGatewayConnectorServiceWrapper_AutoDiscoverNodeAddress(t *testing.T) {
	t.Parallel()

	key1, _ := testutils.NewPrivateKeyAndAddress(t)
	key2, _ := testutils.NewPrivateKeyAndAddress(t)
	keystoreKey, _ := testutils.NewPrivateKeyAndAddress(t)

	key1V2 := ethkey.FromPrivateKey(key1)
	key2V2 := ethkey.FromPrivateKey(key2)
	keystoreKeyV2 := ethkey.FromPrivateKey(keystoreKey)

	orderedKeyProvider := &fakeOrderedKeyProvider{
		keys:    []ethkey.KeyV2{key1V2, key2V2},
		chainID: big.NewInt(1),
	}

	addressToKey := map[string]*ecdsa.PrivateKey{
		key1V2.Address.Hex():        key1,
		keystoreKeyV2.Address.Hex(): keystoreKey,
	}
	wrapper, err := setupAutoDiscoverTest(t, nil, orderedKeyProvider, []ethkey.KeyV2{key1V2, keystoreKeyV2}, addressToKey)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	err = wrapper.Start(ctx)
	require.NoError(t, err)

	testData := []byte("test")
	wrapperSignature, err := wrapper.Sign(ctx, testData)
	require.NoError(t, err, "Sign should succeed with auto-discovered address")

	recoveredAddr, err := utils.GetSignersEthAddress(testData, wrapperSignature)
	require.NoError(t, err, "Should be able to recover address from signature")
	assert.Equal(t, key1V2.Address, recoveredAddr, "Signature should be from key1V2 (the discovered address)")

	t.Cleanup(func() {
		assert.NoError(t, wrapper.Close())
	})
}

func TestGatewayConnectorServiceWrapper_AutoDiscover(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		nodeAddress        *string
		orderedKeyProvider gatewayconnector.OrderedKeyProvider
		keystoreKeyCount   int
		wantErr            bool
		errContains        string
		expectedErr        error
	}{
		{
			name:               "no provider",
			nodeAddress:        nil,
			orderedKeyProvider: nil,
			keystoreKeyCount:   2,
			wantErr:            true,
			errContains:        "NodeAddress must be configured when ordered key provider is not available",
		},
		{
			name: "no keys",
			orderedKeyProvider: &fakeOrderedKeyProvider{
				keys:    []ethkey.KeyV2{},
				chainID: big.NewInt(1),
			},
			keystoreKeyCount: 1,
			wantErr:          true,
			errContains:      "no enabled keys found for auto-discovery",
		},
		{
			name: "provider error",
			orderedKeyProvider: &fakeOrderedKeyProvider{
				keys:    nil,
				err:     assert.AnError,
				chainID: big.NewInt(1),
			},
			keystoreKeyCount: 1,
			wantErr:          true,
			expectedErr:      assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keystoreKeysV2 := make([]ethkey.KeyV2, tt.keystoreKeyCount)
			for i := 0; i < tt.keystoreKeyCount; i++ {
				key, _ := testutils.NewPrivateKeyAndAddress(t)
				keystoreKeysV2[i] = ethkey.FromPrivateKey(key)
			}

			wrapper, err := setupAutoDiscoverTest(t, tt.nodeAddress, tt.orderedKeyProvider, keystoreKeysV2, nil)
			require.NoError(t, err)

			ctx := testutils.Context(t)
			err = wrapper.Start(ctx)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, wrapper.Close())
				})
			}
		})
	}
}
