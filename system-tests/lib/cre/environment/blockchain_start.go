package environment

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

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/adapters"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

const (
	componentTypeBlockchain = "blockchain"
	componentTypeJD         = "jd"
	componentTypeNodeSet    = "nodeset"
	envLocalAgentURL        = "CRE_LOCAL_AGENT_URL"
	envEC2AgentURL          = "CRE_EC2_AGENT_URL"
	envEC2InstanceID        = "CRE_EC2_INSTANCE_ID"
	envEC2AgentPort         = "CRE_EC2_AGENT_PORT"
	envAgentMode            = "CRE_AGENT_MODE"
	ec2Region               = "us-west-2"
	defaultEC2AgentPort     = 8080
)

type startComponentEnvelope = agent.StartComponentEnvelope
type startComponentRequest = agent.StartComponentPayload
type startComponentResult = agent.StartComponentResponse

type componentClient interface {
	StartComponent(ctx context.Context, envelope agent.StartComponentEnvelope) (*agent.StartComponentResponse, error)
}

type httpComponentClient struct {
	baseURL     string
	client      *http.Client
	maxAttempts int
	retryDelay  time.Duration
	checkHealth bool
}

func newHTTPComponentClient(baseURL string) *httpComponentClient {
	return &httpComponentClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 4 * time.Minute,
		},
		maxAttempts: 1,
		retryDelay:  0,
		checkHealth: false,
	}
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
				err = remoteAgentError(startResp.ErrorCode, startResp.Error)
			} else {
				err = remoteAgentError("remote_agent_error", startResp.Error)
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
			return nil, retry.Unrecoverable(remoteAgentError(startResp.ErrorCode, startResp.Error))
		}
		return nil, retry.Unrecoverable(remoteAgentError("remote_agent_error", startResp.Error))
	}

	return &startResp, nil
}

func (c *httpComponentClient) waitForHealth(ctx context.Context) error {
	return retry.Do(
		func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/health", nil)
			if err != nil {
				return retry.Unrecoverable(pkgerrors.Wrap(err, "failed to create EC2 agent health request"))
			}

			resp, err := c.client.Do(req)
			if err != nil {
				return pkgerrors.Wrap(err, describeEC2AgentHealthFailure(c.baseURL))
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
		envEC2AgentPort,
		envEC2AgentURL,
	)
}

func isRetriableStatus(statusCode int) bool {
	return statusCode == http.StatusBadGateway || statusCode == http.StatusServiceUnavailable || statusCode == http.StatusGatewayTimeout
}

func isRetriableNetworkError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}

func newStartComponentClient(testLogger zerolog.Logger, tunnelManager tunnel.Manager) (componentClient, error) {
	agentMode := strings.TrimSpace(os.Getenv(envAgentMode))
	if strings.EqualFold(agentMode, "ec2") {
		baseURL, err := resolveEC2AgentBaseURL(testLogger, tunnelManager)
		if err != nil {
			return nil, err
		}
		return newEC2HTTPComponentClient(baseURL), nil
	}

	baseURL := os.Getenv(envLocalAgentURL)
	if baseURL == "" {
		return nil, fmt.Errorf("%s must be set for remote component startup", envLocalAgentURL)
	}
	return newHTTPComponentClient(baseURL), nil
}

func resolveEC2AgentBaseURL(testLogger zerolog.Logger, tunnelManager tunnel.Manager) (string, error) {
	if configured := os.Getenv(envEC2AgentURL); configured != "" {
		return configured, nil
	}
	if tunnelManager == nil {
		return "", errors.New("tunnel manager is required to auto-open ec2 agent tunnel")
	}

	instanceID := strings.TrimSpace(os.Getenv(envEC2InstanceID))
	if instanceID == "" {
		return "", fmt.Errorf("%s must be set when %s=ec2 and %s is not provided", envEC2InstanceID, envAgentMode, envEC2AgentURL)
	}

	remotePort := defaultEC2AgentPort
	if configuredPort := strings.TrimSpace(os.Getenv(envEC2AgentPort)); configuredPort != "" {
		parsedPort, err := strconv.Atoi(configuredPort)
		if err != nil || parsedPort <= 0 || parsedPort > 65535 {
			return "", fmt.Errorf("invalid %s: %q", envEC2AgentPort, configuredPort)
		}
		remotePort = parsedPort
	}

	bindings, err := tunnelManager.Start(context.Background(), []tunnel.EndpointRef{
		{
			ComponentID:  "agent",
			EndpointName: "api",
			Scheme:       "http",
			Host:         "127.0.0.1",
			Port:         remotePort,
			OriginalURL:  fmt.Sprintf("http://127.0.0.1:%d", remotePort),
		},
	})
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to open ssm tunnel to ec2 agent")
	}
	if len(bindings) == 0 {
		return "", errors.New("failed to open ssm tunnel to ec2 agent: no bindings returned")
	}

	testLogger.Info().
		Str("instanceID", instanceID).
		Int("remotePort", remotePort).
		Int("localPort", bindings[0].LocalPort).
		Msg("Opened SSM tunnel to EC2 agent")

	return bindings[0].LocalURL, nil
}

