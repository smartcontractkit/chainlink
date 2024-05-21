package web

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// CSAKeysController manages CSA keys
type CSAKeysController struct {
	App chainlink.Application
}

// Index lists CSA keys
// Example:
// "GET <application>/keys/csa"
func (ctrl *CSAKeysController) Index(c *gin.Context) {
	keys, err := ctrl.App.GetKeyStore().CSA().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewCSAKeyResources(keys), "csaKeys")
}

// Create and return a CSA key
// Example:
// "POST <application>/keys/csa"
func (ctrl *CSAKeysController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	key, err := ctrl.App.GetKeyStore().CSA().Create(ctx)
	if err != nil {
		if errors.Is(err, keystore.ErrCSAKeyExists) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.App.GetAuditLogger().Audit(audit.CSAKeyCreated, map[string]interface{}{
		"CSAPublicKey": key.PublicKey,
		"CSVersion":    key.Version,
	})

	jsonAPIResponse(c, presenters.NewCSAKeyResource(key), "csaKeys")
}

// Import imports a CSA key
func (ctrl *CSAKeysController) Import(c *gin.Context) {
	defer ctrl.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Import request body")
	ctx := c.Request.Context()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := ctrl.App.GetKeyStore().CSA().Import(ctx, bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.App.GetAuditLogger().Audit(audit.CSAKeyImported, map[string]interface{}{
		"CSAPublicKey": key.PublicKey,
		"CSVersion":    key.Version,
	})

	jsonAPIResponse(c, presenters.NewCSAKeyResource(key), "csaKey")
}

// Export exports a key
func (ctrl *CSAKeysController) Export(c *gin.Context) {
	defer ctrl.App.GetLogger().ErrorIfFn(c.Request.Body.Close, "Error closing Export request body")

	keyID := c.Param("ID")
	newPassword := c.Query("newpassword")

	bytes, err := ctrl.App.GetKeyStore().CSA().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	ctrl.App.GetAuditLogger().Audit(audit.CSAKeyExported, map[string]interface{}{"keyID": keyID})
	c.Data(http.StatusOK, MediaType, bytes)
}
