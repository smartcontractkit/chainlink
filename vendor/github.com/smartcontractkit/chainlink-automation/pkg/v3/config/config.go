package config

import (
	"runtime"
	"time"

	"github.com/goccy/go-json"
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
)

type ReportingFactoryConfig struct {
	CacheExpiration       time.Duration
	CacheEvictionInterval time.Duration
	MaxServiceWorkers     int
	ServiceQueueLength    int
}

// NOTE: Any changes to this struct should keep in mind existing production contracts
// with deployed config. Additionally, offchain node upgrades can happen
// out of sync and nodes should be compatible with each other during the upgrade
// Please ensure to get a proper review along with an upgrade plan before changing this
type OffchainConfig struct {
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

	// MinConfirmations limits registered transmit events to only those that have
	// the provided number of confirmations.
	MinConfirmations int `json:"minConfirmations"`

	// GasLimitPerReport is the max gas that could be spent per one report.
	// This is needed for calculation of how many upkeeps could be within report.
	GasLimitPerReport uint32 `json:"gasLimitPerReport"`

	// GasOverheadPerUpkeep is gas overhead per upkeep taken place in the report.
	GasOverheadPerUpkeep uint32 `json:"gasOverheadPerUpkeep"`

	// MaxUpkeepBatchSize is the max upkeep batch size of the OCR2 report.
	MaxUpkeepBatchSize int `json:"maxUpkeepBatchSize"`

	// LogProviderConfig holds configuration for the log provider
	LogProviderConfig LogProviderConfig `json:"logProviderConfig"`
}

type LogProviderConfig struct {
	// BlockRate is the amount of blocks used together with LogLimitHigh to define the rate limit for each upkeep in the registry.
	BlockRate uint32 `json:"blockRate"`

	// LogLimit is the lower bound / minimum number of logs that CLA is committed to process for each upkeep per BlockRate.
	LogLimit uint32 `json:"logLimit"`
}

// DecodeOffchainConfig decodes bytes into an OffchainConfig
func DecodeOffchainConfig(b []byte) (OffchainConfig, error) {
	var config OffchainConfig

	// we should at minimum have a parsable config
	// if not, throw an error before validation begins
	if err := json.Unmarshal(b, &config); err != nil {
		return config, err
	}

	// ensure the defaults are applied at a minimum, for any values below the acceptable lower bound
	ensureMinimumDefaults(&config)

	return config, nil
}

func ensureMinimumDefaults(conf *OffchainConfig) {
	if conf.PerformLockoutWindow <= 0 {
		// default of 20 minutes (100 blocks on eth)
		conf.PerformLockoutWindow = 20 * 60 * 1000
	}
	if len(conf.TargetProbability) == 0 {
		conf.TargetProbability = "0.99999"
	}
	if conf.TargetInRounds <= 0 {
		conf.TargetInRounds = 1
	}
	if conf.MinConfirmations <= 0 {
		conf.MinConfirmations = 0
	}
	// defined as uint so cannot be < 0
	if conf.GasLimitPerReport == 0 {
		conf.GasLimitPerReport = 5_300_000
	}
	// defined as uint so cannot be < 0
	if conf.GasOverheadPerUpkeep == 0 {
		conf.GasOverheadPerUpkeep = 300_000
	}
	if conf.MaxUpkeepBatchSize <= 0 {
		conf.MaxUpkeepBatchSize = 1
	}
}
