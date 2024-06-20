package cache

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// EvictionGracePeriod defines how long after the messageVisibilityInterval a root is still kept in the cache
	EvictionGracePeriod = 1 * time.Hour
	// CleanupInterval defines how often roots cache is scanned to evict stale roots
	CleanupInterval = 30 * time.Minute
)

type CommitsRootsCache interface {
	// IsSkipped returns true if the root is either executed or snoozed. Snoozing can be temporary based on the configuration
	IsSkipped(merkleRoot [32]byte) bool
	MarkAsExecuted(merkleRoot [32]byte)
	Snooze(merkleRoot [32]byte)

	// OldestRootTimestamp returns the oldest root timestamp that is not executed yet (minus 1 second).
	// If there are no roots in the queue, it returns the messageVisibilityInterval
	OldestRootTimestamp() time.Time
	// AppendUnexecutedRoot appends the root to the unexecuted roots queue to keep track of the roots that are not executed yet
	// Roots has to be added in the order they are fetched from the database
	AppendUnexecutedRoot(merkleRoot [32]byte, blockTimestamp time.Time)
}

type commitRootsCache struct {
	lggr logger.Logger
	// executedRoots is used to keep track of the roots that are executed. Roots that are considered as executed
	// when all messages are executed on the dest and matching execution state change logs are finalized
	executedRoots *cache.Cache
	// snoozedRoots is used to keep track of the roots that are temporary snoozed
	snoozedRoots *cache.Cache
	// unexecutedRootsQueue is used to keep track of the unexecuted roots in the order they are fetched from database (should be ordered by block_number, log_index)
	// First run of Exec will fill the queue with all the roots that are not executed yet within the [now-messageVisibilityInterval, now] window.
	// When a root is executed, it is removed from the queue. Next database query instead of using entire messageVisibilityInterval window
	// will use oldestRootTimestamp as the lower bound filter for block_timestamp.
	// This way we can reduce the number of database rows fetched with every OCR round.
	// We do it this way because roots for most of the cases are executed sequentially.
	// Instead of skipping snoozed roots after we fetch them from the database, we do that on the db level by narrowing the search window.
	//
	// Example
	// messageVisibilityInterval - 10 days, now - 2010-10-15
	// We fetch all the roots that within the [2010-10-05, 2010-10-15] window and load them to the queue
	// [0xA - 2010-10-10, 0xB - 2010-10-11, 0xC - 2010-10-12] -> 0xA is the oldest root
	// We executed 0xA and a couple of rounds later, we mark 0xA as executed and snoozed that forever which removes it from the queue.
	// [0xB - 2010-10-11, 0xC - 2010-10-12]
	// Now the search filter wil be 0xA timestamp -> [2010-10-11, 20-10-15]
	// If roots are executed out of order, it's not going to change anything. However, for most of the cases we have sequential root execution and that is
	// a huge improvement because we don't need to fetch all the roots from the database in every round.
	unexecutedRootsQueue *orderedmap.OrderedMap[string, time.Time]
	oldestRootTimestamp  time.Time
	rootsQueueMu         sync.RWMutex

	// Both rootSnoozedTime and messageVisibilityInterval can be kept in the commitRootsCache without need to be updated.
	// Those config properties are populates via onchain/offchain config. When changed, OCR plugin will be restarted and cache initialized with new config.
	rootSnoozedTime           time.Duration
	messageVisibilityInterval time.Duration
}

func newCommitRootsCache(
	lggr logger.Logger,
	messageVisibilityInterval time.Duration,
	rootSnoozeTime time.Duration,
	evictionGracePeriod time.Duration,
	cleanupInterval time.Duration,
) *commitRootsCache {
	executedRoots := cache.New(messageVisibilityInterval+evictionGracePeriod, cleanupInterval)
	snoozedRoots := cache.New(rootSnoozeTime, cleanupInterval)

	return &commitRootsCache{
		lggr:                      lggr,
		executedRoots:             executedRoots,
		snoozedRoots:              snoozedRoots,
		unexecutedRootsQueue:      orderedmap.New[string, time.Time](),
		rootSnoozedTime:           rootSnoozeTime,
		messageVisibilityInterval: messageVisibilityInterval,
	}
}

func NewCommitRootsCache(
	lggr logger.Logger,
	messageVisibilityInterval time.Duration,
	rootSnoozeTime time.Duration,
) *commitRootsCache {
	return newCommitRootsCache(
		lggr,
		messageVisibilityInterval,
		rootSnoozeTime,
		EvictionGracePeriod,
		CleanupInterval,
	)
}

