package contracts

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/rs/zerolog/log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("OCR Feed @ocr", func() {

	DescribeTable("Deploys and watches an OCR feed @ocr", func(
		suiteInit environment.K8sEnvSpecInit,
	) {
		var (
			suiteSetup     *actions.DefaultSuiteSetup
			chainlinkNodes []client.Chainlink
			adapter        environment.ExternalAdapter
			defaultWallet  client.BlockchainWallet
			ocrInstance    contracts.OffchainAggregator
			em             *client.ExplorerClient
		)

		By("Deploying the environment", func() {
			var err error
			suiteSetup, err = actions.DefaultLocalSetup(
				suiteInit,
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			em, err = environment.GetExplorerMockClient(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			chainlinkNodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			defaultWallet = suiteSetup.Wallets.Default()
			suiteSetup.Client.ParallelTransactions(true)
		})

		By("Funding nodes and deploying OCR contract", func() {
			err := actions.FundChainlinkNodes(
				chainlinkNodes,
				suiteSetup.Client,
				defaultWallet,
				big.NewFloat(0.05),
				big.NewFloat(2),
			)
			Expect(err).ShouldNot(HaveOccurred())

			// Deploy and config OCR contract
			deployer, err := contracts.NewContractDeployer(suiteSetup.Client)
			Expect(err).ShouldNot(HaveOccurred())

			ocrInstance, err = deployer.DeployOffChainAggregator(defaultWallet, contracts.DefaultOffChainAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred())
			err = ocrInstance.SetConfig(
				defaultWallet,
				chainlinkNodes,
				contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes)),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = ocrInstance.Fund(defaultWallet, nil, big.NewFloat(2))
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Sending OCR jobs to chainlink nodes", func() {
			// Initialize bootstrap node
			bootstrapNode := chainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				ContractAddress: ocrInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.CreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred())

			bta := client.BridgeTypeAttributes{
				Name: "variable",
				URL:  fmt.Sprintf("%s/variable", adapter.ClusterURL()),
			}

			// Send OCR job to other nodes
			for index := 1; index < len(chainlinkNodes); index++ {
				nodeP2PIds, err := chainlinkNodes[index].ReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := chainlinkNodes[index].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeys, err := chainlinkNodes[index].ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				err = chainlinkNodes[index].CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
				}
				_, err = chainlinkNodes[index].CreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Checking OCR rounds", func() {
			roundTimeout := time.Minute * 2
			// Set adapter answer to 5
			err := adapter.SetVariable(5)
			Expect(err).ShouldNot(HaveOccurred())
			err = ocrInstance.RequestNewRound(defaultWallet)
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for the first round
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrInstance, big.NewInt(1), roundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(ocrInstance.Address(), ocrRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// Check answer is as expected
			answer, err := ocrInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(5)), "Latest answer from OCR is not as expected")

			// Change adapter answer to 10
			err = adapter.SetVariable(10)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for the second round
			ocrRound = contracts.NewOffchainAggregatorRoundConfirmer(ocrInstance, big.NewInt(2), roundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(ocrInstance.Address(), ocrRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			// Check answer is as expected
			answer, err = ocrInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(10)), "Latest answer from OCR is not as expected")
		})

		By("Checking explorer telemetry", func() {
			mc, err := em.Count()
			log.Debug().Interface("Telemetry", mc).Msg("Explorer messages count")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mc.Errors).Should(Equal(0))
			Expect(mc.Unknown).Should(Equal(0))
			Expect(mc.Broadcast).Should(BeNumerically(">", 1))
			Expect(mc.DHTAnnounce).Should(BeNumerically(">", 1))
			Expect(mc.NewEpoch).Should(BeNumerically(">", 1))
			Expect(mc.ObserveReq).Should(BeNumerically(">", 1))
			Expect(mc.Received).Should(BeNumerically(">", 1))
			Expect(mc.ReportReq).Should(BeNumerically(">", 1))
			Expect(mc.RoundStarted).Should(BeNumerically(">", 1))
			Expect(mc.Sent).Should(BeNumerically(">", 1))
		})

		By("Tearing down the environment", suiteSetup.TearDown())
	},
		Entry("all the same version", environment.NewChainlinkCluster(5)),
		Entry("different versions", environment.NewMixedVersionChainlinkCluster(5, 2)),
	)
})
