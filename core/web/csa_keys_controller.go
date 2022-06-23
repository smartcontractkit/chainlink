package web

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
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
	key, err := ctrl.App.GetKeyStore().CSA().Create()
	if err != nil {
		if errors.Is(err, keystore.ErrCSAKeyExists) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}

		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewCSAKeyResource(key), "csaKeys")
}

// Import imports a CSA key
func (ctrl *CSAKeysController) Import(c *gin.Context) {
	defer ctrl.App.GetLogger().ErrorIfClosing(c.Request.Body, "Import request body")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := ctrl.App.GetKeyStore().CSA().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewCSAKeyResource(key), "csaKey")
}

// Export exports a key
func (ctrl *CSAKeysController) Export(c *gin.Context) {
	defer ctrl.App.GetLogger().ErrorIfClosing(c.Request.Body, "Export request body")

	keyID := c.Param("ID")
	newPassword := c.Query("newpassword")

	bytes, err := ctrl.App.GetKeyStore().CSA().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, MediaType, bytes)
}
