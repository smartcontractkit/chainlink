package smoke

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet"
	"github.com/smartcontractkit/chainlink/integration-tests/gauntlet/configs"
	"math/big"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

const (
	ErrWatchingNewOCRRound = "Error watching for new OCR round"
)

func TestOCRBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, ocrInstances, sethClient := prepareORCv1SmokeTestEnv(t, l, 5)
	nodeClients := env.ClCluster.NodeAPIs()
	workerNodes := nodeClients[1:]

	err := actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err, "Error setting all adapter responses to the same value")

	err = actions_seth.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, ErrWatchingNewOCRRound)

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}

func TestOCRJobReplacement(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	env, ocrInstances, sethClient := prepareORCv1SmokeTestEnv(t, l, 5)
	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	err := actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err, "Error setting all adapter responses to the same value")
	err = actions_seth.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, ErrWatchingNewOCRRound)

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())

	err = actions.DeleteJobs(nodeClients)
	require.NoError(t, err, "Error deleting OCR jobs")

	err = actions.DeleteBridges(nodeClients)
	require.NoError(t, err, "Error deleting OCR bridges")

	//Recreate job
	err = actions.CreateOCRJobsLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, big.NewInt(sethClient.ChainID))
	require.NoError(t, err, "Error creating OCR jobs")

	err = actions_seth.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, ErrWatchingNewOCRRound)

	answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}

func prepareORCv1SmokeTestEnv(t *testing.T, l zerolog.Logger, firstRoundResult int64) (*test_env.CLClusterTestEnv, []contracts.OffchainAggregator, *seth.Client) {
	config, err := tc.GetConfig("Smoke", tc.OCR)
	if err != nil {
		t.Fatal(err)
	}

	network, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(network.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(6).
		WithFunding(big.NewFloat(.5)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	selectedNetwork := networks.MustGetSelectedNetworkConfig(config.Network)[0]
	sethClient, err := env.GetSethClient(selectedNetwork.ChainID)
	require.NoError(t, err, "Error getting seth client")

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	linkContract, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Error deploying link token contract")

	ocrInstances, err := actions_seth.DeployOCRv1Contracts(l, sethClient, 1, common.HexToAddress(linkContract.Address()), contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes))
	require.NoError(t, err, "Error deploying OCR contracts")

	err = actions.CreateOCRJobsLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, big.NewInt(sethClient.ChainID))
	require.NoError(t, err, "Error creating OCR jobs")

	err = actions_seth.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(3*time.Minute))
	require.NoError(t, err, "Error watching for new OCR round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, firstRoundResult, answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	return env, ocrInstances, sethClient
}

type ZKSyncState struct {
	Gauntlet        *gauntlet.Gauntlet
	ChainlinkClient []*client.ChainlinkClient
	ContractLoader  contracts.ContractLoader
	OCRContract     []contracts.OffchainAggregator
	L2RPC           string
}

