package execute

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// Plugin implements the main ocr3 plugin logic.
type Plugin struct {
	reportingCfg ocr3types.ReportingPluginConfig
	cfg          cciptypes.ExecutePluginConfig
	ccipReader   cciptypes.CCIPReader

	//commitRootsCache cache.CommitsRootsCache
	lastReportTS *atomic.Int64
}

func NewPlugin(
	_ context.Context,
	reportingCfg ocr3types.ReportingPluginConfig,
	cfg cciptypes.ExecutePluginConfig,
	ccipReader cciptypes.CCIPReader,
) *Plugin {
	lastReportTS := &atomic.Int64{}
	lastReportTS.Store(time.Now().Add(-cfg.MessageVisibilityInterval).UnixMilli())

	return &Plugin{
		reportingCfg: reportingCfg,
		cfg:          cfg,
		ccipReader:   ccipReader,
		lastReportTS: lastReportTS,
	}
}

func (p *Plugin) Query(ctx context.Context, outctx ocr3types.OutcomeContext) (types.Query, error) {
	return types.Query{}, nil
}

func getPendingExecutedReports(ctx context.Context, ccipReader cciptypes.CCIPReader, dest cciptypes.ChainSelector, ts time.Time) (cciptypes.ExecutePluginCommitObservations, time.Time, error) {
	latestReportTS := time.Time{}
	commitReports, err := ccipReader.CommitReportsGTETimestamp(ctx, dest, ts, 1000)
	if err != nil {
		return nil, time.Time{}, err
	}
	// TODO: this could be more efficient. reports is also traversed in 'filterOutExecutedMessages' function.
	for _, report := range commitReports {
		if report.Timestamp.After(latestReportTS) {
			latestReportTS = report.Timestamp
		}
	}

	// TODO: this could be more efficient. commitReports is also traversed in 'groupByChainSelector'.
	for _, report := range commitReports {
		if report.Timestamp.After(latestReportTS) {
			latestReportTS = report.Timestamp
		}
	}

	groupedCommits := groupByChainSelector(commitReports)

	// Remove fully executed reports.
	for selector, reports := range groupedCommits {
		if len(reports) == 0 {
			continue
		}

		ranges, err := computeRanges(reports)
		if err != nil {
			return nil, time.Time{}, err
		}

		var executedMessages []cciptypes.SeqNumRange
		for _, seqRange := range ranges {
			executedMessagesForRange, err2 := ccipReader.ExecutedMessageRanges(ctx, selector, dest, seqRange)
			if err2 != nil {
				return nil, time.Time{}, err2
			}
			executedMessages = append(executedMessages, executedMessagesForRange...)
		}

		// Remove fully executed reports.
		groupedCommits[selector], err = filterOutExecutedMessages(reports, executedMessages)
		if err != nil {
			return nil, time.Time{}, err
		}
	}

	return groupedCommits, latestReportTS, nil
}

// Observation collects data across two phases which happen in separate rounds.
// These phases happen continuously so that except for the first round, every
// subsequent round can have a new execution report.
//
// Phase 1: Gather commit reports from the destination chain and determine
// which messages are required to build a valid execution report.
//
// Phase 2: Gather messages from the source chains and build the execution
// report.
func (p *Plugin) Observation(ctx context.Context, outctx ocr3types.OutcomeContext, _ types.Query) (types.Observation, error) {
	previousOutcome, err := cciptypes.DecodeExecutePluginOutcome(outctx.PreviousOutcome)
	if err != nil {
		return types.Observation{}, err
	}

	// Phase 1: Gather commit reports from the destination chain and determine which messages are required to build a valid execution report.
	ownConfig := p.cfg.ObserverInfo[p.reportingCfg.OracleID]
	var groupedCommits cciptypes.ExecutePluginCommitObservations
	if slices.Contains(ownConfig.Reads, p.cfg.DestChain) {
		var latestReportTS time.Time
		groupedCommits, latestReportTS, err = getPendingExecutedReports(ctx, p.ccipReader, p.cfg.DestChain, time.UnixMilli(p.lastReportTS.Load()))
		if err != nil {
			return types.Observation{}, err
		}
		// Update timestamp to the last report.
		p.lastReportTS.Store(latestReportTS.UnixMilli())

		// TODO: truncate grouped commits to a maximum observation size.
		//       Cache everything which is not executed.
	}

	// Phase 2: Gather messages from the source chains and build the execution report.
	messages := make(cciptypes.ExecutePluginMessageObservations)
	if len(previousOutcome.PendingCommitReports) == 0 {
		fmt.Println("TODO: No reports to execute. This is expected after a cold start.")
		// No reports to execute.
		// This is expected after a cold start.
	} else {
		commitReportCache := make(map[cciptypes.ChainSelector][]cciptypes.ExecutePluginCommitDataWithMessages)
		for _, report := range previousOutcome.PendingCommitReports {
			commitReportCache[report.SourceChain] = append(commitReportCache[report.SourceChain], report)
		}

		for selector, reports := range commitReportCache {
			if len(reports) == 0 {
				continue
			}

			ranges, err := computeRanges(reports)
			if err != nil {
				return types.Observation{}, err
			}

			// Read messages for each range.
			for _, seqRange := range ranges {
				msgs, err := p.ccipReader.MsgsBetweenSeqNums(ctx, selector, seqRange)
				if err != nil {
					return nil, err
				}
				for _, msg := range msgs {
					messages[selector][msg.SeqNum] = msg
				}
			}
		}
	}

	// TODO: Fire off messages for an attestation check service.

	return cciptypes.NewExecutePluginObservation(groupedCommits, messages).Encode()
}

