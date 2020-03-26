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
	// DatabaseURL is the url of the postgres server where the ingester saves results
	DatabaseURL string `mapstructure:"db-url"`
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
		"db-url":           "postgres://localhost:5432/explorer?sslmode=disable",
	})
}

// DefaultConfig returns an instantiated config with the application defaults for testing
func TestConfig() *Config {
	cfg := DefaultConfig()
	return cfg
}
