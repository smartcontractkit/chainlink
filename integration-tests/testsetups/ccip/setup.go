package ccip

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	cl "github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const DefaultChainNamePrefix = "chain-"

type NetworkSetup struct {
	NumberOfNetworks int      `toml:"number_of_networks" validate:"required"`
	PrivateKeys      []string `toml:"private_keys" validate:"required"`
}

type CommonConfig struct {
	NodeFundingAmount float64 `toml:"chainlink_node_funding" validate:"required"`
}

type CTFV2Conf struct {
	Common             CommonConfig        `toml:"Common" validate:"required"`
	Network            NetworkSetup        `toml:"Network" validate:"required"`
	BlockchainNetworks []*blockchain.Input `toml:"Networks" validate:"required"`
	NodeSet            *ns.Input           `toml:"nodeset" validate:"required"`
	JDDbInput          *postgres.Input     `toml:"jd_db" validate:"required"`
	JD                 *jd.Input           `toml:"jd" validate:"required"`
	Fake               *fake.Input         `toml:"fake" validate:"required"`
}

// DeployedLocalAnvilDevEnvironment is a helper struct for setting up a local anvil dev environment with docker using ctf v2
type DeployedLocalAnvilDevEnvironment struct {
	changeset.DeployedEnv
	bcs           []*blockchain.Output
	DON           *devenv.DON
	devEnvTestCfg tc.TestConfig
	devEnvCfg     *devenv.EnvironmentConfig
	in            *CTFV2Conf
	pvtKeys       []string
	MockAdapter   *fake.Output
}

func (l *DeployedLocalAnvilDevEnvironment) DeployedEnvironment() changeset.DeployedEnv {
	return l.DeployedEnv
}

func (l *DeployedLocalAnvilDevEnvironment) StartChains(t *testing.T, _ *changeset.TestConfigs) {
	lggr := logger.TestLogger(t)
	ctx := testcontext.Get(t)
	envConfig, cfg, bcs, in := createAnvilDockerNetwork(t)
	l.pvtKeys = in.Network.PrivateKeys
	l.devEnvTestCfg = cfg
	l.bcs = bcs
	l.devEnvCfg = envConfig
	l.in = in
	users := make(map[uint64][]*bind.TransactOpts)
	for _, chain := range envConfig.Chains {
		details, found := chainsel.ChainByEvmChainID(chain.ChainID)
		require.Truef(t, found, "chain not found")
		users[details.Selector] = chain.Users
	}
	homeChainSel := l.devEnvTestCfg.CCIP.GetHomeChainSelector()
	require.NotEmpty(t, homeChainSel, "homeChainSel should not be empty")
	feedSel := l.devEnvTestCfg.CCIP.GetFeedChainSelector()
	require.NotEmpty(t, feedSel, "feedSel should not be empty")
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	require.NoError(t, err)
	replayBlocks, err := changeset.LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)
	l.DeployedEnv.Users = users
	l.DeployedEnv.Env.Chains = chains
	l.DeployedEnv.FeedChainSel = feedSel
	l.DeployedEnv.HomeChainSel = homeChainSel
	l.DeployedEnv.ReplayBlocks = replayBlocks
}

func (l *DeployedLocalAnvilDevEnvironment) StartNodes(t *testing.T, _ *changeset.TestConfigs, crConfig deployment.CapabilityRegistryConfig) {
	require.NotEmpty(t, l.devEnvTestCfg, "integration test config is empty, start chains first")
	require.NotNil(t, l.devEnvCfg, "dev environment config is empty, start chains first")
	nodeOut := startCLNodes(t, crConfig, l.bcs, l.in)
	ctx := testcontext.Get(t)
	lggr := logger.TestLogger(t)
	nodeInfo, err := getNodeInfo(nodeOut, pointer.GetInt(l.devEnvTestCfg.CCIP.CLNode.NoOfBootstraps))
	require.NoError(t, err, "failed to get node info")
	l.devEnvCfg.JDConfig.NodeInfo = nodeInfo
	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, *l.devEnvCfg)
	require.NoError(t, err, "failed to create dev environment")
	require.NotNil(t, e)
	l.DON = don
	l.DeployedEnv.Env = *e

	// fund the nodes
	for _, chain := range l.bcs {
		scSrc, err := seth.NewClientBuilder().
			WithRpcUrl(chain.Nodes[0].HostWSUrl).
			WithGasPriceEstimations(true, 0, seth.Priority_Fast).
			WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
			WithPrivateKeys(l.pvtKeys).
			Build()
		require.NoError(t, err, "failed to create source chain seth client")
		nodeClients, err := cl.New(nodeOut.CLNodes)
		require.NoError(t, err, "failed to create node clients")
		err = ns.FundNodes(scSrc.Client, nodeClients, l.pvtKeys[0], l.in.Common.NodeFundingAmount)
		require.NoError(t, err, "failed to fund nodes")
	}
}

