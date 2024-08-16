package cache

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// EvictionGracePeriod defines how long after the messageVisibilityInterval a root is still kept in the cache
	EvictionGracePeriod = 1 * time.Hour
	// CleanupInterval defines how often roots cache is scanned to evict stale roots
	CleanupInterval = 30 * time.Minute
)

type CommitsRootsCache interface {
	RootsEligibleForExecution(ctx context.Context) ([]ccip.CommitStoreReport, error)
	MarkAsExecuted(merkleRoot [32]byte)
	Snooze(merkleRoot [32]byte)
}

func NewCommitRootsCache(
	lggr logger.Logger,
	reader ccip.CommitStoreReader,
	messageVisibilityInterval time.Duration,
	rootSnoozeTime time.Duration,
) CommitsRootsCache {
	return newCommitRootsCache(
		lggr,
		reader,
		messageVisibilityInterval,
		rootSnoozeTime,
		CleanupInterval,
		EvictionGracePeriod,
	)
}

func newCommitRootsCache(
	lggr logger.Logger,
	reader ccip.CommitStoreReader,
	messageVisibilityInterval time.Duration,
	rootSnoozeTime time.Duration,
	cleanupInterval time.Duration,
	evictionGracePeriod time.Duration,
) *commitRootsCache {
	snoozedRoots := cache.New(rootSnoozeTime, cleanupInterval)
	executedRoots := cache.New(messageVisibilityInterval+evictionGracePeriod, cleanupInterval)

	return &commitRootsCache{
		lggr:                        lggr,
		reader:                      reader,
		rootSnoozeTime:              rootSnoozeTime,
		finalizedRoots:              orderedmap.New[string, ccip.CommitStoreReportWithTxMeta](),
		executedRoots:               executedRoots,
		snoozedRoots:                snoozedRoots,
		messageVisibilityInterval:   messageVisibilityInterval,
		latestFinalizedCommitRootTs: time.Now().Add(-messageVisibilityInterval),
		cacheMu:                     sync.RWMutex{},
	}
}

type commitRootsCache struct {
	lggr                      logger.Logger
	reader                    ccip.CommitStoreReader
	messageVisibilityInterval time.Duration
	rootSnoozeTime            time.Duration

	// Mutable state. finalizedRoots is thread-safe by default, but updating latestFinalizedCommitRootTs and finalizedRoots requires locking.
	cacheMu sync.RWMutex
	// finalizedRoots is a map of merkleRoot -> CommitStoreReportWithTxMeta. It stores all the CommitReports that are
	// marked as finalized by LogPoller, but not executed yet. Keeping only finalized reports doesn't require any state sync between LP and the cache.
	// In order to keep this map size under control, we evict stale items every time we fetch new logs from the database.
	// Also, ccip.CommitStoreReportWithTxMeta is a very tiny entity with almost fixed size, so it's not a big deal to keep it in memory.
	// In case of high memory footprint caused by storing roots, we can make these even more lightweight by removing token/gas price updates.
	// Whenever the root is executed (all messages executed and ExecutionStateChange events are finalized), we remove the root from the map.
	finalizedRoots *orderedmap.OrderedMap[string, ccip.CommitStoreReportWithTxMeta]
	// snoozedRoots used only for temporary snoozing roots. It's a cache with TTL (usually around 5 minutes, but this configuration is set up on chain using rootSnoozeTime)
	snoozedRoots *cache.Cache
	// executedRoots is a cache with TTL (usually around 8 hours, but this configuration is set up on chain using messageVisibilityInterval).
	// We keep executed roots there to make sure we don't accidentally try to reprocess already executed CommitReport
	executedRoots *cache.Cache
	// latestFinalizedCommitRootTs is the timestamp of the latest finalized commit root (youngest in terms of timestamp).
	// It's used get only the logs that were considered as unfinalized in a previous run.
	// This way we limit database scans to the minimum and keep polling "unfinalized" part of the ReportAccepted events queue.
	latestFinalizedCommitRootTs time.Time
}

