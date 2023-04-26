package types

import (
	"math/big"

	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// TxStrategy controls how txes are queued and sent
//
//go:generate mockery --quiet --name TxStrategy --output ./mocks/ --case=underscore --structname TxStrategy --filename tx_strategy.go
type TxStrategy interface {
	// Subject will be saved txes.subject if not null
	Subject() uuid.NullUUID
	// PruneQueue is called after tx insertion
	// It accepts the service responsible for deleting
	// unstarted txs and deletion options
	PruneQueue(pruneService UnstartedTxQueuePruner, qopt pg.QOpt) (n int64, err error)
}

type TxAttemptState string

const (
	TxAttemptInProgress = TxAttemptState("in_progress")
	// TODO: Make name chain-agnostic (https://smartcontract-it.atlassian.net/browse/BCI-981)
	TxAttemptInsufficientEth = TxAttemptState("insufficient_eth")
	TxAttemptBroadcast       = TxAttemptState("broadcast")
)

// Transaction is the type that callers get back, when they create a Transaction using the Txm.
// TODO: Remove this with the EthTx type, once that is extracted out to this namespace.
type Transaction interface {
	GetID() string
}

type NewTx[ADDR types.Hashable, TX_HASH types.Hashable] struct {
	FromAddress      ADDR
	ToAddress        ADDR
	EncodedPayload   []byte
	FeeLimit         uint32
	Meta             *TxMeta[ADDR, TX_HASH]
	ForwarderAddress ADDR

	// Pipeline variables - if you aren't calling this from ethtx task within
	// the pipeline, you don't need these variables
	MinConfirmations  clnull.Uint32
	PipelineTaskRunID *uuid.UUID

	Strategy TxStrategy

	// Checker defines the check that should be run before a transaction is submitted on chain.
	Checker TransmitCheckerSpec[ADDR]
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

	// Used for keepers
	UpkeepID *string `json:"UpkeepID,omitempty"`

	// Used only for forwarded txs, tracks the original destination address.
	// When this is set, it indicates tx is forwarded through To address.
	FwdrDestAddress *ADDR `json:"ForwarderDestAddress,omitempty"`

	// MessageIDs is used by CCIP for tx to executed messages correlation in logs
	MessageIDs []string `json:"MessageIDs,omitempty"`
	// SeqNumbers is used by CCIP for tx to committed sequence numbers correlation in logs
	SeqNumbers []uint64 `json:"SeqNumbers,omitempty"`
}
