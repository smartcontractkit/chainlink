package config

import (
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/chains"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestGeneralConfig_Defaults(t *testing.T) {
	config := NewGeneralConfig()
	assert.Equal(t, uint64(10), config.BlockBackfillDepth())
	assert.Equal(t, new(url.URL), config.BridgeResponseURL())
	assert.Equal(t, big.NewInt(1), config.ChainID())
	assert.Equal(t, false, config.EthereumDisabled())
	assert.Equal(t, false, config.FeatureExternalInitiators())
	assert.Equal(t, 15*time.Minute, config.SessionTimeout().Duration())
}

func TestGeneralConfig_sessionSecret(t *testing.T) {
	t.Parallel()
	config := NewGeneralConfig()
	// config.Set("ROOT", path.Join("/tmp/chainlink_test", "TestConfig_sessionSecret"))
	// err := os.MkdirAll(config.RootDir(), os.FileMode(0770))
	// require.NoError(t, err)
	// defer os.RemoveAll(config.RootDir())

	initial, err := config.SessionSecret()
	require.NoError(t, err)
	require.NotEqual(t, "", initial)
	require.NotEqual(t, "clsession_test_secret", initial)

	second, err := config.SessionSecret()
	require.NoError(t, err)
	require.Equal(t, initial, second)
}

func newEVMConfigWithChainID(id string) *evmConfig {
	gcfg := NewGeneralConfig()
	gcfg.(*generalConfig).viper.Set("ETH_CHAIN_ID", id)
	config := NewEVMConfig(gcfg)
	return config.(*evmConfig)
}

func newEVMConfig(f func(c *generalConfig)) *evmConfig {
	gcfg := NewGeneralConfig()
	f(gcfg.(*generalConfig))
	config := NewEVMConfig(gcfg).(*evmConfig)
	return config
}

func TestEVMConfig_ChainSpecificConfig(t *testing.T) {
	t.Parallel()

	t.Run("with unknown chain ID returns generic defaults", func(t *testing.T) {
		config := newEVMConfigWithChainID("0")

		assert.Equal(t, chains.FallbackConfig.GasBumpThreshold, config.EvmGasBumpThreshold())
		assert.Equal(t, chains.FallbackConfig.GasBumpWei, *config.EvmGasBumpWei())
		assert.Equal(t, chains.FallbackConfig.GasPriceDefault, *config.EvmGasPriceDefault())
		assert.Equal(t, chains.FallbackConfig.MaxGasPriceWei, *config.EvmMaxGasPriceWei())
		assert.Equal(t, chains.FallbackConfig.FinalityDepth, config.EvmFinalityDepth())
		assert.Equal(t, chains.FallbackConfig.HeadTrackerHistoryDepth, config.EvmHeadTrackerHistoryDepth())
		assert.Equal(t, chains.FallbackConfig.BalanceMonitorBlockDelay, config.EvmBalanceMonitorBlockDelay())
		assert.Equal(t, chains.FallbackConfig.EthTxResendAfterThreshold, config.EthTxResendAfterThreshold())
		assert.Equal(t, chains.FallbackConfig.BlockHistoryEstimatorBlockDelay, config.BlockHistoryEstimatorBlockDelay())
		assert.Equal(t, chains.FallbackConfig.BlockHistoryEstimatorBlockHistorySize, config.BlockHistoryEstimatorBlockHistorySize())
		assert.Equal(t, chains.FallbackConfig.MinIncomingConfirmations, config.MinIncomingConfirmations())
		assert.Equal(t, chains.FallbackConfig.MinRequiredOutgoingConfirmations, config.MinRequiredOutgoingConfirmations())
	})

	t.Run("with known chain ID returns defaults for that chain", func(t *testing.T) {
		config := newEVMConfigWithChainID("80001")

		assert.Equal(t, chains.PolygonMumbai.Config().GasBumpThreshold, config.EvmGasBumpThreshold())
		assert.Equal(t, chains.PolygonMumbai.Config().GasBumpWei, *config.EvmGasBumpWei())
		assert.Equal(t, chains.PolygonMumbai.Config().GasPriceDefault, *config.EvmGasPriceDefault())
		assert.Equal(t, chains.PolygonMumbai.Config().MaxGasPriceWei, *config.EvmMaxGasPriceWei())
		assert.Equal(t, chains.PolygonMumbai.Config().FinalityDepth, config.EvmFinalityDepth())
		assert.Equal(t, chains.PolygonMumbai.Config().HeadTrackerHistoryDepth, config.EvmHeadTrackerHistoryDepth())
		assert.Equal(t, chains.PolygonMumbai.Config().BalanceMonitorBlockDelay, config.EvmBalanceMonitorBlockDelay())
		assert.Equal(t, chains.PolygonMumbai.Config().EthTxResendAfterThreshold, config.EthTxResendAfterThreshold())
		assert.Equal(t, chains.PolygonMumbai.Config().BlockHistoryEstimatorBlockDelay, config.BlockHistoryEstimatorBlockDelay())
		assert.Equal(t, chains.PolygonMumbai.Config().BlockHistoryEstimatorBlockHistorySize, config.BlockHistoryEstimatorBlockHistorySize())
		assert.Equal(t, chains.PolygonMumbai.Config().MinIncomingConfirmations, config.MinIncomingConfirmations())
		assert.Equal(t, chains.PolygonMumbai.Config().MinRequiredOutgoingConfirmations, config.MinRequiredOutgoingConfirmations())
	})
}

