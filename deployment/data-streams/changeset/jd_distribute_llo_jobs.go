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

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

var _ deployment.ChangeSetV2[CsDistributeLLOJobSpecsConfig] = CsDistributeLLOJobSpecs{}

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

	// Servers is a list of Data Engine Producer endpoints, where the key is the server URL and the value is its public key.
	//
	// Example:
	// 	"mercury-pipeline-testnet-producer.stage-2.cldev.cloud:1340": "11a34b5187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c"
	Servers map[string]string

	// NodeNames specifies on which nodes to distribute the job specs.
	NodeNames []string
}

type CsDistributeLLOJobSpecs struct{}

func (CsDistributeLLOJobSpecs) Apply(e deployment.Environment, cfg CsDistributeLLOJobSpecsConfig) (deployment.ChangesetOutput, error) {
	ctx, cancel := context.WithTimeout(e.GetContext(), defaultJobSpecsTimeout)
	defer cancel()

	chainID, _, err := chainAndAddresses(e, cfg.ChainSelectorEVM)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// Add a label to the job spec to identify the related DON
	labels := append([]*ptypes.Label(nil),
		&ptypes.Label{
			Key: utils.DonIdentifier(cfg.Filter.DONID, cfg.Filter.DONName),
		},
		&ptypes.Label{
			Key:   devenv.LabelJobTypeKey,
			Value: pointer.To(devenv.LabelJobTypeValueLLO),
		},
	)

	bootstrapProposals, err := generateBootstrapProposals(ctx, e, cfg, chainID, labels)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate bootstrap proposals: %w", err)
	}
	oracleProposals, err := generateOracleProposals(ctx, e, cfg, chainID, labels)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate oracle proposals: %w", err)
	}

	proposedJobs, err := proposeAllOrNothing(ctx, e.Offchain, append(bootstrapProposals, oracleProposals...))
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose all jobs: %w", err)
	}

	return deployment.ChangesetOutput{
		Jobs: proposedJobs,
	}, nil
}

func generateBootstrapProposals(ctx context.Context, e deployment.Environment, cfg CsDistributeLLOJobSpecsConfig, chainID string, labels []*ptypes.Label) ([]*jobv1.ProposeJobRequest, error) {
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
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get externalJobID: %w", err)
		}

		bootstrapSpec := jobs.NewBootstrapSpec(
			cfg.ConfiguratorAddress,
			cfg.Filter.DONID,
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

func generateOracleProposals(ctx context.Context, e deployment.Environment, cfg CsDistributeLLOJobSpecsConfig, chainID string, labels []*ptypes.Label) ([]*jobv1.ProposeJobRequest, error) {
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

	bootstrapMultiaddr, err := getBootstrapMultiAddr(ctx, e, cfg)
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
func chainConfigs(ctx context.Context, e deployment.Environment, chainID string, nodes []*node.Node) (map[string]*node.OCR2Config, error) {
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
func getBootstrapMultiAddr(ctx context.Context, e deployment.Environment, cfg CsDistributeLLOJobSpecsConfig) (string, error) {
	// Get all bootstrap nodes for this DON.
	// We fetch these with a custom filter because the filter in the config defines which nodes need to be sent jobs
	// and this might not cover any bootstrap nodes.
	respBoots, err := e.Offchain.ListNodes(ctx, &node.ListNodesRequest{
		Filter: &node.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{
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
	resp, err := e.Offchain.ListNodeChainConfigs(ctx, &node.ListNodeChainConfigsRequest{
		Filter: &node.ListNodeChainConfigsRequest_Filter{
			NodeIds: []string{respBoots.Nodes[0].Id},
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

func (f CsDistributeLLOJobSpecs) VerifyPreconditions(_ deployment.Environment, config CsDistributeLLOJobSpecsConfig) error {
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
