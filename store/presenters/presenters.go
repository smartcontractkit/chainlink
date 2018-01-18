package presenters

import "github.com/smartcontractkit/chainlink/store/models"

type Job struct {
	*models.Job
	Runs []*models.JobRun `json:"runs,omitempty"`
}
