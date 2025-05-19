package logpoller

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
)

type GeneratorType = string

const (
	GeneratorType_WASP   = "wasp"
	GeneratorType_Looped = "looped"
)

type Config struct {
	General      *General      `toml:"General"`
	ChaosConfig  *ChaosConfig  `toml:"Chaos"`
	Wasp         *WaspConfig   `toml:"Wasp"`
	LoopedConfig *LoopedConfig `toml:"Looped"`
}

func (c *Config) Validate() error {
	if c.General == nil {
		return fmt.Errorf("General config must be set")
	}

	err := c.General.Validate()
	if err != nil {
		return fmt.Errorf("General config validation failed: %w", err)
	}

	switch *c.General.Generator {
	case GeneratorType_WASP:
		if c.Wasp == nil {
			return fmt.Errorf("wasp config is nil")
		}
		err = c.Wasp.Validate()
		if err != nil {
			return fmt.Errorf("wasp config validation failed: %w", err)
		}
	case GeneratorType_Looped:
		if c.LoopedConfig == nil {
			return fmt.Errorf("looped config is nil")
		}
		err = c.LoopedConfig.Validate()
		if err != nil {
			return fmt.Errorf("looped config validation failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown generator type: %s", *c.General.Generator)
	}

	if c.ChaosConfig != nil {
		if err := c.ChaosConfig.Validate(); err != nil {
			return fmt.Errorf("chaos config validation failed: %w", err)
		}
	}

	return nil
}

type LoopedConfig struct {
	ExecutionCount    *int `toml:"execution_count"`
	MinEmitWaitTimeMs *int `toml:"min_emit_wait_time_ms"`
	MaxEmitWaitTimeMs *int `toml:"max_emit_wait_time_ms"`
}

func (l *LoopedConfig) Validate() error {
	if l.ExecutionCount == nil || *l.ExecutionCount == 0 {
		return fmt.Errorf("execution_count must be set and > 0")
	}

	if l.MinEmitWaitTimeMs == nil || *l.MinEmitWaitTimeMs == 0 {
		return fmt.Errorf("min_emit_wait_time_ms must be set and > 0")
	}

	if l.MaxEmitWaitTimeMs == nil || *l.MaxEmitWaitTimeMs == 0 {
		return fmt.Errorf("max_emit_wait_time_ms must be set and > 0")
	}

	return nil
}

type General struct {
	Generator      *string     `toml:"generator"`
	EventsToEmit   []abi.Event `toml:"-"`
	Contracts      *int        `toml:"contracts"`
	EventsPerTx    *int        `toml:"events_per_tx"`
	UseFinalityTag *bool       `toml:"use_finality_tag"`
}

func (g *General) Validate() error {
	if g.Generator == nil || *g.Generator == "" {
		return fmt.Errorf("generator is empty")
	}

	if g.Contracts == nil || *g.Contracts == 0 {
		return fmt.Errorf("contracts is 0, but must be > 0")
	}

	if g.EventsPerTx == nil || *g.EventsPerTx == 0 {
		return fmt.Errorf("events_per_tx is 0, but must be > 0")
	}

	return nil
}

type ChaosConfig struct {
	ExperimentCount *int    `toml:"experiment_count"`
	TargetComponent *string `toml:"target_component"`
}

func (c *ChaosConfig) Validate() error {
	if c.ExperimentCount != nil && *c.ExperimentCount == 0 {
		return fmt.Errorf("experiment_count must be > 0")
	}

	return nil
}

type WaspConfig struct {
	RPS                   *int64                  `toml:"rps"`
	LPS                   *int64                  `toml:"lps"`
	RateLimitUnitDuration *blockchain.StrDuration `toml:"rate_limit_unit_duration"`
	Duration              *blockchain.StrDuration `toml:"duration"`
	CallTimeout           *blockchain.StrDuration `toml:"call_timeout"`
}

func (w *WaspConfig) Validate() error {
	if w.RPS == nil && w.LPS == nil {
		return fmt.Errorf("either RPS or LPS needs to be set")
	}
	if *w.RPS == 0 && *w.LPS == 0 {
		return fmt.Errorf("either RPS or LPS needs to be a positive integer")
	}
	if *w.RPS != 0 && *w.LPS != 0 {
		return fmt.Errorf("only one of RPS or LPS can be set")
	}
	if w.Duration == nil || w.Duration.Duration == 0 {
		return fmt.Errorf("duration must be set and > 0")
	}
	if w.CallTimeout == nil || w.CallTimeout.Duration == 0 {
		return fmt.Errorf("call_timeout must be set and > 0")
	}
	if w.RateLimitUnitDuration == nil || w.RateLimitUnitDuration.Duration == 0 {
		return fmt.Errorf("rate_limit_unit_duration  must be set and > 0")
	}

	return nil
}
