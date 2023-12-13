package ocr

import (
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config struct {
	Soak   *SoakConfig `toml:"soak"`
	Load   *Volume     `toml:"Load"`
	Volume *Volume     `toml:"Volume"`
	Common *Common     `toml:"Common"`
}

// TODO extract common fields from SoakConfig and VolumeConfig and Load to Common
// this came from env config!
type SoakConfig struct {
	TestDuration         *blockchain.JSONStrDuration `toml:"test_duration"`          //default:"15m
	NumberOfContracts    *int                        `toml:"number_of_contracts"`    //default:"2"
	ChainlinkNodeFunding *float64                    `toml:"chainlink_node_funding"` //default:".1"
	TimeBetweenRounds    *blockchain.JSONStrDuration `toml:"time_between_rounds"`    //default:"1m"
}

type Common struct {
	ETHFunds int `toml:"eth_funds"`
}

type Volume struct {
	TestDuration          *models.Duration `toml:"test_duration"`
	Rate                  int64            `toml:"rate"`
	VURequestsPerUnit     int              `toml:"vu_requests_per_unit"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *models.Duration `toml:"verification_interval"`
	VerificationTimeout   *models.Duration `toml:"verification_timeout"`
	EAChangeInterval      *models.Duration `toml:"ea_change_interval"`
}

func (o *Config) ApplyOverrides(from *Config) error {
	if from == nil {
		return nil
	}

	if from.Soak != nil && o.Soak == nil {
		o.Soak = from.Soak
	} else if from.Soak != nil && o.Soak != nil {
		if err := o.Soak.ApplyOverrides(from.Soak); err != nil {
			return err
		}
	}

	if from.Load != nil && o.Load == nil {
		o.Load = from.Load
	} else if from.Load != nil && o.Load != nil {
		if err := o.Load.ApplyOverrides(from.Load); err != nil {
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

func (o *SoakConfig) ApplyOverrides(from *SoakConfig) error {
	if from == nil {
		return nil
	}

	if from.TestDuration != nil {
		o.TestDuration = from.TestDuration
	}

	if from.NumberOfContracts != nil {
		o.NumberOfContracts = from.NumberOfContracts
	}

	if from.ChainlinkNodeFunding != nil {
		o.ChainlinkNodeFunding = from.ChainlinkNodeFunding
	}

	if from.TimeBetweenRounds != nil {
		o.TimeBetweenRounds = from.TimeBetweenRounds
	}

	return nil
}

func (o *Volume) ApplyOverrides(from *Volume) error {
	if from == nil {
		return nil
	}
	if from.TestDuration != nil {
		o.TestDuration = from.TestDuration
	}
	if from.Rate != 0 {
		o.Rate = from.Rate
	}
	if from.VURequestsPerUnit != 0 {
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

func (o *Config) Validate() error {
	return nil
}
