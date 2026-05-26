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

const (
	cacheSoakConfigCap      = "/configs/workflow-gateway-don-cache-soak-test.toml"
	cacheSoakConfigWeakRef  = "/configs/workflow-gateway-don-cache-soak-weakref-test.toml"
	cachePrometheusRange    = "5m"
	defaultMetricStep       = 1 * time.Minute
	soakProgressLogInterval = 5 * time.Minute
	numberOfDeploymentKeys  = 20
)

// Cap-pressure profile (default Test_V2_CRE_CacheSoak): high registry/cap churn, mostly disk reloads.
const (
	capPressurePercent        = 1000
	moduleCacheMaxLoaded      = 25
	moduleCacheIdleTimeout    = 5 * time.Minute
	fastCronInterval          = 3 * time.Minute
	slowCronInterval          = 8 * time.Minute
	cacheSoakSchedulePeriod   = 3
	defaultCapSoakDuration    = 4 * time.Hour
	defaultCapSoakNumWorkflow = moduleCacheMaxLoaded * capPressurePercent / 100
)

// Weak-ref profile (Test_V2_CRE_CacheSoak_WeakRef): idle eviction to L2, GOGC=off, cron refire before GC.
// IdleTimeout in workflow-gateway-don-cache-soak-weakref-test.toml is 90s; cron period is weakRefCronInterval.
const (
	weakRefMaxLoaded           = 50 // mirrors workflow-gateway-don-cache-soak-weakref-test.toml MaxLoaded
	weakRefNumWorkflows        = 30
	weakRefCronInterval        = 2 * time.Minute
	defaultWeakRefSoakDuration = 45 * time.Minute
)

var (
	_workflowModuleMiB = crePerWorkflowSizeLimitMiB(
		cresettings.Default.PerWorkflow.WASMBinarySizeLimit.DefaultValue,
	)
	_workflowEngineOverheadMiB = crePerWorkflowSizeLimitMiB(
		cresettings.Default.PerWorkflow.WASMMemoryLimit.DefaultValue,
	)
)

type cacheSoakProfile struct {
	name            string
	configPath      string
	numWorkflows    int
	defaultDuration time.Duration
	schedule        func(workflowIndex int) string
	deploymentKeys  int
	artifactSubdir  string
	requireWeakRef  bool
	logMaxLoaded    int
}

type promQueryDef struct {
	query    string
	filename string
	metric   string
	step     time.Duration
}

type queryRangeExport struct {
	NodeName string `json:"node_name"`
	Metric   string `json:"metric"`
	framework.QueryRangeResponse
}

// Test_V2_CRE_CacheSoak exercises cap pressure and disk reload under heavy registry load.
func Test_V2_CRE_CacheSoak(t *testing.T) {
	runCacheSoak(t, cacheSoakProfile{
		name:            "cap-pressure",
		configPath:      cacheSoakConfigCap,
		numWorkflows:    defaultCapSoakNumWorkflow,
		defaultDuration: defaultCapSoakDuration,
		schedule:        cacheSoakCapPressureSchedule,
		deploymentKeys:  numberOfDeploymentKeys,
		artifactSubdir:  "",
		requireWeakRef:  false,
		logMaxLoaded:    moduleCacheMaxLoaded,
	})
}

// Test_V2_CRE_CacheSoak_WeakRef isolates L2 (weak.Pointer) reloads: no cap churn, idle eviction,
// GOGC=off on workflow nodes, and cron refire after IdleTimeout. Asserts weak_ref reload counters.
//
// Run:
//
//	go test -timeout 2h -run '^Test_V2_CRE_CacheSoak_WeakRef$' -v ./system-tests/tests/soak/cre/...
func Test_V2_CRE_CacheSoak_WeakRef(t *testing.T) {
	runCacheSoak(t, cacheSoakProfile{
		name:            "weak-ref",
		configPath:      cacheSoakConfigWeakRef,
		numWorkflows:    weakRefNumWorkflows,
		defaultDuration: defaultWeakRefSoakDuration,
		schedule:        cacheSoakWeakRefSchedule,
		deploymentKeys:  1,
		artifactSubdir:  "weakref",
		requireWeakRef:  true,
		logMaxLoaded:    weakRefMaxLoaded,
	})
}

