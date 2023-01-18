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
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/stretchr/testify/require"

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

const (
	// ChaosGroupMinorityOCR a group of faulty nodes, even if they fail OCR must work
	ChaosGroupMinorityOCR = "chaosGroupMinority"
	// ChaosGroupMajorityOCR a group of nodes that are working even if minority fails
	ChaosGroupMajorityOCR = "chaosGroupMajority"
	// ChaosGroupMajorityOCRPlus a group of nodes that are majority + 1
	ChaosGroupMajorityOCRPlus = "chaosGroupMajority"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestOCRChaos(t *testing.T) {
	testCases := map[string]struct {
		networkChart environment.ConnectedChart
		clChart      environment.ConnectedChart
		chaosFunc    chaos.ManifestFunc
		chaosProps   *chaos.Props
	}{
		// network-* and pods-* are split intentionally into 2 parallel groups
		// we can't use chaos.NewNetworkPartition and chaos.NewFailPods in parallel
		// because of jsii runtime bug, see Makefile
		"network-chaos-fail-majority-network": {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{ChaosGroupMajorityOCR: a.Str("1")},
				ToLabels:    &map[string]*string{ChaosGroupMinorityOCR: a.Str("1")},
				DurationStr: "1m",
			},
		},
		"network-chaos-fail-blockchain-node": {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewNetworkPartition,
			&chaos.Props{
				FromLabels:  &map[string]*string{"app": a.Str("geth")},
				ToLabels:    &map[string]*string{ChaosGroupMajorityOCRPlus: a.Str("1")},
				DurationStr: "1m",
			},
		},
		"pod-chaos-fail-minority": {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMinorityOCR: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		"pod-chaos-fail-majority": {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajorityOCR: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		"pod-chaos-fail-majority-db": {
			ethereum.New(nil),
			chainlink.New(0, defaultOCRSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMajorityOCR: a.Str("1")},
				DurationStr:    "1m",
				ContainerNames: &[]*string{a.Str("chainlink-db")},
			},
		},
	}

	for n, tst := range testCases {
		name := n
		testCase := tst
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testEnvironment := environment.New(&environment.Config{NamespacePrefix: fmt.Sprintf("chaos-ocr-%s", name)}).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(testCase.networkChart).
				AddHelm(testCase.clChart)
			err := testEnvironment.Run()
			require.NoError(t, err)

			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 2, ChaosGroupMinorityOCR)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 3, 5, ChaosGroupMajorityOCR)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 2, 5, ChaosGroupMajorityOCRPlus)
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
				err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
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

			ocrInstances := actions.DeployOCRContracts(t, 1, lt, cd, chainlinkNodes, chainClient)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)
			actions.SetAllAdapterResponsesToTheSameValue(t, 5, ocrInstances, chainlinkNodes, ms)
			actions.CreateOCRJobs(t, ocrInstances, chainlinkNodes, ms)

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
