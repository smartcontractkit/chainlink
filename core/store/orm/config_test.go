package orm

import (
	"math/big"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"

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

func TestConfig_readFromFile(t *testing.T) {
	v := viper.New()
	v.Set("ROOT", "../../../tools/clroot/")

	config := newConfigWithViper(v)
	assert.Equal(t, config.RootDir(), "../../../tools/clroot/")
	assert.Equal(t, config.MinOutgoingConfirmations(), uint64(2))
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