func (p *Plugin) ValidateObservation(outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	decodedObservation, err := cciptypes.DecodeExecutePluginObservation(ao.Observation)
	if err != nil {
		return fmt.Errorf("decode observation: %w", err)
	}

	if err := validateObserverReadingEligibility(p.reportingCfg.OracleID, p.cfg.ObserverInfo, decodedObservation.Messages); err != nil {
		return fmt.Errorf("validate observer reading eligibility: %w", err)
	}

	if err := validateObservedSequenceNumbers(decodedObservation.CommitReports); err != nil {
		return fmt.Errorf("validate observed sequence numbers: %w", err)
	}

	return nil
}

func (p *Plugin) ObservationQuorum(outctx ocr3types.OutcomeContext, query types.Query) (ocr3types.Quorum, error) {
	// TODO: should we use f+1 (or less) instead of 2f+1 because it is not needed for security?
	return ocr3types.QuorumFPlusOne, nil
}

func (p *Plugin) Outcome(outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	decodedObservations, err := decodeAttributedObservations(aos)
	if err != nil {
		return ocr3types.Outcome{}, err

	}
	if len(decodedObservations) < p.reportingCfg.F {
		return ocr3types.Outcome{}, fmt.Errorf("below F threshold")
	}

	mergedCommitObservations, err := mergeCommitObservations(decodedObservations, p.cfg.FChain)
	if err != nil {
		return ocr3types.Outcome{}, err
	}

	mergedMessageObservations, err := mergeMessageObservations(decodedObservations, p.cfg.FChain)
	if err != nil {
		return ocr3types.Outcome{}, err
	}

	observation := cciptypes.NewExecutePluginObservation(
		mergedCommitObservations,
		mergedMessageObservations)

	// flatten commit reports and sort by timestamp.
	var reports []cciptypes.ExecutePluginCommitDataWithMessages
	for _, report := range observation.CommitReports {
		reports = append(reports, report...)
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Timestamp.Before(reports[j].Timestamp)
	})

	// add messages to their reports.
	for _, report := range reports {
		report.Messages = nil
		for i := report.SequenceNumberRange.Start(); i <= report.SequenceNumberRange.End(); i++ {
			if msg, ok := observation.Messages[report.SourceChain][i]; ok {
				report.Messages = append(report.Messages, msg)
			}
		}
	}

	// TODO: select reports and messages for the final exec report.
	// TODO: may only need the proofs for the final exec report rather than the report and the messages.

	return cciptypes.NewExecutePluginOutcome(reports).Encode()
}

func (p *Plugin) Reports(seqNr uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportWithInfo[[]byte], error) {

	panic("implement me")
}

func (p *Plugin) ShouldAcceptAttestedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	panic("implement me")
}

func (p *Plugin) ShouldTransmitAcceptedReport(ctx context.Context, u uint64, r ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	panic("implement me")
}

func (p *Plugin) Close() error {
	panic("implement me")
}

// Interface compatibility checks.
var _ ocr3types.ReportingPlugin[[]byte] = &Plugin{}
