package loadvrfv2

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"math/big"
	"os"
)

const (
	DefaultConfigFilename = "config.toml"

	ErrReadPerfConfig      = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig = "failed to unmarshal TOML config for performance tests"
)

type PerformanceConfig struct {
	Soak       *Soak       `toml:"Soak"`
	Load       *Load       `toml:"Load"`
	SoakVolume *SoakVolume `toml:"SoakVolume"`
	LoadVolume *LoadVolume `toml:"LoadVolume"`
	Common     *Common     `toml:"Common"`
}

type Common struct {
	Funding
}

type Funding struct {
	NodeFunds *big.Float `toml:"node_funds"`
	SubFunds  *big.Int   `toml:"sub_funds"`
}

type Soak struct {
	RPS      int64            `toml:"rps"`
	Duration *models.Duration `toml:"duration"`
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
	d, err := os.ReadFile(DefaultConfigFilename)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadPerfConfig)
	}
	err = toml.Unmarshal(d, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalPerfConfig)
	}
	log.Debug().Interface("PerformanceConfig", cfg).Msg("Parsed performance config")
	return cfg, nil
}
