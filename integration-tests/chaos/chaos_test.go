package chaos_test

//revive:disable:dot-imports
import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	defaultOCRSettings = map[string]interface{}{
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
)

const (
	// ChaosGroupMinority a group of faulty nodes, even if they fail OCR must work
	ChaosGroupMinority = "chaosGroupMinority"
	// ChaosGroupMajority a group of nodes that are working even if minority fails
	ChaosGroupMajority = "chaosGroupMajority"
	// ChaosGroupMajorityPlus a group of nodes that are majority + 1
	ChaosGroupMajorityPlus = "chaosGroupMajority"
)

var _ = Describe("OCR chaos test @chaos-ocr", func() {
	var (
		testScenarios = []TableEntry{
			Entry("Must survive minority removal for 1m @chaos-ocr-fail-minority",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, defaultOCRSettings),
				chaos.NewFailPods,
				&chaos.Props{
					LabelsSelector: &map[string]*string{ChaosGroupMinority: a.Str("1")},
					DurationStr:    "1m",
				},
			),
			Entry("Must recover from majority removal @chaos-ocr-fail-majority",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, defaultOCRSettings),
				chaos.NewFailPods,
				&chaos.Props{
					LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
					DurationStr:    "1m",
				},
			),
			Entry("Must recover from majority DB failure @chaos-ocr-fail-majority-db",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, defaultOCRSettings),
				chaos.NewFailPods,
				&chaos.Props{
					LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
					DurationStr:    "1m",
					ContainerNames: &[]*string{a.Str("chainlink-db")},
				},
			),
			Entry("Must recover from majority network failure @chaos-ocr-fail-majority-network",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, defaultOCRSettings),
				chaos.NewNetworkPartition,
				&chaos.Props{
					FromLabels:  &map[string]*string{ChaosGroupMajority: a.Str("1")},
					ToLabels:    &map[string]*string{ChaosGroupMinority: a.Str("1")},
					DurationStr: "1m",
				},
			),
			Entry("Must recover from blockchain node network failure @chaos-ocr-fail-blockchain-node",
				blockchain.NewEthereumMultiNodeClientSetup(networks.SimulatedEVM),
				ethereum.New(nil),
				chainlink.New(0, defaultOCRSettings),
				chaos.NewNetworkPartition,
				&chaos.Props{
					FromLabels:  &map[string]*string{"app": a.Str("geth")},
					ToLabels:    &map[string]*string{ChaosGroupMajorityPlus: a.Str("1")},
					DurationStr: "1m",
				},
			),
		}

		testEnvironment *environment.Environment
		chainlinkNodes  []client.Chainlink
		chainClient     blockchain.EVMClient

		chaosStartRound int64 = 1
		chaosEndRound   int64 = 4
		chaosApplied          = false
	)

	AfterEach(func() {
		err := actions.TeardownSuite(testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	DescribeTable("OCR chaos on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		networkChart environment.ConnectedChart,
		clChart environment.ConnectedChart,
		chaosFunc chaos.ManifestFunc,
		chaosProps *chaos.Props,
	) {
		By("Deploying the environment")
		testEnvironment = environment.New(&environment.Config{NamespacePrefix: "chaos-core"}).
			AddHelm(mockservercfg.New(nil)).
			AddHelm(mockserver.New(nil)).
			AddHelm(networkChart).
			AddHelm(clChart)
		err := testEnvironment.Run()
		Expect(err).ShouldNot(HaveOccurred())

		err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 2, ChaosGroupMinority)
		Expect(err).ShouldNot(HaveOccurred())
		err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 3, 5, ChaosGroupMajority)
		Expect(err).ShouldNot(HaveOccurred())
		err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 2, 5, ChaosGroupMajorityPlus)
		Expect(err).ShouldNot(HaveOccurred())

		By("Connecting to launched resources")
		chainClient, err = clientFunc(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
		cd, err := contracts.NewContractDeployer(chainClient)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

		chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		ms, err := ctfClient.ConnectMockServer(testEnvironment)
		Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")

		chainClient.ParallelTransactions(true)
		Expect(err).ShouldNot(HaveOccurred())

		lt, err := cd.DeployLinkTokenContract()
		Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

		By("Funding Chainlink nodes")
		err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(10))
		Expect(err).ShouldNot(HaveOccurred())

		By("Deploying OCR contracts")
		ocrInstances := actions.DeployOCRContracts(1, lt, cd, chainlinkNodes, chainClient)
		err = chainClient.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
		By("Setting adapter responses", actions.SetAllAdapterResponsesToTheSameValue(5, ocrInstances, chainlinkNodes, ms))
		By("Creating OCR jobs", actions.CreateOCRJobs(ocrInstances, chainlinkNodes, ms))

		Eventually(func(g Gomega) {
			for _, ocr := range ocrInstances {
				err := ocr.RequestNewRound()
				Expect(err).ShouldNot(HaveOccurred())
			}
			round, err := ocrInstances[0].GetLatestRound(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred())
			log.Info().Int64("RoundID", round.RoundId.Int64()).Send()
			if round.RoundId.Int64() == chaosStartRound && !chaosApplied {
				chaosApplied = true
				_, err = testEnvironment.Chaos.Run(chaosFunc(testEnvironment.Cfg.Namespace, chaosProps))
				Expect(err).ShouldNot(HaveOccurred())
			}
			g.Expect(round.RoundId.Int64()).Should(BeNumerically(">=", chaosEndRound))
		}, "6m", "3s").Should(Succeed())
	},
		testScenarios,
	)
})
