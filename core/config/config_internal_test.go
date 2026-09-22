package config

import (
	"math/big"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/config/parse"
)

func TestStore_bigIntParser(t *testing.T) {
	val, err := parse.BigInt("0")
	require.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(0), val)

	val, err = parse.BigInt("15")
	require.NoError(t, err)
	assert.Equal(t, new(big.Int).SetInt64(15), val)

	val, err = parse.BigInt("x")
	require.Error(t, err)
	assert.Nil(t, val)

	val, err = parse.BigInt("")
	require.Error(t, err)
	assert.Nil(t, val)
}

func TestStore_levelParser(t *testing.T) {
	val, err := parse.LogLevel("ERROR")
	require.NoError(t, err)
	assert.Equal(t, zapcore.ErrorLevel, val)

	val, err = parse.LogLevel("")
	require.NoError(t, err)
	assert.Equal(t, zapcore.InfoLevel, val)

	val, err = parse.LogLevel("primus sucks")
	require.Error(t, err)
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
	require.NoError(t, err)
	assert.Equal(t, true, val)

	val, err = parse.Bool("false")
	require.NoError(t, err)
	assert.Equal(t, false, val)

	_, err = parse.Bool("")
	assert.Error(t, err)
}
