package ocr2

import (
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ocr"
)

type Config struct {
	Soak      *ocr.SoakConfig `toml:"Soak"`
	Common    *ocr.Common     `toml:"Common"`
	Contracts *ocr.Contracts  `toml:"Contracts"`
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
	if o.Contracts != nil {
		if err := o.Contracts.Validate(); err != nil {
			return err
		}
	}
	return nil
}
