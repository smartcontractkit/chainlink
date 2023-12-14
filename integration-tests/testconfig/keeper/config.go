package keeper

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/net"
)

type Config struct {
	Common *Common `toml:"Common"`
}

func (c *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}
	if from.Common != nil {
		if c.Common == nil {
			c.Common = from.Common
		}
		if err := c.Common.ApplyOverrides(from.Common); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Validate() error {
	if c.Common == nil {
		return nil
	}
	return c.Common.Validate()
}

type Common struct {
	RegistryToTest     *string  `toml:"registry_to_test"`
	NumberOfRegistries *int     `toml:"number_of_registries"`
	NumberOfNodes      *int     `toml:"number_of_nodes"`
	NumberOfUpkeeps    *int     `toml:"number_of_upkeeps"`
	UpkeepGasLimit     *int64   `toml:"upkeep_gas_limit"`
	CheckGasToBurn     *int64   `toml:"check_gas_to_burn"`
	PerformGasToBurn   *int64   `toml:"perform_gas_to_burn"`
	MaxPerformGas      *int64   `toml:"max_perform_gas"`
	BlockRange         *int64   `toml:"block_range"`
	BlockInterval      *int64   `toml:"block_interval"`
	ForceSingleTxKey   *bool    `toml:"forces_single_tx_key"`
	DeleteJobsOnEnd    *bool    `toml:"delete_jobs_on_end"`
	RegistryAddress    *string  `toml:"registry_address"`
	RegistrarAddress   *string  `toml:"registrar_address"`
	LinkTokenAddress   *string  `toml:"link_token_address"`
	EthFeedAddress     *string  `toml:"eth_feed_address"`
	GasFeedAddress     *string  `toml:"gas_feed_address"`
	TestInputs         []string `toml:"test_inputs"`
}

func (c *Common) ApplyOverrides(from *Common) error {
	if from == nil {
		return nil
	}
	if from.RegistryToTest != nil {
		c.RegistryToTest = from.RegistryToTest
	}
	if from.NumberOfRegistries != nil {
		c.NumberOfRegistries = from.NumberOfRegistries
	}
	if from.NumberOfNodes != nil {
		c.NumberOfNodes = from.NumberOfNodes
	}
	if from.NumberOfUpkeeps != nil {
		c.NumberOfUpkeeps = from.NumberOfUpkeeps
	}
	if from.UpkeepGasLimit != nil {
		c.UpkeepGasLimit = from.UpkeepGasLimit
	}
	if from.CheckGasToBurn != nil {
		c.CheckGasToBurn = from.CheckGasToBurn
	}
	if from.PerformGasToBurn != nil {
		c.PerformGasToBurn = from.PerformGasToBurn
	}
	if from.MaxPerformGas != nil {
		c.MaxPerformGas = from.MaxPerformGas
	}
	if from.BlockRange != nil {
		c.BlockRange = from.BlockRange
	}
	if from.BlockInterval != nil {
		c.BlockInterval = from.BlockInterval
	}
	if from.ForceSingleTxKey != nil {
		c.ForceSingleTxKey = from.ForceSingleTxKey
	}
	if from.DeleteJobsOnEnd != nil {
		c.DeleteJobsOnEnd = from.DeleteJobsOnEnd
	}
	if from.RegistryAddress != nil {
		c.RegistryAddress = from.RegistryAddress
	}
	if from.RegistrarAddress != nil {
		c.RegistrarAddress = from.RegistrarAddress
	}
	if from.LinkTokenAddress != nil {
		c.LinkTokenAddress = from.LinkTokenAddress
	}
	if from.EthFeedAddress != nil {
		c.EthFeedAddress = from.EthFeedAddress
	}
	if from.GasFeedAddress != nil {
		c.GasFeedAddress = from.GasFeedAddress
	}
	if from.TestInputs != nil {
		c.TestInputs = from.TestInputs
	}
	return nil
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
	if c.RegistryAddress == nil || *c.RegistryAddress == "" {
		return errors.New("registry_address must be set")
	}
	if !net.IsValidURL(*c.RegistryAddress) {
		return errors.New("registry_address must be a valid address")
	}
	if c.RegistrarAddress == nil || *c.RegistrarAddress == "" {
		return errors.New("registrar_address must be set")
	}
	if !net.IsValidURL(*c.RegistrarAddress) {
		return errors.New("registrar_address must be a valid address")
	}
	if c.LinkTokenAddress == nil || *c.LinkTokenAddress == "" {
		return errors.New("link_token_address must be set")
	}
	if !net.IsValidURL(*c.LinkTokenAddress) {
		return errors.New("link_token_address must be a valid address")
	}
	if c.EthFeedAddress == nil || *c.EthFeedAddress == "" {
		return errors.New("eth_feed_address must be set")
	}
	if !net.IsValidURL(*c.EthFeedAddress) {
		return errors.New("eth_feed_address must be a valid address")
	}
	if c.GasFeedAddress == nil || *c.GasFeedAddress == "" {
		return errors.New("gas_feed_address must be set")
	}
	if !net.IsValidURL(*c.GasFeedAddress) {
		return errors.New("gas_feed_address must be a valid address")
	}
	return nil
}
