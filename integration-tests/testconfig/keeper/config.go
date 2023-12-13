package keeper

type Config struct {
	Common *Common `toml:"Common"`
}

type Common struct {
	RegistryVersion      *string  `toml:"registry_version"`
	RegistryToTest       *string  `toml:"registry_to_test"`
	NumberOfRegistries   *int     `toml:"number_of_registries"`
	NumberOfNodes        *int     `toml:"number_of_nodes"`
	NumberOfUpkeeps      *int     `toml:"number_of_upkeeps"`
	UpkeepGasLimit       *int64   `toml:"upkeep_gas_limit"`
	CheckGasToBurn       *int64   `toml:"check_gas_to_burn"`
	PerformGasToBurn     *int64   `toml:"perform_gas_to_burn"`
	MaxPerformGas        *int64   `toml:"max_perform_gas"`
	BlockRange           *int64   `toml:"block_range"`
	BlockInterval        *int64   `toml:"block_interval"`
	ChainlinkNodeFunding *float64 `toml:"chainlink_node_funding"`
	ForceSingleTxKey     *bool    `toml:"forces_single_tx_key"`
	DeleteJobsOnEnd      *bool    `toml:"delete_jobs_on_end"`
	RegistryAddress      *string  `toml:"registry_address"`
	RegistrarAddress     *string  `toml:"registrar_address"`
	LinkTokenAddress     *string  `toml:"link_token_address"`
	EthFeedAddress       *string  `toml:"eth_feed_address"`
	GasFeedAddress       *string  `toml:"gas_feed_address"`
	TestInputs           []string `toml:"test_inputs"`
}

func (c *Config) ApplyOverrides(from interface{}) error {
	//TODO implement me
	return nil
}
