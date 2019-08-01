package orm

import (
	"math/big"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_ConfigDefaults(t *testing.T) {
	t.Parallel()
	config := NewConfig(NewBootstrapConfigStore())
	assert.Equal(t, uint64(0), config.ChainID())
	assert.Equal(t, big.NewInt(20000000000), config.EthGasPriceDefault())
	assert.Equal(t, "0x514910771AF9Ca656af840dff83E8264EcF986CA", common.HexToAddress(config.LinkContractAddress()).String())
	assert.Equal(t, assets.NewLink(1000000000000000000), config.MinimumContractPayment())
	assert.Equal(t, 15*time.Minute, config.SessionTimeout())
	assert.Equal(t, new(url.URL), config.BridgeResponseURL())
}

func TestConfig_sessionSecret(t *testing.T) {
	t.Parallel()
	store := NewBootstrapConfigStore()
	config := NewConfig(store)

	store.Set("ROOT", path.Join("/tmp/chainlink_test", "TestConfig_sessionSecret"))
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
	store := NewBootstrapConfigStore()
	config := NewConfig(store)

	store.Set("SECURE_COOKIES", false)
	opts := config.SessionOptions()
	require.False(t, opts.Secure)

	store.Set("SECURE_COOKIES", true)
	opts = config.SessionOptions()
	require.True(t, opts.Secure)
}

func TestConfig_readFromFile(t *testing.T) {
	v := viper.New()
	v.Set("ROOT", "../../../tools/clroot/")

	store := newConfigWithViper(v)
	config := NewConfig(store)

	assert.Equal(t, config.RootDir(), "../../../tools/clroot/")
	assert.Equal(t, config.MinOutgoingConfirmations(), uint64(2))
	assert.Equal(t, config.MinimumContractPayment(), assets.NewLink(1000000000000))
	assert.Equal(t, config.Dev(), true)
	assert.Equal(t, config.TLSPort(), uint16(0))
}

func TestConfig_DatabaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, uri, expect string
	}{
		{"default", "", "/root/db.sqlite3"},
		{"root", "/root/db.sqlite3", "/root/db.sqlite3"},
		{"windows root", `C:\root\db.sqlite3`, `C:\root\db.sqlite3`},
		{"garbage", "89324*$*#@(=", "89324*$*#@(="},
		{"relative path", "store/db/here", "store/db/here"},
		{"file uri", "file://host/path", "file://host/path"},
		{"postgres uri", "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full", "postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"},
		{"postgres string", "user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full", "user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := NewBootstrapConfigStore()
			store.Set("ROOT", "/root")
			store.Set("DATABASE_URL", test.uri)

			cfg := NewConfig(store)
			assert.Equal(t, test.expect, cfg.DatabaseURL())
		})
	}
}
