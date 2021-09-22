package web

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// P2PKeysController manages P2P keys
type P2PKeysController struct {
	App chainlink.Application
}

// Index lists P2P keys
// Example:
// "GET <application>/keys/p2p"
func (p2pkc *P2PKeysController) Index(c *gin.Context) {
	keys, err := p2pkc.App.GetKeyStore().P2P().GetAll()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResources(keys), "p2pKey")
}

// Create and return a P2P key
// Example:
// "POST <application>/keys/p2p"
func (p2pkc *P2PKeysController) Create(c *gin.Context) {
	key, err := p2pkc.App.GetKeyStore().P2P().Create()
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Delete a P2P key
// Example:
// "DELETE <application>/keys/p2p/:keyID"
// "DELETE <application>/keys/p2p/:keyID?hard=true"
func (p2pkc *P2PKeysController) Delete(c *gin.Context) {
	keyID := c.Param("keyID")
	key, err := p2pkc.App.GetKeyStore().P2P().Get(keyID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	_, err = p2pkc.App.GetKeyStore().P2P().Delete(key.ID())
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Import imports a P2P key
// Example:
// "Post <application>/keys/p2p/import"
func (p2pkc *P2PKeysController) Import(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	oldPassword := c.Query("oldpassword")
	key, err := p2pkc.App.GetKeyStore().P2P().Import(bytes, oldPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	jsonAPIResponse(c, presenters.NewP2PKeyResource(key), "p2pKey")
}

// Export exports a P2P key
// Example:
// "Post <application>/keys/p2p/export"
func (p2pkc *P2PKeysController) Export(c *gin.Context) {
	defer logger.ErrorIfCalling(c.Request.Body.Close)

	stringID := c.Param("ID")
	newPassword := c.Query("newpassword")
	bytes, err := p2pkc.App.GetKeyStore().P2P().Export(stringID, newPassword)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, MediaType, bytes)
}
