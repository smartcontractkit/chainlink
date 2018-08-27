package web

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web/metric"
	"net/http"
	"strings"
)

// MetricsController fetches high-level information about the nodes usage
type MetricsController struct {
	App *services.ChainlinkApplication
}

// Show displays metrics in a specific format switched on the User-Agent header
// Example:
//  "<application>/metrics"
func (sc *MetricsController) Show(c *gin.Context) {
	if utils.StripBearer(c.Request.Header.Get("Authorization")) != sc.App.Store.Config.MetricsBearerToken {
		publicError(c, http.StatusUnauthorized, fmt.Errorf("incorrect access token for metrics"))
	} else {
		f := false
		for _, mp := range metric.Controllers {
			if strings.Contains(c.GetHeader("User-Agent"), mp.UserAgent()) {
				f = true
				if j, err := sc.App.Store.Jobs(); err != nil {
					c.AbortWithError(500, fmt.Errorf("error getting all jobs: %+v", err))
				} else if jsm, err := services.AllJobSpecMetrics(sc.App.Store, j); err != nil {
					c.AbortWithError(500, fmt.Errorf("error job spec stats: %+v", err))
				} else {
					mp.Show(&jsm, c.Writer, c.Request)
				}
				break
			}
		}
		if !f {
			sc.ShowRaw(c)
		}
	}
}

// ShowRaw displays metrics, one page at a time in the raw JSON format
// Example:
//  "<application>/metrics?size=1&page=2"
func (sc *MetricsController) ShowRaw(c *gin.Context) {
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
	} else if s, err := services.AllJobSpecMetrics(sc.App.Store, jobs); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting job stats: %+v", err))
	} else {
		buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, s)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("failed to marshal document: %+v", err))
		} else {
			c.Data(200, MediaType, buffer)
		}
	}
}
