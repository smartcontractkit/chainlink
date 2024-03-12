package ocr

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

type Config struct {
	Soak   *SoakConfig `toml:"Soak"`
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
