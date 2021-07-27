package web

import (
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/gin-gonic/gin"
)

// ConfigController manages config variables
type ConfigController struct {
	App chainlink.Application
}

// Show returns the whitelist of config variables
// Example:
//  "<application>/config"
func (cc *ConfigController) Show(c *gin.Context) {
	cw, err := presenters.NewConfigPrinter(cc.App.GetStore())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to build config whitelist: %+v", err))
		return
	}

	jsonAPIResponse(c, cw, "config")
}

type configPatchRequest struct {
	EthGasPriceDefault *utils.Big `json:"ethGasPriceDefault"`
}

// ConfigPatchResponse represents the change to the configuration made due to a
// PATCH to the config endpoint
type ConfigPatchResponse struct {
	EthGasPriceDefault Change `json:"ethGasPriceDefault"`
}

// Change represents the old value and the new value after a PATH request has
// been made
type Change struct {
	From string `json:"old"`
	To   string `json:"new"`
}

// GetID returns the jsonapi ID.
func (c ConfigPatchResponse) GetID() string {
	return "configuration"
}

// SetID is used to conform to the UnmarshallIdentifier interface for
// deserializing from jsonapi documents.
func (*ConfigPatchResponse) SetID(string) error {
	return nil
}

// Patch updates one or more configuration options
func (cc *ConfigController) Patch(c *gin.Context) {
	request := &configPatchRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	if err := cc.App.GetConfig().ORM.SetConfigValue("EthGasPriceDefault", request.EthGasPriceDefault); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to set gas price default: %+v", err))
		return
	}

	response := &ConfigPatchResponse{
		EthGasPriceDefault: Change{
			From: cc.App.GetStore().Config.EthGasPriceDefault().String(),
			To:   request.EthGasPriceDefault.String(),
		},
	}
	jsonAPIResponse(c, response, "config")
}
