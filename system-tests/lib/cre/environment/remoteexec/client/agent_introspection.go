package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
)

const agentIntrospectionTimeout = 30 * time.Second

func GetAgentStatus(ctx context.Context, runtime *Runtime) (*agent.AgentStatusResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	var response agent.AgentStatusResponse
	if err := getAgentJSON(ctx, baseURL+"/v1/status", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func GetAgentLocks(ctx context.Context, runtime *Runtime) (*agent.AgentLocksResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	var response agent.AgentLocksResponse
	if err := getAgentJSON(ctx, baseURL+"/v1/locks", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func GetComponentLogs(ctx context.Context, runtime *Runtime, componentKey string, limit int) (*agent.ComponentLogsResponse, error) {
	baseURL, err := runtimeBaseURL(runtime)
	if err != nil {
		return nil, err
	}
	componentKey = strings.TrimSpace(componentKey)
	if componentKey == "" {
		return nil, errors.New("componentKey is required")
	}

	q := url.Values{}
	q.Set("componentKey", componentKey)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	var response agent.ComponentLogsResponse
	endpoint := baseURL + "/v1/components/logs?" + q.Encode()
	if err := getAgentJSON(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func runtimeBaseURL(runtime *Runtime) (string, error) {
	if runtime == nil {
		return "", errors.New("runtime is nil")
	}
	baseURL := strings.TrimSpace(runtime.AgentBaseURL)
	if baseURL == "" {
		return "", errors.New("runtime is missing agent base url")
	}
	return baseURL, nil
}

func getAgentJSON(ctx context.Context, endpoint string, target any) error {
	httpClient := &http.Client{Timeout: agentIntrospectionTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to build agent request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call agent endpoint %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read agent response from %s: %w", endpoint, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var agentErr agent.StartComponentResponse
		if len(body) > 0 && json.Unmarshal(body, &agentErr) == nil && strings.TrimSpace(agentErr.Error) != "" {
			if agentErr.ErrorCode != "" {
				return RemoteAgentError(agentErr.ErrorCode, agentErr.Error)
			}
			return RemoteAgentError("remote_agent_error", agentErr.Error)
		}
		return fmt.Errorf("agent endpoint %s returned %s: %s", endpoint, resp.Status, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode agent response from %s: %w", endpoint, err)
	}
	return nil
}