// If we go this way this test has to be removed and unified with the other ones. Basically based on the selected networks we would either use one or the other action to deploy contracts.
// We'd either gauntlet for zksync or what we have now for everything else.
func TestOCRZkSync(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.OCR)
	if err != nil {
		t.Fatal(err)
	}

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithMockAdapter().
		WithCLNodes(6).
		WithFunding(big.NewFloat(0.1)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	selectedNetwork := networks.MustGetSelectedNetworkConfig(config.Network)[0]
	sethClient, err := env.GetSethClient(selectedNetwork.ChainID)
	require.NoError(t, err, "Error getting seth client")

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	_, b, _, _ := runtime.Caller(0)
	// ProjectRoot Root folder of this project
	ProjectRoot := filepath.Join(filepath.Dir(b), "/../..")
	// SolanaTestsRoot path to starknet e2e tests
	IntegrationTestsRoot := filepath.Join(ProjectRoot, "integration-tests")

	g, err := gauntlet.New("gauntlet:zksync", fmt.Sprintf("%s/", IntegrationTestsRoot))
	require.NoError(t, err)
	testState := &ZKSyncState{
		Gauntlet:        g,
		ChainlinkClient: nil,
		ContractLoader:  nil,
		L2RPC:           "",
	}
	err = testState.Gauntlet.SetupNetwork(config.Network.RpcHttpUrls[config.Network.SelectedNetworks[0]][0], config.Network.WalletKeys[config.Network.SelectedNetworks[0]][0])
	require.NoError(t, err, "Setting up gauntlet network should not fail")

	testNetwork := networks.MustGetSelectedNetworkConfig(config.Network)[0]
	nKeys, _, err := client.CreateNodeKeysBundle(nodeClients, testNetwork.Name, strconv.FormatInt(testNetwork.ChainID, 10))
	require.NoError(t, err)

	nKeys = nKeys[1:]

	chainClient, err := env.GetSethClientForSelectedNetwork()
	require.NoError(t, err)

	var transmitters []string
	var payees []string
	var signers []string
	var peerIds []string
	var ocrConfigPubKeys []string

	for _, key := range nKeys {
		peerIds = append(peerIds, key.PeerID)
		ocrConfigPubKeys = append(ocrConfigPubKeys, strings.Replace(key.OCRKeys.Data[0].Attributes.OffChainPublicKey, "ocroff_", "", 1))
		transmitters = append(transmitters, strings.Replace(key.EthAddress, "0x", "", 1))
		signers = append(signers, strings.Replace(key.OCRKeys.Data[0].Attributes.OnChainSigningAddress, "ocrsad_", "", 1))
		payees = append(payees, strings.Replace(chainClient.MustGetRootKeyAddress().Hex(), "0x", "", 1))
	}

	err = testState.Gauntlet.DeployLinkToken()
	require.NoError(t, err)

	testState.Gauntlet.Contracts.LinkContract.Contract, err = contracts.LoadLinkTokenContract(l, chainClient, common.HexToAddress(testState.Gauntlet.Contracts.LinkContract.Address))
	require.NoError(t, err)

	err = testState.Gauntlet.DeployAccessController()
	require.NoError(t, err)

	ocrConfig := configs.OCRConfig{}
	ocrConfig.DefaultOcrConfig()
	ocrConfig.DefaultOcrContract()

	ocrConfig.Contract.Link = testState.Gauntlet.Contracts.LinkContract.Address
	ocrConfig.Contract.RequesterAccessController = testState.Gauntlet.Contracts.AccessControllerAddress
	ocrConfig.Contract.BillingAccessController = testState.Gauntlet.Contracts.AccessControllerAddress
	ocrConfig.Config.Signers = signers
	ocrConfig.Config.Transmitters = transmitters
	ocrConfig.Config.OcrConfigPublicKeys = ocrConfigPubKeys
	ocrConfig.Config.OperatorsPeerIds = strings.Join(peerIds, ",")
	ocrJsonContract, err := ocrConfig.MarshalContract()

	err = testState.Gauntlet.DeployOCR(ocrJsonContract)
	require.NoError(t, err)

	ocrContract, err := contracts.LoadOffchainAggregator(l, chainClient, common.HexToAddress(testState.Gauntlet.Contracts.OCRContract.Address))
	require.NoError(t, err)
	testState.Gauntlet.Contracts.OCRContract.Contract = &ocrContract

	err = testState.Gauntlet.AddAccess(testState.Gauntlet.Contracts.OCRContract.Address)
	require.NoError(t, err)

	err = testState.Gauntlet.SetPayees(testState.Gauntlet.Contracts.OCRContract.Address, payees, transmitters)
	require.NoError(t, err)

	ocrJsonConfig, err := ocrConfig.MarshalConfig()
	require.NoError(t, err)

	err = testState.Gauntlet.SetConfig(testState.Gauntlet.Contracts.OCRContract.Address, ocrJsonConfig)
	require.NoError(t, err)

	testState.OCRContract = []contracts.OffchainAggregator{
		testState.Gauntlet.Contracts.OCRContract.Contract,
	}
	var chainlinkClients []*client.ChainlinkClient
	for _, k8sClient := range nodeClients {
		chainlinkClients = append(chainlinkClients, k8sClient)
	}

	var transmitterAddresses []common.Address
	for _, addrStr := range transmitters {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress("0x"+addrStr))
	}

	// Exclude the first node, which is the bootstrap node
	err = testState.OCRContract[0].SetConfig(
		contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes),
		contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
		transmitterAddresses,
	)
	require.NoError(t, err)

	firstRoundResult := 5

	err = actions.CreateOCRJobsLocal(testState.OCRContract, bootstrapNode, workerNodes, firstRoundResult, env.MockAdapter, big.NewInt(sethClient.ChainID))
	require.NoError(t, err, "Error creating OCR jobs")

	err = actions_seth.StartNewRound(contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(testState.OCRContract))
	require.NoError(t, err)

	err = actions_seth.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(testState.OCRContract), 3*time.Minute)
	require.NoError(t, err, "Error watching for new OCR round")

	answer, err := testState.OCRContract[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(firstRoundResult), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())
}
