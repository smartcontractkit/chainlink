/*
This file contains test helpers for the changeset package.
The filename has a suffix of "_test.go" in order to not be included in the production build.
*/

package changeset

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/test"
)

// sendTestLLOJobs sends some test LLO jobs, which we can then revoke, retrieve, delete, etc.
func sendTestLLOJobs(t *testing.T, e cldf.Environment, numOracles, numBootstraps int, autoApproveJobs bool) []cldf.ChangesetOutput {
	chainSel := e.AllChainSelectors()[0]
	configurator := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err := e.ExistingAddresses.Save(chainSel, configurator,
		cldf.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  cldf.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	bootstrapNodeNames, oracleNodeNames := collectNodeNames(t, e, numOracles, numBootstraps)

	var labels []*ptypes.Label
	if !autoApproveJobs {
		labels = append(labels, &ptypes.Label{
			Key: test.LabelDoNotAutoApprove,
		})
	}

	config := CsDistributeLLOJobSpecsConfig{
		ChainSelectorEVM: chainSel,
		Filter: &jd.ListFilter{
			DONID:             testutil.TestDON.ID,
			DONName:           testutil.TestDON.Name,
			EnvLabel:          testutil.TestDON.Env,
			NumOracleNodes:    numOracles,
			NumBootstrapNodes: numBootstraps,
		},
		FromBlock:                   0,
		ConfigMode:                  "bluegreen",
		ChannelConfigStoreAddr:      common.HexToAddress("DEAD"),
		ChannelConfigStoreFromBlock: 0,
		ConfiguratorAddress:         configurator,
		Labels:                      labels,
		Servers: map[string]string{
			"mercury-pipeline-testnet-producer.TEST.cldev.cloud:1340": "0000005187b1498c0ccb2e56d5ee8040a03a4955822ed208749b474058fc3f9c",
		},
		NodeNames: append(bootstrapNodeNames, oracleNodeNames...),
	}

	_, out, err := commonChangesets.ApplyChangesetsV2(t,
		e,
		[]commonChangesets.ConfiguredChangeSet{
			commonChangesets.Configure(CsDistributeLLOJobSpecs{}, config),
		},
	)
	require.NoError(t, err)
	return out
}

// sendTestStreamJobs sends some test stream jobs, which we can then revoke, retrieve, delete, etc.
func sendTestStreamJobs(t *testing.T, e cldf.Environment, numOracles int, autoApproveJobs bool) []cldf.ChangesetOutput {
	_, oracleNodeNames := collectNodeNames(t, e, numOracles, 0)

	var labels []*ptypes.Label
	if !autoApproveJobs {
		labels = append(labels, &ptypes.Label{
			Key: test.LabelDoNotAutoApprove,
		})
	}

	config := CsDistributeStreamJobSpecsConfig{
		Filter: &jd.ListFilter{
			DONID:             testutil.TestDON.ID,
			DONName:           testutil.TestDON.Name,
			EnvLabel:          testutil.TestDON.Env,
			NumOracleNodes:    numOracles,
			NumBootstrapNodes: 0,
		},
		Streams: []StreamSpecConfig{
			{
				StreamID:   1000001038,
				Name:       "ICP/USD-RefPrice",
				StreamType: jobs.StreamTypeQuote,
				ReportFields: jobs.QuoteReportFields{
					Bid: jobs.ReportFieldLLO{
						ResultPath: "data,bid",
					},
					Benchmark: jobs.ReportFieldLLO{
						ResultPath: "data,mid",
					},
					Ask: jobs.ReportFieldLLO{
						ResultPath: "data,ask",
					},
				},
				EARequestParams: EARequestParams{
					Endpoint: "cryptolwba",
					From:     "ICP",
					To:       "USD",
				},
				APIs: []string{"api1", "api2", "api3", "api4"},
			},
		},
		Labels:    labels,
		NodeNames: oracleNodeNames,
	}

	_, out, err := commonChangesets.ApplyChangesetsV2(t,
		e,
		[]commonChangesets.ConfiguredChangeSet{
			commonChangesets.Configure(CsDistributeStreamJobSpecs{}, config),
		},
	)
	require.NoError(t, err)
	return out
}

func collectNodeNames(t *testing.T, e cldf.Environment, numOracles, numBootstraps int) (btNames, oracleNames []string) {
	bootstrapNodeNames := make([]string, 0, numBootstraps)
	oracleNodeNames := make([]string, 0, numOracles)
	resp, err := e.Offchain.ListNodes(context.Background(), &node.ListNodesRequest{
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

	return bootstrapNodeNames, oracleNodeNames
}
