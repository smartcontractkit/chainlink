package cache_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func Test_RootsEligibleForExecution(t *testing.T) {
	ctx := testutils.Context(t)
	chainID := testutils.NewRandomEVMChainID()
	orm := logpoller.NewORM(chainID, pgtest.NewSqlxDB(t), logger.TestLogger(t))
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Hour,
		FinalityDepth:            2,
		BackfillBatchSize:        20,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(orm, nil, logger.TestLogger(t), nil, lpOpts)

	commitStoreAddr := utils.RandomAddress()

	block2 := time.Now().Add(-8 * time.Hour)
	block3 := time.Now().Add(-5 * time.Hour)
	block4 := time.Now().Add(-1 * time.Hour)
	newBlock4 := time.Now().Add(-2 * time.Hour)
	block5 := time.Now()

	root1 := utils.RandomBytes32()
	root2 := utils.RandomBytes32()
	root3 := utils.RandomBytes32()
	root4 := utils.RandomBytes32()
	root5 := utils.RandomBytes32()

	inputLogs := []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 2, 1, root1, block2),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 2, 2, root2, block2),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 2, time.Now(), 1)))

	feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

	commitStore, err := v1_2_0.NewCommitStore(logger.TestLogger(t), commitStoreAddr, nil, lp, feeEstimatorConfig)
	require.NoError(t, err)

	rootsCache := cache.NewCommitRootsCache(logger.TestLogger(t), commitStore, 10*time.Hour, time.Second)

	roots, err := rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2)

	rootsCache.Snooze(root1)
	rootsCache.Snooze(root2)

	// Roots are snoozed
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots)

	// Roots are unsnoozed
	require.Eventually(t, func() bool {
		roots, err = rootsCache.RootsEligibleForExecution(ctx)
		require.NoError(t, err)
		return len(roots) == 2
	}, 5*time.Second, 1*time.Second)

	// Marking root as executed doesn't ignore other roots from the same block
	rootsCache.MarkAsExecuted(root1)
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root2)

	// Finality progress, mark all roots as finalized
	require.NoError(t, orm.InsertBlock(ctx, utils.RandomBytes32(), 3, time.Now(), 3))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root2)

	inputLogs = []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 3, 1, root3, block3),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 4, 1, root4, block4),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 5, 1, root5, block5),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 5, time.Now(), 3)))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root2, root3, root4, root5)

	// Mark root in the middle as executed but keep the oldest one still waiting
	rootsCache.MarkAsExecuted(root3)
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root2, root4, root5)

	// Simulate reorg by removing all unfinalized blocks
	require.NoError(t, orm.DeleteLogsAndBlocksAfter(ctx, 4))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root2)

	// Root4 comes back but with the different block_timestamp (before the reorged block)
	inputLogs = []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 4, 1, root4, newBlock4),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 5, time.Now(), 3)))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root2, root4)

	// Mark everything as executed
	rootsCache.MarkAsExecuted(root2)
	rootsCache.MarkAsExecuted(root4)
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots)
}

func Test_RootsEligibleForExecutionWithReorgs(t *testing.T) {
	ctx := testutils.Context(t)
	chainID := testutils.NewRandomEVMChainID()
	orm := logpoller.NewORM(chainID, pgtest.NewSqlxDB(t), logger.TestLogger(t))
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Hour,
		FinalityDepth:            2,
		BackfillBatchSize:        20,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(orm, nil, logger.TestLogger(t), nil, lpOpts)

	commitStoreAddr := utils.RandomAddress()

	block1 := time.Now().Add(-8 * time.Hour)
	block2 := time.Now().Add(-5 * time.Hour)
	block3 := time.Now().Add(-2 * time.Hour)
	block4 := time.Now().Add(-1 * time.Hour)

	root1 := utils.RandomBytes32()
	root2 := utils.RandomBytes32()
	root3 := utils.RandomBytes32()

	// Genesis block
	require.NoError(t, orm.InsertBlock(ctx, utils.RandomBytes32(), 1, block1, 1))
	inputLogs := []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 2, 1, root1, block2),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 2, 2, root2, block2),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 3, 1, root3, block3),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 3, time.Now(), 1)))

	feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

	commitStore, err := v1_2_0.NewCommitStore(logger.TestLogger(t), commitStoreAddr, nil, lp, feeEstimatorConfig)
	require.NoError(t, err)

	rootsCache := cache.NewCommitRootsCache(logger.TestLogger(t), commitStore, 10*time.Hour, time.Second)

	// Get all including finalized and unfinalized
	roots, err := rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2, root3)

	// Reorg everything away
	require.NoError(t, orm.DeleteLogsAndBlocksAfter(ctx, 2))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots)

	// Reinsert the logs, mark first one as finalized
	inputLogs = []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 3, 1, root1, block3),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 4, 1, root2, block4),
		createReportAcceptedLog(t, chainID, commitStoreAddr, 4, 2, root3, block4),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 5, time.Now(), 3)))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2, root3)

	// Reorg away everything except the finalized one
	require.NoError(t, orm.DeleteLogsAndBlocksAfter(ctx, 4))
	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root1)
}

