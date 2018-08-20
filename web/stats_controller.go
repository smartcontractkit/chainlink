package web

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
)

// StatsController fetches high-level information about the nodes usage
type StatsController struct {
	App *services.ChainlinkApplication
}

// Show displays Stats, one page at a time
// Example:
//  "<application>/stats?size=1&page=2"
func (sc *StatsController) Show(c *gin.Context) {
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
	} else if s, err := services.AllJobSpecStats(sc.App.Store, jobs); err != nil {
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
