package reorg

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/logging"
	"github.com/smartcontractkit/chainlink-env/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	baseDRTOML = `
[Feature]
LogPoller = true
`
	networkDRTOML = `
Enabled = true
FinalityDepth = %s
LogPollInterval = '1s'

[EVM.HeadTracker]
HistoryDepth = %s
`
)

const (
	EVMFinalityDepth         = "200"
	EVMTrackerHistoryDepth   = "400"
	reorgBlocks              = 10
	minIncomingConfirmations = "200"
	timeout                  = "15m"
	interval                 = "2s"
)

var (
	networkSettings = blockchain.EVMNetwork{
		Name:      "geth",
		Simulated: true,
		ChainID:   1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func CleanupReorgTest(
	t *testing.T,
	testEnvironment *environment.Environment,
	chainlinkNodes []*client.ChainlinkK8sClient,
	chainClient blockchain.EVMClient,
) {
	if chainClient != nil {
		chainClient.GasStats().PrintStats()
	}
	err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.PanicLevel, chainClient)
	require.NoError(t, err, "Error tearing down environment")
}

func TestDirectRequestReorg(t *testing.T) {
	logging.Init()
	l := logging.GetTestLogger(t)
	testEnvironment := environment.New(&environment.Config{
		TTL:  1 * time.Hour,
		Test: t,
	})
	err := testEnvironment.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddChart(blockscout.New(&blockscout.Props{
			WsURL:   "ws://geth-ethereum-geth:8546",
			HttpURL: "http://geth-ethereum-geth:8544",
		})).
		AddHelm(reorg.New(&reorg.Props{
			NetworkName: "geth",
			NetworkType: "geth-reorg",
			Values: map[string]interface{}{
				"geth": map[string]interface{}{
					"genesis": map[string]interface{}{
						"networkId": "1337",
					},
				},
			},
		})).
		Run()
	require.NoError(t, err, "Error deploying test environment")
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	// related https://app.shortcut.com/chainlinklabs/story/38295/creating-an-evm-chain-via-cli-or-api-immediately-polling-the-nodes-and-returning-an-error
	// node must work and reconnect even if network is not working
	time.Sleep(90 * time.Second)
	activeEVMNetwork = networks.SimulatedEVMNonDev
	netCfg := fmt.Sprintf(networkDRTOML, EVMFinalityDepth, EVMTrackerHistoryDepth)
	chainlinkDeployment := chainlink.New(0, map[string]interface{}{
		"replicas": 1,
		"toml":     client.AddNetworkDetailedConfig(baseDRTOML, netCfg, activeEVMNetwork),
	})

	err = testEnvironment.AddHelm(chainlinkDeployment).Run()
	require.NoError(t, err, "Error adding to test environment")

	chainClient, err := blockchain.NewEVMClient(networkSettings, testEnvironment, l)
	require.NoError(t, err, "Error connecting to blockchain")
	cd, err := contracts.NewContractDeployer(chainClient, l)
	require.NoError(t, err, "Error building contract deployer")
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Chainlink nodes")
	ms, err := ctfClient.ConnectMockServer(testEnvironment)
	require.NoError(t, err, "Error connecting to Mockserver")

	t.Cleanup(func() {
		CleanupReorgTest(t, testEnvironment, chainlinkNodes, chainClient)
	})

	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, big.NewFloat(10))
	require.NoError(t, err, "Error funding Chainlink nodes")

	lt, err := cd.DeployLinkTokenContract()
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")
	oracle, err := cd.DeployOracle(lt.Address())
	require.NoError(t, err, "Deploying Oracle Contract shouldn't fail")
	consumer, err := cd.DeployAPIConsumer(lt.Address())
	require.NoError(t, err, "Deploying Consumer Contract shouldn't fail")
	err = chainClient.SetDefaultWallet(0)
	require.NoError(t, err, "Setting default wallet shouldn't fail")
	err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
	require.NoError(t, err, "Transferring %d to consumer contract shouldn't fail", big.NewInt(2e18))

	err = ms.SetValuePath("/variable", 5)
	require.NoError(t, err, "Setting mockserver value path shouldn't fail")

	jobUUID := uuid.New()

	bta := client.BridgeTypeAttributes{
		Name: fmt.Sprintf("five-%s", jobUUID.String()),
		URL:  fmt.Sprintf("%s/variable", ms.Config.ClusterURL),
	}
	err = chainlinkNodes[0].MustCreateBridge(&bta)
	require.NoError(t, err, "Creating bridge shouldn't fail")

	drps := &client.DirectRequestTxPipelineSpec{
		BridgeTypeAttributes: bta,
		DataPath:             "data,result",
	}
	ost, err := drps.String()
	require.NoError(t, err, "Building observation source spec shouldn't fail")

	_, err = chainlinkNodes[0].MustCreateJob(&client.DirectRequestJobSpec{
		Name:                     "direct_request",
		MinIncomingConfirmations: minIncomingConfirmations,
		ContractAddress:          oracle.Address(),
		EVMChainID:               chainClient.GetChainID().String(),
		ExternalJobID:            jobUUID.String(),
		ObservationSource:        ost,
	})
	require.NoError(t, err, "Creating direct_request job shouldn't fail")

	rc, err := NewReorgController(
		&ReorgConfig{
			FromPodLabel:            reorg.TXNodesAppLabel,
			ToPodLabel:              reorg.MinerNodesAppLabel,
			Network:                 chainClient,
			Env:                     testEnvironment,
			BlockConsensusThreshold: 3,
			Timeout:                 1800 * time.Second,
		},
	)
	require.NoError(t, err, "Error getting reorg controller")
	rc.ReOrg(reorgBlocks)
	rc.WaitReorgStarted()

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
	require.NoError(t, err, "Calling oracle contract shouldn't fail")

	err = rc.WaitDepthReached()
	require.NoError(t, err, "Error waiting for depth to be reached")

	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		d, err := consumer.Data(context.Background())
		g.Expect(err).ShouldNot(gomega.HaveOccurred(), "Getting data from consumer contract shouldn't fail")
		g.Expect(d).ShouldNot(gomega.BeNil(), "Expected the initial on chain data to be nil")
		log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
		g.Expect(d.Int64()).Should(gomega.BeNumerically("==", 5), "Expected the on-chain data to be 5, but found %d", d.Int64())
	}, timeout, interval).Should(gomega.Succeed())
}
