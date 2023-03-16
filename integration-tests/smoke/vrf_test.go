package smoke

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func TestVRFBasic(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupVRFTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err, "Connecting client shouldn't fail")
	cd, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err, "Deploying contracts shouldn't fail")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Connecting to chainlink nodes shouldn't fail")
	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.ErrorLevel, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})
	chainClient.ParallelTransactions(true)

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(.01))
	require.NoError(t, err, "Funding chainlink nodes with ETH shouldn't fail")

	lt, err := cd.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	bhs, err := cd.DeployBlockhashStore()
	require.NoError(t, err, "Deploying Blockhash store shouldn't fail")
	coordinator, err := cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
	require.NoError(t, err, "Deploying VRF coordinator shouldn't fail")
	consumer, err := cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
	require.NoError(t, err, "Deploying VRF consumer contract shouldn't fail")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed to wait for VRF setup contracts to deploy")

	err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
	require.NoError(t, err, "Funding consumer contract shouldn't fail")
	_, err = cd.DeployVRFContract()
	require.NoError(t, err, "Deploying VRF contract shouldn't fail")
	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Waiting for event subscriptions in nodes shouldn't fail")

	for _, n := range chainlinkNodes {
		nodeKey, err := n.MustCreateVRFKey()
		require.NoError(t, err, "Creating VRF key shouldn't fail")
		log.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
		pubKeyCompressed := nodeKey.Data.ID
		jobUUID := uuid.NewV4()
		os := &client.VRFTxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		require.NoError(t, err, "Building observation source spec shouldn't fail")
		job, err := n.MustCreateJob(&client.VRFJobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			MinIncomingConfirmations: 1,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
		})
		require.NoError(t, err, "Creating VRF Job shouldn't fail")

		oracleAddr, err := n.PrimaryEthAddress()
		require.NoError(t, err, "Getting primary ETH address of chainlink node shouldn't fail")
		provingKey, err := actions.EncodeOnChainVRFProvingKey(*nodeKey)
		require.NoError(t, err, "Encoding on-chain VRF Proving key shouldn't fail")
		err = coordinator.RegisterProvingKey(
			big.NewInt(1),
			oracleAddr,
			provingKey,
			actions.EncodeOnChainExternalJobID(jobUUID),
		)
		require.NoError(t, err, "Registering the on-chain VRF Proving key shouldn't fail")
		encodedProvingKeys := make([][2]*big.Int, 0)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)

		requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
		require.NoError(t, err, "Getting Hash of encoded proving keys shouldn't fail")
		err = consumer.RequestRandomness(requestHash, big.NewInt(1))
		require.NoError(t, err, "Requesting randomness shouldn't fail")

		gom := gomega.NewGomegaWithT(t)
		timeout := time.Minute * 2
		gom.Eventually(func(g gomega.Gomega) {
			jobRuns, err := chainlinkNodes[0].MustReadRunsByJob(job.Data.ID)
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Job execution shouldn't fail")

			out, err := consumer.RandomnessOutput(context.Background())
			g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Getting the randomness output of the consumer shouldn't fail")
			// Checks that the job has actually run
			g.Expect(len(jobRuns.Data)).Should(gomega.BeNumerically(">=", 1),
				fmt.Sprintf("Expected the VRF job to run once or more after %s", timeout))

			// TODO: This is an imperfect check, given it's a random number, it CAN be 0, but chances are unlikely.
			// So we're just checking that the answer has changed to something other than the default (0)
			// There's a better formula to ensure that VRF response is as expected, detailed under Technical Walkthrough.
			// https://blog.chain.link/chainlink-vrf-on-chain-verifiable-randomness/
			g.Expect(out.Uint64()).ShouldNot(gomega.BeNumerically("==", 0), "Expected the VRF job give an answer other than 0")
			log.Debug().Uint64("Output", out.Uint64()).Msg("Randomness fulfilled")
		}, timeout, "1s").Should(gomega.Succeed())
	}
}

func setupVRFTest(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	testNetwork = networks.SelectedNetwork
	evmConfig := ethereum.New(nil)
	if !testNetwork.Simulated {
		evmConfig = ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}
	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-vrf-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"toml": client.AddNetworksConfig("", testNetwork),
		}))
	err := testEnvironment.Run()
	require.NoError(t, err, "Error running test environment")
	return testEnvironment, testNetwork
}
