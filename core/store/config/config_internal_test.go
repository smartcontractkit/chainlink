package config

import (
	"math/big"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestStore_ConfigDefaults(t *testing.T) {
	config := NewConfig()
	assert.Equal(t, uint64(10), config.BlockBackfillDepth())
	assert.Equal(t, new(url.URL), config.BridgeResponseURL())
	assert.Equal(t, big.NewInt(1), config.ChainID())
	assert.Equal(t, false, config.EthereumDisabled())
	assert.Equal(t, big.NewInt(20000000000), config.EthGasPriceDefault())
	assert.Equal(t, false, config.FeatureExternalInitiators())
	assert.Equal(t, "0x514910771AF9Ca656af840dff83E8264EcF986CA", common.HexToAddress(config.LinkContractAddress()).String())
	assert.Equal(t, assets.NewLink(1000000000000000000), config.MinimumContractPayment())
	assert.Equal(t, 15*time.Minute, config.SessionTimeout().Duration())
}

func TestConfig_sessionSecret(t *testing.T) {
	t.Parallel()
	config := NewConfig()
	config.Set("ROOT", path.Join("/tmp/chainlink_test", "TestConfig_sessionSecret"))
	err := os.MkdirAll(config.RootDir(), os.FileMode(0770))
	require.NoError(t, err)
	defer os.RemoveAll(config.RootDir())

	initial, err := config.SessionSecret()
	require.NoError(t, err)
	require.NotEqual(t, "", initial)
	require.NotEqual(t, "clsession_test_secret", initial)

	second, err := config.SessionSecret()
	require.NoError(t, err)
	require.Equal(t, initial, second)
}

func TestConfig_sessionOptions(t *testing.T) {
	t.Parallel()
	config := NewConfig()

	config.Set("SECURE_COOKIES", false)
	opts := config.SessionOptions()
	require.False(t, opts.Secure)

	config.Set("SECURE_COOKIES", true)
	opts = config.SessionOptions()
	require.True(t, opts.Secure)
}

func TestConfig_EthereumSecondaryURLs(t *testing.T) {
	t.Parallel()
	config := NewConfig()

	localhost, err := url.Parse("http://localhost")
	require.NoError(t, err)
	readme, err := url.Parse("http://readme.net")
	require.NoError(t, err)

	tests := []struct {
		name     string
		oldInput string
		newInput string
		output   []url.URL
	}{
		{"nothing specified", "", "", []url.URL{}},
		{"old option specified", "http://localhost", "", []url.URL{*localhost}},
		{"new option specified", "", "http://localhost", []url.URL{*localhost}},
		{"multipl new options specified", "", "http://localhost ; http://readme.net", []url.URL{*localhost, *readme}},
		{"multipl new options specified comma", "", "http://localhost,http://readme.net", []url.URL{*localhost, *readme}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config.Set("ETH_SECONDARY_URL", test.oldInput)
			config.Set("ETH_SECONDARY_URLS", test.newInput)

			urls := config.EthereumSecondaryURLs()
			assert.Equal(t, test.output, urls)
		})
	}
}

