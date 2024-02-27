package chaos

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
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

func getDefaultOcrSettings(config *tc.TestConfig) map[string]interface{} {
	defaultOCRSettings["toml"] = networks.AddNetworksConfig(baseTOML, config.Pyroscope, networks.MustGetSelectedNetworkConfig(config.Network)[0])
	return defaultAutomationSettings
}

func TestOCRChaos(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Chaos", tc.OCR)
	require.NoError(t, err, "Error getting config")

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(config.ChainlinkImage, target)
		ctf_config.MightConfigOverridePyroscopeKey(config.Pyroscope, target)
	}

	chainlinkCfg := chainlink.NewWithOverride(0, getDefaultOcrSettings(&config), config.ChainlinkImage, overrideFn)

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
		// https://github.com/smartcontractkit/chainlink-testing-framework/k8s/blob/master/README.md
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

			testEnvironment := environment.New(&environment.Config{
				NamespacePrefix: fmt.Sprintf("chaos-ocr-%s", name),
				Test:            t,
			}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(testCase.networkChart).
				AddHelm(testCase.clChart)
			err := testEnvironment.Run()
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
			readSethCfg := cfg.GetSethConfig()
			require.NotNil(t, readSethCfg, "Seth config shouldn't be nil")

			network := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())[0]
			network = utils.MustReplaceSimulatedNetworkUrlWithK8(l, network, *testEnvironment)

			sethCfg := utils.MergeSethAndEvmNetworkConfigs(l, network, *readSethCfg)
			seth, err := seth.NewClientWithConfig(&sethCfg)
			require.NoError(t, err, "Error creating seth client")

			chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
			bootstrapNode, workerNodes := chainlinkNodes[0], chainlinkNodes[1:]
			t.Cleanup(func() {
				err := actions_seth.TeardownRemoteSuite(t, seth, testEnvironment.Cfg.Namespace, chainlinkNodes, nil, &cfg)
				require.NoError(t, err, "Error tearing down environment")
			})

			ms, err := ctfClient.ConnectMockServer(testEnvironment)
			require.NoError(t, err, "Creating mockserver clients shouldn't fail")

			linkDeploymentData, err := contracts.DeployLinkTokenContract(seth)
			require.NoError(t, err, "Error deploying link token contract")

			err = actions_seth.FundChainlinkNodesFromRootAddress(l, seth, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes), big.NewFloat(10))
			require.NoError(t, err)

			ocrInstances, err := actions_seth.DeployOCRv1Contracts(l, seth, 1, linkDeploymentData.Address, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(workerNodes))
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
