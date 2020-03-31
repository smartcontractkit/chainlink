package service

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config contains the startup parameters to configure the monitor
type Config struct {
	// ResponseTimeout is the duration to wait before the response is treated as timed-out to alert on
	ResponseTimeout time.Duration `mapstructure:"response-timeout"`
	// NetworkID is the Ethereum Chain ID for the contracts you want to listen to
	NetworkID int `mapstructure:"eth-chain-id"`
	// EthereumURL is the websocket endpoint the monitor uses to watch the aggregator contracts
	EthereumURL string `mapstructure:"eth-url"`
	// FeedsUIURL is the base URL for the FeedsTracker UI
	FeedsUIURL string `mapstructure:"feeds-ui-url"`
	// DatabaseHost of the postgres server where the ingester saves results
	DatabaseHost string `mapstructure:"db-host"`
	// DatabaseName of the postgres server where the ingester saves results
	DatabaseName string `mapstructure:"db-name"`
	// DatabasePort of the postgres server where the ingester saves results
	DatabasePort int `mapstructure:"db-port"`
	// DatabaseUsername of the postgres server where the ingester saves results
	DatabaseUsername string `mapstructure:"db-username"`
	// DatabasePassword of the postgres server where the ingester saves results
	DatabasePassword string `mapstructure:"db-password"`
}

// NewConfig will return an instantiated config based on the passed in defaults
// and the environment variables as defined in the config struct
func NewConfig(defaults map[string]interface{}) *Config {
	v := viper.New()
	c := Config{}
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	_ = v.ReadInConfig()
	_ = v.Unmarshal(&c)
	return &c
}

// DefaultConfig returns an instantiated config with the application defaults
func DefaultConfig() *Config {
	return NewConfig(map[string]interface{}{
		"response-timeout": time.Minute * 5,
		"eth-chain-id":     1,
		"eth-url":          "ws://localhost:8545",
		"feeds-ui-url":     "https://feeds.chain.link",
		"db-host":          "localhost",
		"db-name":          "explorer",
		"db-port":          "5432",
		"db-username":      "postgres",
		"db-password":      "postgres",
	})
}

// DefaultConfig returns an instantiated config with the application defaults for testing
func TestConfig() *Config {
	cfg := DefaultConfig()
	return cfg
}
