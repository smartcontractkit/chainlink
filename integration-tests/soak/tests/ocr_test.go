package soak

//revive:disable:dot-imports
import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

var _ = Describe("OCR Soak Test @soak-ocr", func() {
	var (
		err             error
		testEnvironment *environment.Environment
		ocrSoakTest     *testsetups.OCRSoakTest
		soakNetwork     *blockchain.EVMNetwork
	)

	BeforeEach(func() {
		By("Connecting to the soak environment", func() {
			soakNetwork = blockchain.LoadNetworkFromEnvironment()
			testEnvironment = environment.New(&environment.Config{InsideK8s: true})
			err = testEnvironment.
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(&ethereum.Props{
					NetworkName: soakNetwork.Name,
					Simulated:   soakNetwork.Simulated,
				})).
				AddHelm(chainlink.New(0, nil)).
				AddHelm(chainlink.New(1, nil)).
				AddHelm(chainlink.New(2, nil)).
				AddHelm(chainlink.New(3, nil)).
				AddHelm(chainlink.New(4, nil)).
				AddHelm(chainlink.New(5, nil)).
				Run()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setting up Soak Test", func() {
			chainClient, err := blockchain.NewMetisMultiNodeClientSetup(soakNetwork)(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			ocrSoakTest = testsetups.NewOCRSoakTest(&testsetups.OCRSoakTestInputs{
				BlockchainClient:     chainClient,
				TestDuration:         time.Minute * 5,
				NumberOfContracts:    4,
				ChainlinkNodeFunding: big.NewFloat(.1),
				ExpectedRoundTime:    time.Minute,
				RoundTimeout:         time.Minute * 10,
				TimeBetweenRounds:    time.Minute,
				StartingAdapterValue: 5,
			})
			ocrSoakTest.Setup(testEnvironment)
		})
	})

	Describe("With soak test contracts deployed", func() {
		It("runs the soak test until error or timeout", func() {
			ocrSoakTest.Run()
		})
	})

	AfterEach(func() {
		if err = actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals()); err != nil {
			log.Error().Err(err).Msg("Error when tearing down remote suite")
		}
		log.Info().Msg("Soak Test Concluded")
	})
})
