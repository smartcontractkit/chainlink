package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	clnull "github.com/smartcontractkit/chainlink-common/pkg/utils/null"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

// TxStrategy controls how txes are queued and sent
type TxStrategy interface {
	// Subject will be saved txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after tx insertion
	// It accepts the service responsible for deleting
	// unstarted txs and deletion options
	PruneQueue(ctx context.Context, pruneService UnstartedTxQueuePruner) (ids []int64, err error)
}

type TxAttemptState int8

type TxState string

const (
	TxAttemptInProgress TxAttemptState = iota + 1
	TxAttemptInsufficientFunds
	TxAttemptBroadcast
	txAttemptStateCount // always at end to calculate number of states
)

var txAttemptStateStrings = []string{
	"unknown_attempt_state",    // default 0 value
	TxAttemptInProgress:        "in_progress",
	TxAttemptInsufficientFunds: "insufficient_funds",
	TxAttemptBroadcast:         "broadcast",
}

func NewTxAttemptState(state string) (s TxAttemptState) {
	if index := slices.Index(txAttemptStateStrings, state); index != -1 {
		s = TxAttemptState(index)
	}
	return s
}

// String returns string formatted states for logging
func (s TxAttemptState) String() (str string) {
	if s < txAttemptStateCount {
		return txAttemptStateStrings[s]
	}
	return txAttemptStateStrings[0]
}

type TxRequest[ADDR types.Hashable, TX_HASH types.Hashable] struct {
	// IdempotencyKey is a globally unique ID set by the caller, to prevent accidental creation of duplicated Txs during retries or crash recovery.
	// If this field is set, the TXM will first search existing Txs with this field.
	// If found, it will return the existing Tx, without creating a new one. TXM will not validate or ensure that existing Tx is same as the incoming TxRequest.
	// If not found, TXM will create a new Tx.
	// If IdempotencyKey is set to null, TXM will always create a new Tx.
	// Since IdempotencyKey has to be globally unique, consider prepending the service or component's name it is being used by
	// Such as {service}-{ID}. E.g vrf-12345
	IdempotencyKey   *string
	FromAddress      ADDR
	ToAddress        ADDR
	EncodedPayload   []byte
	Value            big.Int
	FeeLimit         uint64
	Meta             *TxMeta[ADDR, TX_HASH]
	ForwarderAddress ADDR

	// Pipeline variables - if you aren't calling this from chain tx task within
	// the pipeline, you don't need these variables
	MinConfirmations  clnull.Uint32
	PipelineTaskRunID *uuid.UUID

	Strategy TxStrategy

	// Checker defines the check that should be run before a transaction is submitted on chain.
	Checker TransmitCheckerSpec[ADDR]

	// Mark tx requiring callback
	SignalCallback bool
}

// TransmitCheckerSpec defines the check that should be performed before a transaction is submitted
// on chain.
type TransmitCheckerSpec[ADDR types.Hashable] struct {
	// CheckerType is the type of check that should be performed. Empty indicates no check.
	CheckerType TransmitCheckerType `json:",omitempty"`

	// VRFCoordinatorAddress is the address of the VRF coordinator that should be used to perform
	// VRF transmit checks. This should be set iff CheckerType is TransmitCheckerTypeVRFV2.
	VRFCoordinatorAddress *ADDR `json:",omitempty"`

	// VRFRequestBlockNumber is the block number in which the provided VRF request has been made.
	// This should be set iff CheckerType is TransmitCheckerTypeVRFV2.
	VRFRequestBlockNumber *big.Int `json:",omitempty"`
}

// TransmitCheckerType describes the type of check that should be performed before a transaction is
// executed on-chain.
type TransmitCheckerType string

