package smoke

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	ccdeploy "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view"
	jobv1 "github.com/smartcontractkit/chainlink/integration-tests/deployment/jd/job/v1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRevert(t *testing.T) {
	ec, err := ethclient.Dial("wss...")
	require.NoError(t, err)
	tx, _, err := ec.TransactionByHash(testcontext.Get(t), common.HexToHash("0x6aae71ad356d383a41d20aad69de4e1cde536c33b0a41488e529df731bf086ab"))
	require.NoError(t, err)
	rec, err := ec.TransactionReceipt(testcontext.Get(t), tx.Hash())
	require.NoError(t, err)
	fromTx, err := deployment.GetErrorReasonFromTx(ec, common.HexToAddress("0xBE7294B7910606845500Ba524FfcAC8917A00F34"), *tx, rec)
	require.NoError(t, err)
	errStr, err := deployment.ParseErrorFromABI(fromTx, capabilities_registry.CapabilitiesRegistryABI)
	require.NoError(t, err)
	fmt.Println(errStr)
}

func Test0002_InitialDeployOnLocal(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testcontext.Get(t)
	tenv := ccdeploy.NewLocalDevEnvironment(t, lggr)

	e := tenv.Env

	state, err := ccdeploy.LoadOnchainState(tenv.Env, tenv.Ab)
	require.NoError(t, err)

	feeds := state.Chains[tenv.FeedChainSel].USDFeeds
	tokenConfig := ccdeploy.NewTokenConfig()
	tokenConfig.UpsertTokenInfo(ccdeploy.LinkSymbol,
		pluginconfig.TokenInfo{
			AggregatorAddress: feeds[ccdeploy.LinkSymbol].Address().String(),
			Decimals:          ccdeploy.LinkDecimals,
			DeviationPPB:      cciptypes.NewBigIntFromInt64(1e9),
		},
	)
	// Apply migration
	output, err := changeset.Apply0002(tenv.Env, ccdeploy.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    tokenConfig,
		// Capreg/config and feeds already exist.
		CCIPOnChainState: state,
	})
	require.NoError(t, err)
	// Get new state after migration.
	state, err = ccdeploy.LoadOnchainState(e, output.AddressBook)
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	ccdeploy.ReplayLogs(t, e.Offchain, tenv.ReplayBlocks)

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
	require.NoError(t, ccdeploy.AddLanesForAll(e, state))
	v, err := state.View(e.AllChainSelectors())
	require.NoError(t, err)
	require.NoError(t, view.SaveView(v))
	// Need to keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	// Send a message from each chain to every other chain.
	expectedSeqNum := make(map[uint64]uint64)
	for src := range e.Chains {
		for dest, destChain := range e.Chains {
			if src == dest {
				continue
			}
			latesthdr, err := destChain.Client.HeaderByNumber(testcontext.Get(t), nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[dest] = &block
			seqNum := ccdeploy.SendRequest(t, e, state, src, dest, false)
			expectedSeqNum[dest] = seqNum
		}
	}

	// Wait for all commit reports to land.
	ccdeploy.ConfirmCommitForAllWithExpectedSeqNums(t, e, state, expectedSeqNum, startBlocks)

	// After commit is reported on all chains, token prices should be updated in FeeQuoter.
	for dest := range e.Chains {
		linkAddress := state.Chains[dest].LinkToken.Address()
		feeQuoter := state.Chains[dest].FeeQuoter
		timestampedPrice, err := feeQuoter.GetTokenPrice(nil, linkAddress)
		require.NoError(t, err)
		require.Equal(t, ccdeploy.MockLinkPrice, timestampedPrice.Value)
	}

	// Wait for all exec reports to land
	ccdeploy.ConfirmExecWithSeqNrForAll(t, e, state, expectedSeqNum, startBlocks)

	// TODO: Apply the proposal.
}
