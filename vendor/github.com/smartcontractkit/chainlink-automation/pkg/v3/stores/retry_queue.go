package stores

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

var (
	// DefaultExpiration is the default expiration time for items in the queue
	DefaultExpiration = 24 * time.Hour
	// RetryInterval is the default time between retries
	RetryInterval = 30 * time.Second
)

type retryQueueRecord struct {
	// payload is the desired unit of work to be retried
	payload commontypes.UpkeepPayload
	// interval is the retry interval for the payload
	interval time.Duration
	// pending is true if the item is currently being retried
	pending bool
	// createdAt is the first time the item was seen by the queue
	createdAt time.Time
	// updatedAt is the last time the item was added to the queue
	updatedAt time.Time
}

func (r retryQueueRecord) elapsed(now time.Time, expr time.Duration) bool {
	return now.Sub(r.updatedAt) > expr
}

func (r retryQueueRecord) expired(now time.Time, expr time.Duration) bool {
	return now.Sub(r.createdAt) > expr
}

type retryQueue struct {
	lggr *log.Logger

	records    map[string]retryQueueRecord
	lock       sync.RWMutex
	expiration time.Duration
	interval   time.Duration
}

var _ types.RetryQueue = (*retryQueue)(nil)

func NewRetryQueue(lggr *log.Logger) *retryQueue {
	return &retryQueue{
		lggr:       log.New(lggr.Writer(), fmt.Sprintf("[%s | retry-queue]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
		records:    map[string]retryQueueRecord{},
		lock:       sync.RWMutex{},
		expiration: DefaultExpiration,
		interval:   RetryInterval,
	}
}

func (q *retryQueue) Enqueue(records ...types.RetryRecord) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	now := time.Now()

	for _, rec := range records {
		payload := rec.Payload
		record, ok := q.records[payload.WorkID]
		if !ok {
			record = retryQueueRecord{
				payload:   payload,
				createdAt: now,
			}
		}
		if payload.Trigger.BlockNumber > record.payload.Trigger.BlockNumber {
			// new item is newer -> replace payload
			q.lggr.Printf("updating payload for workID %s on block %d", payload.WorkID, payload.Trigger.BlockNumber)
			record.payload = payload
		}
		// Enqueue the item with updatedAt = now. It will be dequeue-able after RetryInterval
		// If item was already pending, it will be eligible to get retried again
		// (can happen when the same payload gets retryable error again)
		record.updatedAt = now
		record.pending = false
		// if some custom interval is set for this record, use it.
		// otherwise use the default interval
		if rec.Interval > 0 {
			record.interval = rec.Interval
		} else {
			record.interval = q.interval
		}
		q.records[payload.WorkID] = record
	}

	return nil
}

// Dequeue returns the next n items in the queue, considering retry time schedules
// Returns only non-pending items that are within their retry interval.
//
// NOTE: Items that are expired are removed from the queue.
func (q *retryQueue) Dequeue(n int) ([]commontypes.UpkeepPayload, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	now := time.Now()

	var results []commontypes.UpkeepPayload
	for k, record := range q.records {
		if record.expired(now, q.expiration) {
			q.lggr.Printf("removing expired record %s", k)
			delete(q.records, k)
			continue
		}
		if record.pending {
			continue
		}
		if record.elapsed(now, record.interval) {
			results = append(results, record.payload)
			record.pending = true
			q.records[k] = record
			if len(results) >= n {
				break
			}
		}
	}

	if len(results) > 0 {
		q.lggr.Printf("dequeued %d payloads", len(results))
	}

	return results, nil
}

// Size returns the number of items in the queue that are not expired
func (q *retryQueue) Size() int {
	q.lock.RLock()
	defer q.lock.RUnlock()

	now := time.Now()
	size := 0

	for _, record := range q.records {
		if record.pending || record.expired(now, q.expiration) {
			continue
		}
		size++
	}

	return size
}
