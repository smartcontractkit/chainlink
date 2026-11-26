package cre

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

const (
	// Fires once a year, so workflows stay registered and idle and nothing executes.
	capacityIdleSchedule = "0 0 0 1 1 *"

	capacityMaxDeployKeys = 20
	capacityPollInterval  = 30 * time.Second
	capacityUsageWindow   = 30 * time.Second // CPU is averaged over this window
	// covers the OTel collector keeping a previous environment's series alive (~5m)
	capacityStaleSeriesTimeout = 10 * time.Minute
	capacityResultsFile        = "metrics/capacity_steps.json"
)

type capacityNodeSample struct {
	Node               string  `json:"node"`
	Index              int     `json:"index"`
	RunningWorkflows   float64 `json:"running_workflows"`
	LoadedModules      float64 `json:"loaded_modules"`
	Goroutines         float64 `json:"goroutines"`
	MemoryMiB          float64 `json:"memory_mib"`
	GoHeapMiB          float64 `json:"go_heap_mib"`    // live Go heap of the chainlink process
	NonGoMiB           float64 `json:"non_go_rss_mib"` // chainlink process RSS not held by the Go runtime (cgo/wasmtime)
	CPUPercent         float64 `json:"cpu_percent"`    // 100 = one full core
	ModuleCacheDiskMiB float64 `json:"module_cache_disk_mib"`
	ActivationsLost    float64 `json:"activations_dropped_or_abandoned"`
}

type capacityStepResult struct {
	Step            int                      `json:"step"`
	TotalWorkflows  int                      `json:"total_workflows"` // deployed by this test so far
	RegisterSeconds float64                  `json:"register_seconds"`
	ActivateSeconds float64                  `json:"activate_seconds"`
	AllNodesReached bool                     `json:"all_nodes_reached"`
	Nodes           []capacityNodeSample     `json:"nodes"`
	Postgres        t_helpers.ContainerUsage `json:"postgres"` // shared by all workflow nodes
}

