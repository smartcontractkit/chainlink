package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/config"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
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
//
//	"<application>/config"
func (cc *ConfigController) Show(c *gin.Context) {
	cfg := cc.App.GetConfig()
	if _, ok := cfg.(chainlink.ConfigV2); ok {
		jsonAPIError(c, http.StatusUnprocessableEntity, v2.ErrUnsupported)
		return
	}
	// Legacy config
	cw := config.NewConfigPrinter(cc.App.GetConfig())

	cc.App.GetAuditLogger().Audit(audit.EnvNoncriticalEnvDumped, map[string]interface{}{})
	jsonAPIResponse(c, cw, "config")
}

func (cc *ConfigController) Show2(c *gin.Context) {
	cfg := cc.App.GetConfig()
	cfg2, ok := cfg.(chainlink.ConfigV2)
	if !ok {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.New("unsupported with legacy ENV config"))
		return
	}
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
	user, effective := cfg2.ConfigTOML()
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

func (cc *ConfigController) Dump(c *gin.Context) {
	cfg := cc.App.GetConfig()
	if _, ok := cfg.(chainlink.ConfigV2); ok {
		jsonAPIError(c, http.StatusUnprocessableEntity, v2.ErrUnsupported)
		return
	}
	// Legacy config mode
	userToml, err := cc.App.ConfigDump(c)
	if err != nil {
		cc.App.GetLogger().Errorw("Failed to dump TOML config", "err", err)
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, ConfigV2Resource{userToml}, "config")
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

	cc.App.GetAuditLogger().Audit(audit.ConfigUpdated, map[string]interface{}{"configResponse": response})
	jsonAPIResponse(c, response, "config")
}
