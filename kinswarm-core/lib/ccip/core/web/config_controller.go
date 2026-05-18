package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/gin-gonic/gin"
)

// ConfigController manages config variables
type ConfigController struct {
	App chainlink.Application
}

// Show returns the whitelist of config variables
// Example:
//
//	"<application>/config"
func (cc *ConfigController) Show(c *gin.Context) {
	cfg := cc.App.GetConfig()
	var userOnly bool
	if s, has := c.GetQuery("userOnly"); has {
		var err error
		userOnly, err = strconv.ParseBool(s)
		if err != nil {
			jsonAPIError(c, http.StatusBadRequest, fmt.Errorf("invalid bool for userOnly: %v", err))
			return
		}
	}
	var toml string
	user, effective := cfg.ConfigTOML()
	if userOnly {
		toml = user
	} else {
		toml = effective
	}
	jsonAPIResponse(c, ConfigV2Resource{toml}, "config")
}

type ConfigV2Resource struct {
	Config string `json:"config"`
}

func (c ConfigV2Resource) GetID() string {
	return utils.NewBytes32ID()
}

func (c *ConfigV2Resource) SetID(string) error {
	return nil
}
