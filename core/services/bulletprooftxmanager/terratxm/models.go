package terratxm

import "time"

type State string

var (
	Unstarted State = "unstarted"
	Completed State = "completed"
	Errored   State = "errored"
)

type TerraMsg struct {
	ID         int64
	ContractID string
	State      State
	Msg        []byte
	From       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
