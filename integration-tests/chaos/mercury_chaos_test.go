package chaos

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/onsi/gomega"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"

	"github.com/smartcontractkit/chainlink-env/chaos"
)

var (
	resources = &mercury.ResourcesConfig{
		DONResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
		},
		DONDBResources: map[string]interface{}{
			"stateful": "true",
			"capacity": "1Gi",
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"cpu":    "500m",
					"memory": "1024Mi",
				},
				"limits": map[string]interface{}{
					"cpu":    "500m",
					"memory": "1024Mi",
				},
			},
		},
		MercuryResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
		},
		MercuryDBResources: map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "500m",
				"memory": "1024Mi",
			},
		},
	}
)

func TestMercuryChaos(t *testing.T) {
	t.Parallel()
	l := utils.GetTestLogger(t)
	testCases := map[string]struct {
		chaosFunc  chaos.ManifestFunc
		chaosProps *chaos.Props
	}{
		NetworkChaosFailMajorityNetwork: {
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{ChaosGroupMajority: a.Str("1")},
				ToLabels:    &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr: "1m",
			},
		},
		NetworkChaosFailBlockchainNode: {
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": a.Str("geth")},
				ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: a.Str("1")},
				DurationStr: "1m",
			},
		},
		PodChaosFailMinorityNodes: {
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityNodes: {
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityDB: {
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
				DurationStr:    "1m",
				ContainerNames: &[]*string{a.Str("chainlink-db")},
			},
		},
		PodChaosFailMercury: {
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{"app": a.Str("mercury-server")},
				DurationStr:    "1m",
			},
		},
		NetworkChaosDisruptNetworkDONMercury: {
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": a.Str("mercury-server")},
				ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: a.Str("1")},
				DurationStr: "1m",
			},
		},
	}

	for n, tst := range testCases {
		name := n
		testCase := tst
		t.Run(fmt.Sprintf("Mercury_%s", name), func(t *testing.T) {
			t.Parallel()

			feeds := mercuryactions.GenFeedIds(1)

			env, _, err := mercury.SetupMultiFeedSingleVerifierEnv(t.Name(), "chaos", feeds, resources)
			require.NoError(t, err)
			t.Cleanup(func() {
				env.Cleanup(t)
			})

			err = env.Env.Client.LabelChaosGroup(env.Env.Cfg.Namespace, 1, 2, ChaosGroupMinority)
			require.NoError(t, err)
			err = env.Env.Client.LabelChaosGroup(env.Env.Cfg.Namespace, 3, 5, ChaosGroupMajority)
			require.NoError(t, err)
			err = env.Env.Client.LabelChaosGroup(env.Env.Cfg.Namespace, 2, 5, ChaosGroupMajorityPlus)
			require.NoError(t, err)

			chaosID, err := env.Env.Chaos.Run(testCase.chaosFunc(env.Env.Cfg.Namespace, testCase.chaosProps))
			require.NoError(t, err)
			err = env.Env.Chaos.WaitForAllRecovered(chaosID)
			require.NoError(t, err)
			env.Env.Cfg.NoManifestUpdate = true
			err = env.Env.Run()
			require.NoError(t, err)
			blockAfterChaos, err := env.EvmClient.LatestBlockNumber(context.Background())
			require.NoError(t, err, "Err getting latest block number")

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				report, _, err := env.MSClient.GetReportsByFeedIdStr(mercury.Byte32ToString(feeds[0]), blockAfterChaos)
				l.Info().Interface("Report", report).Msg("Last report received")
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Error getting report from Mercury Server")
				g.Expect(report.ChainlinkBlob).ShouldNot(gomega.BeEmpty(), "Report response does not contain chainlinkBlob")

				reportBytes, err := hex.DecodeString(report.ChainlinkBlob[2:])
				require.NoError(t, err)
				r, err := mercuryactions.DecodeReport(reportBytes)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Error decoding mercury report")
				l.Info().Interface("Report", r).Msg("Validated report received")
			}, "3m", "1s").Should(gomega.Succeed())
		})
	}
}
