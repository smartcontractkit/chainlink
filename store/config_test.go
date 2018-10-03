package store

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestStore_ConfigDefaults(t *testing.T) {
	t.Parallel()
	config := NewConfig()
	assert.Equal(t, uint64(0), config.ChainID)
	assert.Equal(t, *big.NewInt(20000000000), config.EthGasPriceDefault)
	assert.Equal(t, "0x514910771AF9Ca656af840dff83E8264EcF986CA", common.HexToAddress(config.LinkContractAddress).String())
	assert.Equal(t, *assets.NewLink(1000000000000000000), config.MinimumContractPayment)
	assert.Equal(t, 15*time.Minute, config.SessionTimeout.Duration)
	assert.Equal(t, "", config.BridgeResponseURL.String())
}

func TestConfig_sessionSecret(t *testing.T) {
	t.Parallel()
	config := NewConfig()
	config.RootDir = path.Join("/tmp/chainlink_test", fmt.Sprintf("%s", "TestConfig_sessionSecret"))
	err := os.MkdirAll(config.RootDir, os.FileMode(0770))
	require.NoError(t, err)

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

	tests := []struct {
		name string
		dev  bool
		want bool
	}{
		{"dev", true, false},
		{"production", false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config.Dev = test.dev
			opts := config.SessionOptions()
			require.Equal(t, test.want, opts.Secure)
		})
	}
}

func TestStore_DurationMarshalJSON(t *testing.T) {
	t.Parallel()

	d := Duration{
		Duration: time.Millisecond,
	}
	b, err := json.Marshal(d)

	assert.NoError(t, err)
	assert.Equal(t, []byte(`"1ms"`), b)
}

func TestStore_DurationUnmarshalJSON(t *testing.T) {
	t.Parallel()

	da := Duration{}
	err := json.Unmarshal([]byte(`"1ms"`), &da)
	assert.NoError(t, err)
	assert.Equal(t, Duration{Duration: time.Millisecond}, da)
}

func TestStore_addressParser(t *testing.T) {
	zero := &common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fifteen := &common.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15}

	val, err := addressParser("")
	assert.NoError(t, err)
	assert.Equal(t, nil, val)

	val, err = addressParser("0x000000000000000000000000000000000000000F")
	assert.NoError(t, err)
	assert.Equal(t, fifteen, val)

	val, err = addressParser("0X000000000000000000000000000000000000000F")
	assert.NoError(t, err)
	assert.Equal(t, fifteen, val)

	val, err = addressParser("0")
	assert.NoError(t, err)
	assert.Equal(t, zero, val)

	val, err = addressParser("15")
	assert.NoError(t, err)
	assert.Equal(t, fifteen, val)

	val, err = addressParser("0x0")
	assert.Error(t, err)

	val, err = addressParser("x")
	assert.Error(t, err)
}

func TestStore_bigIntParser(t *testing.T) {
	val, err := bigIntParser("0")
	assert.NoError(t, err)
	assert.Equal(t, *new(big.Int).SetInt64(0), val)

	val, err = bigIntParser("15")
	assert.NoError(t, err)
	assert.Equal(t, *new(big.Int).SetInt64(15), val)

	val, err = bigIntParser("x")
	assert.Error(t, err)

	val, err = bigIntParser("")
	assert.Error(t, err)
}

func TestStore_levelParser(t *testing.T) {
	val, err := levelParser("ERROR")
	assert.NoError(t, err)
	assert.Equal(t, LogLevel{zapcore.ErrorLevel}, val)

	val, err = levelParser("")
	assert.NoError(t, err)
	assert.Equal(t, LogLevel{zapcore.InfoLevel}, val)

	val, err = levelParser("primus sucks")
	assert.Error(t, err)
}
