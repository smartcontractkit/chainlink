package environment

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

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
	chains := make([]devenv.ChainConfig, 0)
	for chainSelector, bcOut := range input.BlockchainOutputs {
		cID, err := strconv.ParseUint(bcOut.ChainID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chain ID: %w", err)
		}

		sethClient, ok := input.SethClients[chainSelector]
		if !ok {
			return nil, fmt.Errorf("seth client not found for chain selector: %d", chainSelector)
		}

		chains = append(chains, devenv.ChainConfig{
			ChainID:   cID,
			ChainName: sethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(bcOut.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: bcOut.Nodes[0].ExternalWSUrl,
				Internal: bcOut.Nodes[0].InternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: bcOut.Nodes[0].ExternalHTTPUrl,
				Internal: bcOut.Nodes[0].InternalHTTPUrl,
			}},
			DeployerKey: sethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
		})
	}

	for idx, nodeOutput := range input.NodeSetOutput {
		// check how many bootstrap nodes we have in each DON
		bootstrapNodes, err := libnode.FindManyWithLabel(input.Topology.DonsMetadata[idx].NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.BootstrapNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}

		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, len(bootstrapNodes))
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
			GRPC:     input.JdOutput.ExternalGRPCUrl,
			WSRPC:    input.JdOutput.InternalWSRPCUrl,
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

		envs[idx] = env
		dons[idx] = don
	}

	var nodeIDs []string
	for _, env := range envs {
		nodeIDs = append(nodeIDs, env.NodeIDs...)
	}

	for i, don := range dons {
		for j, node := range input.Topology.DonsMetadata[i].NodesMetadata {
			// required for job proposals, because they need to include the ID of the node in Job Distributor
			node.Labels = append(node.Labels, &types.Label{
				Key:   libnode.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			// required for OCR2/3 job specs
			node.Labels = append(node.Labels, &types.Label{
				Key:   libnode.NodeOCR2KeyBundleIDKey,
				Value: don.Nodes[j].Ocr2KeyBundleID,
			})
		}
	}

	var jd cldf.OffchainClient
	var err error

	if len(input.NodeSetOutput) > 0 {
		// We create a new instance of JD client using `allNodesInfo` instead of `nodeInfo` to ensure that it can interact with all nodes.
		// Otherwise, JD would fail to accept job proposals for unknown nodes, even though it would still propose jobs to them. And that
		// would be happening silently, without any error messages, and we wouldn't know about it until much later.
		jd, err = devenv.NewJDClient(context.Background(), devenv.JDConfig{
			GRPC:     input.JdOutput.ExternalGRPCUrl,
			WSRPC:    input.JdOutput.InternalWSRPCUrl,
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
	donTopology.HomeChainSelector = input.Topology.HomeChainSelector

	for i, donMetadata := range input.Topology.DonsMetadata {
		donTopology.DonsWithMetadata = append(donTopology.DonsWithMetadata, &types.DonWithMetadata{
			DON:         dons[i],
			DonMetadata: donMetadata,
		})
	}

	output.DonTopology = donTopology
	output.DonTopology.GatewayConnectorOutput = input.Topology.GatewayConnectorOutput

	return output, nil
}
