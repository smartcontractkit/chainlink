package services

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// AllJobSpecMetrics returns all Job Spec Stat data for the Job Spec inputs
func AllJobSpecMetrics(store *store.Store, jobs []models.JobSpec) ([]models.JobSpecMetrics, error) {
	var merr error

	jsm := []models.JobSpecMetrics{}
	for _, j := range jobs {
		if jobMetrics, err := jobSpecMetrics(store, j); err != nil {
			merr = multierr.Append(merr, err)
		} else {
			jsm = append(jsm, jobMetrics)
		}
	}
	return jsm, merr
}

func jobSpecMetrics(store *store.Store, job models.JobSpec) (models.JobSpecMetrics, error) {
	jrs, err := store.JobRunsFor(job.ID)
	if err != nil {
		return models.JobSpecMetrics{}, err
	}

	rc := make(map[models.RunStatus]int)
	ac := make(map[models.TaskType]int)
	pc := make(map[string][]models.ParamCount)
	for _, jr := range jrs {
		rc[jr.Status]++
		if len(jr.TaskRuns) > 0 {
			countParams(store, jr, pc)
		}
	}
	for _, t := range job.Tasks {
		ac[t.Type]++
	}
	return models.JobSpecMetrics{
		ID:           job.ID,
		RunCount:     len(jrs),
		AdaptorCount: ac,
		StatusCount:  rc,
		ParamCount:   pc,
	}, nil
}

func countParams(store *store.Store, jobRun models.JobRun, paramCount map[string][]models.ParamCount) {
	tp := store.Config.MetricsParam

	for _, p := range tp {
		trp := jobRun.TaskRuns[0].Task.Params.Get(p)
		if trp.Exists() {
			found := false
			for i := range paramCount[p] {
				if paramCount[p][i].Value == trp.String() {
					found = true
					paramCount[p][i].Count++
					break
				}
			}

			if !found {
				paramCount[p] = append(paramCount[p], models.ParamCount{
					Value: trp.String(),
					Count: 1,
				})
			}
		}
	}
}
