package chaos

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/lib/utils/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

var (
	defaultOCRSettings = map[string]interface{}{
		"replicas": 6,
		"db": map[string]interface{}{
			"stateful": true,
			"capacity": "1Gi",
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
		},
	}
	chaosStartRound int64 = 1
	chaosEndRound   int64 = 4
)

func TestOCRChaos(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Chaos"}, tc.OCR)
	require.NoError(t, err, "Error getting config")

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(config.GetChainlinkImageConfig(), target)
		ctf_config.MightConfigOverridePyroscopeKey(config.GetPyroscopeConfig(), target)
	}

	tomlConfig, err := actions.BuildTOMLNodeConfigForK8s(&config, networks.MustGetSelectedNetworkConfig(config.Network)[0])
	require.NoError(t, err, "Error building TOML config")

	defaultOCRSettings["toml"] = tomlConfig

	chainlinkCfg := chainlink.NewWithOverride(0, defaultOCRSettings, config.ChainlinkImage, overrideFn)

	testCases := map[string]struct {
		networkChart environment.ConnectedChart
		clChart      environment.ConnectedChart
		chaosFunc    chaos.ManifestFunc
		chaosProps   *chaos.Props
	}{
		// network-* and pods-* are split intentionally into 2 parallel groups
		// we can't use chaos.NewNetworkPartition and chaos.NewFailPods in parallel
		// because of jsii runtime bug, see Makefile and please use those targets to run tests
		//
		// We are using two chaos experiments to simulate pods/network faults,
		// check chaos.NewFailPods method (https://chaos-mesh.org/docs/simulate-pod-chaos-on-kubernetes/)
		// and chaos.NewNetworkPartition method (https://chaos-mesh.org/docs/simulate-network-chaos-on-kubernetes/)
		// in order to regenerate Go bindings if k8s version will be updated
		// you can pull new CRD spec from your current cluster and check README here
		// https://github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/blob/master/README.md
		NetworkChaosFailMajorityNetwork: {
			ethereum.New(nil),
			chainlinkCfg,
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{ChaosGroupMajority: ptr.Ptr("1")},
				ToLabels:    &map[string]*string{ChaosGroupMinority: ptr.Ptr("1")},
				DurationStr: "1m",
			},
		},
		NetworkChaosFailBlockchainNode: {
			ethereum.New(nil),
			chainlinkCfg,
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": ptr.Ptr("geth")},
				ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: ptr.Ptr("1")},
				DurationStr: "1m",
			},
		},
		PodChaosFailMinorityNodes: {
			ethereum.New(nil),
			chainlinkCfg,
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMinority: ptr.Ptr("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityNodes: {
			ethereum.New(nil),
			chainlinkCfg,
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: ptr.Ptr("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityDB: {
			ethereum.New(nil),
			chainlinkCfg,
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: ptr.Ptr("1")},
				DurationStr:    "1m",
				ContainerNames: &[]*string{ptr.Ptr("chainlink-db")},
			},
		},
	}

	for n, tst := range testCases {
		name := n
		testCase := tst
		t.Run(fmt.Sprintf("OCR_%s", name), func(t *testing.T) {
			t.Parallel()

			nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels("data-feedsv1.0", "chaos")
			require.NoError(t, err, "Error creating required chain.link labels for namespace")

			workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels("data-feedsv1.0", "chaos")
			require.NoError(t, err, "Error creating required chain.link labels for workloads and pods")

			testEnvironment := environment.New(&environment.Config{
				NamespacePrefix: fmt.Sprintf("chaos-ocr-%s", name),
				Test:            t,
				Labels:          nsLabels,
				WorkloadLabels:  workloadPodLabels,
				PodLabels:       workloadPodLabels,
			}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(testCase.networkChart).
				AddHelm(testCase.clChart)
			err = testEnvironment.Run()
			require.NoError(t, err)
			if testEnvironment.WillUseRemoteRunner() {
				return
			}

			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, "instance=node-", 1, 2, ChaosGroupMinority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, "instance=node-", 3, 5, ChaosGroupMajority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, "instance=node-", 2, 5, ChaosGroupMajorityPlus)
			require.NoError(t, err)

			cfg := config.MustCopy().(tc.TestConfig)

			network := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())[0]
			network = seth_utils.MustReplaceSimulatedNetworkUrlWithK8(l, network, *testEnvironment)

			seth, err := seth_utils.GetChainClient(&cfg, network)
			require.NoError(t, err, "Error creating seth client")

			chainlinkNodes, err := nodeclient.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
			bootstrapNode, workerNodes := chainlinkNodes[0], chainlinkNodes[1:]
			t.Cleanup(func() {
				err := actions.TeardownRemoteSuite(t, seth, testEnvironment.Cfg.Namespace, chainlinkNodes, nil, &cfg)
				require.NoError(t, err, "Error tearing down environment")
			})

			ms := ctfClient.ConnectMockServer(testEnvironment)
			linkContract, err := actions.LinkTokenContract(l, seth, config.OCR)
			require.NoError(t, err, "Error deploying link token contract")

			err = actions.FundChainlinkNodesFromRootAddress(l, seth, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes), big.NewFloat(10))
			require.NoError(t, err)

			ocrInstances, err := actions.SetupOCRv1Contracts(l, seth, config.OCR, common.HexToAddress(linkContract.Address()), contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(workerNodes))
			require.NoError(t, err)
			err = actions.CreateOCRJobs(ocrInstances, bootstrapNode, workerNodes, 5, ms, fmt.Sprint(seth.ChainID))
			require.NoError(t, err)

			chaosApplied := false

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				for _, ocr := range ocrInstances {
					err := ocr.RequestNewRound()
					require.NoError(t, err, "Error requesting new round")
				}
				round, err := ocrInstances[0].GetLatestRound(testcontext.Get(t))
				g.Expect(err).ShouldNot(gomega.HaveOccurred())
				l.Info().Int64("RoundID", round.RoundId.Int64()).Msg("Latest OCR Round")
				if round.RoundId.Int64() == chaosStartRound && !chaosApplied {
					chaosApplied = true
					_, err = testEnvironment.Chaos.Run(testCase.chaosFunc(testEnvironment.Cfg.Namespace, testCase.chaosProps))
					require.NoError(t, err)
				}
				g.Expect(round.RoundId.Int64()).Should(gomega.BeNumerically(">=", chaosEndRound))
			}, "6m", "3s").Should(gomega.Succeed())
		})
	}
}
