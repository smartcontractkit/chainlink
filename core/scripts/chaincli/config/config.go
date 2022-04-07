package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config represents configuration fields
type Config struct {
	NodeURL         string   `mapstructure:"NODE_URL"`
	ChainID         int64    `mapstructure:"CHAIN_ID"`
	PrivateKey      string   `mapstructure:"PRIVATE_KEY"`
	LinkTokenAddr   string   `mapstructure:"LINK_TOKEN_ADDR"`
	Keepers         []string `mapstructure:"KEEPERS"`
	KeeperURLs      []string `mapstructure:"KEEPER_URLS"`
	KeeperEmails    []string `mapstructure:"KEEPER_EMAILS"`
	KeeperPasswords []string `mapstructure:"KEEPER_PASSWORDS"`
	ApproveAmount   string   `mapstructure:"APPROVE_AMOUNT"`
	GasLimit        uint64   `mapstructure:"GAS_LIMIT"`
	FundNodeAmount  string   `mapstructure:"FUND_CHAINLINK_NODE"`

	// Keeper config
	LinkETHFeedAddr      string `mapstructure:"LINK_ETH_FEED"`
	FastGasFeedAddr      string `mapstructure:"FAST_GAS_FEED"`
	PaymentPremiumPBB    uint32 `mapstructure:"PAYMENT_PREMIUM_PBB"`
	FlatFeeMicroLink     uint32 `mapstructure:"FLAT_FEE_MICRO_LINK"`
	BlockCountPerTurn    int64  `mapstructure:"BLOCK_COUNT_PER_TURN"`
	CheckGasLimit        uint32 `mapstructure:"CHECK_GAS_LIMIT"`
	StalenessSeconds     int64  `mapstructure:"STALENESS_SECONDS"`
	GasCeilingMultiplier uint16 `mapstructure:"GAS_CEILING_MULTIPLIER"`
	FallbackGasPrice     int64  `mapstructure:"FALLBACK_GAS_PRICE"`
	FallbackLinkPrice    int64  `mapstructure:"FALLBACK_LINK_PRICE"`

	// Keepers Config
	RegistryAddress                 string `mapstructure:"KEEPER_REGISTRY_ADDRESS"`
	RegistryConfigUpdate            bool   `mapstructure:"KEEPER_CONFIG_UPDATE"`
	KeepersCount                    int    `mapstructure:"KEEPERS_COUNT"`
	UpkeepTestRange                 int64  `mapstructure:"UPKEEP_TEST_RANGE"`
	UpkeepAverageEligibilityCadence int64  `mapstructure:"UPKEEP_AVERAGE_ELIGIBILITY_CADENCE"`
	UpkeepInterval                  int64  `mapstructure:"UPKEEP_INTERVAL"`
	UpkeepCheckData                 string `mapstructure:"UPKEEP_CHECK_DATA"`
	UpkeepGasLimit                  uint32 `mapstructure:"UPKEEP_GAS_LIMIT"`
	UpkeepCount                     int64  `mapstructure:"UPKEEP_COUNT"`
	AddFundsAmount                  string `mapstructure:"UPKEEP_ADD_FUNDS_AMOUNT"`

	// Feeds config
	FeedBaseAddr  string `mapstructure:"FEED_BASE_ADDR"`
	FeedQuoteAddr string `mapstructure:"FEED_QUOTE_ADDR"`
	FeedDecimals  uint8  `mapstructure:"FEED_DECIMALS"`
}

// New is the constructor of Config
func New() *Config {
	var cfg Config
	configFile := viper.GetString("config")
	if configFile != "" {
		log.Println("Using config file", configFile)
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		log.Println("Using config file .env")
		viper.SetConfigFile(".env")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config: ", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("failed to unmarshal config: ", err)
	}

	return &cfg
}

func init() {
	// Represented in WEI, which is 1000 Ether
	viper.SetDefault("APPROVE_AMOUNT", "100000000000000000000000")
	viper.SetDefault("GAS_LIMIT", 8000000)
	viper.SetDefault("PAYMENT_PREMIUM_PBB", 200000000)
	viper.SetDefault("FLAT_FEE_MICRO_LINK", 0)
	viper.SetDefault("BLOCK_COUNT_PER_TURN", 1)
	viper.SetDefault("CHECK_GAS_LIMIT", 650000000)
	viper.SetDefault("STALENESS_SECONDS", 90000)
	viper.SetDefault("GAS_CEILING_MULTIPLIER", 3)
	viper.SetDefault("FALLBACK_GAS_PRICE", 200000000000)
	viper.SetDefault("FALLBACK_LINK_PRICE", 20000000000000000)
	// Represented in WEI, which is 100 Ether
	viper.SetDefault("UPKEEP_ADD_FUNDS_AMOUNT", "100000000000000000000")
	viper.SetDefault("UPKEEP_TEST_RANGE", 1)
	viper.SetDefault("UPKEEP_INTERVAL", 10)
	viper.SetDefault("UPKEEP_CHECK_DATA", "0x00")
	viper.SetDefault("UPKEEP_GAS_LIMIT", 500000)
	viper.SetDefault("UPKEEP_COUNT", 5)
	viper.SetDefault("KEEPERS_COUNT", 2)

	viper.SetDefault("FEED_DECIMALS", 18)
	viper.SetDefault("MUST_TAKE_TURNS", true)
}
