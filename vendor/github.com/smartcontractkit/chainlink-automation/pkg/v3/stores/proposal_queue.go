package stores

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

type proposalQueueRecord struct {
	proposal ocr2keepers.CoordinatedBlockProposal
	// visited is true if the record was already dequeued
	removed bool
	// createdAt is the first time the proposal was seen by the queue
	createdAt time.Time
}

// Default expiry for a proposal in the queue
// Proposals are short lived, only kept for a period of time to
// run pipeline for it once (ObservationProcessLimit)
// The same workID can get coordinated on a new block after that time and it should
// be processed on a new block
const proposalExpiry = 20 * time.Second

func (r proposalQueueRecord) expired(now time.Time, expr time.Duration) bool {
	return now.Sub(r.createdAt) > expr
}

type proposalQueue struct {
	lock    sync.RWMutex
	records map[string]proposalQueueRecord

	typeGetter types.UpkeepTypeGetter
}

var _ types.ProposalQueue = &proposalQueue{}

func NewProposalQueue(typeGetter types.UpkeepTypeGetter) *proposalQueue {
	return &proposalQueue{
		records:    map[string]proposalQueueRecord{},
		typeGetter: typeGetter,
	}
}

func (pq *proposalQueue) Enqueue(newProposals ...ocr2keepers.CoordinatedBlockProposal) error {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	for _, p := range newProposals {
		if existing, ok := pq.records[p.WorkID]; ok {
			if existing.proposal.Trigger.BlockNumber >= p.Trigger.BlockNumber {
				// Only if existing proposal is on newer or equal check block then skip this proposal
				continue
			}
		}
		pq.records[p.WorkID] = proposalQueueRecord{
			proposal:  p,
			createdAt: time.Now(),
		}
	}

	return nil
}

func (pq *proposalQueue) Dequeue(t types.UpkeepType, n int) ([]ocr2keepers.CoordinatedBlockProposal, error) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	var proposals []ocr2keepers.CoordinatedBlockProposal
	for _, record := range pq.records {
		if record.expired(time.Now(), proposalExpiry) {
			delete(pq.records, record.proposal.WorkID)
			continue
		}
		if record.removed {
			continue
		}
		if pq.typeGetter(record.proposal.UpkeepID) == t {
			proposals = append(proposals, record.proposal)
		}
	}
	if len(proposals) < n {
		n = len(proposals)
	}
	// limit the number of proposals returned
	proposals = proposals[:n]
	// mark results as removed
	for _, p := range proposals {
		proposal := pq.records[p.WorkID]
		proposal.removed = true
		pq.records[p.WorkID] = proposal
	}

	return proposals, nil
}

func (pq *proposalQueue) Size() int {
	pq.lock.RLock()
	defer pq.lock.RUnlock()

	now := time.Now()
	size := 0

	for _, record := range pq.records {
		if record.removed || record.expired(now, DefaultExpiration) {
			continue
		}
		size++
	}

	return size
}
