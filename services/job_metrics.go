package services

import (
	"expvar"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"go.uber.org/multierr"
)

// JobMetrics safely stores all job metric data in memory
type JobMetrics interface {
	Start() error
	Add(models.JobSpec) error
	AddRun(models.JobRun)
	Get(string) models.JobSpecMetrics
}

type jobMetrics struct {
	collection map[string]models.JobSpecMetrics
	eMap       *expvar.Map
	store      *store.Store
}

// NewJobMetrics initializes a JobMetrics.
func NewJobMetrics(str *store.Store) JobMetrics {
	n := "jobmetrics"
	if expvar.Get(n) == nil {
		expvar.NewMap(n)
	}
	return &jobMetrics{
		collection: make(map[string]models.JobSpecMetrics),
		eMap:       expvar.Get(n).(*expvar.Map),
		store:      str,
	}
}

// Start fetches all the current jobs and their runs from the store
// then adds them to the collection
func (jm *jobMetrics) Start() error {
	var merr error

	aj, err := jm.store.Jobs()
	if err != nil {
		return err
	}
	for _, j := range aj {
		if err := jm.Add(j); err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	return merr
}

// Add allows a new job to be added to the metrics collection
func (jm *jobMetrics) Add(job models.JobSpec) error {
	jrs, err := jm.store.JobRunsFor(job.ID)
	if err != nil {
		return err
	}

	ac := make(map[models.TaskType]int)
	for _, t := range job.Tasks {
		ac[t.Type]++
	}
	jm.collection[job.ID] = models.JobSpecMetrics{
		ID:           job.ID,
		RunCount:     len(jrs),
		AdaptorCount: ac,
		StatusCount:  make(map[models.RunStatus]int),
		ParamCount:   make(map[string][]models.ParamCount),
	}

	for _, jr := range jrs {
		jm.AddRun(jr)
	}

	jm.eMap.Set(job.ID, jm.collection[job.ID])
	return nil
}

// AddRun adds a new job run to the metrics collection
func (jm *jobMetrics) AddRun(run models.JobRun) {
	e := jm.collection[run.JobID]
	sc := e.StatusCount
	pc := e.ParamCount

	sc[run.Status]++
	if len(run.TaskRuns) > 0 {
		countParams(jm.store, run, pc)
	}

	jm.collection[run.JobID] = models.JobSpecMetrics{
		ID:           e.ID,
		RunCount:     e.RunCount,
		AdaptorCount: e.AdaptorCount,
		StatusCount:  sc,
		ParamCount:   pc,
	}
}

// Get returns a single JobSpecMetrics from a given job id
func (jm *jobMetrics) Get(jobID string) models.JobSpecMetrics {
	return jm.collection[jobID]
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
