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

	testLogger.Info().Msg("Waiting for logs to accumulate...")
	time.Sleep(60 * time.Second)

	lokiURL := "http://localhost:3030"

	beholderLogsCount, err := queryLokiForBeholderLogs(lokiURL, 120)
	require.NoError(t, err, "Failed to query Loki")
	testLogger.Info().Int("beholderLogsCount", beholderLogsCount).Msg("Found logs with beholder_data_type")
	require.Positive(t, beholderLogsCount, "Expected to find logs with beholder_data_type=zap_log_message in Loki, but found none")

	testLogger.Info().Msg("✅ Log Streaming Test PASSED: beholder_data_type logs are flowing to Loki")
}

func queryLokiForBeholderLogs(lokiBaseURL string, lastNSeconds int) (int, error) {
	return queryLoki(lokiBaseURL, `{service_name=~".*chainlink.*"} | json | beholder_data_type="zap_log_message"`, lastNSeconds)
}

// queryLoki is a generic function to query Loki with any query
func queryLoki(lokiBaseURL, query string, lastNSeconds int) (int, error) {
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

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
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
