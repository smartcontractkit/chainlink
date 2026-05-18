package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

func Test_CacheIsInitializedWithFirstCall(t *testing.T) {
	commitStoreReader := mocks.NewCommitStoreReader(t)
	cache := newCommitRootsCache(logger.TestLogger(t), commitStoreReader, time.Hour, time.Hour, time.Hour, time.Hour)
	commitStoreReader.On("GetAcceptedCommitReportsGteTimestamp", mock.Anything, mock.Anything, mock.Anything).Return([]ccip.CommitStoreReportWithTxMeta{}, nil)

	roots, err := cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots)
}

func Test_CacheExpiration(t *testing.T) {
	ts1 := time.Now().Add(-5 * time.Millisecond).Truncate(time.Millisecond)
	ts2 := time.Now().Add(-3 * time.Millisecond).Truncate(time.Millisecond)
	ts3 := time.Now().Add(-1 * time.Millisecond).Truncate(time.Millisecond)

	root1 := utils.RandomBytes32()
	root2 := utils.RandomBytes32()
	root3 := utils.RandomBytes32()

	commitStoreReader := mocks.NewCommitStoreReader(t)
	cache := newCommitRootsCache(logger.TestLogger(t), commitStoreReader, time.Second, time.Hour, time.Hour, time.Hour)
	mockCommitStoreReader(commitStoreReader, time.Time{}, []ccip.CommitStoreReportWithTxMeta{
		createCommitStoreEntry(root1, ts1, true),
		createCommitStoreEntry(root2, ts2, true),
		createCommitStoreEntry(root3, ts3, false),
	})
	roots, err := cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2, root3)

	require.Eventually(t, func() bool {
		mockCommitStoreReader(commitStoreReader, time.Time{}, []ccip.CommitStoreReportWithTxMeta{
			createCommitStoreEntry(root3, ts3, false),
		})
		roots, err = cache.RootsEligibleForExecution(tests.Context(t))
		require.NoError(t, err)
		return len(roots) == 1 && roots[0].MerkleRoot == root3
	}, 5*time.Second, 1*time.Second)
}

func Test_CacheFullEviction(t *testing.T) {
	commitStoreReader := mocks.NewCommitStoreReader(t)
	cache := newCommitRootsCache(logger.TestLogger(t), commitStoreReader, 2*time.Second, 1*time.Second, time.Second, time.Second)

	maxElements := 10000
	commitRoots := make([]ccip.CommitStoreReportWithTxMeta, maxElements)
	for i := 0; i < maxElements; i++ {
		finalized := i >= maxElements/2
		commitRoots[i] = createCommitStoreEntry(utils.RandomBytes32(), time.Now(), finalized)
	}
	mockCommitStoreReader(commitStoreReader, time.Time{}, commitRoots)

	roots, err := cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	require.Len(t, roots, maxElements)

	// Marks some of them as exeucted and some of them as snoozed
	for i := 0; i < maxElements; i++ {
		if i%3 == 0 {
			cache.MarkAsExecuted(commitRoots[i].MerkleRoot)
		}
		if i%3 == 1 {
			cache.Snooze(commitRoots[i].MerkleRoot)
		}
	}
	// Eventually everything should be entirely removed from cache. We need that check to verify if cache doesn't grow indefinitely
	require.Eventually(t, func() bool {
		mockCommitStoreReader(commitStoreReader, time.Time{}, []ccip.CommitStoreReportWithTxMeta{})
		roots1, err1 := cache.RootsEligibleForExecution(tests.Context(t))
		require.NoError(t, err1)

		return len(roots1) == 0 &&
			cache.finalizedRoots.Len() == 0 &&
			len(cache.snoozedRoots.Items()) == 0 &&
			len(cache.executedRoots.Items()) == 0
	}, 10*time.Second, time.Second)
}