func getNodeInfo(nodeOut *ns.Output, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.DockerP2PUrl)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse p2p url")
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("bootstrap-%d", i),
				P2PPort:     p2pURL.Port(),
				CLConfig: clclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("node-%d", i),
				P2PPort:     p2pURL.Port(),
				CLConfig: clclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		}
	}
	return nodeInfo, nil
}

func (l *DeployedLocalAnvilDevEnvironment) MockUSDCAttestationServer(t *testing.T, isFaulty bool) string {
	// SetMockServerWithUSDCAttestation responds with a mock attestation for any msgHash
	// The path is set with regex to match any path that starts with /v1/attestations
	path := "/v1/attestations"
	fakeOut, err := fake.NewFakeDataProvider(l.in.Fake)
	require.NoError(t, err, "failed to create mock data provider")
	l.MockAdapter = fakeOut
	err = fake.JSON(
		"GET",
		fmt.Sprintf("%s/:hash", path),
		map[string]any{
			"status":      "complete",
			"attestation": "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b",
		}, 200,
	)
	require.NoError(t, err, "failed to set up mock usdc attestation server")
	if isFaulty {
		err = fake.JSON(
			"GET",
			fmt.Sprintf("%s/:hash", path),
			map[string]any{
				"status": "pending",
				"error":  "internal error",
			}, 200,
		)
		require.NoError(t, err, "failed to set up mock faulty usdc attestation server")
	}
	return fakeOut.BaseURLDocker
}

func createAnvilDockerNetwork(t *testing.T) (
	*devenv.EnvironmentConfig,
	tc.TestConfig,
	[]*blockchain.Output,
	*CTFV2Conf,
) {
	in, err := framework.Load[CTFV2Conf](t)
	require.NoError(t, err, "failed to load ccip test config")
	if in.Network.NumberOfNetworks > len(in.BlockchainNetworks) {
		// if number of networks is greater than number of blockchain networks, dynamically add the required networks
		networksNeeded := in.Network.NumberOfNetworks - len(in.BlockchainNetworks)
		finalPortID, err := strconv.Atoi(in.BlockchainNetworks[len(in.BlockchainNetworks)-1].Port)
		require.NoError(t, err, "failed to convert port number to int")
		for i := 1; i <= networksNeeded; i++ {
			in.BlockchainNetworks = append(in.BlockchainNetworks, &blockchain.Input{
				ChainID: strconv.Itoa(90000000 + i),
				Type:    "anvil",
				Port:    strconv.Itoa(finalPortID + i),
			})
		}
	}
	// spin up 2 anvils
	var blockchains []*blockchain.Output
	for c := 0; c < in.Network.NumberOfNetworks; c++ {
		bc, err := blockchain.NewBlockchainNetwork(in.BlockchainNetworks[c])
		require.NoError(t, err, "failed to create blockchain network")
		blockchains = append(blockchains, bc)
		// mine periodically if not overridden
		//TODO: Need to handle this check in a better way instead of checking for nil
		if in.BlockchainNetworks[c].DockerCmdParamsOverrides == nil {
			miner := rpc.NewRemoteAnvilMiner(bc.Nodes[0].HostHTTPUrl, nil)
			miner.MinePeriodically(1 * time.Second)
		}
	}
	jdOutput, err := jd.NewJD(in.JD)
	require.NoError(t, err, "failed to create new job distributor")
	jdConfig := devenv.JDConfig{
		GRPC:  jdOutput.HostGRPCUrl,
		WSRPC: jdOutput.DockerWSRPCUrl,
		Creds: insecure.NewCredentials(),
	}
	chains := make([]devenv.ChainConfig, 0, len(blockchains))
	for _, chain := range blockchains {
		chainID, err := strconv.ParseInt(chain.ChainID, 10, 64)
		require.NoError(t, err, "invalid chain id")
		id, err := strconv.ParseUint(chain.ChainID, 10, 64)
		require.NoError(t, err, "invalid chain id")
		pvtKey, err := crypto.HexToECDSA(in.Network.PrivateKeys[0])
		require.NoError(t, err, "invalid private key")
		deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(chainID))
		require.NoError(t, err, "failed to create deployer")
		chainCfg := devenv.ChainConfig{
			ChainID:     id,
			ChainName:   DefaultChainNamePrefix + chain.ChainID,
			ChainType:   devenv.EVMChainType,
			WSRPCs:      []string{chain.Nodes[0].HostWSUrl},
			HTTPRPCs:    []string{chain.Nodes[0].HostHTTPUrl},
			DeployerKey: deployer,
		}
		err = chainCfg.SetUsers(in.Network.PrivateKeys[1:])
		require.NoError(t, err, "Error setting users")
		chains = append(chains, chainCfg)
	}
	cfg, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.CCIP)
	require.NoError(t, err, "Error getting config")

	return &devenv.EnvironmentConfig{
		Chains:   chains,
		JDConfig: jdConfig,
	}, cfg, blockchains, in
}

