package loadvrfv2plus

import (
	"encoding/base64"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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
	Soak              *Soak              `toml:"Soak"`
	Load              *Load              `toml:"Load"`
	SoakVolume        *SoakVolume        `toml:"SoakVolume"`
	LoadVolume        *LoadVolume        `toml:"LoadVolume"`
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
	IsNativePayment      bool   `toml:"is_native_payment"`
	MinimumConfirmations uint16 `toml:"minimum_confirmations"`
}

type Funding struct {
	NodeFunds      float64 `toml:"node_funds"`
	SubFundsLink   int64   `toml:"sub_funds_link"`
	SubFundsNative int64   `toml:"sub_funds_native"`
}

type Soak struct {
	RPS int64 `toml:"rps"`
	//Duration *models.Duration `toml:"duration"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`

	RandomnessRequestCountPerRequest          uint16 `toml:"randomness_request_count_per_request"`
	RandomnessRequestCountPerRequestDeviation uint16 `toml:"randomness_request_count_per_request_deviation"`
}

type SoakVolume struct {
	Products int64            `toml:"products"`
	Pace     *models.Duration `toml:"pace"`
	Duration *models.Duration `toml:"duration"`
}

type Load struct {
	RPSFrom     int64            `toml:"rps_from"`
	RPSIncrease int64            `toml:"rps_increase"`
	RPSSteps    int              `toml:"rps_steps"`
	Duration    *models.Duration `toml:"duration"`
}

type LoadVolume struct {
	ProductsFrom     int64            `toml:"products_from"`
	ProductsIncrease int64            `toml:"products_increase"`
	ProductsSteps    int              `toml:"products_steps"`
	Pace             *models.Duration `toml:"pace"`
	Duration         *models.Duration `toml:"duration"`
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