/*
Test_V2_CRE_WorkflowCapacity finds how many idle workflows the workflow DON can hold.
It first records a baseline (step 0, nothing deployed yet). Then it registers copies of the cron
example workflow in steps and, after each step, waits until every workflow node runs them all,
then records memory, CPU, goroutines and module-cache state per node, plus the DON's Postgres.
The last step where all nodes caught up without a restart or OOM kill is the capacity.
At the end it logs the per-node results and a per-step memory/CPU summary against the baseline.

Prerequisites (soak/cre/prepare_capacity_run.sh sets them up; run it before every run):
  - No local CRE running, or one started with configs/workflow-gateway-don-capacity-test.toml.
    Start from a fresh environment: the on-chain DON limit is 5000 active workflows. Workflows
    already running are tolerated (steps are measured on top of them) but logged as a warning.
    Before deploying, the test fails if any workflow node's module cache doesn't match the
    topology (or CTF_CONFIGS), e.g. because an environment from another topology was reused.
  - Observability stack running (`go run . obs up -f` in core/scripts/cre/environment).

For a 2-node workflow DON (more workflows per machine), set
CTF_CONFIGS=configs/workflow-gateway-don-capacity-test-2-nodes.toml.

Environment variables:
  - CRE_CAPACITY_TARGET: workflows to deploy in total (default 5000).
  - CRE_CAPACITY_STEP: workflows registered per step (default 500, min 2).
  - CRE_CAPACITY_STEP_TIMEOUT: max wait for nodes to run a step's workflows (default 45m).
  - CRE_CAPACITY_SETTLE: wait before each sample, baseline included, so idle modules get
    evicted and memory/CPU settle (default 2m).

Workflows are left registered on exit; tear down with `go run . env stop`.
Per-step results are written to metrics/capacity_steps.json.
*/
//nolint:paralleltest // measures the whole environment; anything running alongside skews the result
func Test_V2_CRE_WorkflowCapacity(t *testing.T) {
	target := envInt(t, "CRE_CAPACITY_TARGET", 5000)
	step := envInt(t, "CRE_CAPACITY_STEP", 500)
	stepTimeout := parseDuration(os.Getenv("CRE_CAPACITY_STEP_TIMEOUT"), 45*time.Minute)
	settle := parseDuration(os.Getenv("CRE_CAPACITY_SETTLE"), 2*time.Minute)
	require.GreaterOrEqual(t, target, 2, "CRE_CAPACITY_TARGET must be at least 2")
	require.GreaterOrEqual(t, step, 2, "CRE_CAPACITY_STEP must be at least 2")

	pc := framework.NewPrometheusQueryClient(framework.LocalPrometheusBaseURL)
	_, promErr := pc.Query("up", time.Now())
	require.NoError(t, promErr, "Prometheus is not reachable at %s; start it with `go run . obs up -f`", framework.LocalPrometheusBaseURL)

	tconf := t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-capacity-test.toml")
	freshEnvironment := !envconfig.LocalCREStateFileExists(tconf.RelativePathToRepoRoot)
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, tconf)
	testLogger := framework.L
	workflowDON := testEnv.Dons.MustWorkflowDON()
	t_helpers.RequireModuleCacheAsConfigured(t, testEnv, workflowDON.Name)
	if freshEnvironment {
		waitForStaleNodeSeries(t, pc, workflowDON)
	}

	var results []capacityStepResult
	t.Cleanup(func() {
		logCapacityResults(results)
		if err := saveJSONFile(capacityResultsFile, results); err != nil {
			testLogger.Error().Err(err).Msg("failed to save capacity results")
		}
	})

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/cron/main.go"
	nodeContainers := t_helpers.SnapshotNodeContainerRestarts(t, testEnv)

	testLogger.Info().Dur("settle", settle).Msg("Recording baseline before deploying workflows")
	time.Sleep(settle)
	baseline := capacityStepResult{AllNodesReached: true}
	baseline.Nodes, baseline.Postgres = sampleCapacityStep(t, pc, workflowDON)
	results = append(results, baseline)
	logStepUsage(baseline)
	for _, n := range baseline.Nodes {
		if n.RunningWorkflows > 0 {
			testLogger.Warn().Str("node", n.Node).Float64("running_workflows", n.RunningWorkflows).
				Msg("Environment already runs workflows; steps are measured on top of them and share the on-chain DON limit. Run `go run . env stop` for a clean run")
		}
	}

	deployed := 0
	for stepIdx := 1; deployed < target; stepIdx++ {
		n := min(step, target-deployed)
		if target-deployed-n == 1 {
			n++ // CompileAndDeployWorkflowNTimes needs at least 2, so absorb a remainder of 1
		}
		offset := deployed

		testLogger.Info().Int("step", stepIdx).Int("deploying", n).Int("total_after", deployed+n).Msg("Capacity step starting")
		registerStart := time.Now()
		t_helpers.CompileAndDeployWorkflowNTimes(t, testEnv, testLogger,
			func(i int) string { return fmt.Sprintf("capacity%05d", offset+i) },
			func(int) *crontypes.WorkflowConfig { return &crontypes.WorkflowConfig{Schedule: capacityIdleSchedule} },
			workflowFileLocation,
			n,
			min(capacityMaxDeployKeys, n),
			t_helpers.WithoutCleanup(),
		)
		deployed += n
		result := capacityStepResult{
			Step:            stepIdx,
			TotalWorkflows:  deployed,
			RegisterSeconds: time.Since(registerStart).Seconds(),
		}

		activateStart := time.Now()
		var failure string
		result.AllNodesReached, failure = waitForRunningWorkflows(t, pc, workflowDON, baseline, offset, deployed, stepTimeout)
		result.ActivateSeconds = time.Since(activateStart).Seconds()

		if result.AllNodesReached {
			time.Sleep(settle)
		}
		result.Nodes, result.Postgres = sampleCapacityStep(t, pc, workflowDON)
		results = append(results, result)
		logStepUsage(result)

		t_helpers.AssertNodeContainersStable(t, nodeContainers)
		require.True(t, result.AllNodesReached, "step %d (%d workflows): %s", stepIdx, deployed, failure)
		testLogger.Info().Int("step", stepIdx).Int("total", deployed).Msg("Capacity step passed")
	}
}

