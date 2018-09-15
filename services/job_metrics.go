package services

import (
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
	"expvar"
)

type JobMetrics interface {
	Start()                  error
	Add(models.JobSpec)      error
	AddRun(models.JobRun)
}

type jobMetrics struct {
	collection    map[string]models.JobSpecMetrics
	eMap          *expvar.Map
	store         *store.Store
}

func NewJobMetrics(str *store.Store) *jobMetrics {
	return &jobMetrics{
		collection: make(map[string]models.JobSpecMetrics),
		eMap: expvar.NewMap("jobmetrics"),
		store: str,
	}
}

func (jsm *jobMetrics) Start() error {
	var merr error

	aj, err := jsm.store.Jobs()
	if err != nil {
		return err
	}
	for _, j := range aj {
		if err := jsm.Add(j); err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	return merr
}

func (jsm *jobMetrics) Add(job models.JobSpec) error {
	jrs, err := jsm.store.JobRunsFor(job.ID)
	if err != nil {
		return err
	}

	ac := make(map[models.TaskType]int)
	for _, t := range job.Tasks {
		ac[t.Type]++
	}
	jsm.collection[job.ID] = models.JobSpecMetrics{
		ID:           job.ID,
		RunCount:     len(jrs),
		AdaptorCount: ac,
		StatusCount:  make(map[models.RunStatus]int),
		ParamCount:   make(map[string][]models.ParamCount),
	}

	for _, jr := range jrs {
		jsm.AddRun(jr)
	}

	jsm.eMap.Set(job.ID, jsm.collection[job.ID])
	return nil
}

func (jsm *jobMetrics) AddRun(run models.JobRun) {
	e := jsm.collection[run.JobID]
	sc := e.StatusCount
	pc := e.ParamCount

	sc[run.Status]++
	if len(run.TaskRuns) > 0 {
		countParams(jsm.store, run, pc)
	}

	jsm.collection[run.JobID] = models.JobSpecMetrics{
		ID: e.ID,
		RunCount: e.RunCount,
		AdaptorCount: e.AdaptorCount,
		StatusCount: sc,
		ParamCount: pc,
	}
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