func (r *commitRootsCache) RootsEligibleForExecution(ctx context.Context) ([]ccip.CommitStoreReport, error) {
	// 1. Fetch all the logs from the database after the latest finalized commit root timestamp.
	// If this is a first run, it will fetch all the logs based on the messageVisibilityInterval.
	// Worst case scenario, it will fetch around 480 reports (OCR Commit 60 seconds (fast chains default) * messageVisibilityInterval set to 8 hours (mainnet default))
	// Even with the larger messageVisibilityInterval window (e.g. 24 hours) it should be acceptable (around 1500 logs).
	// Keep in mind that this potentially heavy operation happens only once during the plugin boot and it's no different from the previous implementation.
	logs, err := r.fetchLogsFromCommitStore(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Iterate over the logs and check if the root is finalized or not. Return finalized and unfinalized reports
	// It promotes finalized roots to the finalizedRoots map and evicts stale roots.
	finalizedReports, unfinalizedReports := r.updateFinalizedRoots(logs)

	// 3. Join finalized commit reports with unfinalized reports and outfilter snoozed roots.
	// Return only the reports that are not snoozed.
	return r.pickReadyToExecute(finalizedReports, unfinalizedReports), nil
}

// MarkAsExecuted marks the root as executed. It means that all the messages from the root were executed and the ExecutionStateChange event was finalized.
// Executed roots are removed from the cache.
func (r *commitRootsCache) MarkAsExecuted(merkleRoot [32]byte) {
	prettyMerkleRoot := merkleRootToString(merkleRoot)
	r.lggr.Infow("Marking root as executed and removing entirely from cache", "merkleRoot", prettyMerkleRoot)

	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()
	r.finalizedRoots.Delete(prettyMerkleRoot)
	r.executedRoots.SetDefault(prettyMerkleRoot, struct{}{})
}

// Snooze temporarily snoozes the root. It means that the root is not eligible for execution for a certain period of time.
// Snoozed roots are skipped when calling RootsEligibleForExecution
func (r *commitRootsCache) Snooze(merkleRoot [32]byte) {
	prettyMerkleRoot := merkleRootToString(merkleRoot)
	r.lggr.Infow("Snoozing root temporarily", "merkleRoot", prettyMerkleRoot, "rootSnoozeTime", r.rootSnoozeTime)
	r.snoozedRoots.SetDefault(prettyMerkleRoot, struct{}{})
}

func (r *commitRootsCache) isSnoozed(merkleRoot [32]byte) bool {
	_, snoozed := r.snoozedRoots.Get(merkleRootToString(merkleRoot))
	return snoozed
}

func (r *commitRootsCache) isExecuted(merkleRoot [32]byte) bool {
	_, executed := r.executedRoots.Get(merkleRootToString(merkleRoot))
	return executed
}

func (r *commitRootsCache) fetchLogsFromCommitStore(ctx context.Context) ([]ccip.CommitStoreReportWithTxMeta, error) {
	r.cacheMu.Lock()
	messageVisibilityWindow := time.Now().Add(-r.messageVisibilityInterval)
	if r.latestFinalizedCommitRootTs.Before(messageVisibilityWindow) {
		r.latestFinalizedCommitRootTs = messageVisibilityWindow
	}
	commitRootsFilterTimestamp := r.latestFinalizedCommitRootTs
	r.cacheMu.Unlock()

	// IO operation, release lock before!
	r.lggr.Infow("Fetching Commit Reports with timestamp greater than or equal to", "blockTimestamp", commitRootsFilterTimestamp)
	return r.reader.GetAcceptedCommitReportsGteTimestamp(ctx, commitRootsFilterTimestamp, 0)
}

func (r *commitRootsCache) updateFinalizedRoots(logs []ccip.CommitStoreReportWithTxMeta) ([]ccip.CommitStoreReportWithTxMeta, []ccip.CommitStoreReportWithTxMeta) {
	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()

	// Assuming logs are properly ordered by block_timestamp, log_index
	var unfinalizedReports []ccip.CommitStoreReportWithTxMeta
	for _, log := range logs {
		prettyMerkleRoot := merkleRootToString(log.MerkleRoot)
		// Defensive check, if something is marked as executed, never allow it to come back to the cache
		if r.isExecuted(log.MerkleRoot) {
			r.lggr.Debugw("Ignoring root marked as executed", "merkleRoot", prettyMerkleRoot, "blockTimestamp", log.BlockTimestampUnixMilli)
			continue
		}

		if log.IsFinalized() {
			r.lggr.Debugw("Adding finalized root to cache", "merkleRoot", prettyMerkleRoot, "blockTimestamp", log.BlockTimestampUnixMilli)
			r.finalizedRoots.Store(prettyMerkleRoot, log)
		} else {
			r.lggr.Debugw("Bypassing unfinalized root", "merkleRoot", prettyMerkleRoot, "blockTimestamp", log.BlockTimestampUnixMilli)
			unfinalizedReports = append(unfinalizedReports, log)
		}
	}

	if newest := r.finalizedRoots.Newest(); newest != nil {
		r.latestFinalizedCommitRootTs = time.UnixMilli(newest.Value.BlockTimestampUnixMilli)
	}

	var finalizedRoots []ccip.CommitStoreReportWithTxMeta
	var rootsToDelete []string

	messageVisibilityWindow := time.Now().Add(-r.messageVisibilityInterval)
	for pair := r.finalizedRoots.Oldest(); pair != nil; pair = pair.Next() {
		// Mark items as stale if they are older than the messageVisibilityInterval
		// SortedMap doesn't allow to iterate and delete, so we mark roots for deletion and remove them in a separate loop
		if time.UnixMilli(pair.Value.BlockTimestampUnixMilli).Before(messageVisibilityWindow) {
			rootsToDelete = append(rootsToDelete, pair.Key)
			continue
		}
		finalizedRoots = append(finalizedRoots, pair.Value)
	}

	// Remove stale items
	for _, root := range rootsToDelete {
		r.finalizedRoots.Delete(root)
	}

	return finalizedRoots, unfinalizedReports
}

func (r *commitRootsCache) pickReadyToExecute(r1 []ccip.CommitStoreReportWithTxMeta, r2 []ccip.CommitStoreReportWithTxMeta) []ccip.CommitStoreReport {
	allReports := append(r1, r2...)
	eligibleReports := make([]ccip.CommitStoreReport, 0, len(allReports))
	for _, report := range allReports {
		if r.isSnoozed(report.MerkleRoot) {
			r.lggr.Debugw("Skipping snoozed root",
				"minSeqNr", report.Interval.Min,
				"maxSeqNr", report.Interval.Max,
				"merkleRoot", merkleRootToString(report.MerkleRoot))
			continue
		}
		eligibleReports = append(eligibleReports, report.CommitStoreReport)
	}
	// safety check, probably not needed
	slices.SortFunc(eligibleReports, func(i, j ccip.CommitStoreReport) int {
		return int(i.Interval.Min - j.Interval.Min)
	})
	return eligibleReports
}

// internal use only for testing
func (r *commitRootsCache) finalizedCachedLogs() []ccip.CommitStoreReport {
	r.cacheMu.RLock()
	defer r.cacheMu.RUnlock()

	var finalizedRoots []ccip.CommitStoreReport
	for pair := r.finalizedRoots.Oldest(); pair != nil; pair = pair.Next() {
		finalizedRoots = append(finalizedRoots, pair.Value.CommitStoreReport)
	}
	return finalizedRoots
}

func merkleRootToString(merkleRoot ccip.Hash) string {
	return merkleRoot.String()
}
