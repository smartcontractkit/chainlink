package ocr2keeper

import (
	"encoding/json"
	"fmt"
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
}

func ValidatePluginConfig(cfg PluginConfig) error {
	if cfg.CacheExpiration < 0 {
		return fmt.Errorf("cache expiration cannot be less than zero")
	}

	if cfg.CacheEvictionInterval < 0 {
		return fmt.Errorf("cache eviction interval cannot be less than zero")
	}

	if cfg.CacheEvictionInterval > 0 && cfg.CacheEvictionInterval.Value() < time.Second {
		return fmt.Errorf("cache eviction interval should be more than every second")
	}

	if cfg.MaxServiceWorkers < 0 {
		return fmt.Errorf("max service workers cannot be less than zero")
	}

	if cfg.ServiceQueueLength < 0 {
		return fmt.Errorf("service queue length cannot be less than zero")
	}

	return nil
}
