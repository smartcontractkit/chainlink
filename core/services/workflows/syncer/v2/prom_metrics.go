package v2

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promModuleCacheReloadTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "platform_workflow_module_cache_reload_total",
		Help: "Module cache reloads by source (weak_ref or disk)",
	}, []string{"source"})
	promModuleCacheEvictionTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "platform_workflow_module_cache_eviction_total",
		Help: "Total module cache evictions",
	})
	promModuleCacheLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform_workflow_module_cache_loaded",
		Help: "Currently loaded modules in the LRU cache",
	})
	promModuleCacheMemorySavedBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform_workflow_module_cache_memory_saved_bytes",
		Help: "Bytes of memory saved by evicting idle modules",
	})
	promModuleCacheVersionMismatchTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "platform_workflow_module_cache_version_mismatch_total",
		Help: "Cached binaries rejected due to engine version mismatch",
	})
	promModuleCachePinExhaustedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "platform_workflow_module_cache_pin_exhausted_total",
		Help: "Execute retries exhausted before a pin succeeds",
	})
	promModuleCacheTryAcquireExhaustedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "platform_workflow_module_cache_try_acquire_exhausted_total",
		Help: "moduleEntry CAS attempts exhausted while pinning",
	})
	promModuleCacheDiskUsageBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "platform_workflow_module_cache_disk_usage_bytes",
		Help: "Total on-disk bytes under the module cache directory",
	})
)
