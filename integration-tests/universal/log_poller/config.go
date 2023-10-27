package logpoller

import (
	"os"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const (
	DefaultConfigFilename = "config.toml"

	ErrReadPerfConfig      = "failed to read TOML config for performance tests"
	ErrUnmarshalPerfConfig = "failed to unmarshal TOML config for performance tests"
)

type GeneratorType = string

const (
	GeneratorType_WASP   = "wasp"
	GeneratorType_Looped = "looped"
)

type Config struct {
	General      *General      `toml:"general"`
	Wasp         *WaspConfig   `toml:"wasp"`
	LoopedConfig *LoopedConfig `toml:"looped"`
}

type LoopedConfig struct {
	ContractConfig `toml:"contract"`
	FuzzConfig     `toml:"fuzz"`
}

type ContractConfig struct {
	ExecutionCount int `toml:"execution_count"`
}

type FuzzConfig struct {
	MinEmitWaitTimeMs int `toml:"min_emit_wait_time_ms"`
	MaxEmitWaitTimeMs int `toml:"max_emit_wait_time_ms"`
}

type General struct {
	Generator    string      `toml:"generator"`
	EventsToEmit []abi.Event `toml:"-"`
	Contracts    int         `toml:"contracts"`
	EventsPerTx  int         `toml:"events_per_tx"`
}

type WaspConfig struct {
	Load *Load `toml:"load"`
}

type Load struct {
	RPS                   int64            `toml:"rps"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`
	Duration              *models.Duration `toml:"duration"`
	CallTimeout           *models.Duration `toml:"call_timeout"`
}

func ReadConfig(configName string) (*Config, error) {
	var cfg *Config
	d, err := os.ReadFile(configName)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadPerfConfig)
	}
	err = toml.Unmarshal(d, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalPerfConfig)
	}

	switch cfg.General.Generator {
	case GeneratorType_WASP:
	case GeneratorType_Looped:
	default:
		panic("Unknown generator type")
	}

	log.Debug().Interface("Config", cfg).Msg("Parsed config")
	return cfg, nil
}
