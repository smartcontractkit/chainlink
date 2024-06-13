package ccipocr3

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExecutePluginReport struct {
	ChainReports []ExecutionPluginReportSingleChain `json:"chainReports"`
}

type ExecutionPluginReportSingleChain struct {
	SourceChainSelector ChainSelector `json:"sourceChainSelector"`
	Messages            []CCIPMsg     `json:"messages"`
	OffchainTokenData   [][][]byte    `json:"offchainTokenData"`
	Proofs              []Bytes32     `json:"proofs"`
	ProofFlagBits       BigInt        `json:"proofFlagBits"`
}

/////////////////////////
// Execute Observation //
/////////////////////////

type ExecutePluginCommitDataWithMessages struct {
	ExecutePluginCommitData
	Messages []CCIPMsg `json:"messages"`
}

// ExecutePluginCommitData is the data that is committed to the chain.
type ExecutePluginCommitData struct {
	// SourceChain of the chain that contains the commit report.
	SourceChain ChainSelector `json:"chainSelector"`
	// Timestamp of the block that contains the commit.
	Timestamp time.Time `json:"timestamp"`
	// BlockNum of the block that contains the commit.
	BlockNum uint64 `json:"blockNum"`
	// MerkleRoot of the messages that are in this commit report.
	MerkleRoot Bytes32 `json:"merkleRoot"`
	// SequenceNumberRange of the messages that are in this commit report.
	SequenceNumberRange SeqNumRange `json:"sequenceNumberRange"`
	// ExecutedMessages are the messages in this report that have already been executed.
	ExecutedMessages []SeqNum `json:"executed"`
}

type ExecutePluginCommitObservations map[ChainSelector][]ExecutePluginCommitDataWithMessages
type ExecutePluginMessageObservations map[ChainSelector]map[SeqNum]CCIPMsg

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

func NewExecutePluginObservation(commitReports ExecutePluginCommitObservations, messages ExecutePluginMessageObservations) ExecutePluginObservation {
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

/////////////////////
// Execute Outcome //
/////////////////////

// ExecutePluginOutcome is the outcome of the ExecutePlugin.
type ExecutePluginOutcome struct {
	// PendingCommitReports are the oldest reports with pending commits. The slice is
	// sorted from oldest to newest.
	PendingCommitReports []ExecutePluginCommitDataWithMessages `json:"commitReports"`
}

func NewExecutePluginOutcome(
	pendingCommits []ExecutePluginCommitDataWithMessages,
) ExecutePluginOutcome {
	return ExecutePluginOutcome{
		PendingCommitReports: pendingCommits,
	}
}

func (o ExecutePluginOutcome) Encode() ([]byte, error) {
	return json.Marshal(o)
}

func DecodeExecutePluginOutcome(b []byte) (ExecutePluginOutcome, error) {
	o := ExecutePluginOutcome{}
	err := json.Unmarshal(b, &o)
	return o, err
}

func (o ExecutePluginOutcome) String() string {
	return fmt.Sprintf("NextCommits: %v", o.PendingCommitReports)
}
