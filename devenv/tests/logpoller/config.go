package logpoller

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type GeneratorType = string

const (
	GeneratorType_WASP   = "wasp"   //nolint: revive //we feel like using underscores
	GeneratorType_Looped = "looped" //nolint: revive //we feel like using underscores
)

type Config struct {
	General      General
	ChaosConfig  *ChaosConfig
	Wasp         *WaspConfig
	LoopedConfig *LoopedConfig
}

func (c *Config) Validate() error {
	err := c.General.Validate()
	if err != nil {
		return fmt.Errorf("General config validation failed: %w", err)
	}

	switch c.General.Generator {
	case GeneratorType_WASP:
		if c.Wasp == nil {
			return errors.New("wasp config is nil")
		}
		err = c.Wasp.Validate()
		if err != nil {
			return fmt.Errorf("wasp config validation failed: %w", err)
		}
	case GeneratorType_Looped:
		if c.LoopedConfig == nil {
			return errors.New("looped config is nil")
		}
		err = c.LoopedConfig.Validate()
		if err != nil {
			return fmt.Errorf("looped config validation failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown generator type: %s", c.General.Generator)
	}

	if c.ChaosConfig != nil {
		if err := c.ChaosConfig.Validate(); err != nil {
			return fmt.Errorf("chaos config validation failed: %w", err)
		}
	}

	return nil
}

type LoopedConfig struct {
	ExecutionCount    int
	MinEmitWaitTimeMs int
	MaxEmitWaitTimeMs int
}

func (l *LoopedConfig) Validate() error {
	if l.ExecutionCount == 0 {
		return errors.New("execution_count must be set and > 0")
	}

	if l.MinEmitWaitTimeMs == 0 {
		return errors.New("min_emit_wait_time_ms must be set and > 0")
	}

	if l.MaxEmitWaitTimeMs == 0 {
		return errors.New("max_emit_wait_time_ms must be set and > 0")
	}

	return nil
}

type General struct {
	Generator        string
	EventsToEmit     []abi.Event
	Contracts        int
	EventsPerTx      int
	FundingAmountEth float64
}

func (g *General) Validate() error {
	if g.Generator == "" {
		return errors.New("generator is empty")
	}

	if g.Contracts == 0 {
		return errors.New("contracts is 0, but must be > 0")
	}

	if g.EventsPerTx == 0 {
		return errors.New("events_per_tx is 0, but must be > 0")
	}

	return nil
}

type ChaosConfig struct {
	ExperimentCount int
	TargetComponent string
}

func (c *ChaosConfig) Validate() error {
	if c.ExperimentCount == 0 {
		return errors.New("experiment_count must be > 0")
	}

	return nil
}

type WaspConfig struct {
	RPS                   int64         `toml:"rps"`
	LPS                   int64         `toml:"lps"`
	RateLimitUnitDuration time.Duration `toml:"rate_limit_unit_duration"`
	Duration              time.Duration `toml:"duration"`
	CallTimeout           time.Duration `toml:"call_timeout"`
}

func (w *WaspConfig) Validate() error {
	if w.RPS == 0 && w.LPS == 0 {
		return errors.New("either RPS or LPS needs to be a positive integer")
	}
	if w.RPS != 0 && w.LPS != 0 {
		return errors.New("only one of RPS or LPS can be set")
	}
	if w.Duration == 0 {
		return errors.New("duration must be set and > 0")
	}
	if w.CallTimeout == 0 {
		return errors.New("call_timeout must be set and > 0")
	}
	if w.RateLimitUnitDuration == 0 {
		return errors.New("rate_limit_unit_duration  must be set and > 0")
	}

	return nil
}
