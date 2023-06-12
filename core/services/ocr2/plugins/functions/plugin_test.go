package functions_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions"
)

func TestNewConnector_Success(t *testing.T) {
	t.Parallel()
	address := "0x00000000DE801ceE9471ADf23370c48b011f82a6"

	gwcCfg := &connector.ConnectorConfig{
		NodeAddress: address,
		DonId:       "my_don",
	}
	chainID := big.NewInt(80001)
	ethKeystore := ksmocks.NewEth(t)
	ethKeystore.On("EnabledKeysForChain", mock.Anything).Return([]ethkey.KeyV2{{Address: common.HexToAddress(address)}}, nil)
	_, err := functions.NewConnector(gwcCfg, ethKeystore, chainID, logger.TestLogger(t))
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
	ethKeystore.On("EnabledKeysForChain", mock.Anything).Return([]ethkey.KeyV2{{Address: common.HexToAddress(addresses[1])}}, nil)
	_, err := functions.NewConnector(gwcCfg, ethKeystore, chainID, logger.TestLogger(t))
	require.Error(t, err)
}
