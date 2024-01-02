package ocr

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config struct {
	Soak   *SoakConfig `toml:"Soak"`
	Volume *Volume     `toml:"Volume"`
	Common *Common     `toml:"Common"`
}

func (o *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}

	if from.Common != nil && o.Common == nil {
		o.Common = from.Common
	} else if err := o.Common.ApplyOverrides(from.Common); err != nil {
		return err
	}

	if from.Soak != nil && o.Soak == nil {
		o.Soak = from.Soak
	} else if from.Soak != nil && o.Soak != nil {
		if err := o.Soak.ApplyOverrides(from.Soak); err != nil {
			return err
		}
	}

	if from.Volume != nil && o.Volume == nil {
		o.Volume = from.Volume
	} else if from.Volume != nil && o.Volume != nil {
		if err := o.Volume.ApplyOverrides(from.Volume); err != nil {
			return err
		}
	}

	return nil
}

func (o *Config) Validate() error {
	if o.Common != nil {
		if err := o.Common.Validate(); err != nil {
			return err
		}
	}
	if o.Soak != nil {
		if err := o.Soak.Validate(); err != nil {
			return err
		}
	}
	if o.Volume != nil {
		if err := o.Volume.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Common struct {
	ETHFunds     *int             `toml:"eth_funds"`
	TestDuration *models.Duration `toml:"test_duration"` //default:"15m
}

func (o *Common) ApplyOverrides(from *Common) error {
	if from == nil {
		return nil
	}
	if from.ETHFunds != nil {
		o.ETHFunds = from.ETHFunds
	}
	if from.TestDuration != nil {
		o.TestDuration = from.TestDuration
	}
	return nil
}

func (o *Common) Validate() error {
	if o.ETHFunds != nil && *o.ETHFunds < 0 {
		return errors.New("eth_funds must be set and cannot be negative")
	}
	if o.TestDuration == nil || o.TestDuration.Duration() == 0 {
		return errors.New("test_duration must be set and be a positive integer")
	}

	return nil
}

type Volume struct {
	Rate                  *int64           `toml:"rate"`
	VURequestsPerUnit     *int             `toml:"vu_requests_per_unit"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *models.Duration `toml:"verification_interval"`
	VerificationTimeout   *models.Duration `toml:"verification_timeout"`
	EAChangeInterval      *models.Duration `toml:"ea_change_interval"`
}

func (o *Volume) ApplyOverrides(from *Volume) error {
	if from == nil {
		return nil
	}
	if from.Rate != nil {
		o.Rate = from.Rate
	}
	if from.VURequestsPerUnit != nil {
		o.VURequestsPerUnit = from.VURequestsPerUnit
	}
	if from.RateLimitUnitDuration != nil {
		o.RateLimitUnitDuration = from.RateLimitUnitDuration
	}
	if from.VerificationInterval != nil {
		o.VerificationInterval = from.VerificationInterval
	}
	if from.VerificationTimeout != nil {
		o.VerificationTimeout = from.VerificationTimeout
	}
	if from.EAChangeInterval != nil {
		o.EAChangeInterval = from.EAChangeInterval
	}

	return nil
}

func (o *Volume) Validate() error {
	if o.Rate == nil || *o.Rate <= 0 {
		return errors.New("rate must be set and be a positive integer")
	}
	if o.VURequestsPerUnit == nil || *o.VURequestsPerUnit <= 0 {
		return errors.New("vu_requests_per_unit must be set and be a positive integer")
	}
	if o.RateLimitUnitDuration == nil || o.RateLimitUnitDuration.Duration() == 0 {
		return errors.New("rate_limit_unit_duration must be set and be a positive integer")
	}
	if o.VerificationInterval == nil || o.VerificationInterval.Duration() == 0 {
		return errors.New("verification_interval must be set and be a positive integer")
	}
	if o.VerificationTimeout == nil || o.VerificationTimeout.Duration() == 0 {
		return errors.New("verification_timeout must be set and be a positive integer")
	}
	if o.EAChangeInterval == nil || o.EAChangeInterval.Duration() == 0 {
		return errors.New("ea_change_interval must be set and be a positive integer")
	}

	return nil
}

type SoakConfig struct {
	OCRVersion        *string                 `toml:"ocr_version"`
	NumberOfContracts *int                    `toml:"number_of_contracts"`
	TimeBetweenRounds *blockchain.StrDuration `toml:"time_between_rounds"`
}

func (o *SoakConfig) ApplyOverrides(from *SoakConfig) error {
	if from == nil {
		return nil
	}

	if from.OCRVersion != nil {
		o.OCRVersion = from.OCRVersion
	}

	if from.NumberOfContracts != nil {
		o.NumberOfContracts = from.NumberOfContracts
	}

	if from.TimeBetweenRounds != nil {
		o.TimeBetweenRounds = from.TimeBetweenRounds
	}

	return nil
}

func (o *SoakConfig) Validate() error {
	if o.OCRVersion == nil || *o.OCRVersion == "" {
		return errors.New("ocr_version must be set to either 1 or 2")
	}
	if o.NumberOfContracts == nil || *o.NumberOfContracts <= 1 {
		return errors.New("number_of_contracts must be set and be greater than 1")
	}
	if o.TimeBetweenRounds == nil || o.TimeBetweenRounds.Duration == 0 {
		return errors.New("time_between_rounds must be set and be a positive integer")
	}
	return nil
}