// TxMeta contains fields of the transaction metadata
// Not all fields are guaranteed to be present
type TxMeta[ADDR types.Hashable, TX_HASH types.Hashable] struct {
	JobID *int32 `json:"JobID,omitempty"`

	// Pipeline fields
	FailOnRevert null.Bool `json:"FailOnRevert,omitempty"`

	// VRF-only fields
	RequestID     *TX_HASH `json:"RequestID,omitempty"`
	RequestTxHash *TX_HASH `json:"RequestTxHash,omitempty"`
	// Batch variants of the above
	RequestIDs      []TX_HASH `json:"RequestIDs,omitempty"`
	RequestTxHashes []TX_HASH `json:"RequestTxHashes,omitempty"`
	// Used for the VRFv2 - max link this tx will bill
	// should it get bumped
	MaxLink *string `json:"MaxLink,omitempty"`
	// Used for the VRFv2 - the subscription ID of the
	// requester of the VRF.
	SubID *uint64 `json:"SubId,omitempty"`
	// Used for the VRFv2Plus - the uint256 subscription ID of the
	// requester of the VRF.
	GlobalSubID *string `json:"GlobalSubId,omitempty"`
	// Used for VRFv2Plus - max native token this tx will bill
	// should it get bumped
	MaxEth *string `json:"MaxEth,omitempty"`

	// Used for keepers
	UpkeepID *string `json:"UpkeepID,omitempty"`

	// Used for VRF to know if the txn is a ForceFulfilment txn
	ForceFulfilled          *bool   `json:"ForceFulfilled,omitempty"`
	ForceFulfillmentAttempt *uint64 `json:"ForceFulfillmentAttempt,omitempty"`

	// Used for Keystone Workflows
	WorkflowExecutionID *string `json:"WorkflowExecutionID,omitempty"`

	// Used only for forwarded txs, tracks the original destination address.
	// When this is set, it indicates tx is forwarded through To address.
	FwdrDestAddress *ADDR `json:"ForwarderDestAddress,omitempty"`

	// MessageIDs is used by CCIP for tx to executed messages correlation in logs
	MessageIDs []string `json:"MessageIDs,omitempty"`
	// SeqNumbers is used by CCIP for tx to committed sequence numbers correlation in logs
	SeqNumbers []uint64 `json:"SeqNumbers,omitempty"`
}

type TxAttempt[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	ID    int64
	TxID  int64
	Tx    Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]
	TxFee FEE
	// ChainSpecificFeeLimit on the TxAttempt is always the same as the on-chain encoded value for fee limit
	ChainSpecificFeeLimit   uint64
	SignedRawTx             []byte
	Hash                    TX_HASH
	CreatedAt               time.Time
	BroadcastBeforeBlockNum *int64
	State                   TxAttemptState
	Receipts                []ChainReceipt[TX_HASH, BLOCK_HASH] `json:"-"`
	TxType                  int
	IsPurgeAttempt          bool
}

func (a *TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) String() string {
	return fmt.Sprintf("TxAttempt(ID:%d,TxID:%d,Fee:%s,TxType:%d", a.ID, a.TxID, a.TxFee, a.TxType)
}

type Tx[
	CHAIN_ID types.ID,
	ADDR types.Hashable,
	TX_HASH, BLOCK_HASH types.Hashable,
	SEQ types.Sequence,
	FEE feetypes.Fee,
] struct {
	ID             int64
	IdempotencyKey *string
	Sequence       *SEQ
	FromAddress    ADDR
	ToAddress      ADDR
	EncodedPayload []byte
	Value          big.Int
	// FeeLimit on the Tx is always the conceptual gas limit, which is not
	// necessarily the same as the on-chain encoded value (i.e. Optimism)
	FeeLimit uint64
	Error    null.String
	// BroadcastAt is updated every time an attempt for this tx is re-sent
	// In almost all cases it will be within a second or so of the actual send time.
	BroadcastAt *time.Time
	// InitialBroadcastAt is recorded once, the first ever time this tx is sent
	InitialBroadcastAt *time.Time
	CreatedAt          time.Time
	State              TxState
	TxAttempts         []TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE] `json:"-"`
	// Marshalled TxMeta
	// Used for additional context around transactions which you want to log
	// at send time.
	Meta    *sqlutil.JSON
	Subject uuid.NullUUID
	ChainID CHAIN_ID

	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  clnull.Uint32

	// TransmitChecker defines the check that should be performed before a transaction is submitted on
	// chain.
	TransmitChecker *sqlutil.JSON

	// Marks tx requiring callback
	SignalCallback bool
	// Marks tx callback as signaled
	CallbackCompleted bool
}

