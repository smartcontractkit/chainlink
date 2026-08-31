package ccip

import (
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

// mixedVersionImageSelector assigns the baseline image to bootstrap and even-indexed
// plugin nodes, and the PR image to odd-indexed plugin nodes.
func mixedVersionImageSelector(t *testing.T) testsetups.NodeImageSelector {
	prImage := os.Getenv("CCIP_PR_IMAGE")
	prVersion := os.Getenv("CCIP_PR_VERSION")
	baselineImage := os.Getenv("CCIP_BASELINE_IMAGE")
	baselineVersion := os.Getenv("CCIP_BASELINE_VERSION")

	require.NotEmpty(t, prImage, "CCIP_PR_IMAGE must be set")
	require.NotEmpty(t, prVersion, "CCIP_PR_VERSION must be set")
	require.NotEmpty(t, baselineImage, "CCIP_BASELINE_IMAGE must be set")
	require.NotEmpty(t, baselineVersion, "CCIP_BASELINE_VERSION must be set")

	return func(nodeIndex int) (string, string) {
		// nodeIndex 0 is the bootstrap node, the rest are plugin nodes.
		if nodeIndex == 0 || nodeIndex%2 == 1 {
			return baselineImage, baselineVersion
		}
		return prImage, prVersion
	}
}

func Test_CCIPMixedVersionDON(t *testing.T) {
	if os.Getenv(testhelpers.ENVTESTTYPE) != string(testhelpers.Docker) {
		t.Skip("skipping mixed version test in non-docker mode, set CCIP_V16_TEST_ENV=docker")
	}

	chains := []chainsel.Chain{
		chainsel.GETH_TESTNET,  // source
		chainsel.GETH_DEVNET_2, // dest
	}
	var chainIDs = []uint64{
		chains[0].Selector,
		chains[1].Selector,
	}

	e, _, _ := testsetups.NewIntegrationEnvironmentWithImageSelector(
		t,
		mixedVersionImageSelector(t),
		testhelpers.WithEVMChainsBySelectors(chainIDs),
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	require.Len(t, allChainSelectors, 2)
	sourceChain := chains[0].Selector
	destChain := chains[1].Selector
	require.Contains(t, allChainSelectors, sourceChain)
	require.Contains(t, allChainSelectors, destChain)
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	// connect a single lane, source to dest
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	t.Run("data message to eoa on mixed version DON", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  &nonce,
				Receiver:               common.HexToAddress("0xdead").Bytes(),
				MsgData:                []byte("hello from mixed version DON"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
			},
		)
	})

	t.Run("message to contract implementing CCIPReceiver on mixed version DON", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  &out.Nonce,
				Receiver:               state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(),
				MsgData:                []byte("hello CCIPReceiver from mixed version DON"),
				ExtraArgs:              nil, // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}

func Test_CCIPRollingUpgrade(t *testing.T) {
	if os.Getenv(testhelpers.ENVTESTTYPE) != string(testhelpers.Docker) {
		t.Skip("skipping rolling upgrade test in non-docker mode, set CCIP_V16_TEST_ENV=docker")
	}

	prImage := os.Getenv("CCIP_PR_IMAGE")
	prVersion := os.Getenv("CCIP_PR_VERSION")
	baselineImage := os.Getenv("CCIP_BASELINE_IMAGE")
	baselineVersion := os.Getenv("CCIP_BASELINE_VERSION")

	require.NotEmpty(t, prImage, "CCIP_PR_IMAGE must be set")
	require.NotEmpty(t, prVersion, "CCIP_PR_VERSION must be set")
	require.NotEmpty(t, baselineImage, "CCIP_BASELINE_IMAGE must be set")
	require.NotEmpty(t, baselineVersion, "CCIP_BASELINE_VERSION must be set")

	chains := []chainsel.Chain{
		chainsel.GETH_TESTNET,  // source
		chainsel.GETH_DEVNET_2, // dest
	}
	var chainIDs = []uint64{
		chains[0].Selector,
		chains[1].Selector,
	}

	// Start all nodes on the baseline image.
	allBaselineImage := func(int) (string, string) {
		return baselineImage, baselineVersion
	}

	e, _, dockerEnv := testsetups.NewIntegrationEnvironmentWithImageSelector(
		t,
		allBaselineImage,
		testhelpers.WithEVMChainsBySelectors(chainIDs),
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	require.Len(t, allChainSelectors, 2)
	sourceChain := chains[0].Selector
	destChain := chains[1].Selector
	require.Contains(t, allChainSelectors, sourceChain)
	require.Contains(t, allChainSelectors, destChain)
	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChain,
		", dest chain selector:", destChain,
	)

	// connect a single lane, source to dest
	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		out    mt.TestCaseOutput
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	// Send a baseline message before any upgrades.
	t.Run("baseline message before upgrade", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  &nonce,
				Receiver:               common.HexToAddress("0xdead").Bytes(),
				MsgData:                []byte("baseline before rolling upgrade"),
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	// Get the CL cluster nodes for upgrades.
	localDevEnv, ok := dockerEnv.(*testsetups.DeployedLocalDevEnvironment)
	require.True(t, ok, "expected DeployedLocalDevEnvironment in docker mode")
	clCluster := localDevEnv.GetCLClusterTestEnv().ClCluster
	require.NotNil(t, clCluster, "CL cluster should not be nil")
	require.NotEmpty(t, clCluster.Nodes, "CL cluster should have nodes")

	// Upgrade plugin nodes one by one, sending a message after each upgrade.
	// Skip the bootstrap node (index 0).
	for i, node := range clCluster.Nodes {
		if i == 0 {
			continue
		}
		t.Run(fmt.Sprintf("upgrade node %d (%s) and send message", i, node.ContainerName), func(t *testing.T) {
			t.Logf("Upgrading node %d (%s) from %s:%s to %s:%s",
				i, node.ContainerName,
				node.ContainerImage, node.ContainerVersion,
				prImage, prVersion)

			err := node.UpgradeVersion(prImage, prVersion)
			require.NoError(t, err, "failed to upgrade node %d", i)

			out = mt.Run(
				t,
				mt.TestCase{
					ValidationType:         mt.ValidationTypeExec,
					TestSetup:              setup,
					Nonce:                  &out.Nonce,
					Receiver:               common.HexToAddress("0xdead").Bytes(),
					MsgData:                []byte(fmt.Sprintf("message after upgrading node %d", i)),
					ExtraArgs:              nil,
					ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
				},
			)
		})
	}

	// Send a final message after all upgrades.
	t.Run("final message after all upgrades", func(t *testing.T) {
		out = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  &out.Nonce,
				Receiver:               state.MustGetEVMChainState(destChain).Receiver.Address().Bytes(),
				MsgData:                []byte("final message after rolling upgrade"),
				ExtraArgs:              nil,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}
