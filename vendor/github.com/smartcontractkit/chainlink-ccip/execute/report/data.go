package report

import (
	"sort"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/plugintypes"
)

// markNewMessagesExecuted compares an execute plugin report with the commit report metadata and marks the new messages
// as executed.
func markNewMessagesExecuted(
	execReport cciptypes.ExecutePluginReportSingleChain, report plugintypes.ExecutePluginCommitData,
) plugintypes.ExecutePluginCommitData {
	// Mark new messages executed.
	for i := 0; i < len(execReport.Messages); i++ {
		report.ExecutedMessages =
			append(report.ExecutedMessages, execReport.Messages[i].Header.SequenceNumber)
	}
	sort.Slice(
		report.ExecutedMessages,
		func(i, j int) bool { return report.ExecutedMessages[i] < report.ExecutedMessages[j] })

	return report
}
