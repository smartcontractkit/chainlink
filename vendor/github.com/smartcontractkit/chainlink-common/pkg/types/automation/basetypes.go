package automation

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const (
	checkResultDelimiter = 0x09
)

// checkResultStringTemplate is a JSON template, used for debugging purposes only.
var checkResultStringTemplate = `{
	"PipelineExecutionState":%d,
	"Retryable":%v,
	"Eligible":%v,
	"IneligibilityReason":%d,
	"UpkeepID":%s,
	"Trigger":%s,
	"WorkID":"%s",
	"GasAllocated":%d,
	"PerformData":"%s",
	"FastGasWei":%s,
	"LinkNative":%s,
	"RetryInterval":%d
}`

func init() {
	checkResultStringTemplate = strings.Replace(checkResultStringTemplate, " ", "", -1)
	checkResultStringTemplate = strings.Replace(checkResultStringTemplate, "\t", "", -1)
	checkResultStringTemplate = strings.Replace(checkResultStringTemplate, "\n", "", -1)
}

type TransmitEventType int

const (
	UnknownEvent TransmitEventType = iota
	PerformEvent
	StaleReportEvent
	ReorgReportEvent
	InsufficientFundsReportEvent
)

// UpkeepState is a final state of some unit of work.
type UpkeepState uint8

const (
	UnknownState UpkeepState = iota
	// Performed means the upkeep was performed
	Performed
	// Ineligible means the upkeep was not eligible to be performed
	Ineligible
)

// UpkeepIdentifier is a unique identifier for the upkeep, represented as uint256 in the contract.
type UpkeepIdentifier [32]byte

// String returns a base 10 numerical string representation of the upkeep identifier.
func (u UpkeepIdentifier) String() string {
	return u.BigInt().String()
}

func (u UpkeepIdentifier) BigInt() *big.Int {
	return big.NewInt(0).SetBytes(u[:])
}

// FromBigInt sets the upkeep identifier from a big.Int,
// returning true if the big.Int is valid and false otherwise.
// in case of an invalid big.Int the upkeep identifier is set to 32 zeros.
func (u *UpkeepIdentifier) FromBigInt(i *big.Int) bool {
	*u = [32]byte{}
	if i.Cmp(big.NewInt(0)) == -1 {
		return false
	}
	b := i.Bytes()
	if len(b) == 0 {
		return true
	}
	if len(b) <= 32 {
		copy(u[32-len(b):], i.Bytes())
		return true
	}
	return false
}

type BlockNumber uint64

// BlockKey represent a block (number and hash)
// NOTE: This struct is sent on the p2p network as part of observations to get quorum
// Any change here should be backwards compatible and should keep validation and
// quorum requirements in mind. Please ensure to get a proper review along with an
// upgrade plan before changing this
type BlockKey struct {
	Number BlockNumber
	Hash   [32]byte
}

type TransmitEvent struct {
	// Type describes the type of event
	Type TransmitEventType
	// TransmitBlock is the block height of the transmit event
	TransmitBlock BlockNumber
	// Confirmations is the block height behind latest
	Confirmations int64
	// TransactionHash is the hash for the transaction where the event originated
	TransactionHash [32]byte
	// UpkeepID uniquely identifies the upkeep in the registry
	UpkeepID UpkeepIdentifier
	// WorkID uniquely identifies the unit of work for the specified upkeep
	WorkID string
	// CheckBlock is the block value that the upkeep was originally checked at
	CheckBlock BlockNumber
}

// NOTE: This struct is sent on the p2p network as part of observations to get quorum
// Any change here should be backwards compatible and should keep validation and
// quorum requirements in mind. Any field that is needed to be encoded should be added
// as well to checkResultMsg struct, and to be encoded/decoded in the MarshalJSON and
// UnmarshalJSON functions. Please ensure to get a proper review along with an upgrade
// plan before changing this.
type CheckResult struct {
	// zero if success, else indicates an error code
	PipelineExecutionState uint8
	// if PipelineExecutionState is non zero, then retryable indicates that the same
	// payload can be processed again in order to get a successful execution
	Retryable bool
	// Rest of these fields are only applicable if PipelineExecutionState is zero
	// Eligible indicates whether this result is eligible to be performed
	Eligible bool
	// If result is not eligible then the reason it failed. Should be 0 if eligible
	IneligibilityReason uint8
	// Upkeep is all the information that identifies the upkeep
	UpkeepID UpkeepIdentifier
	// Trigger is the event that triggered the upkeep to be checked
	Trigger Trigger
	// WorkID represents the unit of work for the check result
	// Exploratory: Make workID an internal field and an external WorkID() function which generates WID
	WorkID string
	// GasAllocated is the gas to provide an upkeep in a report
	GasAllocated uint64
	// PerformData is the raw data returned when simulating an upkeep perform
	PerformData []byte
	// FastGasWei is the fast gas price in wei when performing this upkeep
	FastGasWei *big.Int
	// Link to native ratio to be used when performing this upkeep
	LinkNative *big.Int
	// RetryInterval is the time interval after which the same payload can be retried.
	// This field is used is special cases (such as mercury lookup), where we want to
	// have a different retry interval than the default one (30s)
	// NOTE: this field is not encoded in JSON and is only used internally
	RetryInterval time.Duration
}

// checkResultMsg is used for encoding and decoding check results.
type checkResultMsg struct {
	PipelineExecutionState uint8
	Retryable              bool
	Eligible               bool
	IneligibilityReason    uint8
	UpkeepID               UpkeepIdentifier
	Trigger                Trigger
	WorkID                 string
	GasAllocated           uint64
	PerformData            []byte
	FastGasWei             *big.Int
	LinkNative             *big.Int
}