func runCacheSoak(t *testing.T, profile cacheSoakProfile) {
	t.Helper()

	numWorkflows, err := soakNumWorkflows(profile.numWorkflows)
	require.NoError(t, err)

	soakDuration := parseDuration(os.Getenv("CRE_SOAK_DURATION"), profile.defaultDuration)

	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, profile.configPath))
	testLogger := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"

	testLogger.Info().
		Str("profile", profile.name).
		Int("max_loaded_per_node", profile.logMaxLoaded).
		Int("target_workflows", numWorkflows).
		Int("target_loaded_mib", profile.logMaxLoaded*(_workflowModuleMiB+_workflowEngineOverheadMiB)).
		Msg("Deploying cache soak workflows")
	workflowIDs := t_helpers.CompileAndDeployWorkflowNTimes(t, testEnv, testLogger,
		func(i int) string { return fmt.Sprintf("cachetest%d", i) },
		func(i int) *crontypes.WorkflowConfig {
			return &crontypes.WorkflowConfig{Schedule: profile.schedule(i)}
		},
		workflowFileLocation,
		numWorkflows,
		profile.deploymentKeys,
	)
	testLogger.Info().Int("count", len(workflowIDs)).Str("profile", profile.name).Msg("All cache-test workflows deployed")
	nodeContainers := t_helpers.SnapshotNodeContainerRestarts(t, testEnv)
	startTime := time.Now()

	timeout := 2 * time.Minute
	testLogger.Info().
		Str("profile", profile.name).
		Float64("timeout_minutes", timeout.Minutes()).
		Msg("Waiting for first workflow execution...")
	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, "Amazing workflow user log", timeout)
	testLogger.Info().
		Str("profile", profile.name).
		Dur("duration", soakDuration).
		Msg("First workflow execution confirmed, running cache soak...")

	t_helpers.AssertNodeLogs(t, testEnv, "Module cache enabled")

	testLogger.Info().
		Str("profile", profile.name).
		Float64("duration_minutes", soakDuration.Minutes()).
		Int("workflows", numWorkflows).
		Msg("Observing cache activity")
	observeUntil := time.Now().Add(soakDuration)
	for time.Now().Before(observeUntil) {
		time.Sleep(soakProgressLogInterval)
		testLogger.Info().
			Str("profile", profile.name).
			Dur("remaining", time.Until(observeUntil).Round(time.Second)).
			Msg("Cache soak progress")
	}
	testLogger.Info().Str("profile", profile.name).Msg("Cache soak complete")
	endTime := time.Now()

	pc := framework.NewPrometheusQueryClient(framework.LocalPrometheusBaseURL)

	workflowDONs := testEnv.Dons.DonsWithFlag(cre.WorkflowDON)
	require.NotEmpty(t, workflowDONs, "no workflow DONs found")

	if profile.requireWeakRef {
		assertWeakRefReloadsObserved(t, pc, workflowDONs, startTime, endTime)
	}

	saveCacheSoakMetrics(t, pc, workflowDONs, startTime, endTime, profile.artifactSubdir)

	t_helpers.AssertNodeContainersStable(t, nodeContainers)
	testLogger.Info().Str("profile", profile.name).Msg("Node containers stable. None was restarted or OOM-killed.")
}

func saveCacheSoakMetrics(
	t *testing.T,
	pc *framework.PrometheusQueryClient,
	workflowDONs []*cre.Don,
	startTime, endTime time.Time,
	artifactSubdir string,
) {
	t.Helper()

	for _, m := range cacheSoakMetricDefinitions() {
		filename := metricArtifactPath(artifactSubdir, m.filename)
		results := queryRangeExports(t, pc, workflowDONs, m, startTime, endTime)
		require.NoError(t, saveJSONFile(filename, results))
		framework.L.Info().Str("filename", filename).Msg("Saved JSON file for metric")
	}
}

