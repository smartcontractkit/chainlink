package logpoller

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

const (
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
	ChaosConfig  *ChaosConfig  `toml:"chaos"`
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
	Generator      string      `toml:"generator"`
	EventsToEmit   []abi.Event `toml:"-"`
	Contracts      int         `toml:"contracts"`
	EventsPerTx    int         `toml:"events_per_tx"`
	UseFinalityTag bool        `toml:"use_finality_tag"`
}

type ChaosConfig struct {
	ExperimentCount int `toml:"experiment_count"`
}

type WaspConfig struct {
	Load *Load `toml:"load"`
}

type Load struct {
	RPS                   int64            `toml:"rps"`
	LPS                   int64            `toml:"lps"`
	RateLimitUnitDuration *models.Duration `toml:"rate_limit_unit_duration"`
	Duration              *models.Duration `toml:"duration"`
	CallTimeout           *models.Duration `toml:"call_timeout"`
}

func (c *Config) ApplyOverrides(_ *Config) error {
	//TODO implement me
	return nil
}

func (c *Config) Validate() error {
	if c.General == nil {
		return fmt.Errorf("General config is nil")
	}

	err := c.General.validate()
	if err != nil {
		return fmt.Errorf("General config validation failed: %w", err)
	}

	switch c.General.Generator {
	case GeneratorType_WASP:
		if c.Wasp == nil {
			return fmt.Errorf("wasp config is nil")
		}
		if c.Wasp.Load == nil {
			return fmt.Errorf("wasp load config is nil")
		}

		err = c.Wasp.validate()
		if err != nil {
			return fmt.Errorf("wasp config validation failed: %w", err)
		}
	case GeneratorType_Looped:
		if c.LoopedConfig == nil {
			return fmt.Errorf("looped config is nil")
		}

		err = c.LoopedConfig.validate()
		if err != nil {
			return fmt.Errorf("looped config validation failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown generator type: %s", c.General.Generator)
	}

	return nil
}

func (g *General) validate() error {
	if g.Generator == "" {
		return fmt.Errorf("generator is empty")
	}

	if g.Contracts == 0 {
		return fmt.Errorf("contracts is 0, but must be > 0")
	}

	if g.EventsPerTx == 0 {
		return fmt.Errorf("events_per_tx is 0, but must be > 0")
	}

	return nil
}

func (w *WaspConfig) validate() error {
	if w.Load == nil {
		return fmt.Errorf("Load config is nil")
	}

	err := w.Load.validate()
	if err != nil {
		return fmt.Errorf("Load config validation failed: %w", err)
	}

	return nil
}

func (l *Load) validate() error {
	if l.RPS == 0 && l.LPS == 0 {
		return fmt.Errorf("either RPS or LPS needs to be set")
	}

	if l.RPS != 0 && l.LPS != 0 {
		return fmt.Errorf("only one of RPS or LPS can be set")
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

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func mustParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return b
}
