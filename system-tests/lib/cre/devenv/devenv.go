package environment

import (
	"context"
	"slices"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_offchain "github.com/smartcontractkit/chainlink-deployments-framework/offchain"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

func BuildFullCLDEnvironment(ctx context.Context, lgr logger.Logger, input *cre.FullCLDEnvironmentInput, credentials credentials.TransportCredentials) (*cre.FullCLDEnvironmentOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	envs := make([]*cldf.Environment, len(input.NodeSetOutput))
	dons := make([]*devenv.DON, len(input.NodeSetOutput))

	var allNodesInfo []devenv.NodeInfo
	for idx, nodeOutput := range input.NodeSetOutput {
		// check how many bootstrap nodes we have in each DON
		bootstrapNodes, err := libnode.FindManyWithLabel(input.Topology.DonsMetadata[idx].NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.BootstrapNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}

		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, input.Topology.DonsMetadata[idx].ID, len(bootstrapNodes))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)
		chains := make([]devenv.ChainConfig, 0)

		// for each nodeSet create only chains that the DON supports
		for chainSelector, bcOut := range input.BlockchainOutputs {
			if len(input.Topology.DonsMetadata[idx].SupportedChains) > 0 && !slices.Contains(input.Topology.DonsMetadata[idx].SupportedChains, bcOut.ChainID) {
				continue
			}

			cfg, cfgErr := cre.ChainConfigFromWrapped(bcOut)
			if cfgErr != nil {
				return nil, errors.Wrapf(err, "failed to build chain config for chain selector %d", chainSelector)
			}

			chains = append(chains, cfg)
		}

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

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Minute)
		env, don, envErr := devenv.NewEnvironment(func() context.Context {
			return ctxWithTimeout
		}, lgr, devenvConfig)
		if envErr != nil {
			cancel()
			return nil, errors.Wrap(envErr, "failed to create environment")
		}
		cancel()

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
			node.Labels = append(node.Labels, &cre.Label{
				Key:   libnode.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			// required for OCR2/3 job specs
			node.Labels = append(node.Labels, &cre.Label{
				Key:   libnode.NodeOCR2KeyBundleIDKey,
				Value: don.Nodes[j].Ocr2KeyBundleID,
			})
		}
	}

	var jd cldf_offchain.Client

	if len(input.NodeSetOutput) > 0 {
		// We create a new instance of JD client using `allNodesInfo` instead of `nodeInfo` to ensure that it can interact with all nodes.
		// Otherwise, JD would fail to accept job proposals for unknown nodes, even though it would still propose jobs to them. And that
		// would be happening silently, without any error messages, and we wouldn't know about it until much later.
		var jdErr error
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
		jd, jdErr = devenv.NewJDClient(ctxWithTimeout, devenv.JDConfig{
			GRPC:     input.JdOutput.ExternalGRPCUrl,
			WSRPC:    input.JdOutput.InternalWSRPCUrl,
			Creds:    credentials,
			NodeInfo: allNodesInfo,
		})
		if jdErr != nil {
			cancel()
			return nil, errors.Wrap(jdErr, "failed to create JD client")
		}
		cancel()
	} else {
		jd = envs[0].Offchain
	}

	// create chains for all chains that are supported by any of the DONs, so that changeset can be applied to all chains
	allChainsConfigs := make([]devenv.ChainConfig, 0)
	for chainSelector, bcOut := range input.BlockchainOutputs {
		cfg, cfgErr := cre.ChainConfigFromWrapped(bcOut)
		if cfgErr != nil {
			return nil, errors.Wrapf(cfgErr, "failed to build chain config for chain selector %d", chainSelector)
		}

		allChainsConfigs = append(allChainsConfigs, cfg)
	}

	blockChains, allChainsErr := devenv.NewChains(lgr, allChainsConfigs)
	if allChainsErr != nil {
		return nil, errors.Wrap(allChainsErr, "failed to create chains")
	}

	// we take stateless fields from the first environment, because they are not environment specific
	output := &cre.FullCLDEnvironmentOutput{
		Environment: &cldf.Environment{
			Name:              envs[0].Name,
			Logger:            envs[0].Logger,
			ExistingAddresses: input.ExistingAddresses,
			DataStore:         input.Datastore,
			Offchain:          jd,
			OCRSecrets:        envs[0].OCRSecrets,
			GetContext:        envs[0].GetContext,
			NodeIDs:           nodeIDs,
			BlockChains:       blockChains,
			OperationsBundle:  input.OperationsBundle,
		},
	}

	donTopology := &cre.DonTopology{}
	donTopology.WorkflowDonID = input.Topology.WorkflowDONID
	donTopology.HomeChainSelector = input.Topology.HomeChainSelector
	donTopology.CapabilitiesPeeringData = input.Topology.CapabilitiesPeeringData
	donTopology.OCRPeeringData = input.Topology.OCRPeeringData

	for i, donMetadata := range input.Topology.DonsMetadata {
		donTopology.DonsWithMetadata = append(donTopology.DonsWithMetadata, &cre.DonWithMetadata{
			DON:         dons[i],
			DonMetadata: donMetadata,
		})
	}

	output.DonTopology = donTopology
	output.DonTopology.GatewayConnectorOutput = input.Topology.GatewayConnectorOutput

	return output, nil
}
