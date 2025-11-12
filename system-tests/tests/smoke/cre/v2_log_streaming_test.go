package cre

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	testutils "github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// LokiQueryResponse represents the response from Loki's query_range API
type LokiQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// ExecuteLogStreamingTest validates that logs with beholder_data_type are flowing to Loki
func ExecuteLogStreamingTest(t *testing.T, testEnv *ttypes.TestEnvironment) {
	testLogger := framework.L
	testLogger.Info().Msg("Starting Log Streaming Test")

	lokiURL := "http://localhost:3030"

	testLogger.Info().Msg("Waiting for Beholder logs to appear in Loki...")

	var beholderLogsCount int
	require.Eventually(t, func() bool {
		var err error
		beholderLogsCount, err = queryLokiForBeholderLogs(t.Context(), lokiURL, 120)
		if err != nil {
			testLogger.Debug().Err(err).Msg("Error querying Loki")
			return false
		}
		return beholderLogsCount > 0
	}, testutils.WaitTimeout(t), 5*time.Second, "Expected to find logs with beholder_data_type=zap_log_message in Loki within timeout")

	testLogger.Info().Int("beholderLogsCount", beholderLogsCount).Msg("Found logs with beholder_data_type")
	testLogger.Info().Msg("✅ Log Streaming Test PASSED: beholder_data_type logs are flowing to Loki")
}

func queryLokiForBeholderLogs(ctx context.Context, lokiBaseURL string, lastNSeconds int) (int, error) {
	return queryLoki(ctx, lokiBaseURL, `{service_name=~".*chainlink.*"} | json | beholder_data_type="zap_log_message"`, lastNSeconds)
}

// queryLoki is a generic function to query Loki with any query
func queryLoki(ctx context.Context, lokiBaseURL, query string, lastNSeconds int) (int, error) {
	end := time.Now()
	start := end.Add(-time.Duration(lastNSeconds) * time.Second)

	startNano := strconv.FormatInt(start.UnixNano(), 10)
	endNano := strconv.FormatInt(end.UnixNano(), 10)

	queryURL := lokiBaseURL + "/loki/api/v1/query_range"

	u, err := url.Parse(queryURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("start", startNano)
	q.Set("end", endNano)
	q.Set("limit", "10")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to build Loki request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to query Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return 0, fmt.Errorf("Failed to read response body: %w", readErr)
		}
		return 0, fmt.Errorf("Loki query failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var lokiResp LokiQueryResponse
	if err := json.Unmarshal(body, &lokiResp); err != nil {
		return 0, fmt.Errorf("failed to parse Loki response: %w", err)
	}

	count := 0
	for _, result := range lokiResp.Data.Result {
		count += len(result.Values)
	}

	return count, nil
}
