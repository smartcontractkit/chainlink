package testconfig

import (
	"errors"
	"time"
)

const (
	ErrDeviationShouldBeLessThanOriginal = "`RandomnessRequestCountPerRequestDeviation` should be less than `RandomnessRequestCountPerRequest`"
)

type Config struct {
	Common            *Common            `toml:"Common"`
	General           *General           `toml:"General"`
	ExistingEnvConfig *ExistingEnvConfig `toml:"ExistingEnvConfig"`
	NewEnvConfig      *NewEnvConfig      `toml:"NewEnvConfig"`
	Performance       *PerformanceConfig `toml:"PerformanceConfig"`
}

type Common struct {
	CancelSubsAfterTestRun bool `toml:"cancel_subs_after_test_run"`
}

type PerformanceConfig struct {
	TestDuration          time.Duration `toml:"test_duration"`            // How long to run the test for  default:"3m"
	RPS                   int64         `toml:"rps"`                      // How many requests per second to send default:"1"
	RateLimitUnitDuration time.Duration `toml:"rate_limit_unit_duration"` // default:"1m"

	// Using existing environment and contracts
	UseExistingEnv     bool   `toml:"use_existing_env"`    // Whether to use an existing environment or create a new one  default:"false"
	CoordinatorAddress string `toml:"coordinator_address"` // Coordinator address
	ConsumerAddress    string `toml:"consumer_address"`    // Consumer address
	LinkAddress        string `toml:"link_address"`        // Link address
	SubID              uint64 `toml:"sub_id"`              // Subscription ID
	KeyHash            string `toml:"key_hash"`
}

type ExistingEnvConfig struct {
	CoordinatorAddress            string   `toml:"coordinator_address"`
	ConsumerAddress               string   `toml:"consumer_address"`
	LinkAddress                   string   `toml:"link_address"`
	SubID                         uint64   `toml:"sub_id"`
	KeyHash                       string   `toml:"key_hash"`
	CreateFundSubsAndAddConsumers bool     `toml:"create_fund_subs_and_add_consumers"`
	NodeSendingKeys               []string `toml:"node_sending_keys"`
	Funding
}

type NewEnvConfig struct {
	Funding
}

type Funding struct {
	SubFunding
	NodeSendingKeyFunding    float64 `toml:"node_sending_key_funding"`
	NodeSendingKeyFundingMin float64 `toml:"node_sending_key_funding_min"`
}

type SubFunding struct {
	SubFundsLink float64 `toml:"sub_funds_link"`
}

type General struct {
	CLNodeMaxGasPriceGWei int64 `toml:"max_gas_price_gwei"` // Max gas price in GWei for the chainlink node default:"1000"
	// IsNativePayment                bool    `toml:"is_native_payment"`                   // Whether to use native payment or LINK token default:"false"
	LinkNativeFeedResponse         int64   `toml:"link_native_feed_response"`           // Response of the LINK/ETH feed default:"1000000000000000000"
	MinimumConfirmations           uint16  `toml:"minimum_confirmations" `              // Minimum number of confirmations for the VRF Coordinator default:"3"
	SubscriptionFundingAmountLink  float64 `toml:"subscription_funding_amount_link"`    // Amount of LINK to fund the subscription with default:"5"
	NumberOfWords                  uint32  `toml:"number_of_words" `                    // Number of words to request default:"3"
	CallbackGasLimit               uint32  `toml:"callback_gas_limit" `                 // Gas limit for the callback default:"1000000"
	MaxGasLimitCoordinatorConfig   uint32  `toml:"max_gas_limit_coordinator_config"`    // Max gas limit for the VRF Coordinator config  default:"2500000"
	FallbackWeiPerUnitLink         int64   `toml:"fallback_wei_per_unit_link"`          // Fallback wei per unit LINK for the VRF Coordinator config  default:"60000000000000000"
	StalenessSeconds               uint32  `toml:"staleness_seconds" `                  // Staleness in seconds for the VRF Coordinator config default:"86400"
	GasAfterPaymentCalculation     uint32  `toml:"gas_after_payment_calculation" `      // Gas after payment calculation for the VRF Coordinator config default:"33825"
	FulfillmentFlatFeeLinkPPMTier1 uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_1"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier2 uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_2"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier3 uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_3"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier4 uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_4"` //default:"500"
	FulfillmentFlatFeeLinkPPMTier5 uint32  `toml:"fulfilment_flat_fee_link_ppm_tier_5"` //default:"500"
	ReqsForTier2                   int64   `toml:"reqs_for_tier_2"`                     // default:"0"
	ReqsForTier3                   int64   `toml:"reqs_for_tier_2"`                     // default:"0"
	ReqsForTier4                   int64   `toml:"reqs_for_tier_3"`                     // default:"0"
	ReqsForTier5                   int64   `toml:"reqs_for_tier_4"`                     // default:"0"

	NumberOfSubToCreate int `toml:"number_of_sub_to_create"` // Number of subscriptions to create default:"1"

	RandomnessRequestCountPerRequest          uint16 `toml:"randomness_request_count_per_request"`           // How many randomness requests to send per request default:"1"
	RandomnessRequestCountPerRequestDeviation uint16 `toml:"randomness_request_count_per_request_deviation"` // How many randomness requests to send per request  default:"0"

	RandomWordsFulfilledEventTimeout time.Duration `toml:"random_words_fulfilled_event_timeout"` // How long to wait for the RandomWordsFulfilled event to be emitted default:"2m"

	// Wrapper Config
	WrapperGasOverhead                      uint32  `toml:"wrapped_gas_overhead"`                         // default:"50000"
	CoordinatorGasOverhead                  uint32  `toml:"coordinator_gas_overhead"`                     // default:"52000"
	WrapperPremiumPercentage                uint8   `toml:"wrapper_premium_percentage"`                   // default:"25"
	WrapperMaxNumberOfWords                 uint8   `toml:"wrapper_max_number_of_words"`                  // default:"10"
	WrapperConsumerFundingAmountNativeToken float64 `toml:"wrapper_consumer_funding_amount_native_token"` // default:"1"
	WrapperConsumerFundingAmountLink        int64   `toml:"wrapper_consumer_funding_amount_link"`         // default:"10"
}

func (c *Config) ApplyOverrides(_ *Config) error {
	//TODO implement me
	return nil
}

func (c *Config) Validate() error {
	if c.General.RandomnessRequestCountPerRequest <= c.General.RandomnessRequestCountPerRequestDeviation {
		return errors.New(ErrDeviationShouldBeLessThanOriginal)
	}

	return nil
}
