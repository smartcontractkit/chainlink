package terratxm

import (
	"time"

	"github.com/smartcontractkit/terra.go/msg"
)

type State string

var (
	Unstarted   State = "unstarted"
	Broadcasted State = "broadcasted"
	Confirmed   State = "confirmed"
	Errored     State = "errored"
)

type TerraMsg struct {
	ID         int64
	ContractID string
	State      State
	Msg        []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// In memory only
	ExecuteContract *msg.ExecuteContract
}

func GetMsgs(tms []TerraMsg) []msg.Msg {
	var msgs []msg.Msg
	for i := range tms {
		msgs = append(msgs, tms[i].ExecuteContract)
	}
	return msgs
}

func GetIDs(tms []TerraMsg) []int64 {
	var ids []int64
	for i := range tms {
		ids = append(ids, tms[i].ID)
	}
	return ids
}
