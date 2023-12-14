package testconfig

import (
	"errors"

	vrfv2 "github.com/smartcontractkit/chainlink/integration-tests/testconfig/vrfv2"
)

type Config struct {
	Common            *Common                  `toml:"Common"`
	General           *General                 `toml:"General"`
	ExistingEnvConfig *ExistingEnvConfig       `toml:"ExistingEnvConfig"`
	NewEnvConfig      *NewEnvConfig            `toml:"NewEnvConfig"`
	Performance       *vrfv2.PerformanceConfig `toml:"Performance"`
}

func (c *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}
	if c.Common == nil && from.Common != nil {
		c.Common = from.Common
	} else if c.Common != nil && from.Common != nil {
		if err := c.Common.ApplyOverrides(from.Common); err != nil {
			return err
		}
	}
	if c.General == nil && from.General != nil {
		c.General = from.General
	} else if c.General != nil && from.General != nil {
		if err := c.General.ApplyOverrides(from.General); err != nil {
			return err
		}
	}
	if c.ExistingEnvConfig == nil && from.ExistingEnvConfig != nil {
		c.ExistingEnvConfig = from.ExistingEnvConfig
	} else if c.ExistingEnvConfig != nil && from.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.ApplyOverrides(from.ExistingEnvConfig); err != nil {
			return err
		}
	}
	if c.NewEnvConfig == nil && from.NewEnvConfig != nil {
		c.NewEnvConfig = from.NewEnvConfig
	} else if c.NewEnvConfig != nil && from.NewEnvConfig != nil {
		if err := c.NewEnvConfig.ApplyOverrides(from.NewEnvConfig); err != nil {
			return err
		}
	}
	if c.Performance == nil && from.Performance != nil {
		c.Performance = from.Performance
	} else if c.Performance != nil && from.Performance != nil {
		if err := c.Performance.ApplyOverrides(from.Performance); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Common != nil {
		if err := c.Common.Validate(); err != nil {
			return err
		}
	}
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
	if c.NewEnvConfig != nil {
		if err := c.NewEnvConfig.Validate(); err != nil {
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

type Common struct {
	*vrfv2.Common
}

func (c *Common) ApplyOverrides(from *Common) error {
	if from == nil {
		return nil
	}
	return c.Common.ApplyOverrides(from.Common)
}

func (c *Common) Validate() error {
	if c.Common == nil {
		return nil
	}
	return c.Common.Validate()
}

type General struct {
	*vrfv2.General
	SubscriptionFundingAmountNative *float64 `toml:"subscription_funding_amount_native"` // Amount of LINK to fund the subscription with default:"1"
	FulfillmentFlatFeeLinkPPM       *uint32  `toml:"fulfillment_flat_fee_link_ppm"`      // Flat fee in ppm for LINK for the VRF Coordinator config default:"500"
	FulfillmentFlatFeeNativePPM     *uint32  `toml:"fulfillment_flat_fee_native_ppm"`    // Flat fee in ppm for native currency for the VRF Coordinator config default:"500"
}

func (c *General) ApplyOverrides(from *General) error {
	if from == nil {
		return nil
	}
	if err := c.General.ApplyOverrides(from.General); err != nil {
		return err
	}
	if from.SubscriptionFundingAmountNative != nil {
		c.SubscriptionFundingAmountNative = from.SubscriptionFundingAmountNative
	}
	if from.FulfillmentFlatFeeLinkPPM != nil {
		c.FulfillmentFlatFeeLinkPPM = from.FulfillmentFlatFeeLinkPPM
	}
	if from.FulfillmentFlatFeeNativePPM != nil {
		c.FulfillmentFlatFeeNativePPM = from.FulfillmentFlatFeeNativePPM
	}

	return nil
}

func (c *General) Validate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.SubscriptionFundingAmountNative == nil || *c.SubscriptionFundingAmountNative <= 0 {
		return errors.New("subscription_funding_amount_native must be greater than 0")
	}
	if c.FulfillmentFlatFeeLinkPPM == nil || *c.FulfillmentFlatFeeLinkPPM <= 0 {
		return errors.New("fulfillment_flat_fee_link_ppm must be greater than 0")
	}
	if c.FulfillmentFlatFeeNativePPM == nil || *c.FulfillmentFlatFeeNativePPM <= 0 {
		return errors.New("fulfillment_flat_fee_native_ppm must be greater than 0")
	}

	return nil
}

type NewEnvConfig struct {
	*Funding
}

func (c *NewEnvConfig) ApplyOverrides(from *NewEnvConfig) error {
	if from == nil {
		return nil
	}
	if c.Funding == nil {
		return nil
	}
	return c.Funding.ApplyOverrides(from.Funding)
}

func (c *NewEnvConfig) Validate() error {
	if c.Funding == nil {
		return nil
	}

	return c.Funding.Validate()
}

type ExistingEnvConfig struct {
	*vrfv2.ExistingEnvConfig
	*Funding
}

func (c *ExistingEnvConfig) ApplyOverrides(from *ExistingEnvConfig) error {
	if from == nil {
		return nil
	}
	if from.ExistingEnvConfig != nil && c.ExistingEnvConfig == nil {
		c.ExistingEnvConfig = from.ExistingEnvConfig
	} else if from.ExistingEnvConfig != nil && c.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.ApplyOverrides(from.ExistingEnvConfig); err != nil {
			return err
		}
	}
	if from.Funding != nil && c.Funding == nil {
		c.Funding = from.Funding
	} else if from.Funding != nil && c.Funding != nil {
		if err := c.Funding.ApplyOverrides(from.Funding); err != nil {
			return err
		}
	}

	return nil
}

func (c *ExistingEnvConfig) Validate() error {
	if c.ExistingEnvConfig != nil {
		if err := c.ExistingEnvConfig.Validate(); err != nil {
			return err
		}
	}
	if c.Funding != nil {
		if err := c.Funding.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Funding struct {
	*SubFunding
	NodeSendingKeyFunding    *float64 `toml:"node_sending_key_funding"`
	NodeSendingKeyFundingMin *float64 `toml:"node_sending_key_funding_min"`
}

func (c *Funding) ApplyOverrides(from *Funding) error {
	if from == nil {
		return nil
	}
	if from.NodeSendingKeyFunding != nil {
		c.NodeSendingKeyFunding = from.NodeSendingKeyFunding
	}
	if from.NodeSendingKeyFundingMin != nil {
		c.NodeSendingKeyFundingMin = from.NodeSendingKeyFundingMin
	}

	return c.SubFunding.ApplyOverrides(from.SubFunding)
}

func (c *Funding) Validate() error {
	if c.NodeSendingKeyFunding == nil || *c.NodeSendingKeyFunding <= 0 {
		return errors.New("node_sending_key_funding must be greater than 0")
	}
	if c.NodeSendingKeyFundingMin == nil || *c.NodeSendingKeyFundingMin <= 0 {
		return errors.New("node_sending_key_funding_min must be greater than 0")
	}
	if err := c.SubFunding.Validate(); err != nil {
		return err
	}
	if *c.NodeSendingKeyFunding < *c.NodeSendingKeyFundingMin {
		return errors.New("node_sending_key_funding must be greater than or equal to node_sending_key_funding_min")
	}

	return nil
}

type SubFunding struct {
	SubFundsLink   *float64 `toml:"sub_funds_link"`
	SubFundsNative *float64 `toml:"sub_funds_native"`
}

func (c *SubFunding) ApplyOverrides(from *SubFunding) error {
	if from == nil {
		return nil
	}
	if from.SubFundsLink != nil {
		c.SubFundsLink = from.SubFundsLink
	}
	if from.SubFundsNative != nil {
		c.SubFundsNative = from.SubFundsNative
	}

	return nil
}

func (c *SubFunding) Validate() error {
	if c.SubFundsLink == nil || *c.SubFundsLink <= 0 {
		return errors.New("sub_funds_link must be greater than 0")
	}
	if c.SubFundsNative == nil || *c.SubFundsNative <= 0 {
		return errors.New("sub_funds_native must be greater than 0")
	}

	return nil
}