// waitForRunningWorkflows polls until every node runs its baseline count plus deployed workflows.
// It stops early if a node drops or abandons an activation, since those workflows are not
// retried until the node restarts, or if a node reports fewer workflows than the previous step
// reached, which means its metrics stopped arriving.
func waitForRunningWorkflows(t *testing.T, pc *framework.PrometheusQueryClient, don *cre.Don, baseline capacityStepResult, previouslyDeployed, deployed int, timeout time.Duration) (bool, string) {
	t.Helper()

	baselineRunning := make(map[int]float64, len(baseline.Nodes))
	baselineLost := make(map[int]float64, len(baseline.Nodes))
	for _, n := range baseline.Nodes {
		baselineRunning[n.Index] = n.RunningWorkflows
		baselineLost[n.Index] = n.ActivationsLost
	}

	deadline := time.Now().Add(timeout)
	for {
		reached := true
		var progress []string
		for _, node := range don.Nodes {
			want := baselineRunning[node.Index] + float64(deployed)
			running := promScalar(t, pc, nodeQuery("max(platform_workflow_registry_syncer_running_workflows{%s})", don, node.Index))
			if lost := promScalar(t, pc, nodeQuery(activationsLostQuery, don, node.Index)) - baselineLost[node.Index]; lost > 0 {
				return false, fmt.Sprintf("node %d dropped or abandoned %.0f workflow activations", node.Index, lost)
			}
			if alreadyRunning := baselineRunning[node.Index] + float64(previouslyDeployed); running < alreadyRunning {
				return false, fmt.Sprintf("node %d reports %.0f running workflows, fewer than the %.0f it already ran: its metrics stopped "+
					"reaching Prometheus, either because the OTel collector rejects its exports (node logs: \"Telemetry error\" ... larger than max) "+
					"or because Prometheus scrapes of the collector time out (up{job=%q}). soak/cre/prepare_capacity_run.sh configures the "+
					"collector for both", node.Index, running, alreadyRunning, "compose-otel-collector-1")
			}
			if running < want {
				reached = false
			}
			progress = append(progress, fmt.Sprintf("node%d=%.0f/%.0f", node.Index, running, want))
		}
		framework.L.Info().Int("deployed", deployed).Str("running/want", strings.Join(progress, " ")).Msg("Waiting for workflows to start")
		if reached {
			return true, ""
		}
		if time.Now().After(deadline) {
			return false, fmt.Sprintf("not all nodes reached their running workflow count within %s (%s)", timeout, strings.Join(progress, " "))
		}
		time.Sleep(capacityPollInterval)
	}
}

// waitForStaleNodeSeries waits until no node of a freshly started DON reports running workflows or
// loaded modules. Node metrics reach Prometheus through the OTel collector, which keeps exporting
// the last values of a previous environment's nodes (same node_don/node_index labels) for a few
// minutes; a baseline taken before they expire is wrong.
func waitForStaleNodeSeries(t *testing.T, pc *framework.PrometheusQueryClient, don *cre.Don) {
	t.Helper()

	deadline := time.Now().Add(capacityStaleSeriesTimeout)
	for {
		var stale []string
		for _, node := range don.Nodes {
			running := promScalar(t, pc, nodeQuery("max(platform_workflow_registry_syncer_running_workflows{%s})", don, node.Index))
			loaded := promScalar(t, pc, nodeQuery("max(platform_workflow_module_cache_loaded{%s})", don, node.Index))
			if running > 0 || loaded > 0 {
				stale = append(stale, fmt.Sprintf("node%d running=%.0f loaded=%.0f", node.Index, running, loaded))
			}
		}
		if len(stale) == 0 {
			return
		}
		require.Falsef(t, time.Now().After(deadline),
			"a fresh environment still reports workflows after %s: %s", capacityStaleSeriesTimeout, strings.Join(stale, ", "))
		framework.L.Info().Str("stale", strings.Join(stale, ", ")).Msg("Waiting for metrics of a previous environment to expire")
		time.Sleep(capacityPollInterval)
	}
}

const activationsLostQuery = "(sum(platform_workflow_registry_syncer_activation_dropped_total{%[1]s}) or vector(0)) + (sum(platform_workflow_registry_syncer_activation_abandoned_total{%[1]s}) or vector(0))"

