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

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/stretchr/testify/require"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

func Test_V2_CRE_CacheSoak(t *testing.T) {
	numWorkflows := 10
	if os.Getenv("CRE_SOAK_NUM_WORKFLOWS") != "" {
		var err error
		numWorkflows, err = strconv.Atoi(os.Getenv("CRE_SOAK_NUM_WORKFLOWS"))
		if err != nil {
			t.Fatalf("failed to parse CRE_SOAK_NUM_WORKFLOWS: %v", err)
		}
	}

	cacheObserveWindow := 2 * time.Minute
	if os.Getenv("CRE_SOAK_CACHE_OBSERVE_WINDOW") != "" {
		var err error
		cacheObserveWindow, err = time.ParseDuration(os.Getenv("CRE_SOAK_CACHE_OBSERVE_WINDOW"))
		if err != nil {
			t.Fatalf("failed to parse CRE_SOAK_CACHE_OBSERVE_WINDOW: %v", err)
		}
	}

	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, "/configs/workflow-gateway-don-cache-test.toml"))
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

	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *",
	}

	startTime := time.Now()
	for i := range numWorkflows {
		t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, fmt.Sprintf("cachetest%d", i), &workflowConfig, workflowFileLocation)
	}
	testLogger.Info().Int("count", numWorkflows).Msg("All cache-test workflows deployed")

	t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, "Amazing workflow user log", 2*time.Minute)
	testLogger.Info().Dur("window", cacheObserveWindow).Msg("First workflow execution confirmed, observing cache activity...")

	// t_helpers.AssertNodeLogs(t, testEnv, "Module cache enabled")
	endTime := time.Now()

	// Check Prometheus metrics
	pc := framework.NewPrometheusQueryClient(framework.LocalPrometheusBaseURL)

	workflowDONs := testEnv.Dons.DonsWithFlag(cre.WorkflowDON)
	require.NotEmpty(t, workflowDONs, "no workflow DONs found")

	type wrappedQueryRangeResponse struct {
		NodeName string `json:"node_name"`
		framework.QueryRangeResponse
	}

	type metric struct {
		query    string
		filename string
	}

	metrics := []metric{
		{
			query:    "increase(platform_workflow_module_cache_eviction_total{node_don=\"%s\", node_index=\"%d\"}[1m])",
			filename: "metrics/cache_eviction_increase.json",
		},
		{
			query:    "sum by (source) (increase(platform_workflow_module_cache_reload_total{node_don=\"%s\", node_index=\"%d\"}[1m]))",
			filename: "metrics/cache_reload_increase.json",
		},
		{
			query:    "increase(platform_workflow_module_cache_memory_saved_bytes{node_don=\"%s\", node_index=\"%d\"}[1m])",
			filename: "metrics/cache_memory_saved_bytes.json",
		},
	}

	for _, metric := range metrics {
		results := make([]wrappedQueryRangeResponse, 0)
		for _, don := range workflowDONs {
			for _, node := range don.Nodes {
				query := fmt.Sprintf(metric.query, don.Name, node.Index)
				fmt.Println("query:", query)
				queryResponse, err := pc.QueryRange(framework.QueryRangeParams{
					Query: query,
					Start: startTime,
					End:   endTime,
					Step:  1 * time.Minute,
				})
				require.NoError(t, err, "failed to query Prometheus metrics, query:", query)
				results = append(results, wrappedQueryRangeResponse{
					NodeName:           node.Name,
					QueryRangeResponse: *queryResponse,
				})
			}
		}

		require.NoError(t, saveJSONFile(metric.filename, results), "failed to save JSON file for metric:", metric.filename)
		testLogger.Info().Str("filename", metric.filename).Msg("Saved JSON file for metric")
	}
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
