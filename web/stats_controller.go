package web

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"go.uber.org/multierr"
)

// StatsController fetches high-level information about the nodes usage
type StatsController struct {
	App *services.ChainlinkApplication
}

// Show displays Stats, one page at a time
// Example:
//  "<application>/stats?size=1&page=2"
func (sc *StatsController) Show(c *gin.Context) {
	store := sc.App.Store
	jobs := []models.JobSpec{}

	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	skip := storm.Skip(offset)
	limit := storm.Limit(size)

	if count, err := sc.App.Store.Count(&models.JobSpec{}); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting count of JobSpec: %+v", err))
	} else if err := sc.App.Store.AllByIndex("CreatedAt", &jobs, skip, limit); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting Jobs: %+v", err))
	} else if jss, err := sc.allJobStats(jobs); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting job stats: %+v", err))
	} else if account, err := store.KeyStore.GetAccount(); err != nil {
		publicError(c, 400, err)
	} else {
		a := presenters.AccountBalance{
			Address: account.Address.Hex(),
		}
		s := presenters.Stats{
			Account:      a,
			JobSpecStats: jss,
		}

		buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, s)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("failed to marshal document: %+v", err))
		} else {
			c.Data(200, MediaType, buffer)
		}
	}
}

func (sc *StatsController) allJobStats(jobs []models.JobSpec) ([]presenters.JobSpecStats, error) {
	js := []presenters.JobSpecStats{}

	var merr error
	for _, j := range jobs {
		if jobStats, err := sc.jobStats(j); err != nil {
			merr = multierr.Append(merr, err)
		} else {
			js = append(js, jobStats)
		}
	}
	return js, merr
}

func (sc *StatsController) jobStats(job models.JobSpec) (presenters.JobSpecStats, error) {
	jrs, err := sc.App.Store.JobRunsFor(job.ID)
	if err != nil {
		return presenters.JobSpecStats{}, err
	}

	rc := make(map[models.RunStatus]int)
	uc := make(map[string]int)
	ac := make(map[models.TaskType]int)
	for _, jr := range jrs {
		rc[jr.Status]++
		if len(jr.TaskRuns) > 0 && jr.TaskRuns[0].Task.Params.Get("url").Exists() {
			uc[jr.TaskRuns[0].Task.Params.Get("url").Str]++
		}
	}
	for _, t := range job.Tasks {
		ac[t.Type]++
	}
	return presenters.JobSpecStats{
		ID:           job.ID,
		RunCount:     len(jrs),
		AdaptorCount: ac,
		StatusCount:  rc,
		URLCount:     uc,
	}, nil
}
