package ocr2keeper

import (
	"encoding/json"
	"errors"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	p, err := time.ParseDuration(raw)
	if err != nil {
		return err
	}

	*d = Duration(p)
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

// NOTE: This plugin config is shared among different versions of keepers
// Any changes to this config should keep in mind existing production
// deployments of all versions of keepers and should be backwards compatible
// with existing job specs.
type PluginConfig struct {
	// CacheExpiration is the duration of time a cached key is available. Use
	// this value to balance memory usage and RPC calls. A new set of keys is
	// generated with every block so a good setting might come from block time
	// times number of blocks of history to support not replaying reports.
	CacheExpiration Duration `json:"cacheExpiration"`
	// CacheEvictionInterval is a parameter for how often the cache attempts to
	// evict expired keys. This value should be short enough to ensure key
	// eviction doesn't block for too long, and long enough that it doesn't
	// cause frequent blocking.
	CacheEvictionInterval Duration `json:"cacheEvictionInterval"`
	// MaxServiceWorkers is the total number of go-routines allowed to make RPC
	// simultaneous calls on behalf of the sampling operation. This parameter
	// is 10x the number of available CPUs by default. The RPC calls are memory
	// heavy as opposed to CPU heavy as most of the work involves waiting on
	// network responses.
	MaxServiceWorkers int `json:"maxServiceWorkers"`
	// ServiceQueueLength is the buffer size for the RPC service queue. Fewer
	// workers or slower RPC responses will cause this queue to build up.
	// Adding new items to the queue will block if the queue becomes full.
	ServiceQueueLength int `json:"serviceQueueLength"`
	// ContractVersion is the contract version
	ContractVersion string `json:"contractVersion"`
	// CaptureAutomationCustomTelemetry is a bool flag to toggle Custom Telemetry Service
	CaptureAutomationCustomTelemetry *bool `json:"captureAutomationCustomTelemetry,omitempty"`
}

func ValidatePluginConfig(cfg PluginConfig) error {
	if cfg.CacheExpiration < 0 {
		return errors.New("cache expiration cannot be less than zero")
	}

	if cfg.CacheEvictionInterval < 0 {
		return errors.New("cache eviction interval cannot be less than zero")
	}

	if cfg.CacheEvictionInterval > 0 && cfg.CacheEvictionInterval.Value() < time.Second {
		return errors.New("cache eviction interval should be more than every second")
	}

	if cfg.MaxServiceWorkers < 0 {
		return errors.New("max service workers cannot be less than zero")
	}

	if cfg.ServiceQueueLength < 0 {
		return errors.New("service queue length cannot be less than zero")
	}

	return nil
}
