package cre

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

// Memory sizing targets ~50% of a 128GiB soak host for concurrently loaded WASM modules across the
// 4-node workflow DON (see workflow-gateway-don-cache-soak-test.toml nodes + MaxLoaded).
// moduleCacheMaxLoaded is per node; _defaultSoakNumWorkflows exceeds cap slots (~20%)
// so enforceCap stays active after workflows are being triggered across nodes.
// On-chain registry limits (SetDONLimit) and node [Workflows.Limits] must exceed _defaultSoakNumWorkflows.
const (
	capPressurePercent   = 200 // 200% of MaxLoaded
	moduleCacheMaxLoaded = 100 // mirrors workflow-gateway-don-cache-soak-test.toml MaxLoaded

	moduleCacheIdleTimeout = 5 * time.Minute
	fastCronInterval       = 3 * time.Minute
	slowCronInterval       = 8 * time.Minute

	defaultSoakDuration     = 4 * time.Hour
	defaultMetricStep       = 1 * time.Minute
	cachePrometheusRange    = "5m" // increase() window; align with defaultMetricStep
	soakProgressLogInterval = 5 * time.Minute

	// One of every cacheSoakSchedulePeriod workflows uses slowCronInterval (~1/3 idle-eviction tier).
	cacheSoakSchedulePeriod = 3

	numberOfDeploymentKeys = 20
)

var (
	_workflowModuleMiB = crePerWorkflowSizeLimitMiB(
		cresettings.Default.PerWorkflow.WASMBinarySizeLimit.DefaultValue,
	)
	_workflowEngineOverheadMiB = crePerWorkflowSizeLimitMiB(
		cresettings.Default.PerWorkflow.WASMMemoryLimit.DefaultValue,
	)
	_defaultSoakNumWorkflows = moduleCacheMaxLoaded * capPressurePercent / 100
)

