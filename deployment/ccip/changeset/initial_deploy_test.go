package changeset

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
)

func TestInitialDeploy(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := Context(t)
	tenv := NewMemoryEnvironment(t, lggr, 3, 4, MockLinkPrice, MockWethPrice)
	e := tenv.Env

	state, err := LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	output, err := DeployPrerequisites(e, DeployPrerequisiteConfig{
		ChainSelectors: tenv.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))

	cfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, chain := range e.AllChainSelectors() {
		cfg[chain] = commontypes.MCMSWithTimelockConfig{
			Canceller:         commonchangeset.SingleGroupMCMS(t),
			Bypasser:          commonchangeset.SingleGroupMCMS(t),
			Proposer:          commonchangeset.SingleGroupMCMS(t),
			TimelockExecutors: e.AllDeployerKeys(),
			TimelockMinDelay:  big.NewInt(0),
		}
	}
	output, err = commonchangeset.DeployMCMSWithTimelock(e, cfg)
	require.NoError(t, err)
	require.NoError(t, e.ExistingAddresses.Merge(output.AddressBook))

	output, err = InitialDeploy(tenv.Env, DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	// Get new state after migration.
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))
	state, err = LoadOnchainState(e)
	require.NoError(t, err)
	require.NotNil(t, state.Chains[tenv.HomeChainSel].LinkToken)
	// Ensure capreg logs are up to date.
	ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

	// Apply the jobs.
	for nodeID, jobs := range output.JobSpecs {
		for _, job := range jobs {
			// Note these auto-accept
			_, err := e.Offchain.ProposeJob(ctx,
				&jobv1.ProposeJobRequest{
					NodeId: nodeID,
					Spec:   job,
				})
			require.NoError(t, err)
		}
	}

	// Add all lanes
	require.NoError(t, AddLanesForAll(e, state))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[SourceDestPair]uint64)

	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			msgSentEvent := TestSendRequest(t, e, state, src, dest, false, router.ClientEVM2AnyMessage{
				Receiver:     common.LeftPadBytes(state.Chains[dest].Receiver.Address().Bytes(), 32),
				Data:         []byte("hello"),
				TokenAmounts: nil,
				FeeToken:     common.HexToAddress("0x0"),
				ExtraArgs:    nil,
			})
			expectedSeqNum[SourceDestPair{
				SourceChainSelector: src,
				DestChainSelector:   dest,
			}] = msgSentEvent.SequenceNumber
		}
	}

	// Wait for all commit reports to land.
	ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// Confirm token and gas prices are updated
	ConfirmTokenPriceUpdatedForAll(t, e, state, startBlocks,
		DefaultInitialPrices.LinkPrice, DefaultInitialPrices.WethPrice)
	// TODO: Fix gas prices?
	//ConfirmGasPriceUpdatedForAll(t, e, state, startBlocks)
	//
	//// Wait for all exec reports to land
	ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)
}
