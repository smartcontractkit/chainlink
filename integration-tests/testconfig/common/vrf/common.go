package vrf

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

type Config struct {
	General           *General           `toml:"General"`
	ExistingEnvConfig *ExistingEnvConfig `toml:"ExistingEnv"`
	Performance       *PerformanceConfig `toml:"Performance"`
}

const (
	ErrDeviationShouldBeLessThanOriginal = "`RandomnessRequestCountPerRequestDeviation` should be less than `RandomnessRequestCountPerRequest`"
)

func (c *Config) Validate() error {
	if c.General != nil {
		if err := c.General.Validate(); err != nil {
			return err
		}
	}
	if c.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.Validate(); err != nil {
			return err
		}
	}
	if c.Performance != nil {
		if err := c.Performance.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type PerformanceConfig struct {
	TestDuration          *blockchain.StrDuration `toml:"test_duration"`
	RPS                   *int64                  `toml:"rps"`
	RateLimitUnitDuration *blockchain.StrDuration `toml:"rate_limit_unit_duration"`

	BHSTestDuration              *blockchain.StrDuration `toml:"bhs_test_duration"`
	BHSTestRPS                   *int64                  `toml:"bhs_test_rps"`
	BHSTestRateLimitUnitDuration *blockchain.StrDuration `toml:"bhs_test_rate_limit_unit_duration"`
}

func (c *PerformanceConfig) Validate() error {
	if c.TestDuration == nil || c.TestDuration.Duration == 0 {
		return errors.New("test_duration must be set to a positive value")
	}
	if c.RPS == nil || *c.RPS == 0 {
		return errors.New("rps must be set to a positive value")
	}
	if c.RateLimitUnitDuration == nil {
		return errors.New("rate_limit_unit_duration must be set ")
	}
	if c.BHSTestDuration == nil || c.BHSTestDuration.Duration == 0 {
		return errors.New("bhs_test_duration must be set to a positive value")
	}
	if c.BHSTestRPS == nil || *c.BHSTestRPS == 0 {
		return errors.New("bhs_test_rps must be set to a positive value")
	}
	if c.BHSTestRateLimitUnitDuration == nil {
		return errors.New("bhs_test_rate_limit_unit_duration must be set ")
	}

	return nil
}

type ExistingEnvConfig struct {
	CoordinatorAddress            *string  `toml:"coordinator_address"`
	ConsumerAddress               *string  `toml:"consumer_address"`
	LinkAddress                   *string  `toml:"link_address"`
	KeyHash                       *string  `toml:"key_hash"`
	CreateFundSubsAndAddConsumers *bool    `toml:"create_fund_subs_and_add_consumers"`
	NodeSendingKeys               []string `toml:"node_sending_keys"`
	Funding
}

func (c *ExistingEnvConfig) Validate() error {
	if c.CreateFundSubsAndAddConsumers == nil {
		return errors.New("create_fund_subs_and_add_consumers must be set ")
	}
	if c.CoordinatorAddress == nil {
		return errors.New("coordinator_address must be set when using existing environment")
	}
	if !common.IsHexAddress(*c.CoordinatorAddress) {
		return errors.New("coordinator_address must be a valid hex address")
	}
	if c.KeyHash == nil {
		return errors.New("key_hash must be set when using existing environment")
	}
	if *c.KeyHash == "" {
		return errors.New("key_hash must be a non-empty string")
	}
	if *c.CreateFundSubsAndAddConsumers {
		if err := c.Funding.Validate(); err != nil {
			return err
		}
	} else {
		if c.ConsumerAddress == nil || *c.ConsumerAddress == "" {
			return errors.New("consumer_address must be set when using existing environment")
		}
		if !common.IsHexAddress(*c.ConsumerAddress) {
			return errors.New("consumer_address must be a valid hex address")
		}
	}

	if c.NodeSendingKeys != nil {
		for _, key := range c.NodeSendingKeys {
			if !common.IsHexAddress(key) {
				return errors.New("node_sending_keys must be a valid hex address")
			}
		}
	}

	return nil
}

type Funding struct {
	NodeSendingKeyFundingMin *float64 `toml:"node_sending_key_funding_min"`
}

func (c *Funding) Validate() error {
	if c.NodeSendingKeyFundingMin != nil && *c.NodeSendingKeyFundingMin <= 0 {
		return errors.New("when set node_sending_key_funding_min must be a positive value")
	}

	return nil
}

type General struct {
	UseExistingEnv                  *bool    `toml:"use_existing_env"`
	CancelSubsAfterTestRun          *bool    `toml:"cancel_subs_after_test_run"`
	CLNodeMaxGasPriceGWei           *int64   `toml:"cl_node_max_gas_price_gwei"`         // Max gas price in GWei for the chainlink node
	LinkNativeFeedResponse          *int64   `toml:"link_native_feed_response"`          // Response of the LINK/ETH feed
	MinimumConfirmations            *uint16  `toml:"minimum_confirmations"`              // Minimum number of confirmations for the VRF Coordinator
	SubscriptionFundingAmountLink   *float64 `toml:"subscription_funding_amount_link"`   // Amount of LINK to fund the subscription with
	SubscriptionRefundingAmountLink *float64 `toml:"subscription_refunding_amount_link"` // Amount of LINK to fund the subscription with
	NumberOfWords                   *uint32  `toml:"number_of_words"`                    // Number of words to request
	CallbackGasLimit                *uint32  `toml:"callback_gas_limit"`                 // Gas limit for the callback
	MaxGasLimitCoordinatorConfig    *uint32  `toml:"max_gas_limit_coordinator_config"`   // Max gas limit for the VRF Coordinator config
	FallbackWeiPerUnitLink          *string  `toml:"fallback_wei_per_unit_link"`         // Fallback wei per unit LINK for the VRF Coordinator config
	StalenessSeconds                *uint32  `toml:"staleness_seconds"`                  // Staleness in seconds for the VRF Coordinator config
	GasAfterPaymentCalculation      *uint32  `toml:"gas_after_payment_calculation"`      // Gas after payment calculation for the VRF Coordinator

	NumberOfSubToCreate         *int `toml:"number_of_sub_to_create"`          // Number of subscriptions to create
	NumberOfSendingKeysToCreate *int `toml:"number_of_sending_keys_to_create"` // Number of sending keys to create

	RandomnessRequestCountPerRequest          *uint16 `toml:"randomness_request_count_per_request"`           // How many randomness requests to send per request
	RandomnessRequestCountPerRequestDeviation *uint16 `toml:"randomness_request_count_per_request_deviation"` // How many randomness requests to send per request

	RandomWordsFulfilledEventTimeout *blockchain.StrDuration `toml:"random_words_fulfilled_event_timeout"` // How long to wait for the RandomWordsFulfilled event to be emitted
	WaitFor256BlocksTimeout          *blockchain.StrDuration `toml:"wait_for_256_blocks_timeout"`          // How long to wait for 256 blocks to be mined

	// Wrapper Config
	WrapperGasOverhead      *uint32 `toml:"wrapped_gas_overhead"`
	WrapperMaxNumberOfWords *uint8  `toml:"wrapper_max_number_of_words"`

	WrapperConsumerFundingAmountNativeToken *float64 `toml:"wrapper_consumer_funding_amount_native_token"`
	WrapperConsumerFundingAmountLink        *int64   `toml:"wrapper_consumer_funding_amount_link"`

	//VRF Job Config
	VRFJobForwardingAllowed             *bool                   `toml:"vrf_job_forwarding_allowed"`
	VRFJobEstimateGasMultiplier         *float64                `toml:"vrf_job_estimate_gas_multiplier"`
	VRFJobBatchFulfillmentEnabled       *bool                   `toml:"vrf_job_batch_fulfillment_enabled"`
	VRFJobBatchFulfillmentGasMultiplier *float64                `toml:"vrf_job_batch_fulfillment_gas_multiplier"`
	VRFJobPollPeriod                    *blockchain.StrDuration `toml:"vrf_job_poll_period"`
	VRFJobRequestTimeout                *blockchain.StrDuration `toml:"vrf_job_request_timeout"`
	VRFJobSimulationBlock               *string                 `toml:"vrf_job_simulation_block"`

	//BHS Job Config
	BHSJobWaitBlocks     *int                    `toml:"bhs_job_wait_blocks"`
	BHSJobLookBackBlocks *int                    `toml:"bhs_job_lookback_blocks"`
	BHSJobPollPeriod     *blockchain.StrDuration `toml:"bhs_job_poll_period"`
	BHSJobRunTimeout     *blockchain.StrDuration `toml:"bhs_job_run_timeout"`

	//BHF Job Config
	BHFJobWaitBlocks     *int                    `toml:"bhf_job_wait_blocks"`
	BHFJobLookBackBlocks *int                    `toml:"bhf_job_lookback_blocks"`
	BHFJobPollPeriod     *blockchain.StrDuration `toml:"bhf_job_poll_period"`
	BHFJobRunTimeout     *blockchain.StrDuration `toml:"bhf_job_run_timeout"`
}

func (c *General) Validate() error {
	if c.UseExistingEnv == nil {
		return errors.New("use_existing_env must not be nil")
	}
	if c.CLNodeMaxGasPriceGWei == nil || *c.CLNodeMaxGasPriceGWei == 0 {
		return errors.New("cl_node_max_gas_price_gwei must be set to a positive value")
	}
	if c.LinkNativeFeedResponse == nil || *c.LinkNativeFeedResponse == 0 {
		return errors.New("link_native_feed_response must be set to a positive value")
	}
	if c.MinimumConfirmations == nil {
		return errors.New("minimum_confirmations must be set to a non-negative value")
	}
	if c.SubscriptionFundingAmountLink == nil || *c.SubscriptionFundingAmountLink < 0 {
		return errors.New("subscription_funding_amount_link must be set to non-negative value")
	}
	if c.SubscriptionRefundingAmountLink == nil || *c.SubscriptionRefundingAmountLink < 0 {
		return errors.New("subscription_refunding_amount_link must be set to non-negative value")
	}
	if c.NumberOfWords == nil || *c.NumberOfWords == 0 {
		return errors.New("number_of_words must be set to a positive value")
	}
	if c.CallbackGasLimit == nil || *c.CallbackGasLimit == 0 {
		return errors.New("callback_gas_limit must be set to a positive value")
	}
	if c.MaxGasLimitCoordinatorConfig == nil || *c.MaxGasLimitCoordinatorConfig == 0 {
		return errors.New("max_gas_limit_coordinator_config must be set to a positive value")
	}
	if c.FallbackWeiPerUnitLink == nil {
		return errors.New("fallback_wei_per_unit_link must be set")
	}
	if c.StalenessSeconds == nil || *c.StalenessSeconds == 0 {
		return errors.New("staleness_seconds must be set to a positive value")
	}
	if c.GasAfterPaymentCalculation == nil || *c.GasAfterPaymentCalculation == 0 {
		return errors.New("gas_after_payment_calculation must be set to a positive value")
	}
	if c.NumberOfSubToCreate == nil || *c.NumberOfSubToCreate == 0 {
		return errors.New("number_of_sub_to_create must be set to a positive value")
	}

	if c.NumberOfSendingKeysToCreate == nil || *c.NumberOfSendingKeysToCreate < 0 {
		return errors.New("number_of_sending_keys_to_create must be set to 0 or a positive value")
	}

	if c.RandomnessRequestCountPerRequest == nil || *c.RandomnessRequestCountPerRequest == 0 {
		return errors.New("randomness_request_count_per_request must be set to a positive value")
	}
	if c.RandomnessRequestCountPerRequestDeviation == nil {
		return errors.New("randomness_request_count_per_request_deviation must be set to a non-negative value")
	}
	if c.RandomWordsFulfilledEventTimeout == nil || c.RandomWordsFulfilledEventTimeout.Duration == 0 {
		return errors.New("random_words_fulfilled_event_timeout must be set to a positive value")
	}
	if c.WaitFor256BlocksTimeout == nil || c.WaitFor256BlocksTimeout.Duration == 0 {
		return errors.New("wait_for_256_blocks_timeout must be set to a positive value")
	}
	if c.WrapperGasOverhead == nil {
		return errors.New("wrapped_gas_overhead must be set to a non-negative value")
	}
	if c.WrapperMaxNumberOfWords == nil || *c.WrapperMaxNumberOfWords == 0 {
		return errors.New("wrapper_max_number_of_words must be set to a positive value")
	}
	if c.WrapperConsumerFundingAmountNativeToken == nil || *c.WrapperConsumerFundingAmountNativeToken < 0 {
		return errors.New("wrapper_consumer_funding_amount_native_token must be set to a non-negative value")
	}
	if c.WrapperConsumerFundingAmountLink == nil || *c.WrapperConsumerFundingAmountLink < 0 {
		return errors.New("wrapper_consumer_funding_amount_link must be set to a non-negative value")
	}
	if *c.RandomnessRequestCountPerRequest <= *c.RandomnessRequestCountPerRequestDeviation {
		return errors.New(ErrDeviationShouldBeLessThanOriginal)
	}

	if c.VRFJobForwardingAllowed == nil {
		return errors.New("vrf_job_forwarding_allowed must be set")
	}

	if c.VRFJobBatchFulfillmentEnabled == nil {
		return errors.New("vrf_job_batch_fulfillment_enabled must be set")
	}
	if c.VRFJobEstimateGasMultiplier == nil || *c.VRFJobEstimateGasMultiplier < 0 {
		return errors.New("vrf_job_estimate_gas_multiplier must be set to a non-negative value")
	}
	if c.VRFJobBatchFulfillmentGasMultiplier == nil || *c.VRFJobBatchFulfillmentGasMultiplier < 0 {
		return errors.New("vrf_job_batch_fulfillment_gas_multiplier must be set to a non-negative value")
	}

	if c.VRFJobPollPeriod == nil || c.VRFJobPollPeriod.Duration == 0 {
		return errors.New("vrf_job_poll_period must be set to a non-negative value")
	}

	if c.VRFJobRequestTimeout == nil || c.VRFJobRequestTimeout.Duration == 0 {
		return errors.New("vrf_job_request_timeout must be set to a non-negative value")
	}

	if c.VRFJobSimulationBlock != nil && (*c.VRFJobSimulationBlock != "latest" && *c.VRFJobSimulationBlock != "pending") {
		return errors.New("simulation_block must be nil or \"latest\" or \"pending\"")
	}

	if c.BHSJobLookBackBlocks == nil || *c.BHSJobLookBackBlocks < 0 {
		return errors.New("bhs_job_lookback_blocks must be set to a non-negative value")
	}

	if c.BHSJobPollPeriod == nil || c.BHSJobPollPeriod.Duration == 0 {
		return errors.New("bhs_job_poll_period must be set to a non-negative value")
	}

	if c.BHSJobRunTimeout == nil || c.BHSJobRunTimeout.Duration == 0 {
		return errors.New("bhs_job_run_timeout must be set to a non-negative value")
	}

	if c.BHSJobWaitBlocks == nil || *c.BHSJobWaitBlocks < 0 {
		return errors.New("bhs_job_wait_blocks must be set to a non-negative value")
	}

	return nil
}