// Cron timing (below) keeps cap vs idle eviction visible in 5m Prometheus buckets; schedules are staggered.
func Test_V2_CRE_CacheSoak(t *testing.T) {
	numWorkflows := _defaultSoakNumWorkflows
	if os.Getenv("CRE_SOAK_NUM_WORKFLOWS") != "" {
		var err error
		numWorkflows, err = strconv.Atoi(os.Getenv("CRE_SOAK_NUM_WORKFLOWS"))
		if err != nil {
			t.Fatalf("failed to parse CRE_SOAK_NUM_WORKFLOWS: %v", err)
		}
	}

	soakDuration := parseDuration(os.Getenv("CRE_SOAK_DURATION"), defaultSoakDuration)

	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-cache-soak-test.toml"))
	testLogger := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// Do not use t.Context() here: it is cancelled before cleanup runs, which breaks chip-router
		// unregister and can leave gRPC Publish blocked on full log channels after WatchWorkflowLogs returns.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"

	testLogger.Info().
		Int("max_loaded_per_node", moduleCacheMaxLoaded).
		Int("target_workflows", numWorkflows).
		Int("target_loaded_mib", moduleCacheMaxLoaded*(_workflowModuleMiB+_workflowEngineOverheadMiB)).
		Msg("Deploying cache soak workflows")
	workflowIDs := t_helpers.CompileAndDeployWorkflowNTimes(t, testEnv, testLogger,
		func(i int) string { return fmt.Sprintf("cachetest%d", i) },
		func(i int) *crontypes.WorkflowConfig {
			return &crontypes.WorkflowConfig{Schedule: cacheSoakWorkflowSchedule(i)}
		},
		workflowFileLocation,
		numWorkflows,
		numberOfDeploymentKeys,
	)
	testLogger.Info().Int("count", len(workflowIDs)).Msg("All cache-test workflows deployed")
	nodeContainers := t_helpers.SnapshotNodeContainerRestarts(t, testEnv)
	startTime := time.Now()

	timeout := 2 * time.Minute
	testLogger.Info().
		Float64("timeout_minutes", timeout.Minutes()).
		Msg("Waiting for first workflow execution...")
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, "Amazing workflow user log", timeout)
	testLogger.Info().Dur("duration", soakDuration).Msg("First workflow execution confirmed, running cache soak...")

	t_helpers.AssertNodeLogs(t, testEnv, "Module cache enabled")

	testLogger.Info().
		Float64("duration_minutes", soakDuration.Minutes()).
		Float64("fast_interval_seconds", fastCronInterval.Seconds()).
		Float64("slow_interval_seconds", slowCronInterval.Seconds()).
		Float64("idle_timeout_seconds", moduleCacheIdleTimeout.Seconds()).
		Int("workflows", numWorkflows).
		Msg("Observing cache activity")
	observeUntil := time.Now().Add(soakDuration)
	for time.Now().Before(observeUntil) {
		time.Sleep(soakProgressLogInterval)
		testLogger.Info().Dur("remaining", time.Until(observeUntil).Round(time.Second)).Msg("Cache soak progress")
	}
	testLogger.Info().Msg("Cache soak complete")
	endTime := time.Now()

	// Check Prometheus metrics
	pc := framework.NewPrometheusQueryClient(framework.LocalPrometheusBaseURL)

	workflowDONs := testEnv.Dons.DonsWithFlag(cre.WorkflowDON)
	require.NotEmpty(t, workflowDONs, "no workflow DONs found")

	type wrappedQueryRangeResponse struct {
		NodeName string `json:"node_name"`
		Metric   string `json:"metric"`
		framework.QueryRangeResponse
	}

	type metric struct {
		query    string
		filename string
		metric   string
		step     time.Duration // interval between query points
	}

	metrics := []metric{
		{
			metric:   "platform_workflow_module_cache_reload_total",
			query:    fmt.Sprintf("sum by (source) (increase(platform_workflow_module_cache_reload_total{node_don=\"%%s\", node_index=\"%%d\"}[%s]))", cachePrometheusRange),
			filename: "metrics/cache_reload_increase.json",
			step:     defaultMetricStep,
		},
		{
			// Counter: WASM re-instantiated from on-disk cache (cold load path).
			metric:   "platform_workflow_module_cache_reload_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_reload_total{node_don=\"%%s\", node_index=\"%%d\", source=\"disk\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_reload_disk_increase.json",
			step:     defaultMetricStep,
		},
		{
			// Counter: WASM resurrected from in-memory weak ref (warm load path).
			metric:   "platform_workflow_module_cache_reload_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_reload_total{node_don=\"%%s\", node_index=\"%%d\", source=\"weak_ref\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_reload_memory_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_eviction_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_eviction_total{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_eviction_increase.json",
			step:     defaultMetricStep,
		},
		{
			// Gauge: peak loaded modules per step (cap pressure vs MaxLoaded).
			metric:   "platform_workflow_module_cache_loaded",
			query:    fmt.Sprintf("max_over_time(platform_workflow_module_cache_loaded{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_loaded.json",
			step:     defaultMetricStep,
		},
		{
			// Gauge: average evicted-module bytes not held in RAM per step.
			metric:   "platform_workflow_module_cache_memory_saved_bytes",
			query:    fmt.Sprintf("avg_over_time(platform_workflow_module_cache_memory_saved_bytes{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_memory_saved_bytes.json",
			step:     defaultMetricStep,
		},
		{
			// Gauge: total on-disk bytes under the module cache dir (DiskMonitor, ~1m tick).
			// max_over_time: peak footprint per step; correlates with deploy count and disk reloads.
			metric:   "platform_workflow_module_cache_disk_usage_bytes",
			query:    fmt.Sprintf("max_over_time(platform_workflow_module_cache_disk_usage_bytes{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_disk_usage_bytes.json",
			step:     defaultMetricStep,
		},
		{
			// Gauge: typical on-disk cache footprint during the step (smoothed).
			metric:   "platform_workflow_module_cache_disk_usage_bytes",
			query:    fmt.Sprintf("avg_over_time(platform_workflow_module_cache_disk_usage_bytes{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_disk_usage_avg_bytes.json",
			step:     defaultMetricStep,
		},
		{
			// Gauge: workflows fetched from registry on last sync tick (registered on this node).
			metric:   "platform_workflow_registry_syncer_fetched_workflows",
			query:    fmt.Sprintf("max_over_time(platform_workflow_registry_syncer_fetched_workflows{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/registry_fetched_workflows.json",
			step:     defaultMetricStep,
		},
		{
			// Gauge: workflow engines currently running on this node.
			metric:   "platform_workflow_registry_syncer_running_workflows",
			query:    fmt.Sprintf("max_over_time(platform_workflow_registry_syncer_running_workflows{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/registry_running_workflows.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_registry_syncer_completed_syncs_total",
			query:    fmt.Sprintf("increase(platform_workflow_registry_syncer_completed_syncs_total{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/registry_completed_syncs_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_registry_syncer_reconcile_events_backoff_total",
			query:    fmt.Sprintf("increase(platform_workflow_registry_syncer_reconcile_events_backoff_total{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/registry_reconcile_backoff_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric: "platform_engine_workflow_execution_started_count",
			query: fmt.Sprintf(
				"sum(increase(platform_engine_workflow_execution_started_count{node_don=\"%%s\", node_index=\"%%d\"}[%s]))",
				cachePrometheusRange,
			),
			filename: "metrics/engine_execution_started_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric: "platform_engine_workflow_execution_succeeded_count",
			query: fmt.Sprintf(
				"sum(increase(platform_engine_workflow_execution_succeeded_count{node_don=\"%%s\", node_index=\"%%d\"}[%s]))",
				cachePrometheusRange,
			),
			filename: "metrics/engine_execution_succeeded_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric: "platform_engine_workflow_execution_failed_count",
			query: fmt.Sprintf(
				"sum(increase(platform_engine_workflow_execution_failed_count{node_don=\"%%s\", node_index=\"%%d\"}[%s]))",
				cachePrometheusRange,
			),
			filename: "metrics/engine_execution_failed_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric: "platform_engine_trigger_event_received_total",
			query: fmt.Sprintf(
				"sum(increase(platform_engine_trigger_event_received_total{node_don=\"%%s\", node_index=\"%%d\"}[%s]))",
				cachePrometheusRange,
			),
			filename: "metrics/engine_trigger_event_received_increase.json",
			step:     defaultMetricStep,
		},
		// Engine schedule skew by module cache path (source: loaded | weak_ref | disk).
		{
			metric:   "platform_engine_trigger_queue_to_execution_start_seconds",
			query:    histogramQuantileQuery("platform_engine_trigger_queue_to_execution_start_seconds", 0.50),
			filename: "metrics/engine_trigger_skew_p50.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_trigger_queue_to_execution_start_seconds",
			query:    histogramQuantileQuery("platform_engine_trigger_queue_to_execution_start_seconds", 0.95),
			filename: "metrics/engine_trigger_skew_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_trigger_queue_to_execution_start_seconds",
			query:    histogramQuantileQueryBySource("platform_engine_trigger_queue_to_execution_start_seconds", 0.95, "loaded"),
			filename: "metrics/engine_trigger_skew_loaded_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_trigger_queue_to_execution_start_seconds",
			query:    histogramQuantileQueryBySource("platform_engine_trigger_queue_to_execution_start_seconds", 0.95, "weak_ref"),
			filename: "metrics/engine_trigger_skew_weak_ref_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_trigger_queue_to_execution_start_seconds",
			query:    histogramQuantileQueryBySource("platform_engine_trigger_queue_to_execution_start_seconds", 0.95, "disk"),
			filename: "metrics/engine_trigger_skew_disk_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_trigger_event_queue_wait_seconds",
			query:    histogramQuantileQuery("platform_engine_trigger_event_queue_wait_seconds", 0.95),
			filename: "metrics/engine_trigger_queue_wait_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_execution_semaphore_wait_seconds",
			query:    histogramQuantileQuery("platform_engine_execution_semaphore_wait_seconds", 0.95),
			filename: "metrics/engine_execution_semaphore_wait_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_engine_workflow_completed_time_seconds",
			query:    histogramQuantileQuery("platform_engine_workflow_completed_time_seconds", 0.95),
			filename: "metrics/engine_execution_duration_p95.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_version_mismatch_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_version_mismatch_total{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_version_mismatch.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_pin_exhausted_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_pin_exhausted_total{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_pin_exhausted.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_try_acquire_exhausted_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_try_acquire_exhausted_total{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_try_acquire_exhausted.json",
			step:     defaultMetricStep,
		},
		// average memory usage of the container over the last 10 minutes, unit:MBs
		// queried every 5 minutes
		// name is the Docker container name, this metric is gathered by cAdvisor
		{
			metric:   "container_memory_rss",
			query:    "avg_over_time(container_memory_rss{name=\"%s-node%d\"}[10m]) / 1024 / 1024",
			filename: "metrics/container_memory_rss.json",
			step:     5 * time.Minute,
		},
		// average CPU usage of the container over the last 10 minutes, unit:%
		// queried every 5 minutes
		// name is the Docker container name, this metric is gathered by cAdvisor
		{
			metric:   "container_cpu_usage_seconds_total",
			query:    "sum(rate(container_cpu_usage_seconds_total{name=\"%s-node%d\"}[10m])) * 100",
			filename: "metrics/container_cpu_usage_seconds_total.json",
			step:     5 * time.Minute,
		},
	}

	for _, metric := range metrics {
		results := make([]wrappedQueryRangeResponse, 0)
		for _, don := range workflowDONs {
			for _, node := range don.Nodes {
				query := fmt.Sprintf(metric.query, don.Name, node.Index)
				queryResponse, err := pc.QueryRange(framework.QueryRangeParams{
					Query: query,
					Start: startTime,
					End:   endTime,
					Step:  metric.step,
				})
				require.NoError(t, err, "failed to query Prometheus metrics, query:", query)
				results = append(results, wrappedQueryRangeResponse{
					NodeName:           node.Name,
					QueryRangeResponse: *queryResponse,
					Metric:             metric.metric,
				})
			}
		}

		require.NoError(t, saveJSONFile(metric.filename, results), "failed to save JSON file for metric:", metric.filename)
		testLogger.Info().Str("filename", metric.filename).Msg("Saved JSON file for metric")
	}

	t_helpers.AssertNodeContainersStable(t, nodeContainers)
	testLogger.Info().Msg("Node containers stable. None was restarted or OOM-killed.")
}

func crePerWorkflowSizeLimitMiB(size config.Size) int {
	return int(size / config.MByte)
}

// histogramQuantileQuery aggregates per-workflow engine histograms on a node (sum by le).
func histogramQuantileQuery(metric string, quantile float64) string {
	return fmt.Sprintf(
		`histogram_quantile(%g, sum by (le) (rate(%s_bucket{node_don="%%s", node_index="%%d"}[%s])))`,
		quantile, metric, cachePrometheusRange,
	)
}

func histogramQuantileQueryBySource(metric string, quantile float64, source string) string {
	return fmt.Sprintf(
		`histogram_quantile(%g, sum by (le) (rate(%s_bucket{node_don="%%s", node_index="%%d", source="%s"}[%s])))`,
		quantile, metric, source, cachePrometheusRange,
	)
}

// cacheSoakWorkflowSchedule returns a minute-granularity cron schedule aligned with the soak topology.
// Workflows with index divisible by cacheSoakSchedulePeriod use slowCronInterval (idle eviction tier);
// the rest use fastCronInterval to keep MaxLoaded full and drive cap eviction.
// Offsets stagger fires so cap and idle events land in different 5m Prometheus buckets.
func cacheSoakWorkflowSchedule(workflowIndex int) string {
	if workflowIndex%cacheSoakSchedulePeriod == 0 {
		offset := (workflowIndex / cacheSoakSchedulePeriod) % int(slowCronInterval.Minutes())
		return fmt.Sprintf("0 %d/%d * * * *", offset, int(slowCronInterval.Minutes()))
	}
	offset := workflowIndex % int(fastCronInterval.Minutes())
	return fmt.Sprintf("0 %d/%d * * * *", offset, int(fastCronInterval.Minutes()))
}

func saveJSONFile(path string, v any) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory for %q: %w", path, err)
		}
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON for %q: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil { //nolint:gosec // test artifact
		return fmt.Errorf("write file %q: %w", path, err)
	}
	return nil
}
