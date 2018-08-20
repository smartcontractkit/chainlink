package services

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// AllJobSpecStats returns all Job Spec Stat data for the Job Spec inputs
func AllJobSpecStats(store *store.Store, jobs []models.JobSpec) (models.JobSpecStats, error) {
	var merr error

	s := models.JobSpecStats{}
	if a, err := store.KeyStore.GetAccount(); err != nil {
		merr = multierr.Append(merr, err)
	} else {
		s.Address = a.Address.Hex()
	}

	for _, j := range jobs {
		if jobStats, err := jobSpecCounts(store, j); err != nil {
			merr = multierr.Append(merr, err)
		} else {
			s.JobSpecCounts = append(s.JobSpecCounts, jobStats)
		}
	}
	return s, merr
}

func jobSpecCounts(store *store.Store, job models.JobSpec) (models.JobSpecCounts, error) {
	jrs, err := store.JobRunsFor(job.ID)
	if err != nil {
		return models.JobSpecCounts{}, err
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
	return models.JobSpecCounts{
		ID:           job.ID,
		RunCount:     len(jrs),
		AdaptorCount: ac,
		StatusCount:  rc,
		ParamCount:   pc,
	}, nil
}

func countParams(store *store.Store, jobRun models.JobRun, paramCount map[string][]models.ParamCount) {
	tp := store.Config.StatsParam

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
