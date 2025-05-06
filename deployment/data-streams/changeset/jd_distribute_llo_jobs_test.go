package changeset

import (
	"context"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
)

func TestDistributeLLOJobSpecs(t *testing.T) {
	t.Parallel()

	const donID = 1
	const donName = "don"
	const envName = "env"

	env := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      false,
		ShouldDeployLinkToken: false,
		NumNodes:              2,
		NumBootstrapNodes:     1,
		NodeLabels: []*ptypes.Label{
			{
				Key:   devenv.LabelProductKey,
				Value: pointer.To(utils.ProductLabel),
			},
			{
				Key:   devenv.LabelEnvironmentKey,
				Value: pointer.To(envName),
			},
			{
				Key: utils.DonIdentifier(donID, donName),
			},
		},
	}).Environment

	// Collect the names of the nodes.
	bootstrapNodeNames := make([]string, 0, 1)
	oracleNodeNames := make([]string, 0, 2)
	resp, err := env.Offchain.ListNodes(context.Background(), &node.ListNodesRequest{
		Filter: &node.ListNodesRequest_Filter{},
	})
	require.NoError(t, err)
	for _, n := range resp.Nodes {
		for _, label := range n.Labels {
			if label.Key == devenv.LabelNodeTypeKey {
				switch *label.Value {
				case devenv.LabelNodeTypeValueBootstrap:
					bootstrapNodeNames = append(bootstrapNodeNames, n.Name)
				case devenv.LabelNodeTypeValuePlugin:
					oracleNodeNames = append(oracleNodeNames, n.Name)
				default:
					t.Fatalf("unexpected n type: %s", *label.Value)
				}
			}
		}
	}

	// pick the first EVM chain selector
	chainSelector := env.AllChainSelectors()[0]

	// insert a Configurator address for the given DON
	configuratorAddr := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err = env.ExistingAddresses.Save(chainSelector, configuratorAddr, //nolint: staticcheck // I don't care that ExistingAddresses is deprecated. We will fix it later.
		deployment.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  deployment.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	oracleSpec := `name = 'don | 1'
type = 'offchainreporting2'
schemaVersion = 1
contractID = '0x4170ed0880ac9a755fd29b2688956bd959f923f4'
ocrKeyBundleID = 'cee9d802bf0e28bc74c78d7512e44b25ce6580bf5c45ed15186ae871a3437eb1'
maxTaskDuration = '1s'
contractConfigTrackerPollInterval = '1s'
relay = 'evm'
pluginType = 'llo'

[relayConfig]
chainID = '90000001'
lloConfigMode = 'bluegreen'
lloDonID = 1

[pluginConfig]
channelDefinitionsContractAddress = '0x000000000000000000000000000000000000dEaD'
channelDefinitionsContractFromBlock = 0
donID = 1
servers = {'mercury-pipeline-testnet-producer.TEST.cldev.cloud:1340' = '0000005187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c'}
`

	bootstrapSpec := `name = 'bootstrap'
type = 'bootstrap'
schemaVersion = 1
contractID = '0x4170ed0880ac9a755fd29b2688956bd959f923f4'
donID = 1
relay = 'evm'

[relayConfig]
chainID = '90000001'
`

	config := CsDistributeLLOJobSpecsConfig{
		ChainSelectorEVM: chainSelector,
		Filter: &jd.ListFilter{
			DONID:             donID,
			DONName:           donName,
			EnvLabel:          envName,
			NumOracleNodes:    2,
			NumBootstrapNodes: 1,
		},
		FromBlock:                   0,
		ConfigMode:                  "bluegreen",
		ChannelConfigStoreAddr:      common.HexToAddress("DEAD"),
		ChannelConfigStoreFromBlock: 0,
		ConfiguratorAddress:         configuratorAddr,
		Servers: map[string]string{
			"mercury-pipeline-testnet-producer.TEST.cldev.cloud:1340": "0000005187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c",
		},
		NodeNames: append(bootstrapNodeNames, oracleNodeNames...),
	}

	tests := []struct {
		name                 string
		prepConfFn           func(CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig
		wantErr              *string
		wantOracleSpec       string
		wantBootstrapSpec    string
		wantNumOracleJobs    int
		wantNumBootstrapJobs int
	}{
		{
			name:                 "success",
			wantOracleSpec:       oracleSpec,
			wantBootstrapSpec:    bootstrapSpec,
			wantNumOracleJobs:    2,
			wantNumBootstrapJobs: 1,
		},
		{
			// This test only makes sense when run after "success" because the two use the same ExternalJobID.
			name:                 "success proposing updates to existing jobs",
			wantOracleSpec:       oracleSpec,
			wantBootstrapSpec:    bootstrapSpec,
			wantNumOracleJobs:    2,
			wantNumBootstrapJobs: 1,
		},
		{
			name: "success when sending jobs to a subset of nodes",
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.NodeNames = append(bootstrapNodeNames, oracleNodeNames[:1]...) //nolint: gocritic // I want a combined list. GoCritic doesn't like it.
				c.Filter = &jd.ListFilter{
					DONID:             donID,
					DONName:           donName,
					EnvLabel:          envName,
					NumOracleNodes:    1,
					NumBootstrapNodes: 1,
				}
				return c
			},
			wantOracleSpec:       oracleSpec,
			wantBootstrapSpec:    bootstrapSpec,
			wantNumOracleJobs:    1,
			wantNumBootstrapJobs: 1,
		},
		{
			name: "success when sending jobs to the remaining nodes",
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.NodeNames = []string{oracleNodeNames[0]}
				c.Filter = &jd.ListFilter{
					DONID:             donID,
					DONName:           donName,
					EnvLabel:          envName,
					NumOracleNodes:    1,
					NumBootstrapNodes: 0,
				}
				return c
			},
			wantOracleSpec:       oracleSpec,
			wantBootstrapSpec:    "",
			wantNumOracleJobs:    1,
			wantNumBootstrapJobs: 0,
		},
		{
			name: "missing channel config store",
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.ChannelConfigStoreAddr = common.Address{}
				return c
			},
			wantErr: pointer.To("channel config store address is required"),
		},
		{
			name: "missing servers",
			prepConfFn: func(c CsDistributeLLOJobSpecsConfig) CsDistributeLLOJobSpecsConfig {
				c.Servers = nil
				return c
			},
			wantErr: pointer.To("servers map is required"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conf := config
			if tc.prepConfFn != nil {
				conf = tc.prepConfFn(conf)
			}
			_, out, err := changeset.ApplyChangesetsV2(t,
				env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(CsDistributeLLOJobSpecs{}, conf),
				},
			)

			if tc.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), *tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, tc.wantNumOracleJobs+tc.wantNumBootstrapJobs)

			// These are lines with dynamic values which we cannot compare.
			linesToStrip := []string{"externalJobID", "transmitterID", "p2pv2Bootstrappers", "ocrKeyBundleID"}
			wantBootstrapSpec := testutil.StripLineContaining(tc.wantBootstrapSpec, linesToStrip)
			wantOracleSpec := testutil.StripLineContaining(tc.wantOracleSpec, linesToStrip)

			foundBootstrapJobs := 0
			foundOracleJobs := 0
			for _, j := range out[0].Jobs {
				spec := testutil.StripLineContaining(j.Spec, linesToStrip)
				if strings.Contains(spec, "bootstrap") {
					require.Equal(t, wantBootstrapSpec, spec)
					foundBootstrapJobs++
				} else {
					require.Equal(t, wantOracleSpec, spec)
					foundOracleJobs++
				}
			}
			require.Equal(t, tc.wantNumBootstrapJobs, foundBootstrapJobs)
			require.Equal(t, tc.wantNumOracleJobs, foundOracleJobs)
		})
	}
}
