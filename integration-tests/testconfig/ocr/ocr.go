package ocr

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

type Config struct {
	Soak   *SoakConfig `toml:"Soak"`
	Load   *Load       `toml:"Load"`
	Volume *Volume     `toml:"Volume"`
	Common *Common     `toml:"Common"`
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
	ETHFunds     *int                    `toml:"eth_funds"`
	TestDuration *blockchain.StrDuration `toml:"test_duration"`
}

func (o *Common) Validate() error {
	if o.ETHFunds != nil && *o.ETHFunds < 0 {
		return errors.New("eth_funds must be set and cannot be negative")
	}
	return nil
}

type Load struct {
	Rate                  *int64                  `toml:"rate"`
	RequestsPerUnit       *int                    `toml:"requests_per_unit"`
	RateLimitUnitDuration *blockchain.StrDuration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *blockchain.StrDuration `toml:"verification_interval"`
	VerificationTimeout   *blockchain.StrDuration `toml:"verification_timeout"`
	EAChangeInterval      *blockchain.StrDuration `toml:"ea_change_interval"`
	TestDuration          *blockchain.StrDuration `toml:"test_duration"`
}

func (o *Load) Validate() error {
	if o.TestDuration == nil {
		return errors.New("load test duration must be set")
	}
	if o.Rate == nil || *o.Rate <= 0 {
		return errors.New("rate must be set and be a positive integer")
	}
	if o.RequestsPerUnit == nil || *o.RequestsPerUnit <= 0 {
		return errors.New("vu_requests_per_unit must be set and be a positive integer")
	}
	if o.RateLimitUnitDuration == nil || o.RateLimitUnitDuration.Duration == 0 {
		return errors.New("rate_limit_unit_duration must be set and be a positive integer")
	}
	if o.VerificationInterval == nil || o.VerificationInterval.Duration == 0 {
		return errors.New("verification_interval must be set and be a positive integer")
	}
	if o.VerificationTimeout == nil || o.VerificationTimeout.Duration == 0 {
		return errors.New("verification_timeout must be set and be a positive integer")
	}
	if o.EAChangeInterval == nil || o.EAChangeInterval.Duration == 0 {
		return errors.New("ea_change_interval must be set and be a positive integer")
	}

	return nil
}

type Volume struct {
	Rate                  *int64                  `toml:"rate"`
	VURequestsPerUnit     *int                    `toml:"vu_requests_per_unit"`
	RateLimitUnitDuration *blockchain.StrDuration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *blockchain.StrDuration `toml:"verification_interval"`
	VerificationTimeout   *blockchain.StrDuration `toml:"verification_timeout"`
	EAChangeInterval      *blockchain.StrDuration `toml:"ea_change_interval"`
	TestDuration          *blockchain.StrDuration `toml:"test_duration"`
}

func (o *Volume) Validate() error {
	if o.TestDuration == nil {
		return errors.New("volume test duration must be set")
	}
	if o.Rate == nil || *o.Rate <= 0 {
		return errors.New("rate must be set and be a positive integer")
	}
	if o.VURequestsPerUnit == nil || *o.VURequestsPerUnit <= 0 {
		return errors.New("vu_requests_per_unit must be set and be a positive integer")
	}
	if o.RateLimitUnitDuration == nil || o.RateLimitUnitDuration.Duration == 0 {
		return errors.New("rate_limit_unit_duration must be set and be a positive integer")
	}
	if o.VerificationInterval == nil || o.VerificationInterval.Duration == 0 {
		return errors.New("verification_interval must be set and be a positive integer")
	}
	if o.VerificationTimeout == nil || o.VerificationTimeout.Duration == 0 {
		return errors.New("verification_timeout must be set and be a positive integer")
	}
	if o.EAChangeInterval == nil || o.EAChangeInterval.Duration == 0 {
		return errors.New("ea_change_interval must be set and be a positive integer")
	}

	return nil
}

type SoakConfig struct {
	OCRVersion        *string                 `toml:"ocr_version"`
	NumberOfContracts *int                    `toml:"number_of_contracts"`
	TimeBetweenRounds *blockchain.StrDuration `toml:"time_between_rounds"`
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
