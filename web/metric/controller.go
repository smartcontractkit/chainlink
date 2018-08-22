package metric

import (
	"github.com/smartcontractkit/chainlink/store/models"
	"net/http"
)

// Controllers holds all the supported proprietary metric controllers
var Controllers []Controller

// Set the supported proprietary metric controllers
func init() {
	Controllers = []Controller{
		&PromController{},
	}
}

// Controller is a interface for showing metrics in proprietary formats
type Controller interface {
	Show(jss *models.JobSpecMetrics, w http.ResponseWriter, r *http.Request)
	UserAgent() string
}