func (e *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetError() error {
	if e.Error.Valid {
		return errors.New(e.Error.String)
	}
	return nil
}

// GetID allows Tx to be used as jsonapi.MarshalIdentifier
func (e *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetID() string {
	return fmt.Sprintf("%d", e.ID)
}

// GetMeta returns an Tx's meta in struct form, unmarshalling it from JSON first.
func (e *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetMeta() (*TxMeta[ADDR, TX_HASH], error) {
	if e.Meta == nil {
		return nil, nil
	}
	var m TxMeta[ADDR, TX_HASH]
	if err := json.Unmarshal(*e.Meta, &m); err != nil {
		return nil, fmt.Errorf("unmarshalling meta: %w", err)
	}

	return &m, nil
}

// GetLogger returns a new logger with metadata fields.
func (e *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetLogger(lgr logger.Logger) logger.SugaredLogger {
	lgr = logger.With(lgr,
		"txID", e.ID,
		"sequence", e.Sequence,
		"checker", e.TransmitChecker,
		"feeLimit", e.FeeLimit,
	)

	meta, err := e.GetMeta()
	if err != nil {
		lgr.Errorw("failed to get meta of the transaction", "err", err)
		return logger.Sugared(lgr)
	}

	if meta != nil {
		lgr = logger.With(lgr, "jobID", meta.JobID)

		if meta.RequestTxHash != nil {
			lgr = logger.With(lgr, "requestTxHash", *meta.RequestTxHash)
		}

		if meta.RequestID != nil {
			id := *meta.RequestID
			lgr = logger.With(lgr, "requestID", new(big.Int).SetBytes(id.Bytes()).String())
		}

		if len(meta.RequestIDs) != 0 {
			var ids []string
			for _, id := range meta.RequestIDs {
				ids = append(ids, new(big.Int).SetBytes(id.Bytes()).String())
			}
			lgr = logger.With(lgr, "requestIDs", strings.Join(ids, ","))
		}

		if meta.UpkeepID != nil {
			lgr = logger.With(lgr, "upkeepID", *meta.UpkeepID)
		}

		if meta.SubID != nil {
			lgr = logger.With(lgr, "subID", *meta.SubID)
		}

		if meta.MaxLink != nil {
			lgr = logger.With(lgr, "maxLink", *meta.MaxLink)
		}

		if meta.FwdrDestAddress != nil {
			lgr = logger.With(lgr, "FwdrDestAddress", *meta.FwdrDestAddress)
		}

		if len(meta.MessageIDs) > 0 {
			for _, mid := range meta.MessageIDs {
				lgr = logger.With(lgr, "messageID", mid)
			}
		}

		if len(meta.SeqNumbers) > 0 {
			lgr = logger.With(lgr, "SeqNumbers", meta.SeqNumbers)
		}
	}

	return logger.Sugared(lgr)
}

// GetChecker returns an Tx's transmit checker spec in struct form, unmarshalling it from JSON
// first.
func (e *Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, SEQ, FEE]) GetChecker() (TransmitCheckerSpec[ADDR], error) {
	if e.TransmitChecker == nil {
		return TransmitCheckerSpec[ADDR]{}, nil
	}
	var t TransmitCheckerSpec[ADDR]
	if err := json.Unmarshal(*e.TransmitChecker, &t); err != nil {
		return t, fmt.Errorf("unmarshalling transmit checker: %w", err)
	}

	return t, nil
}

// Provides error classification to external components in a chain agnostic way
// Only exposes the error types that could be set in the transaction error field
type ErrorClassifier interface {
	error
	IsFatal() bool
}
