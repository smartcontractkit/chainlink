package presenters

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
)

type Job struct {
	*models.Job
	Runs []*models.JobRun `json:"runs,omitempty"`
}

func (job Job) FriendlyCreatedAt() string {
	return job.CreatedAt.HumanString()
}

func (job Job) FriendlyEndAt() string {
	if job.EndAt.Valid {
		return utils.HumanTimeString(job.EndAt.Time)
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

func (i Initiator) FriendlyRunAt() string {
	if i.Type == "cron" {
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
