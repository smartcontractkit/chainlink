package terratxm

import (
	"time"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/client"

	"github.com/smartcontractkit/terra.go/msg"
)

// State represents the state of a given terra msg
// Happy path: Unstarted->Broadcasted->Confirmed
type State string

var (
	// Unstarted means queued but not processed.
	// Valid next states: Broadcasted, Errored (sim fails)
	Unstarted State = "unstarted"
	// Broadcasted means included in the mempool of a node.
	// Valid next states: Confirmed (found onchain), Errored (tx expired waiting for confirmation)
	Broadcasted State = "broadcasted"
	// Confirmed means we're able to retrieve the txhash of the tx which broadcasted the msg.
	// Valid next states: none, terminal state
	Confirmed State = "confirmed"
	// Errored means the msg reverted in simulation OR the tx containing the message timed out waiting to be confirmed
	// TODO: when we add gas bumping, we'll address that timeout case
	// Valid next states, none, terminal state
	Errored State = "errored"
)

// TerraMsg a terra msg
type TerraMsg struct {
	ID         int64
	ChainID    string `db:"terra_chain_id"`
	ContractID string
	State      State
	Msg        []byte
	TxHash     *string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// In memory only
	ExecuteContract *msg.ExecuteContract
}

func getSimMsgs(tms []TerraMsg) []client.SimMsg {
	var msgs []client.SimMsg
	for i := range tms {
		msgs = append(msgs, client.SimMsg{
			ID:  tms[i].ID,
			Msg: tms[i].ExecuteContract,
		})
	}
	return msgs
}

func getSimMsgsIDs(msgs []client.SimMsg) []int64 {
	var ids []int64
	for i := range msgs {
		ids = append(ids, msgs[i].ID)
	}
	return ids
}

func getMsgs(simMsgs []client.SimMsg) []msg.Msg {
	var msgs []msg.Msg
	for i := range simMsgs {
		msgs = append(msgs, simMsgs[i].Msg)
	}
	return msgs
}

func getIDs(tms []TerraMsg) []int64 {
	var ids []int64
	for i := range tms {
		ids = append(ids, tms[i].ID)
	}
	return ids
}
