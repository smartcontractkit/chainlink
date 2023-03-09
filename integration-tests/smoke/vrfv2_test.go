package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"go.uber.org/zap/zapcore"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"

	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func TestVRFv2Basic(t *testing.T) {
	linkEthFeedResponse := big.NewInt(1e18)
	minimumConfirmations := 3

	t.Parallel()
	testEnvironment, testNetwork := setupVRFv2Test(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err)
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err)
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})
	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err)
	bhs, err := contractDeployer.DeployBlockhashStore()
	require.NoError(t, err)
	mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
	require.NoError(t, err)
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkToken.Address(), bhs.Address(), mf.Address())
	require.NoError(t, err)
	consumer, err := contractDeployer.DeployVRFConsumerV2(linkToken.Address(), coordinator.Address())
	require.NoError(t, err)
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(1))
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	// https://docs.chain.link/docs/chainlink-vrf/#subscription-limits
	linkFunding := big.NewInt(100)

	err = linkToken.Transfer(consumer.Address(), big.NewInt(0).Mul(linkFunding, big.NewInt(1e18)))
	require.NoError(t, err)
	err = coordinator.SetConfig(
		uint16(minimumConfirmations),
		2.5e6,
		86400,
		33825,
		linkEthFeedResponse,
		ethereum.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: 1,
			FulfillmentFlatFeeLinkPPMTier2: 1,
			FulfillmentFlatFeeLinkPPMTier3: 1,
			FulfillmentFlatFeeLinkPPMTier4: 1,
			FulfillmentFlatFeeLinkPPMTier5: 1,
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40)},
	)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	err = consumer.CreateFundedSubscription(big.NewInt(0).Mul(linkFunding, big.NewInt(1e18)))
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	var (
		job                *client.Job
		encodedProvingKeys = make([][2]*big.Int, 0)
	)
	for _, n := range chainlinkNodes {
		vrfKey, err := n.MustCreateVRFKey()
		require.NoError(t, err)
		log.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.NewV4()
		os := &client.VRFV2TxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		require.NoError(t, err)
		oracleAddr, err := n.PrimaryEthAddress()
		require.NoError(t, err)
		job, err = n.MustCreateJob(&client.VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddresses:            []string{oracleAddr},
			EVMChainID:               fmt.Sprint(chainClient.GetNetworkConfig().ChainID),
			MinIncomingConfirmations: minimumConfirmations,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		require.NoError(t, err)
		provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
		require.NoError(t, err)
		err = coordinator.RegisterProvingKey(
			oracleAddr,
			provingKey,
		)
		require.NoError(t, err)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)
	}

	words := uint32(3)
	keyHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
	require.NoError(t, err)
	err = consumer.RequestRandomness(keyHash, 1, uint16(minimumConfirmations), 1000000, words)
	require.NoError(t, err)

	gom := gomega.NewGomegaWithT(t)
	timeout := time.Minute * 2
	gom.Eventually(func(g gomega.Gomega) {
		jobRuns, err := chainlinkNodes[0].MustReadRunsByJob(job.Data.ID)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
		randomness, err := consumer.GetAllRandomWords(context.Background(), int(words))
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		for _, w := range randomness {
			log.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
			g.Expect(w.Uint64()).ShouldNot(gomega.BeNumerically("==", 0), "Expected the VRF job give an answer other than 0")
		}
	}, timeout, "1s").Should(gomega.Succeed())
}

func setupVRFv2Test(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	testNetwork = networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	networkDetailTOML := `[EVM.GasEstimator]
LimitDefault = 3_500_000
PriceMax = 100000000000
FeeCapDefault = 100000000000`
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-vrfv2-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]any{
			"toml": client.AddNetworkDetailedConfig("", networkDetailTOML, testNetwork),
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
