package web

import (
	"fmt"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

const (
	// PaginationDefault is the number of records to supply from a paginated
	// request when no size param is supplied.
	PaginationDefault = 20
)

// JobSpecsController manages JobSpec requests.
type JobSpecsController struct {
	App *services.ChainlinkApplication
}

// ParsePaginatedRequest parses the parameters that control pagination for a
// collection request, returning the size and offset if specified, or a
// sensible default.
func ParsePaginatedRequest(sizeParam, offsetParam string) (int, int, error) {
	var err error
	var offset int
	size := PaginationDefault

	if sizeParam != "" {
		if size, err = strconv.Atoi(sizeParam); err != nil {
			return 0, 0, fmt.Errorf("invalid size param, error: %+v", err)
		}
	}

	if offsetParam != "" {
		if offset, err = strconv.Atoi(offsetParam); err != nil {
			return 0, 0, fmt.Errorf("invalid offset param, error: %+v", err)
		}
	}

	return size, offset, nil
}

// Index lists JobSpecs, one page at a time.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) Index(c *gin.Context) {
	sizeParam := c.Query("size")
	offsetParam := c.Query("offset")
	size, offset, err := ParsePaginatedRequest(sizeParam, offsetParam)
	if err != nil {
		c.JSON(422, gin.H{
			"errors": []string{err.Error()},
		})
	}

	var jobs []models.JobSpec
	if count, err := jsc.App.Store.Count(&models.JobSpec{}); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("error getting count of JobSpec: %+v", err).Error()},
		})
	} else if err := jsc.App.Store.AllByIndex("CreatedAt", &jobs, storm.Limit(size), storm.Skip(offset)); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{fmt.Errorf("erorr fetching All JobSpecs: %+v", err).Error()},
		})
	} else {
		pjs := make([]presenters.JobSpec, len(jobs))
		for i, j := range jobs {
			pjs[i] = presenters.JobSpec{JobSpec: j}
		}
		c.JSON(200, presenters.NewHALResponse(c.Request.URL.Path, size, offset, count, pjs))
	}
}

// Index lists *ALL* of the existing JobSpecs.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) IndexV2(c *gin.Context) {
	var jobs []models.JobSpec
	if err := jsc.App.Store.AllByIndex("CreatedAt", &jobs); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		pjs := make([]presenters.JobSpec, len(jobs))
		for i, j := range jobs {
			pjs[i] = presenters.JobSpec{JobSpec: j}
		}
		c.JSON(200, pjs)
	}
}

// Create adds validates, saves, and starts a new JobSpec.
// Example:
//  "<application>/specs"
func (jsc *JobSpecsController) Create(c *gin.Context) {
	j := models.NewJob()

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = services.ValidateJob(j, jsc.App.Store); err != nil {
		c.JSON(400, gin.H{
			"errors": []string{err.Error()},
		})
	} else if err = jsc.App.AddJob(j); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
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
		c.JSON(404, gin.H{
			"errors": []string{"JobSpec not found."},
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else if runs, err := jsc.App.Store.JobRunsFor(j.ID); err != nil {
		c.JSON(500, gin.H{
			"errors": []string{err.Error()},
		})
	} else {
		c.JSON(200, presenters.JobSpec{JobSpec: j, Runs: runs})
	}
}
