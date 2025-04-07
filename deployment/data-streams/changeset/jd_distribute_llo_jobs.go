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
)

var _ deployment.ChangeSetV2[CsDistributeLLOJobSpecsConfig] = CsDistributeLLOJobSpecs{}

const (
	lloJobMaxTaskDuration             = time.Second
	contractConfigTrackerPollInterval = time.Second
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
		})

	// nils will be filled out later with n-specific values:
	lloSpec := &jobs.LLOJobSpec{
		Base: jobs.Base{
			Name:          fmt.Sprintf("%s | %d", cfg.Filter.DONName, cfg.Filter.DONID),
			Type:          jobs.JobSpecTypeLLO,
			SchemaVersion: 1,
			ExternalJobID: uuid.New(),
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

	oracleNodes, err := jd.FetchDONOraclesFromJD(ctx, e.Offchain, cfg.Filter)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get workflow don nodes: %w", err)
	}

	nodeConfigMap, err := chainConfigs(ctx, e, chainID, oracleNodes)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get node chain configs: %w", err)
	}

	var proposals []*jobv1.ProposeJobRequest
	for _, n := range oracleNodes {
		lloSpec.TransmitterID = n.GetPublicKey() // CSAKey
		lloSpec.OCRKeyBundleID = &nodeConfigMap[n.Id].OcrKeyBundle.BundleId

		p2p := nodeConfigMap[n.Id].P2PKeyBundle
		lloSpec.P2PV2Bootstrappers = []string{fmt.Sprintf("%s:%s", p2p.GetPublicKey(), p2p.GetPeerId())}
		lloSpec.PluginConfig.Servers = cfg.Servers

		renderedSpec, err := lloSpec.MarshalTOML()
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal llo spec: %w", err)
		}

		proposals = append(proposals, &jobv1.ProposeJobRequest{
			NodeId: n.Id,
			Spec:   string(renderedSpec),
			Labels: labels,
		})
	}
	proposedJobs, err := proposeAllOrNothing(ctx, e.Offchain, proposals)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to propose all oracle jobs: %w", err)
	}

	return deployment.ChangesetOutput{
		Jobs: proposedJobs,
	}, nil
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

	return nil
}
