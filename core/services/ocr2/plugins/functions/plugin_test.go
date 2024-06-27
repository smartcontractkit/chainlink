package functions_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	gfmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
	s4mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"
)

func TestNewConnector_Success(t *testing.T) {
	t.Parallel()
	keyV2, err := ethkey.NewV2()
	require.NoError(t, err)

	gwcCfg := &connector.ConnectorConfig{
		NodeAddress: keyV2.Address.String(),
		DonId:       "my_don",
	}
	chainID := big.NewInt(80001)
	ethKeystore := ksmocks.NewEth(t)
	s4Storage := s4mocks.NewStorage(t)
	allowlist := gfmocks.NewOnchainAllowlist(t)
	subscriptions := gfmocks.NewOnchainSubscriptions(t)
	rateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	ethKeystore.On("EnabledKeysForChain", mock.Anything).Return([]ethkey.KeyV2{keyV2}, nil)
	_, err = functions.NewConnector(gwcCfg, ethKeystore, chainID, s4Storage, allowlist, rateLimiter, subscriptions, *assets.NewLinkFromJuels(0), logger.TestLogger(t))
	require.NoError(t, err)
}

func TestNewConnector_NoKeyForConfiguredAddress(t *testing.T) {
	t.Parallel()
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
	allowlist := gfmocks.NewOnchainAllowlist(t)
	subscriptions := gfmocks.NewOnchainSubscriptions(t)
	rateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	ethKeystore.On("EnabledKeysForChain", mock.Anything).Return([]ethkey.KeyV2{{Address: common.HexToAddress(addresses[1])}}, nil)
	_, err = functions.NewConnector(gwcCfg, ethKeystore, chainID, s4Storage, allowlist, rateLimiter, subscriptions, *assets.NewLinkFromJuels(0), logger.TestLogger(t))
	require.Error(t, err)
}
