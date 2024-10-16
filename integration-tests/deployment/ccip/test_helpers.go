package ccipdeployment

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/devenv"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
)

const (
	HomeChainIndex = 0
	FeedChainIndex = 1
)

// Context returns a context with the test's deadline, if available.
func Context(tb testing.TB) context.Context {
	ctx := context.Background()
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}

type DeployedEnv struct {
	Env               deployment.Environment
	Ab                deployment.AddressBook
	HomeChainSel      uint64
	FeedChainSel      uint64
	ReplayBlocks      map[uint64]uint64
	FeeTokenContracts map[uint64]FeeTokenContracts
}

func (e *DeployedEnv) SetupJobs(t *testing.T) {
	ctx := testcontext.Get(t)
	jbs, err := NewCCIPJobSpecs(e.Env.NodeIDs, e.Env.Offchain)
	require.NoError(t, err)
	for nodeID, jobs := range jbs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Env.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}
	// Wait for plugins to register filters?
	// TODO: Investigate how to avoid.
	time.Sleep(30 * time.Second)
	ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)
}

func ReplayLogs(t *testing.T, oc deployment.OffchainClient, replayBlocks map[uint64]uint64) {
	switch oc := oc.(type) {
	case *memory.JobClient:
		require.NoError(t, oc.ReplayLogs(replayBlocks))
	case *devenv.JobDistributor:
		require.NoError(t, oc.ReplayLogs(replayBlocks))
	default:
		t.Fatalf("unsupported offchain client type %T", oc)
	}
}

func DeployTestContracts(t *testing.T,
	lggr logger.Logger,
	ab deployment.AddressBook,
	homeChainSel,
	feedChainSel uint64,
	chains map[uint64]deployment.Chain,
) (map[uint64]FeeTokenContracts, deployment.CapabilityRegistryConfig) {
	capReg, err := DeployCapReg(lggr, ab, chains[homeChainSel])
	require.NoError(t, err)
	_, err = DeployFeeds(lggr, ab, chains[feedChainSel])
	require.NoError(t, err)
	feeTokenContracts, err := DeployFeeTokensToChains(lggr, ab, chains)
	require.NoError(t, err)
	evmChainID, err := chainsel.ChainIdFromSelector(homeChainSel)
	require.NoError(t, err)
	return feeTokenContracts, deployment.CapabilityRegistryConfig{
		EVMChainID: evmChainID,
		Contract:   capReg.Address,
	}
}

func LatestBlocksByChain(ctx context.Context, chains map[uint64]deployment.Chain) (map[uint64]uint64, error) {
	latestBlocks := make(map[uint64]uint64)
	for _, chain := range chains {
		latesthdr, err := chain.Client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get latest header for chain %d", chain.Selector)
		}
		block := latesthdr.Number.Uint64()
		latestBlocks[chain.Selector] = block
	}
	return latestBlocks, nil
}

func allocateCCIPChainSelectors(chains map[uint64]deployment.Chain) (homeChainSel uint64, feeChainSel uint64) {
	// Lower chainSel is home chain.
	var chainSels []uint64
	// Say first chain is home chain.
	for chainSel := range chains {
		chainSels = append(chainSels, chainSel)
	}
	sort.Slice(chainSels, func(i, j int) bool {
		return chainSels[i] < chainSels[j]
	})
	// Take lowest for determinism.
	return chainSels[HomeChainIndex], chainSels[FeedChainIndex]
}

// NewMemoryEnvironment creates a new CCIP environment
// with capreg, fee tokens, feeds and nodes set up.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, numChains int) DeployedEnv {
	require.GreaterOrEqual(t, numChains, 2, "numChains must be at least 2 for home and feed chains")
	ctx := testcontext.Get(t)
	chains := memory.NewMemoryChains(t, numChains)
	homeChainSel, feedSel := allocateCCIPChainSelectors(chains)
	replayBlocks, err := LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)

	ab := deployment.NewMemoryAddressBook()
	feeTokenContracts, crConfig := DeployTestContracts(t, lggr, ab, homeChainSel, feedSel, chains)
	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, crConfig)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	return DeployedEnv{
		Ab:                ab,
		Env:               e,
		HomeChainSel:      homeChainSel,
		FeedChainSel:      feedSel,
		ReplayBlocks:      replayBlocks,
		FeeTokenContracts: feeTokenContracts,
	}
}

