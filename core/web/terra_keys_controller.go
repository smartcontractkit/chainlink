package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// TerraKeysController manages Terra keys
type TerraKeysController struct {
	App chainlink.Application
}

// Index lists Terra keys
// Example:
// "GET <application>/keys/terra"
func (terkc *TerraKeysController) Index(c *gin.Context) {
	keys, err := terkc.App.GetKeyStore().Terra().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewTerraKeyResources(keys), "terraKey")
}

// Create and return a Terra key
// Example:
// "POST <application>/keys/terra"
func (terkc *TerraKeysController) Create(c *gin.Context) {
	key, err := terkc.App.GetKeyStore().Terra().Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewTerraKeyResource(key), "terraKey")
}

// Delete a Terra key
// Example:
// "DELETE <application>/keys/terra/:keyID"
// "DELETE <application>/keys/terra/:keyID?hard=true"
func (terkc *TerraKeysController) Delete(c *gin.Context) {
	keyID := c.Param("keyID")
	key, err := terkc.App.GetKeyStore().Terra().Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = terkc.App.GetKeyStore().Terra().Delete(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewTerraKeyResource(key), "terraKey")
}

// Import imports a Terra key
// Example:
// "Post <application>/keys/terra/import"
func (terkc *TerraKeysController) Import(c *gin.Context) {
	defer terkc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Import ")

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := terkc.App.GetKeyStore().Terra().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewTerraKeyResource(key), "terraKey")
}

// Export exports a Terra key
// Example:
// "Post <application>/keys/terra/export"
func (terkc *TerraKeysController) Export(c *gin.Context) {
	defer terkc.App.GetLogger().ErrorIfClosing(c.Request.Body, "Export request body")

	keyID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := terkc.App.GetKeyStore().Terra().Export(keyID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