// UniqueID returns a unique identifier for the check result.
// It is used to achieve quorum on results before being sent within a report.
func (r CheckResult) UniqueID() string {
	var resultBytes []byte

	resultBytes = append(resultBytes, r.PipelineExecutionState)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, []byte(fmt.Sprintf("%+v", r.Retryable))...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, []byte(fmt.Sprintf("%+v", r.Eligible))...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, r.IneligibilityReason)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, r.UpkeepID[:]...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, r.Trigger.BlockHash[:]...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, big.NewInt(int64(r.Trigger.BlockNumber)).Bytes()...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	if r.Trigger.LogTriggerExtension != nil {
		// Note: We encode the whole trigger extension so the behaiour of
		// LogTriggerExtentsion.BlockNumber and LogTriggerExtentsion.BlockHash should be
		// consistent across nodes when sending observations
		resultBytes = append(resultBytes, []byte(fmt.Sprintf("%+v", r.Trigger.LogTriggerExtension))...)
	}
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, r.WorkID[:]...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, big.NewInt(int64(r.GasAllocated)).Bytes()...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	resultBytes = append(resultBytes, r.PerformData[:]...)
	resultBytes = append(resultBytes, checkResultDelimiter)

	if r.FastGasWei != nil {
		resultBytes = append(resultBytes, r.FastGasWei.Bytes()...)
	}
	resultBytes = append(resultBytes, checkResultDelimiter)

	if r.LinkNative != nil {
		resultBytes = append(resultBytes, r.LinkNative.Bytes()...)
	}
	resultBytes = append(resultBytes, checkResultDelimiter)

	return fmt.Sprintf("%x", resultBytes)
}

// NOTE: this function is used for debugging purposes only.
// for encoding check results, please use the Encoder interface
func (r CheckResult) String() string {
	return fmt.Sprintf(
		checkResultStringTemplate, r.PipelineExecutionState, r.Retryable, r.Eligible,
		r.IneligibilityReason, r.UpkeepID, r.Trigger, r.WorkID, r.GasAllocated,
		hex.EncodeToString(r.PerformData), r.FastGasWei, r.LinkNative, r.RetryInterval,
	)
}

func (r CheckResult) MarshalJSON() ([]byte, error) {
	crm := &checkResultMsg{
		PipelineExecutionState: r.PipelineExecutionState,
		Retryable:              r.Retryable,
		Eligible:               r.Eligible,
		IneligibilityReason:    r.IneligibilityReason,
		UpkeepID:               r.UpkeepID,
		Trigger:                r.Trigger,
		WorkID:                 r.WorkID,
		GasAllocated:           r.GasAllocated,
		PerformData:            r.PerformData,
		FastGasWei:             r.FastGasWei,
		LinkNative:             r.LinkNative,
	}

	return json.Marshal(crm)
}

func (r *CheckResult) UnmarshalJSON(data []byte) error {
	var crm checkResultMsg

	if err := json.Unmarshal(data, &crm); err != nil {
		return err
	}

	r.PipelineExecutionState = crm.PipelineExecutionState
	r.Retryable = crm.Retryable
	r.Eligible = crm.Eligible
	r.IneligibilityReason = crm.IneligibilityReason
	r.UpkeepID = crm.UpkeepID
	r.Trigger = crm.Trigger
	r.WorkID = crm.WorkID
	r.GasAllocated = crm.GasAllocated
	r.PerformData = crm.PerformData
	r.FastGasWei = crm.FastGasWei
	r.LinkNative = crm.LinkNative

	return nil
}

// BlockHistory is a list of block keys
type BlockHistory []BlockKey

func (bh BlockHistory) Latest() (BlockKey, error) {
	if len(bh) == 0 {
		return BlockKey{}, fmt.Errorf("empty block history")
	}

	return bh[0], nil
}

type UpkeepPayload struct {
	// Upkeep is all the information that identifies the upkeep
	UpkeepID UpkeepIdentifier
	// Trigger is the event that triggered the upkeep to be checked
	Trigger Trigger
	// WorkID uniquely identifies the unit of work for the specified upkeep
	WorkID string
	// CheckData is the data used to check the upkeep
	CheckData []byte
}

// Determines whether the payload is empty, used within filtering
func (p UpkeepPayload) IsEmpty() bool {
	return p.WorkID == ""
}

// CoordinatedBlockProposal is used to represent a unit of work that can be performed
// after a check block has been coordinated between nodes.
// NOTE: This struct is sent on the p2p network as part of observations to get quorum
// Any change here should be backwards compatible and should keep validation and
// quorum requirements in mind. Please ensure to get a proper review along with an
// upgrade plan before changing this
// NOTE: Only the trigger.BlockHash and trigger.BlockNumber are coordinated across
// the network to get a quorum. WorkID is guaranteed to be correctly generated.
// Rest of the fields here SHOULD NOT BE TRUSTED as they can be manipulated by
// a single malicious node.
type CoordinatedBlockProposal struct {
	// UpkeepID is the id of the proposed upkeep
	UpkeepID UpkeepIdentifier
	// Trigger represents the event that triggered the upkeep to be checked
	Trigger Trigger
	// WorkID represents the unit of work for the coordinated proposal
	WorkID string
}

// ReportedUpkeep contains details of an upkeep for which a report was generated.
type ReportedUpkeep struct {
	// UpkeepID id of the underlying upkeep
	UpkeepID UpkeepIdentifier
	// Trigger data for the upkeep
	Trigger Trigger
	// WorkID represents the unit of work for the reported upkeep
	WorkID string
}
