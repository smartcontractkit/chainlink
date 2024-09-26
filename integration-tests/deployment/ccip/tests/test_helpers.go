package tests

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/devenv"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
)

const (
	HomeChainIndex = 0
	FeedChainIndex = 1
)

type DeployedEnv struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	FeedChainSel uint64
	TokenConfig  ccipdeployment.TokenConfig
	ReplayBlocks map[uint64]uint64
}

func (e *DeployedEnv) SetupJobs(t *testing.T) {
	ctx := testcontext.Get(t)
	jbs, err := ccipdeployment.NewCCIPJobSpecs(e.Env.NodeIDs, e.Env.Offchain)
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

	// Ensure job related logs are up to date.
	require.NoError(t, e.Env.Offchain.ReplayLogs(e.ReplayBlocks))
}

func SetUpHomeAndFeedChains(t *testing.T, lggr logger.Logger, homeChainSel, feedChainSel uint64, chains map[uint64]deployment.Chain) (deployment.AddressBook, deployment.CapabilityRegistryConfig) {
	homeChainEVM, _ := chainsel.ChainIdFromSelector(homeChainSel)
	ab, capReg, err := ccipdeployment.DeployCapReg(lggr, chains[homeChainSel])
	require.NoError(t, err)

	feedAb, _, err := ccipdeployment.DeployMockFeeds(lggr, chains[feedChainSel])
	require.NoError(t, err)
	require.NoError(t, ab.Merge(feedAb))

	return ab, deployment.CapabilityRegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
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

// NewMemoryEnvironment creates a new CCIP environment
// with capreg, feeds and nodes set up.
func NewMemoryEnvironment(t *testing.T, lggr logger.Logger, numChains int) DeployedEnv {
	require.GreaterOrEqual(t, numChains, 2, "numChains must be at least 2 for home and feed chains")
	ctx := testcontext.Get(t)
	chains := memory.NewMemoryChains(t, numChains)
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
	homeChainSel := chainSels[HomeChainIndex]
	feedSel := chainSels[FeedChainIndex]
	replayBlocks, err := LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)
	ab, capReg := SetUpHomeAndFeedChains(t, lggr, homeChainSel, feedSel, chains)

	tokenConfig := ccipdeployment.NewTokenConfig()

	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, capReg)
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(lggr, chains, nodes)
	return DeployedEnv{
		Ab:           ab,
		Env:          e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		TokenConfig:  tokenConfig,
		ReplayBlocks: replayBlocks,
	}
}

func NewMemoryEnvironmentWithJobs(t *testing.T, lggr logger.Logger, numChains int) DeployedEnv {
	e := NewMemoryEnvironment(t, lggr, numChains)
	e.SetupJobs(t)
	return e
}

func NewLocalDevEnvironment(t *testing.T, lggr logger.Logger) DeployedEnv {
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
	ab, capReg := SetUpHomeAndFeedChains(t, lggr, homeChainSel, feedSel, chains)

	// start the chainlink nodes with the CR address
	err = devenv.StartChainlinkNodes(t, envConfig, capReg, testEnv, cfg)
	require.NoError(t, err)

	e, don, err := devenv.NewEnvironment(ctx, lggr, *envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	zeroLogLggr := logging.GetTestLogger(t)
	// fund the nodes
	devenv.FundNodes(t, zeroLogLggr, testEnv, cfg, don.PluginNodes())

	return DeployedEnv{
		Ab:           ab,
		Env:          *e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
		TokenConfig:  ccipdeployment.NewTokenConfig(),
	}
}

func SendRequest(t *testing.T, e deployment.Environment, state ccipdeployment.CCIPOnChainState, src, dest uint64, testRouter bool) uint64 {
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
	}, []uint64{dest})
	require.NoError(t, err)
	require.True(t, it.Next())
	seqNum := it.Event.Message.Header.SequenceNumber
	t.Logf("CCIP message sent from chain selector %d to chain selector %d tx %s seqNum %d", src, dest, tx.Hash().String(), seqNum)
	return seqNum
}

// AddLanesForAll adds densely connected lanes for all chains in the environment so that each chain
// is connected to every other chain except itself.
func AddLanesForAll(e deployment.Environment, state ccipdeployment.CCIPOnChainState) error {
	for source := range e.Chains {
		for dest := range e.Chains {
			if source != dest {
				err := ccipdeployment.AddLane(e, state, source, dest)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
