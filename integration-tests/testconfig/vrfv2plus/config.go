package testconfig

import (
	vrfv2 "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
)

type Config struct {
	Common            *Common                 `toml:"Common"`
	General           *General                `toml:"General"`
	ExistingEnvConfig *ExistingEnvConfig      `toml:"ExistingEnvConfig"`
	NewEnvConfig      *NewEnvConfig           `toml:"NewEnvConfig"`
	Performance       vrfv2.PerformanceConfig `toml:"Performance"`
}

type Common struct {
	vrfv2.Common
}

type General struct {
	vrfv2.General
	SubscriptionFundingAmountNative *float64 `toml:"subscription_funding_amount_native"` // Amount of LINK to fund the subscription with default:"1"
	FulfillmentFlatFeeLinkPPM       *uint32  `toml:"fulfillment_flat_fee_link_ppm"`      // Flat fee in ppm for LINK for the VRF Coordinator config default:"500"
	FulfillmentFlatFeeNativePPM     *uint32  `toml:"fulfillment_flat_fee_native_ppm"`    // Flat fee in ppm for native currency for the VRF Coordinator config default:"500"
}

type NewEnvConfig struct {
	Funding
}

type ExistingEnvConfig struct {
	*vrfv2.ExistingEnvConfig
	Funding
}

type Funding struct {
	SubFunding
	NodeSendingKeyFunding    *float64 `toml:"node_sending_key_funding"`
	NodeSendingKeyFundingMin *float64 `toml:"node_sending_key_funding_min"`
}

type SubFunding struct {
	SubFundsLink   *float64 `toml:"sub_funds_link"`
	SubFundsNative *float64 `toml:"sub_funds_native"`
}

func (c *Config) ApplyOverrides(_ *Config) error {
	//TODO implement me
	return nil
}

func (c *Config) Validate() error {
	//TODO implement me
	return nil
}