func blockchainFromOutput(testLogger zerolog.Logger, output *blockchain.Output) (blockchains.Blockchain, error) {
	if output == nil {
		return nil, pkgerrors.New("blockchain output is nil")
	}

	if output.Type != blockchain.TypeAnvil {
		return nil, fmt.Errorf("remote blockchain reconstruction supports only %s in phase 2A, got %s", blockchain.TypeAnvil, output.Type)
	}

	return evm.FromOutput(testLogger, output)
}

func prettifyAgentLogLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return ""
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return trimmed
	}

	message, _ := payload["message"].(string)
	if message == "" {
		return trimmed
	}

	level, _ := payload["level"].(string)
	if level == "" {
		level = "info"
	}

	cmd, _ := payload["Cmd"].(string)
	if cmd != "" {
		return fmt.Sprintf("[%s] %s (cmd=%s)", level, message, cmd)
	}

	return fmt.Sprintf("[%s] %s", level, message)
}

func validatePhase2ARemoteBlockchainInput(input *blockchain.Input) error {
	if input == nil {
		return errors.New("blockchain input is nil")
	}
	if input.Type != blockchain.TypeAnvil {
		return fmt.Errorf("remote target in phase 2A supports only %s, got %s", blockchain.TypeAnvil, input.Type)
	}
	return nil
}

func startBlockchainsWithTargets(
	ctx context.Context,
	testLogger zerolog.Logger,
	configuredBlockchains []*config.Blockchain,
	deployers map[blockchain.ChainFamily]blockchains.Deployer,
	tunnelManager tunnel.Manager,
	rewriteInternalForLocalNodes bool,
) (*blockchains.DeployedBlockchains, error) {
	blockchainInputs, err := config.ResolveBlockchainInputs(configuredBlockchains)
	if err != nil {
		return nil, err
	}

	localIdx := make([]int, 0, len(configuredBlockchains))
	localInputs := make([]*blockchain.Input, 0, len(configuredBlockchains))
	remoteIdx := make([]int, 0, len(configuredBlockchains))
	for idx, configuredBlockchain := range configuredBlockchains {
		if configuredBlockchain.Target == config.TargetRemote {
			remoteIdx = append(remoteIdx, idx)
			continue
		}
		localIdx = append(localIdx, idx)
		localInputs = append(localInputs, configuredBlockchain.InputRef())
	}

	outputs := make([]blockchains.Blockchain, len(configuredBlockchains))

	if len(localInputs) > 0 {
		for i, idx := range localIdx {
			deployedOutput, err := agent.DeployBlockchainComponent(ctx, deployers, localInputs[i])
			if err != nil {
				return nil, err
			}
			reconstructedBlockchain, err := blockchainFromOutput(testLogger, deployedOutput)
			if err != nil {
				return nil, err
			}
			outputs[idx] = reconstructedBlockchain
		}
	}

	if len(remoteIdx) > 0 {
		startClient, err := newStartComponentClient(testLogger, tunnelManager)
		if err != nil {
			return nil, err
		}

		for _, idx := range remoteIdx {
			input := blockchainInputs[idx]
			configured := configuredBlockchains[idx]
			if err := validatePhase2ARemoteBlockchainInput(input); err != nil {
				return nil, err
			}

			payload := agent.StartComponentPayload{
				ComponentType: componentTypeBlockchain,
				Blockchain:    input,
				ReusePolicy:   string(configured.RemoteStartPolicy),
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to encode blockchain payload")
			}

			response, err := startClient.StartComponent(ctx, agent.StartComponentEnvelope{
				SchemaVersion: agent.SchemaVersionV1,
				Operation:     agent.OperationStartComponent,
				Payload:       payloadBytes,
			})
			if err != nil {
				return nil, err
			}
			if response.ComponentType != componentTypeBlockchain {
				return nil, fmt.Errorf("unexpected component type in start response: %s", response.ComponentType)
			}
			for _, logLine := range response.AgentLogs {
				pretty := prettifyAgentLogLine(logLine)
				if pretty == "" {
					continue
				}
				testLogger.Info().Msgf("[agent] %s", pretty)
			}
			blockchainOutput, err := agent.DecodeFromTransport[blockchain.Output](response.Output)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to decode blockchain transport payload")
			}

			if err := rewriteRemoteBlockchainOutputForLocalAccess(ctx, testLogger, tunnelManager, idx, input, blockchainOutput, rewriteInternalForLocalNodes); err != nil {
				return nil, err
			}

			reconstructedBlockchain, err := blockchainFromOutput(testLogger, blockchainOutput)
			if err != nil {
				return nil, err
			}
			outputs[idx] = reconstructedBlockchain
		}
	}

	cldfBlockchains := make([]cldf_chain.BlockChain, 0, len(outputs))
	for _, db := range outputs {
		if db == nil {
			return nil, pkgerrors.New("blockchain output is nil")
		}
		chain, chainErr := db.ToCldfChain()
		if chainErr != nil {
			return nil, pkgerrors.Wrap(chainErr, "failed to create cldf chain from blockchain")
		}
		cldfBlockchains = append(cldfBlockchains, chain)
	}

	return &blockchains.DeployedBlockchains{
		Outputs:         outputs,
		CldfBlockChains: cldf_chain.NewBlockChainsFromSlice(cldfBlockchains),
	}, nil
}

