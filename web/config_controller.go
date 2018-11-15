package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/presenters"
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
		c.AbortWithError(500, fmt.Errorf("failed to build config whitelist: %+v", err))
	} else if json, err := jsonapi.Marshal(cw); err != nil {
		c.AbortWithError(500, fmt.Errorf("failed to marshal config using jsonapi: %+v", err))
	} else {
		c.Data(200, MediaType, json)
	}
}
