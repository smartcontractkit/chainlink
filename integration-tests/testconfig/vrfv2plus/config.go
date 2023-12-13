package testconfig

import (
	"errors"

	vrfv2 "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
)

type Config struct {
	Common            *Common            `toml:"Common"`
	ExistingEnvConfig *ExistingEnvConfig `toml:"ExistingEnvConfig"`
	NewEnvConfig      *NewEnvConfig      `toml:"NewEnvConfig"`
}

type Common struct {
	vrfv2.Common
	SubscriptionFundingAmountNative float64 `toml:"subscription_funding_amount_native"` // Amount of LINK to fund the subscription with default:"1"
	FulfillmentFlatFeeLinkPPM       uint32  `toml:"fulfillment_flat_fee_link_ppm"`      // Flat fee in ppm for LINK for the VRF Coordinator config default:"500"
	FulfillmentFlatFeeNativePPM     uint32  `toml:"fulfillment_flat_fee_native_ppm"`    // Flat fee in ppm for native currency for the VRF Coordinator config default:"500"
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
	SubFundsLink   float64 `toml:"sub_funds_link"`
	SubFundsNative float64 `toml:"sub_funds_native"`
}

func (c *Config) ApplyOverrides(from interface{}) error {
	//TODO implement me
	return nil
}

func (c *Config) Validate() error {
	if c.Common.RandomnessRequestCountPerRequest <= c.Common.RandomnessRequestCountPerRequestDeviation {
		return errors.New(vrfv2.ErrDeviationShouldBeLessThanOriginal)
	}

	return nil
}