func NewMemoryEnvironmentWithJobs(t *testing.T, lggr logger.Logger, numChains int) DeployedEnv {
	e := NewMemoryEnvironment(t, lggr, numChains)
	e.SetupJobs(t)
	return e
}

func SendRequest(t *testing.T, e deployment.Environment, state CCIPOnChainState, src, dest uint64, testRouter bool) uint64 {
	msg := router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
		Data:         []byte("hello"),
		TokenAmounts: nil, // TODO: no tokens for now
		// Pay native.
		FeeToken:  common.HexToAddress("0x0"),
		ExtraArgs: nil, // TODO: no extra args for now, falls back to default
	}
	router := state.Chains[src].Router
	if testRouter {
		router = state.Chains[src].TestRouter
	}
	fee, err := router.GetFee(
		&bind.CallOpts{Context: context.Background()}, dest, msg)
	require.NoError(t, err, deployment.MaybeDataErr(err))

	t.Logf("Sending CCIP request from chain selector %d to chain selector %d",
		src, dest)
	e.Chains[src].DeployerKey.Value = fee
	tx, err := router.CcipSend(
		e.Chains[src].DeployerKey,
		dest,
		msg)
	require.NoError(t, err)
	blockNum, err := e.Chains[src].Confirm(tx)
	require.NoError(t, err)
	it, err := state.Chains[src].OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	}, []uint64{dest}, []uint64{})
	require.NoError(t, err)
	require.True(t, it.Next())
	seqNum := it.Event.Message.Header.SequenceNumber
	t.Logf("CCIP message sent from chain selector %d to chain selector %d tx %s seqNum %d", src, dest, tx.Hash().String(), seqNum)
	return seqNum
}

// DeployedLocalDevEnvironment is a helper struct for setting up a local dev environment with docker
type DeployedLocalDevEnvironment struct {
	DeployedEnv
	testEnv *test_env.CLClusterTestEnv
	DON     *devenv.DON
}

func (d DeployedLocalDevEnvironment) RestartChainlinkNodes(t *testing.T) error {
	errGrp := errgroup.Group{}
	for _, n := range d.testEnv.ClCluster.Nodes {
		n := n
		errGrp.Go(func() error {
			if err := n.Container.Terminate(testcontext.Get(t)); err != nil {
				return err
			}
			err := n.RestartContainer()
			if err != nil {
				return err
			}
			return nil
		})

	}
	return errGrp.Wait()
}

func NewLocalDevEnvironment(t *testing.T, lggr logger.Logger) (DeployedEnv, *test_env.CLClusterTestEnv, testconfig.TestConfig) {
	ctx := testcontext.Get(t)
	// create a local docker environment with simulated chains and job-distributor
	// we cannot create the chainlink nodes yet as we need to deploy the capability registry first
	envConfig, testEnv, cfg := devenv.CreateDockerEnv(t)
	require.NotNil(t, envConfig)
	require.NotEmpty(t, envConfig.Chains, "chainConfigs should not be empty")
	require.NotEmpty(t, envConfig.JDConfig, "jdUrl should not be empty")
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	require.NoError(t, err)
	// locate the home chain
	homeChainSel := envConfig.HomeChainSelector
	require.NotEmpty(t, homeChainSel, "homeChainSel should not be empty")
	feedSel := envConfig.FeedChainSelector
	require.NotEmpty(t, feedSel, "feedSel should not be empty")
	replayBlocks, err := LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)

	ab := deployment.NewMemoryAddressBook()
	feeContracts, crConfig := DeployTestContracts(t, lggr, ab, homeChainSel, feedSel, chains)

	// start the chainlink nodes with the CR address
	err = devenv.StartChainlinkNodes(t, envConfig,
		crConfig,
		testEnv, cfg)
	require.NoError(t, err)

	e, don, err := devenv.NewEnvironment(ctx, lggr, *envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	zeroLogLggr := logging.GetTestLogger(t)
	// fund the nodes
	devenv.FundNodes(t, zeroLogLggr, testEnv, cfg, don.PluginNodes())

	return DeployedEnv{
		Ab:                ab,
		Env:               *e,
		HomeChainSel:      homeChainSel,
		FeedChainSel:      feedSel,
		ReplayBlocks:      replayBlocks,
		FeeTokenContracts: feeContracts,
	}, testEnv, cfg
}

