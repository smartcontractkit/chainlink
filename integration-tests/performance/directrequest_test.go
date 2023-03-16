package performance

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
	"github.com/stretchr/testify/require"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func TestDirectRequestPerformance(t *testing.T) {
	testEnvironment := setupDirectRequestTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(blockchain.SimulatedEVMNetwork, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	mockServerClient, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Error connecting to mock server")

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.01))
	require.NoError(t, err, "Funding chainlink nodes with ETH shouldn't fail")

	lt, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	oracle, err := contractDeployer.DeployOracle(lt.Address())
	require.NoError(t, err, "Deploying Oracle Contract shouldn't fail")
	consumer, err := contractDeployer.DeployAPIConsumer(lt.Address())
	require.NoError(t, err, "Deploying Consumer Contract shouldn't fail")
	err = chainClient.SetDefaultWallet(0)
	require.NoError(t, err, "Setting default wallet shouldn't fail")
	err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
	require.NoError(t, err, "Transferring %d to consumer contract shouldn't fail", big.NewInt(2e18))

	err = mockServerClient.SetValuePath("/variable", 5)
	require.NoError(t, err, "Setting mockserver value path shouldn't fail")

	jobUUID := uuid.NewV4()

	bta := client.BridgeTypeAttributes{
		Name: fmt.Sprintf("five-%s", jobUUID.String()),
		URL:  fmt.Sprintf("%s/variable", mockServerClient.Config.ClusterURL),
	}
	err = chainlinkNodes[0].MustCreateBridge(&bta)
	require.NoError(t, err, "Creating bridge shouldn't fail")

	os := &client.DirectRequestTxPipelineSpec{
		BridgeTypeAttributes: bta,
		DataPath:             "data,result",
	}
	ost, err := os.String()
	require.NoError(t, err, "Building observation source spec shouldn't fail")

	_, err = chainlinkNodes[0].MustCreateJob(&client.DirectRequestJobSpec{
		Name:                     "direct_request",
		MinIncomingConfirmations: "1",
		ContractAddress:          oracle.Address(),
		ExternalJobID:            jobUUID.String(),
		ObservationSource:        ost,
	})
	require.NoError(t, err, "Creating direct_request job shouldn't fail")

	profileFunction := func(chainlinkNode *client.Chainlink) {
		if chainlinkNode != chainlinkNodes[len(chainlinkNodes)-1] {
			// Not the last node, hence not all nodes started profiling yet.
			return
		}
		jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
		var jobID [32]byte
		copy(jobID[:], jobUUIDReplaces)
		err = consumer.CreateRequestTo(
			oracle.Address(),
			jobID,
			big.NewInt(1e18),
			fmt.Sprintf("%s/variable", mockServerClient.Config.ClusterURL),
			"data,result",
			big.NewInt(100),
		)
		require.NoError(t, err, "Calling oracle contract shouldn't fail")

		gom := gomega.NewGomegaWithT(t)
		gom.Eventually(func(g gomega.Gomega) {
			d, err := consumer.Data(context.Background())
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Getting data from consumer contract shouldn't fail")
			g.Expect(d).ShouldNot(gomega.BeNil(), "Expected the initial on chain data to be nil")
			log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
			g.Expect(d.Int64()).Should(gomega.BeNumerically("==", 5), "Expected the on-chain data to be 5, but found %d", d.Int64())
		}, "2m", "1s").Should(gomega.Succeed())
	}

	profileTest := testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
		ProfileFunction: profileFunction,
		ProfileDuration: 30 * time.Second,
		ChainlinkNodes:  chainlinkNodes,
	})
	profileTest.Setup(testEnvironment)
	t.Cleanup(func() {
		CleanupPerformanceTest(t, testEnvironment, chainlinkNodes, profileTest.TestReporter, chainClient)
	})
	profileTest.Run()
}

func setupDirectRequestTest(t *testing.T) (testEnvironment *environment.Environment) {
	network := networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !network.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})
	}
	baseTOML := `[WebServer]
HTTPWriteTimout = '300s'`
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("performance-cron-%s", strings.ReplaceAll(strings.ToLower(network.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"toml": client.AddNetworksConfig(baseTOML, network),
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment
}
