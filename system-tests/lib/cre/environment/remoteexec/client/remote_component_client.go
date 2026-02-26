package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	retry "github.com/avast/retry-go/v4"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/remoteexec/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/runtimecfg"
)

const (
	ComponentTypeBlockchain = "blockchain"
	ComponentTypeJD         = "jd"
	ComponentTypeNodeSet    = "nodeset"
	EnvRemoteAgentURL       = "CRE_REMOTE_AGENT_URL"
	EnvRemoteAgentPort      = "CRE_REMOTE_AGENT_PORT"
	defaultRemoteAgentPort  = 18080
)

type ComponentClient interface {
	StartComponent(ctx context.Context, envelope agent.StartComponentEnvelope) (*agent.StartComponentResponse, error)
}

type httpComponentClient struct {
	baseURL     string
	client      *http.Client
	maxAttempts int
	retryDelay  time.Duration
	checkHealth bool
}

type Runtime struct {
	AgentBaseURL string
	RemoteHostIP string
	Client       ComponentClient
}

type RuntimeInput struct {
	AgentBaseURL string
	RemoteHostIP string
	AgentPort    int
}

func newRemoteHTTPComponentClient(baseURL string) *httpComponentClient {
	return &httpComponentClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 4 * time.Minute,
		},
		maxAttempts: 3,
		retryDelay:  2 * time.Second,
		checkHealth: true,
	}
}

func ResolveRuntime(testLogger zerolog.Logger) (*Runtime, error) {
	return ResolveRuntimeWithInput(testLogger, RuntimeInput{})
}

func ResolveRuntimeWithInput(testLogger zerolog.Logger, input RuntimeInput) (*Runtime, error) {
	baseURL, err := resolveRemoteAgentBaseURL(testLogger, input)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve remote agent base URL: %w", err)
	}
	remoteHostIP, err := resolveRemoteHostIP(input, baseURL)
	if err != nil {
		return nil, err
	}
	client := newRemoteHTTPComponentClient(baseURL)
	runtime := &Runtime{
		AgentBaseURL: baseURL,
		RemoteHostIP: remoteHostIP,
		Client:       client,
	}

	// Best-effort compatibility check: fail on definitive protocol incompatibility,
	// but do not fail runtime resolution if status endpoint is temporarily unavailable.
	compatCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	status, statusErr := GetAgentStatus(compatCtx, runtime)
	if statusErr != nil {
		testLogger.Warn().Err(statusErr).Msg("skipping remote agent compatibility check (status unavailable)")
		return runtime, nil
	}
	if compatErr := checkCompatibilityStatus(status, nil); compatErr != nil {
		return nil, compatErr
	}

	return runtime, nil
}

func NewComponentClient(runtime *Runtime) (ComponentClient, error) {
	if runtime == nil {
		return nil, errors.New("resolved runtime is nil")
	}
	if runtime.Client != nil {
		return runtime.Client, nil
	}
	if strings.TrimSpace(runtime.AgentBaseURL) == "" {
		return nil, errors.New("resolved runtime is missing agent base url")
	}
	return newRemoteHTTPComponentClient(runtime.AgentBaseURL), nil
}

