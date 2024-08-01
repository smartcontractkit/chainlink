package hooks

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	ocr2keepersv3 "github.com/smartcontractkit/chainlink-automation/pkg/v3"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/random"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
)

type AddFromStagingHook struct {
	store  types.ResultStore
	logger *log.Logger
	coord  types.Coordinator
	sorter stagedResultSorter
}

func NewAddFromStagingHook(store types.ResultStore, coord types.Coordinator, logger *log.Logger) AddFromStagingHook {
	return AddFromStagingHook{
		store:  store,
		coord:  coord,
		logger: log.New(logger.Writer(), fmt.Sprintf("[%s | build hook:add-from-staging]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
		sorter: stagedResultSorter{
			shuffledIDs: make(map[string]string),
		},
	}
}

// RunHook adds results from the store to the observation.
// It sorts by a shuffled workID. workID for all items is shuffled using a pseudorandom source
// that is the same across all nodes for a given round. This ensures that all nodes try to
// send the same subset of workIDs if they are available, while giving different priority
// to workIDs in different rounds.
func (hook *AddFromStagingHook) RunHook(obs *ocr2keepersv3.AutomationObservation, limit int, rSrc [16]byte) error {
	results, err := hook.store.View()
	if err != nil {
		return err
	}
	results, err = hook.coord.FilterResults(results)
	if err != nil {
		return err
	}

	b, err := obs.Encode()
	if err != nil {
		return err
	}

	results = hook.sorter.orderResults(results, rSrc)
	added, _ := hook.addByPercentageExceeded(obs, limit, results, len(b))

	hook.logger.Printf("skipped %d available results in staging", len(results)-added)

	hook.logger.Printf("adding %d results to observation", added)

	return nil
}

func (hook *AddFromStagingHook) addByPercentageExceeded(obs *ocr2keepersv3.AutomationObservation, limit int, results []automation.CheckResult, baseSize int) (int, int) {
	if limit > len(results) {
		limit = len(results)
	}

	if limit <= 0 {
		return len(obs.Performable), 0
	}

	obs.Performable = results[:limit]

	encodingCalls := 1
	b, _ := obs.Encode()

	if observationSize := len(b); observationSize > ocr2keepersv3.MaxObservationLength {
		performablesSize := observationSize - baseSize
		avgPerformableSize := performablesSize / limit
		exceededBy := observationSize - ocr2keepersv3.MaxObservationLength
		avgPerformablesExceeded := int(math.Ceil(float64(exceededBy) / float64(avgPerformableSize)))
		limit -= avgPerformablesExceeded + 1 // ensure we always remove at least one performable on the next call
		if limit <= 0 {
			return len(obs.Performable), encodingCalls
		}
		added, numEncodings := hook.addByPercentageExceeded(obs, limit, results, baseSize)
		return added, numEncodings + encodingCalls
	}

	return len(obs.Performable), encodingCalls
}

type stagedResultSorter struct {
	lastRandSrc [16]byte
	shuffledIDs map[string]string
	lock        sync.Mutex
}

// orderResults orders the results by the shuffled workID
func (sorter *stagedResultSorter) orderResults(results []automation.CheckResult, rSrc [16]byte) []automation.CheckResult {
	sorter.lock.Lock()
	defer sorter.lock.Unlock()

	shuffledIDs := sorter.updateShuffledIDs(results, rSrc)
	// sort by the shuffled workID
	sort.Slice(results, func(i, j int) bool {
		return shuffledIDs[results[i].WorkID] < shuffledIDs[results[j].WorkID]
	})

	return results
}

// updateShuffledIDs updates the shuffledIDs cache with the new random source or items.
// NOTE: This function is not thread-safe and should be called with a lock
func (sorter *stagedResultSorter) updateShuffledIDs(results []automation.CheckResult, rSrc [16]byte) map[string]string {
	// once the random source changes, the workIDs needs to be shuffled again with the new source
	if !bytes.Equal(sorter.lastRandSrc[:], rSrc[:]) {
		sorter.lastRandSrc = rSrc
		sorter.shuffledIDs = make(map[string]string)
	}

	for _, result := range results {
		if _, ok := sorter.shuffledIDs[result.WorkID]; !ok {
			sorter.shuffledIDs[result.WorkID] = random.ShuffleString(result.WorkID, rSrc)
		}
	}

	return sorter.shuffledIDs
}
