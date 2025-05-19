package changeset

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

var _ cldf.ChangeSetV2[CsDistributeLLOJobSpecsConfig] = CsDistributeLLOJobSpecs{}

const (
	lloJobMaxTaskDuration             = jobs.TOMLDuration(time.Second)
	contractConfigTrackerPollInterval = jobs.TOMLDuration(time.Second)
)

type CsDistributeLLOJobSpecsConfig struct {
	ChainSelectorEVM uint64
	Filter           *jd.ListFilter

	FromBlock  uint64
	ConfigMode string // e.g. bluegreen

	ChannelConfigStoreAddr      common.Address
	ChannelConfigStoreFromBlock uint64
	ConfiguratorAddress         string
	Labels                      []*ptypes.Label

	// Servers is a list of Data Engine Producer endpoints, where the key is the server URL and the value is its public key.
	//
	// Example:
	// 	"mercury-pipeline-testnet-producer.stage-2.cldev.cloud:1340": "11a34b5187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c"
	Servers map[string]string

	// NodeNames specifies on which nodes to distribute the job specs.
	NodeNames []string
}

type CsDistributeLLOJobSpecs struct{}

func (CsDistributeLLOJobSpecs) Apply(e cldf.Environment, cfg CsDistributeLLOJobSpecsConfig) (cldf.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultJobSpecsTimeout)
	defer cancel()

	chainID, _, err := chainAndAddresses(e, cfg.ChainSelectorEVM)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Add a label to the job spec to identify the related DON
	cfg.Labels = append(cfg.Labels,
		&ptypes.Label{
			Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
		},
		&ptypes.Label{
			Key:   devenv.LabelJobTypeKey,
			Value: pointer.To(devenv.LabelJobTypeValueLLO),
		},
	)

	bootstrapProposals, err := generateBootstrapProposals(ctx, e, cfg, chainID, cfg.Labels)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate bootstrap proposals: %w", err)
	}
	// These will be empty when we send only oracle jobs. In that case we'll fetch the bootstrappers by the don
	// identifier label.
	boostrapNodeIDs := make([]string, 0, len(bootstrapProposals))
	for _, p := range bootstrapProposals {
		boostrapNodeIDs = append(boostrapNodeIDs, p.NodeId)
	}
	oracleProposals, err := generateOracleProposals(ctx, e, cfg, chainID, cfg.Labels, boostrapNodeIDs)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate oracle proposals: %w", err)
	}
	allProposals := append(bootstrapProposals, oracleProposals...) //nolint: gocritic // ignore a silly rule
	proposedJobs, err := proposeAllOrNothing(ctx, e.Offchain, allProposals)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose all jobs: %w", err)
	}

	err = labelNodesForProposals(e.GetContext(), e.Offchain, allProposals, utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to label nodes for proposals: %w", err)
	}

	return cldf.ChangesetOutput{
		Jobs: proposedJobs,
	}, nil
}

