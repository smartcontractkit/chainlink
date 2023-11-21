package loadvrfv2plus

import (
	"encoding/base64"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"os"
)

const (
	DefaultConfigFilename = "config.toml"

	ErrReadPerfConfig                    = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig               = "failed to unmarshal TOML config for performance tests"
	ErrDeviationShouldBeLessThanOriginal = "`RandomnessRequestCountPerRequestDeviation` should be less than `RandomnessRequestCountPerRequest`"
)

type PerformanceConfig struct {
	Soak   *Soak   `toml:"Soak"`
	Load   *Load   `toml:"Load"`
	Stress *Stress `toml:"Stress"`
	Spike  *Spike  `toml:"Spike"`

	Common            *Common            `toml:"Common"`
	ExistingEnvConfig *ExistingEnvConfig `toml:"ExistingEnvConfig"`
	NewEnvConfig      *NewEnvConfig      `toml:"NewEnvConfig"`
}

type ExistingEnvConfig struct {
	CoordinatorAddress string `toml:"coordinator_address"`
	ConsumerAddress    string `toml:"consumer_address"`
	SubID              string `toml:"sub_id"`
	KeyHash            string `toml:"key_hash"`
}

type NewEnvConfig struct {
	Funding
	NumberOfSubToCreate int `toml:"number_of_sub_to_create"`
}

type Common struct {
	MinimumConfirmations uint16 `toml:"minimum_confirmations"`
}

type Funding struct {
	NodeFunds      float64 `toml:"node_funds"`
	SubFundsLink   int64   `toml:"sub_funds_link"`
	SubFundsNative int64   `toml:"sub_funds_native"`
}

type Soak struct {
	PerformanceTestConfig
}

type Load struct {
	PerformanceTestConfig
}

type Stress struct {
	PerformanceTestConfig
}

type Spike struct {
	PerformanceTestConfig
}

type PerformanceTestConfig struct {
	RPS int64 `toml:"rps"`
	//Duration *models.Duration `toml:"duration"`
	RateLimitUnitDuration                     *models.Duration `toml:"rate_limit_unit_duration"`
	RandomnessRequestCountPerRequest          uint16           `toml:"randomness_request_count_per_request"`
	RandomnessRequestCountPerRequestDeviation uint16           `toml:"randomness_request_count_per_request_deviation"`
}

func ReadConfig() (*PerformanceConfig, error) {
	var cfg *PerformanceConfig
	rawConfig := os.Getenv("CONFIG")
	var d []byte
	var err error
	if rawConfig == "" {
		d, err = os.ReadFile(DefaultConfigFilename)
		if err != nil {
			return nil, errors.Wrap(err, ErrReadPerfConfig)
		}
	} else {
		d, err = base64.StdEncoding.DecodeString(rawConfig)
	}
	err = toml.Unmarshal(d, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalPerfConfig)
	}

	if cfg.Soak.RandomnessRequestCountPerRequest <= cfg.Soak.RandomnessRequestCountPerRequestDeviation {
		return nil, errors.Wrap(err, ErrDeviationShouldBeLessThanOriginal)
	}

	log.Debug().Interface("Config", cfg).Msg("Parsed config")
	return cfg, nil
}

func SetPerformanceTestConfig(vrfv2PlusConfig *vrfv2plus_config.VRFV2PlusConfig, cfg *PerformanceConfig) {
	switch os.Getenv("TEST_TYPE") {
	case "Soak":
		vrfv2PlusConfig.RPS = cfg.Soak.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Soak.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Soak.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Soak.RandomnessRequestCountPerRequestDeviation
	case "Load":
		vrfv2PlusConfig.RPS = cfg.Load.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Load.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Load.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Load.RandomnessRequestCountPerRequestDeviation
	case "Stress":
		vrfv2PlusConfig.RPS = cfg.Stress.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Stress.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Stress.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Stress.RandomnessRequestCountPerRequestDeviation
	case "Spike":
		vrfv2PlusConfig.RPS = cfg.Spike.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Spike.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Spike.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Spike.RandomnessRequestCountPerRequestDeviation
	}
}