func startCLNodes(
	t *testing.T,
	crConfig deployment.CapabilityRegistryConfig,
	blockchains []*blockchain.Output,
	in *CTFV2Conf,
) *ns.Output {
	tomlNodeConfig := in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides
	tomlNodeConfig += getChainSpecificNodeSpec(blockchains)
	tomlNodeConfig += fmt.Sprintf(`
		# This is needed for external registry
		[Capabilities]
		[Capabilities.ExternalRegistry]
		Address = '%s'
		NetworkID = 'evm'
		ChainID = '%s'`,
		crConfig.Contract.String(),
		strconv.FormatUint(crConfig.EVMChainID, 10))
	in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = tomlNodeConfig

	nodeOut, err := ns.NewSharedDBNodeSet(in.NodeSet, blockchains[0])
	require.NoError(t, err, "failed to create node set")
	require.NotNil(t, nodeOut)
	return nodeOut
}

func getChainSpecificNodeSpec(bcs []*blockchain.Output) string {
	//TODO: Can improve this based on number of chains instead of hardcoding for two chains
	var tomlNodeConfig string
	for i, bc := range bcs {
		if i == 0 {
			tomlNodeConfig += fmt.Sprintf(`
				[[EVM]]
				ChainID = '%s'
				AutoCreateKey = true
				FinalityDepth = 1
				MinContractPayment = '0'
				
				[EVM.GasEstimator]
				PriceMax = '200 gwei'
				LimitDefault = 6000000
				FeeCapDefault = '200 gwei'
				
				[[EVM.Nodes]]
				Name = '%s'
				WSURL = '%s'
				HTTPURL = '%s'`,
				bc.ChainID,
				"chain-"+bc.ChainID,
				bc.Nodes[0].DockerInternalWSUrl,
				bc.Nodes[0].DockerInternalHTTPUrl)
		} else {
			tomlNodeConfig += fmt.Sprintf(`
				[[EVM]]
				ChainID = '%s'
				LogPollInterval = '500ms'
		
				[EVM.Transactions]
				ForwardersEnabled = false
				
				[EVM.GasEstimator]
				LimitDefault = 5000000
		
				[[EVM.Nodes]]
				Name = '%s'
				WSURL = '%s'
				HTTPURL = '%s'`,
				bc.ChainID,
				"chain-"+bc.ChainID,
				bc.Nodes[0].DockerInternalWSUrl,
				bc.Nodes[0].DockerInternalHTTPUrl)
		}
	}
	return tomlNodeConfig
}
