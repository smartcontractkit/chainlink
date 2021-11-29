package config

import (
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestGeneralConfig_Defaults(t *testing.T) {
	config := NewGeneralConfig()
	assert.Equal(t, uint64(10), config.BlockBackfillDepth())
	assert.Equal(t, new(url.URL), config.BridgeResponseURL())
	assert.Nil(t, config.DefaultChainID())
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

func TestConfig_readFromFile(t *testing.T) {
	v := viper.New()
	v.Set("ROOT", "../../tools/clroot/")

	config := newGeneralConfigWithViper(v)
	assert.Equal(t, config.RootDir(), "../../tools/clroot/")
	assert.Equal(t, config.Dev(), true)
	assert.Equal(t, config.TLSPort(), uint16(0))
}

func TestStore_bigIntParser(t *testing.T) {
	val, err := ParseBigInt("0")
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(0), val)

	val, err = ParseBigInt("15")
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(15), val)

	val, err = ParseBigInt("x")
	assert.Error(t, err)
	assert.Nil(t, val)

	val, err = ParseBigInt("")
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestStore_levelParser(t *testing.T) {
	val, err := ParseLogLevel("ERROR")
	assert.NoError(t, err)
	assert.Equal(t, LogLevel{zapcore.ErrorLevel}, val)

	val, err = ParseLogLevel("")
	assert.NoError(t, err)
	assert.Equal(t, LogLevel{zapcore.InfoLevel}, val)

	val, err = ParseLogLevel("primus sucks")
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
			i, err := ParseURL(test.input)

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
	val, err := ParseBool("true")
	assert.NoError(t, err)
	assert.Equal(t, true, val)

	val, err = ParseBool("false")
	assert.NoError(t, err)
	assert.Equal(t, false, val)

	_, err = ParseBool("")
	assert.Error(t, err)
}
