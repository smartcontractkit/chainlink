package cltest

import (
	"encoding"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/orm"
	"go.uber.org/zap/zapcore"
)

// Config is a configuration store used for testing
type Config struct {
	t testing.TB
	orm.Configger
	orm.ConfigStore
}

// NewConfig returns a new Config
func NewConfig(t testing.TB) *Config {
	store := orm.NewBootstrapConfigStore()
	config := orm.NewConfig(store)

	c := &Config{
		Configger:   config,
		ConfigStore: store,
	}

	c.Set("BRIDGE_RESPONSE_URL", "http://localhost:6688")
	c.Set("ETH_CHAIN_ID", 3)
	c.Set("CHAINLINK_DEV", true)
	c.Set("ETH_GAS_BUMP_THRESHOLD", 3)
	c.Set("LOG_LEVEL", orm.LogLevel{Level: zapcore.DebugLevel})
	c.Set("MINIMUM_SERVICE_DURATION", "24h")
	c.Set("MIN_INCOMING_CONFIRMATIONS", 1)
	c.Set("MIN_OUTGOING_CONFIRMATIONS", 6)
	c.Set("MINIMUM_CONTRACT_PAYMENT", minimumContractPayment.Text(10))
	//c.Set("ROOT", rootdir)
	c.Set("SESSION_TIMEOUT", "2m")

	return c
}

// Shutdown cleans up any resources allocated by the Config
func (c Config) Shutdown() {
}

// Set saves any type to the config store, only to be used in tests
func (c Config) Set(name string, value interface{}) {
	switch v := value.(type) {
	case encoding.TextMarshaler:
		c.SetMarshaler(name, v)
	case fmt.Stringer:
		c.SetStringer(name, v)
	case string:
		c.SetString(name, v)
	}
}

// SessionSecret returns a static session secret
func (c Config) SessionSecret() ([]byte, error) {
	return []byte("clsession_test_secret"), nil
}
