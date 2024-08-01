package plugintypes

import (
	"encoding/json"
	"sort"
	"time"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// ///////////////////////
// Execute Observation //
// ///////////////////////

// ExecutePluginCommitData is the data that is committed to the chain.
type ExecutePluginCommitData struct {
	// SourceChain of the chain that contains the commit report.
	SourceChain cciptypes.ChainSelector `json:"chainSelector"`
	// Timestamp of the block that contains the commit.
	Timestamp time.Time `json:"timestamp"`
	// BlockNum of the block that contains the commit.
	BlockNum uint64 `json:"blockNum"`
	// MerkleRoot of the messages that are in this commit report.
	MerkleRoot cciptypes.Bytes32 `json:"merkleRoot"`
	// SequenceNumberRange of the messages that are in this commit report.
	SequenceNumberRange cciptypes.SeqNumRange `json:"sequenceNumberRange"`

	// Messages that are part of the commit report.
	Messages []cciptypes.Message `json:"messages"`

	// ExecutedMessages are the messages in this report that have already been executed.
	ExecutedMessages []cciptypes.SeqNum `json:"executedMessages"`

	// TokenData for each message.
	TokenData [][][]byte `json:"-"`
}

type ExecutePluginCommitObservations map[cciptypes.ChainSelector][]ExecutePluginCommitData
type ExecutePluginMessageObservations map[cciptypes.ChainSelector]map[cciptypes.SeqNum]cciptypes.Message

// ExecutePluginObservation is the observation of the ExecutePlugin.
// TODO: revisit observation types. The maps used here are more space efficient and easier to work
// with but require more transformations compared to the on-chain representations.
type ExecutePluginObservation struct {
	// CommitReports are determined during the first phase of execute.
	// It contains the commit reports we would like to execute in the following round.
	CommitReports ExecutePluginCommitObservations `json:"commitReports"`
	// Messages are determined during the second phase of execute.
	// Ideally, it contains all the messages identified by the previous outcome's
	// NextCommits. With the previous outcome, and these messsages, we can build the
	// execute report.
	Messages ExecutePluginMessageObservations `json:"messages"`
	// TODO: some of the nodes configuration may need to be included here.
}

func NewExecutePluginObservation(
	commitReports ExecutePluginCommitObservations, messages ExecutePluginMessageObservations) ExecutePluginObservation {
	return ExecutePluginObservation{
		CommitReports: commitReports,
		Messages:      messages,
	}
}

func (obs ExecutePluginObservation) Encode() ([]byte, error) {
	return json.Marshal(obs)
}

func DecodeExecutePluginObservation(b []byte) (ExecutePluginObservation, error) {
	obs := ExecutePluginObservation{}
	err := json.Unmarshal(b, &obs)
	return obs, err
}

// ///////////////////
// Execute Outcome //
// ///////////////////

// ExecutePluginOutcome is the outcome of the ExecutePlugin.
type ExecutePluginOutcome struct {
	// PendingCommitReports are the oldest reports with pending commits. The slice is
	// sorted from oldest to newest.
	PendingCommitReports []ExecutePluginCommitData `json:"commitReports"`

	// Report is built from the oldest pending commit reports.
	Report cciptypes.ExecutePluginReport `json:"report"`
}

func (o ExecutePluginOutcome) IsEmpty() bool {
	return len(o.PendingCommitReports) == 0 && len(o.Report.ChainReports) == 0
}

func NewExecutePluginOutcome(
	pendingCommits []ExecutePluginCommitData,
	report cciptypes.ExecutePluginReport,
) ExecutePluginOutcome {
	return newSortedExecuteOutcome(pendingCommits, report)
}

// Encode encodes the outcome by first sorting the pending commit reports and the chain reports
// and then JSON marshalling.
// The encoding MUST be deterministic.
func (o ExecutePluginOutcome) Encode() ([]byte, error) {
	// We sort again here in case construction is not via the constructor.
	return json.Marshal(newSortedExecuteOutcome(o.PendingCommitReports, o.Report))
}

func newSortedExecuteOutcome(
	pendingCommits []ExecutePluginCommitData,
	report cciptypes.ExecutePluginReport) ExecutePluginOutcome {
	pendingCommitsCP := append([]ExecutePluginCommitData{}, pendingCommits...)
	reportCP := append([]cciptypes.ExecutePluginReportSingleChain{}, report.ChainReports...)
	sort.Slice(
		pendingCommitsCP,
		func(i, j int) bool {
			return pendingCommitsCP[i].SourceChain < pendingCommitsCP[j].SourceChain
		})
	sort.Slice(
		reportCP,
		func(i, j int) bool {
			return reportCP[i].SourceChainSelector < reportCP[j].SourceChainSelector
		})
	return ExecutePluginOutcome{
		PendingCommitReports: pendingCommitsCP,
		Report:               cciptypes.ExecutePluginReport{ChainReports: reportCP},
	}
}

func DecodeExecutePluginOutcome(b []byte) (ExecutePluginOutcome, error) {
	o := ExecutePluginOutcome{}
	err := json.Unmarshal(b, &o)
	return o, err
}
