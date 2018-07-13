package web

import (
	"errors"
	"fmt"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// JobSpecsController manages JobSpec requests.
type JobSpecsController struct {
	App *services.ChainlinkApplication
}

// Index lists JobSpecs, one page at a time.
// Example:
//  "<application>/specs?size=1&page=2"
func (jsc *JobSpecsController) Index(c *gin.Context) {
	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		publicError(c, 422, err)
		return
	}

	skip := storm.Skip(offset)
	limit := storm.Limit(size)

	var jobs []models.JobSpec
	if count, err := jsc.App.Store.Count(&models.JobSpec{}); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting count of JobSpec: %+v", err))
	} else if err := jsc.App.Store.AllByIndex("CreatedAt", &jobs, skip, limit); err != nil {
		c.AbortWithError(500, fmt.Errorf("erorr fetching All JobSpecs: %+v", err))
	} else {
		pjs := make([]presenters.JobSpec, len(jobs))
		for i, j := range jobs {
			pjs[i] = presenters.JobSpec{JobSpec: j}
		}

		buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, pjs)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("failed to marshal document: %+v", err))
		} else {
			c.Data(200, MediaType, buffer)
		}
	}
}

// Create adds validates, saves, and starts a new JobSpec.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) Create(c *gin.Context) {

	var jr models.JobSpecRequest
	if err := c.ShouldBindJSON(&jr); err != nil {
		publicError(c, 400, err)
	} else if j, err := models.NewJobFromRequest(jr); err != nil {
		c.AbortWithError(500, err)
	} else if err := services.ValidateJob(j, jsc.App.Store); err != nil {
		publicError(c, 400, err)
	} else if err = jsc.App.AddJob(j); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, presenters.JobSpec{JobSpec: j})
	}
}

// Show returns the details of a JobSpec.
// Example:
//  "<application>/specs/:SpecID"
func (jsc *JobSpecsController) Show(c *gin.Context) {
	id := c.Param("SpecID")
	if j, err := jsc.App.Store.FindJob(id); err == storm.ErrNotFound {
		publicError(c, 404, errors.New("JobSpec not found"))
	} else if err != nil {
		c.AbortWithError(500, err)
	} else if runs, err := jsc.App.Store.JobRunsFor(j.ID); err != nil {
		c.AbortWithError(500, err)
	} else if doc, err := marshalSpecFromJSONAPI(j, runs); err != nil {
		c.AbortWithError(500, err)
	} else {
		c.JSON(200, doc)
	}
}

func marshalSpecFromJSONAPI(j models.JobSpec, runs []models.JobRun) (*jsonapi.Document, error) {
	pruns := make([]presenters.JobRun, len(runs))
	for i, r := range runs {
		pruns[i] = presenters.JobRun{r}
	}
	p := presenters.JobSpec{JobSpec: j, Runs: pruns}
	doc, err := jsonapi.MarshalToStruct(p, nil)
	return doc, err
}
