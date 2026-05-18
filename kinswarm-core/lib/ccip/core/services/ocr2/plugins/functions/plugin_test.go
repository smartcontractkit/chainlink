package functions_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	sfmocks "github.com/smartcontractkit/chainlink/v2/core/services/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	gfaMocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist/mocks"
	gfsMocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	s4mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"
)

func TestNewConnector_Success(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	keyV2, err := ethkey.NewV2()
	require.NoError(t, err)

	gwcCfg := &connector.ConnectorConfig{
		NodeAddress: keyV2.Address.String(),
		DonId:       "my_don",
	}
	chainID := big.NewInt(80001)
	ethKeystore := ksmocks.NewEth(t)
	s4Storage := s4mocks.NewStorage(t)
	allowlist := gfaMocks.NewOnchainAllowlist(t)
	subscriptions := gfsMocks.NewOnchainSubscriptions(t)
	rateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	listener := sfmocks.NewFunctionsListener(t)
	offchainTransmitter := sfmocks.NewOffchainTransmitter(t)
	ethKeystore.On("EnabledKeysForChain", mock.Anything, mock.Anything).Return([]ethkey.KeyV2{keyV2}, nil)
	config := &config.PluginConfig{
		GatewayConnectorConfig: gwcCfg,
	}
	_, _, err = functions.NewConnector(ctx, config, ethKeystore, chainID, s4Storage, allowlist, rateLimiter, subscriptions, listener, offchainTransmitter, logger.TestLogger(t))
	require.NoError(t, err)
}

func TestNewConnector_NoKeyForConfiguredAddress(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	addresses := []string{
		"0x00000000DE801ceE9471ADf23370c48b011f82a6",
		"0x11111111DE801ceE9471ADf23370c48b011f82a6",
	}

	gwcCfg := &connector.ConnectorConfig{
		NodeAddress: addresses[0],
		DonId:       "my_don",
	}
	chainID := big.NewInt(80001)
	ethKeystore := ksmocks.NewEth(t)
	s4Storage := s4mocks.NewStorage(t)
	allowlist := gfaMocks.NewOnchainAllowlist(t)
	subscriptions := gfsMocks.NewOnchainSubscriptions(t)
	rateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	listener := sfmocks.NewFunctionsListener(t)
	offchainTransmitter := sfmocks.NewOffchainTransmitter(t)
	ethKeystore.On("EnabledKeysForChain", mock.Anything, mock.Anything).Return([]ethkey.KeyV2{{Address: common.HexToAddress(addresses[1])}}, nil)
	config := &config.PluginConfig{
		GatewayConnectorConfig: gwcCfg,
	}
	_, _, err = functions.NewConnector(ctx, config, ethKeystore, chainID, s4Storage, allowlist, rateLimiter, subscriptions, listener, offchainTransmitter, logger.TestLogger(t))
	require.Error(t, err)
}
