// Package presenters allow for the specification and result
// of a Job, its associated Tasks, and every JobRun and TaskRun
// to be returned in a consistent manner, typically as a string.
package presenters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

// Job holds the Job definition and each run associated with that Job.
type Job struct {
	models.Job
	Runs []models.JobRun `json:"runs,omitempty"`
}

// MarshalJSON returns the JSON data of the Job and its Initiators.
func (j Job) MarshalJSON() ([]byte, error) {
	type Alias Job
	pis := make([]Initiator, len(j.Initiators))
	for i, modelInitr := range j.Initiators {
		pis[i] = Initiator{modelInitr}
	}
	return json.Marshal(&struct {
		Initiators []Initiator `json:"initiators"`
		Alias
	}{
		pis,
		Alias(j),
	})
}

// FriendlyCreatedAt returns a human-readable string of the Job's
// CreatedAt field.
func (job Job) FriendlyCreatedAt() string {
	return job.CreatedAt.HumanString()
}

// FriendlyEndAt returns a human-readable string of the Job's
// EndAt field.
func (job Job) FriendlyEndAt() string {
	if job.EndAt.Valid {
		return utils.ISO8601UTC(job.EndAt.Time)
	}
	return ""
}

// FriendlyInitiators returns the list of Initiator types as
// a comma separated string.
func (job Job) FriendlyInitiators() string {
	var initrs []string
	for _, i := range job.Initiators {
		initrs = append(initrs, i.Type)
	}
	return strings.Join(initrs, ",")
}

// FriendlyTasks returns the list of Task types as a comma
// separated string.
func (job Job) FriendlyTasks() string {
	var tasks []string
	for _, t := range job.Tasks {
		tasks = append(tasks, t.Type)
	}

	return strings.Join(tasks, ",")
}

// Initiator holds the Job definition's Initiator.
type Initiator struct {
	models.Initiator
}

// MarshalJSON returns the JSON data of the Initiator based
// on its Initiator Type.
func (i Initiator) MarshalJSON() ([]byte, error) {
	switch i.Type {
	case models.InitiatorWeb:
		return json.Marshal(&struct {
			Type string `json:"type"`
		}{
			models.InitiatorWeb,
		})
	case models.InitiatorCron:
		return json.Marshal(&struct {
			Type     string      `json:"type"`
			Schedule models.Cron `json:"schedule"`
		}{
			models.InitiatorCron,
			i.Schedule,
		})
	case models.InitiatorRunAt:
		return json.Marshal(&struct {
			Type string      `json:"type"`
			Time models.Time `json:"time"`
			Ran  bool        `json:"ran"`
		}{
			models.InitiatorRunAt,
			i.Time,
			i.Ran,
		})
	case models.InitiatorEthLog:
		return json.Marshal(&struct {
			Type    string         `json:"type"`
			Address common.Address `json:"address"`
		}{
			models.InitiatorEthLog,
			i.Address,
		})
	default:
		return nil, fmt.Errorf("Cannot marshal unsupported initiator type %v", i.Type)
	}
}

// FriendlyRunAt returns a human-readable string for Cron Initiator types.
func (i Initiator) FriendlyRunAt() string {
	if i.Type == models.InitiatorCron {
		return i.Time.HumanString()
	}
	return ""
}

var empty_address = common.Address{}.String()

// FriendlyAddress returns the Ethereum address if present, and a blank
// string if not.
func (i Initiator) FriendlyAddress() string {
	if empty_address == i.Address.String() {
		return ""
	}
	return i.Address.String()
}

// Task holds a task specified in the Job definition.
type Task struct {
	models.Task
}

// FriendlyParams returns a string of the parameters specified
// for the Task.
func (t Task) FriendlyParams() (string, error) {
	j, err := json.Marshal(&t.Params)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
