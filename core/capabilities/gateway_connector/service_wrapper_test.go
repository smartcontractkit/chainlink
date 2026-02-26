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
	// Verify chainID matches
	if f.chainID != nil && f.chainID.Cmp(chainID) != 0 {
		return nil, assert.AnError
	}
	return f.keys, nil
}

func generateWrapper(t *testing.T, privateKey *ecdsa.PrivateKey, keystoreKey *ecdsa.PrivateKey) (*gatewayconnector.ServiceWrapper, error) {
	logger := logger.Test(t)
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
	ethKeystore := &keystest.FakeChainStore{Addresses: keystest.Addresses{keystoreKeyV2.Address}}
	gc := config.Capabilities().GatewayConnector()
	chainID := big.NewInt(1)
	wrapper := gatewayconnector.NewGatewayConnectorServiceWrapper(gc, ethKeystore, nil, chainID, clockwork.NewFakeClock(), logger)
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

// setupAutoDiscoverTest creates a wrapper with auto-discovery configuration
func setupAutoDiscoverTest(
	t *testing.T,
	nodeAddress *string,
	orderedKeyProvider gatewayconnector.OrderedKeyProvider,
	keystoreAddresses []ethkey.KeyV2,
	addressToPrivateKey map[string]*ecdsa.PrivateKey,
) (*gatewayconnector.ServiceWrapper, error) {
	logger := logger.Test(t)
	chainID := big.NewInt(1)

	config, err := chainlink.GeneralConfigOpts{
		Config: chainlink.Config{
			Core: toml.Core{
				Capabilities: toml.Capabilities{
					GatewayConnector: toml.GatewayConnector{
						ChainIDForNodeKey:         ptr("1"),
						NodeAddress:               nodeAddress,
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
	require.NoError(t, err)

	var ethKeystore keys.Store
	if addressToPrivateKey != nil {
		// Use signing keystore that actually signs
		// addressToPrivateKey is already map[string]*ecdsa.PrivateKey
		ethKeystore = evmtestutils.NewSigningKeystore(addressToPrivateKey, keystoreAddresses)
	} else {
		// Use fake keystore for tests that don't need actual signing
		addresses := make(keystest.Addresses, len(keystoreAddresses))
		for i, key := range keystoreAddresses {
			addresses[i] = key.Address
		}
		ethKeystore = &keystest.FakeChainStore{Addresses: addresses}
	}
	gc := config.Capabilities().GatewayConnector()

	wrapper := gatewayconnector.NewGatewayConnectorServiceWrapper(
		gc,
		ethKeystore,
		orderedKeyProvider,
		chainID,
		clockwork.NewFakeClock(),
		logger,
	)

	return wrapper, nil
}

func TestGatewayConnectorServiceWrapper_AutoDiscoverNodeAddress(t *testing.T) {
	t.Parallel()

	key1, _ := testutils.NewPrivateKeyAndAddress(t)
	key2, _ := testutils.NewPrivateKeyAndAddress(t)
	keystoreKey, _ := testutils.NewPrivateKeyAndAddress(t)

	key1V2 := ethkey.FromPrivateKey(key1)
	key2V2 := ethkey.FromPrivateKey(key2)
	keystoreKeyV2 := ethkey.FromPrivateKey(keystoreKey)

	chainID := big.NewInt(1)
	orderedKeyProvider := &fakeOrderedKeyProvider{
		keys:    []ethkey.KeyV2{key1V2, key2V2},
		chainID: chainID,
	}

	// Create address to private key mapping for signing keystore
	addressToKey := map[string]*ecdsa.PrivateKey{
		key1V2.Address.Hex():        key1,
		keystoreKeyV2.Address.Hex(): keystoreKey,
	}
	wrapper, err := setupAutoDiscoverTest(t, nil, orderedKeyProvider, []ethkey.KeyV2{key1V2, keystoreKeyV2}, addressToKey)
	require.NoError(t, err)

	ctx := testutils.Context(t)
	err = wrapper.Start(ctx)
	require.NoError(t, err)

	// Verify that Sign() uses the discovered address (key1V2) by verifying the signature
	testData := []byte("test")
	wrapperSignature, err := wrapper.Sign(ctx, testData)
	require.NoError(t, err, "Sign should succeed with auto-discovered address")

	// Verify the signature was created with key1's address using utils
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
			name:        "no keys",
			nodeAddress: nil,
			orderedKeyProvider: &fakeOrderedKeyProvider{
				keys:    []ethkey.KeyV2{},
				chainID: big.NewInt(1),
			},
			keystoreKeyCount: 1,
			wantErr:          true,
			errContains:      "no enabled keys found for auto-discovery",
		},
		{
			name:        "provider error",
			nodeAddress: nil,
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
