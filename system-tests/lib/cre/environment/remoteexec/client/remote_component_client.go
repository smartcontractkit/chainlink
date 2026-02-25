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
	EnvEC2AgentURL          = "CRE_EC2_AGENT_URL"
	EnvEC2AgentPort         = "CRE_EC2_AGENT_PORT"
	defaultEC2AgentPort     = 18080
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
	EC2HostIP    string
	Client       ComponentClient
}

type RuntimeInput struct {
	AgentBaseURL string
	EC2HostIP    string
	AgentPort    int
}

func newEC2HTTPComponentClient(baseURL string) *httpComponentClient {
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
	baseURL, err := resolveEC2AgentBaseURL(testLogger, input)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve EC2 agent base URL: %w", err)
	}
	ec2HostIP, err := resolveEC2HostIP(input)
	if err != nil {
		return nil, err
	}
	client := newEC2HTTPComponentClient(baseURL)
	runtime := &Runtime{
		AgentBaseURL: baseURL,
		EC2HostIP:    ec2HostIP,
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
	return newEC2HTTPComponentClient(runtime.AgentBaseURL), nil
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
			return fmt.Errorf("%s: status %s", describeEC2AgentHealthFailure(c.baseURL), resp.Status)
		},
		retry.Attempts(uint(c.maxAttempts)),
		retry.Delay(c.retryDelay),
		retry.Context(ctx),
		retry.LastErrorOnly(true),
	)
}

func describeEC2AgentHealthFailure(baseURL string) string {
	return fmt.Sprintf(
		"failed EC2 CRE agent health check (%s/v1/health); verify the agent process is running and %s matches its listen port (or set %s explicitly)",
		baseURL,
		EnvEC2AgentPort,
		EnvEC2AgentURL,
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

func resolveEC2AgentBaseURL(testLogger zerolog.Logger, input RuntimeInput) (string, error) {
	if configured := strings.TrimSpace(input.AgentBaseURL); configured != "" {
		return configured, nil
	}
	if configured := strings.TrimSpace(os.Getenv(EnvEC2AgentURL)); configured != "" {
		return configured, nil
	}
	remotePort, err := resolveEC2AgentPort(input)
	if err != nil {
		return "", err
	}
	ec2HostIP, err := resolveEC2HostIP(input)
	if err != nil {
		return "", err
	}
	testLogger.Debug().Str("ec2HostIP", ec2HostIP).Int("port", remotePort).Msg("resolved EC2 CRE agent base URL")
	return fmt.Sprintf("http://%s:%d", ec2HostIP, remotePort), nil
}

func resolveEC2AgentPort(input RuntimeInput) (int, error) {
	if input.AgentPort > 0 {
		return input.AgentPort, nil
	}
	remotePort := defaultEC2AgentPort
	if configuredPort := strings.TrimSpace(os.Getenv(EnvEC2AgentPort)); configuredPort != "" {
		parsedPort, err := strconv.Atoi(configuredPort)
		if err != nil || parsedPort <= 0 || parsedPort > 65535 {
			return 0, fmt.Errorf("invalid %s: %q", EnvEC2AgentPort, configuredPort)
		}
		remotePort = parsedPort
	}
	return remotePort, nil
}

func resolveEC2HostIP(input RuntimeInput) (string, error) {
	if configured := strings.TrimSpace(input.EC2HostIP); configured != "" {
		return configured, nil
	}
	return runtimecfg.DirectHostIP()
}