// sampleCapacityStep reads node metrics from Prometheus and memory/CPU of the node containers and
// the DON's Postgres from Docker (cAdvisor on Docker Desktop has no per-container series).
func sampleCapacityStep(t *testing.T, pc *framework.PrometheusQueryClient, don *cre.Don) ([]capacityNodeSample, t_helpers.ContainerUsage) {
	t.Helper()

	postgres := don.Name + "-ns-postgresql"
	containers := make([]string, 0, 1+len(don.Nodes))
	containers = append(containers, postgres)
	for _, node := range don.Nodes {
		containers = append(containers, fmt.Sprintf("%s-node%d", don.Name, node.Index))
	}
	usage := t_helpers.SampleContainerUsage(t, containers, capacityUsageWindow)

	samples := make([]capacityNodeSample, 0, len(don.Nodes))
	for _, node := range don.Nodes {
		container := fmt.Sprintf("%s-node%d", don.Name, node.Index)
		samples = append(samples, capacityNodeSample{
			Node:               container,
			Index:              node.Index,
			RunningWorkflows:   promScalar(t, pc, nodeQuery("max(platform_workflow_registry_syncer_running_workflows{%s})", don, node.Index)),
			LoadedModules:      promScalar(t, pc, nodeQuery("max(platform_workflow_module_cache_loaded{%s})", don, node.Index)),
			ModuleCacheDiskMiB: promScalar(t, pc, nodeQuery("max(platform_workflow_module_cache_disk_usage_bytes{%s})", don, node.Index)) / 1024 / 1024,
			ActivationsLost:    promScalar(t, pc, nodeQuery(activationsLostQuery, don, node.Index)),
			Goroutines:         promScalar(t, pc, fmt.Sprintf(`max(go_goroutines{job=%q})`, container)),
			GoHeapMiB:          promScalar(t, pc, fmt.Sprintf(`max(go_memstats_heap_inuse_bytes{job=%q})`, container)) / 1024 / 1024,
			NonGoMiB: promScalar(t, pc, fmt.Sprintf(
				`max(process_resident_memory_bytes{job=%[1]q}) - (max(go_memstats_sys_bytes{job=%[1]q}) - max(go_memstats_heap_released_bytes{job=%[1]q}))`,
				container)) / 1024 / 1024,
			MemoryMiB:  usage[container].MemoryMiB,
			CPUPercent: usage[container].CPUPercent,
		})
	}
	return samples, usage[postgres]
}

// nodeQuery fills the node selector into every %s / %[1]s verb of query.
func nodeQuery(query string, don *cre.Don, nodeIndex int) string {
	return fmt.Sprintf(query, fmt.Sprintf(`node_don=%q, node_index="%d"`, don.Name, nodeIndex))
}

// promScalar returns the first sample of an instant query, or 0 when the series doesn't exist yet.
func promScalar(t *testing.T, pc *framework.PrometheusQueryClient, query string) float64 {
	t.Helper()
	resp, err := pc.Query(query, time.Now())
	require.NoError(t, err, "Prometheus query failed: %s", query)
	if resp.Data == nil || len(resp.Data.Result) == 0 || len(resp.Data.Result[0].Value) < 2 {
		return 0
	}
	raw, ok := resp.Data.Result[0].Value[1].(string)
	require.True(t, ok, "unexpected Prometheus value type for query %s", query)
	v, err := strconv.ParseFloat(raw, 64)
	require.NoError(t, err, "failed to parse Prometheus value %q for query %s", raw, query)
	return v
}

// stepUsage aggregates a step's per-node memory and CPU.
type stepUsage struct {
	totalMemMiB, avgMemMiB, maxMemMiB float64
	totalCPU, maxCPU                  float64
}

func usageOf(r capacityStepResult) stepUsage {
	var u stepUsage
	for _, n := range r.Nodes {
		u.totalMemMiB += n.MemoryMiB
		u.maxMemMiB = max(u.maxMemMiB, n.MemoryMiB)
		u.totalCPU += n.CPUPercent
		u.maxCPU = max(u.maxCPU, n.CPUPercent)
	}
	if len(r.Nodes) > 0 {
		u.avgMemMiB = u.totalMemMiB / float64(len(r.Nodes))
	}
	return u
}

