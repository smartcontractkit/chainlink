package logpoller

import (
	"fmt"
	"os"
	"strconv"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
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
	RPS                   int64                  `toml:"rps"`
	LPS                   int64                  `toml:"lps"`
	RateLimitUnitDuration *commonconfig.Duration `toml:"rate_limit_unit_duration"`
	Duration              *commonconfig.Duration `toml:"duration"`
	CallTimeout           *commonconfig.Duration `toml:"call_timeout"`
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

func (c *Config) OverrideFromEnv() error {
	if contr := os.Getenv("CONTRACTS"); contr != "" {
		c.General.Contracts = mustParseInt(contr)
	}

	if eventsPerTx := os.Getenv("EVENTS_PER_TX"); eventsPerTx != "" {
		c.General.EventsPerTx = mustParseInt(eventsPerTx)
	}

	if useFinalityTag := os.Getenv("USE_FINALITY_TAG"); useFinalityTag != "" {
		c.General.UseFinalityTag = mustParseBool(useFinalityTag)
	}

	if duration := os.Getenv("LOAD_DURATION"); duration != "" {
		d, err := commonconfig.ParseDuration(duration)
		if err != nil {
			return err
		}

		if c.General.Generator == GeneratorType_WASP {
			c.Wasp.Load.Duration = &d
		} else {
			// this is completely arbitrary and practice shows that even with this values
			// test executes much longer than specified, probably due to network latency
			c.LoopedConfig.FuzzConfig.MinEmitWaitTimeMs = 400
			c.LoopedConfig.FuzzConfig.MaxEmitWaitTimeMs = 600
			// divide by 4 based on past runs, but we should do it in a better way
			c.LoopedConfig.ContractConfig.ExecutionCount = int(d.Duration().Seconds() / 4)
		}
	}

	return nil
}

func (c *Config) validate() error {
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
