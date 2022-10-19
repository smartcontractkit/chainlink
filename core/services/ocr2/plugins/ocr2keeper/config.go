package ocr2keeper

import (
	"encoding/json"
	"fmt"
	"time"
)

type PluginConfig struct {
	// CacheExpiration is the duration of time a cached key is available. Use
	// this value to balance memory usage and RPC calls. A new set of keys is
	// generated with every block so a good setting might come from block time
	// times number of blocks of history to support not replaying reports.
	CacheExpiration time.Duration `json:"cacheExpiration"`
	// CacheEvictionInterval is a parameter for how often the cache attempts to
	// evict expired keys. This value should be short enough to ensure key
	// eviction doesn't block for too long, and long enough that it doesn't
	// cause frequent blocking.
	CacheEvictionInterval time.Duration `json:"cacheEvictionInterval"`
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

type rawStruct struct {
	cacheExpiration       string
	cacheEvictionInterval string
	maxServiceWorkers     int
	serviceQueueLength    int
}

func (c *PluginConfig) UnmarshalJSON(b []byte) error {

	var raw rawStruct
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	conf := PluginConfig{
		MaxServiceWorkers:  raw.maxServiceWorkers,
		ServiceQueueLength: raw.serviceQueueLength,
	}

	d, err := time.ParseDuration(raw.cacheExpiration)
	if err != nil {
		return err
	}

	conf.CacheExpiration = d

	d, err = time.ParseDuration(raw.cacheEvictionInterval)
	if err != nil {
		return err
	}

	conf.CacheEvictionInterval = d
	*c = conf
	return nil
}

func (c PluginConfig) MarshalJSON() ([]byte, error) {
	raw := rawStruct{
		cacheExpiration:       c.CacheExpiration.String(),
		cacheEvictionInterval: c.CacheEvictionInterval.String(),
		maxServiceWorkers:     c.MaxServiceWorkers,
		serviceQueueLength:    c.ServiceQueueLength,
	}

	return json.Marshal(raw)
}

func ValidatePluginConfig(cfg PluginConfig) error {
	if cfg.CacheExpiration < 0 {
		return fmt.Errorf("cache expiration cannot be less than zero")
	}

	if cfg.CacheEvictionInterval < 0 {
		return fmt.Errorf("cache eviction interval cannot be less than zero")
	}

	if cfg.CacheEvictionInterval > 0 && cfg.CacheEvictionInterval < time.Second {
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
