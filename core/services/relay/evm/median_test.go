package evm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	configmocks "github.com/smartcontractkit/chainlink-evm/pkg/config/mocks"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	logmocks "github.com/smartcontractkit/chainlink/v2/common/log/mocks"
	medianconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/median/config"
)

func TestNewMedianProvider(t *testing.T) {
	lggr := logger.Test(t)

	chain := mocks.NewChain(t)
	chainID := testutils.NewRandomEVMChainID()
	chain.EXPECT().ID().Return(chainID)

	contractID := testutils.NewAddress()
	transmitterAddr := testutils.NewAddress()

	keystore := &keystest.FakeChainStore{Addresses: keystest.Addresses{transmitterAddr}}
	relayer := &Relayer{
		lggr:        logger.Sugared(lggr),
		chain:       chain,
		evmKeystore: keystore,
	}

	pargs := commontypes.PluginArgs{}

	t.Run("wrong chainID", func(t *testing.T) {
		relayConfigBadChainID := config.RelayConfig{}
		rc, err2 := json.Marshal(&relayConfigBadChainID)
		rargs2 := commontypes.RelayArgs{ContractID: contractID.String(), RelayConfig: rc}
		require.NoError(t, err2)
		_, err2 = relayer.NewMedianProvider(testutils.Context(t), rargs2, pargs)
		assert.ErrorContains(t, err2, "chain id in spec does not match")
	})

	t.Run("invalid contractID", func(t *testing.T) {
		relayConfig := config.RelayConfig{ChainID: big.New(chainID)}
		rc, err2 := json.Marshal(&relayConfig)
		require.NoError(t, err2)
		rargsBadContractID := commontypes.RelayArgs{ContractID: "NotAContractID", RelayConfig: rc}
		_, err2 = relayer.NewMedianProvider(testutils.Context(t), rargsBadContractID, pargs)
		assert.ErrorContains(t, err2, "invalid contractID")
	})

	t.Run("plugin config contains gas limit", func(t *testing.T) {
		mocks, relayer := setupMocksAndRelayer(t)
		relayer.evmKeystore = keystore
		relayer.lggr = logger.Sugared(lggr)

		mocks.Chain.EXPECT().ID().Return(chainID)
		mocks.EVM.EXPECT().ChainID().Return(chainID)
		mocks.Poller.EXPECT().RegisterFilter(mock.Anything, mock.Anything).Return(nil)
		mocks.EVM.EXPECT().ChainType().Return("")
		mocks.Chain.EXPECT().LogBroadcaster().Return(logmocks.NewBroadcaster(t))

		mockGasEstimator := configmocks.NewGasEstimator(t)
		mocks.EVM.EXPECT().GasEstimator().Return(mockGasEstimator)
		mockGasEstimator.EXPECT().LimitDefault().Return(uint64(500000))
		mockGasEstimator.EXPECT().LimitJobType().Return(&txmgr.TestLimitJobTypeConfig{})

		relayConfig := config.RelayConfig{
			ChainID:                big.New(chainID),
			EffectiveTransmitterID: null.StringFrom(transmitterAddr.String()),
			SendingKeys:            []string{transmitterAddr.String()},
		}
		rc, err := json.Marshal(&relayConfig)
		require.NoError(t, err)

		pluginCfg := medianconfig.PluginConfig{GasLimit: uint32Ptr(99)}
		pluginConfigBytes, err := json.Marshal(pluginCfg)
		require.NoError(t, err)
		pargs.PluginConfig = pluginConfigBytes

		_, err = relayer.NewMedianProvider(testutils.Context(t), commontypes.RelayArgs{ContractID: contractID.String(), RelayConfig: rc}, pargs)
		assert.NoError(t, err)
	})
}

func uint32Ptr(v uint32) *uint32 {
	return &v
}