func NewLocalDevEnvironmentWithRMN(t *testing.T, lggr logger.Logger) (DeployedEnv, devenv.RMNCluster) {
	tenv, dockerenv, _ := NewLocalDevEnvironment(t, lggr)
	state, err := LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := NewTokenConfig()
	tokenConfig.UpsertTokenInfo(LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[LinkSymbol].Address().String(),
			Decimals:          LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	// Deploy CCIP contracts.
	err = DeployCCIPContracts(tenv.Env, tenv.Ab, DeployCCIPContractConfig{
		HomeChainSel:       tenv.HomeChainSel,
		FeedChainSel:       tenv.FeedChainSel,
		ChainsToDeploy:     tenv.Env.AllChainSelectors(),
		TokenConfig:        tokenConfig,
		MCMSConfig:         NewTestMCMSConfig(t, tenv.Env),
		CapabilityRegistry: state.Chains[tenv.HomeChainSel].CapabilityRegistry.Address(),
		FeeTokenContracts:  tenv.FeeTokenContracts,
	})
	require.NoError(t, err)
	l := logging.GetTestLogger(t)
	config := GenerateTestRMNConfig(t, 1, tenv, MustNetworksToRPCMap(dockerenv.EVMNetworks))
	rmnCluster, err := devenv.NewRMNCluster(
		t, l,
		[]string{dockerenv.DockerNetwork.Name},
		config,
		"rageproxy",
		"latest",
		"afn2proxy",
		"latest",
		dockerenv.LogStream,
	)
	require.NoError(t, err)
	return tenv, *rmnCluster
}

func MustNetworksToRPCMap(evmNetworks []*blockchain.EVMNetwork) map[uint64]string {
	rpcs := make(map[uint64]string)
	for _, network := range evmNetworks {
		sel, err := chainsel.SelectorFromChainId(uint64(network.ChainID))
		if err != nil {
			panic(err)
		}
		rpcs[sel] = network.HTTPURLs[0]
	}
	return rpcs
}

func MustCCIPNameToRMNName(a string) string {
	m := map[string]string{
		chainsel.GETH_TESTNET.Name:  "DevnetAlpha",
		chainsel.GETH_DEVNET_2.Name: "DevnetBeta",
		// TODO: Add more as needed.
	}
	v, ok := m[a]
	if !ok {
		panic(fmt.Sprintf("no mapping for %s", a))
	}
	return v
}

func GenerateTestRMNConfig(t *testing.T, nRMNNodes int, tenv DeployedEnv, rpcMap map[uint64]string) map[string]devenv.RMNConfig {
	// Find the bootstrappers.
	nodes, err := deployment.NodeInfo(tenv.Env.NodeIDs, tenv.Env.Offchain)
	require.NoError(t, err)
	bootstrappers := nodes.BootstrapLocators()

	// Just set all RMN nodes to support all chains.
	state, err := LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)
	var remoteChains []devenv.RemoteChain
	var rpcs []devenv.Chain
	for chainSel, chain := range state.Chains {
		c, _ := chainsel.ChainBySelector(chainSel)
		rmnName := MustCCIPNameToRMNName(c.Name)
		remoteChains = append(remoteChains, devenv.RemoteChain{
			Name:             rmnName,
			Stability:        devenv.Stability{Type: "stable"},
			StartBlockNumber: 0,
			OffRamp:          chain.OffRamp.Address().String(),
			RMNRemote:        chain.RMNRemote.Address().String(),
		})
		rpcs = append(rpcs, devenv.Chain{
			Name: rmnName,
			RPC:  rpcMap[chainSel],
		})
	}
	hc, _ := chainsel.ChainBySelector(tenv.HomeChainSel)
	shared := devenv.SharedConfig{
		Networking: devenv.Networking{
			RageProxy:     devenv.DefaultRageProxy,
			Bootstrappers: bootstrappers,
		},
		HomeChain: devenv.HomeChain{
			Name:                 MustCCIPNameToRMNName(hc.Name),
			CapabilitiesRegistry: state.Chains[tenv.HomeChainSel].CapabilityRegistry.Address().String(),
			CCIPHome:             state.Chains[tenv.HomeChainSel].CCIPHome.Address().String(),
			// TODO: RMNHome
		},
		RemoteChains: remoteChains,
	}

	rmnConfig := make(map[string]devenv.RMNConfig)
	for i := 0; i < nRMNNodes; i++ {
		// Listen addresses _should_ be able to operator on the same port since
		// they are inside the docker network.
		proxyLocal := devenv.ProxyLocalConfig{
			ListenAddresses:   []string{devenv.DefaultProxyListenAddress},
			AnnounceAddresses: []string{},
			ProxyAddress:      devenv.DefaultRageProxy,
			DiscovererDbPath:  devenv.DefaultDiscovererDbPath,
		}
		rmnConfig[fmt.Sprintf("rmn_%d", i)] = devenv.RMNConfig{
			Shared:      shared,
			Local:       devenv.LocalConfig{Chains: rpcs},
			ProxyShared: devenv.DefaultRageProxySharedConfig,
			ProxyLocal:  proxyLocal,
		}
	}
	return rmnConfig
}

