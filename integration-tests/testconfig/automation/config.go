package automation

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
)

type Config struct {
	General          *General          `toml:"General"`
	Load             []Load            `toml:"Load"`
	DataStreams      *DataStreams      `toml:"DataStreams"`
	AutomationConfig *AutomationConfig `toml:"AutomationConfig"`
	Resiliency       *ResiliencyConfig `toml:"Resiliency"`
	Benchmark        *Benchmark        `toml:"Benchmark"`
	Contracts        *Contracts        `toml:"Contracts"`
}

func (c *Config) Validate() error {
	if c.General != nil {
		if err := c.General.Validate(); err != nil {
			return err
		}
	}
	if len(c.Load) > 0 {
		for _, load := range c.Load {
			if err := load.Validate(); err != nil {
				return err
			}
		}
	}
	if c.DataStreams != nil {
		if err := c.DataStreams.Validate(); err != nil {
			return err
		}
	}

	if c.AutomationConfig != nil {
		if err := c.AutomationConfig.Validate(); err != nil {
			return err
		}
	}
	if c.Resiliency != nil {
		if err := c.Resiliency.Validate(); err != nil {
			return err
		}
	}
	if c.Benchmark != nil {
		if err := c.Benchmark.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Benchmark struct {
	RegistryToTest     *string `toml:"registry_to_test"`
	NumberOfRegistries *int    `toml:"number_of_registries"`
	NumberOfUpkeeps    *int    `toml:"number_of_upkeeps"`
	UpkeepGasLimit     *int64  `toml:"upkeep_gas_limit"`
	CheckGasToBurn     *int64  `toml:"check_gas_to_burn"`
	PerformGasToBurn   *int64  `toml:"perform_gas_to_burn"`
	BlockRange         *int64  `toml:"block_range"`
	BlockInterval      *int64  `toml:"block_interval"`
	ForceSingleTxKey   *bool   `toml:"forces_single_tx_key"`
	DeleteJobsOnEnd    *bool   `toml:"delete_jobs_on_end"`
}

func (c *Benchmark) Validate() error {
	if c.RegistryToTest == nil || *c.RegistryToTest == "" {
		return errors.New("registry_to_test must be set")
	}
	if c.NumberOfRegistries == nil || *c.NumberOfRegistries <= 0 {
		return errors.New("number_of_registries must be a positive integer")
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
	if c.BlockRange == nil || *c.BlockRange <= 0 {
		return errors.New("block_range must be a positive integer")
	}
	if c.BlockInterval == nil || *c.BlockInterval <= 0 {
		return errors.New("block_interval must be a positive integer")
	}
	return nil
}

// General is a common configuration for all automation performance tests
type General struct {
	NumberOfNodes         *int    `toml:"number_of_nodes"`
	Duration              *int    `toml:"duration"`
	BlockTime             *int    `toml:"block_time"`
	SpecType              *string `toml:"spec_type"`
	ChainlinkNodeLogLevel *string `toml:"chainlink_node_log_level"`
	UsePrometheus         *bool   `toml:"use_prometheus"`
	RemoveNamespace       *bool   `toml:"remove_namespace"`
}

func (c *General) Validate() error {
	if c.NumberOfNodes == nil || *c.NumberOfNodes < 1 {
		return errors.New("number_of_nodes must be set to a positive integer")
	}
	if c.Duration == nil || *c.Duration < 1 {
		return errors.New("duration must be set to a positive integer")
	}
	if c.BlockTime == nil || *c.BlockTime < 1 {
		return errors.New("block_time must be set to a positive integer")
	}
	if c.SpecType == nil {
		return errors.New("spec_type must be set")
	}
	if c.ChainlinkNodeLogLevel == nil {
		return errors.New("chainlink_node_log_level must be set")
	}
	if c.UsePrometheus == nil {
		return errors.New("use_prometheus must be set")
	}
	if c.RemoveNamespace == nil {
		return errors.New("remove_namespace must be set")
	}

	return nil
}

type Load struct {
	NumberOfUpkeeps               *int     `toml:"number_of_upkeeps"`
	NumberOfEvents                *int     `toml:"number_of_events"`
	NumberOfSpamMatchingEvents    *int     `toml:"number_of_spam_matching_events"`
	NumberOfSpamNonMatchingEvents *int     `toml:"number_of_spam_non_matching_events"`
	CheckBurnAmount               *big.Int `toml:"check_burn_amount"`
	PerformBurnAmount             *big.Int `toml:"perform_burn_amount"`
	SharedTrigger                 *bool    `toml:"shared_trigger"`
	UpkeepGasLimit                *uint32  `toml:"upkeep_gas_limit"`
	IsStreamsLookup               *bool    `toml:"is_streams_lookup"`
	Feeds                         []string `toml:"feeds"`
}

func (c *Load) Validate() error {
	if c.NumberOfUpkeeps == nil || *c.NumberOfUpkeeps < 1 {
		return errors.New("number_of_upkeeps must be set to a positive integer")
	}
	if c.NumberOfEvents == nil || *c.NumberOfEvents < 0 {
		return errors.New("number_of_events must be set to a non-negative integer")
	}
	if c.NumberOfSpamMatchingEvents == nil || *c.NumberOfSpamMatchingEvents < 0 {
		return errors.New("number_of_spam_matching_events must be set to a non-negative integer")
	}
	if c.NumberOfSpamNonMatchingEvents == nil || *c.NumberOfSpamNonMatchingEvents < 0 {
		return errors.New("number_of_spam_non_matching_events must be set to a non-negative integer")
	}
	if c.CheckBurnAmount == nil || c.CheckBurnAmount.Cmp(big.NewInt(0)) < 0 {
		return errors.New("check_burn_amount must be set to a non-negative integer")
	}
	if c.PerformBurnAmount == nil || c.PerformBurnAmount.Cmp(big.NewInt(0)) < 0 {
		return errors.New("perform_burn_amount must be set to a non-negative integer")
	}
	if c.SharedTrigger == nil {
		return errors.New("shared_trigger must be set")
	}
	if c.UpkeepGasLimit == nil || *c.UpkeepGasLimit < 1 {
		return errors.New("upkeep_gas_limit must be set to a positive integer")
	}
	if c.IsStreamsLookup == nil {
		return errors.New("is_streams_lookup must be set")
	}
	if *c.IsStreamsLookup {
		if len(c.Feeds) == 0 {
			return errors.New("feeds must be set")
		}
	}

	return nil
}

type DataStreams struct {
	Enabled       *bool   `toml:"enabled"`
	URL           *string `toml:"-"`
	Username      *string `toml:"-"`
	Password      *string `toml:"-"`
	DefaultFeedID *string `toml:"default_feed_id"`
}

func (c *DataStreams) Validate() error {
	if c.Enabled != nil && *c.Enabled {
		if c.URL == nil {
			return errors.New("data_streams_url must be set")
		}
		if c.Username == nil {
			return errors.New("data_streams_username must be set")
		}
		if c.Password == nil {
			return errors.New("data_streams_password must be set")
		}
		if c.DefaultFeedID == nil {
			return errors.New("data_streams_feed_id must be set")
		}
	} else {
		c.Enabled = new(bool)
		*c.Enabled = false
	}
	return nil
}

type AutomationConfig struct {
	PluginConfig     *PluginConfig     `toml:"PluginConfig"`
	PublicConfig     *PublicConfig     `toml:"PublicConfig"`
	RegistrySettings *RegistrySettings `toml:"RegistrySettings"`
}

func (c *AutomationConfig) Validate() error {
	if err := c.PluginConfig.Validate(); err != nil {
		return err
	}
	if err := c.PublicConfig.Validate(); err != nil {
		return err
	}
	return c.RegistrySettings.Validate()
}

type PluginConfig struct {
	PerformLockoutWindow *int64             `toml:"perform_lockout_window"`
	TargetProbability    *string            `toml:"target_probability"`
	TargetInRounds       *int               `toml:"target_in_rounds"`
	MinConfirmations     *int               `toml:"min_confirmations"`
	GasLimitPerReport    *uint32            `toml:"gas_limit_per_report"`
	GasOverheadPerUpkeep *uint32            `toml:"gas_overhead_per_upkeep"`
	MaxUpkeepBatchSize   *int               `toml:"max_upkeep_batch_size"`
	LogProviderConfig    *LogProviderConfig `toml:"LogProviderConfig"`
}

type LogProviderConfig struct {
	BlockRate *uint32 `toml:"block_rate"`
	LogLimit  *uint32 `toml:"log_limit"`
}

func (c *PluginConfig) Validate() error {
	if err := c.LogProviderConfig.Validate(); err != nil {
		return err
	}
	if c.PerformLockoutWindow == nil || *c.PerformLockoutWindow < 0 {
		return errors.New("perform_lockout_window must be set to a non-negative integer")
	}
	if c.TargetProbability == nil || *c.TargetProbability == "" {
		return errors.New("target_probability must be set")
	}
	if c.TargetInRounds == nil || *c.TargetInRounds < 1 {
		return errors.New("target_in_rounds must be set to a positive integer")
	}
	if c.MinConfirmations == nil || *c.MinConfirmations < 0 {
		return errors.New("min_confirmations must be set to a non-negative integer")
	}
	if c.GasLimitPerReport == nil || *c.GasLimitPerReport < 1 {
		return errors.New("gas_limit_per_report must be set to a positive integer")
	}
	if c.GasOverheadPerUpkeep == nil || *c.GasOverheadPerUpkeep < 1 {
		return errors.New("gas_overhead_per_upkeep must be set to a positive integer")
	}
	if c.MaxUpkeepBatchSize == nil || *c.MaxUpkeepBatchSize < 1 {
		return errors.New("max_upkeep_batch_size must be set to a positive integer")
	}
	return nil

}

func (c *LogProviderConfig) Validate() error {
	if c.BlockRate == nil || *c.BlockRate < 1 {
		return errors.New("block_rate must be set to a positive integer")
	}
	if c.LogLimit == nil || *c.LogLimit < 1 {
		return errors.New("log_limit must be set to a positive integer")
	}
	return nil

}

type PublicConfig struct {
	DeltaProgress                           *time.Duration `toml:"delta_progress"`
	DeltaResend                             *time.Duration `toml:"delta_resend"`
	DeltaInitial                            *time.Duration `toml:"delta_initial"`
	DeltaRound                              *time.Duration `toml:"delta_round"`
	DeltaGrace                              *time.Duration `toml:"delta_grace"`
	DeltaCertifiedCommitRequest             *time.Duration `toml:"delta_certified_commit_request"`
	DeltaStage                              *time.Duration `toml:"delta_stage"`
	RMax                                    *uint64        `toml:"r_max"`
	F                                       *int           `toml:"f"`
	MaxDurationQuery                        *time.Duration `toml:"max_duration_query"`
	MaxDurationObservation                  *time.Duration `toml:"max_duration_observation"`
	MaxDurationShouldAcceptAttestedReport   *time.Duration `toml:"max_duration_should_accept_attested_report"`
	MaxDurationShouldTransmitAcceptedReport *time.Duration `toml:"max_duration_should_transmit_accepted_report"`
}

func (c *PublicConfig) Validate() error {
	if c.DeltaProgress == nil || *c.DeltaProgress < 0 {
		return errors.New("delta_progress must be set to a non-negative duration")
	}
	if c.DeltaResend == nil || *c.DeltaResend < 0 {
		return errors.New("delta_resend must be set to a non-negative duration")
	}
	if c.DeltaInitial == nil || *c.DeltaInitial < 0 {
		return errors.New("delta_initial must be set to a non-negative duration")
	}
	if c.DeltaRound == nil || *c.DeltaRound < 0 {
		return errors.New("delta_round must be set to a non-negative duration")
	}
	if c.DeltaGrace == nil || *c.DeltaGrace < 0 {
		return errors.New("delta_grace must be set to a non-negative duration")
	}
	if c.DeltaCertifiedCommitRequest == nil || *c.DeltaCertifiedCommitRequest < 0 {
		return errors.New("delta_certified_commit_request must be set to a non-negative duration")
	}
	if c.DeltaStage == nil || *c.DeltaStage < 0 {
		return errors.New("delta_stage must be set to a non-negative duration")
	}
	if c.RMax == nil || *c.RMax < 1 {
		return errors.New("r_max must be set to a positive integer")
	}
	if c.F == nil || *c.F < 1 {
		return errors.New("f must be set to a positive integer")
	}
	if c.MaxDurationQuery == nil || *c.MaxDurationQuery < 0 {
		return errors.New("max_duration_query must be set to a non-negative duration")
	}
	if c.MaxDurationObservation == nil || *c.MaxDurationObservation < 0 {
		return errors.New("max_duration_observation must be set to a non-negative duration")
	}
	if c.MaxDurationShouldAcceptAttestedReport == nil || *c.MaxDurationShouldAcceptAttestedReport < 0 {
		return errors.New("max_duration_should_accept_attested_report must be set to a non-negative duration")
	}
	if c.MaxDurationShouldTransmitAcceptedReport == nil || *c.MaxDurationShouldTransmitAcceptedReport < 0 {
		return errors.New("max_duration_should_transmit_accepted_report must be set to a non-negative duration")
	}
	return nil

}

type RegistrySettings struct {
	PaymentPremiumPPB    *uint32  `toml:"payment_premium_ppb"`
	FlatFeeMicroLINK     *uint32  `toml:"flat_fee_micro_link"`
	CheckGasLimit        *uint32  `toml:"check_gas_limit"`
	StalenessSeconds     *big.Int `toml:"staleness_seconds"`
	GasCeilingMultiplier *uint16  `toml:"gas_ceiling_multiplier"`
	MaxPerformGas        *uint32  `toml:"max_perform_gas"`
	MinUpkeepSpend       *big.Int `toml:"min_upkeep_spend"`
	FallbackGasPrice     *big.Int `toml:"fallback_gas_price"`
	FallbackLinkPrice    *big.Int `toml:"fallback_link_price"`
	FallbackNativePrice  *big.Int `toml:"fallback_native_price"`
	MaxCheckDataSize     *uint32  `toml:"max_check_data_size"`
	MaxPerformDataSize   *uint32  `toml:"max_perform_data_size"`
	MaxRevertDataSize    *uint32  `toml:"max_revert_data_size"`
}

func (c *RegistrySettings) Validate() error {
	if c.PaymentPremiumPPB == nil {
		return errors.New("payment_premium_ppb must be set to a non-negative integer")
	}
	if c.FlatFeeMicroLINK == nil {
		return errors.New("flat_fee_micro_link must be set to a non-negative integer")
	}
	if c.CheckGasLimit == nil || *c.CheckGasLimit < 1 {
		return errors.New("check_gas_limit must be set to a positive integer")
	}
	if c.StalenessSeconds == nil || c.StalenessSeconds.Cmp(big.NewInt(0)) < 0 {
		return errors.New("staleness_seconds must be set to a non-negative integer")
	}
	if c.GasCeilingMultiplier == nil {
		return errors.New("gas_ceiling_multiplier must be set to a non-negative integer")
	}
	if c.MaxPerformGas == nil || *c.MaxPerformGas < 1 {
		return errors.New("max_perform_gas must be set to a positive integer")
	}
	if c.MinUpkeepSpend == nil || c.MinUpkeepSpend.Cmp(big.NewInt(0)) < 0 {
		return errors.New("min_upkeep_spend must be set to a non-negative integer")
	}
	if c.FallbackGasPrice == nil || c.FallbackGasPrice.Cmp(big.NewInt(0)) < 0 {
		return errors.New("fallback_gas_price must be set to a non-negative integer")
	}
	if c.FallbackLinkPrice == nil || c.FallbackLinkPrice.Cmp(big.NewInt(0)) < 0 {
		return errors.New("fallback_link_price must be set to a non-negative integer")
	}
	if c.FallbackNativePrice == nil || c.FallbackNativePrice.Cmp(big.NewInt(0)) < 0 {
		return errors.New("fallback_native_price must be set to a non-negative integer")
	}
	if c.MaxCheckDataSize == nil || *c.MaxCheckDataSize < 1 {
		return errors.New("max_check_data_size must be set to a positive integer")
	}
	if c.MaxPerformDataSize == nil || *c.MaxPerformDataSize < 1 {
		return errors.New("max_perform_data_size must be set to a positive integer")
	}
	if c.MaxRevertDataSize == nil || *c.MaxRevertDataSize < 1 {
		return errors.New("max_revert_data_size must be set to a positive integer")
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

type Contracts struct {
	ShouldBeUsed            *bool                      `toml:"use"`
	LinkTokenAddress        *string                    `toml:"link_token"`
	WethAddress             *string                    `toml:"weth"`
	TranscoderAddress       *string                    `toml:"transcoder"`
	ChainModuleAddress      *string                    `toml:"chain_module"`
	RegistryAddress         *string                    `toml:"registry"`
	RegistrarAddress        *string                    `toml:"registrar"`
	LinkEthFeedAddress      *string                    `toml:"link_eth_feed"`
	EthGasFeedAddress       *string                    `toml:"eth_gas_feed"`
	EthUSDFeedAddress       *string                    `toml:"eth_usd_feed"`
	LinkUSDFeedAddress      *string                    `toml:"link_usd_feed"`
	UpkeepContractAddresses []string                   `toml:"upkeep_contracts"`
	MultiCallAddress        *string                    `toml:"multicall"`
	Settings                map[string]ContractSetting `toml:"Settings"`
}

func (o *Contracts) Validate() error {
	if o.LinkTokenAddress != nil && !common.IsHexAddress(*o.LinkTokenAddress) {
		return errors.New("link_token must be a valid ethereum address")
	}
	if o.WethAddress != nil && !common.IsHexAddress(*o.WethAddress) {
		return errors.New("weth must be a valid ethereum address")
	}
	if o.TranscoderAddress != nil && !common.IsHexAddress(*o.TranscoderAddress) {
		return errors.New("transcoder must be a valid ethereum address")
	}
	if o.ChainModuleAddress != nil && !common.IsHexAddress(*o.ChainModuleAddress) {
		return errors.New("chain_module must be a valid ethereum address")
	}
	if o.RegistryAddress != nil && !common.IsHexAddress(*o.RegistryAddress) {
		return errors.New("registry must be a valid ethereum address")
	}
	if o.RegistrarAddress != nil && !common.IsHexAddress(*o.RegistrarAddress) {
		return errors.New("registrar must be a valid ethereum address")
	}
	if o.LinkEthFeedAddress != nil && !common.IsHexAddress(*o.LinkEthFeedAddress) {
		return errors.New("link_eth_feed must be a valid ethereum address")
	}
	if o.EthGasFeedAddress != nil && !common.IsHexAddress(*o.EthGasFeedAddress) {
		return errors.New("eth_gas_feed must be a valid ethereum address")
	}
	if o.EthUSDFeedAddress != nil && !common.IsHexAddress(*o.EthUSDFeedAddress) {
		return errors.New("eth_usd_feed must be a valid ethereum address")
	}
	if o.LinkUSDFeedAddress != nil && !common.IsHexAddress(*o.LinkUSDFeedAddress) {
		return errors.New("link_usd_feed must be a valid ethereum address")
	}
	if o.MultiCallAddress != nil && !common.IsHexAddress(*o.MultiCallAddress) {
		return errors.New("multicall must be a valid ethereum address")
	}
	if o.UpkeepContractAddresses != nil {
		allEnabled := make(map[bool]int)
		allConfigure := make(map[bool]int)
		for _, address := range o.UpkeepContractAddresses {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("upkeep_contracts must be valid ethereum addresses, but %s is not", address)
			}

			if v, ok := o.Settings[address]; ok {
				if v.ShouldBeUsed != nil {
					allEnabled[*v.ShouldBeUsed]++
				} else {
					allEnabled[true]++
				}
				if v.Configure != nil {
					allConfigure[*v.Configure]++
				} else {
					allConfigure[true]++
				}
			}
		}

		if allEnabled[true] > 0 && allEnabled[false] > 0 {
			return errors.New("either all or none offchain_aggregators must be used")
		}

		if allConfigure[true] > 0 && allConfigure[false] > 0 {
			return errors.New("either all or none offchain_aggregators must be configured")
		}
	}

	return nil
}

func (c *Config) UseExistingContracts() bool {
	if c.Contracts == nil {
		return false
	}

	if c.Contracts.ShouldBeUsed != nil {
		return *c.Contracts.ShouldBeUsed
	}

	return false
}

func (c *Config) LinkTokenContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.LinkTokenAddress != nil {
		return common.HexToAddress(*c.Contracts.LinkTokenAddress), nil
	}

	return common.Address{}, errors.New("link token address must be set")
}

func (c *Config) WethContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.WethAddress != nil {
		return common.HexToAddress(*c.Contracts.WethAddress), nil
	}

	return common.Address{}, errors.New("weth address must be set")
}

func (c *Config) TranscoderContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.TranscoderAddress != nil {
		return common.HexToAddress(*c.Contracts.TranscoderAddress), nil
	}

	return common.Address{}, errors.New("transcoder address must be set")
}

func (c *Config) ChainModuleContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.ChainModuleAddress != nil {
		return common.HexToAddress(*c.Contracts.ChainModuleAddress), nil
	}

	return common.Address{}, errors.New("chain module address must be set")
}

