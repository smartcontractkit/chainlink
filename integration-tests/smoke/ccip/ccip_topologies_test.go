package ccip

import (
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/cciptesthelpertypes"
	mt "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIPTopologies_EVM2EVM_RoleDON_AllSupportSource_SomeSupportDest(t *testing.T) {
	runCCIPTopologiesTest(t, 2, 1)
}

func Test_CCIPTopologies_EVM2EVM_RoleDON_AllSupportDest_SomeSupportSource(t *testing.T) {
	runCCIPTopologiesTest(t, 1, 2)
}

func runCCIPTopologiesTest(t *testing.T, fChainSource, fChainDest int) {
	// fix the chain ids for the test so we can appropriately set finality depth numbers on the destination chain.
	chains := []chainsel.Chain{
		chainsel.TEST_90000001,
		chainsel.TEST_90000002,
		chainsel.TEST_90000003,
	}
	sort.Slice(chains, func(i, j int) bool { return chains[i].Selector < chains[j].Selector })
	homeChainSel := chains[0].Selector
	sourceChainSel := chains[1].Selector
	destChainSel := chains[2].Selector

	const (
		fRoleDON = 2
		nRoleDON = 3*fRoleDON + 1
	)

	// Setup 3 chains and a single lane.
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(len(chains)),
		testhelpers.WithNumOfNodes(nRoleDON),
		testhelpers.WithRoleDONTopology(cciptesthelpertypes.NewRandomTopology(
			cciptesthelpertypes.RandomTopologyArgs{
				FChainToNumChains: map[int]int{
					fChainSource: 1, // 1 chain with fChain fChainSource
					fChainDest:   1, // 1 chain with fChain fChainDest
				},
				Seed: 42, // for reproducible setups.
			},
		)),
	)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	allChainSelectors := maps.Keys(e.Env.BlockChains.EVMChains())
	require.Len(t, allChainSelectors, 3)

	require.Contains(t, allChainSelectors, homeChainSel)
	require.Contains(t, allChainSelectors, sourceChainSel)
	require.Contains(t, allChainSelectors, destChainSel)

	t.Log("All chain selectors:", allChainSelectors,
		", home chain selector:", e.HomeChainSel,
		", feed chain selector:", e.FeedChainSel,
		", source chain selector:", sourceChainSel,
		", dest chain selector:", destChainSel,
	)
	// connect a single lane, source to dest
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(
		t, &e, state, sourceChainSel, destChainSel, false)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChainSel].DeployerKey.From.Bytes(), 32)
		setup  = mt.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChainSel,
			destChainSel,
			sender,
			false, // testRouter
		)
	)

	// Wait for filter registration for CCIPMessageSent (onramp), CommitReportAccepted (offramp), and ExecutionStateChanged (offramp)
	testhelpers.WaitForEventFilterRegistrationOnLane(t, state, e.Env.Offchain, sourceChainSel, destChainSel)

	t.Run("data message to eoa", func(t *testing.T) {
		_ = mt.Run(
			t,
			mt.TestCase{
				ValidationType:         mt.ValidationTypeExec,
				TestSetup:              setup,
				Nonce:                  &nonce,
				Receiver:               common.HexToAddress("0xdead").Bytes(),
				MsgData:                []byte("hello eoa"),
				ExtraArgs:              nil,                                 // default extraArgs
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS, // success because offRamp won't call an EOA
				ExtraAssertions: []func(t *testing.T){
					func(t *testing.T) {
					},
				},
			},
		)
	})
}