func (s *commitRootsCache) IsSkipped(merkleRoot [32]byte) bool {
	_, snoozed := s.snoozedRoots.Get(merkleRootToString(merkleRoot))
	_, executed := s.executedRoots.Get(merkleRootToString(merkleRoot))
	return snoozed || executed
}

func (s *commitRootsCache) MarkAsExecuted(merkleRoot [32]byte) {
	prettyMerkleRoot := merkleRootToString(merkleRoot)
	s.executedRoots.SetDefault(prettyMerkleRoot, struct{}{})

	s.rootsQueueMu.Lock()
	defer s.rootsQueueMu.Unlock()
	// if there is only one root in the queue, we put its block_timestamp as oldestRootTimestamp
	if s.unexecutedRootsQueue.Len() == 1 {
		s.oldestRootTimestamp = s.unexecutedRootsQueue.Oldest().Value
	}
	s.unexecutedRootsQueue.Delete(prettyMerkleRoot)
	if head := s.unexecutedRootsQueue.Oldest(); head != nil {
		s.oldestRootTimestamp = head.Value
	}
	s.lggr.Debugw("Deleting executed root from the queue",
		"merkleRoot", prettyMerkleRoot,
		"oldestRootTimestamp", s.oldestRootTimestamp,
	)
}

func (s *commitRootsCache) Snooze(merkleRoot [32]byte) {
	s.snoozedRoots.SetDefault(merkleRootToString(merkleRoot), struct{}{})
}

func (s *commitRootsCache) OldestRootTimestamp() time.Time {
	messageVisibilityInterval := time.Now().Add(-s.messageVisibilityInterval)
	timestamp, ok := s.pickOldestRootBlockTimestamp(messageVisibilityInterval)

	if ok {
		return timestamp
	}

	s.rootsQueueMu.Lock()
	defer s.rootsQueueMu.Unlock()

	// If rootsSearchFilter is before messageVisibilityInterval, it means that we have roots that are stuck forever and will never be executed
	// In that case, we wipe out the entire queue. Next round should start from the messageVisibilityInterval and rebuild cache from scratch.
	s.unexecutedRootsQueue = orderedmap.New[string, time.Time]()
	return messageVisibilityInterval
}

func (s *commitRootsCache) pickOldestRootBlockTimestamp(messageVisibilityInterval time.Time) (time.Time, bool) {
	s.rootsQueueMu.RLock()
	defer s.rootsQueueMu.RUnlock()

	// If there are no roots in the queue, we can return the messageVisibilityInterval
	if s.oldestRootTimestamp.IsZero() {
		return messageVisibilityInterval, true
	}

	if s.oldestRootTimestamp.After(messageVisibilityInterval) {
		// Query used for fetching roots from the database is exclusive (block_timestamp > :timestamp)
		// so we need to subtract 1 second from the head timestamp to make sure that this root is included in the results
		return s.oldestRootTimestamp.Add(-time.Second), true
	}
	return time.Time{}, false
}
func (s *commitRootsCache) AppendUnexecutedRoot(merkleRoot [32]byte, blockTimestamp time.Time) {
	prettyMerkleRoot := merkleRootToString(merkleRoot)

	s.rootsQueueMu.Lock()
	defer s.rootsQueueMu.Unlock()

	// If the root is already in the queue, we must not add it to the queue
	if _, found := s.unexecutedRootsQueue.Get(prettyMerkleRoot); found {
		return
	}
	// If the root is already executed, we must not add it to the queue
	if _, executed := s.executedRoots.Get(prettyMerkleRoot); executed {
		return
	}
	// Initialize the search filter with the first root that is added to the queue
	if s.unexecutedRootsQueue.Len() == 0 {
		s.oldestRootTimestamp = blockTimestamp
	}
	s.unexecutedRootsQueue.Set(prettyMerkleRoot, blockTimestamp)
	s.lggr.Debugw("Adding unexecuted root to the queue",
		"merkleRoot", prettyMerkleRoot,
		"blockTimestamp", blockTimestamp,
		"oldestRootTimestamp", s.oldestRootTimestamp,
	)
}

func merkleRootToString(merkleRoot [32]byte) string {
	return hex.EncodeToString(merkleRoot[:])
}
