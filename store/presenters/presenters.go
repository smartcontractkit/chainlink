package presenters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

type Job struct {
	models.Job
	Runs []models.JobRun `json:"runs,omitempty"`
}

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

func (job Job) FriendlyCreatedAt() string {
	return job.CreatedAt.HumanString()
}

func (job Job) FriendlyEndAt() string {
	if job.EndAt.Valid {
		return utils.ISO8601UTC(job.EndAt.Time)
	}
	return ""
}

func (job Job) FriendlyInitiators() string {
	var initrs []string
	for _, i := range job.Initiators {
		initrs = append(initrs, i.Type)
	}
	return strings.Join(initrs, ",")
}

func (job Job) FriendlyTasks() string {
	var tasks []string
	for _, t := range job.Tasks {
		tasks = append(tasks, t.Type)
	}

	return strings.Join(tasks, ",")
}

type Initiator struct {
	models.Initiator
}

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
	case models.InitiatorChainlinkLog:
		return json.Marshal(&struct {
			Type    string         `json:"type"`
			Address common.Address `json:"address"`
		}{
			models.InitiatorChainlinkLog,
			i.Address,
		})
	default:
		return nil, fmt.Errorf("Cannot marshal unsupported initiator type %v", i.Type)
	}
}

func (i Initiator) FriendlyRunAt() string {
	if i.Type == models.InitiatorCron {
		return i.Time.HumanString()
	}
	return ""
}

var empty_address = common.Address{}.String()

func (i Initiator) FriendlyAddress() string {
	if empty_address == i.Address.String() {
		return ""
	}
	return i.Address.String()
}

type Task struct {
	models.Task
}

func (t Task) FriendlyParams() (string, error) {
	j, err := json.Marshal(&t.Params)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