// AddLanesForAll adds densely connected lanes for all chains in the environment so that each chain
// is connected to every other chain except itself.
func AddLanesForAll(e deployment.Environment, state CCIPOnChainState) error {
	for source := range e.Chains {
		for dest := range e.Chains {
			if source != dest {
				err := AddLane(e, state, source, dest)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

const (
	// MockLinkAggregatorDescription This is the description of the MockV3Aggregator.sol contract
	// nolint:lll
	// https://github.com/smartcontractkit/chainlink/blob/a348b98e90527520049c580000a86fb8ceff7fa7/contracts/src/v0.8/tests/MockV3Aggregator.sol#L76-L76
	MockLinkAggregatorDescription = "v0.8/tests/MockV3Aggregator.sol"
	// MockWETHAggregatorDescription WETH use description from MockETHUSDAggregator.sol
	// nolint:lll
	// https://github.com/smartcontractkit/chainlink/blob/a348b98e90527520049c580000a86fb8ceff7fa7/contracts/src/v0.8/automation/testhelpers/MockETHUSDAggregator.sol#L19-L19
	MockWETHAggregatorDescription = "MockETHUSDAggregator"
)

var (
	MockLinkPrice = big.NewInt(5e18)
	// MockDescriptionToTokenSymbol maps a mock feed description to token descriptor
	MockDescriptionToTokenSymbol = map[string]TokenSymbol{
		MockLinkAggregatorDescription: LinkSymbol,
		MockWETHAggregatorDescription: WethSymbol,
	}
)

func DeployFeeds(lggr logger.Logger, ab deployment.AddressBook, chain deployment.Chain) (map[string]common.Address, error) {
	linkTV := deployment.NewTypeAndVersion(PriceFeed, deployment.Version1_0_0)
	mockLinkFeed, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface] {
			linkFeed, tx, _, err1 := mock_v3_aggregator_contract.DeployMockV3Aggregator(
				chain.DeployerKey,
				chain.Client,
				LinkDecimals,  // decimals
				MockLinkPrice, // initialAnswer
			)
			aggregatorCr, err2 := aggregator_v3_interface.NewAggregatorV3Interface(linkFeed, chain.Client)

			return ContractDeploy[*aggregator_v3_interface.AggregatorV3Interface]{
				Address: linkFeed, Contract: aggregatorCr, Tv: linkTV, Tx: tx, Err: multierr.Append(err1, err2),
			}
		})

	if err != nil {
		lggr.Errorw("Failed to deploy link feed", "err", err)
		return nil, err
	}

	lggr.Infow("deployed mockLinkFeed", "addr", mockLinkFeed.Address)

	desc, err := mockLinkFeed.Contract.Description(&bind.CallOpts{})
	if err != nil {
		lggr.Errorw("Failed to get description", "err", err)
		return nil, err
	}

	if desc != MockLinkAggregatorDescription {
		lggr.Errorw("Unexpected description for Link token", "desc", desc)
		return nil, fmt.Errorf("unexpected description: %s", desc)
	}

	tvToAddress := map[string]common.Address{
		desc: mockLinkFeed.Address,
	}
	return tvToAddress, nil
}