func logStepUsage(r capacityStepResult) {
	u := usageOf(r)
	framework.L.Info().
		Int("step", r.Step).
		Int("total_workflows", r.TotalWorkflows).
		Float64("nodes_mem_total_mib", u.totalMemMiB).
		Float64("nodes_mem_max_mib", u.maxMemMiB).
		Float64("nodes_cpu_total_pct", u.totalCPU).
		Float64("postgres_mem_mib", r.Postgres.MemoryMiB).
		Float64("postgres_cpu_pct", r.Postgres.CPUPercent).
		Msg("Workflow DON resource usage")
}

func logCapacityResults(results []capacityStepResult) {
	if len(results) == 0 {
		return
	}

	var b strings.Builder
	fmt.Fprintf(&b, "\n%-5s %-6s %-9s %-9s %-16s %-8s %-7s %-10s %-9s %-9s %-10s %-8s %-9s\n",
		"step", "total", "reg_s", "act_s", "node", "running", "loaded", "goroutine", "mem_mib", "go_heap", "non_go_mib", "cpu_pct", "cache_mib")
	for _, r := range results {
		for _, n := range r.Nodes {
			fmt.Fprintf(&b, "%-5d %-6d %-9.0f %-9.0f %-16s %-8.0f %-7.0f %-10.0f %-9.0f %-9.0f %-10.0f %-8.1f %-9.0f\n",
				r.Step, r.TotalWorkflows, r.RegisterSeconds, r.ActivateSeconds, n.Node,
				n.RunningWorkflows, n.LoadedModules, n.Goroutines, n.MemoryMiB, n.GoHeapMiB, n.NonGoMiB, n.CPUPercent, n.ModuleCacheDiskMiB)
		}
	}
	framework.L.Info().Msg("Workflow capacity results per node:" + b.String())

	// Memory/CPU per step for the whole workflow DON; step 0 is the baseline before deploying.
	// Per node: mib_per_wf is the average node's growth over the baseline divided by all workflows
	// deployed; step_mib_per_wf only counts the step's own workflows, so a one-time jump (e.g. memory
	// the allocator keeps after the first compile burst) doesn't hide the marginal cost.
	baseline := usageOf(results[0])
	b.Reset()
	fmt.Fprintf(&b, "\n%-5s %-6s %-14s %-13s %-13s %-13s %-11s %-16s %-12s %-10s %-8s %-8s\n",
		"step", "total", "mem_total_mib", "mem_avg_mib", "mem_max_mib", "mem_vs_base", "mib_per_wf", "step_mib_per_wf", "cpu_total_%", "cpu_max_%", "pg_mib", "pg_cpu_%")
	for i, r := range results {
		u := usageOf(r)
		perWorkflow, stepPerWorkflow := "-", "-"
		if r.TotalWorkflows > 0 {
			perWorkflow = fmt.Sprintf("%.2f", (u.avgMemMiB-baseline.avgMemMiB)/float64(r.TotalWorkflows))
		}
		if i > 0 && r.TotalWorkflows > results[i-1].TotalWorkflows {
			prev := usageOf(results[i-1])
			stepPerWorkflow = fmt.Sprintf("%.2f", (u.avgMemMiB-prev.avgMemMiB)/float64(r.TotalWorkflows-results[i-1].TotalWorkflows))
		}
		fmt.Fprintf(&b, "%-5d %-6d %-14.0f %-13.0f %-13.0f %-13s %-11s %-16s %-12.1f %-10.1f %-8.0f %-8.1f\n",
			r.Step, r.TotalWorkflows, u.totalMemMiB, u.avgMemMiB, u.maxMemMiB,
			fmt.Sprintf("%+.0f", u.totalMemMiB-baseline.totalMemMiB), perWorkflow, stepPerWorkflow, u.totalCPU, u.maxCPU,
			r.Postgres.MemoryMiB, r.Postgres.CPUPercent)
	}
	framework.L.Info().Msg("Workflow DON memory/CPU by step (step 0 = baseline):" + b.String())
}

func envInt(t *testing.T, name string, defaultValue int) int {
	t.Helper()
	raw := os.Getenv(name)
	if raw == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(raw)
	require.NoError(t, err, "failed to parse %s", name)
	return v
}