func (c *httpComponentClient) StartComponent(ctx context.Context, envelope agent.StartComponentEnvelope) (*agent.StartComponentResponse, error) {
	if c.checkHealth {
		if err := c.waitForHealth(ctx); err != nil {
			return nil, err
		}
	}

	var result *agent.StartComponentResponse
	err := retry.Do(
		func() error {
			var err error
			result, err = c.startComponentOnce(ctx, envelope)
			return err
		},
		retry.Attempts(uint(c.maxAttempts)),
		retry.Delay(c.retryDelay),
		retry.Context(ctx),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *httpComponentClient) startComponentOnce(ctx context.Context, envelope agent.StartComponentEnvelope) (*agent.StartComponentResponse, error) {
	body, err := json.Marshal(envelope)
	if err != nil {
		return nil, retry.Unrecoverable(pkgerrors.Wrap(err, "failed to encode start component envelope"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/components/start", bytes.NewReader(body))
	if err != nil {
		return nil, retry.Unrecoverable(pkgerrors.Wrap(err, "failed to create start component request"))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		if isRetriableNetworkError(err) {
			return nil, pkgerrors.Wrap(err, "failed to execute start component request")
		}
		return nil, retry.Unrecoverable(pkgerrors.Wrap(err, "failed to execute start component request"))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, retry.Unrecoverable(pkgerrors.Wrap(err, "failed to read start component response"))
	}

	var startResp agent.StartComponentResponse
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &startResp); err != nil {
			return nil, retry.Unrecoverable(pkgerrors.Wrap(err, "failed to decode start component response"))
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if startResp.Error != "" {
			if startResp.ErrorCode != "" {
				err = RemoteAgentError(startResp.ErrorCode, startResp.Error)
			} else {
				err = RemoteAgentError("remote_agent_error", startResp.Error)
			}
		} else {
			err = fmt.Errorf("start component request failed with status %s: %s", resp.Status, string(respBody))
		}

		if isRetriableStatus(resp.StatusCode) {
			return nil, err
		}
		return nil, retry.Unrecoverable(err)
	}
	if startResp.Error != "" {
		if startResp.ErrorCode != "" {
			return nil, retry.Unrecoverable(RemoteAgentError(startResp.ErrorCode, startResp.Error))
		}
		return nil, retry.Unrecoverable(RemoteAgentError("remote_agent_error", startResp.Error))
	}

	return &startResp, nil
}

func (c *httpComponentClient) waitForHealth(ctx context.Context) error {
	healthURL := c.baseURL + "/v1/health"
	return retry.Do(
		func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
			if err != nil {
				return retry.Unrecoverable(err)
			}
			resp, err := c.client.Do(req)
			if err != nil {
				return err
			}
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			return fmt.Errorf("%s: status %s", describeRemoteAgentHealthFailure(c.baseURL), resp.Status)
		},
		retry.Attempts(uint(c.maxAttempts)),
		retry.Delay(c.retryDelay),
		retry.Context(ctx),
		retry.LastErrorOnly(true),
	)
}

func describeRemoteAgentHealthFailure(baseURL string) string {
	return fmt.Sprintf(
		"failed remote CRE agent health check (%s/v1/health); verify the agent process is running and %s matches its listen port (or set %s explicitly)",
		baseURL,
		EnvRemoteAgentPort,
		EnvRemoteAgentURL,
	)
}

func isRetriableStatus(statusCode int) bool {
	return statusCode == http.StatusBadGateway || statusCode == http.StatusServiceUnavailable || statusCode == http.StatusGatewayTimeout
}

func isRetriableNetworkError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}

func RemoteAgentError(code, message string) error {
	return fmt.Errorf("remote agent error (%s): %s", code, message)
}

func resolveRemoteAgentBaseURL(testLogger zerolog.Logger, input RuntimeInput) (string, error) {
	if configured := strings.TrimSpace(input.AgentBaseURL); configured != "" {
		return configured, nil
	}
	if configured := strings.TrimSpace(os.Getenv(EnvRemoteAgentURL)); configured != "" {
		return configured, nil
	}
	remotePort, err := resolveRemoteAgentPort(input)
	if err != nil {
		return "", err
	}
	remoteHostIP, err := resolveRemoteHostIP(input, "")
	if err != nil {
		return "", err
	}
	testLogger.Debug().Str("remoteHostIP", remoteHostIP).Int("port", remotePort).Msg("resolved remote CRE agent base URL")
	return fmt.Sprintf("http://%s:%d", remoteHostIP, remotePort), nil
}

func resolveRemoteAgentPort(input RuntimeInput) (int, error) {
	if input.AgentPort > 0 {
		return input.AgentPort, nil
	}
	remotePort := defaultRemoteAgentPort
	if configuredPort := strings.TrimSpace(os.Getenv(EnvRemoteAgentPort)); configuredPort != "" {
		parsedPort, err := strconv.Atoi(configuredPort)
		if err != nil || parsedPort <= 0 || parsedPort > 65535 {
			return 0, fmt.Errorf("invalid %s: %q", EnvRemoteAgentPort, configuredPort)
		}
		remotePort = parsedPort
	}
	return remotePort, nil
}

func resolveRemoteHostIP(input RuntimeInput, baseURL string) (string, error) {
	if configured := strings.TrimSpace(input.RemoteHostIP); configured != "" {
		return configured, nil
	}
	if host, ok := hostFromBaseURL(baseURL); ok {
		return host, nil
	}
	return runtimecfg.DirectHostIP()
}

func hostFromBaseURL(baseURL string) (string, bool) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return "", false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", false
	}
	host := strings.TrimSpace(parsed.Hostname())
	if host == "" {
		return "", false
	}
	return host, true
}
