package don

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"

	ctf_jd "github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
)

func LinkToJobDistributor(ctx context.Context, input *cre.LinkDonsToJDInput) (*cldf.Environment, []*devenv.DON, error) {
	if input == nil {
		return nil, nil, errors.New("input is nil")
	}
	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input validation failed")
	}

	dons := make([]*devenv.DON, len(input.NodeSetOutput))
	var allNodesInfo []devenv.NodeInfo

	for idx, nodeOutput := range input.NodeSetOutput {
		bootstrapNodes, err := libnode.FindManyWithLabel(input.Topology.DonsMetadata[idx].NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.BootstrapNode}, libnode.EqualLabels)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to find bootstrap nodes")
		}

		nodeInfo, err := libnode.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, input.Topology.DonsMetadata[idx].ID, len(bootstrapNodes))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)

		supportedChains, schErr := findSupportedChainsForDON(input.Topology.DonsMetadata[idx], input.BlockchainOutputs)
		if schErr != nil {
			return nil, nil, errors.Wrap(schErr, "failed to find supported chains for DON")
		}

		var regErr error
		dons[idx], regErr = configureJDForDON(ctx, nodeInfo, supportedChains, input.JdOutput)
		if regErr != nil {
			return nil, nil, fmt.Errorf("failed to configure JD for DON: %w", regErr)
		}
	}

	var nodeIDs []string
	for _, don := range dons {
		nodeIDs = append(nodeIDs, don.NodeIds()...)
	}

	dons = addOCRKeyLabelsToNodeMetadata(dons, input.Topology)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	jd, jdErr := devenv.NewJDClient(ctxWithTimeout, devenv.JDConfig{
		GRPC:     input.JdOutput.ExternalGRPCUrl,
		WSRPC:    input.JdOutput.InternalWSRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: allNodesInfo,
	})

	if jdErr != nil {
		return nil, nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	input.CldfEnvironment.Offchain = jd
	input.CldfEnvironment.NodeIDs = nodeIDs

	return input.CldfEnvironment, dons, nil
}

func configureJDForDON(ctx context.Context, nodeInfo []devenv.NodeInfo, supportedChains []devenv.ChainConfig, jdOutput *ctf_jd.Output) (*devenv.DON, error) {
	jdConfig := jd.JDConfig{
		GRPC:  jdOutput.ExternalGRPCUrl,
		WSRPC: jdOutput.InternalWSRPCUrl,
		Creds: insecure.NewCredentials(),
	}

	jdClient, jdErr := jd.NewJDClient(jdConfig)
	if jdErr != nil {
		return nil, errors.Wrap(jdErr, "failed to create JD client")
	}

	donJDClient := &devenv.JobDistributor{
		JobDistributor: jdClient,
	}

	don, regErr := devenv.NewRegisteredDON(ctx, nodeInfo, *donJDClient)
	if regErr != nil {
		return nil, fmt.Errorf("failed to create registered DON: %w", regErr)
	}

	if err := don.CreateSupportedChains(ctx, supportedChains, *donJDClient); err != nil {
		return nil, fmt.Errorf("failed to create supported chains: %w", err)
	}

	return don, nil
}

func findSupportedChainsForDON(donMetadata *cre.DonMetadata, blockchainOutputs []*cre.WrappedBlockchainOutput) ([]devenv.ChainConfig, error) {
	chains := make([]devenv.ChainConfig, 0)

	for chainSelector, bcOut := range blockchainOutputs {
		if len(donMetadata.SupportedChains) > 0 && !slices.Contains(donMetadata.SupportedChains, bcOut.ChainID) {
			continue
		}

		cfg, cfgErr := cre.ChainConfigFromWrapped(bcOut)
		if cfgErr != nil {
			return nil, errors.Wrapf(cfgErr, "failed to build chain config for chain selector %d", chainSelector)
		}

		chains = append(chains, cfg)
	}

	return chains, nil
}

func addOCRKeyLabelsToNodeMetadata(dons []*devenv.DON, topology *cre.Topology) []*devenv.DON {
	for i, don := range dons {
		for j, node := range topology.DonsMetadata[i].NodesMetadata {
			// required for job proposals, because they need to include the ID of the node in Job Distributor
			node.Labels = append(node.Labels, &cre.Label{
				Key:   libnode.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			ocrSupportedFamilies := make([]string, 0)
			for family, key := range don.Nodes[j].ChainsOcr2KeyBundlesID {
				node.Labels = append(node.Labels, &cre.Label{
					Key:   libnode.CreateNodeOCR2KeyBundleIDKey(family),
					Value: key,
				})
				ocrSupportedFamilies = append(ocrSupportedFamilies, family)
			}

			node.Labels = append(node.Labels, &cre.Label{
				Key:   libnode.NodeOCRFamiliesKey,
				Value: libnode.CreateNodeOCRFamiliesListValue(ocrSupportedFamilies),
			})
		}
	}

	return dons
}
