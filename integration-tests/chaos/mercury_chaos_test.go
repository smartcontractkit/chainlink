package chaos

import (
	"context"
	"fmt"
	"testing"

	"github.com/onsi/gomega"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/chaos"
)

var (
	defaultDONSettings = map[string]interface{}{
		"stateful": true,
		"capacity": "10Gi",
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "250m",
				"memory": "256Mi",
			},
			"limits": map[string]interface{}{
				"cpu":    "250m",
				"memory": "256Mi",
			},
		},
	}
	defaultMercuryDBSettings = map[string]interface{}{
		"stateful": "true",
		"capacity": "10Gi",
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2048Mi",
			},
		},
	}
)

func TestMercuryChaos(t *testing.T) {
	t.Parallel()
	l := actions.GetTestLogger(t)
	testCases := map[string]struct {
		chaosFunc  chaos.ManifestFunc
		chaosProps *chaos.Props
	}{
		//NetworkChaosFailMajorityNetwork: {
		//	chaos.NewNetworkPartition,
		//	&chaos.Props{
		//		FromLabels:  &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		ToLabels:    &map[string]*string{ChaosGroupMinority: a.Str("1")},
		//		DurationStr: "1m",
		//	},
		//},
		//NetworkChaosFailBlockchainNode: {
		//	chaos.NewNetworkPartition,
		//	&chaos.Props{
		//		FromLabels:  &map[string]*string{"app": a.Str("geth")},
		//		ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: a.Str("1")},
		//		DurationStr: "1m",
		//	},
		//},
		//
		//PodChaosFailMinorityNodes: {
		//	chaos.NewFailPods,
		//	&chaos.Props{
		//		LabelsSelector: &map[string]*string{ChaosGroupMinority: a.Str("1")},
		//		DurationStr:    "1m",
		//	},
		//},
		//PodChaosFailMajorityNodes: {
		//	chaos.NewFailPods,
		//	&chaos.Props{
		//		LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		DurationStr:    "1m",
		//	},
		//},
		//PodChaosFailMajorityDB: {
		//	chaos.NewFailPods,
		//	&chaos.Props{
		//		LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		DurationStr:    "1m",
		//		ContainerNames: &[]*string{a.Str("chainlink-db")},
		//	},
		//},
		PodChaosFailMercury: {
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{"app": a.Str("mercury-server")},
				DurationStr:    "1m",
			},
		},
	}

	for n, tst := range testCases {
		name := n
		testCase := tst
		t.Run(fmt.Sprintf("Mercury_%s", name), func(t *testing.T) {
			t.Parallel()

			feeds := []string{"feed-1"}
			env, err := mercury.SetupMercuryMultiFeedEnv(t.Name(), "chaos", feeds, defaultDONSettings, defaultMercuryDBSettings)
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
			err = env.Env.WaitHealthy()
			require.NoError(t, err)
			env.Env.Cfg.NoManifestUpdate = true
			err = env.Env.Run()
			require.NoError(t, err)
			blockAfterChaos, err := env.EvmClient.LatestBlockNumber(context.Background())
			require.NoError(t, err, "Err getting latest block number")

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				report, _, err := env.MSClient.GetReports(feeds[0], blockAfterChaos)
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Error getting report from Mercury Server")
				l.Info().Interface("Report", report).Msg("Last report received")
				g.Expect(report.ChainlinkBlob).ShouldNot(gomega.BeEmpty(), "Report response does not contain chainlinkBlob")
				err = mercuryactions.ValidateReport([]byte(report.ChainlinkBlob))
				g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Error validating mercury report")
				l.Info().Interface("Report", report).Msg("Validated report received")
			}, "3m", "1s").Should(gomega.Succeed())
		})
	}
}