func TestConfig_ChainSpecificConfig(t *testing.T) {
	t.Parallel()

	t.Run("with unknown chain ID returns generic defaults", func(t *testing.T) {
		config := NewConfig()
		config.Set("ETH_CHAIN_ID", "0")

		assert.Equal(t, chains.FallbackConfig.EthGasBumpThreshold, config.EthGasBumpThreshold())
		assert.Equal(t, chains.FallbackConfig.EthGasBumpWei, *config.EthGasBumpWei())
		assert.Equal(t, chains.FallbackConfig.EthGasPriceDefault, *config.EthGasPriceDefault())
		assert.Equal(t, chains.FallbackConfig.EthMaxGasPriceWei, *config.EthMaxGasPriceWei())
		assert.Equal(t, chains.FallbackConfig.EthFinalityDepth, config.EthFinalityDepth())
		assert.Equal(t, chains.FallbackConfig.EthHeadTrackerHistoryDepth, config.EthHeadTrackerHistoryDepth())
		assert.Equal(t, chains.FallbackConfig.EthBalanceMonitorBlockDelay, config.EthBalanceMonitorBlockDelay())
		assert.Equal(t, chains.FallbackConfig.EthTxResendAfterThreshold, config.EthTxResendAfterThreshold())
		assert.Equal(t, chains.FallbackConfig.BlockHistoryEstimatorBlockDelay, config.BlockHistoryEstimatorBlockDelay())
		assert.Equal(t, chains.FallbackConfig.BlockHistoryEstimatorBlockHistorySize, config.BlockHistoryEstimatorBlockHistorySize())
		assert.Equal(t, chains.FallbackConfig.MinIncomingConfirmations, config.MinIncomingConfirmations())
		assert.Equal(t, chains.FallbackConfig.MinRequiredOutgoingConfirmations, config.MinRequiredOutgoingConfirmations())
	})

	t.Run("with known chain ID returns defaults for that chain", func(t *testing.T) {
		config := NewConfig()
		config.Set("ETH_CHAIN_ID", "80001")

		assert.Equal(t, chains.PolygonMumbai.Config().EthGasBumpThreshold, config.EthGasBumpThreshold())
		assert.Equal(t, chains.PolygonMumbai.Config().EthGasBumpWei, *config.EthGasBumpWei())
		assert.Equal(t, chains.PolygonMumbai.Config().EthGasPriceDefault, *config.EthGasPriceDefault())
		assert.Equal(t, chains.PolygonMumbai.Config().EthMaxGasPriceWei, *config.EthMaxGasPriceWei())
		assert.Equal(t, chains.PolygonMumbai.Config().EthFinalityDepth, config.EthFinalityDepth())
		assert.Equal(t, chains.PolygonMumbai.Config().EthHeadTrackerHistoryDepth, config.EthHeadTrackerHistoryDepth())
		assert.Equal(t, chains.PolygonMumbai.Config().EthBalanceMonitorBlockDelay, config.EthBalanceMonitorBlockDelay())
		assert.Equal(t, chains.PolygonMumbai.Config().EthTxResendAfterThreshold, config.EthTxResendAfterThreshold())
		assert.Equal(t, chains.PolygonMumbai.Config().BlockHistoryEstimatorBlockDelay, config.BlockHistoryEstimatorBlockDelay())
		assert.Equal(t, chains.PolygonMumbai.Config().BlockHistoryEstimatorBlockHistorySize, config.BlockHistoryEstimatorBlockHistorySize())
		assert.Equal(t, chains.PolygonMumbai.Config().MinIncomingConfirmations, config.MinIncomingConfirmations())
		assert.Equal(t, chains.PolygonMumbai.Config().MinRequiredOutgoingConfirmations, config.MinRequiredOutgoingConfirmations())
	})

	t.Run("setting env var overrides", func(t *testing.T) {
		config := NewConfig()
		config.Set("ETH_GAS_BUMP_THRESHOLD", "42")
		config.Set("ETH_GAS_BUMP_WEI", "42")
		config.Set("ETH_GAS_PRICE_DEFAULT", "42")
		config.Set("ETH_MAX_GAS_PRICE_WEI", "42")
		config.Set("ETH_FINALITY_DEPTH", "42")
		config.Set("ETH_HEAD_TRACKER_HISTORY_DEPTH", "42")
		config.Set("ETH_BALANCE_MONITOR_BLOCK_DELAY", "42")
		config.Set("ETH_TX_RESEND_AFTER_THRESHOLD", "42s")
		config.Set("GAS_UPDATER_BLOCK_DELAY", "42")
		config.Set("GAS_UPDATER_BLOCK_HISTORY_SIZE", "42")
		config.Set("MIN_INCOMING_CONFIRMATIONS", "42")
		config.Set("MIN_OUTGOING_CONFIRMATIONS", "42")

		assert.Equal(t, 42, int(config.EthGasBumpThreshold()))
		assert.Equal(t, "42", config.EthGasBumpWei().String())
		assert.Equal(t, "42", config.EthGasPriceDefault().String())
		assert.Equal(t, "42", config.EthMaxGasPriceWei().String())
		assert.Equal(t, 42, int(config.EthFinalityDepth()))
		assert.Equal(t, 42, int(config.EthHeadTrackerHistoryDepth()))
		assert.Equal(t, 42, int(config.EthBalanceMonitorBlockDelay()))
		assert.Equal(t, 42*time.Second, config.EthTxResendAfterThreshold())
		assert.Equal(t, 42, int(config.BlockHistoryEstimatorBlockDelay()))
		assert.Equal(t, 42, int(config.BlockHistoryEstimatorBlockHistorySize()))
		assert.Equal(t, 42, int(config.MinIncomingConfirmations()))
		assert.Equal(t, 42, int(config.MinRequiredOutgoingConfirmations()))
	})
}

func TestConfig_readFromFile(t *testing.T) {
	v := viper.New()
	v.Set("ROOT", "../../../tools/clroot/")

	config := newConfigWithViper(v)
	assert.Equal(t, config.RootDir(), "../../../tools/clroot/")
	assert.Equal(t, config.MinRequiredOutgoingConfirmations(), uint64(2))
	assert.Equal(t, config.MinimumContractPayment(), assets.NewLink(1000000000000))
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