// Not very likely, but let's be more defensive here and verify if cache works properly and can deal with duplicates
func Test_BlocksWithTheSameTimestamps(t *testing.T) {
	ctx := testutils.Context(t)
	chainID := testutils.NewRandomEVMChainID()
	orm := logpoller.NewORM(chainID, pgtest.NewSqlxDB(t), logger.TestLogger(t))
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Hour,
		FinalityDepth:            2,
		BackfillBatchSize:        20,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(orm, nil, logger.TestLogger(t), nil, lpOpts)

	commitStoreAddr := utils.RandomAddress()

	block := time.Now().Add(-1 * time.Hour).Truncate(time.Second)
	root1 := utils.RandomBytes32()
	root2 := utils.RandomBytes32()

	inputLogs := []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 2, 1, root1, block),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 2, time.Now(), 2)))

	feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

	commitStore, err := v1_2_0.NewCommitStore(logger.TestLogger(t), commitStoreAddr, nil, lp, feeEstimatorConfig)
	require.NoError(t, err)

	rootsCache := cache.NewCommitRootsCache(logger.TestLogger(t), commitStore, 10*time.Hour, time.Second)
	roots, err := rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root1)

	inputLogs = []logpoller.Log{
		createReportAcceptedLog(t, chainID, commitStoreAddr, 3, 1, root2, block),
	}
	require.NoError(t, orm.InsertLogsWithBlock(ctx, inputLogs, logpoller.NewLogPollerBlock(utils.RandomBytes32(), 3, time.Now(), 3)))

	roots, err = rootsCache.RootsEligibleForExecution(ctx)
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2)
}

func assertRoots(t *testing.T, roots []cciptypes.CommitStoreReport, root ...[32]byte) {
	require.Len(t, roots, len(root))
	for i, r := range root {
		require.Equal(t, r, roots[i].MerkleRoot)
	}
}

func createReportAcceptedLog(t testing.TB, chainID *big.Int, address common.Address, blockNumber int64, logIndex int64, merkleRoot common.Hash, blockTimestamp time.Time) logpoller.Log {
	tAbi, err := commit_store_1_2_0.CommitStoreMetaData.GetAbi()
	require.NoError(t, err)
	eseEvent, ok := tAbi.Events["ReportAccepted"]
	require.True(t, ok)

	gasPriceUpdates := make([]commit_store_1_2_0.InternalGasPriceUpdate, 100)
	tokenPriceUpdates := make([]commit_store_1_2_0.InternalTokenPriceUpdate, 100)

	for i := 0; i < 100; i++ {
		gasPriceUpdates[i] = commit_store_1_2_0.InternalGasPriceUpdate{
			DestChainSelector: uint64(i),
			UsdPerUnitGas:     big.NewInt(int64(i)),
		}
		tokenPriceUpdates[i] = commit_store_1_2_0.InternalTokenPriceUpdate{
			SourceToken: utils.RandomAddress(),
			UsdPerToken: big.NewInt(int64(i)),
		}
	}

	message := commit_store_1_2_0.CommitStoreCommitReport{
		PriceUpdates: commit_store_1_2_0.InternalPriceUpdates{
			TokenPriceUpdates: tokenPriceUpdates,
			GasPriceUpdates:   gasPriceUpdates,
		},
		Interval:   commit_store_1_2_0.CommitStoreInterval{Min: 1, Max: 10},
		MerkleRoot: merkleRoot,
	}

	logData, err := eseEvent.Inputs.Pack(message)
	require.NoError(t, err)

	topic0 := commit_store_1_2_0.CommitStoreReportAccepted{}.Topic()

	return logpoller.Log{
		Topics: [][]byte{
			topic0[:],
		},
		Data:           logData,
		LogIndex:       logIndex,
		BlockHash:      utils.RandomBytes32(),
		BlockNumber:    blockNumber,
		BlockTimestamp: blockTimestamp.Truncate(time.Millisecond),
		EventSig:       topic0,
		Address:        address,
		TxHash:         utils.RandomBytes32(),
		EvmChainId:     ubig.New(chainID),
	}
}
