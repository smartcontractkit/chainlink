package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

// BuildInfoController has the build_info endpoint.
type BuildInfoController struct {
	App chainlink.Application
}

// Show returns the build info.
func (eic *BuildInfoController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": static.Version, "versionTag": static.VersionTag, "commitSHA": static.Sha})
}
