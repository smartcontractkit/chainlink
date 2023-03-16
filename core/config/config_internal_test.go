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

	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func TestGeneralConfig_Defaults(t *testing.T) {
	config := NewGeneralConfig(logger.TestLogger(t))
	assert.Equal(t, uint64(10), config.BlockBackfillDepth())
	assert.Equal(t, (*url.URL)(nil), config.BridgeResponseURL())
	assert.Nil(t, config.DefaultChainID())
	assert.True(t, config.EVMRPCEnabled())
	assert.True(t, config.EVMEnabled())
	assert.False(t, config.SolanaEnabled())
	assert.False(t, config.StarkNetEnabled())
	assert.Equal(t, false, config.FeatureExternalInitiators())
	assert.Equal(t, 15*time.Minute, config.SessionTimeout().Duration())
}

func TestGeneralConfig_GlobalOCRDatabaseTimeout(t *testing.T) {
	t.Setenv(envvar.Name("OCRDatabaseTimeout"), "3s")
	config := NewGeneralConfig(logger.TestLogger(t))

	timeout, ok := config.GlobalOCRDatabaseTimeout()
	require.True(t, ok)
	require.Equal(t, 3*time.Second, timeout)
}

func TestGeneralConfig_GlobalOCRObservationGracePeriod(t *testing.T) {
	t.Setenv(envvar.Name("OCRObservationGracePeriod"), "3s")
	config := NewGeneralConfig(logger.TestLogger(t))

	timeout, ok := config.GlobalOCRObservationGracePeriod()
	require.True(t, ok)
	require.Equal(t, 3*time.Second, timeout)
}

func TestGeneralConfig_GlobalOCRContractTransmitterTransmitTimeout(t *testing.T) {
	t.Setenv(envvar.Name("OCRContractTransmitterTransmitTimeout"), "3s")
	config := NewGeneralConfig(logger.TestLogger(t))

	timeout, ok := config.GlobalOCRContractTransmitterTransmitTimeout()
	require.True(t, ok)
	require.Equal(t, 3*time.Second, timeout)
}

func TestConfig_readFromFile(t *testing.T) {
	v := viper.New()
	v.Set("ROOT", "../../tools/clroot/")

	config := newGeneralConfigWithViper(v, logger.TestLogger(t))
	assert.Equal(t, config.RootDir(), "../../tools/clroot/")
	assert.Equal(t, config.Dev(), true)
	assert.Equal(t, config.TLSPort(), uint16(0))
}

func TestStore_bigIntParser(t *testing.T) {
	val, err := parse.BigInt("0")
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(0), val)

	val, err = parse.BigInt("15")
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(15), val)

	val, err = parse.BigInt("x")
	assert.Error(t, err)
	assert.Nil(t, val)

	val, err = parse.BigInt("")
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestStore_levelParser(t *testing.T) {
	val, err := parse.LogLevel("ERROR")
	assert.NoError(t, err)
	assert.Equal(t, zapcore.ErrorLevel, val)

	val, err = parse.LogLevel("")
	assert.NoError(t, err)
	assert.Equal(t, zapcore.InfoLevel, val)

	val, err = parse.LogLevel("primus sucks")
	assert.Error(t, err)
	assert.Equal(t, val, zapcore.Level(0))
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
			i, err := parse.URL(test.input)

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
	val, err := parse.Bool("true")
	assert.NoError(t, err)
	assert.Equal(t, true, val)

	val, err = parse.Bool("false")
	assert.NoError(t, err)
	assert.Equal(t, false, val)

	_, err = parse.Bool("")
	assert.Error(t, err)
}