func metricArtifactPath(artifactSubdir, filename string) string {
	if artifactSubdir == "" {
		return filename
	}
	return filepath.Join(artifactSubdir, filename)
}

func queryRangeExports(
	t *testing.T,
	pc *framework.PrometheusQueryClient,
	workflowDONs []*cre.Don,
	m promQueryDef,
	startTime, endTime time.Time,
) []queryRangeExport {
	t.Helper()

	capacity := 0
	for _, don := range workflowDONs {
		capacity += len(don.Nodes)
	}
	results := make([]queryRangeExport, 0, capacity)
	for _, don := range workflowDONs {
		for _, node := range don.Nodes {
			query := fmt.Sprintf(m.query, don.Name, node.Index)
			queryResponse, err := pc.QueryRange(framework.QueryRangeParams{
				Query: query,
				Start: startTime,
				End:   endTime,
				Step:  m.step,
			})
			require.NoError(t, err, "query Prometheus: %s", query)
			results = append(results, queryRangeExport{
				NodeName:           node.Name,
				QueryRangeResponse: *queryResponse,
				Metric:             m.metric,
			})
		}
	}
	return results
}

func cacheSoakMetricDefinitions() []promQueryDef {
	return []promQueryDef{
		{
			metric:   "platform_workflow_module_cache_reload_total",
			query:    fmt.Sprintf("sum by (source) (increase(platform_workflow_module_cache_reload_total{node_don=\"%%s\", node_index=\"%%d\"}[%s]))", cachePrometheusRange),
			filename: "metrics/cache_reload_increase.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_reload_total",
			query:    fmt.Sprintf("increase(platform_workflow_module_cache_reload_total{node_don=\"%%s\", node_index=\"%%d\", source=\"disk\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_reload_disk_increase.json",
			step:     defaultMetricStep,
		},
		{
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
			metric:   "platform_workflow_module_cache_loaded",
			query:    fmt.Sprintf("max_over_time(platform_workflow_module_cache_loaded{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_loaded.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_memory_saved_bytes",
			query:    fmt.Sprintf("avg_over_time(platform_workflow_module_cache_memory_saved_bytes{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_memory_saved_bytes.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_disk_usage_bytes",
			query:    fmt.Sprintf("max_over_time(platform_workflow_module_cache_disk_usage_bytes{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_disk_usage_bytes.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_module_cache_disk_usage_bytes",
			query:    fmt.Sprintf("avg_over_time(platform_workflow_module_cache_disk_usage_bytes{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/cache_disk_usage_avg_bytes.json",
			step:     defaultMetricStep,
		},
		{
			metric:   "platform_workflow_registry_syncer_fetched_workflows",
			query:    fmt.Sprintf("max_over_time(platform_workflow_registry_syncer_fetched_workflows{node_don=\"%%s\", node_index=\"%%d\"}[%s])", cachePrometheusRange),
			filename: "metrics/registry_fetched_workflows.json",
			step:     defaultMetricStep,
		},
		{
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
		{
			metric:   "container_memory_rss",
			query:    "avg_over_time(container_memory_rss{name=\"%s-node%d\"}[5m]) / 1024 / 1024",
			filename: "metrics/container_memory_rss.json",
			step:     5 * time.Minute,
		},
		{
			metric:   "container_memory_rss",
			query:    "max_over_time(container_memory_rss{name=\"%s-node%d\"}[5m]) / 1024 / 1024",
			filename: "metrics/container_memory_rss_max.json",
			step:     5 * time.Minute,
		},
		{
			metric:   "container_cpu_usage_seconds_total",
			query:    "sum(rate(container_cpu_usage_seconds_total{name=\"%s-node%d\"}[5m])) * 100",
			filename: "metrics/container_cpu_usage_seconds_total.json",
			step:     5 * time.Minute,
		},
	}
}

const weakRefReloadQuery = `sum(increase(platform_workflow_module_cache_reload_total{node_don="%s", node_index="%d", source="weak_ref"}[%s]))`

// assertWeakRefReloadsObserved requires L2 weak.Pointer reloads during the soak window.
func assertWeakRefReloadsObserved(
	t *testing.T,
	pc *framework.PrometheusQueryClient,
	workflowDONs []*cre.Don,
	startTime, endTime time.Time,
) {
	t.Helper()

	window := promqlDuration(endTime.Sub(startTime))
	for _, don := range workflowDONs {
		for _, node := range don.Nodes {
			query := fmt.Sprintf(weakRefReloadQuery, don.Name, node.Index, window)
			resp, err := pc.Query(query, endTime)
			require.NoError(t, err, "weak_ref query: %s", query)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data)
			require.NotEmpty(t, resp.Data.Result,
				"weak_ref series missing for don=%s node_index=%d", don.Name, node.Index)

			total, err := prometheusScalarValue(resp.Data.Result[0].Value)
			require.NoError(t, err)
			require.Greater(t, total, 0.0,
				"weak_ref reloads=0 for don=%s node_index=%d window=%s (GOGC=off, IdleTimeout < cron, workflows <= MaxLoaded)",
				don.Name, node.Index, window)
			t.Logf("weak_ref reloads don=%s node_index=%d total=%g window=%s", don.Name, node.Index, total, window)
		}
	}
}

func soakNumWorkflows(defaultN int) (int, error) {
	raw := os.Getenv("CRE_SOAK_NUM_WORKFLOWS")
	if raw == "" {
		return defaultN, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("CRE_SOAK_NUM_WORKFLOWS: %w", err)
	}
	return n, nil
}

func promqlDuration(d time.Duration) string {
	mins := int(d.Round(time.Minute) / time.Minute)
	if mins < 5 {
		return "5m"
	}
	return strconv.Itoa(mins) + "m"
}

func prometheusScalarValue(value []interface{}) (float64, error) {
	if len(value) < 2 {
		return 0, fmt.Errorf("prometheus value: want 2 elements, got %d", len(value))
	}
	if s, ok := value[1].(string); ok {
		return strconv.ParseFloat(s, 64)
	}
	if f, ok := value[1].(float64); ok {
		return f, nil
	}
	return 0, fmt.Errorf("prometheus value: unexpected type %T", value[1])
}

func crePerWorkflowSizeLimitMiB(size config.Size) int {
	return int(size / config.MByte)
}

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

// cacheSoakCapPressureSchedule: fast cron + slow tier for cap vs idle eviction under registry pressure.
func cacheSoakCapPressureSchedule(workflowIndex int) string {
	if workflowIndex%cacheSoakSchedulePeriod != 0 {
		offset := workflowIndex % int(fastCronInterval.Minutes())
		return fmt.Sprintf("0 %d/%d * * * *", offset, int(fastCronInterval.Minutes()))
	}
	offset := (workflowIndex / cacheSoakSchedulePeriod) % int(slowCronInterval.Minutes())
	return fmt.Sprintf("0 %d/%d * * * *", offset, int(slowCronInterval.Minutes()))
}

// cacheSoakWeakRefSchedule: every weakRefCronInterval, staggered; idle eviction fires before next tick.
func cacheSoakWeakRefSchedule(workflowIndex int) string {
	periodMin := int(weakRefCronInterval.Minutes())
	offset := workflowIndex % periodMin
	return fmt.Sprintf("0 %d/%d * * * *", offset, periodMin)
}

func saveJSONFile(path string, v any) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %q: %w", dir, err)
		}
	}

	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %q: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil { //nolint:gosec // test artifact
		return fmt.Errorf("write %q: %w", path, err)
	}
	return nil
}
