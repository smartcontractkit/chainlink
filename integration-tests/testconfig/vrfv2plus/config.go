package testconfig

import (
	"errors"

	vrf_common_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/common/vrf"
)

type BillingType string

const (
	BillingType_Link            BillingType = "LINK"
	BillingType_Native          BillingType = "NATIVE"
	BillingType_Link_and_Native BillingType = "LINK_AND_NATIVE"
)

type Config struct {
	General           *General                             `toml:"General"`
	ExistingEnvConfig *ExistingEnvConfig                   `toml:"ExistingEnv"`
	Performance       *vrf_common_config.PerformanceConfig `toml:"Performance"`
}

func (c *Config) Validate() error {
	if c.General != nil {
		if err := c.General.Validate(); err != nil {
			return err
		}
	}
	if c.Performance != nil {
		if err := c.Performance.Validate(); err != nil {
			return err
		}
	}
	if c.ExistingEnvConfig != nil && *c.General.UseExistingEnv {
		if err := c.ExistingEnvConfig.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type General struct {
	*vrf_common_config.General
	SubscriptionBillingType           *string  `toml:"subscription_billing_type"`              // Billing type for the subscription
	SubscriptionFundingAmountNative   *float64 `toml:"subscription_funding_amount_native"`     // Amount of LINK to fund the subscription with
	SubscriptionRefundingAmountNative *float64 `toml:"subscription_refunding_amount_native"`   // Amount of LINK to fund the subscription with
	FulfillmentFlatFeeNativePPM       *uint32  `toml:"fulfillment_flat_fee_native_ppm"`        // Flat fee in ppm for native currency for the VRF Coordinator config
	FulfillmentFlatFeeLinkDiscountPPM *uint32  `toml:"fulfillment_flat_fee_link_discount_ppm"` // Flat fee discount in ppm for LINK for the VRF Coordinator config
	NativePremiumPercentage           *uint8   `toml:"native_premium_percentage"`              // Native Premium Percentage
	LinkPremiumPercentage             *uint8   `toml:"link_premium_percentage"`                // LINK Premium Percentage

	//Wrapper config
	CoordinatorGasOverheadPerWord      *uint16 `toml:"coordinator_gas_overhead_per_word"`
	CoordinatorGasOverheadNative       *uint32 `toml:"coordinator_gas_overhead_native"`
	CoordinatorGasOverheadLink         *uint32 `toml:"coordinator_gas_overhead_link"`
	CoordinatorNativePremiumPercentage *uint8  `toml:"coordinator_native_premium_percentage"`
	CoordinatorLinkPremiumPercentage   *uint8  `toml:"coordinator_link_premium_percentage"`
}

func (c *General) Validate() error {
	if err := c.General.Validate(); err != nil {
		return err
	}
	if c.SubscriptionBillingType == nil || *c.SubscriptionBillingType == "" {
		return errors.New("subscription_billing_type must be set to either: LINK, NATIVE, LINK_AND_NATIVE")
	}
	if c.SubscriptionFundingAmountNative == nil || *c.SubscriptionFundingAmountNative <= 0 {
		return errors.New("subscription_funding_amount_native must be greater than 0")
	}
	if c.SubscriptionRefundingAmountNative == nil || *c.SubscriptionRefundingAmountNative <= 0 {
		return errors.New("subscription_refunding_amount_native must be greater than 0")
	}
	if c.FulfillmentFlatFeeNativePPM == nil {
		return errors.New("fulfillment_flat_fee_native_ppm must not be nil")
	}
	if c.FulfillmentFlatFeeLinkDiscountPPM == nil {
		return errors.New("fulfillment_flat_fee_link_discount_ppm must not be nil")
	}
	if c.NativePremiumPercentage == nil {
		return errors.New("native_premium_percentage must not be nil")
	}
	if c.LinkPremiumPercentage == nil {
		return errors.New("link_premium_percentage must not be nil")
	}
	if c.CoordinatorGasOverheadPerWord == nil {
		return errors.New("coordinator_gas_overhead_per_word must not be nil")
	}
	if c.CoordinatorGasOverheadNative == nil || *c.CoordinatorGasOverheadNative == 0 {
		return errors.New("coordinator_gas_overhead_native must be set to a non-negative value")
	}
	if c.CoordinatorGasOverheadLink == nil || *c.CoordinatorGasOverheadLink == 0 {
		return errors.New("coordinator_gas_overhead_link must be set to a non-negative value")
	}
	if c.CoordinatorNativePremiumPercentage == nil {
		return errors.New("coordinator_native_premium_percentage must not be nil")
	}
	if c.CoordinatorLinkPremiumPercentage == nil {
		return errors.New("coordinator_link_premium_percentage must not be nil")
	}
	return nil
}

type ExistingEnvConfig struct {
	*vrf_common_config.ExistingEnvConfig
	SubID *string `toml:"sub_id"`
}

func (c *ExistingEnvConfig) Validate() error {
	if c.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.Validate(); err != nil {
			return err
		}
	}
	if !*c.CreateFundSubsAndAddConsumers {
		if c.SubID == nil && *c.SubID == "" {
			return errors.New("sub_id must be set when using existing environment")
		}
	}
	return c.Funding.Validate()
}
