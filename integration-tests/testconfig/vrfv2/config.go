package testconfig

import (
	"errors"

	vrf_common_config "github.com/smartcontractkit/chainlink/integration-tests/testconfig/common/vrf"
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

type ExistingEnvConfig struct {
	*vrf_common_config.ExistingEnvConfig
	SubID *uint64 `toml:"sub_id"`
}

func (c *ExistingEnvConfig) Validate() error {
	if c.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.Validate(); err != nil {
			return err
		}
	}
	if !*c.CreateFundSubsAndAddConsumers {
		if c.SubID == nil {
			return errors.New("sub_id must be set when using existing environment")
		}

		if *c.SubID == 0 {
			return errors.New("sub_id must be positive value")
		}
	}

	return c.Funding.Validate()
}

type General struct {
	*vrf_common_config.General
	FulfillmentFlatFeeLinkPPMTier1 *uint32 `toml:"fulfilment_flat_fee_link_ppm_tier_1"`
	FulfillmentFlatFeeLinkPPMTier2 *uint32 `toml:"fulfilment_flat_fee_link_ppm_tier_2"`
	FulfillmentFlatFeeLinkPPMTier3 *uint32 `toml:"fulfilment_flat_fee_link_ppm_tier_3"`
	FulfillmentFlatFeeLinkPPMTier4 *uint32 `toml:"fulfilment_flat_fee_link_ppm_tier_4"`
	FulfillmentFlatFeeLinkPPMTier5 *uint32 `toml:"fulfilment_flat_fee_link_ppm_tier_5"`
	ReqsForTier2                   *int64  `toml:"reqs_for_tier_2"`
	ReqsForTier3                   *int64  `toml:"reqs_for_tier_3"`
	ReqsForTier4                   *int64  `toml:"reqs_for_tier_4"`
	ReqsForTier5                   *int64  `toml:"reqs_for_tier_5"`
	CoordinatorGasOverhead         *uint32 `toml:"coordinator_gas_overhead"`
	WrapperPremiumPercentage       *uint8  `toml:"wrapper_premium_percentage"`
}

func (c *General) Validate() error {
	if c.General != nil {
		if err := c.General.Validate(); err != nil {
			return err
		}
	}
	if c.FulfillmentFlatFeeLinkPPMTier1 == nil || *c.FulfillmentFlatFeeLinkPPMTier1 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_1 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier2 == nil || *c.FulfillmentFlatFeeLinkPPMTier2 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_2 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier3 == nil || *c.FulfillmentFlatFeeLinkPPMTier3 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_3 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier4 == nil || *c.FulfillmentFlatFeeLinkPPMTier4 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_4 must be set to a positive value")
	}
	if c.FulfillmentFlatFeeLinkPPMTier5 == nil || *c.FulfillmentFlatFeeLinkPPMTier5 == 0 {
		return errors.New("fulfilment_flat_fee_link_ppm_tier_5 must be set to a positive value")
	}
	if c.ReqsForTier2 == nil || *c.ReqsForTier2 < 0 {
		return errors.New("reqs_for_tier_2 must be set to a non-negative value")
	}
	if c.ReqsForTier3 == nil || *c.ReqsForTier3 < 0 {
		return errors.New("reqs_for_tier_3 must be set to a non-negative value")
	}
	if c.ReqsForTier4 == nil || *c.ReqsForTier4 < 0 {
		return errors.New("reqs_for_tier_4 must be set to a non-negative value")
	}
	if c.ReqsForTier5 == nil || *c.ReqsForTier5 < 0 {
		return errors.New("reqs_for_tier_5 must be set to a non-negative value")
	}
	if c.CoordinatorGasOverhead == nil || *c.CoordinatorGasOverhead == 0 {
		return errors.New("coordinator_gas_overhead must be set to a non-negative value")
	}
	if c.WrapperPremiumPercentage == nil || *c.WrapperPremiumPercentage == 0 {
		return errors.New("wrapper_premium_percentage must be set to a positive value")
	}
	return nil
}