func TestConfig_readFromFile(t *testing.T) {
	v := viper.New()
	v.Set("ROOT", "../../../tools/clroot/")

	config := newGeneralConfigWithViper(v)
	assert.Equal(t, config.RootDir(), "../../../tools/clroot/")
	assert.Equal(t, config.Dev(), true)
	assert.Equal(t, config.TLSPort(), uint16(0))
}

func TestStore_addressParser(t *testing.T) {
	zero := &common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fifteen := &common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15}

	val, err := parseAddress("")
	assert.NoError(t, err)
	assert.Equal(t, nil, val)

	val, err = parseAddress("0x000000000000000000000000000000000000000F")
	assert.NoError(t, err)
	assert.Equal(t, fifteen, val)

	val, err = parseAddress("0X000000000000000000000000000000000000000F")
	assert.NoError(t, err)
	assert.Equal(t, fifteen, val)

	val, err = parseAddress("0")
	assert.NoError(t, err)
	assert.Equal(t, zero, val)

	val, err = parseAddress("15")
	assert.NoError(t, err)
	assert.Equal(t, fifteen, val)

	val, err = parseAddress("0x0")
	assert.Error(t, err)
	assert.Nil(t, val)

	val, err = parseAddress("x")
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestStore_bigIntParser(t *testing.T) {
	val, err := parseBigInt("0")
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(0), val)

	val, err = parseBigInt("15")
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(15), val)

	val, err = parseBigInt("x")
	assert.Error(t, err)
	assert.Nil(t, val)

	val, err = parseBigInt("")
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestStore_levelParser(t *testing.T) {
	val, err := parseLogLevel("ERROR")
	assert.NoError(t, err)
	assert.Equal(t, LogLevel{zapcore.ErrorLevel}, val)

	val, err = parseLogLevel("")
	assert.NoError(t, err)
	assert.Equal(t, LogLevel{zapcore.InfoLevel}, val)

	val, err = parseLogLevel("primus sucks")
	assert.Error(t, err)
	assert.Equal(t, val, LogLevel{})
}

func TestStore_urlParser(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid URL", "http://localhost:3000", false},
		{"invalid URL", ":", true},
		{"empty URL", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i, err := parseURL(test.input)

			if test.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				w, ok := i.(*url.URL)
				require.True(t, ok)
				assert.Equal(t, test.input, w.String())
			}
		})
	}
}

func TestStore_boolParser(t *testing.T) {
	val, err := parseBool("true")
	assert.NoError(t, err)
	assert.Equal(t, true, val)

	val, err = parseBool("false")
	assert.NoError(t, err)
	assert.Equal(t, false, val)

	_, err = parseBool("")
	assert.Error(t, err)
}
