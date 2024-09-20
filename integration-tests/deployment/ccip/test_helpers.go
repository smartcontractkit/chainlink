package ccipdeployment

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"

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

type DeployedTestEnvironment struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	FeedChainSel uint64
	Nodes        map[string]memory.Node
}

// NewEnvironmentWithCRAndFeeds creates a new CCIP environment
// with capreg, feeds and nodes set up.
func NewEnvironmentWithCRAndFeeds(t *testing.T, lggr logger.Logger, numChains int) DeployedTestEnvironment {
	require.GreaterOrEqual(t, numChains, 2, "numChains must be at least 2 for home and feed chains")
	ctx := Context(t)
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
	homeChainEVM, _ := chainsel.ChainIdFromSelector(homeChainSel)
	ab, capReg, err := DeployCapReg(lggr, chains[homeChainSel])
	require.NoError(t, err)

	feedSel := chainSels[FeedChainIndex]
	feedAb, _, err := DeployFeeds(lggr, chains[feedSel])
	require.NoError(t, err)
	require.NoError(t, ab.Merge(feedAb))

	nodes := memory.NewNodes(t, zapcore.InfoLevel, chains, 4, 1, deployment.CapabilityRegistryConfig{
		EVMChainID: homeChainEVM,
		Contract:   capReg,
	})
	for _, node := range nodes {
		require.NoError(t, node.App.Start(ctx))
		t.Cleanup(func() {
			require.NoError(t, node.App.Stop())
		})
	}

	e := memory.NewMemoryEnvironmentFromChainsNodes(t, lggr, chains, nodes)
	return DeployedTestEnvironment{
		Ab:           ab,
		Env:          e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		Nodes:        nodes,
	}
}

func NewEnvironmentWithCRAndJobs(t *testing.T, lggr logger.Logger, numChains int) DeployedTestEnvironment {
	ctx := Context(t)
	e := NewEnvironmentWithCRAndFeeds(t, lggr, numChains)
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

	// Ensure job related logs are up to date.
	require.NoError(t, ReplayAllLogs(e.Nodes, e.Env.Chains))
	return e
}

func ReplayAllLogs(nodes map[string]memory.Node, chains map[uint64]deployment.Chain) error {
	for _, node := range nodes {
		for sel := range chains {
			if err := node.ReplayLogs(map[uint64]uint64{sel: 1}); err != nil {
				return err
			}
		}
	}
	return nil
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
	}, []uint64{dest})
	require.NoError(t, err)
	require.True(t, it.Next())
	seqNum := it.Event.Message.Header.SequenceNumber
	t.Logf("CCIP message sent from chain selector %d to chain selector %d tx %s seqNum %d", src, dest, tx.Hash().String(), seqNum)
	return seqNum
}

// DeployedLocalDevEnvironment is a helper struct for setting up a local dev environment with docker
type DeployedLocalDevEnvironment struct {
	Ab           deployment.AddressBook
	Env          deployment.Environment
	HomeChainSel uint64
	Nodes        []devenv.Node
}

func NewDeployedLocalDevEnvironment(t *testing.T, lggr logger.Logger) DeployedLocalDevEnvironment {
	ctx := Context(t)
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
	homeChainEVM, err := chainsel.ChainIdFromSelector(homeChainSel)
	require.NoError(t, err)
	require.NotEmpty(t, homeChainEVM, "homeChainEVM should not be empty")

	// deploy the capability registry
	ab, capReg, err := DeployCapReg(lggr, chains[homeChainSel])
	require.NoError(t, err)

	// start the chainlink nodes with the CR address
	err = devenv.StartChainlinkNodes(t,
		envConfig, deployment.CapabilityRegistryConfig{
			EVMChainID: homeChainEVM,
			Contract:   capReg,
		},
		testEnv, cfg)
	require.NoError(t, err)

	e, don, err := devenv.NewEnvironment(ctx, lggr, *envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	require.NotNil(t, don)
	zeroLogLggr := logging.GetTestLogger(t)
	// fund the nodes
	devenv.FundNodes(t, zeroLogLggr, testEnv, cfg, don.PluginNodes())
	return DeployedLocalDevEnvironment{
		Ab:           ab,
		Env:          *e,
		HomeChainSel: homeChainSel,
		Nodes:        don.Nodes,
	}
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

func DeployFeeds(lggr logger.Logger, chain deployment.Chain) (deployment.AddressBook, map[string]common.Address, error) {
	ab := deployment.NewMemoryAddressBook()
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
		return ab, nil, err
	}

	lggr.Infow("deployed mockLinkFeed", "addr", mockLinkFeed.Address)

	desc, err := mockLinkFeed.Contract.Description(&bind.CallOpts{})
	if err != nil {
		lggr.Errorw("Failed to get description", "err", err)
		return ab, nil, err
	}

	if desc != MockLinkAggregatorDescription {
		lggr.Errorw("Unexpected description for Link token", "desc", desc)
		return ab, nil, fmt.Errorf("unexpected description: %s", desc)
	}

	tvToAddress := map[string]common.Address{
		desc: mockLinkFeed.Address,
	}
	return ab, tvToAddress, nil
}
