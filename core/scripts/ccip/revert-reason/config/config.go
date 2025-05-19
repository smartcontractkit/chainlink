package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config represents configuration fields
type Config struct {
	NodeURL     string `mapstructure:"NODE_URL"`
	FromAddress string `mapstructure:"FROM_ADDRESS"`
}

// New creates a new config
func New() *Config {
	var cfg Config
	configFile := viper.GetString("config")
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile(".env")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config: ", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("failed to unmarshal config: ", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatal("failed to validate config: ", err)
	}

	return &cfg
}

// Validate validates the given config
func (c *Config) Validate() error {
	return nil
}
