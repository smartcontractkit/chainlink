package logprovider

import (
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"sync"
	"time"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

type recoveryQueue struct {
	// queue is a fifo queue of workIDs
	queue []string

	// workIDs is a map of workIDs to their corresponding payloads
	workIDs map[string]ocr2keepers.UpkeepPayload

	// counters is a map of upkeepIDs to the number of payloads for that upkeep
	counters map[string]int

	// visited is a map of workIDs to their corresponding visitedRecord, marked as visited when the payload was removed from the q
	visited map[string]visitedRecord

	// seen is a map of workIDs to their corresponding visitedRecord, marked as seen when the payload is added to the q
	seen map[string]visitedRecord

	// size is the maximum capacity of the queue
	//size int
	//
	//// maxSizePerUpkeep is the maximum number of payloads that can be pending for a single upkeep
	//maxSizePerUpkeep int

	maxPendingPayloadsPerUpkeep int

	lock sync.RWMutex

	lggr logger.Logger
}

func NewRecoveryQueue(lgr logger.Logger, maxPendingPayloadsPerUpkeep int) *recoveryQueue {
	return &recoveryQueue{
		queue:                       make([]string, 0),
		workIDs:                     make(map[string]ocr2keepers.UpkeepPayload),
		counters:                    make(map[string]int),
		visited:                     make(map[string]visitedRecord),
		seen:                        make(map[string]visitedRecord),
		maxPendingPayloadsPerUpkeep: maxPendingPayloadsPerUpkeep,
		lock:                        sync.RWMutex{},
		lggr:                        lgr.Named("RecoveryQueue"),
	}
}

func (q *recoveryQueue) has(workIDs ...string) []bool {
	var res []bool
	for _, workID := range workIDs {
		_, ok := q.workIDs[workID]
		res = append(res, ok)
	}
	return res
}

func (q *recoveryQueue) add(payloads ...ocr2keepers.UpkeepPayload) (int, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	added := 0

	for _, payload := range payloads {
		if _, ok := q.workIDs[payload.WorkID]; !ok {
			upkeepID := payload.UpkeepID.String()
			count := q.counters[upkeepID]

			if count >= q.maxPendingPayloadsPerUpkeep {
				continue
			}

			count++
			q.counters[upkeepID] = count

			q.queue = append(q.queue, payload.WorkID)
			q.workIDs[payload.WorkID] = payload

			q.seen[payload.WorkID] = visitedRecord{
				visitedAt: time.Now(),
				payload:   payload,
			}

			added++
		}
	}

	q.lggr.Debugw("added payloads", "numberOfPayloads", len(q.queue))

	return added, nil
}

func (q *recoveryQueue) remove(workID string) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	indexToRemove := 0

	for queueIdx, qWorkID := range q.queue {
		if qWorkID == workID {
			payload := q.workIDs[workID]

			upkeepID := payload.UpkeepID.String()

			indexToRemove = queueIdx

			delete(q.workIDs, workID)
			delete(q.seen, workID)

			count := q.counters[upkeepID]
			count--
			if count == 0 {
				delete(q.counters, upkeepID)
			} else {
				q.counters[upkeepID] = count
			}

			q.visited[workID] = visitedRecord{
				visitedAt: time.Now(),
				payload:   payload,
			}
			break
		}
	}

	q.queue = append(q.queue[:indexToRemove], q.queue[indexToRemove+1:]...)

	return nil
}

func (q *recoveryQueue) getPayloads(maxPayloads, allowedLogsPerUpkeep int) ([]ocr2keepers.UpkeepPayload, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	var res []ocr2keepers.UpkeepPayload

	upkeepCounts := map[string]int{}
	indicesToRemove := []int{}

	for queueIdx, workID := range q.queue {
		payload := q.workIDs[workID]

		upkeepID := payload.UpkeepID.String()

		upkeepCount := upkeepCounts[upkeepID]

		if upkeepCount == allowedLogsPerUpkeep {
			continue
		}

		upkeepCounts[upkeepID] = upkeepCount + 1

		indicesToRemove = append(indicesToRemove, queueIdx)

		res = append(res, payload)

		delete(q.workIDs, workID)
		delete(q.seen, workID)

		count := q.counters[upkeepID]
		count--
		if count == 0 {
			delete(q.counters, upkeepID)
		} else {
			q.counters[upkeepID] = count
		}

		q.visited[workID] = visitedRecord{
			visitedAt: time.Now(),
			payload:   payload,
		}

		if len(res) == maxPayloads {
			break
		}
	}

	for i := len(indicesToRemove) - 1; i >= 0; i-- {
		indexToRemove := indicesToRemove[i]
		q.queue = append(q.queue[:indexToRemove], q.queue[indexToRemove+1:]...)
	}

	return res, nil
}