func newEC2TunnelManager(testLogger zerolog.Logger) (tunnel.Manager, error) {
	if os.Getenv(envAgentMode) != "ec2" {
		return tunnel.NewNoopManager(), nil
	}

	instanceID := strings.TrimSpace(os.Getenv(envEC2InstanceID))
	if instanceID == "" {
		// Keep compatibility with pure manual-tunneling mode.
		return tunnel.NewNoopManager(), nil
	}

	return tunnel.NewManager(tunnel.NewSSMProvider(instanceID, ec2Region, testLogger)), nil
}

func NewEC2TunnelManager(testLogger zerolog.Logger) (tunnel.Manager, error) {
	return newEC2TunnelManager(testLogger)
}

func rewriteRemoteBlockchainOutputForLocalAccess(
	ctx context.Context,
	testLogger zerolog.Logger,
	tunnelManager tunnel.Manager,
	configuredIndex int,
	input *blockchain.Input,
	output *blockchain.Output,
	rewriteInternalForLocalNodes bool,
) error {
	if output == nil {
		return nil
	}

	componentID := tunnel.CanonicalComponentID(tunnel.KindBlockchain, configuredIndex, input.Type)
	adapter := adapters.NewBlockchainAdapter()

	refs, err := adapter.DescribeEndpoints(componentID, output)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to describe blockchain tunnel endpoints")
	}

	bindings, err := tunnelManager.Start(ctx, refs)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to start tunnels for blockchain output")
	}
	for _, binding := range bindings {
		testLogger.Info().
			Str("componentID", binding.ComponentID).
			Str("endpointName", binding.EndpointName).
			Str("originalURL", binding.OriginalURL).
			Str("localURL", binding.LocalURL).
			Msg("Established endpoint tunnel")
	}

	if err := adapter.RewriteWithBindings(output, bindings); err != nil {
		return pkgerrors.Wrap(err, "failed to rewrite blockchain output with local tunnel bindings")
	}
	if rewriteInternalForLocalNodes {
		if err := rewriteBlockchainInternalURLsForLocalNodes(output); err != nil {
			return pkgerrors.Wrap(err, "failed to rewrite blockchain internal urls for local node containers")
		}
	}

	return nil
}

func rewriteBlockchainInternalURLsForLocalNodes(output *blockchain.Output) error {
	if output == nil {
		return nil
	}

	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")
	for _, node := range output.Nodes {
		if node == nil {
			continue
		}

		if node.ExternalHTTPUrl != "" {
			internal, err := rewriteURLHost(node.ExternalHTTPUrl, dockerHost)
			if err != nil {
				return err
			}
			node.InternalHTTPUrl = internal
		}

		if node.ExternalWSUrl != "" {
			internal, err := rewriteURLHost(node.ExternalWSUrl, dockerHost)
			if err != nil {
				return err
			}
			node.InternalWSUrl = internal
		}
	}

	return nil
}

func rewriteURLHost(rawURL, host string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %q: %w", rawURL, err)
	}
	if parsed.Port() != "" {
		parsed.Host = net.JoinHostPort(host, parsed.Port())
		return parsed.String(), nil
	}
	parsed.Host = host
	return parsed.String(), nil
}

func remoteAgentError(code, message string) error {
	return fmt.Errorf("remote agent error (%s): %s", code, message)
}
