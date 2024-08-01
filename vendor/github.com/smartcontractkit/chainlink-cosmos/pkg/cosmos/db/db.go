package db

import (
	"time"
)

type Node struct {
	ID            int32
	Name          string
	CosmosChainID string
	TendermintURL string `db:"tendermint_url"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// State represents the state of a given cosmos msg
// Happy path: Unstarted->Broadcasted->Confirmed
type State string

var (
	// Unstarted means queued but not processed.
	// Valid next states: Started, Errored (cancelled)
	Unstarted State = "unstarted"
	// Started means included in a batch about to be broadcast.
	// Valid next states: Broadcasted, Errored (sim fails)
	Started State = "started"
	// Broadcasted means included in the mempool of a node.
	// Valid next states: Confirmed (found onchain), Errored (tx expired waiting for confirmation)
	Broadcasted State = "broadcasted"
	// Confirmed means we're able to retrieve the txhash of the tx which broadcasted the msg.
	// Valid next states: none, terminal state
	Confirmed State = "confirmed"
	// Errored means the msg:
	//  - reverted in simulation
	//  - the tx containing the message timed out waiting to be confirmed
	//  - the msg was cancelled
	// TODO: when we add gas bumping, we'll address that timeout case
	// Valid next states, none, terminal state
	Errored State = "errored"
)

type Msg struct {
	ID         int64
	ChainID    string `db:"cosmos_chain_id"`
	ContractID string
	State      State
	Type       string // cosmos-sdk/types.MsgTypeURL()
	Raw        []byte // proto.Marshal()
	TxHash     *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