func generateBootstrapProposals(
	ctx context.Context,
	e cldf.Environment,
	cfg CsDistributeLLOJobSpecsConfig,
	chainID string,
	labels []*ptypes.Label,
) ([]*jobv1.ProposeJobRequest, error) {
	bootstrapNodes, err := jd.FetchDONBootstrappersFromJD(ctx, e.Offchain, cfg.Filter, cfg.NodeNames)
	if err != nil {
		return nil, fmt.Errorf("failed to get bootstrap nodes: %w", err)
	}

	localLabels := append(labels, //nolint: gocritic // obvious and readable locally modified copy of labels
		&ptypes.Label{
			Key:   devenv.LabelNodeTypeKey,
			Value: pointer.To(devenv.LabelNodeTypeValueBootstrap),
		},
	)

	var proposals []*jobv1.ProposeJobRequest
	for _, btNode := range bootstrapNodes {
		externalJobID, err := fetchExternalJobID(e, btNode.Id, []*ptypes.Selector{
			{
				Key:   devenv.LabelJobTypeKey,
				Value: pointer.To(devenv.LabelJobTypeValueLLO),
				Op:    ptypes.SelectorOp_EQ,
			},
			{
				Key:   devenv.LabelNodeTypeKey,
				Value: pointer.To(devenv.LabelNodeTypeValueBootstrap),
				Op:    ptypes.SelectorOp_EQ,
			},
			{
				Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
				Op:  ptypes.SelectorOp_EXIST,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get externalJobID: %w", err)
		}

		bootstrapSpec := jobs.NewBootstrapSpec(
			cfg.ConfiguratorAddress,
			cfg.Filter.DONID,
			cfg.Filter.DONName,
			jobs.RelayTypeEVM,
			jobs.RelayConfig{
				ChainID: chainID,
			},
			externalJobID,
		)

		renderedSpec, err := bootstrapSpec.MarshalTOML()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal bootstrap spec: %w", err)
		}

		proposals = append(proposals, &jobv1.ProposeJobRequest{
			NodeId: btNode.Id,
			Spec:   string(renderedSpec),
			Labels: localLabels,
		})
	}

	return proposals, nil
}

func generateOracleProposals(
	ctx context.Context,
	e cldf.Environment,
	cfg CsDistributeLLOJobSpecsConfig,
	chainID string,
	labels []*ptypes.Label,
	boostrapNodeIDs []string,
) ([]*jobv1.ProposeJobRequest, error) {
	// nils will be filled out later with n-specific values:
	lloSpec := &jobs.LLOJobSpec{
		Base: jobs.Base{
			Name:          fmt.Sprintf("%s | %d", cfg.Filter.DONName, cfg.Filter.DONID),
			Type:          jobs.JobSpecTypeLLO,
			SchemaVersion: 1,
			// We intentionally do not set ExternalJobID here - we'll set it separately for each node.
		},
		ContractID:                        cfg.ConfiguratorAddress,
		P2PV2Bootstrappers:                nil,
		OCRKeyBundleID:                    nil,
		MaxTaskDuration:                   lloJobMaxTaskDuration,
		ContractConfigTrackerPollInterval: contractConfigTrackerPollInterval,
		Relay:                             jobs.RelayTypeEVM,
		PluginType:                        jobs.PluginTypeLLO,
		RelayConfig: jobs.RelayConfigLLO{
			ChainID:       chainID,
			FromBlock:     cfg.FromBlock,
			LLOConfigMode: cfg.ConfigMode,
			LLODonID:      cfg.Filter.DONID,
		},
		PluginConfig: jobs.PluginConfigLLO{
			ChannelDefinitionsContractAddress:   cfg.ChannelConfigStoreAddr.Hex(),
			ChannelDefinitionsContractFromBlock: cfg.ChannelConfigStoreFromBlock,
			DonID:                               cfg.Filter.DONID,
			Servers:                             nil,
		},
	}

	oracleNodes, err := jd.FetchDONOraclesFromJD(ctx, e.Offchain, cfg.Filter, cfg.NodeNames)
	if err != nil {
		return nil, fmt.Errorf("failed to get oracle nodes: %w", err)
	}

	nodeConfigMap, err := chainConfigs(ctx, e, chainID, oracleNodes)
	if err != nil {
		return nil, fmt.Errorf("failed to get node chain configs: %w", err)
	}

	bootstrapMultiaddr, err := getBootstrapMultiAddr(ctx, e, cfg, boostrapNodeIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get bootstrap bootstrapMultiaddr: %w", err)
	}

	localLabels := append(labels, //nolint: gocritic // obvious and readable locally modified copy of labels
		&ptypes.Label{
			Key:   devenv.LabelNodeTypeKey,
			Value: pointer.To(devenv.LabelNodeTypeValuePlugin),
		},
	)

	var proposals []*jobv1.ProposeJobRequest
	for _, n := range oracleNodes {
		externalJobID, err := fetchExternalJobID(e, n.Id, []*ptypes.Selector{
			{
				Key:   devenv.LabelJobTypeKey,
				Value: pointer.To(devenv.LabelJobTypeValueLLO),
				Op:    ptypes.SelectorOp_EQ,
			},
			{
				Key:   devenv.LabelNodeTypeKey,
				Value: pointer.To(devenv.LabelNodeTypeValuePlugin),
				Op:    ptypes.SelectorOp_EQ,
			},
			{
				Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
				Op:  ptypes.SelectorOp_EXIST,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get externalJobID: %w", err)
		}
		if externalJobID == uuid.Nil {
			externalJobID = uuid.New()
		}

		lloSpec.Base.ExternalJobID = externalJobID
		lloSpec.TransmitterID = n.GetPublicKey() // CSAKey
		lloSpec.OCRKeyBundleID = &nodeConfigMap[n.Id].OcrKeyBundle.BundleId

		lloSpec.P2PV2Bootstrappers = []string{bootstrapMultiaddr}
		lloSpec.PluginConfig.Servers = cfg.Servers

		renderedSpec, err := lloSpec.MarshalTOML()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal llo spec: %w", err)
		}

		proposals = append(proposals, &jobv1.ProposeJobRequest{
			NodeId: n.Id,
			Spec:   string(renderedSpec),
			Labels: localLabels,
		})
	}

	return proposals, nil
}

// chainConfigs returns a map of node IDs to their chain configs for the given chain ID.
func chainConfigs(ctx context.Context, e cldf.Environment, chainID string, nodes []*node.Node) (map[string]*node.OCR2Config, error) {
	nodeConfigMap := make(map[string]*node.OCR2Config)
	for _, n := range nodes {
		ncf, err := e.Offchain.ListNodeChainConfigs(ctx,
			&node.ListNodeChainConfigsRequest{
				Filter: &node.ListNodeChainConfigsRequest_Filter{
					NodeIds: []string{n.Id},
				},
			})
		if err != nil {
			return nil, fmt.Errorf("failed to get chain config: %w", err)
		}
		for _, nc := range ncf.GetChainConfigs() {
			if nc.GetChain().Id == chainID {
				nodeConfigMap[nc.GetNodeId()] = nc.GetOcr2Config()
			}
		}
	}

	return nodeConfigMap, nil
}

// getBootstrapMultiAddr fetches the bootstrap node from Job Distributor and returns its multiaddr.
// If boostrapNodeIDs is empty, it will return the first bootstrap node found for this DON.
func getBootstrapMultiAddr(ctx context.Context, e cldf.Environment, cfg CsDistributeLLOJobSpecsConfig, boostrapNodeIDs []string) (string, error) {
	if len(boostrapNodeIDs) == 0 {
		// Get all bootstrap nodes for this DON.
		// We fetch these with a custom filter because the filter in the config defines which nodes need to be sent jobs
		// and this might not cover any bootstrap nodes.
		respBoots, err := e.Offchain.ListNodes(ctx, &node.ListNodesRequest{
			Filter: &node.ListNodesRequest_Filter{
				Selectors: []*ptypes.Selector{
					// We can afford to filter by DonIdentifier here because if the caller didn't provide any bootstrap node IDs,
					// then they are updating an existing job spec and the bootstrap nodes are already labeled with the DON ID.
					{
						Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
						Op:  ptypes.SelectorOp_EXIST,
					},
					{
						Key:   devenv.LabelNodeTypeKey,
						Op:    ptypes.SelectorOp_EQ,
						Value: pointer.To(devenv.LabelNodeTypeValueBootstrap),
					},
					{
						Key:   devenv.LabelEnvironmentKey,
						Op:    ptypes.SelectorOp_EQ,
						Value: &cfg.Filter.EnvLabel,
					},
					{
						Key:   devenv.LabelProductKey,
						Op:    ptypes.SelectorOp_EQ,
						Value: pointer.To(utils.ProductLabel),
					},
				},
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to list bootstrap nodes for DON %d - %s: %w", cfg.Filter.DONID, cfg.Filter.DONName, err)
		}
		if len(respBoots.Nodes) == 0 {
			return "", errors.New("no bootstrap nodes found")
		}
		for _, n := range respBoots.Nodes {
			boostrapNodeIDs = append(boostrapNodeIDs, n.Id)
		}
	}

	resp, err := e.Offchain.ListNodeChainConfigs(ctx, &node.ListNodeChainConfigsRequest{
		Filter: &node.ListNodeChainConfigsRequest_Filter{
			NodeIds: boostrapNodeIDs,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get chain config: %w", err)
	}
	if resp == nil || len(resp.ChainConfigs) == 0 {
		return "", errors.New("no chain configs found")
	}
	return resp.ChainConfigs[0].Ocr2Config.Multiaddr, nil
}

func (f CsDistributeLLOJobSpecs) VerifyPreconditions(_ cldf.Environment, config CsDistributeLLOJobSpecsConfig) error {
	if config.ChainSelectorEVM == 0 {
		return errors.New("chain selector is required")
	}
	if config.Filter == nil {
		return errors.New("filter is required")
	}
	if config.ConfigMode != "bluegreen" {
		return fmt.Errorf("invalid config mode: %s", config.ConfigMode)
	}
	if config.ChannelConfigStoreAddr == (common.Address{}) {
		return errors.New("channel config store address is required")
	}
	if len(config.Servers) == 0 {
		return errors.New("servers map is required")
	}
	if len(config.NodeNames) == 0 {
		return errors.New("node names are required")
	}

	return nil
}

// labelNodesForProposals adds a DON Identifier label to the nodes for the given proposals.
func labelNodesForProposals(ctx context.Context, jd cldf.OffchainClient, props []*jobv1.ProposeJobRequest, donIdentifier string) error {
	for _, p := range props {
		nodeResp, err := jd.GetNode(ctx, &node.GetNodeRequest{Id: p.NodeId})
		if err != nil {
			return fmt.Errorf("failed to get node %s: %w", p.NodeId, err)
		}
		newLabels := append(nodeResp.Node.Labels, &ptypes.Label{ //nolint: gocritic // local copy
			Key: donIdentifier,
		})

		_, err = jd.UpdateNode(ctx, &node.UpdateNodeRequest{
			Id:        p.NodeId,
			Name:      nodeResp.Node.Name,
			PublicKey: nodeResp.Node.PublicKey,
			Labels:    newLabels,
		})
		if err != nil {
			return fmt.Errorf("failed to label node %s: %w", p.NodeId, err)
		}
	}
	return nil
}
