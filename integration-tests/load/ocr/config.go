package ocr

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

const (
	DefaultConfigFilename  = "config.toml"
	ErrReadPerfConfig      = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig = "failed to unmarshal TOML config for performance tests"
)

type PerformanceConfig struct {
	Load   *Load   `toml:"Load"`
	Volume *Volume `toml:"Volume"`
	Common *Common `toml:"Common"`
}

type Common struct {
	ETHFunds int `toml:"eth_funds"`
}

type Load struct {
	TestDuration          *commonconfig.Duration `toml:"test_duration"`
	Rate                  int64                  `toml:"rate"`
	RateLimitUnitDuration *commonconfig.Duration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *commonconfig.Duration `toml:"verification_interval"`
	VerificationTimeout   *commonconfig.Duration `toml:"verification_timeout"`
	EAChangeInterval      *commonconfig.Duration `toml:"ea_change_interval"`
}

type Volume struct {
	TestDuration          *commonconfig.Duration `toml:"test_duration"`
	Rate                  int64                  `toml:"rate"`
	VURequestsPerUnit     int                    `toml:"vu_requests_per_unit"`
	RateLimitUnitDuration *commonconfig.Duration `toml:"rate_limit_unit_duration"`
	VerificationInterval  *commonconfig.Duration `toml:"verification_interval"`
	VerificationTimeout   *commonconfig.Duration `toml:"verification_timeout"`
	EAChangeInterval      *commonconfig.Duration `toml:"ea_change_interval"`
}

func ReadConfig() (*PerformanceConfig, error) {
	var cfg *PerformanceConfig
	rawConfig := os.Getenv("CONFIG")
	var d []byte
	var err error
	if rawConfig == "" {
		d, err = os.ReadFile(DefaultConfigFilename)
		if err != nil {
			return nil, fmt.Errorf("%s, err: %w", ErrReadPerfConfig, err)
		}
	} else {
		d, err = base64.StdEncoding.DecodeString(rawConfig)
		if err != nil {
			return nil, fmt.Errorf("%s, err: %w", ErrReadPerfConfig, err)
		}
	}
	err = toml.Unmarshal(d, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%s, err: %w", ErrUnmarshalPerfConfig, err)
	}

	log.Debug().Interface("Config", cfg).Msg("Parsed config")
	return cfg, nil
}
