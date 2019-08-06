package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
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

type configPatchRequest struct {
	EthGasPriceDefault *models.Big `json:"ethGasPriceDefault"`
}

// ConfigPatchResponse represents the change to the configuration made due to a
// PATCH to the config endpoint
type ConfigPatchResponse struct {
	OldEthGasPriceDefault *models.Big `json:"oldEthGasPriceDefault"`
	NewEthGasPriceDefault *models.Big `json:"newEthGasPriceDefault"`
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
	response := &ConfigPatchResponse{
		OldEthGasPriceDefault: models.NewBig(cc.App.GetStore().Config.EthGasPriceDefault()),
	}

	if err := c.ShouldBindJSON(request); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
	} else if err := cc.App.GetStore().SetConfigValue("EthGasPriceDefault", request.EthGasPriceDefault); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to set gas price default: %+v", err))
	} else {
		response.NewEthGasPriceDefault = request.EthGasPriceDefault
		jsonAPIResponse(c, response, "config")
	}
}