func (c *Config) RegistryContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.RegistryAddress != nil {
		return common.HexToAddress(*c.Contracts.RegistryAddress), nil
	}

	return common.Address{}, errors.New("registry address must be set")
}

func (c *Config) RegistrarContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.RegistrarAddress != nil {
		return common.HexToAddress(*c.Contracts.RegistrarAddress), nil
	}

	return common.Address{}, errors.New("registrar address must be set")
}

func (c *Config) LinkEthFeedContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.LinkEthFeedAddress != nil {
		return common.HexToAddress(*c.Contracts.LinkEthFeedAddress), nil
	}

	return common.Address{}, errors.New("link eth feed address must be set")
}

func (c *Config) EthGasFeedContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.EthGasFeedAddress != nil {
		return common.HexToAddress(*c.Contracts.EthGasFeedAddress), nil
	}

	return common.Address{}, errors.New("eth gas feed address must be set")
}

func (c *Config) EthUSDFeedContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.EthUSDFeedAddress != nil {
		return common.HexToAddress(*c.Contracts.EthUSDFeedAddress), nil
	}

	return common.Address{}, errors.New("eth usd feed address must be set")
}

func (c *Config) LinkUSDFeedContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.LinkUSDFeedAddress != nil {
		return common.HexToAddress(*c.Contracts.LinkUSDFeedAddress), nil
	}

	return common.Address{}, errors.New("link usd feed address must be set")
}

