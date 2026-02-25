package environment

import (
	"context"
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
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/tron"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func blockchainFromOutput(testLogger zerolog.Logger, input *blockchain.Input, output *blockchain.Output) (blockchains.Blockchain, error) {
	if output == nil {
		return nil, pkgerrors.New("blockchain output is nil")
	}

	switch output.Type {
	case blockchain.TypeAnvil:
		return evm.From(testLogger, output)
	case blockchain.TypeTron:
		return tron.From(testLogger, output)
	case blockchain.TypeSolana:
		return solana.From(input, output)
	default:
		return nil, fmt.Errorf("unsupported blockchain type for reconstruction: %s", output.Type)
	}
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

func startBlockchains(
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

	outputs := make([]blockchains.Blockchain, len(configuredBlockchains))
	for idx, configured := range configuredBlockchains {
		input := blockchainInputs[idx]
		var deployedOutput *blockchain.Output

		if configured.Placement == config.PlacementRemote {
			if err := validateRemoteBlockchainInput(input); err != nil {
				return nil, err
			}
			payload := agent.StartComponentPayload{
				ComponentType: componentTypeBlockchain,
				Blockchain:    input,
				ReusePolicy:   string(configured.RemoteStartPolicy),
			}
			deployedOutput, err = startRemoteComponent[blockchain.Output](
				ctx,
				testLogger,
				remoteRuntime.Client,
				payload,
				componentTypeBlockchain,
			)
			if err != nil {
				return nil, err
			}
			if rewriteInternalForLocalNodes {
				// direct mode keeps internal URLs unchanged
			}
			if err := rewriteRemoteBlockchainOutputForDirectAccess(deployedOutput, remoteRuntime.EC2HostIP); err != nil {
				return nil, err
			}
		} else {
			deployedOutput, err = blockchains.StartChain(ctx, deployers, input)
			if err != nil {
				return nil, err
			}
		}

		reconstructedBlockchain, err := blockchainFromOutput(testLogger, input, deployedOutput)
		if err != nil {
			return nil, err
		}
		outputs[idx] = reconstructedBlockchain
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