func Test_CacheProgression_Internal(t *testing.T) {
	ts1 := time.Now().Add(-5 * time.Hour).Truncate(time.Millisecond)
	ts2 := time.Now().Add(-3 * time.Hour).Truncate(time.Millisecond)
	ts3 := time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)

	root1 := utils.RandomBytes32()
	root2 := utils.RandomBytes32()
	root3 := utils.RandomBytes32()

	commitStoreReader := mocks.NewCommitStoreReader(t)

	cache := newCommitRootsCache(logger.TestLogger(t), commitStoreReader, 10*time.Hour, time.Hour, time.Hour, time.Hour)

	// Empty cache, no results from the reader
	mockCommitStoreReader(commitStoreReader, time.Time{}, []ccip.CommitStoreReportWithTxMeta{})
	roots, err := cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots)
	assertRoots(t, cache.finalizedCachedLogs())

	// Single unfinalized root returned
	mockCommitStoreReader(commitStoreReader, time.Time{}, []ccip.CommitStoreReportWithTxMeta{createCommitStoreEntry(root1, ts1, false)})
	roots, err = cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots, root1)
	assertRoots(t, cache.finalizedCachedLogs())

	// Finalized and unfinalized roots returned
	mockCommitStoreReader(commitStoreReader, time.Time{}, []ccip.CommitStoreReportWithTxMeta{
		createCommitStoreEntry(root1, ts1, true),
		createCommitStoreEntry(root2, ts2, false),
	})
	roots, err = cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2)
	assertRoots(t, cache.finalizedCachedLogs(), root1)

	// Returning the same data should not impact cache state (no duplicates)
	mockCommitStoreReader(commitStoreReader, ts1, []ccip.CommitStoreReportWithTxMeta{
		createCommitStoreEntry(root1, ts1, true),
		createCommitStoreEntry(root2, ts2, false),
	})
	roots, err = cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots, root1, root2)
	assertRoots(t, cache.finalizedCachedLogs(), root1)

	// Snoozing oldest root
	cache.Snooze(root1)
	mockCommitStoreReader(commitStoreReader, ts1, []ccip.CommitStoreReportWithTxMeta{
		createCommitStoreEntry(root2, ts2, false),
		createCommitStoreEntry(root3, ts3, false),
	})
	roots, err = cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots, root2, root3)
	assertRoots(t, cache.finalizedCachedLogs(), root1)

	// Snoozing everything
	cache.Snooze(root2)
	cache.Snooze(root3)
	mockCommitStoreReader(commitStoreReader, ts1, []ccip.CommitStoreReportWithTxMeta{
		createCommitStoreEntry(root2, ts2, true),
		createCommitStoreEntry(root3, ts3, true),
	})
	roots, err = cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots)
	assertRoots(t, cache.finalizedCachedLogs(), root1, root2, root3)

	// Marking everything as executed removes it entirely, even if root is returned from the CommitStore
	cache.MarkAsExecuted(root1)
	cache.MarkAsExecuted(root2)
	cache.MarkAsExecuted(root3)
	mockCommitStoreReader(commitStoreReader, ts3, []ccip.CommitStoreReportWithTxMeta{
		createCommitStoreEntry(root2, ts2, true),
		createCommitStoreEntry(root3, ts3, true),
	})
	roots, err = cache.RootsEligibleForExecution(tests.Context(t))
	require.NoError(t, err)
	assertRoots(t, roots)
	assertRoots(t, cache.finalizedCachedLogs())
}

func assertRoots(t *testing.T, reports []ccip.CommitStoreReport, expectedRoots ...[32]byte) {
	require.Len(t, reports, len(expectedRoots))
	for i, report := range reports {
		assert.Equal(t, expectedRoots[i], report.MerkleRoot)
	}
}

func mockCommitStoreReader(reader *mocks.CommitStoreReader, blockTimestamp time.Time, roots []ccip.CommitStoreReportWithTxMeta) {
	if blockTimestamp.IsZero() {
		reader.On("GetAcceptedCommitReportsGteTimestamp", mock.Anything, mock.Anything, mock.Anything).
			Return(roots, nil).Once()
	} else {
		reader.On("GetAcceptedCommitReportsGteTimestamp", mock.Anything, blockTimestamp, mock.Anything).
			Return(roots, nil).Once()
	}
}

func createCommitStoreEntry(root [32]byte, ts time.Time, finalized bool) ccip.CommitStoreReportWithTxMeta {
	status := ccip.FinalizedStatusNotFinalized
	if finalized {
		status = ccip.FinalizedStatusFinalized
	}
	return ccip.CommitStoreReportWithTxMeta{
		CommitStoreReport: ccip.CommitStoreReport{
			MerkleRoot: root,
		},
		TxMeta: ccip.TxMeta{
			BlockTimestampUnixMilli: ts.UnixMilli(),
			Finalized:               status,
		},
	}
}
