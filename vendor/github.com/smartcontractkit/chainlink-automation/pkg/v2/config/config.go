package config

import (
	"encoding/json"
	"runtime"
	"time"
)

const (
	// DefaultCacheExpiration is the default amount of time a key can remain
	// in the cache before being eligible to be cleared
	DefaultCacheExpiration = 20 * time.Minute
	// DefaultCacheClearInterval is the default setting for the interval at
	// which the cache attempts to evict expired keys
	DefaultCacheClearInterval = 30 * time.Second
	// DefaultServiceQueueLength is the default buffer size for the RPC worker
	// queue.
	DefaultServiceQueueLength = 1000
)

var (
	// DefaultMaxServiceWorkers is the max number of workers allowed to make
	// simultaneous RPC calls. The default is based on the number of CPUs
	// available to the current process.
	DefaultMaxServiceWorkers = 10 * runtime.GOMAXPROCS(0)
	// each field should be validated and default values set, if any
	// if a value is invalid, return an error but don't override it with the
	// default
	validators = []validator{
		validatePerformLockoutWindow,
		validateTargetProbability,
		validateTargetInRounds,
		validateSamplingJobDuration,
		validateMinConfirmations,
		validateGasLimitPerReport,
		validateGasOverheadPerUpkeep,
		validateMaxUpkeepBatchSize,
		validateReportBlockLag,
	}
)

type ReportingFactoryConfig struct {
	CacheExpiration       time.Duration
	CacheEvictionInterval time.Duration
	MaxServiceWorkers     int
	ServiceQueueLength    int
}

// OffchainConfig ...
type OffchainConfig struct {
	// Version is used by the plugin to switch feature sets based on the intent
	// of the off-chain config
	Version string `json:"version"`

	// PerformLockoutWindow is the window in which a single upkeep cannot be
	// performed again while waiting for a confirmation. Standard setting is
	// 100 blocks * average block time. Units are in milliseconds
	PerformLockoutWindow int64 `json:"performLockoutWindow"`

	// TargetProbability is the probability that all upkeeps will be checked
	// within the provided number rounds
	TargetProbability string `json:"targetProbability"`

	// TargetInRounds is the number of rounds for the above probability to be
	// calculated
	TargetInRounds int `json:"targetInRounds"`

	// SamplingJobDuration is the time allowed for a sampling run to complete
	// before forcing a new job on the latest block. Units are in milliseconds.
	SamplingJobDuration int64 `json:"samplingJobDuration"`

	// MinConfirmations limits registered log events to only those that have
	// the provided number of confirmations.
	MinConfirmations int `json:"minConfirmations"`

	// GasLimitPerReport is the max gas that could be spent per one report.
	// This is needed for calculation of how many upkeeps could be within report.
	GasLimitPerReport uint32 `json:"gasLimitPerReport"`

	// GasOverheadPerUpkeep is gas overhead per upkeep taken place in the report.
	GasOverheadPerUpkeep uint32 `json:"gasOverheadPerUpkeep"`

	// MaxUpkeepBatchSize is the max upkeep batch size of the OCR2 report.
	MaxUpkeepBatchSize int `json:"maxUpkeepBatchSize"`

	// ReportBlockLag is the number to subtract from median block number during report phase.
	ReportBlockLag int `json:"reportBlockLag"`

	// MercuryLookup is a flag to use mercury lookup in the plugin
	MercuryLookup bool `json:"mercuryLookup"`
}

// DecodeOffchainConfig decodes bytes into an OffchainConfig
func DecodeOffchainConfig(b []byte) (OffchainConfig, error) {
	var config OffchainConfig

	// we should at minimum have a parsable config
	// if not, throw an error before validation begins
	if err := json.Unmarshal(b, &config); err != nil {
		return config, err
	}

	// go through all validators and return an error immediately if encountered
	for _, v := range validators {
		if err := v(&config); err != nil {
			return config, err
		}
	}

	return config, nil
}

type validator func(*OffchainConfig) error

func validatePerformLockoutWindow(conf *OffchainConfig) error {
	if conf.PerformLockoutWindow <= 0 {
		// default of 20 minutes (100 blocks on eth)
		conf.PerformLockoutWindow = 20 * 60 * 1000
	}

	return nil
}

func validateTargetProbability(conf *OffchainConfig) error {
	if len(conf.TargetProbability) == 0 {
		conf.TargetProbability = "0.99999"
	}

	return nil
}

func validateTargetInRounds(conf *OffchainConfig) error {
	if conf.TargetInRounds <= 0 {
		conf.TargetInRounds = 1
	}

	return nil
}

func validateSamplingJobDuration(conf *OffchainConfig) error {
	if conf.SamplingJobDuration <= 0 {
		// default of 3 seconds if not set
		conf.SamplingJobDuration = 3000
	}

	return nil
}

func validateMinConfirmations(conf *OffchainConfig) error {
	if conf.MinConfirmations <= 0 {
		conf.MinConfirmations = 0
	}

	return nil
}

func validateGasLimitPerReport(conf *OffchainConfig) error {
	// defined as uint so cannot be < 0
	if conf.GasLimitPerReport == 0 {
		conf.GasLimitPerReport = 5_300_000
	}

	return nil
}

func validateGasOverheadPerUpkeep(conf *OffchainConfig) error {
	// defined as uint so cannot be < 0
	if conf.GasOverheadPerUpkeep == 0 {
		conf.GasOverheadPerUpkeep = 300_000
	}

	return nil
}

func validateMaxUpkeepBatchSize(conf *OffchainConfig) error {
	if conf.MaxUpkeepBatchSize <= 0 {
		conf.MaxUpkeepBatchSize = 1
	}

	return nil
}

func validateReportBlockLag(conf *OffchainConfig) error {
	if conf.ReportBlockLag < 0 {
		conf.ReportBlockLag = 0
	}

	return nil
}
