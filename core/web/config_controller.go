package web

import (
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
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
	cw := config.NewConfigPrinter(cc.App.GetConfig())

	jsonAPIResponse(c, cw, "config")
}

type configPatchRequest struct {
	EvmGasPriceDefault *utils.Big `json:"ethGasPriceDefault"`
	EVMChainID         *utils.Big `json:"evmChainID"`
}

// ConfigPatchResponse represents the change to the configuration made due to a
// PATCH to the config endpoint
type ConfigPatchResponse struct {
	EvmGasPriceDefault Change     `json:"ethGasPriceDefault"`
	EVMChainID         *utils.Big `json:"evmChainID"`
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

	chain, err := getChain(cc.App.GetChains().EVM, request.EVMChainID.String())
	switch err {
	case ErrInvalidChainID, ErrMultipleChains, ErrMissingChainID:
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	case nil:
		break
	default:
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if err := chain.Config().SetEvmGasPriceDefault(request.EvmGasPriceDefault.ToInt()); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, fmt.Errorf("failed to set gas price default: %+v", err))
		return
	}
	response := &ConfigPatchResponse{
		EvmGasPriceDefault: Change{
			From: chain.Config().EvmGasPriceDefault().String(),
			To:   request.EvmGasPriceDefault.String(),
		}, EVMChainID: utils.NewBig(chain.ID()),
	}
	jsonAPIResponse(c, response, "config")
}
