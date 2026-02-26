package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
)

const chipSinkLifecycleTimeout = 30 * time.Second

func StartRemoteChipTestSink(ctx context.Context, runtime *Runtime, req agent.ChipTestSinkStartRequest) (*agent.ChipTestSinkStartResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	var out agent.ChipTestSinkStartResponse
	if err := postAgentJSON(ctx, baseURL+"/v1/chip/sink/start", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func StopRemoteChipTestSink(ctx context.Context, runtime *Runtime) (*agent.ChipTestSinkStopResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	var out agent.ChipTestSinkStopResponse
	if err := postAgentJSON(ctx, baseURL+"/v1/chip/sink/stop", map[string]any{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetRemoteChipTestSinkStatus(ctx context.Context, runtime *Runtime) (*agent.ChipTestSinkStatusResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	var out agent.ChipTestSinkStatusResponse
	if err := getAgentJSON(ctx, baseURL+"/v1/chip/sink/status", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func GetRemoteChipTestSinkEvents(ctx context.Context, runtime *Runtime, since time.Time, limit int) (*agent.ChipTestSinkEventsResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	endpoint := baseURL + "/v1/chip/sink/events"
	query := make([]string, 0, 2)
	if limit > 0 {
		query = append(query, fmt.Sprintf("limit=%d", limit))
	}
	if !since.IsZero() {
		query = append(query, "since="+since.UTC().Format(time.RFC3339Nano))
	}
	if len(query) > 0 {
		endpoint += "?" + strings.Join(query, "&")
	}
	var out agent.ChipTestSinkEventsResponse
	if err := getAgentJSON(ctx, endpoint, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func postAgentJSON(ctx context.Context, endpoint string, payload any, target any) error {
	httpClient := &http.Client{Timeout: chipSinkLifecycleTimeout}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal agent request body for %s: %w", endpoint, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to build agent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call agent endpoint %s: %w", endpoint, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read agent response from %s: %w", endpoint, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var agentErr agent.StartComponentResponse
		if len(respBody) > 0 && json.Unmarshal(respBody, &agentErr) == nil && strings.TrimSpace(agentErr.Error) != "" {
			if agentErr.ErrorCode != "" {
				return RemoteAgentError(agentErr.ErrorCode, agentErr.Error)
			}
			return RemoteAgentError("remote_agent_error", agentErr.Error)
		}
		return fmt.Errorf("agent endpoint %s returned %s: %s", endpoint, resp.Status, strings.TrimSpace(string(respBody)))
	}
	if err := json.Unmarshal(respBody, target); err != nil {
		return fmt.Errorf("failed to decode agent response from %s: %w", endpoint, err)
	}
	return nil
}
