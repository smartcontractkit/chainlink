package logpoller

import (
	"fmt"
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

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	log.Debug().Interface("Config", cfg).Msg("Parsed config")
	return cfg, nil
}

func (c *Config) validate() error {
	if c.General == nil {
		return fmt.Errorf("General config is nil")
	}

	err := c.General.validate()
	if err != nil {
		return fmt.Errorf("General config validation failed: %v", err)
	}

	switch c.General.Generator {
	case GeneratorType_WASP:
		if c.Wasp == nil {
			return fmt.Errorf("Wasp config is nil")
		}
		if c.Wasp.Load == nil {
			return fmt.Errorf("Wasp load config is nil")
		}

		err = c.Wasp.validate()
		if err != nil {
			return fmt.Errorf("Wasp config validation failed: %v", err)
		}
	case GeneratorType_Looped:
		if c.LoopedConfig == nil {
			return fmt.Errorf("Looped config is nil")
		}

		err = c.LoopedConfig.validate()
		if err != nil {
			return fmt.Errorf("Looped config validation failed: %v", err)
		}
	default:
		return fmt.Errorf("Unknown generator type: %s", c.General.Generator)
	}

	return nil
}

func (g *General) validate() error {
	if g.Generator == "" {
		return fmt.Errorf("Generator is empty")
	}

	if g.Contracts == 0 {
		return fmt.Errorf("Contracts is 0, but must be > 0")
	}

	if g.EventsPerTx == 0 {
		return fmt.Errorf("Events_per_tx is 0, but must be > 0")
	}

	return nil
}

func (w *WaspConfig) validate() error {
	if w.Load == nil {
		return fmt.Errorf("Load config is nil")
	}

	err := w.Load.validate()
	if err != nil {
		return fmt.Errorf("Load config validation failed: %v", err)
	}

	return nil
}

func (l *Load) validate() error {
	if l.RPS == 0 {
		return fmt.Errorf("RPS is 0, but must be > 0")

	}

	if l.Duration == nil {
		return fmt.Errorf("duration is nil")
	}

	if l.CallTimeout == nil {
		return fmt.Errorf("call_timeout is nil")
	}
	if l.RateLimitUnitDuration == nil {
		return fmt.Errorf("rate_limit_unit_duration is nil")
	}

	return nil
}

func (l *LoopedConfig) validate() error {
	if l.ExecutionCount == 0 {
		return fmt.Errorf("execution_count is 0, but must be > 0")
	}

	if l.MinEmitWaitTimeMs == 0 {
		return fmt.Errorf("min_emit_wait_time_ms is 0, but must be > 0")
	}

	if l.MaxEmitWaitTimeMs == 0 {
		return fmt.Errorf("max_emit_wait_time_ms is 0, but must be > 0")
	}

	return nil
}