func (c *Config) UpkeepContractAddresses() ([]common.Address, error) {
	if c.Contracts != nil && c.Contracts.UpkeepContractAddresses != nil {
		addresses := make([]common.Address, len(c.Contracts.UpkeepContractAddresses))
		for i, address := range c.Contracts.UpkeepContractAddresses {
			addresses[i] = common.HexToAddress(address)
		}
		return addresses, nil
	}

	return nil, errors.New("upkeep contract addresses must be set")
}

func (c *Config) MultiCallContractAddress() (common.Address, error) {
	if c.Contracts != nil && c.Contracts.MultiCallAddress != nil {
		return common.HexToAddress(*c.Contracts.MultiCallAddress), nil
	}

	return common.Address{}, errors.New("multicall address must be set")
}

func (c *Config) UseExistingLinkTokenContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.LinkTokenAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.LinkTokenAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingWethContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.WethAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.WethAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingTranscoderContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.TranscoderAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.TranscoderAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingRegistryContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.RegistryAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.RegistryAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingRegistrarContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.RegistrarAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.RegistrarAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingLinkEthFeedContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.LinkEthFeedAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.LinkEthFeedAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingEthGasFeedContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.EthGasFeedAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.EthGasFeedAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingEthUSDFeedContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.EthUSDFeedAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.EthUSDFeedAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingLinkUSDFeedContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.LinkUSDFeedAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.LinkUSDFeedAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

func (c *Config) UseExistingUpkeepContracts() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.UpkeepContractAddresses == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	for _, address := range c.Contracts.UpkeepContractAddresses {
		if v, ok := c.Contracts.Settings[address]; ok {
			if v.ShouldBeUsed != nil && *v.ShouldBeUsed {
				return true
			}
		}
	}

	return false
}

func (c *Config) UseExistingMultiCallContract() bool {
	if !c.UseExistingContracts() {
		return false
	}

	if c.Contracts.MultiCallAddress == nil {
		return false
	}

	if len(c.Contracts.Settings) == 0 {
		return true
	}

	if v, ok := c.Contracts.Settings[*c.Contracts.MultiCallAddress]; ok {
		return v.ShouldBeUsed != nil && *v.ShouldBeUsed
	}

	return true
}

type ContractSetting struct {
	ShouldBeUsed *bool `toml:"use"`
	Configure    *bool `toml:"configure"`
}
