package testcommon

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

func ConfigLocation() string {
	relPath := "../config.yml"
	path, err := filepath.Abs(relPath)
	if err != nil {
		log.Err(err).Str("Relative Path", relPath).Msg("Error finding the local config")
		return ""
	}
	return path
}

type OCRSetupInputs struct {
	SuiteSetup     *actions.DefaultSuiteSetup
	ChainlinkNodes []client.Chainlink
	Adapter        environment.ExternalAdapter
	DefaultWallet  client.BlockchainWallet
	OCRInstance    contracts.OffchainAggregator
	Em             *client.ExplorerClient
}

func DeployOCRForEnv(i *OCRSetupInputs, envInit environment.K8sEnvSpecInit) {
	ginkgo.By("Deploying the environment", func() {
		var err error
		i.SuiteSetup, err = actions.DefaultLocalSetup(
			envInit,
			client.NewNetworkFromConfig,
			ConfigLocation(),
		)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		i.Em, err = environment.GetExplorerMockClient(i.SuiteSetup.Env)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		i.Adapter, err = environment.GetExternalAdapter(i.SuiteSetup.Env)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		i.ChainlinkNodes, err = environment.GetChainlinkClients(i.SuiteSetup.Env)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		i.DefaultWallet = i.SuiteSetup.Wallets.Default()
		i.SuiteSetup.Client.ParallelTransactions(true)
	})
}

func SetupOCRTest(i *OCRSetupInputs) {
	ginkgo.By("Funding nodes and deploying OCR contract", func() {
		err := actions.FundChainlinkNodes(
			i.ChainlinkNodes,
			i.SuiteSetup.Client,
			i.DefaultWallet,
			big.NewFloat(0.05),
			big.NewFloat(2),
		)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// Deploy and config OCR contract
		deployer, err := contracts.NewContractDeployer(i.SuiteSetup.Client)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		i.OCRInstance, err = deployer.DeployOffChainAggregator(i.DefaultWallet, contracts.DefaultOffChainAggregatorOptions())
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		err = i.OCRInstance.SetConfig(
			i.DefaultWallet,
			i.ChainlinkNodes,
			contracts.DefaultOffChainAggregatorConfig(len(i.ChainlinkNodes)),
		)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		err = i.OCRInstance.Fund(i.DefaultWallet, nil, big.NewFloat(2))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		err = i.SuiteSetup.Client.WaitForEvents()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})

	ginkgo.By("Sending OCR jobs to chainlink nodes", func() {
		// Initialize bootstrap node
		bootstrapNode := i.ChainlinkNodes[0]
		bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
			ContractAddress: i.OCRInstance.Address(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.CreateJob(bootstrapSpec)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		bta := client.BridgeTypeAttributes{
			Name: "variable",
			URL:  fmt.Sprintf("%s/variable", i.Adapter.ClusterURL()),
		}

		// Send OCR job to other nodes
		for index := 1; index < len(i.ChainlinkNodes); index++ {
			nodeP2PIds, err := i.ChainlinkNodes[index].ReadP2PKeys()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := i.ChainlinkNodes[index].PrimaryEthAddress()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			nodeOCRKeys, err := i.ChainlinkNodes[index].ReadOCRKeys()
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			err = i.ChainlinkNodes[index].CreateBridge(&bta)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    i.OCRInstance.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
			}
			_, err = i.ChainlinkNodes[index].CreateJob(ocrSpec)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		}
	})
}

func CheckRound(i *OCRSetupInputs) {
	ginkgo.By("Checking OCR rounds", func() {
		roundTimeout := time.Minute * 2
		// Set adapter answer to 5
		err := i.Adapter.SetVariable(5)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		err = i.OCRInstance.RequestNewRound(i.DefaultWallet)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		err = i.SuiteSetup.Client.WaitForEvents()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// Wait for the first round
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(i.OCRInstance, big.NewInt(1), roundTimeout)
		i.SuiteSetup.Client.AddHeaderEventSubscription(i.OCRInstance.Address(), ocrRound)
		err = i.SuiteSetup.Client.WaitForEvents()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// Check answer is as expected
		answer, err := i.OCRInstance.GetLatestAnswer(context.Background())
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(answer.Int64()).Should(gomega.Equal(int64(5)), "Latest answer from OCR is not as expected")

		// Change adapter answer to 10
		err = i.Adapter.SetVariable(10)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// Wait for the second round
		ocrRound = contracts.NewOffchainAggregatorRoundConfirmer(i.OCRInstance, big.NewInt(2), roundTimeout)
		i.SuiteSetup.Client.AddHeaderEventSubscription(i.OCRInstance.Address(), ocrRound)
		err = i.SuiteSetup.Client.WaitForEvents()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// Check answer is as expected
		answer, err = i.OCRInstance.GetLatestAnswer(context.Background())
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(answer.Int64()).Should(gomega.Equal(int64(10)), "Latest answer from OCR is not as expected")
	})
}

func CheckTelemetry(i *OCRSetupInputs) {
	ginkgo.By("Checking explorer telemetry", func() {
		mc, err := i.Em.Count()
		log.Debug().Interface("Telemetry", mc).Msg("Explorer messages count")
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(mc.Errors).Should(gomega.Equal(0))
		gomega.Expect(mc.Unknown).Should(gomega.Equal(0))
		gomega.Expect(mc.Broadcast).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.DHTAnnounce).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.NewEpoch).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.ObserveReq).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.Received).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.ReportReq).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.RoundStarted).Should(gomega.BeNumerically(">", 1))
		gomega.Expect(mc.Sent).Should(gomega.BeNumerically(">", 1))
	})
}
