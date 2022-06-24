package smoke

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	it "github.com/smartcontractkit/chainlink/integration-tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

var _ = Describe("Direct request suite @runlog", func() {
	DescribeTable("Direct request suite on different EVM networks", func(
		clientFunc func(*environment.Environment) (blockchain.EVMClient, error),
		networkChart environment.ConnectedChart,
		clChart environment.ConnectedChart,
	) {
		var (
			err            error
			c              blockchain.EVMClient
			cd             contracts.ContractDeployer
			chainlinkNodes []client.Chainlink
			oracle         contracts.Oracle
			consumer       contracts.APIConsumer
			jobUUID        uuid.UUID
			ms             *client.MockserverClient
			e              *environment.Environment
		)

		By("Deploying the environment", func() {
			e = environment.New(nil).
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, nil))
			err = e.Run()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Connecting to launched resources", func() {
			c, err = blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings)(e)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			cd, err = contracts.NewContractDeployer(c)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(e)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			ms, err = client.ConnectMockServer(e)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Funding Chainlink nodes", func() {
			ethAmount, err := c.EstimateCostForChainlinkOperations(1)
			Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
			err = actions.FundChainlinkNodes(chainlinkNodes, c, ethAmount)
			Expect(err).ShouldNot(HaveOccurred(), "Funding chainlink nodes with ETH shouldn't fail")
		})

		By("Deploying contracts", func() {
			lt, err := cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
			oracle, err = cd.DeployOracle(lt.Address())
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Oracle Contract shouldn't fail")
			consumer, err = cd.DeployAPIConsumer(lt.Address())
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Consumer Contract shouldn't fail")
			err = c.SetDefaultWallet(0)
			Expect(err).ShouldNot(HaveOccurred(), "Setting default wallet shouldn't fail")
			err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred(), "Transferring %d to consumer contract shouldn't fail", big.NewInt(2e18))
		})

		By("Creating directrequest job", func() {
			err = ms.SetValuePath("/variable", 5)
			Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")

			jobUUID = uuid.NewV4()

			bta := client.BridgeTypeAttributes{
				Name: fmt.Sprintf("five-%s", jobUUID.String()),
				URL:  fmt.Sprintf("%s/variable", ms.Config.ClusterURL),
			}
			err = chainlinkNodes[0].CreateBridge(&bta)
			Expect(err).ShouldNot(HaveOccurred(), "Creating bridge shouldn't fail")

			os := &client.DirectRequestTxPipelineSpec{
				BridgeTypeAttributes: bta,
				DataPath:             "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred(), "Building observation source spec shouldn't fail")

			_, err = chainlinkNodes[0].CreateJob(&client.DirectRequestJobSpec{
				Name:                     "direct_request",
				MinIncomingConfirmations: "1",
				ContractAddress:          oracle.Address(),
				ExternalJobID:            jobUUID.String(),
				ObservationSource:        ost,
			})
			Expect(err).ShouldNot(HaveOccurred(), "Creating direct_request job shouldn't fail")
		})

		By("Calling oracle contract", func() {
			jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
			var jobID [32]byte
			copy(jobID[:], jobUUIDReplaces)
			err = consumer.CreateRequestTo(
				oracle.Address(),
				jobID,
				big.NewInt(1e18),
				fmt.Sprintf("%s/variable", ms.Config.ClusterURL),
				"data,result",
				big.NewInt(100),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Calling oracle contract shouldn't fail")
		})
		By("receives API call data on-chain", func() {
			Eventually(func(g Gomega) {
				d, err := consumer.Data(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Getting data from consumer contract shouldn't fail")
				g.Expect(d).ShouldNot(BeNil(), "Expected the initial on chain data to be nil")
				log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
				g.Expect(d.Int64()).Should(BeNumerically("==", 5), "Expected the on-chain data to be 5, but found %d", d.Int64())
			}, "2m", "1s").Should(Succeed())
		})
		By("Tearing down the environment", func() {
			c.GasStats().PrintStats()
			err = actions.TeardownSuite(e, utils.ProjectRoot, chainlinkNodes, nil, c)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	},
		Entry("VRF on Geth @geth",
			blockchain.NewEthereumMultiNodeClientSetup(it.DefaultGethSettings),
			ethereum.New(nil),
			chainlink.New(0, nil),
		),
	)
})
