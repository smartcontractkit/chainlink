package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	clnull "github.com/smartcontractkit/chainlink-common/pkg/utils/null"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
)

type TxState string

const (
	TxUnstarted   = TxState("unstarted")
	TxUnconfirmed = TxState("unconfirmed")
	TxConfirmed   = TxState("confirmed")

	TxFatalError = TxState("fatal")
	TxFinalized  = TxState("finalized")
)

type Transaction struct {
	ID                uint64
	IdempotencyKey    *string
	ChainID           *big.Int
	Nonce             *uint64
	FromAddress       common.Address
	ToAddress         common.Address
	Value             *big.Int
	Data              []byte
	SpecifiedGasLimit uint64

	CreatedAt       time.Time
	LastBroadcastAt time.Time

	State        TxState
	IsPurgeable  bool
	Attempts     []*Attempt
	AttemptCount uint16 // AttempCount is strictly kept in memory and prevents indefinite retrying
	Meta         *sqlutil.JSON
	Subject      uuid.NullUUID

	// Pipeline variables - if you aren't calling this from chain tx task within
	// the pipeline, you don't need these variables
	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  clnull.Uint32
	SignalCallback    bool
	CallbackCompleted bool
}

//	func (t *Transaction) String() string {
//		return fmt.Sprintf(`{"ID":%d, "IdempotencyKey":%v, "ChainID":%v, "Nonce":%d, "FromAddress":%v, "ToAddress":%v, "Value":%v, `+
//			`"Data":%v, "SpecifiedGasLimit":%d, "CreatedAt":%v, "LastBroadcastAt":%v, "State":%v, "IsPurgeable":%v, "AttemptCount":%d, `+
//			`"Meta":%v, "Subject":%v, "PipelineTaskRunID":%v, "MinConfirmations":%v, "SignalCallback":%v, "CallbackCompleted":%v`,
//			t.ID, *t.IdempotencyKey, t.ChainID, t.Nonce, t.FromAddress, t.ToAddress, t.Value, t.Data, t.SpecifiedGasLimit, t.CreatedAt, t.LastBroadcastAt,
//			t.State, t.IsPurgeable, t.AttemptCount, t.Meta, t.Subject, t.PipelineTaskRunID, t.MinConfirmations, t.SignalCallback, t.CallbackCompleted)
//	}
func (t *Transaction) FindAttemptByHash(attemptHash common.Hash) (*Attempt, error) {
	for _, a := range t.Attempts {
		if a.Hash == attemptHash {
			return a, nil
		}
	}
	return nil, fmt.Errorf("attempt with hash: %v was not found", attemptHash)
}

func (t *Transaction) DeepCopy() *Transaction {
	txCopy := *t
	attemptsCopy := make([]*Attempt, 0, len(t.Attempts))
	for _, attempt := range t.Attempts {
		attemptsCopy = append(attemptsCopy, attempt.DeepCopy())
	}
	txCopy.Attempts = attemptsCopy
	return &txCopy
}

func (t *Transaction) GetMeta() (*TxMeta, error) {
	if t.Meta == nil {
		return nil, nil
	}
	var m TxMeta
	if err := json.Unmarshal(*t.Meta, &m); err != nil {
		return nil, fmt.Errorf("unmarshalling meta: %w", err)
	}
	return &m, nil
}

type Attempt struct {
	ID                uint64
	TxID              uint64
	Hash              common.Hash
	Fee               gas.EvmFee
	GasLimit          uint64
	Type              byte
	SignedTransaction *types.Transaction

	CreatedAt   time.Time
	BroadcastAt time.Time
}

func (a *Attempt) DeepCopy() *Attempt {
	txCopy := *a
	if a.SignedTransaction != nil {
		txCopy.SignedTransaction = a.SignedTransaction.WithoutBlobTxSidecar()
	}
	return &txCopy
}

func (a *Attempt) String() string {
	return fmt.Sprintf(`{"ID":%d, "TxID":%d, "Hash":%v, "Fee":%v, "GasLimit":%d, "Type":%v, "CreatedAt":%v, "BroadcastAt":%v}`,
		a.ID, a.TxID, a.Hash, a.Fee, a.GasLimit, a.Type, a.CreatedAt, a.BroadcastAt)
}

type TxRequest struct {
	IdempotencyKey    *string
	ChainID           *big.Int
	FromAddress       common.Address
	ToAddress         common.Address
	Value             *big.Int
	Data              []byte
	SpecifiedGasLimit uint64

	Meta             *sqlutil.JSON // TODO: *TxMeta after migration
	ForwarderAddress common.Address
	// QueueingTxStrategy QueueingTxStrategy

	// Pipeline variables - if you aren't calling this from chain tx task within
	// the pipeline, you don't need these variables
	PipelineTaskRunID uuid.NullUUID
	MinConfirmations  clnull.Uint32
	SignalCallback    bool
}

type TxMeta struct {
	// Pipeline
	JobID        *int32    `json:"JobID,omitempty"`
	FailOnRevert null.Bool `json:"FailOnRevert,omitempty"`

	// VRF
	RequestID               *common.Hash  `json:"RequestID,omitempty"`
	RequestTxHash           *common.Hash  `json:"RequestTxHash,omitempty"`
	RequestIDs              []common.Hash `json:"RequestIDs,omitempty"`
	RequestTxHashes         []common.Hash `json:"RequestTxHashes,omitempty"`
	MaxLink                 *string       `json:"MaxLink,omitempty"`
	SubID                   *uint64       `json:"SubId,omitempty"`
	GlobalSubID             *string       `json:"GlobalSubId,omitempty"`
	MaxEth                  *string       `json:"MaxEth,omitempty"`
	ForceFulfilled          *bool         `json:"ForceFulfilled,omitempty"`
	ForceFulfillmentAttempt *uint64       `json:"ForceFulfillmentAttempt,omitempty"`

	// Used for keepers
	UpkeepID *string `json:"UpkeepID,omitempty"`

	// Used for Keystone Workflows
	WorkflowExecutionID *string `json:"WorkflowExecutionID,omitempty"`

	// Forwarders
	FwdrDestAddress *common.Address `json:"ForwarderDestAddress,omitempty"`

	// CCIP
	MessageIDs []string `json:"MessageIDs,omitempty"`
	SeqNumbers []uint64 `json:"SeqNumbers,omitempty"`

	// Dual Broadcast
	DualBroadcast       *bool   `json:"DualBroadcast,omitempty"`
	DualBroadcastParams *string `json:"DualBroadcastParams,omitempty"`
}

type QueueingTxStrategy struct {
	QueueSize uint32
	Subject   uuid.NullUUID
}
