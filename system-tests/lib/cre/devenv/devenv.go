package environment

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func BuildFullCLDEnvironment(lgr logger.Logger, input *types.FullCLDEnvironmentInput, credentials credentials.TransportCredentials) (*types.FullCLDEnvironmentOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	envs := make([]*deployment.Environment, len(input.NodeSetOutput))
	dons := make([]*devenv.DON, len(input.NodeSetOutput))

	var allNodesInfo []devenv.NodeInfo
	chains := []devenv.ChainConfig{
		{
			ChainID:   input.SethClient.Cfg.Network.ChainID,
			ChainName: input.SethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(input.BlockchainOutput.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: input.BlockchainOutput.Nodes[0].HostWSUrl,
				Internal: input.BlockchainOutput.Nodes[0].DockerInternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: input.BlockchainOutput.Nodes[0].HostHTTPUrl,
				Internal: input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl,
			}},
			DeployerKey: input.SethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
		},
	}

	for i, nodeOutput := range input.NodeSetOutput {
		// assume that each nodeset has only one bootstrap node
		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, 1)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)

		// if DON has no capabilities we don't need to create chain configs (e.g. for gateway nodes)
		// we indicate to `devenv.NewEnvironment` that it should skip chain creation by passing an empty chain config
		if len(nodeOutput.Capabilities) == 0 {
			chains = []devenv.ChainConfig{}
		}

		jdConfig := devenv.JDConfig{
			GRPC:     input.JdOutput.HostGRPCUrl,
			WSRPC:    input.JdOutput.DockerWSRPCUrl,
			Creds:    credentials,
			NodeInfo: nodeInfo,
		}

		devenvConfig := devenv.EnvironmentConfig{
			JDConfig: jdConfig,
			Chains:   chains,
		}

		env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create environment")
		}

		envs[i] = env
		dons[i] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	for i, don := range dons {
		for j, node := range input.Topology.DonsMetadata[i].NodesMetadata {
			// both are required for job creation
			node.Labels = append(node.Labels, &types.Label{
				Key:   libnode.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			node.Labels = append(node.Labels, &types.Label{
				Key:   libnode.NodeOCR2KeyBundleIDKey,
				Value: don.Nodes[j].Ocr2KeyBundleID,
			})

			node.Labels = append(node.Labels, &types.Label{
				Key:   libnode.NodeOCR2KeyBundleIDKey,
				Value: don.Nodes[j].Ocr2KeyBundleID,
			})
		}
	}

	var jd deployment.OffchainClient
	var err error

	if len(input.NodeSetOutput) > 0 {
		// We create a new instance of JD client using `allNodesInfo` instead of `nodeInfo` to ensure that it can interact with all nodes.
		// Otherwise, JD would fail to accept job proposals for unknown nodes, even though it would still propose jobs to them. And that
		// would be happening silently, without any error messages, and we wouldn't know about it until much later.
		jd, err = devenv.NewJDClient(context.Background(), devenv.JDConfig{
			GRPC:     input.JdOutput.HostGRPCUrl,
			WSRPC:    input.JdOutput.DockerWSRPCUrl,
			Creds:    credentials,
			NodeInfo: allNodesInfo,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create JD client")
		}
	} else {
		jd = envs[0].Offchain
	}

	// we assume that all DONs run on the same chain and that there's only one chain
	output := &types.FullCLDEnvironmentOutput{
		Environment: &deployment.Environment{
			Name:              envs[0].Name,
			Logger:            envs[0].Logger,
			ExistingAddresses: input.ExistingAddresses,
			Chains:            envs[0].Chains,
			Offchain:          jd,
			OCRSecrets:        envs[0].OCRSecrets,
			GetContext:        envs[0].GetContext,
			NodeIDs:           nodeIDs,
		},
	}

	donTopology := &types.DonTopology{}
	donTopology.WorkflowDonID = input.Topology.WorkflowDONID

	for i, donMetadata := range input.Topology.DonsMetadata {
		donTopology.DonsWithMetadata = append(donTopology.DonsWithMetadata, &types.DonWithMetadata{
			DON:         dons[i],
			DonMetadata: donMetadata,
		})
	}

	output.DonTopology = donTopology

	return output, nil
}
