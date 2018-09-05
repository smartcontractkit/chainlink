package web

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
)

// MetricsController fetches high-level information about the nodes usage
type MetricsController struct {
	App *services.ChainlinkApplication
}

// Show displays metrics, one page at a time
// Example:
//  "<application>/metrics?size=1&page=2"
func (sc *MetricsController) Show(c *gin.Context) {
	jobs := []models.JobSpec{}

	size, page, offset, err := ParsePaginatedRequest(c.Query("size"), c.Query("page"))
	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	skip := storm.Skip(offset)
	limit := storm.Limit(size)

	store := sc.App.Store
	if count, err := store.Count(&models.JobSpec{}); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting count of JobSpec: %+v", err))
	} else if err := store.AllByIndex("CreatedAt", &jobs, skip, limit); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting Jobs: %+v", err))
	} else if jsm, err := services.AllJobSpecMetrics(sc.App.Store, jobs); err != nil {
		c.AbortWithError(500, fmt.Errorf("error getting job stats: %+v", err))
	} else if account, err := store.KeyStore.GetAccount(); err != nil {
		publicError(c, 400, err)
	} else {
		jsmc := presenters.JobSpecMetricsCollection{
			Address:        account.Address.Hex(),
			JobSpecMetrics: jsm,
		}
		buffer, err := NewPaginatedResponse(*c.Request.URL, size, page, count, jsmc)
		if err != nil {
			c.AbortWithError(500, fmt.Errorf("failed to marshal document: %+v", err))
		} else {
			c.Data(200, MediaType, buffer)
		}
	}
}
