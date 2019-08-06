package web

import (
	"fmt"
	"math/big"
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

type PatchRequest struct {
	EthGasPriceDefault big.Int
}

// Patch updates one or more configuration options
func (cc *ConfigController) Patch(c *gin.Context) {
	request := &PatchRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if err := cc.App.GetStore().SetConfigValue("EthGasPriceDefault", &request.EthGasPriceDefault); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to set gas price default: %+v", err))
	} else if cw, err := presenters.NewConfigWhitelist(cc.App.GetStore()); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to build config response: %+v", err))
	} else {
		jsonAPIResponse(c, cw, "config")
	}
}
