package loadvrfv2plus

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"
)

const (
	DefaultConfigFilename = "config.toml"
	SoakTestType          = "Soak"
	LoadTestType          = "Load"
	StressTestType        = "Stress"
	SpikeTestType         = "Spike"

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
	LinkAddress        string `toml:"link_address"`
	SubID              string `toml:"sub_id"`
	KeyHash            string `toml:"key_hash"`
	Funding
	CreateFundSubsAndAddConsumers bool     `toml:"create_fund_subs_and_add_consumers"`
	NodeSendingKeys               []string `toml:"node_sending_keys"`
}

type NewEnvConfig struct {
	Funding
}

type Common struct {
	MinimumConfirmations   uint16 `toml:"minimum_confirmations"`
	CancelSubsAfterTestRun bool   `toml:"cancel_subs_after_test_run"`
}

type Funding struct {
	SubFunding
	NodeSendingKeyFunding    float64 `toml:"node_sending_key_funding"`
	NodeSendingKeyFundingMin float64 `toml:"node_sending_key_funding_min"`
}

type SubFunding struct {
	SubFundsLink   float64 `toml:"sub_funds_link"`
	SubFundsNative float64 `toml:"sub_funds_native"`
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
	NumberOfSubToCreate int `toml:"number_of_sub_to_create"`

	RPS int64 `toml:"rps"`
	//Duration *commonconfig.Duration `toml:"duration"`
	RateLimitUnitDuration                     *commonconfig.Duration `toml:"rate_limit_unit_duration"`
	RandomnessRequestCountPerRequest          uint16                 `toml:"randomness_request_count_per_request"`
	RandomnessRequestCountPerRequestDeviation uint16                 `toml:"randomness_request_count_per_request_deviation"`
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

	if cfg.Soak.RandomnessRequestCountPerRequest <= cfg.Soak.RandomnessRequestCountPerRequestDeviation {
		return nil, fmt.Errorf("%s, err: %w", ErrDeviationShouldBeLessThanOriginal, err)
	}

	log.Debug().Interface("Config", cfg).Msg("Parsed config")
	return cfg, nil
}

func SetPerformanceTestConfig(testType string, vrfv2PlusConfig *vrfv2plus_config.VRFV2PlusConfig, cfg *PerformanceConfig) {
	switch testType {
	case SoakTestType:
		vrfv2PlusConfig.NumberOfSubToCreate = cfg.Soak.NumberOfSubToCreate
		vrfv2PlusConfig.RPS = cfg.Soak.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Soak.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Soak.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Soak.RandomnessRequestCountPerRequestDeviation
	case LoadTestType:
		vrfv2PlusConfig.NumberOfSubToCreate = cfg.Load.NumberOfSubToCreate
		vrfv2PlusConfig.RPS = cfg.Load.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Load.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Load.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Load.RandomnessRequestCountPerRequestDeviation
	case StressTestType:
		vrfv2PlusConfig.NumberOfSubToCreate = cfg.Stress.NumberOfSubToCreate
		vrfv2PlusConfig.RPS = cfg.Stress.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Stress.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Stress.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Stress.RandomnessRequestCountPerRequestDeviation
	case SpikeTestType:
		vrfv2PlusConfig.NumberOfSubToCreate = cfg.Spike.NumberOfSubToCreate
		vrfv2PlusConfig.RPS = cfg.Spike.RPS
		vrfv2PlusConfig.RateLimitUnitDuration = cfg.Spike.RateLimitUnitDuration.Duration()
		vrfv2PlusConfig.RandomnessRequestCountPerRequest = cfg.Spike.RandomnessRequestCountPerRequest
		vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation = cfg.Spike.RandomnessRequestCountPerRequestDeviation
	}
}
