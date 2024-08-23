package keeper

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

type Config struct {
	Common     *Common           `toml:"Common"`
	Resiliency *ResiliencyConfig `toml:"Resiliency"`
}

func (c *Config) Validate() error {
	if c.Common == nil {
		return nil
	}
	if err := c.Common.Validate(); err != nil {
		return err
	}
	if c.Resiliency == nil {
		return nil
	}
	return c.Resiliency.Validate()
}

type Common struct {
	RegistryToTest     *string `toml:"registry_to_test"`
	NumberOfRegistries *int    `toml:"number_of_registries"`
	NumberOfNodes      *int    `toml:"number_of_nodes"`
	NumberOfUpkeeps    *int    `toml:"number_of_upkeeps"`
	UpkeepGasLimit     *int64  `toml:"upkeep_gas_limit"`
	CheckGasToBurn     *int64  `toml:"check_gas_to_burn"`
	PerformGasToBurn   *int64  `toml:"perform_gas_to_burn"`
	MaxPerformGas      *int64  `toml:"max_perform_gas"`
	BlockRange         *int64  `toml:"block_range"`
	BlockInterval      *int64  `toml:"block_interval"`
	ForceSingleTxKey   *bool   `toml:"forces_single_tx_key"`
	DeleteJobsOnEnd    *bool   `toml:"delete_jobs_on_end"`
	RegistryAddress    *string `toml:"registry_address"`
	RegistrarAddress   *string `toml:"registrar_address"`
	LinkTokenAddress   *string `toml:"link_token_address"`
	EthFeedAddress     *string `toml:"eth_feed_address"`
	GasFeedAddress     *string `toml:"gas_feed_address"`
}

func (c *Common) Validate() error {
	if c.RegistryToTest == nil || *c.RegistryToTest == "" {
		return errors.New("registry_to_test must be set")
	}
	if c.NumberOfRegistries == nil || *c.NumberOfRegistries <= 0 {
		return errors.New("number_of_registries must be a positive integer")
	}
	if c.NumberOfNodes == nil || *c.NumberOfNodes <= 0 {
		return errors.New("number_of_nodes must be a positive integer")
	}
	if c.NumberOfUpkeeps == nil || *c.NumberOfUpkeeps <= 0 {
		return errors.New("number_of_upkeeps must be a positive integer")
	}
	if c.UpkeepGasLimit == nil || *c.UpkeepGasLimit <= 0 {
		return errors.New("upkeep_gas_limit must be a positive integer")
	}
	if c.CheckGasToBurn == nil || *c.CheckGasToBurn <= 0 {
		return errors.New("check_gas_to_burn must be a positive integer")
	}
	if c.PerformGasToBurn == nil || *c.PerformGasToBurn <= 0 {
		return errors.New("perform_gas_to_burn must be a positive integer")
	}
	if c.MaxPerformGas == nil || *c.MaxPerformGas <= 0 {
		return errors.New("max_perform_gas must be a positive integer")
	}
	if c.BlockRange == nil || *c.BlockRange <= 0 {
		return errors.New("block_range must be a positive integer")
	}
	if c.BlockInterval == nil || *c.BlockInterval <= 0 {
		return errors.New("block_interval must be a positive integer")
	}
	if c.RegistryAddress == nil {
		c.RegistryAddress = new(string)
	}
	if c.RegistrarAddress == nil {
		c.RegistrarAddress = new(string)
	}
	if c.LinkTokenAddress == nil {
		c.LinkTokenAddress = new(string)
	}
	if c.EthFeedAddress == nil {
		c.EthFeedAddress = new(string)
	}
	if c.GasFeedAddress == nil {
		c.GasFeedAddress = new(string)
	}
	return nil
}

type ResiliencyConfig struct {
	ContractCallLimit    *uint                   `toml:"contract_call_limit"`
	ContractCallInterval *blockchain.StrDuration `toml:"contract_call_interval"`
}

func (c *ResiliencyConfig) Validate() error {
	if c.ContractCallLimit == nil {
		return errors.New("contract_call_limit must be set")
	}
	if c.ContractCallInterval == nil {
		return errors.New("contract_call_interval must be set")
	}

	return nil
}
