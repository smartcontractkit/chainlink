package web

import (
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/static"

	"github.com/gin-gonic/gin"
)

// BuildVersonController has the build_info endpoint.
type BuildInfoController struct {
	App chainlink.Application
}

// Show returns the build info.
func (eic *BuildInfoController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": static.Version, "commitSHA": static.Sha})
}
