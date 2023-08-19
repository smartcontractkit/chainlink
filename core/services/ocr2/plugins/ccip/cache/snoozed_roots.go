package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	// EvictionGracePeriod defines how long after the permissionless execution threshold a root is still kept in the cache
	EvictionGracePeriod = 1 * time.Hour
	// CleanupInterval defines how often roots have to be evicted
	CleanupInterval = 30 * time.Minute
)

type SnoozedRoots interface {
	IsSnoozed(merkleRoot [32]byte) bool
	MarkAsExecuted(merkleRoot [32]byte)
	Snooze(merkleRoot [32]byte)
}

type snoozedRoots struct {
	cache *cache.Cache
	// Both rootSnoozedTime and permissionLessExecutionThresholdDuration can be kept in the snoozedRoots without need to be updated.
	// Those config properties are populates via onchain/offchain config. When changed, OCR plugin will be restarted and cache initialized with new config.
	rootSnoozedTime                          time.Duration
	permissionLessExecutionThresholdDuration time.Duration
}

func newSnoozedRoots(
	permissionLessExecutionThresholdDuration time.Duration,
	rootSnoozeTime time.Duration,
	evictionGracePeriod time.Duration,
	cleanupInterval time.Duration,
) *snoozedRoots {
	evictionTime := permissionLessExecutionThresholdDuration + evictionGracePeriod
	internalCache := cache.New(evictionTime, cleanupInterval)

	return &snoozedRoots{
		cache:                                    internalCache,
		rootSnoozedTime:                          rootSnoozeTime,
		permissionLessExecutionThresholdDuration: permissionLessExecutionThresholdDuration,
	}
}

func NewSnoozedRoots(permissionLessExecutionThresholdDuration time.Duration, rootSnoozeTime time.Duration) *snoozedRoots {
	return newSnoozedRoots(permissionLessExecutionThresholdDuration, rootSnoozeTime, EvictionGracePeriod, CleanupInterval)
}

func (s *snoozedRoots) IsSnoozed(merkleRoot [32]byte) bool {
	rawValue, found := s.cache.Get(merkleRootToString(merkleRoot))
	return found && time.Now().Before(rawValue.(time.Time))
}

func (s *snoozedRoots) MarkAsExecuted(merkleRoot [32]byte) {
	s.cache.SetDefault(merkleRootToString(merkleRoot), time.Now().Add(s.permissionLessExecutionThresholdDuration))
}

func (s *snoozedRoots) Snooze(merkleRoot [32]byte) {
	s.cache.SetDefault(merkleRootToString(merkleRoot), time.Now().Add(s.rootSnoozedTime))
}

func merkleRootToString(merkleRoot [32]byte) string {
	return string(merkleRoot[:])
}
