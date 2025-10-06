package don

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"google.golang.org/grpc/credentials/insecure"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	ctf_jd "github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func CreateJobs(ctx context.Context, testLogger zerolog.Logger, input cre.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, donMetadata := range input.DonTopology.ToDonMetadata() {
		if jobSpecs, ok := input.DonToJobSpecs[donMetadata.ID]; ok {
			createErr := jobs.Create(ctx, input.CldEnv.Offchain, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", donMetadata.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", donMetadata.ID)
		}
	}

	return nil
}

func AnyDonHasCapability(donMetadata []*cre.DonMetadata, capability cre.CapabilityFlag) bool {
	for _, don := range donMetadata {
		if flags.HasFlagForAnyChain(don.Flags, capability) {
			return true
		}
	}

	return false
}

func NodeNeedsAnyGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability) ||
		flags.HasFlag(nodeFlags, cre.VaultCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPActionCapability) ||
		flags.HasFlag(nodeFlags, cre.HTTPTriggerCapability)
}

func NodeNeedsWebAPIGateway(nodeFlags []cre.CapabilityFlag) bool {
	return flags.HasFlag(nodeFlags, cre.CustomComputeCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITriggerCapability) ||
		flags.HasFlag(nodeFlags, cre.WebAPITargetCapability)
}

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
		// a maximum of 1 bootstrap is supported due to environment constraints
		bootstrapNodeCount := 0
		if input.Topology.DonsMetadata.List()[idx].ContainsBootstrapNode() {
			bootstrapNodeCount = 1
		}

		nodeInfo, err := node.GetNodeInfo(nodeOutput.Output, nodeOutput.NodeSetName, input.Topology.DonsMetadata.List()[idx].ID, bootstrapNodeCount)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get node info")
		}
		allNodesInfo = append(allNodesInfo, nodeInfo...)

		supportedChains, schErr := findSupportedChainsForDON(input.Topology.DonsMetadata.List()[idx], input.BlockchainOutputs)
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
		hasEVMChainEnabled := slices.Contains(donMetadata.EVMChains(), bcOut.ChainID)
		hasSolanaWriteCapability := flags.HasFlagForAnyChain(donMetadata.Flags, cre.WriteSolanaCapability)
		chainIsSolana := strings.EqualFold(bcOut.BlockchainOutput.Family, chainselectors.FamilySolana)
		if !hasEVMChainEnabled && (!hasSolanaWriteCapability || !chainIsSolana) {
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
		for j, donNode := range topology.DonsMetadata.List()[i].NodesMetadata {
			// required for job proposals, because they need to include the ID of the node in Job Distributor
			donNode.Labels = append(donNode.Labels, &cre.Label{
				Key:   cre.NodeIDKey,
				Value: don.NodeIds()[j],
			})

			ocrSupportedFamilies := make([]string, 0)
			for family, key := range don.Nodes[j].ChainsOcr2KeyBundlesID {
				donNode.Labels = append(donNode.Labels, &cre.Label{
					Key:   node.CreateNodeOCR2KeyBundleIDKey(family),
					Value: key,
				})
				ocrSupportedFamilies = append(ocrSupportedFamilies, family)
			}

			donNode.Labels = append(donNode.Labels, &cre.Label{
				Key:   cre.NodeOCRFamiliesKey,
				Value: node.CreateNodeOCRFamiliesListValue(ocrSupportedFamilies),
			})
		}
	}

	return dons
}
