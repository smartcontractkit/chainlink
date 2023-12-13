package ocr

import (
	pkg_errors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type Config struct {
	Soak   *SoakConfig `toml:"soak"`
	Load   *Load       `toml:"Load"`
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

type Load struct {
	TestDuration          *models.Duration `toml:"test_duration"`
	Rate                  int64            `toml:"rate"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *models.Duration `toml:"verification_interval"`
	VerificationTimeout   *models.Duration `toml:"verification_timeout"`
	EAChangeInterval      *models.Duration `toml:"ea_change_interval"`
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

func (o *Config) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case *Config:
		if asCfg == nil {
			return nil
		}

		if asCfg.Soak != nil && o.Soak == nil {
			o.Soak = asCfg.Soak
		}

		if asCfg.Soak != nil && o.Soak != nil {
			if err := o.Soak.ApplyOverrides(asCfg.Soak); err != nil {
				return err
			}
		}

		return nil
	default:
		return pkg_errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (o *SoakConfig) ApplyOverrides(from interface{}) error {
	switch asCfg := (from).(type) {
	case *SoakConfig:
		if asCfg == nil {
			return nil
		}

		if asCfg.TestDuration != nil {
			o.TestDuration = asCfg.TestDuration
		}

		if asCfg.NumberOfContracts != nil {
			o.NumberOfContracts = asCfg.NumberOfContracts
		}

		if asCfg.ChainlinkNodeFunding != nil {
			o.ChainlinkNodeFunding = asCfg.ChainlinkNodeFunding
		}

		if asCfg.TimeBetweenRounds != nil {
			o.TimeBetweenRounds = asCfg.TimeBetweenRounds
		}

		return nil
	default:
		return pkg_errors.Errorf("cannot apply overrides from unknown type %T", from)
	}
}

func (o *Config) Validate() error {
	return nil
}
