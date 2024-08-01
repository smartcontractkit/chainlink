package plugin

import (
	"log"
	"sort"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/random"
)

type resultAndCount struct {
	result ocr2keepers.CheckResult
	count  int
}

type performables struct {
	limit           int
	keyRandSource   [16]byte
	quorumThreshold int
	logger          *log.Logger
	resultCount     map[string]resultAndCount
}

// Performables gets quorum on agreed check results which should ultimately be
// performed within a report. It assumes only valid observations are added to it
// and simply adds all results which achieve the quorumThreshold.
// Results are agreed upon by their UniqueID() which contains all the data
// within the result.
func newPerformables(quorumThreshold int, limit int, rSrc [16]byte, logger *log.Logger) *performables {
	return &performables{
		quorumThreshold: quorumThreshold,
		limit:           limit,
		keyRandSource:   rSrc,
		logger:          logger,
		resultCount:     make(map[string]resultAndCount),
	}
}

func (p *performables) add(observation ocr2keepersv3.AutomationObservation) {
	initialCount := len(p.resultCount)
	for _, result := range observation.Performable {
		uid := result.UniqueID()
		payloadCount, ok := p.resultCount[uid]
		if !ok {
			payloadCount = resultAndCount{
				result: result,
				count:  1,
			}
		} else {
			payloadCount.count++
		}

		p.resultCount[uid] = payloadCount
	}
	p.logger.Printf("Added %d new results from %d performables", len(p.resultCount)-initialCount, len(observation.Performable))
}

func (p *performables) set(outcome *ocr2keepersv3.AutomationOutcome) {
	performable := make([]ocr2keepers.CheckResult, 0)

	// Added workIDs
	addedWid := make(map[string]bool)
	uids := make([]string, 0, len(p.resultCount))
	for uid := range p.resultCount {
		uids = append(uids, uid)
	}
	sort.Strings(uids)
	for _, uid := range uids {
		// Traverse in sorted order of UID
		payload := p.resultCount[uid]
		// For every payload that reaches threshold and workID has not been added before, add it to performables
		if payload.count >= p.quorumThreshold && !addedWid[payload.result.WorkID] {
			addedWid[payload.result.WorkID] = true
			performable = append(performable, payload.result)
		}
	}
	p.logger.Printf("Adding %d agreed performables reaching quorumThreshold %d", len(performable), p.quorumThreshold)

	// Sort by a shuffled workID.
	sort.Slice(performable, func(i, j int) bool {
		return random.ShuffleString(performable[i].WorkID, p.keyRandSource) < random.ShuffleString(performable[j].WorkID, p.keyRandSource)
	})

	if len(performable) > p.limit {
		p.logger.Printf("Limiting new performables in outcome to %d", p.limit)
		performable = performable[:p.limit]
	}
	p.logger.Printf("Setting outcome.AgreedPerformables with %d performables", len(performable))
	outcome.AgreedPerformables = performable
}
