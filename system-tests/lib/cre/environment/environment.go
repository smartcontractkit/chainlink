package environment

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func BuildTopologyAndCLDEnvironment(lgr logger.Logger, nodeSetInput []*types.CapabilitiesAwareNodeSet, jdOutput *jd.Output, nodeSetOutput []*types.WrappedNodeOutput, blockchainOutput *blockchain.Output, sethClient *seth.Client) (*deployment.Environment, *types.DonTopology, error) {
	env, dons, err := buildChainlinkDeploymentEnv(lgr, jdOutput, nodeSetOutput, blockchainOutput, sethClient)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build chainlink deployment environment")
	}
	donTopology, err := libdon.BuildDONTopology(dons, nodeSetInput, nodeSetOutput)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build DON topology")
	}

	return env, donTopology, nil
}

func buildChainlinkDeploymentEnv(lgr logger.Logger, jdOutput *jd.Output, nodeSetOutput []*types.WrappedNodeOutput, blockchainOutput *blockchain.Output, sethClient *seth.Client) (*deployment.Environment, []*devenv.DON, error) {
	envs := make([]*deployment.Environment, len(nodeSetOutput))
	dons := make([]*devenv.DON, len(nodeSetOutput))

	for i, nodeOutput := range nodeSetOutput {
		// assume that each nodeset has only one bootstrap node
		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get node info")
		}

		jdConfig := devenv.JDConfig{
			GRPC:     jdOutput.HostGRPCUrl,
			WSRPC:    jdOutput.DockerWSRPCUrl,
			Creds:    insecure.NewCredentials(),
			NodeInfo: nodeInfo,
		}

		devenvConfig := devenv.EnvironmentConfig{
			JDConfig: jdConfig,
			Chains: []devenv.ChainConfig{
				{
					ChainID:   sethClient.Cfg.Network.ChainID,
					ChainName: sethClient.Cfg.Network.Name,
					ChainType: strings.ToUpper(blockchainOutput.Family),
					WSRPCs: []devenv.CribRPCs{{
						External: blockchainOutput.Nodes[0].HostWSUrl,
						Internal: blockchainOutput.Nodes[0].DockerInternalWSUrl,
					}},
					HTTPRPCs: []devenv.CribRPCs{{
						External: blockchainOutput.Nodes[0].HostHTTPUrl,
						Internal: blockchainOutput.Nodes[0].DockerInternalHTTPUrl,
					}},
					DeployerKey: sethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
				},
			},
		}

		env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to create environment")
		}

		envs[i] = env
		dons[i] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	// also, we don't care which instance of offchain client we use, because we have
	// only one instance of offchain client and we have just configured it to work
	// with nodes from all DONs
	return &deployment.Environment{
		Name:              envs[0].Name,
		Logger:            envs[0].Logger,
		ExistingAddresses: envs[0].ExistingAddresses,
		Chains:            envs[0].Chains,
		Offchain:          envs[0].Offchain,
		OCRSecrets:        envs[0].OCRSecrets,
		GetContext:        envs[0].GetContext,
		NodeIDs:           nodeIDs,
	}, dons, nil
}
