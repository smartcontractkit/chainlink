package commit

import (
	"errors"
	"fmt"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/libocr/commontypes"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// validateObserverReadingEligibility checks if the observer is eligible to observe the messages it observed.
func validateObserverReadingEligibility(
	observer commontypes.OracleID,
	observerCfg map[commontypes.OracleID]cciptypes.ObserverInfo,
	observedMsgs cciptypes.ExecutePluginMessageObservations,
) error {
	observerInfo, exists := observerCfg[observer]
	if !exists {
		return fmt.Errorf("observer not found in config")
	}

	observerReadChains := mapset.NewSet(observerInfo.Reads...)

	for chainSel, msgs := range observedMsgs {
		if len(msgs) == 0 {
			continue
		}

		if !observerReadChains.Contains(chainSel) {
			return fmt.Errorf("observer not allowed to read from chain %d", chainSel)
		}
	}

	return nil
}

// validateObservedSequenceNumbers checks if the sequence numbers of the provided messages are unique for each chain and
// that they match the observed max sequence numbers.
func validateObservedSequenceNumbers(observedData map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitData) error {
	for _, commitData := range observedData {
		// observed commitData must not contain duplicates

		observedMerkleRoots := mapset.NewSet[string]()
		observedRanges := make([]cciptypes.SeqNumRange, 0)

		for _, data := range commitData {
			rootStr := data.MerkleRoot.String()
			if observedMerkleRoots.Contains(rootStr) {
				return fmt.Errorf("duplicate merkle root %s observed", rootStr)
			}
			observedMerkleRoots.Add(rootStr)

			for _, rng := range observedRanges {
				if rng.Overlaps(data.SequenceNumberRange) {
					return fmt.Errorf("sequence number range %v overlaps with %v", data.SequenceNumberRange, rng)
				}
			}
			observedRanges = append(observedRanges, data.SequenceNumberRange)

			// Executed sequence numbers should belong in the observed range.
			for _, seqNum := range data.ExecutedMessages {
				if !data.SequenceNumberRange.Contains(seqNum) {
					return fmt.Errorf("executed message %d not in observed range %v", seqNum, data.SequenceNumberRange)
				}
			}
		}
	}

	return nil
}

var errOverlappingRanges = errors.New("overlapping sequence numbers in reports")

// computeRanges takes a slice of reports and computes the smallest number of contiguous ranges
// that cover all the sequence numbers in the reports.
func computeRanges(reports []cciptypes.ExecutePluginCommitData) ([]cciptypes.SeqNumRange, error) {
	var ranges []cciptypes.SeqNumRange

	if len(reports) == 0 {
		return nil, nil
	}

	var seqRange cciptypes.SeqNumRange
	for i, report := range reports {
		if i == 0 {
			// initialize
			seqRange = cciptypes.NewSeqNumRange(report.SequenceNumberRange.Start(), report.SequenceNumberRange.End())
		} else if seqRange.End()+1 == report.SequenceNumberRange.Start() {
			// extend the contiguous range
			seqRange.SetEnd(report.SequenceNumberRange.End())
		} else if report.SequenceNumberRange.Start() < seqRange.End() {
			return nil, errOverlappingRanges
		} else {
			ranges = append(ranges, seqRange)

			// Reset the range.
			seqRange = cciptypes.NewSeqNumRange(report.SequenceNumberRange.Start(), report.SequenceNumberRange.End())
		}
	}
	// add final range
	ranges = append(ranges, seqRange)

	return ranges, nil
}

func groupByChainSelector(reports []cciptypes.CommitPluginReportWithMeta) cciptypes.ExecutePluginCommitObservations {
	commitReportCache := make(map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitData)
	for _, report := range reports {
		for _, singleReport := range report.Report.MerkleRoots {
			commitReportCache[singleReport.ChainSel] = append(commitReportCache[singleReport.ChainSel], cciptypes.ExecutePluginCommitData{
				Timestamp:           report.Timestamp,
				BlockNum:            report.BlockNum,
				MerkleRoot:          singleReport.MerkleRoot,
				SequenceNumberRange: singleReport.SeqNumsRange,
				ExecutedMessages:    nil,
			})
		}
	}
	return commitReportCache
}

// filterOutExecutedMessages returns a new reports slice with fully executed messages removed.
// Unordered inputs are supported.
func filterOutExecutedMessages(reports []cciptypes.ExecutePluginCommitData, executedMessages []cciptypes.SeqNumRange) ([]cciptypes.ExecutePluginCommitData, error) {
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].SequenceNumberRange.Start() < reports[j].SequenceNumberRange.Start()
	})

	// If none are executed, return the (sorted) input.
	if len(executedMessages) == 0 {
		return reports, nil
	}

	sort.Slice(executedMessages, func(i, j int) bool {
		return executedMessages[i].Start() < executedMessages[j].Start()
	})

	// Make sure they do not overlap
	previousMax := cciptypes.SeqNum(0)
	for _, seqRange := range executedMessages {
		if seqRange.Start() < previousMax {
			return nil, errOverlappingRanges
		}
		previousMax = seqRange.End()
	}

	var filtered []cciptypes.ExecutePluginCommitData

	reportIdx := 0
	for _, executed := range executedMessages {
		for i := reportIdx; i < len(reports); i++ {
			reportRange := reports[i].SequenceNumberRange
			if executed.End() < reportRange.Start() {
				// need to go to the next set of executed messages.
				break
			}

			if executed.End() < reportRange.Start() {
				// add report that has non-executed messages.
				reportIdx++
				filtered = append(filtered, reports[i])
				continue
			}

			if reportRange.Start() >= executed.Start() && reportRange.End() <= executed.End() {
				// skip fully executed report.
				reportIdx++
				continue
			}

			s := executed.Start()
			if reportRange.Start() > executed.Start() {
				s = reportRange.Start()
			}
			for ; s <= executed.End(); s++ {
				// This range runs into the next report.
				if s > reports[i].SequenceNumberRange.End() {
					reportIdx++
					filtered = append(filtered, reports[i])
					break
				}
				reports[i].ExecutedMessages = append(reports[i].ExecutedMessages, s)
			}
		}
	}

	// Add any remaining reports that were not fully executed.
	for i := reportIdx; i < len(reports); i++ {
		filtered = append(filtered, reports[i])
	}

	return filtered, nil
}
