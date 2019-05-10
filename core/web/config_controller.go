package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
)

// ConfigController manages config variables
type ConfigController struct {
	App services.Application
}

// Show returns the whitelist of config variables
// Example:
//  "<application>/config"
func (cc *ConfigController) Show(c *gin.Context) {
	cw, err := presenters.NewConfigWhitelist(cc.App.GetStore())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to build config whitelist: %+v", err))
	} else {
		jsonAPIResponse(c, cw, "config")
	}
}
