package environment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"

	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/agent"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

type startComponentEnvelope = agent.StartComponentEnvelope
type startComponentRequest = agent.StartComponentPayload

func blockchainFromOutput(testLogger zerolog.Logger, output *blockchain.Output) (blockchains.Blockchain, error) {
	if output == nil {
		return nil, pkgerrors.New("blockchain output is nil")
	}

	if output.Type != blockchain.TypeAnvil {
		return nil, fmt.Errorf("remote blockchain reconstruction supports only %s in phase 2A, got %s", blockchain.TypeAnvil, output.Type)
	}

	return evm.FromOutput(testLogger, output)
}

func validateRemoteBlockchainInput(input *blockchain.Input) error {
	if input == nil {
		return errors.New("blockchain input is nil")
	}
	if input.Type != blockchain.TypeAnvil {
		return fmt.Errorf("remote target supports only %s, got %s", blockchain.TypeAnvil, input.Type)
	}
	return nil
}

func startBlockchainsWithTargets(
	ctx context.Context,
	testLogger zerolog.Logger,
	configuredBlockchains []*config.Blockchain,
	deployers map[blockchain.ChainFamily]blockchains.Deployer,
	remoteRuntime *resolvedRemoteRuntime,
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
		if configuredBlockchain.Placement == config.PlacementRemote {
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
		if remoteRuntime == nil {
			return nil, errors.New("remote runtime is required when starting remote blockchains")
		}
		startClient, err := newRemoteComponentClient(remoteRuntime)
		if err != nil {
			return nil, err
		}

		for _, idx := range remoteIdx {
			input := blockchainInputs[idx]
			configured := configuredBlockchains[idx]
			if err := validateRemoteBlockchainInput(input); err != nil {
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

			if err := rewriteRemoteBlockchainOutputForLocalAccess(blockchainOutput, remoteRuntime.EC2HostIP, rewriteInternalForLocalNodes); err != nil {
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

func rewriteRemoteBlockchainOutputForLocalAccess(
	output *blockchain.Output,
	ec2HostIP string,
	rewriteInternalForLocalNodes bool,
) error {
	_ = rewriteInternalForLocalNodes // direct mode keeps internal URLs unchanged
	if output == nil {
		return nil
	}
	return rewriteRemoteBlockchainOutputForDirectAccess(output, ec2HostIP)
}

func rewriteRemoteBlockchainOutputForDirectAccess(output *blockchain.Output, ec2HostIP string) error {
	if output == nil {
		return nil
	}
	for _, node := range output.Nodes {
		if node == nil {
			continue
		}
		if node.ExternalHTTPUrl != "" {
			rewritten, err := rewriteURLHost(node.ExternalHTTPUrl, ec2HostIP)
			if err != nil {
				return err
			}
			node.ExternalHTTPUrl = rewritten
		}
		if node.ExternalWSUrl != "" {
			rewritten, err := rewriteURLHost(node.ExternalWSUrl, ec2HostIP)
			if err != nil {
				return err
			}
			node.ExternalWSUrl = rewritten
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
