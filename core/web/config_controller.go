package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/orm"
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
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to build config response: %+v", err))
	} else {
		jsonAPIResponse(c, cw, "config")
	}
}

// Patch updates one or more configuration options
func (cc *ConfigController) Patch(c *gin.Context) {
	request := &orm.Request{}

	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if err := cc.App.GetStore().SetConfigValue(orm.EnvVarName("EthGasPriceDefault"), request.EthGasPriceDefault.String()); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to set gas price default: %+v", err))
	} else if cw, err := presenters.NewConfigWhitelist(cc.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to build config response: %+v", err))
	} else {
		jsonAPIResponse(c, cw, "config")
	}
}
