package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

const (
	// StatusInProgress is used for when a run is actively being executed.
	StatusInProgress = "in progress"
	// StatusPending is used for when a run is waiting on the completion
	// of another event.
	StatusPending = "pending"
	// StatusErrored is used for when a run has errored and will not complete.
	StatusErrored = "errored"
	// StatusCompleted is used for when a run has successfully completed execution.
	StatusCompleted = "completed"
)

// JobSpec is the definition for all the work to be carried out by the node
// for a given contract. It contains the Initiators, Tasks (which are the
// individual steps to be carried out), StartAt, EndAt, and CreatedAt fields.
type JobSpec struct {
	ID         string      `json:"id" storm:"id,unique"`
	Initiators []Initiator `json:"initiators"`
	Tasks      []TaskSpec  `json:"tasks" storm:"inline"`
	StartAt    null.Time   `json:"startAt" storm:"index"`
	EndAt      null.Time   `json:"endAt" storm:"index"`
	CreatedAt  Time        `json:"createdAt" storm:"index"`
}

// NewJob initializes a new job by generating a unique ID and setting
// the CreatedAt field to the time of invokation.
func NewJob() JobSpec {
	return JobSpec{
		ID:        utils.NewBytes32ID(),
		CreatedAt: Time{Time: time.Now()},
	}
}

// NewRun initializes the job by creating the IDs for the job
// and all associated tasks, and setting the CreatedAt field.
func (j JobSpec) NewRun() JobRun {
	jrid := utils.NewBytes32ID()
	taskRuns := make([]TaskRun, len(j.Tasks))
	for i, task := range j.Tasks {
		taskRuns[i] = TaskRun{
			ID:     utils.NewBytes32ID(),
			Task:   task,
			Result: RunResult{JobRunID: jrid},
		}
	}

	return JobRun{
		ID:        jrid,
		JobID:     j.ID,
		CreatedAt: time.Now(),
		TaskRuns:  taskRuns,
	}
}

// InitiatorsFor returns an array of Initiators for the given list of
// Initiator types.
func (j JobSpec) InitiatorsFor(types ...string) []Initiator {
	list := []Initiator{}
	for _, initr := range j.Initiators {
		for _, t := range types {
			if initr.Type == t {
				list = append(list, initr)
			}
		}
	}
	return list
}

// WebAuthorized returns true if the "web" initiator is present.
func (j JobSpec) WebAuthorized() bool {
	for _, initr := range j.Initiators {
		if initr.Type == InitiatorWeb {
			return true
		}
	}
	return false
}

// Returns true if any of the job's initiators are triggered by event logs.
func (j JobSpec) IsLogInitiated() bool {
	for _, initr := range j.Initiators {
		if initr.IsLogInitiated() {
			return true
		}
	}
	return false
}

// Ended returns true if the job has ended.
func (j JobSpec) Ended(t time.Time) bool {
	if !j.EndAt.Valid {
		return false
	}
	return t.After(j.EndAt.Time)
}

// Started returns true if the job has started.
func (j JobSpec) Started(t time.Time) bool {
	if !j.StartAt.Valid {
		return true
	}
	return t.After(j.StartAt.Time) || t.Equal(j.StartAt.Time)
}

const (
	// InitiatorRunLog for tasks in a job to watch an ethereum address
	// and expect a JSON payload from a log event.
	InitiatorRunLog = "runlog"
	// InitiatorCron for tasks in a job to be ran on a schedule.
	InitiatorCron = "cron"
	// InitiatorEthLog for tasks in a job to use the Ethereum blockchain.
	InitiatorEthLog = "ethlog"
	// InitiatorRunAt for tasks in a job to be ran once.
	InitiatorRunAt = "runat"
	// InitiatorWeb for tasks in a job making a web request.
	InitiatorWeb = "web"
)

// Initiator could be though of as a trigger, define how a Job can be
// started, or rather, how a JobRun can be created from a Job.
// Initiators will have their own unique ID, but will be assocated
// to a parent JobID.
type Initiator struct {
	ID       int            `json:"id" storm:"id,increment"`
	JobID    string         `json:"jobId" storm:"index"`
	Type     string         `json:"type" storm:"index"`
	Schedule Cron           `json:"schedule,omitempty"`
	Time     Time           `json:"time,omitempty"`
	Ran      bool           `json:"ran,omitempty"`
	Address  common.Address `json:"address,omitempty" storm:"index"`
}

// UnmarshalJSON parses the raw initiator data and updates the
// initiator as long as the type is valid.
func (i *Initiator) UnmarshalJSON(input []byte) error {
	type Alias Initiator
	var aux Alias
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}

	*i = Initiator(aux)
	i.Type = strings.ToLower(aux.Type)
	return nil
}

// Returns true if triggered by event logs.
func (i Initiator) IsLogInitiated() bool {
	return i.Type == InitiatorEthLog || i.Type == InitiatorRunLog
}

// TaskSpec is the definition of work to be carried out. The
// Type will be an adapter, and the Params will contain any
// additional information that adapter would need to operate.
type TaskSpec struct {
	Type   string `json:"type" storm:"index"`
	Params JSON
}

// UnmarshalJSON parses the given input and updates the TaskSpec.
func (t *TaskSpec) UnmarshalJSON(input []byte) error {
	type Alias TaskSpec
	var aux Alias
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}

	t.Type = strings.ToLower(aux.Type)
	var params json.RawMessage
	if err := json.Unmarshal(input, &params); err != nil {
		return err
	}

	t.Params = JSON{gjson.ParseBytes(params)}
	return nil
}

// MarshalJSON returns the JSON-encoded TaskSpec Params.
func (t TaskSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Params)
}

// BridgeType is used for external adapters and has fields for
// the name of the adapter and its URL.
type BridgeType struct {
	Name string `json:"name" storm:"id,unique"`
	URL  WebURL `json:"url"`
}

// UnmarshalJSON parses the given input and updates the BridgeType
// Name and URL.
func (bt *BridgeType) UnmarshalJSON(input []byte) error {
	type Alias BridgeType
	var aux Alias
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}
	bt.Name = strings.ToLower(aux.Name)
	bt.URL = aux.URL
	return nil
}
