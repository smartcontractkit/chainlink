package environment

import (
	"bytes"
	"context"
	"errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const (
	componentTypeBlockchain = "blockchain"
	envLocalAgentURL        = "CRE_LOCAL_AGENT_URL"
	envAgentMode            = "CRE_AGENT_MODE"
)

type startComponentEnvelope struct {
	SchemaVersion string          `json:"schemaVersion"`
	Operation     string          `json:"operation"`
	Payload       json.RawMessage `json:"payload"`
}

type startBlockchainRequest struct {
	ComponentType string            `json:"componentType"`
	Blockchain    *blockchain.Input `json:"blockchain"`
}

type startBlockchainResult struct {
	BlockchainOutput map[string]any `json:"blockchainOutput"`
	AgentLogs        []string       `json:"agentLogs"`
	ErrorCode        string         `json:"errorCode"`
	Error            string         `json:"error"`
}

type componentClient interface {
	StartComponent(ctx context.Context, envelope startComponentEnvelope) (*startBlockchainResult, error)
}

type httpComponentClient struct {
	baseURL string
	client  *http.Client
}

func newHTTPComponentClient(baseURL string) *httpComponentClient {
	return &httpComponentClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 4 * time.Minute,
		},
	}
}

func (c *httpComponentClient) StartComponent(ctx context.Context, envelope startComponentEnvelope) (*startBlockchainResult, error) {
	body, err := json.Marshal(envelope)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to encode start component envelope")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/components/start", bytes.NewReader(body))
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create start component request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to execute start component request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to read start component response")
	}

	var startResp startBlockchainResult
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &startResp); err != nil {
			return nil, pkgerrors.Wrap(err, "failed to decode start component response")
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if startResp.Error != "" {
			if startResp.ErrorCode != "" {
				return nil, fmt.Errorf("%s: %s", startResp.ErrorCode, startResp.Error)
			}
			return nil, errors.New(startResp.Error)
		}
		return nil, fmt.Errorf("start component request failed with status %s: %s", resp.Status, string(respBody))
	}
	if startResp.Error != "" {
		if startResp.ErrorCode != "" {
			return nil, fmt.Errorf("%s: %s", startResp.ErrorCode, startResp.Error)
		}
		return nil, errors.New(startResp.Error)
	}

	return &startResp, nil
}

func newStartComponentClient() (componentClient, error) {
	if os.Getenv(envAgentMode) == "ec2" {
		return &ec2ComponentClient{}, nil
	}

	baseURL := os.Getenv(envLocalAgentURL)
	if baseURL == "" {
		return nil, fmt.Errorf("%s must be set for remote component startup", envLocalAgentURL)
	}
	return newHTTPComponentClient(baseURL), nil
}

type ec2ComponentClient struct{}

func (c *ec2ComponentClient) StartComponent(ctx context.Context, envelope startComponentEnvelope) (*startBlockchainResult, error) {
	return nil, errors.New("ec2 agent client is not implemented yet")
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
		startClient, err := newStartComponentClient()
		if err != nil {
			return nil, err
		}

		for _, idx := range remoteIdx {
			input := blockchainInputs[idx]
			if err := validatePhase2ARemoteBlockchainInput(input); err != nil {
				return nil, err
			}

			payload := startBlockchainRequest{
				ComponentType: componentTypeBlockchain,
				Blockchain:    input,
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to encode blockchain payload")
			}

			response, err := startClient.StartComponent(ctx, startComponentEnvelope{
				SchemaVersion: agent.SchemaVersionV1,
				Operation:     agent.OperationStartComponent,
				Payload:       payloadBytes,
			})
			if err != nil {
				return nil, err
			}
			for _, logLine := range response.AgentLogs {
				pretty := prettifyAgentLogLine(logLine)
				if pretty == "" {
					continue
				}
				testLogger.Info().Msgf("[agent] %s", pretty)
			}

			blockchainOutput, err := agent.DecodeFromTransport[blockchain.Output](response.BlockchainOutput)
			if err != nil {
				return nil, pkgerrors.Wrap(err, "failed to decode blockchain transport payload")
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
