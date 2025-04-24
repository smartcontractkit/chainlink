/*
This file contains test helpers for the changeset package.
The filename has a suffix of "_test.go" in order to not be included in the production build.
*/

package changeset

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonChangesets "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
)

// sendTestLLOJobs sends some test LLO jobs, which we can then revoke, retrieve, delete, etc.
func sendTestLLOJobs(t *testing.T, e deployment.Environment, numOracles, numBootstraps int) []deployment.ChangesetOutput {
	chainSel := e.AllChainSelectors()[0]
	configurator := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err := e.ExistingAddresses.Save(chainSel, configurator, //nolint: staticcheck // I don't care that ExistingAddresses is deprecated. We will fix it later.
		deployment.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  deployment.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	bootstrapNodeNames, oracleNodeNames := collectNodeNames(t, e, numOracles, numBootstraps)

	config := CsDistributeLLOJobSpecsConfig{
		ChainSelectorEVM: chainSel,
		Filter: &jd.ListFilter{
			DONID:             1,
			DONName:           "don",
			EnvLabel:          e.Name,
			NumOracleNodes:    numBootstraps,
			NumBootstrapNodes: 1,
		},
		FromBlock:                   0,
		ConfigMode:                  "bluegreen",
		ChannelConfigStoreAddr:      common.HexToAddress("DEAD"),
		ChannelConfigStoreFromBlock: 0,
		ConfiguratorAddress:         configurator,
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

func collectNodeNames(t *testing.T, e deployment.Environment, numOracles, numBootstraps int) (btNames, oracleNames []string) {
	// Collect the names of the nodes.
	bootstrapNodeNames := make([]string, 0, numBootstraps)
	oracleNodeNames := make([]string, 0, numOracles)
	resp, err := e.Offchain.ListNodes(context.Background(), &node.ListNodesRequest{
		Filter: &node.ListNodesRequest_Filter{},
	})
	require.NoError(t, err)
	for _, n := range resp.Nodes {
		for _, label := range n.Labels {
			if label.Key == utils.LabelNodeType {
				switch *label.Value {
				case jd.NodeTypeBootstrap.String():
					bootstrapNodeNames = append(bootstrapNodeNames, n.Name)
				case jd.NodeTypeOracle.String():
					oracleNodeNames = append(oracleNodeNames, n.Name)
				default:
					t.Fatalf("unexpected n type: %s", *label.Value)
				}
			}
		}
	}

	return bootstrapNodeNames, oracleNodeNames
}
