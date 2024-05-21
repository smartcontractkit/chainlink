package automation

import (
	"errors"
	"math/big"
	"time"
)

type Config struct {
	General          *General          `toml:"General"`
	Load             []Load            `toml:"Load"`
	DataStreams      *DataStreams      `toml:"DataStreams"`
	AutomationConfig *AutomationConfig `toml:"AutomationConfig"`
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
	URL           *string `toml:"url"`
	Username      *string `toml:"username"`
	Password      *string `toml:"password"`
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
	UseLogBufferV1   *bool             `toml:"use_log_buffer_v1"`
}

func (c *AutomationConfig) Validate() error {
	if err := c.PluginConfig.Validate(); err != nil {
		return err
	}
	if err := c.PublicConfig.Validate(); err != nil {
		return err
	}
	if err := c.RegistrySettings.Validate(); err != nil {
		return err
	}
	if c.UseLogBufferV1 == nil {
		return errors.New("use_log_buffer_v1 must be set")
	}
	return nil
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
