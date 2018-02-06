package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

const (
	StatusInProgress = "in progress"
	StatusPending    = "pending"
	StatusErrored    = "errored"
	StatusCompleted  = "completed"
)

type Job struct {
	ID         string      `json:"id" storm:"id,index,unique"`
	Initiators []Initiator `json:"initiators"`
	Tasks      []Task      `json:"tasks" storm:"inline"`
	StartAt    null.Time   `json:"startAt" storm:"index"`
	EndAt      null.Time   `json:"endAt" storm:"index"`
	CreatedAt  Time        `json:"createdAt" storm:"index"`
}

func NewJob() *Job {
	return &Job{
		ID:        utils.NewBytes32ID(),
		CreatedAt: Time{Time: time.Now()},
	}
}

func (j *Job) NewRun() *JobRun {
	taskRuns := make([]TaskRun, len(j.Tasks))
	for i, task := range j.Tasks {
		taskRuns[i] = TaskRun{
			ID:   utils.NewBytes32ID(),
			Task: task,
		}
	}

	return &JobRun{
		ID:        utils.NewBytes32ID(),
		JobID:     j.ID,
		CreatedAt: time.Now(),
		TaskRuns:  taskRuns,
	}
}

func (j *Job) InitiatorsFor(types ...string) []Initiator {
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

func (j *Job) WebAuthorized() bool {
	for _, initr := range j.Initiators {
		if initr.Type == InitiatorWeb {
			return true
		}
	}
	return false
}

func (j *Job) Ended(t time.Time) bool {
	if !j.EndAt.Valid {
		return false
	}
	return t.After(j.EndAt.Time)
}

func (j *Job) Started(t time.Time) bool {
	if !j.StartAt.Valid {
		return true
	}
	return t.After(j.StartAt.Time) || t.Equal(j.StartAt.Time)
}

const (
	InitiatorChainlinkLog = "chainlinklog"
	InitiatorCron         = "cron"
	InitiatorEthLog       = "ethlog"
	InitiatorRunAt        = "runat"
	InitiatorWeb          = "web"
)

var initiatorWhitelist = map[string]bool{
	InitiatorChainlinkLog: true,
	InitiatorCron:         true,
	InitiatorEthLog:       true,
	InitiatorRunAt:        true,
	InitiatorWeb:          true,
}

type Initiator struct {
	ID       int            `json:"id" storm:"id,increment"`
	JobID    string         `json:"jobId" storm:"index"`
	Type     string         `json:"type" storm:"index"`
	Schedule Cron           `json:"schedule,omitempty"`
	Time     Time           `json:"time,omitempty"`
	Ran      bool           `json:"ran,omitempty"`
	Address  common.Address `json:"address,omitempty" storm:"index"`
}

func (i *Initiator) UnmarshalJSON(input []byte) error {
	type Alias Initiator
	var aux Alias
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}

	*i = Initiator(aux)
	i.Type = strings.ToLower(aux.Type)
	if _, valid := initiatorWhitelist[i.Type]; !valid {
		return fmt.Errorf("Initiator %v does not exist", aux.Type)
	}
	return nil
}

type Task struct {
	Type   string `json:"type" storm:"index"`
	Params JSON
}

func (t *Task) UnmarshalJSON(input []byte) error {
	type Alias Task
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

func (t Task) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Params)
}

type BridgeType struct {
	Name string `json:"name" storm:"id,index,unique"`
	URL  WebURL `json:"url"`
}

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
