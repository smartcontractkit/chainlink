package soak

//import (
//	"context"
//	"fmt"
//	"github.com/onsi/gomega"
//	uuid "github.com/satori/go.uuid"
//	ctf_ethereum "github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
//	"github.com/smartcontractkit/chainlink-testing-framework/utils"
//	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
//	"math/big"
//	"os"
//	"strings"
//	"testing"
//	"time"
//
//	"github.com/kelseyhightower/envconfig"
//	"github.com/stretchr/testify/require"
//
//	"github.com/smartcontractkit/chainlink-env/environment"
//	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
//	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
//	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
//
//	networks "github.com/smartcontractkit/chainlink/integration-tests"
//	"github.com/smartcontractkit/chainlink/integration-tests/actions"
//	"github.com/smartcontractkit/chainlink/integration-tests/client"
//	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
//)
//
//func TestVRFV2Soak(t *testing.T) {
//	l := utils.GetTestLogger(t)
//	testEnvironment, testNetwork, testInputs := SetupVRFV2SoakEnv(t)
//	if testEnvironment.WillUseRemoteRunner() {
//		return
//	}
//
//	//##################################################
//	linkEthFeedResponse := big.NewInt(1e18)
//	minimumConfirmations := 3
//	subID := uint64(1)
//	linkFundingAmount := big.NewInt(100)
//	numberOfWords := uint32(3)
//
//	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
//	require.NoError(t, err)
//	contractDeployer, err := contracts.NewContractDeployer(chainClient)
//	require.NoError(t, err)
//	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
//	require.NoError(t, err)
//	//t.Cleanup(func() {
//	//	err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
//	//	require.NoError(t, err, "Error tearing down environment")
//	//})
//	chainClient.ParallelTransactions(true)
//
//	linkToken, err := contractDeployer.DeployLinkTokenContract()
//	require.NoError(t, err)
//	bhs, err := contractDeployer.DeployBlockhashStore()
//	require.NoError(t, err)
//	mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
//	require.NoError(t, err)
//	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkToken.Address(), bhs.Address(), mf.Address())
//	require.NoError(t, err)
//
//	consumer, err := contractDeployer.DeployVRFv2Consumer(coordinator.Address())
//	require.NoError(t, err)
//	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(1))
//	require.NoError(t, err)
//	err = chainClient.WaitForEvents()
//	require.NoError(t, err)
//
//	err = coordinator.SetConfig(
//		uint16(minimumConfirmations),
//		2.5e6,
//		86400,
//		33825,
//		linkEthFeedResponse,
//		ctf_ethereum.VRFCoordinatorV2FeeConfig{
//			FulfillmentFlatFeeLinkPPMTier1: 1,
//			FulfillmentFlatFeeLinkPPMTier2: 1,
//			FulfillmentFlatFeeLinkPPMTier3: 1,
//			FulfillmentFlatFeeLinkPPMTier4: 1,
//			FulfillmentFlatFeeLinkPPMTier5: 1,
//			ReqsForTier2:                   big.NewInt(10),
//			ReqsForTier3:                   big.NewInt(20),
//			ReqsForTier4:                   big.NewInt(30),
//			ReqsForTier5:                   big.NewInt(40)},
//	)
//	require.NoError(t, err)
//	err = chainClient.WaitForEvents()
//	require.NoError(t, err)
//
//	err = coordinator.CreateSubscription()
//	require.NoError(t, err)
//	err = chainClient.WaitForEvents()
//	require.NoError(t, err)
//
//	err = coordinator.AddConsumer(subID, consumer.Address())
//	require.NoError(t, err, "Error adding a consumer to a subscription in VRFCoordinator contract")
//
//	actions.FundVRFCoordinatorV2Subscription(
//		t,
//		linkToken,
//		coordinator,
//		chainClient,
//		subID,
//		linkFundingAmount,
//	)
//
//	var (
//		job                *client.Job
//		encodedProvingKeys = make([][2]*big.Int, 0)
//		oracleAddress      string
//	)
//	for _, n := range chainlinkNodes {
//		vrfKey, err := n.MustCreateVRFKey()
//		require.NoError(t, err)
//		l.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
//		pubKeyCompressed := vrfKey.Data.ID
//		jobUUID := uuid.NewV4()
//		os := &client.VRFV2TxPipelineSpec{
//			Address: coordinator.Address(),
//		}
//		ost, err := os.String()
//		require.NoError(t, err)
//		oracleAddress, err = n.PrimaryEthAddress()
//		require.NoError(t, err)
//		job, err = n.MustCreateJob(&client.VRFV2JobSpec{
//			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
//			CoordinatorAddress:       coordinator.Address(),
//			FromAddresses:            []string{oracleAddress},
//			EVMChainID:               fmt.Sprint(chainClient.GetNetworkConfig().ChainID),
//			MinIncomingConfirmations: minimumConfirmations,
//			PublicKey:                pubKeyCompressed,
//			ExternalJobID:            jobUUID.String(),
//			ObservationSource:        ost,
//			BatchFulfillmentEnabled:  false,
//		})
//		require.NoError(t, err)
//		provingKey, err := actions.EncodeOnChainVRFProvingKey(*vrfKey)
//		require.NoError(t, err)
//		err = coordinator.RegisterProvingKey(
//			oracleAddress,
//			provingKey,
//		)
//		require.NoError(t, err)
//		encodedProvingKeys = append(encodedProvingKeys, provingKey)
//	}
//
//	keyHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
//	require.NoError(t, err)
//
//	//##################################################
//
//	//chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
//	//require.NoError(t, err, "Error connecting to testNetwork")
//
//	//contractLoader, err := contracts.NewContractLoader(client)
//
//	vrfV2SoakTest := testsetups.NewVRFV2SoakTest(&testsetups.VRFV2SoakTestInputs{
//		BlockchainClient:     chainClient,
//		TestDuration:         testInputs.TestDuration,
//		ChainlinkNodeFunding: testInputs.ChainlinkNodeFunding,
//		SubscriptionFunding:  testInputs.SubscriptionFunding,
//		StopTestOnError:      testInputs.StopTestOnError,
//		RequestsPerMinute:    testInputs.RequestsPerMinute,
//		TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int) error {
//
//			// request randomness
//			err := consumer.RequestRandomness(keyHash, 1, 1, 300000, numberOfWords)
//			if err != nil {
//				return err
//			}
//
//			//todo - how to make assertions in soak test?
//			gom := gomega.NewGomegaWithT(t)
//			timeout := time.Minute * 2
//			var lastRequestID *big.Int
//			gom.Eventually(func(g gomega.Gomega) {
//				jobRuns, err := chainlinkNodes[0].MustReadRunsByJob(job.Data.ID)
//				g.Expect(err).ShouldNot(gomega.HaveOccurred())
//				g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically("==", 1))
//				lastRequestID, err = consumer.GetLastRequestId(context.Background())
//				l.Debug().Interface("Last Request ID", lastRequestID).Msg("Last Request ID Received")
//
//				g.Expect(err).ShouldNot(gomega.HaveOccurred())
//				status, err := consumer.GetRequestStatus(context.Background(), lastRequestID)
//				g.Expect(err).ShouldNot(gomega.HaveOccurred())
//				g.Expect(status.Fulfilled).Should(gomega.BeTrue())
//				l.Debug().Interface("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
//
//				g.Expect(err).ShouldNot(gomega.HaveOccurred())
//				for _, w := range status.RandomWords {
//					l.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
//					g.Expect(w.Uint64()).Should(gomega.BeNumerically(">", 0), "Expected the VRF job give an answer bigger than 0")
//				}
//			}, timeout, "1s").Should(gomega.Succeed())
//
//			return nil
//		},
//	})
//
//	//				// load the previously deployed consumer contract
//	//				evmClient := vrfv2SoakTest.Inputs.BlockchainClient.(*blockchain.EthereumMultinodeClient).DefaultClient
//	//				tempConsumer := &contracts.EthereumVRFConsumerV2{}
//	//				consumerAddress := os.Getenv("CONSUMER_ADDRESS")
//	//				err = tempConsumer.LoadExistingConsumer(consumerAddress, evmClient)
//	//				Expect(err).ShouldNot(HaveOccurred())
//	//				consumer = tempConsumer
//	//
//	//				// load the previously deployed link token contract
//	//				tempLinkToken := &contracts.EthereumLinkToken{}
//	//				linkTokenContractAddress := os.Getenv("LINK_TOKEN_CONTRACT_ADDRESS")
//	//				err = tempLinkToken.LoadExistingLinkToken(linkTokenContractAddress, evmClient)
//	//				Expect(err).ShouldNot(HaveOccurred())
//	//				linkTokenContract = tempLinkToken
//
//	//todo - how to add report?
//	t.Cleanup(func() {
//		if err := actions.TeardownRemoteSuite(vrfV2SoakTest.TearDownVals(t)); err != nil {
//			l.Error().Err(err).Msg("Error tearing down environment")
//		}
//	})
//	//todo - what should be in Setup?
//	vrfV2SoakTest.Setup(t, testEnvironment)
//	l.Info().Msg("Set up soak test")
//	vrfV2SoakTest.Run(t)
//}
//
//func SetupVRFV2SoakEnv(t *testing.T) (*environment.Environment, blockchain.EVMNetwork, testsetups.VRFV2SoakTestInputs) {
//	var testInputs testsetups.VRFV2SoakTestInputs
//	err := envconfig.Process("VRFV2", &testInputs)
//	require.NoError(t, err, "Error reading OCR soak test inputs")
//	testInputs.SetForRemoteRunner()
//
//	network := networks.SelectedNetwork // Environment currently being used to soak test on
//	var vrfV2EnvVars = map[string]any{
//		"P2P_LISTEN_IP":   "0.0.0.0",
//		"P2P_LISTEN_PORT": "6690",
//	}
//	// For if we end up using env vars
//	vrfV2EnvVars["ETH_URL"] = network.URLs[0]
//	vrfV2EnvVars["ETH_HTTP_URL"] = network.HTTPURLs[0]
//	vrfV2EnvVars["ETH_CHAIN_ID"] = fmt.Sprint(network.ChainID)
//
//	baseEnvironmentConfig := &environment.Config{
//		TTL: time.Hour * 720, // 30 days,
//		NamespacePrefix: fmt.Sprintf(
//			"soak-vrfv2-%s",
//			strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"),
//		),
//		Test: t,
//	}
//
//	testEnvironment := environment.New(baseEnvironmentConfig).
//		AddHelm(ethereum.New(&ethereum.Props{
//			NetworkName: network.Name,
//			Simulated:   network.Simulated,
//			WsURLs:      network.URLs,
//		}))
//	useEnvVars := strings.ToLower(os.Getenv("TEST_USE_ENV_VAR_CONFIG"))
//	if useEnvVars == "true" {
//		testEnvironment.AddHelm(chainlink.NewVersioned(0, "0.0.11", map[string]any{
//			"env": vrfV2EnvVars,
//		}))
//	} else {
//		networkDetailTOML := `[EVM.GasEstimator]
//			LimitDefault = 3_500_000
//			PriceMax = 100000000000
//			FeeCapDefault = 100000000000`
//		testEnvironment.AddHelm(chainlink.New(0, map[string]any{
//			"toml": client.AddNetworkDetailedConfig("", networkDetailTOML, network),
//		}))
//	}
//	err = testEnvironment.Run()
//	require.NoError(t, err, "Error launching test environment")
//	return testEnvironment, network, testInputs
//}
