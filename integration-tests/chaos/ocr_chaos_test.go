package chaos

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/config"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/logging"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var (
	defaultOCRSettings = map[string]interface{}{
		"toml":     client.AddNetworksConfig(config.BaseOCRP2PV1Config, networks.SelectedNetwork),
		"replicas": "6",
		"db": map[string]interface{}{
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
		},
	}
	chaosStartRound int64 = 1
	chaosEndRound   int64 = 4
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestOCRChaos(t *testing.T) {
	t.Parallel()
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
		// https://github.com/smartcontractkit/chainlink-env/blob/master/README.md
		NetworkChaosFailMajorityNetwork: {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{ChaosGroupMajority: a.Str("1")},
				ToLabels:    &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr: "1m",
			},
		},
		NetworkChaosFailBlockchainNode: {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": a.Str("geth")},
				ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: a.Str("1")},
				DurationStr: "1m",
			},
		},
		PodChaosFailMinorityNodes: {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityNodes: {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		PodChaosFailMajorityDB: {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
				DurationStr:    "1m",
				ContainerNames: &[]*string{a.Str("chainlink-db")},
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

			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 2, ChaosGroupMinority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 3, 5, ChaosGroupMajority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 2, 5, ChaosGroupMajorityPlus)
			require.NoError(t, err)

			chainClient, err := blockchain.NewEVMClient(blockchain.SimulatedEVMNetwork, testEnvironment)
			require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
			cd, err := contracts.NewContractDeployer(chainClient)
			require.NoError(t, err, "Deploying contracts shouldn't fail")

			chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
			t.Cleanup(func() {
				if chainClient != nil {
					chainClient.GasStats().PrintStats()
				}
				err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.PanicLevel, chainClient)
				require.NoError(t, err, "Error tearing down environment")
			})

			ms, err := ctfClient.ConnectMockServer(testEnvironment)
			require.NoError(t, err, "Creating mockserver clients shouldn't fail")

			chainClient.ParallelTransactions(true)
			require.NoError(t, err)

			lt, err := cd.DeployLinkTokenContract()
			require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(10))
			require.NoError(t, err)

			ocrInstances, err := actions.DeployOCRContracts(1, lt, cd, chainlinkNodes, chainClient)
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)
			err = actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, ms)
			require.NoError(t, err)
			err = actions.CreateOCRJobs(ocrInstances, chainlinkNodes, ms)
			require.NoError(t, err)

			chaosApplied := false

			gom := gomega.NewGomegaWithT(t)
			gom.Eventually(func(g gomega.Gomega) {
				for _, ocr := range ocrInstances {
					err := ocr.RequestNewRound()
					require.NoError(t, err, "Error requesting new round")
				}
				round, err := ocrInstances[0].GetLatestRound(context.Background())
				g.Expect(err).ShouldNot(gomega.HaveOccurred())
				log.Info().Int64("RoundID", round.RoundId.Int64()).Msg("Latest OCR Round")
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
