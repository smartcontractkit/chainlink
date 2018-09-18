package old

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/store/migrations/migration0"
	"github.com/tidwall/gjson"
)

type JobSpec struct {
	ID         migration0.Unchanged `json:"id" storm:"id,unique"`
	CreatedAt  migration0.Unchanged `json:"createdAt" storm:"index"`
	Initiators migration0.Unchanged `json:"initiators"`
	Tasks      []TaskSpec           `json:"tasks" storm:"inline"`
	StartAt    migration0.Unchanged `json:"startAt" storm:"index"`
	EndAt      migration0.Unchanged `json:"endAt" storm:"index"`
}

type TaskSpec struct {
	Type          migration0.Unchanged `json:"type" storm:"index"`
	Confirmations migration0.Unchanged `json:"confirmations"`
	Params        migration0.JSON      `json:"-"`
}

type JobRun struct {
	ID             migration0.Unchanged `json:"id" storm:"id,unique"`
	JobID          migration0.Unchanged `json:"jobId" storm:"index"`
	Result         migration0.Unchanged `json:"result" storm:"inline"`
	Status         migration0.Unchanged `json:"status" storm:"index"`
	TaskRuns       []TaskRun            `json:"taskRuns" storm:"inline"`
	CreatedAt      migration0.Unchanged `json:"createdAt" storm:"index"`
	CompletedAt    migration0.Unchanged `json:"completedAt"`
	Initiator      migration0.Unchanged `json:"initiator"`
	CreationHeight migration0.Unchanged `json:"creationHeight"`
	Overrides      migration0.Unchanged `json:"overrides"`
}

type TaskRun struct {
	ID     migration0.Unchanged `json:"id" storm:"id,unique"`
	Result migration0.Unchanged `json:"result"`
	Status migration0.Unchanged `json:"status"`
	Task   TaskSpec             `json:"task"`
}

func (t *TaskSpec) UnmarshalJSON(input []byte) error {
	type Alias TaskSpec
	var aux Alias
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}

	t.Confirmations = aux.Confirmations
	t.Type = aux.Type
	var params json.RawMessage
	if err := json.Unmarshal(input, &params); err != nil {
		return err
	}

	t.Params = migration0.JSON{gjson.ParseBytes(params)}
	return nil
}

// MarshalJSON returns the JSON-encoded TaskSpec Params.
func (t TaskSpec) MarshalJSON() ([]byte, error) {
	type Alias TaskSpec
	var aux Alias
	aux = Alias(t)
	b, err := json.Marshal(aux)
	if err != nil {
		return b, err
	}

	js := gjson.ParseBytes(b)
	merged, err := t.Params.Merge(migration0.JSON{js})
	if err != nil {
		return nil, err
	}
	return json.Marshal(merged)
}
