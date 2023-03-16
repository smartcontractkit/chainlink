package reorg

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
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

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

const (
	EVMFinalityDepth         = "200"
	EVMTrackerHistoryDepth   = "400"
	reorgBlocks              = 50
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
		Timeout:                   blockchain.JSONStrDuration{2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func CleanupReorgTest(
	t *testing.T,
	testEnvironment *environment.Environment,
	chainlinkNodes []*client.Chainlink,
	chainClient blockchain.EVMClient,
) {
	if chainClient != nil {
		chainClient.GasStats().PrintStats()
	}
	err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.PanicLevel, chainClient)
	require.NoError(t, err, "Error tearing down environment")
}

func TestDirectRequestReorg(t *testing.T) {
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
	err = testEnvironment.AddHelm(chainlink.New(0, map[string]interface{}{
		"env": map[string]interface{}{
			"eth_url":                        "ws://geth-ethereum-geth:8546",
			"eth_http_url":                   "http://geth-ethereum-geth:8544",
			"eth_chain_id":                   "1337",
			"ETH_FINALITY_DEPTH":             EVMFinalityDepth,
			"ETH_HEAD_TRACKER_HISTORY_DEPTH": EVMTrackerHistoryDepth,
		},
	})).Run()
	require.NoError(t, err, "Error adding to test environment")

	chainClient, err := blockchain.NewEVMClient(networkSettings, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	cd, err := contracts.NewContractDeployer(chainClient)
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

	jobUUID := uuid.NewV4()

	bta := client.BridgeTypeAttributes{
		Name: fmt.Sprintf("five-%s", jobUUID.String()),
		URL:  fmt.Sprintf("%s/variable", ms.Config.ClusterURL),
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
		MinIncomingConfirmations: minIncomingConfirmations,
		ContractAddress:          oracle.Address(),
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
}
