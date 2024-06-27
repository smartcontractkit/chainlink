package vrfv2plus_config

import "time"

type VRFV2PlusConfig struct {
	ChainlinkNodeFunding            float64 `envconfig:"CHAINLINK_NODE_FUNDING" default:".1"`                     // Amount of native currency to fund each chainlink node with
	IsNativePayment                 bool    `envconfig:"IS_NATIVE_PAYMENT" default:"false"`                       // Whether to use native payment or LINK token
	LinkNativeFeedResponse          int64   `envconfig:"LINK_NATIVE_FEED_RESPONSE" default:"1000000000000000000"` // Response of the LINK/ETH feed
	MinimumConfirmations            uint16  `envconfig:"MINIMUM_CONFIRMATIONS" default:"3"`                       // Minimum number of confirmations for the VRF Coordinator
	SubscriptionFundingAmountLink   int64   `envconfig:"SUBSCRIPTION_FUNDING_AMOUNT_LINK" default:"10"`           // Amount of LINK to fund the subscription with
	SubscriptionFundingAmountNative int64   `envconfig:"SUBSCRIPTION_FUNDING_AMOUNT_NATIVE" default:"1"`          // Amount of native currency to fund the subscription with
	NumberOfWords                   uint32  `envconfig:"NUMBER_OF_WORDS" default:"3"`                             // Number of words to request
	CallbackGasLimit                uint32  `envconfig:"CALLBACK_GAS_LIMIT" default:"1000000"`                    // Gas limit for the callback
	MaxGasLimitCoordinatorConfig    uint32  `envconfig:"MAX_GAS_LIMIT_COORDINATOR_CONFIG" default:"2500000"`      // Max gas limit for the VRF Coordinator config
	FallbackWeiPerUnitLink          int64   `envconfig:"FALLBACK_WEI_PER_UNIT_LINK" default:"60000000000000000"`  // Fallback wei per unit LINK for the VRF Coordinator config
	StalenessSeconds                uint32  `envconfig:"STALENESS_SECONDS" default:"86400"`                       // Staleness in seconds for the VRF Coordinator config
	GasAfterPaymentCalculation      uint32  `envconfig:"GAS_AFTER_PAYMENT_CALCULATION" default:"33825"`           // Gas after payment calculation for the VRF Coordinator config
	FulfillmentFlatFeeLinkPPM       uint32  `envconfig:"FULFILLMENT_FLAT_FEE_LINK_PPM" default:"500"`             // Flat fee in ppm for LINK for the VRF Coordinator config
	FulfillmentFlatFeeNativePPM     uint32  `envconfig:"FULFILLMENT_FLAT_FEE_NATIVE_PPM" default:"500"`           // Flat fee in ppm for native currency for the VRF Coordinator config

	RandomnessRequestCountPerRequest          uint16 `envconfig:"RANDOMNESS_REQUEST_COUNT_PER_REQUEST" default:"1"`           // How many randomness requests to send per request
	RandomnessRequestCountPerRequestDeviation uint16 `envconfig:"RANDOMNESS_REQUEST_COUNT_PER_REQUEST_DEVIATION" default:"0"` // How many randomness requests to send per request

	//Wrapper Config
	WrapperGasOverhead                      uint32  `envconfig:"WRAPPER_GAS_OVERHEAD" default:"50000"`
	CoordinatorGasOverhead                  uint32  `envconfig:"COORDINATOR_GAS_OVERHEAD" default:"52000"`
	WrapperPremiumPercentage                uint8   `envconfig:"WRAPPER_PREMIUM_PERCENTAGE" default:"25"`
	WrapperMaxNumberOfWords                 uint8   `envconfig:"WRAPPER_MAX_NUMBER_OF_WORDS" default:"10"`
	WrapperConsumerFundingAmountNativeToken float64 `envconfig:"WRAPPER_CONSUMER_FUNDING_AMOUNT_NATIVE_TOKEN" default:"1"`
	WrapperConsumerFundingAmountLink        int64   `envconfig:"WRAPPER_CONSUMER_FUNDING_AMOUNT_LINK" default:"10"`

	//LOAD/SOAK Test Config
	TestDuration          time.Duration `envconfig:"TEST_DURATION" default:"3m"` // How long to run the test for
	RPS                   int64         `envconfig:"RPS" default:"1"`            // How many requests per second to send
	RateLimitUnitDuration time.Duration `envconfig:"RATE_LIMIT_UNIT_DURATION" default:"1m"`
	//Using existing environment and contracts
	UseExistingEnv     bool   `envconfig:"USE_EXISTING_ENV" default:"false"` // Whether to use an existing environment or create a new one
	CoordinatorAddress string `envconfig:"COORDINATOR_ADDRESS" default:""`   // Coordinator address
	ConsumerAddress    string `envconfig:"CONSUMER_ADDRESS" default:""`      // Consumer address
	SubID              string `envconfig:"SUB_ID" default:""`                // Subscription ID
	KeyHash            string `envconfig:"KEY_HASH" default:""`
}
