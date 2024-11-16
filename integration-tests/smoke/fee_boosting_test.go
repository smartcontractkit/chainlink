package smoke

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"golang.org/x/exp/maps"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type feeboostTestCase struct {
	t                      *testing.T
	sender                 []byte
	deployedEnv            ccdeploy.DeployedEnv
	onchainState           ccdeploy.CCIPOnChainState
	initialPrices          ccdeploy.InitialPrices
	priceFeedPrices        priceFeedPrices
	sourceChain, destChain uint64
}

type priceFeedPrices struct {
	linkPrice *big.Int
	wethPrice *big.Int
}

// TODO: find a way to reuse the same test setup for all tests
func Test_CCIPFeeBoosting(t *testing.T) {
	ctx := ccdeploy.Context(t)

	setupTestEnv := func(t *testing.T, numChains int) (ccdeploy.DeployedEnv, ccdeploy.CCIPOnChainState, []uint64) {
		e, _, _ := testsetups.NewLocalDevEnvironment(
			t, logger.TestLogger(t),
			deployment.E18Mult(5),
			big.NewInt(9e8))

		state, err := ccdeploy.LoadOnchainState(e.Env)
		require.NoError(t, err)

		allChainSelectors := maps.Keys(e.Env.Chains)
		require.Len(t, allChainSelectors, numChains)

		output, err := changeset.DeployPrerequisites(e.Env, changeset.DeployPrerequisiteConfig{
			ChainSelectors: e.Env.AllChainSelectors(),
		})
		require.NoError(t, err)
		require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))

		tokenConfig := ccdeploy.NewTestTokenConfig(state.Chains[e.FeedChainSel].USDFeeds)
		// Apply migration
		output, err = changeset.InitialDeploy(e.Env, ccdeploy.DeployCCIPContractConfig{
			HomeChainSel:   e.HomeChainSel,
			FeedChainSel:   e.FeedChainSel,
			ChainsToDeploy: allChainSelectors,
			TokenConfig:    tokenConfig,
			MCMSConfig:     ccdeploy.NewTestMCMSConfig(t, e.Env),
			OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
		})
		require.NoError(t, err)
		require.NoError(t, e.Env.ExistingAddresses.Merge(output.AddressBook))
		state, err = ccdeploy.LoadOnchainState(e.Env)
		require.NoError(t, err)

		// Ensure capreg logs are up to date.
		ccdeploy.ReplayLogs(t, e.Env.Offchain, e.ReplayBlocks)

		// Apply the jobs.
		for nodeID, jobs := range output.JobSpecs {
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

		return e, state, allChainSelectors
	}

	t.Run("boost needed due to WETH price increase (also covering gas price inscrease)", func(t *testing.T) {
		e, state, chains := setupTestEnv(t, 2)
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			initialPrices: ccdeploy.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  ccdeploy.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: deployment.E18Mult(5),
				wethPrice: big.NewInt(9.9e8), // increase from 9e8 to 9.9e8
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})

	t.Run("boost needed due to LINK price decrease", func(t *testing.T) {
		e, state, chains := setupTestEnv(t, 2)
		runFeeboostTestCase(feeboostTestCase{
			t:            t,
			sender:       common.LeftPadBytes(e.Env.Chains[chains[0]].DeployerKey.From.Bytes(), 32),
			deployedEnv:  e,
			onchainState: state,
			initialPrices: ccdeploy.InitialPrices{
				LinkPrice: deployment.E18Mult(5),
				WethPrice: deployment.E18Mult(9),
				GasPrice:  ccdeploy.ToPackedFee(big.NewInt(1.8e11), big.NewInt(0)),
			},
			priceFeedPrices: priceFeedPrices{
				linkPrice: big.NewInt(4.5e18), // decrease from 5e18 to 4.5e18
				wethPrice: big.NewInt(9e8),
			},
			sourceChain: chains[0],
			destChain:   chains[1],
		})
	})
}

func runFeeboostTestCase(tc feeboostTestCase) {
	require.NoError(tc.t, ccdeploy.AddLane(tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, tc.initialPrices))

	startBlocks := make(map[uint64]*uint64)
	expectedSeqNum := make(map[uint64]uint64)
	msgSentEvent := ccdeploy.TestSendRequest(tc.t, tc.deployedEnv.Env, tc.onchainState, tc.sourceChain, tc.destChain, false, router.ClientEVM2AnyMessage{
		Receiver:     common.LeftPadBytes(tc.onchainState.Chains[tc.destChain].Receiver.Address().Bytes(), 32),
		Data:         []byte("message that needs fee boosting"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	})
	expectedSeqNum[tc.destChain] = msgSentEvent.SequenceNumber

	// hack
	time.Sleep(30 * time.Second)
	replayBlocks := make(map[uint64]uint64)
	replayBlocks[tc.sourceChain] = 1
	replayBlocks[tc.destChain] = 1
	ccdeploy.ReplayLogs(tc.t, tc.deployedEnv.Env.Offchain, replayBlocks)

	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
	ccdeploy.ConfirmExecWithSeqNrForAll(tc.t, tc.deployedEnv.Env, tc.onchainState, expectedSeqNum, startBlocks)
}
