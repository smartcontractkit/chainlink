package testcommon

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

type RunlogSetupInputs struct {
	S             *actions.DefaultSuiteSetup
	Adapter       environment.ExternalAdapter
	Nodes         []client.Chainlink
	NodeAddresses []common.Address
	Oracle        contracts.Oracle
	Consumer      contracts.APIConsumer
	JobUUID       uuid.UUID
	Err           error
}

func SetupRunlogTest(i *RunlogSetupInputs) {
	ginkgo.By("Deploying the environment", func() {
		i.S, i.Err = actions.DefaultLocalSetup(
			environment.NewChainlinkCluster(1),
			client.NewNetworkFromConfig,
			ConfigLocation(),
		)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
		i.Adapter, i.Err = environment.GetExternalAdapter(i.S.Env)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
	})
	ginkgo.By("Funding Chainlink nodes", func() {
		i.Nodes, i.Err = environment.GetChainlinkClients(i.S.Env)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
		i.NodeAddresses, i.Err = actions.ChainlinkNodeAddresses(i.Nodes)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
		i.Err = actions.FundChainlinkNodes(i.Nodes, i.S.Client, i.S.Wallets.Default(), big.NewFloat(2), nil)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
	})
	ginkgo.By("Deploying and funding the contracts", func() {
		i.Oracle, i.Err = i.S.Deployer.DeployOracle(i.S.Wallets.Default(), i.S.Link.Address())
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
		i.Consumer, i.Err = i.S.Deployer.DeployAPIConsumer(i.S.Wallets.Default(), i.S.Link.Address())
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
		i.Err = i.Consumer.Fund(i.S.Wallets.Default(), nil, big.NewFloat(2))
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
	})
	ginkgo.By("Permitting node to fulfill request", func() {
		i.Err = i.Oracle.SetFulfillmentPermission(i.S.Wallets.Default(), i.NodeAddresses[0].Hex(), true)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
	})
	ginkgo.By("Creating directrequest job", func() {
		i.JobUUID = uuid.NewV4()

		bta := client.BridgeTypeAttributes{
			Name: "five",
			URL:  fmt.Sprintf("%s/five", i.Adapter.ClusterURL()),
		}
		i.Err = i.Nodes[0].CreateBridge(&bta)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())

		os := &client.DirectRequestTxPipelineSpec{
			BridgeTypeAttributes: bta,
			DataPath:             "data,result",
		}
		ost, err := os.String()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		_, err = i.Nodes[0].CreateJob(&client.DirectRequestJobSpec{
			Name:              "direct_request",
			ContractAddress:   i.Oracle.Address(),
			ExternalJobID:     i.JobUUID.String(),
			ObservationSource: ost,
		})
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
	})
	ginkgo.By("Calling oracle contract", func() {
		jobUUIDReplaces := strings.Replace(i.JobUUID.String(), "-", "", 4)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
		var jobID [32]byte
		copy(jobID[:], jobUUIDReplaces)
		i.Err = i.Consumer.CreateRequestTo(
			i.S.Wallets.Default(),
			i.Oracle.Address(),
			jobID,
			big.NewInt(1e18),
			i.Adapter.ClusterURL()+"/five",
			"data,result",
			big.NewInt(100),
		)
		gomega.Expect(i.Err).ShouldNot(gomega.HaveOccurred())
	})
}
